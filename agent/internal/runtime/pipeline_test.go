package runtime

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"vogo-agent/internal/protocol"
)

type fakeSource struct {
	records []Record
}

func (s fakeSource) Run(ctx context.Context, emit func(Record) error) error {
	for _, record := range s.records {
		if err := ctx.Err(); err != nil {
			return err
		}
		if err := emit(record); err != nil {
			return err
		}
	}
	return nil
}

type fakeSink struct {
	writes *[]map[string]any
}

func (s fakeSink) Write(ctx context.Context, payload map[string]any) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	*s.writes = append(*s.writes, cloneMap(payload))
	return nil
}

func TestDefaultConnectorRegistryExposesRequiredSourcesAndSinks(t *testing.T) {
	registry := NewDefaultConnectorRegistry()

	for _, sourceType := range []string{"kafka", "http", "s3", "log"} {
		if !registry.HasSource(sourceType) {
			t.Fatalf("expected source connector %q to be registered", sourceType)
		}
	}
	for _, sinkType := range []string{"kafka", "http", "s3", "elasticsearch", "opensearch", "splunk", "feishu", "lark", "strategy", "log"} {
		if !registry.HasSink(sinkType) {
			t.Fatalf("expected sink connector %q to be registered", sinkType)
		}
	}
}

func TestTaskRunnerFiltersCleansCollectsFieldsAndWritesMultipleSinks(t *testing.T) {
	ctx := context.Background()
	registry := NewConnectorRegistry()
	sourceRecords := []Record{
		{Payload: map[string]any{"status": "drop", "name": " ignore ", "age": "20"}},
		{Payload: map[string]any{"status": "ok", "name": " alice ", "age": "30", "secret": "hidden"}},
	}
	var primaryWrites []map[string]any
	var secondaryWrites []map[string]any
	registry.RegisterSource("fake", func(_ protocol.IntegrationConfig, _ protocol.TaskShard) (Source, error) {
		return fakeSource{records: sourceRecords}, nil
	})
	registry.RegisterSink("primary", func(_ protocol.OutputConfig) (Sink, error) {
		return fakeSink{writes: &primaryWrites}, nil
	})
	registry.RegisterSink("secondary", func(_ protocol.OutputConfig) (Sink, error) {
		return fakeSink{writes: &secondaryWrites}, nil
	})

	var reports []protocol.FieldsReportPayload
	runner := NewTaskRunner("agent-a", registry)
	runner.SetFieldReporter(func(report protocol.FieldsReportPayload) {
		reports = append(reports, report)
	})

	err := runner.RunTaskOnce(ctx, protocol.TaskConfigPayload{
		TenantId:      "default",
		TaskId:        1001,
		ConfigVersion: 1,
		Integration: protocol.IntegrationConfig{
			SourceType: "fake",
		},
		Clean: protocol.CleanConfig{
			Enabled: true,
			Filter: map[string]any{
				"enabled": true,
				"logic":   "and",
				"conditions": []any{
					map[string]any{"field": "status", "operator": "eq", "value": "ok"},
				},
			},
			Steps: []map[string]any{
				{"type": "trim", "field": "name"},
				{"type": "to_int", "field": "age"},
				{"type": "remove", "field": "secret"},
			},
		},
		Outputs: []protocol.OutputConfig{
			{Type: "primary"},
			{Type: "secondary"},
		},
	})
	if err != nil {
		t.Fatalf("run task: %v", err)
	}

	if len(primaryWrites) != 1 || !payloadHasValues(primaryWrites[0], map[string]any{"status": "ok", "name": "alice", "age": int64(30), "tenantId": "default", "ingestTaskId": int64(1001)}) {
		t.Fatalf("unexpected primary writes: %#v", primaryWrites)
	}
	if len(secondaryWrites) != 1 || !payloadHasValues(secondaryWrites[0], map[string]any{"status": "ok", "name": "alice", "age": int64(30), "tenantId": "default", "ingestTaskId": int64(1001)}) {
		t.Fatalf("unexpected secondary writes: %#v", secondaryWrites)
	}
	if len(reports) == 0 || reports[0].TenantId != "default" || reports[0].TaskId != 1001 {
		t.Fatalf("expected field report for task, got %#v", reports)
	}
}

