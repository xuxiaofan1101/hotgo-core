package agent

import (
	"fmt"
	"os"
	"strings"

	"go.yaml.in/yaml/v3"
)

var Version = "0.1.0"

const DefaultConfigPath = "manifest/config/config.yaml"

type Config struct {
	AgentId      string
	ServerURL    string
	Hostname     string
	Capabilities []string
	Slots        int
}

type fileConfig struct {
	Agent struct {
		Id           string   `yaml:"id"`
		ServerURL    string   `yaml:"serverUrl"`
		Hostname     string   `yaml:"hostname"`
		Capabilities []string `yaml:"capabilities"`
		Slots        int      `yaml:"slots"`
	} `yaml:"agent"`
}

func LoadConfigFromFile(path string) (Config, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("读取 agent 配置文件失败: %w", err)
	}
	var raw fileConfig
	if err := yaml.Unmarshal(content, &raw); err != nil {
		return Config{}, fmt.Errorf("解析 agent 配置文件失败: %w", err)
	}

	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	if strings.TrimSpace(raw.Agent.Hostname) != "" {
		hostname = strings.TrimSpace(raw.Agent.Hostname)
	}
	agentId := strings.TrimSpace(raw.Agent.Id)
	if agentId == "" {
		agentId = deriveStableAgentId(readStableMachineId(), readPhysicalMacAddresses())
	}
	if agentId == "" {
		agentId = hostname
	}
	config := Config{
		AgentId:      agentId,
		ServerURL:    strings.TrimSpace(raw.Agent.ServerURL),
		Hostname:     hostname,
		Capabilities: normalizeList(raw.Agent.Capabilities, []string{"kafka", "s3", "http", "log"}),
		Slots:        raw.Agent.Slots,
	}
	if config.ServerURL == "" {
		config.ServerURL = "ws://127.0.0.1:8888/api/v1/data-integration/agent/ws"
	}
	if config.Slots <= 0 {
		config.Slots = 4
	}
	return config, nil
}

func normalizeList(items []string, fallback []string) []string {
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	if len(result) == 0 {
		return append([]string{}, fallback...)
	}
	return result
}
