// Package server implements general TCP server
package server

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/kirill-a-belov/test_task_framework/internal/app/server/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/reservation"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
)

func New(ctx context.Context, config *config.Config) *Server {
	_, span := tracer.Start(ctx, "internal.app.server.New")
	defer span.End()

	return &Server{
		config:   config,
		stopChan: make(chan struct{}),
		logger:   logger.New("server"),
	}
}

type Server struct {
	config   *config.Config
	stopChan chan struct{}
	logger   logger.Logger
}

func (s *Server) Start(ctx context.Context) error {
	_, span := tracer.Start(ctx, "internal.app.server.Server.Start")
	defer span.End()

	s.logger.Info(fmt.Sprintf("config: %+v", *s.config))

	reservationModule := reservation.New()

	r := gin.Default()
	reservationModule.RegisterRoutes(r)

	go func() {
		if err := r.Run(); err != nil {
			// TODO(KB): Graceful shutdown
			panic(err)
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) {
	_, span := tracer.Start(ctx, "internal.app.server.Server.Stop")
	defer span.End()

	close(s.stopChan)
}
