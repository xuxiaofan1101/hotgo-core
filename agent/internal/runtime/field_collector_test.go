package runtime

import (
	"testing"

	"vogo-agent/internal/protocol"
)

func TestFieldCollectorSkipsKnownFieldsFromServerSnapshot(t *testing.T) {
	collector := NewFieldCollector(protocol.FieldSnapshot{
		Version:          3,
		KnownFieldHashes: []string{protocol.FieldHash("payload.user.id", "string")},
	})

	samples := collector.Collect(map[string]any{
		"payload": map[string]any{
			"user": map[string]any{
				"id":   "1001",
				"name": "alice",
			},
		},
	})
	second := collector.Collect(map[string]any{
		"payload": map[string]any{
			"user": map[string]any{
				"id":   "1002",
				"name": "bob",
			},
		},
	})

	if len(samples) != 1 {
		t.Fatalf("expected only unknown field, got %+v", samples)
	}
	if samples[0].Path != "payload.user.name" || samples[0].Type != "string" {
		t.Fatalf("unexpected sample: %+v", samples[0])
	}
	if len(second) != 0 {
		t.Fatalf("expected duplicate local field to be skipped, got %+v", second)
	}
}

func TestFieldCollectorRespectsMaxParseDepth(t *testing.T) {
	collector := NewFieldCollector(protocol.FieldSnapshot{MaxParseDepth: 2})

	samples := collector.Collect(map[string]any{
		"level1": map[string]any{
			"level2": map[string]any{
				"level3": "too-deep",
			},
		},
	})

	if len(samples) != 1 {
		t.Fatalf("expected boundary object sample only, got %+v", samples)
	}
	if samples[0].Path != "level1.level2" || samples[0].Type != "object" {
		t.Fatalf("unexpected depth-limited sample: %+v", samples[0])
	}
}
