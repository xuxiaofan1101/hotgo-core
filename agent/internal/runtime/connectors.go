package runtime

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/tls"
	"encoding/hex"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
	"time"

	"github.com/IBM/sarama"
	xdgscram "github.com/xdg-go/scram"

	"vogo-agent/internal/protocol"
)

func registerDefaultConnectors(registry *ConnectorRegistry) {
	registry.RegisterSource("kafka", newKafkaSource)
	registry.RegisterSource("http", newHTTPSource)
	registry.RegisterSource("s3", newS3Source)
	registry.RegisterSource("log", newLogSource)
	registry.RegisterSource("file", newLogSource)

	registry.RegisterSink("kafka", newKafkaSink)
	registry.RegisterSink("http", newHTTPSink)
	registry.RegisterSink("webhook", newHTTPSink)
	registry.RegisterSink("s3", newS3Sink)
	registry.RegisterSink("elasticsearch", newSearchSink)
	registry.RegisterSink("opensearch", newSearchSink)
	registry.RegisterSink("splunk", newSplunkSink)
	registry.RegisterSink("feishu", newFeishuSink)
	registry.RegisterSink("lark", newFeishuSink)
	registry.RegisterSink("strategy", func(config protocol.OutputConfig) (Sink, error) {
		return newStrategySink(config, registry)
	})
	registry.RegisterSink("log", newLogSink)
}

type logSource struct {
	filePath string
	follow   bool
	interval time.Duration
	offset   int64
	maxBytes int64
}

func newLogSource(config protocol.IntegrationConfig, shard protocol.TaskShard) (Source, error) {
	values := mergeConfig(config.Config, shard.Config)
	filePath := configString(values, "path", "file", "filePath")
	if filePath == "" {
		return nil, fmt.Errorf("日志输入需要配置 path")
	}
	return &logSource{
		filePath: filePath,
		follow:   configBool(values, "follow"),
		interval: configDuration(values, "pollIntervalSeconds", 3*time.Second),
		offset:   configInt64(values, "offset", 0),
		maxBytes: configInt64(values, "maxBytes", 10*1024*1024),
	}, nil
}

func (s *logSource) Run(ctx context.Context, emit func(Record) error) error {
	for {
		if err := s.readOnce(ctx, emit); err != nil {
			return err
		}
		if !s.follow {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.interval):
		}
	}
}

func (s *logSource) readOnce(ctx context.Context, emit func(Record) error) error {
	file, err := os.Open(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		return err
	}
	if s.offset > stat.Size() {
		s.offset = 0
	}
	if _, err := file.Seek(s.offset, io.SeekStart); err != nil {
		return err
	}
	limit := s.maxBytes
	if limit <= 0 {
		limit = 10 * 1024 * 1024
	}
	data, err := io.ReadAll(io.LimitReader(file, limit))
	if err != nil {
		return err
	}
	if len(data) == 0 {
		return nil
	}
	s.offset += int64(len(data))
	return emit(Record{Raw: data, Metadata: map[string]any{"path": s.filePath, "offset": s.offset}})
}

type httpSource struct {
	client   *http.Client
	url      string
	method   string
	headers  map[string]string
	body     []byte
	interval time.Duration
	repeat   bool
}

func newHTTPSource(config protocol.IntegrationConfig, shard protocol.TaskShard) (Source, error) {
	values := mergeConfig(config.Config, shard.Config)
	endpoint := configString(values, "url", "endpoint")
	if endpoint == "" {
		return nil, fmt.Errorf("HTTP 输入需要配置 url")
	}
	repeat := configBool(values, "repeat") || configBool(values, "poll")
	return &httpSource{
		client:   &http.Client{Timeout: configDuration(values, "timeoutSeconds", 15*time.Second)},
		url:      endpoint,
		method:   strings.ToUpper(defaultString(configString(values, "method"), http.MethodGet)),
		headers:  configStringMap(values, "headers"),
		body:     []byte(configString(values, "body")),
		interval: configDuration(values, "pollIntervalSeconds", 5*time.Second),
		repeat:   repeat,
	}, nil
}

func (s *httpSource) Run(ctx context.Context, emit func(Record) error) error {
	for {
		if err := s.readOnce(ctx, emit); err != nil {
			return err
		}
		if !s.repeat {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(s.interval):
		}
	}
}

