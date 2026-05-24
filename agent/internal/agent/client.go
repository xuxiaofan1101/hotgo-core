package agent

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"vogo-agent/internal/protocol"
	"vogo-agent/internal/runtime"
)

type TaskConfigRunner interface {
	ApplyTaskConfig(ctx context.Context, config protocol.TaskConfigPayload) error
	StopTask(ctx context.Context, taskId int64) error
	Metrics() protocol.RuntimeMetrics
}

type statsDrainer interface {
	DrainStats() protocol.StatsReportPayload
}

type Client struct {
	config Config
	dialer *websocket.Dialer
	runner TaskConfigRunner
}

func NewClient(config Config) *Client {
	return NewClientWithRunner(config, runtime.NewTaskRunner(config.AgentId, runtime.NewDefaultConnectorRegistry()))
}

func NewClientWithRunner(config Config, runner TaskConfigRunner) *Client {
	if runner == nil {
		runner = runtime.NewTaskRunner(config.AgentId, runtime.NewDefaultConnectorRegistry())
	}
	return &Client{
		config: config,
		dialer: websocket.DefaultDialer,
		runner: runner,
	}
}

func (c *Client) Run(ctx context.Context) error {
	conn, _, err := c.dialer.DialContext(ctx, c.config.ServerURL, nil)
	if err != nil {
		return fmt.Errorf("连接 rule server 失败: %w", err)
	}
	defer conn.Close()
	var writeMu sync.Mutex
	writeJSON := func(payload any) error {
		writeMu.Lock()
		defer writeMu.Unlock()
		return conn.WriteJSON(payload)
	}
	if taskRunner, ok := c.runner.(*runtime.TaskRunner); ok {
		taskRunner.SetFieldReporter(func(report protocol.FieldsReportPayload) {
			_ = writeJSON(protocol.NewFieldsReportEnvelope(report))
		})
	}

	slots := protocol.SlotStatus{Total: c.config.Slots, Available: c.config.Slots}
	if err := writeJSON(protocol.NewHelloEnvelope(
		c.config.AgentId,
		c.config.Hostname,
		c.config.Capabilities,
		localAgentIPs(),
		slots,
		Version,
	)); err != nil {
		return fmt.Errorf("发送 agent hello 失败: %w", err)
	}

	readErr := make(chan error, 1)
	go func() {
		for {
			var envelope protocol.Envelope
			if err := conn.ReadJSON(&envelope); err != nil {
				readErr <- err
				return
			}
			if err := c.handleEnvelope(ctx, envelope); err != nil {
				readErr <- err
				return
			}
		}
	}()

	heartbeatTicker := time.NewTicker(10 * time.Second)
	defer heartbeatTicker.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-readErr:
			return fmt.Errorf("读取 server 消息失败: %w", err)
		case <-heartbeatTicker.C:
			if drainer, ok := c.runner.(statsDrainer); ok {
				stats := drainer.DrainStats()
				if len(stats.Metrics) > 0 {
					if err := writeJSON(protocol.NewStatsReportEnvelope(stats)); err != nil {
						return fmt.Errorf("发送 agent stats 失败: %w", err)
					}
				}
			}
			if err := writeJSON(protocol.NewHeartbeatEnvelope(
				c.config.AgentId,
				slots,
				c.runner.Metrics(),
				localAgentIPs(),
				time.Now(),
			)); err != nil {
				return fmt.Errorf("发送 agent heartbeat 失败: %w", err)
			}
		}
	}
}

func (c *Client) handleEnvelope(ctx context.Context, envelope protocol.Envelope) error {
	switch envelope.Type {
	case protocol.MessageTaskConfig:
		var payload protocol.TaskConfigPayload
		if err := protocol.DecodePayload(envelope, &payload); err != nil {
			return fmt.Errorf("解析任务配置失败: %w", err)
		}
		if err := c.runner.ApplyTaskConfig(ctx, payload); err != nil {
			return fmt.Errorf("应用任务配置失败: %w", err)
		}
	case protocol.MessageTaskStop:
		var payload protocol.TaskStopPayload
		if err := protocol.DecodePayload(envelope, &payload); err != nil {
			return fmt.Errorf("解析任务停止配置失败: %w", err)
		}
		if err := c.runner.StopTask(ctx, payload.TaskId); err != nil {
			return fmt.Errorf("停止任务失败: %w", err)
		}
	}
	return nil
}
