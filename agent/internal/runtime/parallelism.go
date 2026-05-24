package runtime

import "vogo-agent/internal/protocol"

type WorkerPlan struct {
	ReaderConcurrency int
	CleanWorkers      int
	OutputWorkers     int
	MaxInFlight       int
	BatchSize         int
}

func BuildWorkerPlan(config protocol.ParallelismConfig) WorkerPlan {
	return WorkerPlan{
		ReaderConcurrency: positiveOrDefault(config.ReaderConcurrency, 1),
		CleanWorkers:      positiveOrDefault(config.CleanWorkers, 4),
		OutputWorkers:     positiveOrDefault(config.OutputWorkers, 2),
		MaxInFlight:       positiveOrDefault(config.MaxInFlight, 1000),
		BatchSize:         positiveOrDefault(config.BatchSize, 100),
	}
}

func positiveOrDefault(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}
