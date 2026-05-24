package runtime

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"vogo-agent/internal/protocol"
)

type Record struct {
	Raw      []byte
	Payload  map[string]any
	Metadata map[string]any
}

type Source interface {
	Run(ctx context.Context, emit func(Record) error) error
}

type Sink interface {
	Write(ctx context.Context, payload map[string]any) error
}

type SourceFactory func(config protocol.IntegrationConfig, shard protocol.TaskShard) (Source, error)
type SinkFactory func(config protocol.OutputConfig) (Sink, error)

type ConnectorRegistry struct {
	sources map[string]SourceFactory
	sinks   map[string]SinkFactory
}

func NewConnectorRegistry() *ConnectorRegistry {
	return &ConnectorRegistry{
		sources: map[string]SourceFactory{},
		sinks:   map[string]SinkFactory{},
	}
}

func NewDefaultConnectorRegistry() *ConnectorRegistry {
	registry := NewConnectorRegistry()
	registerDefaultConnectors(registry)
	return registry
}

func (r *ConnectorRegistry) RegisterSource(sourceType string, factory SourceFactory) {
	if r == nil || factory == nil {
		return
	}
	r.sources[normalizeConnectorType(sourceType)] = factory
}

func (r *ConnectorRegistry) RegisterSink(sinkType string, factory SinkFactory) {
	if r == nil || factory == nil {
		return
	}
	r.sinks[normalizeConnectorType(sinkType)] = factory
}

func (r *ConnectorRegistry) HasSource(sourceType string) bool {
	_, ok := r.sources[normalizeConnectorType(sourceType)]
	return ok
}

func (r *ConnectorRegistry) HasSink(sinkType string) bool {
	_, ok := r.sinks[normalizeConnectorType(sinkType)]
	return ok
}

func (r *ConnectorRegistry) BuildSource(config protocol.IntegrationConfig, shard protocol.TaskShard) (Source, error) {
	factory, ok := r.sources[normalizeConnectorType(config.SourceType)]
	if !ok {
		return nil, fmt.Errorf("未注册数据输入 connector: %s", config.SourceType)
	}
	return factory(config, shard)
}

func (r *ConnectorRegistry) BuildSink(config protocol.OutputConfig) (Sink, error) {
	factory, ok := r.sinks[normalizeConnectorType(config.Type)]
	if !ok {
		return nil, fmt.Errorf("未注册数据输出 connector: %s", config.Type)
	}
	return factory(config)
}

func normalizeConnectorType(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

type TaskRunner struct {
	agentId       string
	registry      *ConnectorRegistry
	mu            sync.Mutex
	cancels       map[int64]runningTask
	versions      map[int64]int64
	fieldReporter func(protocol.FieldsReportPayload)
	statsMu       sync.Mutex
	statsWindow   time.Time
	stats         map[string]protocol.StatsMetric
	readCount     atomic.Int64
	cleanCount    atomic.Int64
	outputCount   atomic.Int64
	runningShards atomic.Int64
}

type runningTask struct {
	cancel  context.CancelFunc
	version int64
}

func NewTaskRunner(agentId string, registry *ConnectorRegistry) *TaskRunner {
	if registry == nil {
		registry = NewDefaultConnectorRegistry()
	}
	return &TaskRunner{
		agentId:  agentId,
		registry: registry,
		cancels:  map[int64]runningTask{},
		versions: map[int64]int64{},
		stats:    map[string]protocol.StatsMetric{},
	}
}

func (r *TaskRunner) SetFieldReporter(reporter func(protocol.FieldsReportPayload)) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fieldReporter = reporter
}

func (r *TaskRunner) ApplyTaskConfig(ctx context.Context, config protocol.TaskConfigPayload) error {
	if r == nil {
		return errors.New("task runner 未初始化")
	}
	r.mu.Lock()
	if r.versions[config.TaskId] == config.ConfigVersion {
		r.mu.Unlock()
		return ctx.Err()
	}
	if running, ok := r.cancels[config.TaskId]; ok && running.cancel != nil {
		running.cancel()
	}
	taskCtx, cancel := context.WithCancel(ctx)
	r.cancels[config.TaskId] = runningTask{cancel: cancel, version: config.ConfigVersion}
	r.versions[config.TaskId] = config.ConfigVersion
	r.mu.Unlock()

	go func() {
		_ = r.runTask(taskCtx, config)
		r.mu.Lock()
		if running, ok := r.cancels[config.TaskId]; ok && running.version == config.ConfigVersion {
			delete(r.cancels, config.TaskId)
		}
		r.mu.Unlock()
	}()
	return ctx.Err()
}

func (r *TaskRunner) StopTask(ctx context.Context, taskId int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if running, ok := r.cancels[taskId]; ok && running.cancel != nil {
		running.cancel()
		delete(r.cancels, taskId)
	}
	delete(r.versions, taskId)
	return ctx.Err()
}

