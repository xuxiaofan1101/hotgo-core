package protocol

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"
	"time"
)

const (
	MessageAgentHello     = "agent.hello"
	MessageAgentHeartbeat = "agent.heartbeat"
	MessageFieldsReport   = "agent.fields.report"
	MessageStatsReport    = "agent.stats.report"
	MessageAgentError     = "agent.error"
	MessageTaskConfig     = "server.task.config"
	MessageTaskStop       = "server.task.stop"
)

type Envelope struct {
	Type      string          `json:"type"`
	RequestId string          `json:"requestId,omitempty"`
	SentAt    time.Time       `json:"sentAt,omitempty"`
	Payload   json.RawMessage `json:"payload,omitempty"`
}

type SlotStatus struct {
	Total     int `json:"total"`
	Available int `json:"available"`
}

type RuntimeMetrics struct {
	RunningShards int     `json:"runningShards"`
	ReadTps       float64 `json:"readTps"`
	CleanTps      float64 `json:"cleanTps"`
	OutputTps     float64 `json:"outputTps"`
}

type HelloPayload struct {
	AgentId      string     `json:"agentId"`
	Hostname     string     `json:"hostname"`
	Capabilities []string   `json:"capabilities"`
	AgentIp      []string   `json:"agentIp"`
	Slots        SlotStatus `json:"slots"`
	Version      string     `json:"version"`
}

type HeartbeatPayload struct {
	AgentId string         `json:"agentId"`
	AgentIp []string       `json:"agentIp"`
	Slots   SlotStatus     `json:"slots"`
	Metrics RuntimeMetrics `json:"metrics"`
	Time    time.Time      `json:"time"`
}

type TaskConfigPayload struct {
	TenantId      string            `json:"tenantId"`
	TaskId        int64             `json:"taskId"`
	ConfigVersion int64             `json:"configVersion"`
	Shards        []TaskShard       `json:"shards"`
	Integration   IntegrationConfig `json:"integration"`
	Clean         CleanConfig       `json:"clean"`
	Outputs       []OutputConfig    `json:"outputs"`
	FieldSnapshot FieldSnapshot     `json:"fieldSnapshot,omitempty"`
	Parallelism   ParallelismConfig `json:"parallelism"`
}

type TaskStopPayload struct {
	TenantId      string `json:"tenantId"`
	TaskId        int64  `json:"taskId"`
	ConfigVersion int64  `json:"configVersion,omitempty"`
}

type TaskShard struct {
	ShardKey string         `json:"shardKey"`
	Config   map[string]any `json:"config,omitempty"`
}

type IntegrationConfig struct {
	SourceType string         `json:"sourceType"`
	Config     map[string]any `json:"config"`
}

type CleanConfig struct {
	Enabled    bool             `json:"enabled"`
	Filter     map[string]any   `json:"filter,omitempty"`
	ParseDepth int              `json:"parseDepth,omitempty"`
	Steps      []map[string]any `json:"steps"`
}

type OutputConfig struct {
	Type   string         `json:"type"`
	Config map[string]any `json:"config"`
}

type ParallelismConfig struct {
	ReaderConcurrency int `json:"readerConcurrency"`
	CleanWorkers      int `json:"cleanWorkers"`
	OutputWorkers     int `json:"outputWorkers"`
	MaxInFlight       int `json:"maxInFlight"`
	BatchSize         int `json:"batchSize"`
}

type FieldSnapshot struct {
	TenantId         string   `json:"tenantId,omitempty"`
	TaskId           int64    `json:"taskId,omitempty"`
	Version          int64    `json:"version"`
	KnownFieldHashes []string `json:"knownFieldHashes"`
	MaxParseDepth    int      `json:"maxParseDepth,omitempty"`
}

type FieldSample struct {
	Path      string `json:"path"`
	Type      string `json:"type"`
	Sample    string `json:"sample"`
	Hash      string `json:"hash"`
	Count     int64  `json:"count"`
	NullCount int64  `json:"nullCount"`
}

type FieldsReportPayload struct {
	TenantId string        `json:"tenantId"`
	TaskId   int64         `json:"taskId"`
	AgentId  string        `json:"agentId"`
	Fields   []FieldSample `json:"fields"`
}

type StatsReportPayload struct {
	AgentId     string        `json:"agentId"`
	WindowStart time.Time     `json:"windowStart"`
	WindowEnd   time.Time     `json:"windowEnd"`
	Metrics     []StatsMetric `json:"metrics"`
}

type StatsMetric struct {
	TenantId string `json:"tenantId"`
	TaskId   int64  `json:"taskId"`
	Name     string `json:"name"`
	GroupKey string `json:"groupKey,omitempty"`
	Value    int64  `json:"value"`
}

type ErrorPayload struct {
	AgentId string `json:"agentId"`
	TaskId  int64  `json:"taskId,omitempty"`
	Message string `json:"message"`
}

func FieldHash(path, fieldType string) string {
	sum := sha256.Sum256([]byte(strings.TrimSpace(path) + "\x00" + strings.TrimSpace(fieldType)))
	return hex.EncodeToString(sum[:])
}

func NewHelloEnvelope(agentId, hostname string, capabilities []string, agentIp []string, slots SlotStatus, version string) Envelope {
	return Envelope{
		Type:   MessageAgentHello,
		SentAt: time.Now(),
		Payload: mustPayload(HelloPayload{
			AgentId:      agentId,
			Hostname:     hostname,
			Capabilities: append([]string{}, capabilities...),
			AgentIp:      append([]string{}, agentIp...),
			Slots:        slots,
			Version:      version,
		}),
	}
}

func NewHeartbeatEnvelope(agentId string, slots SlotStatus, metrics RuntimeMetrics, agentIp []string, now time.Time) Envelope {
	return Envelope{
		Type:   MessageAgentHeartbeat,
		SentAt: now,
		Payload: mustPayload(HeartbeatPayload{
			AgentId: agentId,
			AgentIp: append([]string{}, agentIp...),
			Slots:   slots,
			Metrics: metrics,
			Time:    now,
		}),
	}
}

func NewFieldsReportEnvelope(payload FieldsReportPayload) Envelope {
	return Envelope{
		Type:    MessageFieldsReport,
		SentAt:  time.Now(),
		Payload: mustPayload(payload),
	}
}

func NewStatsReportEnvelope(payload StatsReportPayload) Envelope {
	return Envelope{
		Type:    MessageStatsReport,
		SentAt:  time.Now(),
		Payload: mustPayload(payload),
	}
}

func NewErrorEnvelope(agentId string, taskId int64, message string) Envelope {
	return Envelope{
		Type:   MessageAgentError,
		SentAt: time.Now(),
		Payload: mustPayload(ErrorPayload{
			AgentId: agentId,
			TaskId:  taskId,
			Message: message,
		}),
	}
}

func DecodePayload[T any](envelope Envelope, target *T) error {
	return json.Unmarshal(envelope.Payload, target)
}

func mustPayload(payload any) json.RawMessage {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	return data
}
