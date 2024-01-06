package server

import (
	"bytes"
	"context"
	"encoding/gob"
	"github.com/kirill-a-belov/test_task_framework/internal/app/server/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/protocol"
	"github.com/kirill-a-belov/test_task_framework/pkg/test_helper"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"io"
	"net"
	"sync"
	"testing"
	"time"
)

func TestServer_Start(t *testing.T) {
	testCaseList := []struct {
		name      string
		args      func() *listenerStarterMock
		wantError bool
	}{
		{
			name: "Success",
			args: func() *listenerStarterMock {
				lsm := &listenerStarterMock{}
				lsm.On("mockFunc").Return(&net.TCPListener{}, nil)

				return lsm
			},
			wantError: false,
		},
		{
			name: "Listen start failed",
			args: func() *listenerStarterMock {
				lsm := &listenerStarterMock{}
				lsm.On("mockFunc").Return(&net.TCPListener{}, errors.New("example error"))

				return lsm
			},
			wantError: true,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := New(ctx, &config.Config{})
			s.listenerStarter = tc.args().mockFunc

			err := s.Start(ctx)
			defer s.Stop(ctx)

			if tc.wantError {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}

type listenerStarterMock struct {
	mock.Mock
}

func (lsm *listenerStarterMock) mockFunc() (net.Listener, error) {
	args := lsm.Called()

	return args.Get(0).(net.Listener), args.Error(1)
}

func TestServer_processor(t *testing.T) {
	testCaseList := []struct {
		name string
		args func() (*Server, func(io.ReadWriter) error)
	}{
		{
			name: "Regular stop",
			args: func() (*Server, func(io.ReadWriter) error) {
				f := func(writer io.ReadWriter) error {
					return nil
				}

				loggerMock := &test_helper.LoggerMock{}
				loggerMock.On("Info", mock.Anything)

				s := New(context.Background(), &config.Config{
					ConnTTL: time.Second,
				})
				s.logger = loggerMock
				s.listener = &listenerStub{}

				return s, f
			},
		},
		{
			name: "Conn limit exceeded",
			args: func() (*Server, func(io.ReadWriter) error) {
				f := func(writer io.ReadWriter) error {
					for {
					}
				}

				loggerMock := &test_helper.LoggerMock{}
				loggerMock.On("Info", mock.Anything).Twice()

				s := New(context.Background(), &config.Config{
					ConnTTL:      time.Second,
					ConnPoolSize: 1,
				})
				s.logger = loggerMock
				s.listener = &listenerStub{}

				return s, f
			},
		},
		{
			name: "Serv func error",
			args: func() (*Server, func(io.ReadWriter) error) {
				f := func(writer io.ReadWriter) error {
					return errors.New("example error")
				}

				loggerMock := &test_helper.LoggerMock{}
				loggerMock.On("Error", mock.Anything).Once()
				loggerMock.On("Info", mock.Anything)

				s := New(context.Background(), &config.Config{
					ConnTTL: time.Second,
				})
				s.logger = loggerMock
				s.listener = &listenerStub{}

				return s, f
			},
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			testServer, testFunc := tc.args()

			ctx := context.Background()

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				testServer.processor(ctx, testFunc)
				wg.Done()
			}()

			time.Sleep(time.Millisecond)
			testServer.Stop(ctx)
			wg.Wait()
		})
	}
}

func TestServer_Stop(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		s := New(ctx, &config.Config{})
		s.listener = &net.TCPListener{}

		s.Stop(ctx)

		_, closed := <-s.stopChan
		assert.False(t, closed)

	})
}

func Test_serv(t *testing.T) {
	testCaseList := []struct {
		name      string
		args      func() io.ReadWriter
		wantError bool
	}{
		{
			name: "Success",
			args: func() io.ReadWriter {
				request := protocol.Request{
					Message: protocol.Message{
						Type: protocol.MessageTypeRequest,
					},
					Payload: []int{1, 2, 3},
				}
				requestBuffer := &bytes.Buffer{}
				require.NoError(t, gob.NewEncoder(requestBuffer).Encode(request))

				rwm := &test_helper.ReadWriterMock{}
				rwm.On("Read", mock.Anything).
					Run(func(args mock.Arguments) {
						b := args[0].([]byte)
						copy(b[:], requestBuffer.Bytes())
					}).
					Return(1000, nil)

				rwm.On("Write", []byte{0x2f, 0xff, 0x85, 0x3, 0x1, 0x1, 0x8, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x1, 0xff, 0x86, 0x0, 0x1, 0x2, 0x1, 0x7, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1, 0xff, 0x82, 0x0, 0x1, 0x7, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x1, 0x4, 0x0, 0x0, 0x0}).
					Return(100, nil)
				rwm.On("Write", []byte{0x1e, 0xff, 0x81, 0x3, 0x1, 0x1, 0x7, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1, 0xff, 0x82, 0x0, 0x1, 0x1, 0x1, 0x4, 0x54, 0x79, 0x70, 0x65, 0x1, 0xc, 0x0, 0x0, 0x0}).
					Return(100, nil)
				rwm.On("Write", []byte{0x11, 0xff, 0x86, 0x1, 0x1, 0x8, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x0, 0x1, 0xc, 0x0}).
					Return(100, nil)

				return rwm
			},
		},
		{
			name: "Wrong msg type",
			args: func() io.ReadWriter {
				request := protocol.Request{
					Message: protocol.Message{
						Type: protocol.MessageTypeResponse,
					},
					Payload: []int{1, 2, 3},
				}
				requestBuffer := &bytes.Buffer{}
				require.NoError(t, gob.NewEncoder(requestBuffer).Encode(request))

				rwm := &test_helper.ReadWriterMock{}
				rwm.On("Read", mock.Anything).
					Run(func(args mock.Arguments) {
						b := args[0].([]byte)
						copy(b[:], requestBuffer.Bytes())
					}).
					Return(1000, nil)

				return rwm
			},
			wantError: true,
		},
		{
			name: "Conn error",
			args: func() io.ReadWriter {
				rwm := &test_helper.ReadWriterMock{}
				rwm.On("Read", mock.Anything).
					Return(1000, errors.New("example error"))

				return rwm
			},
			wantError: true,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			err := serv(tc.args())
			if tc.wantError {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}

type listenerStub struct {
	net.Listener
}

func (*listenerStub) Accept() (net.Conn, error) {
	return &net.TCPConn{}, nil
}
func (*listenerStub) Close() error {
	return nil
}