func (r *TaskRunner) StopAll() {
	r.mu.Lock()
	defer r.mu.Unlock()
	for taskId, running := range r.cancels {
		if running.cancel != nil {
			running.cancel()
		}
		delete(r.cancels, taskId)
	}
	clear(r.versions)
}

func (r *TaskRunner) Metrics() protocol.RuntimeMetrics {
	return protocol.RuntimeMetrics{
		RunningShards: int(r.runningShards.Load()),
		ReadTps:       float64(r.readCount.Load()),
		CleanTps:      float64(r.cleanCount.Load()),
		OutputTps:     float64(r.outputCount.Load()),
	}
}

func (r *TaskRunner) DrainStats() protocol.StatsReportPayload {
	now := time.Now()
	r.statsMu.Lock()
	defer r.statsMu.Unlock()
	windowStart := r.statsWindow
	if windowStart.IsZero() {
		windowStart = now
	}
	metrics := make([]protocol.StatsMetric, 0, len(r.stats))
	for _, metric := range r.stats {
		metrics = append(metrics, metric)
	}
	r.stats = map[string]protocol.StatsMetric{}
	r.statsWindow = now
	return protocol.StatsReportPayload{
		AgentId:     r.agentId,
		WindowStart: windowStart,
		WindowEnd:   now,
		Metrics:     metrics,
	}
}

func (r *TaskRunner) RunTaskOnce(ctx context.Context, config protocol.TaskConfigPayload) error {
	return r.runTask(ctx, config)
}

func (r *TaskRunner) runTask(ctx context.Context, config protocol.TaskConfigPayload) error {
	shards := config.Shards
	if len(shards) == 0 {
		shards = []protocol.TaskShard{{ShardKey: "default"}}
	}
	sinks, err := r.buildSinks(config.Outputs)
	if err != nil {
		return err
	}
	collector := NewFieldCollector(config.FieldSnapshot)
	var collectorMu sync.Mutex
	errCh := make(chan error, len(shards))
	var wg sync.WaitGroup
	for _, shard := range shards {
		if err := ctx.Err(); err != nil {
			return err
		}
		shard := shard
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := r.runShard(ctx, config, shard, collector, &collectorMu, sinks); err != nil && !errors.Is(err, context.Canceled) {
				errCh <- err
			}
		}()
	}
	wg.Wait()
	close(errCh)
	for err := range errCh {
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *TaskRunner) runShard(ctx context.Context, config protocol.TaskConfigPayload, shard protocol.TaskShard, collector *FieldCollector, collectorMu *sync.Mutex, sinks []Sink) error {
	sourceConfig := config.Integration
	sourceConfig.Config = mergeConfig(config.Integration.Config, shard.Config)
	sourceConfig.Config["_tenantId"] = config.TenantId
	sourceConfig.Config["_taskId"] = config.TaskId
	sourceConfig.Config["_configVersion"] = config.ConfigVersion
	source, err := r.registry.BuildSource(sourceConfig, shard)
	if err != nil {
		return err
	}
	r.runningShards.Add(1)
	err = source.Run(ctx, func(record Record) error {
		return r.handleRecord(ctx, config, collector, collectorMu, sinks, record)
	})
	r.runningShards.Add(-1)
	if err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func (r *TaskRunner) buildSinks(outputs []protocol.OutputConfig) ([]Sink, error) {
	sinks := make([]Sink, 0, len(outputs))
	for _, output := range outputs {
		sink, err := r.registry.BuildSink(output)
		if err != nil {
			return nil, err
		}
		sinks = append(sinks, sink)
	}
	return sinks, nil
}

func (r *TaskRunner) handleRecord(ctx context.Context, config protocol.TaskConfigPayload, collector *FieldCollector, collectorMu *sync.Mutex, sinks []Sink, record Record) error {
	payloads, err := ParseRecordPayloads(record, config.Integration.Config)
	if err != nil {
		return err
	}
	for _, payload := range payloads {
		r.readCount.Add(1)
		r.addMetric(config, "records.read", "", 1)
		collectorMu.Lock()
		fields := collector.Collect(payload)
		collectorMu.Unlock()
		r.reportFields(config, fields)
		cleaned, err := ApplyCleanPipeline(payload, config.Clean)
		if err != nil {
			return err
		}
		if cleaned.Dropped {
			r.addMetric(config, "records.dropped", "", 1)
			continue
		}
		r.cleanCount.Add(1)
		r.addMetric(config, "records.cleaned", "", 1)
		outputPayload := enrichOutputPayload(config, cleaned.Payload)
		for _, sink := range sinks {
			if err := sink.Write(ctx, outputPayload); err != nil {
				return err
			}
			r.outputCount.Add(1)
			r.addMetric(config, "records.output", "", 1)
		}
	}
	return nil
}

func (r *TaskRunner) addMetric(config protocol.TaskConfigPayload, name, groupKey string, value int64) {
	r.statsMu.Lock()
	defer r.statsMu.Unlock()
	if r.statsWindow.IsZero() {
		r.statsWindow = time.Now()
	}
	key := fmt.Sprintf("%s:%d:%s:%s", config.TenantId, config.TaskId, name, groupKey)
	metric := r.stats[key]
	metric.TenantId = config.TenantId
	metric.TaskId = config.TaskId
	metric.Name = name
	metric.GroupKey = groupKey
	metric.Value += value
	r.stats[key] = metric
}

func enrichOutputPayload(config protocol.TaskConfigPayload, payload map[string]any) map[string]any {
	enriched := clonePayload(payload)
	if _, ok := enriched["tenantId"]; !ok && config.TenantId != "" {
		enriched["tenantId"] = config.TenantId
	}
	if _, ok := enriched["ingestTaskId"]; !ok && config.TaskId > 0 {
		enriched["ingestTaskId"] = config.TaskId
	}
	if _, ok := enriched["eventId"]; !ok {
		enriched["eventId"] = time.Now().UnixNano()
	}
	if _, ok := enriched["eventType"]; !ok {
		eventType := configString(config.Integration.Config, "eventType", "_eventType")
		if eventType == "" {
			eventType = "custom.event"
		}
		enriched["eventType"] = eventType
	}
	sourceCode := configString(config.Integration.Config, "sourceCode", "_sourceCode")
	if sourceCode != "" {
		enriched["sourceCode"] = sourceCode
	}
	return enriched
}

func (r *TaskRunner) reportFields(config protocol.TaskConfigPayload, fields []protocol.FieldSample) {
	if len(fields) == 0 {
		return
	}
	r.mu.Lock()
	reporter := r.fieldReporter
	r.mu.Unlock()
	if reporter == nil {
		return
	}
	reporter(protocol.FieldsReportPayload{
		TenantId: config.TenantId,
		TaskId:   config.TaskId,
		AgentId:  r.agentId,
		Fields:   fields,
	})
}

func ParseRecordPayloads(record Record, config map[string]any) ([]map[string]any, error) {
	if record.Payload != nil {
		return []map[string]any{clonePayload(record.Payload)}, nil
	}
	raw := bytes.TrimSpace(record.Raw)
	if len(raw) == 0 {
		return nil, nil
	}
	format := strings.ToLower(strings.TrimSpace(configString(config, "format", "payloadFormat")))
	if format == "" || format == "json" {
		payloads, err := parseJSONPayloads(raw)
		if err == nil {
			return payloads, nil
		}
		if format == "json" {
			return nil, err
		}
	}
	switch format {
	case "jsonl", "json_lines", "ndjson":
		return parseJSONLines(raw)
	case "text", "delimited", "csv", "tsv":
		return parseDelimitedText(raw, config)
	default:
		return nil, fmt.Errorf("不支持的数据解析格式: %s", format)
	}
}

func parseJSONPayloads(raw []byte) ([]map[string]any, error) {
	var one map[string]any
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&one); err == nil {
		return []map[string]any{one}, nil
	}
	var many []map[string]any
	decoder = json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	if err := decoder.Decode(&many); err != nil {
		return nil, err
	}
	return many, nil
}

