package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"vogo-agent/internal/agent"
)

func main() {
	configPath := flag.String("config", agent.DefaultConfigPath, "agent 配置文件路径")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	config, err := agent.LoadConfigFromFile(*configPath)
	if err != nil {
		log.Fatalf("加载 agent 配置失败: %v", err)
	}

	client := agent.NewClient(config)
	if err := client.Run(ctx); err != nil && !errors.Is(err, context.Canceled) {
		log.Fatalf("rule-clean-agent stopped: %v", err)
	}
}
