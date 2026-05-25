package sysin

import (
	"context"
	"testing"
)

func TestDataCleanTaskSampleInpFilterDefaultsAndCapsTimeout(t *testing.T) {
	in := &DataCleanTaskSampleInp{SourceId: 1}
	if err := in.Filter(context.Background()); err != nil {
		t.Fatalf("expected default sample input to pass, got %v", err)
	}
	if in.Limit != 50 {
		t.Fatalf("expected default limit 50, got %d", in.Limit)
	}
	if in.TimeoutSeconds != 10 {
		t.Fatalf("expected default timeout 10 seconds, got %d", in.TimeoutSeconds)
	}

	in = &DataCleanTaskSampleInp{SourceId: 1, Limit: 1000, TimeoutSeconds: 120}
	if err := in.Filter(context.Background()); err != nil {
		t.Fatalf("expected capped sample input to pass, got %v", err)
	}
	if in.Limit != 500 {
		t.Fatalf("expected max limit 500, got %d", in.Limit)
	}
	if in.TimeoutSeconds != 60 {
		t.Fatalf("expected max timeout 60 seconds, got %d", in.TimeoutSeconds)
	}
}