func parseJSONLines(raw []byte) ([]map[string]any, error) {
	lines := bytes.Split(raw, []byte("\n"))
	payloads := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		var payload map[string]any
		decoder := json.NewDecoder(bytes.NewReader(line))
		decoder.UseNumber()
		if err := decoder.Decode(&payload); err != nil {
			return nil, err
		}
		payloads = append(payloads, payload)
	}
	return payloads, nil
}

func parseDelimitedText(raw []byte, config map[string]any) ([]map[string]any, error) {
	header := configString(config, "header", "headers")
	delimiter := configString(config, "delimiter")
	lines := strings.Split(strings.TrimSpace(string(raw)), "\n")
	if len(lines) == 0 {
		return nil, nil
	}
	var headers []string
	startLine := 0
	if header != "" {
		headers = splitDelimitedLine(header, delimiter)
		if len(lines) > 0 && equalStringSlices(headers, splitDelimitedLine(lines[0], delimiter)) {
			startLine = 1
		}
	} else {
		headers = splitDelimitedLine(lines[0], delimiter)
		startLine = 1
	}
	payloads := make([]map[string]any, 0, len(lines)-startLine)
	for _, line := range lines[startLine:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		values := splitDelimitedLine(line, delimiter)
		payload := map[string]any{}
		for index, header := range headers {
			if index < len(values) {
				payload[header] = values[index]
			}
		}
		payloads = append(payloads, payload)
	}
	return payloads, nil
}

func splitDelimitedLine(line, delimiter string) []string {
	if delimiter == "space" || delimiter == "" {
		return strings.Fields(line)
	}
	if delimiter == "\\t" || delimiter == "tab" {
		delimiter = "\t"
	}
	if delimiter == "," || delimiter == "csv" {
		reader := csv.NewReader(strings.NewReader(line))
		records, err := reader.Read()
		if err == nil {
			return records
		}
	}
	return strings.Split(line, delimiter)
}

func equalStringSlices(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if strings.TrimSpace(left[index]) != strings.TrimSpace(right[index]) {
			return false
		}
	}
	return true
}
