package runtime

import (
	"testing"

	"vogo-agent/internal/protocol"
)

func TestBuildWorkerPlanUsesConfiguredParallelism(t *testing.T) {
	plan := BuildWorkerPlan(protocol.ParallelismConfig{
		ReaderConcurrency: 4,
		CleanWorkers:      16,
		OutputWorkers:     8,
		MaxInFlight:       5000,
		BatchSize:         200,
	})

	if plan.ReaderConcurrency != 4 || plan.CleanWorkers != 16 || plan.OutputWorkers != 8 {
		t.Fatalf("unexpected workers: %+v", plan)
	}
	if plan.MaxInFlight != 5000 || plan.BatchSize != 200 {
		t.Fatalf("unexpected flow control: %+v", plan)
	}
}

func TestBuildWorkerPlanHasSafeDefaults(t *testing.T) {
	plan := BuildWorkerPlan(protocol.ParallelismConfig{})

	if plan.ReaderConcurrency != 1 {
		t.Fatalf("unexpected reader concurrency: %d", plan.ReaderConcurrency)
	}
	if plan.CleanWorkers <= 0 || plan.OutputWorkers <= 0 {
		t.Fatalf("workers should be positive: %+v", plan)
	}
	if plan.MaxInFlight < 100 || plan.BatchSize < 1 {
		t.Fatalf("flow control should be bounded: %+v", plan)
	}
}
