// TCP Server general implementation
package server

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/kirill-a-belov/temp_test_task/utils/sha256"
	"github.com/kirill-a-belov/temp_test_task/utils/tracing"
	"github.com/kirill-a-belov/test_task_framework/internal/app/server/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/internal/app/server/pkg/protocol"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"math/rand"
	"net"
	"sync/atomic"

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
	span, _ := tracing.NewSpan(ctx, "server.New")
	defer span.Close()

	return &Server{
		config:    config,
		stopChan:  make(chan struct{}),
		logger:    logger.New("server"),
		quoteList: []string{"example quote one", "example quote two", "example quote three"},
	}
}

type Server struct {
	config   *config.Config
	stopChan chan struct{}
	logger   logger.Logger
	connCnt  atomic.Int32

	listner   net.Listener
	quoteList []string
}

func (s *Server) Start(ctx context.Context) error {
	span, _ := tracing.NewSpan(ctx, "internal.app.server.Start")
	defer span.Close()

	var err error
	if s.listner, err = net.Listen(protocol.TCPType, fmt.Sprintf("localhost:%d", s.config.serverPort)); err != nil {
		return errors.Wrap(err, "start listener")
	}

	go s.processor(ctx)

	return nil
}

func (s *Server) processor(ctx context.Context) {
	span, _ := tracing.NewSpan(ctx, "internal.app.server.processor")
	defer span.Close()

	for {
		select {
		case <-s.stopChan:
		default:
			conn, err := s.listner.Accept()
			if err != nil {
				s.logger.Println("connection processing",
					"error", err,
				)
			}

			if s.connCnt.Load() > int32(s.config.maxConns) {
				if _, err := conn.Write([]byte("max conns exceeded")); err != nil {
					s.logger.Println("max conn response sending",
						"error", err,
					)
				}
				conn.Close()
				continue
			}

			s.connCnt.Add(1)
			go func() {
				defer conn.Close()
				defer s.connCnt.Add(-1)
				if err := s.serv(ctx, conn); err != nil {
					s.logger.Println("connection serving",
						"error", err,
					)
				}
			}()
		}
	}
}

func (s *Server) Stop(ctx context.Context) {
	span, _ := tracing.NewSpan(ctx, "internal.app.server.Stop")
	defer span.Close()

	close(s.stopChan)
}

func (s *Server) serv(ctx context.Context, conn net.Conn) error {
	span, _ := tracing.NewSpan(ctx, "internal.app.server.serv")
	defer span.Close()

	clientWelcome := make([]byte, 1024)
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

	return nil
}
