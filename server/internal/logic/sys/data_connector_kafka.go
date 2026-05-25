// Package sys
// @Link  https://github.com/bufanyun/hotgo
// @Copyright  Copyright (c) 2026 HotGo
// @License  https://github.com/bufanyun/hotgo/blob/master/LICENSE
package sys

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"hotgo/internal/consts"
	"hotgo/internal/dao"
	"hotgo/internal/model/entity"
	"hotgo/internal/model/input/sysin"

	"github.com/IBM/sarama"
	"github.com/gogf/gf/v2/errors/gerror"
	xdgscram "github.com/xdg-go/scram"
)

type dataConnectorKafkaAccessTokenProvider struct {
	token string
}

type dataConnectorSaramaSCRAMClient struct {
	conversation  *xdgscram.ClientConversation
	hashGenerator xdgscram.HashGeneratorFcn
}

func (p dataConnectorKafkaAccessTokenProvider) Token() (*sarama.AccessToken, error) {
	if strings.TrimSpace(p.token) == "" {
		return nil, fmt.Errorf("Kafka oauthbearer 认证需要配置 token")
	}
	return &sarama.AccessToken{Token: p.token}, nil
}

func (c *dataConnectorSaramaSCRAMClient) Begin(userName, password, authzID string) error {
	client, err := c.hashGenerator.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	c.conversation = client.NewConversation()
	return nil
}

func (c *dataConnectorSaramaSCRAMClient) Step(challenge string) (string, error) {
	if c.conversation == nil {
		return "", fmt.Errorf("SCRAM 会话未初始化")
	}
	return c.conversation.Step(challenge)
}

func (c *dataConnectorSaramaSCRAMClient) Done() bool {
	return c.conversation != nil && c.conversation.Done()
}

// KafkaTopics 获取Kafka Topic列表
func (s *sSysDataConnector) KafkaTopics(ctx context.Context, in *sysin.DataConnectorKafkaTopicsInp) (res *sysin.DataConnectorKafkaTopicsModel, err error) {
	connector, err := s.getKafkaTopicConnector(ctx, in.Id)
	if err != nil {
		return nil, err
	}

	config := jsonToMap(connector.Config)
	brokers := dataConnectorKafkaStringSlice(config, "brokers")
	if len(brokers) == 0 {
		return nil, gerror.New("Kafka Brokers不能为空")
	}

	kafkaConfig, err := buildDataConnectorKafkaConfig(config)
	if err != nil {
		return nil, gerror.Wrap(err, "Kafka配置无效")
	}

	admin, err := sarama.NewClusterAdmin(brokers, kafkaConfig)
	if err != nil {
		return nil, gerror.Wrap(err, "连接Kafka失败")
	}
	defer admin.Close()

	topicDetails, err := admin.ListTopics()
	if err != nil {
		return nil, gerror.Wrap(err, "读取Kafka Topic列表失败")
	}

	topics := make([]string, 0, len(topicDetails))
	for topic := range topicDetails {
		if isDataConnectorKafkaInternalTopic(topic) {
			continue
		}
		topics = append(topics, topic)
	}
	sort.Strings(topics)

	return &sysin.DataConnectorKafkaTopicsModel{Topics: topics}, nil
}

func isDataConnectorKafkaInternalTopic(topic string) bool {
	return strings.HasPrefix(strings.TrimSpace(topic), "__")
}

func (s *sSysDataConnector) getKafkaTopicConnector(ctx context.Context, id int64) (connector entity.DataConnector, err error) {
	cols := dao.DataConnector.Columns()
	if err = s.Model(ctx).
		Where(cols.Id, id).
		Where(cols.Direction, consts.DataConnectorDirectionSource).
		Where(cols.ConnectorType, consts.DataConnectorTypeKafka).
		Where(cols.Status, consts.StatusEnabled).
		Scan(&connector); err != nil {
		return connector, gerror.Wrap(err, "获取Kafka数据源失败")
	}
	if connector.Id <= 0 {
		return connector, gerror.New("Kafka数据源不存在或已禁用")
	}
	return connector, nil
}

