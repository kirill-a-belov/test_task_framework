// TCP Server general implementation
package server

import (
	"context"
	"fmt"
	"github.com/kirill-a-belov/test_task_framework/internal/app/server/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/network"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/protocol"
	"github.com/kirill-a-belov/test_task_framework/pkg/context_helper"

	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/math"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
	"io"
	"net"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

func New(ctx context.Context, config *config.Config) *Server {
	_, span := tracer.Start(ctx, "internal.app.server.New")
	defer span.End()

	return &Server{
		config:   config,
		stopChan: make(chan struct{}),
		logger:   logger.New("server"),
		listenerStarter: func() (net.Listener, error) {
			return net.Listen(protocol.TCPType, fmt.Sprintf("localhost:%d", config.Port))
		},
	}
}

type Server struct {
	config   *config.Config
	stopChan chan struct{}
	logger   logger.Logger
	connCnt  atomic.Int32

	listener        net.Listener
	listenerStarter func() (net.Listener, error)
}

func (s *Server) Start(ctx context.Context) error {
	_, span := tracer.Start(ctx, "internal.app.server.Server.Start")
	defer span.End()

	s.logger.Info(fmt.Sprintf("config: %+v", *s.config))

	var err error
	if s.listener, err = s.listenerStarter(); err != nil {
		return errors.Wrap(err, "start listener")
	}

	go s.processor(ctx, serv)

	return nil
}

func (s *Server) processor(ctx context.Context, servFunc func(io.ReadWriter) error) {
	_, span := tracer.Start(ctx, "internal.app.server.Server.processor")
	defer span.End()

	for {
		select {
		case <-s.stopChan:
			s.logger.Info("processor terminated")

			return
		default:
			if s.connCnt.Load() >= int32(s.config.ConnPoolSize) {
				s.logger.Info("max conn pool size exceeded")

				const connReleaseWaitTime = 100 * time.Millisecond
				time.Sleep(connReleaseWaitTime)

				continue
			}

			conn, err := s.listener.Accept()
			if err != nil {
				s.logger.Error(err, "connection processing")

				continue
			}
			s.connCnt.Add(1)

			go func() {

				defer conn.Close()
				defer s.connCnt.Add(-1)
				if err := context_helper.RunWithTimeout(s.config.ConnTTL, func() error {
					return servFunc(conn)
				}); err != nil {
					s.logger.Error(err, "connection serving")
				}
			}()
		}
	}
}

func (s *Server) Stop(ctx context.Context) {
	_, span := tracer.Start(ctx, "internal.app.server.Server.Stop")
	defer span.End()

	_ = s.listener.Close()
	close(s.stopChan)
}

func serv(conn io.ReadWriter) error {
	ctx, span := tracer.Start(context.Background(), "internal.app.server.Server.serv")
	defer span.End()

	request, err := network.Receive[protocol.Request](ctx, conn)
	if err != nil {
		return errors.Wrap(err, "receiving server request")
	}
	if request.Type != protocol.MessageTypeRequest {
		return errors.Errorf("server requrest: received wrong message (%v)", request)
	}

	if err := network.Send(ctx, conn, protocol.Response{
		Message: protocol.Message{
			Type: protocol.MessageTypeResponse,
		},
		Payload: math.Sum(request.Payload...),
	}); err != nil {
		return errors.Wrap(err, "sending response")
	}

	return nil
}
