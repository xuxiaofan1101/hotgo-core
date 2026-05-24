package consts

import (
	"testing"

	"hotgo/internal/library/dict"
)

func TestDataIntegrationConfigDictionariesAreRegistered(t *testing.T) {
	requiredTypes := []string{
		"DataConnectorSourceTypeOptions",
		"DataConnectorSinkTypeOptions",
		"DataKafkaAuthModeOptions",
		"DataS3CredentialModeOptions",
		"DataLogLevelOptions",
		"DataKafkaStartModeOptions",
		"DataSampleModeOptions",
		"DataPayloadFormatOptions",
		"DataSinkModeOptions",
		"DataSinkFailurePolicyOptions",
	}

	for _, typ := range requiredTypes {
		if options := dict.GetEnumsOptions(typ); len(options) == 0 {
			t.Fatalf("dictionary %s should be registered with options", typ)
		}
	}
}
