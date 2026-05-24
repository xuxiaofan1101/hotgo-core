package agent

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigFromFileReadsDirectWebSocketConfig(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`
agent:
  id: agent-a
  serverUrl: ws://server.example.com/api/v1/data-integration/agent/ws
  hostname: worker-a
  capabilities:
    - kafka
    - s3
  slots: 8
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := LoadConfigFromFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if config.AgentId != "agent-a" || config.ServerURL != "ws://server.example.com/api/v1/data-integration/agent/ws" {
		t.Fatalf("unexpected config identity: %#v", config)
	}
	if config.Hostname != "worker-a" {
		t.Fatalf("unexpected hostname: %#v", config)
	}
	if config.Slots != 8 || len(config.Capabilities) != 2 || config.Capabilities[0] != "kafka" || config.Capabilities[1] != "s3" {
		t.Fatalf("unexpected config capacity: %#v", config)
	}
}

func TestLoadConfigFromFileDerivesStableAgentIdWhenConfigIdIsEmpty(t *testing.T) {
	oldMachineIdReader := readStableMachineId
	oldMacReader := readPhysicalMacAddresses
	readStableMachineId = func() string { return "machine-abc" }
	readPhysicalMacAddresses = func() []string { return []string{"02:42:ac:11:00:02", "00:16:3e:00:00:01"} }
	t.Cleanup(func() {
		readStableMachineId = oldMachineIdReader
		readPhysicalMacAddresses = oldMacReader
	})

	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	content := []byte(`
agent:
  id: ""
  serverUrl: ws://server.example.com/api/v1/data-integration/agent/ws
  hostname: worker-a
`)
	if err := os.WriteFile(path, content, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	config, err := LoadConfigFromFile(path)
	if err != nil {
		t.Fatalf("load config: %v", err)
	}
	if config.AgentId == "" || config.AgentId == "worker-a" {
		t.Fatalf("agent id should be derived from stable machine fingerprint, got %#v", config.AgentId)
	}

	first := deriveStableAgentId("machine-abc", []string{"02:42:ac:11:00:02", "00:16:3e:00:00:01"})
	second := deriveStableAgentId("machine-abc", []string{"00:16:3e:00:00:01", "02:42:ac:11:00:02"})
	if first == "" || first != second || config.AgentId != first {
		t.Fatalf("agent id should be stable across mac order, config=%q first=%q second=%q", config.AgentId, first, second)
	}
}

func TestAgentVersionCanBeInjectedAtBuildTime(t *testing.T) {
	oldVersion := Version
	Version = "9.9.9-test"
	t.Cleanup(func() { Version = oldVersion })

	if Version != "9.9.9-test" {
		t.Fatalf("agent version should be mutable for ldflags injection")
	}
}