func buildDataConnectorKafkaConfig(config map[string]interface{}) (*sarama.Config, error) {
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.ClientID = "hotgo-data-connector"
	kafkaConfig.Version = sarama.V3_5_0_0
	kafkaConfig.Net.DialTimeout = dataConnectorKafkaDuration(config, "dialTimeoutSeconds", 10*time.Second)
	kafkaConfig.Net.ReadTimeout = dataConnectorKafkaDuration(config, "readTimeoutSeconds", 30*time.Second)
	kafkaConfig.Net.WriteTimeout = dataConnectorKafkaDuration(config, "writeTimeoutSeconds", 30*time.Second)
	kafkaConfig.Metadata.Timeout = dataConnectorKafkaDuration(config, "metadataTimeoutSeconds", 10*time.Second)

	tlsConfig, err := buildDataConnectorKafkaTLSConfig(config)
	if err != nil {
		return nil, err
	}
	if tlsConfig != nil {
		kafkaConfig.Net.TLS.Enable = true
		kafkaConfig.Net.TLS.Config = tlsConfig
	}
	if err := applyDataConnectorKafkaSASLConfig(kafkaConfig, config); err != nil {
		return nil, err
	}
	if err := kafkaConfig.Validate(); err != nil {
		return nil, err
	}
	return kafkaConfig, nil
}

