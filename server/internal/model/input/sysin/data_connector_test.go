package sysin

import (
	"hotgo/internal/model/entity"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gctx"
)

func TestDataConnectorEditInpFilter(t *testing.T) {
	ctx := gctx.New()

	valid := &DataConnectorEditInp{
		DataConnector: entity.DataConnector{
			Name:          "Kafka 输入",
			Code:          "kafka_source",
			Direction:     "source",
			ConnectorType: "kafka",
			Status:        1,
		},
	}
	if err := valid.Filter(ctx); err != nil {
		t.Fatalf("valid connector should pass filter: %v", err)
	}

	invalidDirection := &DataConnectorEditInp{
		DataConnector: entity.DataConnector{
			Name:          "错误方向",
			Code:          "bad_direction",
			Direction:     "unknown",
			ConnectorType: "kafka",
			Status:        1,
		},
	}
	if err := invalidDirection.Filter(ctx); err == nil || !strings.Contains(err.Error(), "连接方向不正确") {
		t.Fatalf("expected invalid direction error, got: %v", err)
	}

	invalidType := &DataConnectorEditInp{
		DataConnector: entity.DataConnector{
			Name:          "错误类型",
			Code:          "bad_type",
			Direction:     "source",
			ConnectorType: "unknown",
			Status:        1,
		},
	}
	if err := invalidType.Filter(ctx); err == nil || !strings.Contains(err.Error(), "连接类型不正确") {
		t.Fatalf("expected invalid connector type error, got: %v", err)
	}
}
