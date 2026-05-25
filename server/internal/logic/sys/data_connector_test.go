package sys

import (
	"crypto/tls"
	"strings"
	"testing"

	"github.com/IBM/sarama"
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

func TestDataConnectorKafkaBrokersSupportsCommaString(t *testing.T) {
	brokers := dataConnectorKafkaStringSlice(map[string]interface{}{
		"brokers": "127.0.0.1:9092, kafka:9092",
	}, "brokers")

	if len(brokers) != 2 {
		t.Fatalf("expected two brokers, got %#v", brokers)
	}
	if brokers[0] != "127.0.0.1:9092" || brokers[1] != "kafka:9092" {
		t.Fatalf("unexpected brokers: %#v", brokers)
	}
}

func TestDataConnectorKafkaInternalTopic(t *testing.T) {
	if !isDataConnectorKafkaInternalTopic("__consumer_offsets") {
		t.Fatal("expected __consumer_offsets to be treated as internal topic")
	}
	if !isDataConnectorKafkaInternalTopic(" __transaction_state ") {
		t.Fatal("expected trimmed __transaction_state to be treated as internal topic")
	}
	if isDataConnectorKafkaInternalTopic("hotgo-test-events") {
		t.Fatal("expected user topic to be kept")
	}
}

func TestBuildDataConnectorKafkaConfigUsesSASLAndTLS(t *testing.T) {
	config, err := buildDataConnectorKafkaConfig(map[string]interface{}{
		"authMode":           "plain",
		"insecureSkipVerify": true,
		"password":           "hotgo-secret",
		"tls":                true,
		"username":           "hotgo",
	})
	if err != nil {
		t.Fatalf("expected kafka config to be valid, got %v", err)
	}

	if config.Net.SASL.Enable != true {
		t.Fatal("expected SASL to be enabled")
	}
	if config.Net.SASL.Mechanism != sarama.SASLTypePlaintext {
		t.Fatalf("expected SASL/PLAIN mechanism, got %s", config.Net.SASL.Mechanism)
	}
	if config.Net.SASL.User != "hotgo" || config.Net.SASL.Password != "hotgo-secret" {
		t.Fatalf("unexpected SASL credentials: %s/%s", config.Net.SASL.User, config.Net.SASL.Password)
	}
	if config.Net.TLS.Enable != true || config.Net.TLS.Config == nil {
		t.Fatal("expected TLS config to be enabled")
	}
	if config.Net.TLS.Config.MinVersion != tls.VersionTLS12 {
		t.Fatalf("expected TLS 1.2 minimum, got %d", config.Net.TLS.Config.MinVersion)
	}
	if config.Net.TLS.Config.InsecureSkipVerify != true {
		t.Fatal("expected insecureSkipVerify to be carried from connector config")
	}
}
