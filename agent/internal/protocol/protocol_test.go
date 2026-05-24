package protocol

import (
	"encoding/json"
	"testing"
	"time"
)

func TestHelloEnvelopeCarriesCapabilitiesAndSlots(t *testing.T) {
	msg := NewHelloEnvelope("agent-a", "node-a", []string{"kafka", "s3"}, []string{"10.0.0.5"}, SlotStatus{
		Total:     16,
		Available: 12,
	}, "1.0.0")

	if msg.Type != MessageAgentHello {
		t.Fatalf("unexpected type: %s", msg.Type)
	}
	var payload HelloPayload
	if err := DecodePayload(msg, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.AgentId != "agent-a" || payload.Hostname != "node-a" {
		t.Fatalf("unexpected identity: %+v", payload)
	}
	if payload.Slots.Available != 12 || payload.Slots.Total != 16 {
		t.Fatalf("unexpected slots: %+v", payload.Slots)
	}
	if len(payload.Capabilities) != 2 {
		t.Fatalf("unexpected capabilities: %+v", payload.Capabilities)
	}
	if len(payload.AgentIp) != 1 || payload.AgentIp[0] != "10.0.0.5" {
		t.Fatalf("unexpected agent ip: %+v", payload.AgentIp)
	}
}

func TestTaskConfigEnvelopeCarriesIntegrationCleanOutputsAndParallelism(t *testing.T) {
	raw := []byte(`{
		"type":"server.task.config",
		"requestId":"req-1",
		"sentAt":"2026-05-22T00:00:00Z",
		"payload":{
			"tenantId":"default",
			"taskId":1001,
			"configVersion":7,
			"shards":[{"shardKey":"kafka:risk-events:0"}],
			"integration":{"sourceType":"kafka","config":{"topic":"risk-events"}},
			"clean":{"enabled":true,"parseDepth":3,"steps":[{"type":"trim","field":"name"}]},
			"outputs":[{"type":"kafka","config":{"topic":"clean-events"}}],
			"parallelism":{"readerConcurrency":4,"cleanWorkers":16,"outputWorkers":8,"maxInFlight":5000,"batchSize":200}
		}
	}`)

	var msg Envelope
	if err := json.Unmarshal(raw, &msg); err != nil {
		t.Fatalf("decode envelope: %v", err)
	}
	var payload TaskConfigPayload
	if err := DecodePayload(msg, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.TaskId != 1001 || payload.ConfigVersion != 7 {
		t.Fatalf("unexpected task identity: %+v", payload)
	}
	if payload.Integration.SourceType != "kafka" {
		t.Fatalf("unexpected integration: %+v", payload.Integration)
	}
	if !payload.Clean.Enabled || payload.Clean.ParseDepth != 3 || len(payload.Outputs) != 1 {
		t.Fatalf("unexpected clean/output: %+v", payload)
	}
	if payload.Parallelism.CleanWorkers != 16 || payload.Parallelism.OutputWorkers != 8 {
		t.Fatalf("unexpected parallelism: %+v", payload.Parallelism)
	}
}

func TestHeartbeatEnvelopeUsesCurrentSlotStatus(t *testing.T) {
	now := time.Date(2026, 5, 22, 0, 0, 0, 0, time.UTC)
	msg := NewHeartbeatEnvelope("agent-a", SlotStatus{Total: 8, Available: 5}, RuntimeMetrics{
		RunningShards: 3,
		ReadTps:       120,
		CleanTps:      118,
		OutputTps:     115,
	}, []string{"10.0.0.5"}, now)

	var payload HeartbeatPayload
	if err := DecodePayload(msg, &payload); err != nil {
		t.Fatalf("decode payload: %v", err)
	}
	if payload.AgentId != "agent-a" || payload.Slots.Available != 5 {
		t.Fatalf("unexpected heartbeat: %+v", payload)
	}
	if payload.Metrics.RunningShards != 3 || !payload.Time.Equal(now) {
		t.Fatalf("unexpected heartbeat metrics: %+v", payload)
	}
	if len(payload.AgentIp) != 1 || payload.AgentIp[0] != "10.0.0.5" {
		t.Fatalf("unexpected agent ip: %+v", payload.AgentIp)
	}
}
