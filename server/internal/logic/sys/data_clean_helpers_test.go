package sys

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func TestValidateDataCleanConfigRejectsInvalidStep(t *testing.T) {
	err := validateDataCleanConfig(gjson.New(map[string]interface{}{
		"steps": []map[string]interface{}{
			{"type": "unsupported_action", "field": "payload.user.id"},
		},
	}))
	if err == nil {
		t.Fatal("expected invalid clean step to be rejected")
	}
	if !strings.Contains(err.Error(), "清洗步骤类型无效") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestNormalizeDataSinkConfigDefaultsStrategyEngine(t *testing.T) {
	config := normalizeDataSinkConfig(nil)
	targets := config.Get("targets").Array()
	if len(targets) != 1 {
		t.Fatalf("expected one default target, got %d", len(targets))
	}

	target := gjson.New(targets[0])
	if target.Get("type").String() != "strategy_engine" {
		t.Fatalf("unexpected default target type: %s", target.Get("type").String())
	}
	if target.Get("enabled").Bool() != true {
		t.Fatal("expected default target to be enabled")
	}
}

func TestProfileDataFieldsExtractsNestedPaths(t *testing.T) {
	fields := profileDataFields([]map[string]interface{}{
		{
			"payload": map[string]interface{}{
				"user": map[string]interface{}{
					"id": "1001",
				},
				"risk": map[string]interface{}{
					"score": float64(92),
				},
			},
		},
	}, 4)

	fieldMap := make(map[string]dataFieldProfile)
	for _, field := range fields {
		fieldMap[field.FieldPath] = field
	}

	if fieldMap["payload.user.id"].FieldType != "string" {
		t.Fatalf("expected payload.user.id string field, got %#v", fieldMap["payload.user.id"])
	}
	if fieldMap["payload.risk.score"].FieldType != "number" {
		t.Fatalf("expected payload.risk.score number field, got %#v", fieldMap["payload.risk.score"])
	}
}

func TestLimitDataCleanSamplePayloadsLimitsRowsBeforeFieldMerge(t *testing.T) {
	payloads := []map[string]interface{}{
		{"alpha": "a", "beta": "b", "gamma": "c"},
		{"delta": "d"},
	}

	fields := profileDataFields(limitDataCleanSamplePayloads(payloads, 1), 4)

	fieldMap := make(map[string]dataFieldProfile)
	for _, field := range fields {
		fieldMap[field.FieldPath] = field
	}

	if len(fieldMap) != 3 {
		t.Fatalf("expected all fields from first sampled row only, got %#v", fieldMap)
	}
	if _, ok := fieldMap["delta"]; ok {
		t.Fatalf("expected second row to be excluded by sample count, got %#v", fieldMap["delta"])
	}
}

func TestAppendDataCleanSamplePayloadsStopsAtLimit(t *testing.T) {
	payloads := appendDataCleanSamplePayloads(
		[]map[string]interface{}{{"first": true}},
		[]map[string]interface{}{{"second": true}, {"third": true}},
		2,
	)

	if len(payloads) != 2 {
		t.Fatalf("expected two sampled rows, got %#v", payloads)
	}
	if _, ok := payloads[1]["second"]; !ok {
		t.Fatalf("expected second row to be kept before limit, got %#v", payloads[1])
	}
}
