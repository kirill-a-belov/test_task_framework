// TCP Server general implementation
package server

import (
	"context"
	"fmt"
	"github.com/kirill-a-belov/test_task_framework/internal/app/server/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/protocol"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/tracer"
	"net"
	"sync/atomic"
	"time"

	"github.com/pkg/errors"
)

/** Нет тайм-аутов, злоумышленник легко может занять все 100 доступных соединений
* Контекст в хендлере соединения из конструктора сервера
* Использование fmt.Sprintf для построения хеша
* Текстовый протокол
* math.Rand для рандома
* нет обработки ошибок в сервер внутри процессора
* God Object из функций с логикой расчета челенжа
* Нет таймаутов*/

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
	connCnt  atomic.Int32

	listener net.Listener
}

func (s *Server) Start(ctx context.Context) error {
	_, span := tracer.Start(ctx, "internal.app.server.Server.Start")
	defer span.End()

	var err error
	if s.listener, err = net.Listen(protocol.TCPType, fmt.Sprintf("localhost:%d", s.config.Port)); err != nil {
		return errors.Wrap(err, "start listener")
	}

	go s.processor(ctx)

	return nil
}

func (s *Server) processor(ctx context.Context) {
	_, span := tracer.Start(ctx, "internal.app.server.Server.processor")
	defer span.End()

	for {
		select {
		case <-s.stopChan:
		default:
			if s.connCnt.Load() > int32(s.config.ConnPoolSize) {
				s.logger.Info("max conn pool size exceeded")
				const connReleaseWaitTime = time.Millisecond
				time.Sleep(connReleaseWaitTime)

				continue
			}

			conn, err := s.listener.Accept()
			if err != nil {
				s.logger.Error(err, "connection processing")
			}

			s.connCnt.Add(1)
			go func() {
				defer conn.Close()
				defer s.connCnt.Add(-1)
				if err := s.serv(conn); err != nil {
					s.logger.Error(err, "connection serving")
				}
			}()
		}
	}
}

func (s *Server) Stop(ctx context.Context) {
	_, span := tracer.Start(ctx, "internal.app.server.Server.Stop")
	defer span.End()

	close(s.stopChan)
}

func (s *Server) serv(conn net.Conn) error {
	ctx, span := tracer.Start(context.Background(), "internal.app.server.Server.serv")
	defer span.End()

	ctx, _ = context.WithDeadline(ctx, time.Now().Add(s.config.ConnTTL))

	/*/clientWelcome := make([]byte, 1024)
	n, err := conn.Read(clientWelcome)
	if err != nil {
		return errors.Wrap(err, "reading client welcome")
	}
	clientWelcome = clientWelcome[:n]
	clientWelcomeRequest := &protocol.ClientWelcomeRequest{}
	if err := json.Unmarshal(clientWelcome, clientWelcomeRequest); err != nil {
		return errors.Wrap(err, "unmarshalling client welcome request")
	}
	if clientWelcomeRequest.Type != protocol.MessageTypeClientWelcome {
		return errors.Errorf("client welcome request: received wrong message (%s)", clientWelcomeRequest)
	}

	sendingPrefix := rand.Int63()
	serverQuestionResponse, err := json.Marshal(protocol.ServerQuestionRequest{
		Message: protocol.Message{
			Type: protocol.MessageTypeServerQuestion,
		},
		Prefix:     sendingPrefix,
		Difficulty: s.config.difficulty,
	})
	if err != nil {
		return errors.Wrap(err, "marshalling server question response")
	}
	if _, err := conn.Write(serverQuestionResponse); err != nil {
		return errors.Wrap(err, "sending server question response")
	}

	clientAnswer := make([]byte, 1024)
	n, err = conn.Read(clientAnswer)
	if err != nil {
		return errors.Wrap(err, "reading client welcome")
	}
	clientAnswer = clientAnswer[:n]
	clientAnswerResponse := &protocol.ClientAnswerResponse{}
	if err := json.Unmarshal(clientAnswer, clientAnswerResponse); err != nil {
		return errors.Wrap(err, "unmarshalling client welcome response")
	}
	if clientAnswerResponse.Type != protocol.MessageTypeClientAnswer {
		return errors.Errorf("client welcome response: received wrong message (%s)", clientAnswerResponse)
	}

	resp := protocol.ServerResultResponse{
		Message: protocol.Message{
			Type: protocol.MessageTypeServerResult,
		},
	}

	var ok bool

	fakeResponse := clientAnswerResponse.Prefix != sendingPrefix || clientAnswerResponse.Difficulty != s.config.difficulty
	if !fakeResponse {
		ok = sha256.Check(
			ctx,
			clientAnswerResponse.Nonce,
			clientAnswerResponse.Prefix,
			clientAnswerResponse.Difficulty,
		)
	}

	switch {
	case ok:
		resp.Success = true
		resp.Payload = s.quoteList[rand.Intn(2)]
	default:
		resp.Payload = "invalid result"
	}

	serverResultResponse, err := json.Marshal(resp)
	if err != nil {
		return errors.Wrap(err, "marshalling server result response")
	}
	if _, err := conn.Write(serverResultResponse); err != nil {
		return errors.Wrap(err, "sending server result response")
	}
	*/
	return nil
}
