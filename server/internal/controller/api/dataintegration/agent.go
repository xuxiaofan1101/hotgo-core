package dataintegration

import (
	"encoding/json"
	"time"

	"hotgo/internal/model/input/sysin"
	"hotgo/internal/service"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

var Agent = cAgent{}

type cAgent struct{}

func (c *cAgent) ServeWebSocket(r *ghttp.Request) {
	ctx := r.Context()
	ws, err := r.WebSocket()
	if err != nil {
		g.Log().Warningf(ctx, "Agent WebSocket upgrade失败: %v", err)
		r.Exit()
		return
	}
	defer ws.Close()
	g.Log().Info(ctx, "Agent WebSocket connected")

	for {
		var envelope sysin.DataAgentEnvelope
		if err = ws.ReadJSON(&envelope); err != nil {
			g.Log().Infof(ctx, "Agent WebSocket closed: %v", err)
			return
		}
		if envelope.Type == sysin.DataAgentMessageHello {
			g.Log().Infof(ctx, "Agent WebSocket received hello requestId:%s", envelope.RequestId)
		}
		response, err := service.SysDataAgent().HandleAgentEnvelope(ctx, envelope)
		if err != nil {
			g.Log().Warningf(ctx, "Agent WebSocket处理失败 type:%s requestId:%s err:%+v", envelope.Type, envelope.RequestId, err)
			response = dataAgentErrorEnvelope(envelope.RequestId, err.Error())
		}
		if response.Type == "" {
			continue
		}
		if err = ws.WriteJSON(response); err != nil {
			g.Log().Warningf(ctx, "Agent WebSocket写入失败: %v", err)
			return
		}
	}
}

func dataAgentErrorEnvelope(requestId string, message string) sysin.DataAgentEnvelope {
	payload, _ := json.Marshal(map[string]any{"message": message})
	return sysin.DataAgentEnvelope{
		Type:      sysin.DataAgentMessageServerError,
		RequestId: requestId,
		SentAt:    time.Now(),
		Payload:   payload,
	}
}