func buildDataConnectorKafkaTLSConfig(config map[string]interface{}) (*tls.Config, error) {
	securityProtocol := strings.ToLower(dataConnectorKafkaString(config, "securityProtocol"))
	caCert := dataConnectorKafkaString(config, "caCert", "caPem", "caCertificate")
	clientCert := dataConnectorKafkaString(config, "clientCert", "clientCertificate")
	clientKey := dataConnectorKafkaString(config, "clientKey", "clientKeyPem")
	enabled := dataConnectorKafkaBool(config, "tls") ||
		strings.Contains(securityProtocol, "ssl") ||
		caCert != "" ||
		clientCert != "" ||
		clientKey != ""
	if !enabled {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: dataConnectorKafkaBool(config, "insecureSkipVerify"),
	}
	if caCert != "" {
		pool := x509.NewCertPool()
		if !pool.AppendCertsFromPEM([]byte(caCert)) {
			return nil, fmt.Errorf("Kafka TLS CA证书格式无效")
		}
		tlsConfig.RootCAs = pool
	}
	if clientCert != "" || clientKey != "" {
		if clientCert == "" || clientKey == "" {
			return nil, fmt.Errorf("Kafka mTLS需要同时配置客户端证书和私钥")
		}
		cert, err := tls.X509KeyPair([]byte(clientCert), []byte(clientKey))
		if err != nil {
			return nil, fmt.Errorf("Kafka mTLS客户端证书或私钥无效: %w", err)
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	return tlsConfig, nil
}

func applyDataConnectorKafkaSASLConfig(kafkaConfig *sarama.Config, config map[string]interface{}) error {
	mode := dataConnectorKafkaAuthMode(config)
	if mode == "" || mode == consts.DataKafkaAuthModeNone {
		return nil
	}

	kafkaConfig.Net.SASL.Enable = true
	if mode == consts.DataKafkaAuthModeOAuthBearer {
		token := dataConnectorKafkaString(config, "token", "oauthToken", "bearerToken")
		if token == "" {
			return fmt.Errorf("Kafka oauthbearer认证需要配置token")
		}
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeOAuth
		kafkaConfig.Net.SASL.TokenProvider = dataConnectorKafkaAccessTokenProvider{token: token}
		return nil
	}

	username := dataConnectorKafkaString(config, "username", "saslUsername")
	password := dataConnectorKafkaString(config, "password", "saslPassword")
	if username == "" || password == "" {
		return fmt.Errorf("Kafka %s认证需要配置username、password", mode)
	}
	kafkaConfig.Net.SASL.User = username
	kafkaConfig.Net.SASL.Password = password

	switch mode {
	case consts.DataKafkaAuthModePlain:
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	case consts.DataKafkaAuthModeScramSha256:
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA256
		kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = dataConnectorKafkaSCRAMClientGenerator(xdgscram.SHA256)
	case consts.DataKafkaAuthModeScramSha512:
		kafkaConfig.Net.SASL.Mechanism = sarama.SASLTypeSCRAMSHA512
		kafkaConfig.Net.SASL.SCRAMClientGeneratorFunc = dataConnectorKafkaSCRAMClientGenerator(xdgscram.SHA512)
	default:
		return fmt.Errorf("Kafka认证模式无效: %s", mode)
	}
	return nil
}

func dataConnectorKafkaSCRAMClientGenerator(hashGenerator xdgscram.HashGeneratorFcn) func() sarama.SCRAMClient {
	return func() sarama.SCRAMClient {
		return &dataConnectorSaramaSCRAMClient{hashGenerator: hashGenerator}
	}
}

func dataConnectorKafkaAuthMode(config map[string]interface{}) string {
	mode := strings.ToLower(strings.TrimSpace(dataConnectorKafkaString(config, "authMode", "mechanism")))
	if mode == "" {
		mode = strings.ToLower(strings.TrimSpace(dataConnectorKafkaString(config, "saslMechanism")))
	}
	mode = strings.ReplaceAll(mode, "_", "-")
	switch mode {
	case "", "none", "no-auth", "anonymous":
		return consts.DataKafkaAuthModeNone
	case "sasl/plain", "sasl-plain", "plain":
		return consts.DataKafkaAuthModePlain
	case "scram-sha256", "scram-sha-256", "sasl/scram-sha-256":
		return consts.DataKafkaAuthModeScramSha256
	case "scram-sha512", "scram-sha-512", "sasl/scram-sha-512":
		return consts.DataKafkaAuthModeScramSha512
	case "oauth", "oauth-bearer", "oauthbearer", "sasl/oauthbearer":
		return consts.DataKafkaAuthModeOAuthBearer
	default:
		return mode
	}
}

func dataConnectorKafkaString(config map[string]interface{}, keys ...string) string {
	for _, key := range keys {
		if value := dataConnectorKafkaSingleString(config[key]); value != "" {
			return value
		}
	}

	nested, _ := config["sasl"].(map[string]interface{})
	for _, key := range keys {
		if value := dataConnectorKafkaSingleString(nested[key]); value != "" {
			return value
		}
	}
	return ""
}

func dataConnectorKafkaSingleString(value interface{}) string {
	if value == nil {
		return ""
	}
	switch typed := value.(type) {
	case string:
		return strings.TrimSpace(typed)
	default:
		text := strings.TrimSpace(fmt.Sprint(typed))
		if text == "<nil>" {
			return ""
		}
		return text
	}
}

func dataConnectorKafkaStringSlice(config map[string]interface{}, keys ...string) []string {
	for _, key := range keys {
		value, ok := config[key]
		if !ok || value == nil {
			continue
		}
		switch typed := value.(type) {
		case []string:
			return dataConnectorKafkaCompactStrings(typed)
		case []interface{}:
			items := make([]string, 0, len(typed))
			for _, item := range typed {
				items = append(items, fmt.Sprint(item))
			}
			return dataConnectorKafkaCompactStrings(items)
		case string:
			return dataConnectorKafkaCompactStrings(strings.Split(typed, ","))
		default:
			text := strings.TrimSpace(fmt.Sprint(typed))
			if text != "" && text != "<nil>" {
				return []string{text}
			}
		}
	}
	return nil
}

func dataConnectorKafkaCompactStrings(items []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func dataConnectorKafkaBool(config map[string]interface{}, key string) bool {
	value, ok := config[key]
	if !ok || value == nil {
		return false
	}
	switch typed := value.(type) {
	case bool:
		return typed
	case string:
		parsed, _ := strconv.ParseBool(strings.TrimSpace(typed))
		return parsed
	default:
		return fmt.Sprint(typed) == "1"
	}
}

func dataConnectorKafkaDuration(config map[string]interface{}, key string, fallback time.Duration) time.Duration {
	value, ok := config[key]
	if !ok || value == nil {
		return fallback
	}

	seconds := int64(0)
	switch typed := value.(type) {
	case int:
		seconds = int64(typed)
	case int64:
		seconds = typed
	case float64:
		seconds = int64(typed)
	case string:
		parsed, _ := strconv.ParseInt(strings.TrimSpace(typed), 10, 64)
		seconds = parsed
	default:
		parsed, _ := strconv.ParseInt(strings.TrimSpace(fmt.Sprint(typed)), 10, 64)
		seconds = parsed
	}
	if seconds <= 0 {
		return fallback
	}
	return time.Duration(seconds) * time.Second
}
