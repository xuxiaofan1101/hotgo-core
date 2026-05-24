package sys

import (
	"strings"
	"testing"

	"github.com/gogf/gf/v2/encoding/gjson"
)

func TestValidateS3ConnectorRequiresRegionNotEndpoint(t *testing.T) {
	err := validateConnectorConfig("source", "s3", gjson.New(map[string]interface{}{
		"credentialMode": "role",
		"region":         "us-east-1",
	}))
	if err != nil {
		t.Fatalf("expected S3 config without endpoint to pass, got %v", err)
	}

	err = validateConnectorConfig("source", "s3", gjson.New(map[string]interface{}{
		"credentialMode": "role",
		"endpoint":       "https://s3.amazonaws.com",
	}))
	if err == nil || !strings.Contains(err.Error(), "S3 Region") {
		t.Fatalf("expected missing region error, got %v", err)
	}
}