func (s *httpSource) readOnce(ctx context.Context, emit func(Record) error) error {
	req, err := http.NewRequestWithContext(ctx, s.method, s.url, bytes.NewReader(s.body))
	if err != nil {
		return err
	}
	for key, value := range s.headers {
		req.Header.Set(key, value)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		data, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("HTTP 输入请求失败: status=%d body=%s", resp.StatusCode, string(data))
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return emit(Record{Raw: data, Metadata: map[string]any{"url": s.url, "status": resp.StatusCode}})
}

type httpSink struct {
	client  *http.Client
	url     string
	method  string
	headers map[string]string
}

func newHTTPSink(config protocol.OutputConfig) (Sink, error) {
	endpoint := configString(config.Config, "url", "webhook", "endpoint")
	if endpoint == "" {
		return nil, fmt.Errorf("HTTP 输出需要配置 url")
	}
	return &httpSink{
		client:  &http.Client{Timeout: configDuration(config.Config, "timeoutSeconds", 10*time.Second)},
		url:     endpoint,
		method:  strings.ToUpper(defaultString(configString(config.Config, "method"), http.MethodPost)),
		headers: configStringMap(config.Config, "headers"),
	}, nil
}

func (s *httpSink) Write(ctx context.Context, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, s.method, s.url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range s.headers {
		req.Header.Set(key, value)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("HTTP 输出失败: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

type logSink struct {
	filePath string
}

func newLogSink(config protocol.OutputConfig) (Sink, error) {
	return &logSink{filePath: configString(config.Config, "path", "file", "filePath")}, nil
}

func (s *logSink) Write(ctx context.Context, payload map[string]any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if s.filePath == "" || s.filePath == "-" {
		_, err = os.Stdout.Write(data)
		return err
	}
	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write(data)
	return err
}

type searchSink struct {
	client  *http.Client
	baseURL string
	index   string
	headers map[string]string
}

func newSearchSink(config protocol.OutputConfig) (Sink, error) {
	baseURL := strings.TrimRight(configString(config.Config, "baseUrl", "url", "endpoint"), "/")
	index := strings.Trim(configString(config.Config, "index", "indexName"), "/")
	if baseURL == "" || index == "" {
		return nil, fmt.Errorf("%s 输出需要配置 baseUrl 和 index", strings.ToUpper(config.Type))
	}
	return &searchSink{
		client:  &http.Client{Timeout: configDuration(config.Config, "timeoutSeconds", 10*time.Second)},
		baseURL: baseURL,
		index:   index,
		headers: configStringMap(config.Config, "headers"),
	}, nil
}

func (s *searchSink) Write(ctx context.Context, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.baseURL+"/"+s.index+"/_doc", bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for key, value := range s.headers {
		req.Header.Set(key, value)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("搜索输出失败: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

type splunkSink struct {
	client  *http.Client
	hecURL  string
	token   string
	index   string
	source  string
	headers map[string]string
}

type feishuSink struct {
	client   *http.Client
	webhook  string
	template string
}

func newFeishuSink(config protocol.OutputConfig) (Sink, error) {
	webhook := configString(config.Config, "webhook", "url")
	if webhook == "" {
		return nil, fmt.Errorf("飞书输出需要配置 webhook")
	}
	template := configString(config.Config, "messageTemplate")
	if template == "" {
		template = "Vogo 数据清洗告警：{{eventType}}"
	}
	return &feishuSink{
		client:   &http.Client{Timeout: configDuration(config.Config, "timeoutSeconds", 10*time.Second)},
		webhook:  webhook,
		template: template,
	}, nil
}

func (s *feishuSink) Write(ctx context.Context, payload map[string]any) error {
	message := s.template
	for key, value := range payload {
		message = strings.ReplaceAll(message, "{{"+key+"}}", fmt.Sprint(value))
	}
	body := map[string]any{
		"msg_type": "text",
		"content":  map[string]any{"text": message},
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.webhook, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("飞书输出失败: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

func newSplunkSink(config protocol.OutputConfig) (Sink, error) {
	hecURL := configString(config.Config, "hecUrl", "url", "endpoint")
	token := configString(config.Config, "token", "hecToken")
	if hecURL == "" || token == "" {
		return nil, fmt.Errorf("Splunk 输出需要配置 hecUrl 和 token")
	}
	return &splunkSink{
		client:  &http.Client{Timeout: configDuration(config.Config, "timeoutSeconds", 10*time.Second)},
		hecURL:  hecURL,
		token:   token,
		index:   configString(config.Config, "index"),
		source:  configString(config.Config, "source"),
		headers: configStringMap(config.Config, "headers"),
	}, nil
}

func (s *splunkSink) Write(ctx context.Context, payload map[string]any) error {
	body := map[string]any{"event": payload}
	if s.index != "" {
		body["index"] = s.index
	}
	if s.source != "" {
		body["source"] = s.source
	}
	data, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.hecURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Splunk "+s.token)
	req.Header.Set("Content-Type", "application/json")
	for key, value := range s.headers {
		req.Header.Set(key, value)
	}
	resp, err := s.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("Splunk 输出失败: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

type s3Source struct {
	client *s3Client
	bucket string
	key    string
	prefix string
	maxKey int
}

func newS3Source(config protocol.IntegrationConfig, shard protocol.TaskShard) (Source, error) {
	values := mergeConfig(config.Config, shard.Config)
	client, err := newS3Client(values)
	if err != nil {
		return nil, err
	}
	bucket := configString(values, "bucket")
	if bucket == "" {
		return nil, fmt.Errorf("S3 输入需要配置 bucket")
	}
	return &s3Source{
		client: client,
		bucket: bucket,
		key:    configString(values, "key", "objectKey"),
		prefix: configString(values, "prefix"),
		maxKey: configInt(values, "maxKeys", 100),
	}, nil
}

func (s *s3Source) Run(ctx context.Context, emit func(Record) error) error {
	keys := []string{s.key}
	var err error
	if s.key == "" {
		keys, err = s.client.listObjects(ctx, s.bucket, s.prefix, s.maxKey)
		if err != nil {
			return err
		}
	}
	for _, key := range keys {
		key = strings.TrimSpace(key)
		if key == "" {
			continue
		}
		data, err := s.client.getObject(ctx, s.bucket, key)
		if err != nil {
			return err
		}
		if err := emit(Record{Raw: data, Metadata: map[string]any{"bucket": s.bucket, "key": key}}); err != nil {
			return err
		}
	}
	return nil
}

type s3Sink struct {
	client *s3Client
	bucket string
	key    string
	prefix string
}

func newS3Sink(config protocol.OutputConfig) (Sink, error) {
	client, err := newS3Client(config.Config)
	if err != nil {
		return nil, err
	}
	bucket := configString(config.Config, "bucket")
	if bucket == "" {
		return nil, fmt.Errorf("S3 输出需要配置 bucket")
	}
	return &s3Sink{
		client: client,
		bucket: bucket,
		key:    configString(config.Config, "key", "objectKey"),
		prefix: configString(config.Config, "prefix"),
	}, nil
}

func (s *s3Sink) Write(ctx context.Context, payload map[string]any) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	key := s.key
	if key == "" {
		key = strings.Trim(strings.TrimSpace(s.prefix), "/")
		if key != "" {
			key += "/"
		}
		key += fmt.Sprintf("agent-output-%d.jsonl", time.Now().UnixNano())
	}
	return s.client.putObject(ctx, s.bucket, key, data)
}

type s3Client struct {
	httpClient   *http.Client
	endpoint     string
	region       string
	accessKey    string
	secretKey    string
	sessionToken string
	pathStyle    bool
}

func newS3Client(config map[string]any) (*s3Client, error) {
	region := defaultString(configString(config, "region"), "us-east-1")
	endpoint := strings.TrimRight(configString(config, "endpoint", "baseUrl"), "/")
	accessKey := configString(config, "accessKey", "accessKeyId")
	secretKey := configString(config, "secretKey", "secretAccessKey")
	sessionToken := configString(config, "sessionToken", "token")
	if accessKey == "" {
		accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	}
	if secretKey == "" {
		secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	}
	if sessionToken == "" {
		sessionToken = os.Getenv("AWS_SESSION_TOKEN")
	}
	return &s3Client{
		httpClient:   &http.Client{Timeout: configDuration(config, "timeoutSeconds", 30*time.Second)},
		endpoint:     endpoint,
		region:       region,
		accessKey:    accessKey,
		secretKey:    secretKey,
		sessionToken: sessionToken,
		pathStyle:    configBool(config, "pathStyle") || configBool(config, "forcePathStyle") || endpoint != "",
	}, nil
}

type s3ListResult struct {
	Contents []struct {
		Key string `xml:"Key"`
	} `xml:"Contents"`
}

func (c *s3Client) listObjects(ctx context.Context, bucket, prefix string, maxKeys int) ([]string, error) {
	query := url.Values{}
	query.Set("list-type", "2")
	if prefix != "" {
		query.Set("prefix", prefix)
	}
	if maxKeys > 0 {
		query.Set("max-keys", fmt.Sprint(maxKeys))
	}
	req, err := c.newS3Request(ctx, http.MethodGet, bucket, "", query, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("S3 列表读取失败: status=%d body=%s", resp.StatusCode, string(body))
	}
	var result s3ListResult
	if err := xml.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(result.Contents))
	for _, item := range result.Contents {
		keys = append(keys, item.Key)
	}
	return keys, nil
}

func (c *s3Client) getObject(ctx context.Context, bucket, objectKey string) ([]byte, error) {
	req, err := c.newS3Request(ctx, http.MethodGet, bucket, objectKey, nil, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return nil, fmt.Errorf("S3 对象读取失败: status=%d body=%s", resp.StatusCode, string(body))
	}
	return io.ReadAll(resp.Body)
}

func (c *s3Client) putObject(ctx context.Context, bucket, objectKey string, body []byte) error {
	req, err := c.newS3Request(ctx, http.MethodPut, bucket, objectKey, nil, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 4096))
		return fmt.Errorf("S3 对象写入失败: status=%d body=%s", resp.StatusCode, string(body))
	}
	return nil
}

func (c *s3Client) newS3Request(ctx context.Context, method, bucket, objectKey string, query url.Values, body io.Reader) (*http.Request, error) {
	endpoint := c.endpoint
	if endpoint == "" {
		endpoint = "https://s3." + c.region + ".amazonaws.com"
	}
	baseURL, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	if c.pathStyle {
		baseURL.Path = path.Join(baseURL.Path, bucket, objectKey)
	} else {
		baseURL.Host = bucket + "." + baseURL.Host
		baseURL.Path = path.Join(baseURL.Path, objectKey)
	}
	if query != nil {
		baseURL.RawQuery = query.Encode()
	}
	req, err := http.NewRequestWithContext(ctx, method, baseURL.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Host", req.URL.Host)
	if c.accessKey != "" && c.secretKey != "" {
		c.sign(req)
	}
	return req, nil
}

func (c *s3Client) sign(req *http.Request) {
	now := time.Now().UTC()
	date := now.Format("20060102")
	amzDate := now.Format("20060102T150405Z")
	req.Header.Set("X-Amz-Date", amzDate)
	if c.sessionToken != "" {
		req.Header.Set("X-Amz-Security-Token", c.sessionToken)
	}
	payloadHash := "UNSIGNED-PAYLOAD"
	req.Header.Set("X-Amz-Content-Sha256", payloadHash)
	signedHeaders := signedHeaderNames(req.Header)
	canonicalRequest := strings.Join([]string{
		req.Method,
		s3CanonicalURI(req.URL.EscapedPath()),
		req.URL.Query().Encode(),
		canonicalHeaders(req.Header),
		strings.Join(signedHeaders, ";"),
		payloadHash,
	}, "\n")
	scope := date + "/" + c.region + "/s3/aws4_request"
	stringToSign := strings.Join([]string{
		"AWS4-HMAC-SHA256",
		amzDate,
		scope,
		sha256Hex([]byte(canonicalRequest)),
	}, "\n")
	signingKey := awsSigningKey(c.secretKey, date, c.region, "s3")
	signature := hex.EncodeToString(hmacSHA256(signingKey, stringToSign))
	req.Header.Set("Authorization", fmt.Sprintf(
		"AWS4-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		c.accessKey,
		scope,
		strings.Join(signedHeaders, ";"),
		signature,
	))
}

func signedHeaderNames(headers http.Header) []string {
	names := make([]string, 0, len(headers))
	for name := range headers {
		names = append(names, strings.ToLower(name))
	}
	sort.Strings(names)
	return names
}

func canonicalHeaders(headers http.Header) string {
	names := signedHeaderNames(headers)
	var builder strings.Builder
	for _, name := range names {
		builder.WriteString(name)
		builder.WriteByte(':')
		builder.WriteString(strings.Join(headers.Values(name), ","))
		builder.WriteByte('\n')
	}
	return builder.String()
}

func s3CanonicalURI(value string) string {
	if value == "" {
		return "/"
	}
	return value
}

func awsSigningKey(secret, date, region, service string) []byte {
	kDate := hmacSHA256([]byte("AWS4"+secret), date)
	kRegion := hmacSHA256(kDate, region)
	kService := hmacSHA256(kRegion, service)
	return hmacSHA256(kService, "aws4_request")
}

func hmacSHA256(key []byte, value string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(value))
	return mac.Sum(nil)
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

type kafkaSource struct {
	brokers []string
	topic   string
	groupId string
	config  *sarama.Config
	offsets map[int32]int64
}

func newKafkaSource(config protocol.IntegrationConfig, shard protocol.TaskShard) (Source, error) {
	values := mergeConfig(config.Config, shard.Config)
	brokers := configStringSlice(values, "brokers")
	topic := configString(values, "topic")
	if len(brokers) == 0 || topic == "" {
		return nil, fmt.Errorf("Kafka 输入需要配置 brokers 和 topic")
	}
	kafkaConfig, err := buildKafkaConfig(values)
	if err != nil {
		return nil, err
	}
	groupId := configString(values, "groupId", "consumerGroup")
	if groupId == "" {
		groupId = fmt.Sprintf("vogo-agent-%s-%d", configString(values, "_tenantId"), configInt64(values, "_taskId", time.Now().UnixNano()))
	}
	return &kafkaSource{
		brokers: brokers,
		topic:   topic,
		groupId: groupId,
		config:  kafkaConfig,
		offsets: kafkaSpecificOffsets(values),
	}, nil
}

func (s *kafkaSource) Run(ctx context.Context, emit func(Record) error) error {
	group, err := sarama.NewConsumerGroup(s.brokers, s.groupId, s.config)
	if err != nil {
		return err
	}
	defer group.Close()
	handler := &kafkaConsumerHandler{emit: emit, offsets: s.offsets}
	for ctx.Err() == nil {
		if err := group.Consume(ctx, []string{s.topic}, handler); err != nil {
			return err
		}
	}
	return ctx.Err()
}

type kafkaConsumerHandler struct {
	emit    func(Record) error
	offsets map[int32]int64
}

func (h *kafkaConsumerHandler) Setup(session sarama.ConsumerGroupSession) error {
	for topic, partitions := range session.Claims() {
		for _, partition := range partitions {
			if offset, ok := h.offsets[partition]; ok {
				session.ResetOffset(topic, partition, offset, "")
			}
		}
	}
	return nil
}

func (h *kafkaConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *kafkaConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case <-session.Context().Done():
			return nil
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			if err := h.emit(Record{
				Raw: msg.Value,
				Metadata: map[string]any{
					"topic":     msg.Topic,
					"partition": msg.Partition,
					"offset":    msg.Offset,
				},
			}); err != nil {
				return err
			}
			session.MarkMessage(msg, "")
		}
	}
}

type kafkaSink struct {
	producer sarama.SyncProducer
	topic    string
	key      string
}

func newKafkaSink(config protocol.OutputConfig) (Sink, error) {
	brokers := configStringSlice(config.Config, "brokers")
	topic := configString(config.Config, "topic")
	if len(brokers) == 0 || topic == "" {
		return nil, fmt.Errorf("Kafka 输出需要配置 brokers 和 topic")
	}
	kafkaConfig, err := buildKafkaConfig(config.Config)
	if err != nil {
		return nil, err
	}
	kafkaConfig.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, kafkaConfig)
	if err != nil {
		return nil, err
	}
	return &kafkaSink{producer: producer, topic: topic, key: configString(config.Config, "key")}, nil
}

func (s *kafkaSink) Write(ctx context.Context, payload map[string]any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{Topic: s.topic, Value: sarama.ByteEncoder(data)}
	if s.key != "" {
		msg.Key = sarama.StringEncoder(s.key)
	}
	_, _, err = s.producer.SendMessage(msg)
	return err
}

type saramaSCRAMClient struct {
	conversation  *xdgscram.ClientConversation
	hashGenerator xdgscram.HashGeneratorFcn
}

func (c *saramaSCRAMClient) Begin(userName, password, authzID string) error {
	client, err := c.hashGenerator.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	c.conversation = client.NewConversation()
	return nil
}

func (c *saramaSCRAMClient) Step(challenge string) (string, error) {
	if c.conversation == nil {
		return "", fmt.Errorf("SCRAM 会话未初始化")
	}
	return c.conversation.Step(challenge)
}

func (c *saramaSCRAMClient) Done() bool {
	return c.conversation != nil && c.conversation.Done()
}

type kafkaTokenProvider struct {
	token string
}

func (p kafkaTokenProvider) Token() (*sarama.AccessToken, error) {
	if p.token == "" {
		return nil, fmt.Errorf("Kafka OAuthBearer 认证需要配置 token")
	}
	return &sarama.AccessToken{Token: p.token}, nil
}

func buildKafkaConfig(config map[string]any) (*sarama.Config, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Net.DialTimeout = configDuration(config, "dialTimeoutSeconds", 10*time.Second)
	kafkaConfig.Net.ReadTimeout = configDuration(config, "readTimeoutSeconds", 30*time.Second)
	kafkaConfig.Net.WriteTimeout = configDuration(config, "writeTimeoutSeconds", 30*time.Second)
	kafkaConfig.Consumer.Fetch.Min = int32(configInt(config, "minBytes", 1))
	kafkaConfig.Consumer.Fetch.Default = int32(configInt(config, "maxBytes", 5*1024*1024))
	kafkaConfig.Consumer.Fetch.Max = int32(configInt(config, "maxBytes", 5*1024*1024))
	kafkaConfig.Consumer.MaxWaitTime = configDuration(config, "maxWaitSeconds", 5*time.Second)
	kafkaConfig.Consumer.Offsets.Initial = kafkaInitialOffset(config)
	kafkaConfig.Producer.Timeout = configDuration(config, "timeoutSeconds", 10*time.Second)
	kafkaConfig.Version = sarama.V3_5_0_0
	if configBool(config, "tls") {
		kafkaConfig.Net.TLS.Enable = true
		kafkaConfig.Net.TLS.Config = &tls.Config{MinVersion: tls.VersionTLS12}
		if configBool(config, "insecureSkipVerify") {
			kafkaConfig.Net.TLS.Config.InsecureSkipVerify = true
		}
	}
	if err := configureKafkaSASL(kafkaConfig, config); err != nil {
		return nil, err
	}
	return kafkaConfig, nil
}

func configureKafkaSASL(kafkaConfig *sarama.Config, config map[string]any) error {
	authMode := strings.ToLower(strings.TrimSpace(configString(config, "authMode", "saslMechanism")))
	username := configString(config, "username")
	password := configString(config, "password")
	token := configString(config, "token", "oauthToken")
	switch authMode {
	case "", "none", "plaintext", "plain":
		if username == "" && password == "" && authMode != "plain" {
			return nil
		}
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
	case "scram-sha-256", "scram_sha_256":
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
		kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &saramaSCRAMClient{hashGenerator: xdgscram.SHA256}
		}
	case "scram-sha-512", "scram_sha_512":
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		kafkaConfig.Net.SASL.User = username
		kafkaConfig.Net.SASL.Password = password
		kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = func() sarama.SCRAMClient {
			return &saramaSCRAMClient{hashGenerator: xdgscram.SHA512}
		}
	case "oauthbearer", "oauth":
		kafkaConfig.Net.SASL.Enable = true
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		kafkaConfig.Net.SASL.TokenProvider = kafkaTokenProvider{token: token}
	default:
		return fmt.Errorf("不支持的 Kafka 认证模式: %s", authMode)
	}
	return nil
}

func kafkaInitialOffset(config map[string]any) int64 {
	switch strings.ToLower(strings.TrimSpace(configString(config, "startMode", "offsetMode"))) {
	case "earliest", "oldest", "beginning":
		return sarama.OffsetOldest
	default:
		return sarama.OffsetNewest
	}
}

func kafkaSpecificOffsets(config map[string]any) map[int32]int64 {
	if strings.ToLower(configString(config, "startMode", "offsetMode")) != "specific" {
		return nil
	}
	result := map[int32]int64{}
	if offset := configInt64(config, "offset", -1); offset >= 0 {
		result[0] = offset
	}
	raw, ok := config["offsets"].(map[string]any)
	if !ok {
		return result
	}
	for partition, offset := range raw {
		parsedPartition, err := parseInt32(partition)
		if err == nil {
			result[parsedPartition] = toInt64(offset)
		}
	}
	return result
}

func parseInt32(value string) (int32, error) {
	var parsed int64
	_, err := fmt.Sscan(value, &parsed)
	return int32(parsed), err
}

func configStringMap(config map[string]any, key string) map[string]string {
	value, ok := config[key]
	if !ok || value == nil {
		return nil
	}
	result := map[string]string{}
	if mapped, ok := value.(map[string]any); ok {
		for k, v := range mapped {
			result[k] = fmt.Sprint(v)
		}
	}
	if mapped, ok := value.(map[string]string); ok {
		for k, v := range mapped {
			result[k] = v
		}
	}
	return result
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}