func TestTaskRunnerExecutesStrategyOutputLocallyWithoutServerPayloadCallback(t *testing.T) {
	ctx := context.Background()
	registry := NewConnectorRegistry()
	var alertWrites []map[string]any
	registry.RegisterSource("fake", func(_ protocol.IntegrationConfig, _ protocol.TaskShard) (Source, error) {
		return fakeSource{records: []Record{{Payload: map[string]any{"eventType": "risk.event", "riskScore": int64(90), "username": "admin"}}}}, nil
	})
	registry.RegisterSink("alert", func(_ protocol.OutputConfig) (Sink, error) {
		return fakeSink{writes: &alertWrites}, nil
	})
	registry.RegisterSink("strategy", func(config protocol.OutputConfig) (Sink, error) {
		return newStrategySink(config, registry)
	})

	runner := NewTaskRunner("agent-a", registry)
	err := runner.RunTaskOnce(ctx, protocol.TaskConfigPayload{
		TenantId:      "default",
		TaskId:        1002,
		ConfigVersion: 1,
		Integration:   protocol.IntegrationConfig{SourceType: "fake"},
		Outputs: []protocol.OutputConfig{{
			Type: "strategy",
			Config: map[string]any{
				"rules": []any{
					map[string]any{
						"id":                  float64(1),
						"name":                "high risk",
						"eventType":           "risk.event",
						"decision":            "review",
						"conditionExpression": "riskScore >= 80",
						"actions": []any{
							map[string]any{
								"actionType": "output.send",
								"actionParams": map[string]any{
									"sinkCode": "risk-alert",
									"message":  "高危事件",
								},
							},
						},
					},
				},
				"sinks": map[string]any{
					"risk-alert": map[string]any{
						"type":   "alert",
						"config": map[string]any{},
					},
				},
			},
		}},
	})
	if err != nil {
		t.Fatalf("run task: %v", err)
	}
	if len(alertWrites) != 1 || alertWrites[0]["message"] != "高危事件" || toInt64(alertWrites[0]["riskScore"]) != 90 {
		t.Fatalf("expected local strategy output, got %#v", alertWrites)
	}
}

func TestTaskRunnerDrainsWindowStatsWithoutPayload(t *testing.T) {
	ctx := context.Background()
	registry := NewConnectorRegistry()
	var writes []map[string]any
	registry.RegisterSource("fake", func(_ protocol.IntegrationConfig, _ protocol.TaskShard) (Source, error) {
		return fakeSource{records: []Record{{Payload: map[string]any{"status": "ok"}}}}, nil
	})
	registry.RegisterSink("out", func(_ protocol.OutputConfig) (Sink, error) {
		return fakeSink{writes: &writes}, nil
	})
	runner := NewTaskRunner("agent-a", registry)
	if err := runner.RunTaskOnce(ctx, protocol.TaskConfigPayload{
		TenantId:    "default",
		TaskId:      1003,
		Integration: protocol.IntegrationConfig{SourceType: "fake"},
		Outputs:     []protocol.OutputConfig{{Type: "out"}},
	}); err != nil {
		t.Fatalf("run task: %v", err)
	}

	report := runner.DrainStats()
	if report.AgentId != "agent-a" || len(report.Metrics) == 0 {
		t.Fatalf("expected stats report, got %#v", report)
	}
	rawReport, _ := json.Marshal(report)
	if strings.Contains(string(rawReport), "payload") || strings.Contains(string(rawReport), "username") {
		t.Fatalf("stats report must not carry business payload: %s", rawReport)
	}
	if findMetricValue(report.Metrics, "records.read") != 1 || findMetricValue(report.Metrics, "records.output") != 1 {
		t.Fatalf("unexpected metrics: %#v", report.Metrics)
	}
}

func TestParseRecordPayloadsSupportsJSONLinesAndHeaderText(t *testing.T) {
	jsonPayloads, err := ParseRecordPayloads(Record{Raw: []byte("{\"id\":1}\n{\"id\":2}")}, map[string]any{"format": "jsonl"})
	if err != nil {
		t.Fatalf("parse json lines: %v", err)
	}
	if len(jsonPayloads) != 2 || toInt64(jsonPayloads[1]["id"]) != 2 {
		t.Fatalf("unexpected json lines payloads: %#v", jsonPayloads)
	}

	textPayloads, err := ParseRecordPayloads(Record{Raw: []byte("version account-id bytes\n2 123456789012 42")}, map[string]any{
		"format":    "text",
		"delimiter": "space",
		"header":    "version account-id bytes",
	})
	if err != nil {
		t.Fatalf("parse header text: %v", err)
	}
	if len(textPayloads) != 1 || textPayloads[0]["account-id"] != "123456789012" {
		t.Fatalf("unexpected text payloads: %#v", textPayloads)
	}
}

func findMetricValue(metrics []protocol.StatsMetric, name string) int64 {
	for _, metric := range metrics {
		if metric.Name == name {
			return metric.Value
		}
	}
	return 0
}

func cloneMap(input map[string]any) map[string]any {
	output := make(map[string]any, len(input))
	for key, value := range input {
		output[key] = value
	}
	return output
}

func payloadHasValues(payload map[string]any, expected map[string]any) bool {
	for key, value := range expected {
		if _, ok := value.(int64); ok {
			if toInt64(payload[key]) != value.(int64) {
				return false
			}
			continue
		}
		if !reflect.DeepEqual(payload[key], value) {
			return false
		}
	}
	return true
}
