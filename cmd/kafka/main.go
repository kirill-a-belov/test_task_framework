package main

import (
	"context"

	"github.com/kirill-a-belov/test_task_framework/internal/app/cmd/kafka"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
)

func main() {
	log := logger.New("cmd.kafka")
	cmd := kafka.New(context.Background())

	if err := cmd.Execute(); err != nil {
		log.Error(err, "running kafka services set")
	}
}
