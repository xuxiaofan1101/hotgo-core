package agent

import (
	"context"
	"encoding/json"
	"testing"

	"vogo-agent/internal/protocol"
)

type fakeRunner struct {
	configs []protocol.TaskConfigPayload
	stopped []int64
}

func (r *fakeRunner) ApplyTaskConfig(ctx context.Context, config protocol.TaskConfigPayload) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.configs = append(r.configs, config)
	return nil
}

func (r *fakeRunner) Metrics() protocol.RuntimeMetrics {
	return protocol.RuntimeMetrics{RunningShards: len(r.configs)}
}

func (r *fakeRunner) StopTask(ctx context.Context, taskId int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.stopped = append(r.stopped, taskId)
	return nil
}

func TestClientHandleTaskConfigEnvelopeAppliesRunnerConfig(t *testing.T) {
	runner := &fakeRunner{}
	client := NewClientWithRunner(Config{AgentId: "agent-a", Slots: 2}, runner)
	payload := protocol.TaskConfigPayload{
		TenantId:      "default",
		TaskId:        2001,
		ConfigVersion: 3,
		Integration:   protocol.IntegrationConfig{SourceType: "kafka"},
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if err := client.handleEnvelope(context.Background(), protocol.Envelope{
		Type:    protocol.MessageTaskConfig,
		Payload: raw,
	}); err != nil {
		t.Fatalf("handle envelope: %v", err)
	}

	if len(runner.configs) != 1 {
		t.Fatalf("expected one runner config, got %#v", runner.configs)
	}
	if runner.configs[0].TaskId != 2001 || runner.configs[0].Integration.SourceType != "kafka" {
		t.Fatalf("unexpected runner config: %#v", runner.configs[0])
	}
}

func TestClientHandleTaskStopEnvelopeStopsRunnerTask(t *testing.T) {
	runner := &fakeRunner{}
	client := NewClientWithRunner(Config{AgentId: "agent-a", Slots: 2}, runner)
	payload := protocol.TaskStopPayload{TenantId: "default", TaskId: 2001}
	raw, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("marshal payload: %v", err)
	}

	if err := client.handleEnvelope(context.Background(), protocol.Envelope{
		Type:    protocol.MessageTaskStop,
		Payload: raw,
	}); err != nil {
		t.Fatalf("handle envelope: %v", err)
	}

	if len(runner.stopped) != 1 || runner.stopped[0] != 2001 {
		t.Fatalf("expected stopped task, got %#v", runner.stopped)
	}
}
