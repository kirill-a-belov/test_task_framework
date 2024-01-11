package client

import (
	"bytes"
	"context"
	"encoding/gob"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kirill-a-belov/test_task_framework/internal/app/client/pkg/config"
	"github.com/kirill-a-belov/test_task_framework/internal/pkg/protocol"
	"github.com/kirill-a-belov/test_task_framework/pkg/test_helper"
)

func TestClient_Start(t *testing.T) {
	testCaseList := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "Success",
			wantError: false,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()
			s := New(ctx, &config.Config{})

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

func TestClient_processor(t *testing.T) {
	testCaseList := []struct {
		name string
		args func() (*Client, func(io.ReadWriter) error)
	}{
		{
			name: "Regular stop",
			args: func() (*Client, func(io.ReadWriter) error) {
				f := func(writer io.ReadWriter) error {
					return nil
				}

				loggerMock := &test_helper.LoggerMock{}
				loggerMock.On("Info", mock.Anything)

				c := New(context.Background(), &config.Config{
					ConnTTL: time.Second,
				})
				c.logger = loggerMock

				diallerMock := &diallerMock{}
				diallerMock.On("mockFunc").Return(&net.TCPConn{}, nil)
				c.dialler = diallerMock.mockFunc

				return c, f
			},
		},
		{
			name: "Handling func error",
			args: func() (*Client, func(io.ReadWriter) error) {
				f := func(writer io.ReadWriter) error {
					return errors.New("example error")
				}

				loggerMock := &test_helper.LoggerMock{}
				loggerMock.On("Error", mock.Anything, mock.Anything)
				loggerMock.On("Info", mock.Anything)

				c := New(context.Background(), &config.Config{
					ConnTTL: time.Second,
				})
				c.logger = loggerMock

				diallerMock := &diallerMock{}
				diallerMock.On("mockFunc").Return(&net.TCPConn{}, nil)
				c.dialler = diallerMock.mockFunc

				return c, f
			},
		},
		{
			name: "Dialing error",
			args: func() (*Client, func(io.ReadWriter) error) {
				f := func(writer io.ReadWriter) error {
					return nil
				}

				loggerMock := &test_helper.LoggerMock{}
				loggerMock.On("Error", mock.Anything, mock.Anything) //.Once()
				loggerMock.On("Info", mock.Anything)

				c := New(context.Background(), &config.Config{
					ConnTTL: time.Second,
				})
				c.logger = loggerMock

				diallerMock := &diallerMock{}
				diallerMock.On("mockFunc").Return(&net.TCPConn{}, errors.New("example error"))
				diallerMock.On("mockFunc").Return(&net.TCPConn{}, nil)

				return c, f
			},
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			testClient, testFunc := tc.args()

			ctx := context.Background()

			var wg sync.WaitGroup
			wg.Add(1)
			go func() {
				testClient.processor(ctx, testFunc)
				wg.Done()
			}()

			time.Sleep(time.Millisecond)
			testClient.Stop(ctx)
			wg.Wait()
		})
	}
}

type diallerMock struct {
	mock.Mock
}

func (dm *diallerMock) mockFunc() (net.Conn, error) {
	args := dm.Called()

	return args.Get(0).(net.Conn), args.Error(1)
}

func TestClient_Stop(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		s := New(ctx, &config.Config{})

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
					Payload: []int64{1, 2, 3},
				}
				requestBuffer := &bytes.Buffer{}
				require.NoError(t, gob.NewEncoder(requestBuffer).Encode(request))

				response := protocol.Response{
					Message: protocol.Message{
						Type: protocol.MessageTypeResponse,
					},
					Payload: 6,
				}
				responseBuffer := &bytes.Buffer{}
				require.NoError(t, gob.NewEncoder(responseBuffer).Encode(response))

				rwm := &test_helper.ReadWriterMock{}
				rwm.On("Write", mock.Anything).
					Run(func(args mock.Arguments) {
						b := args[0].([]byte)
						copy(b[:], requestBuffer.Bytes())
					}).
					Return(1000, nil)

				rwm.On("Read", mock.Anything).
					Run(func(args mock.Arguments) {
						b := args[0].([]byte)
						copy(b[:], responseBuffer.Bytes())
					}).
					Return(1000, nil)

				return rwm
			},
		},
		{
			name: "Wrong message type",
			args: func() io.ReadWriter {
				request := protocol.Request{
					Message: protocol.Message{
						Type: protocol.MessageTypeRequest,
					},
					Payload: []int64{1, 2, 3},
				}
				requestBuffer := &bytes.Buffer{}
				require.NoError(t, gob.NewEncoder(requestBuffer).Encode(request))

				response := protocol.Response{
					Message: protocol.Message{
						Type: protocol.MessageTypeRequest,
					},
					Payload: 6,
				}
				responseBuffer := &bytes.Buffer{}
				require.NoError(t, gob.NewEncoder(responseBuffer).Encode(response))

				rwm := &test_helper.ReadWriterMock{}
				rwm.On("Write", mock.Anything).
					Run(func(args mock.Arguments) {
						b := args[0].([]byte)
						copy(b[:], requestBuffer.Bytes())
					}).
					Return(1000, nil)

				rwm.On("Read", mock.Anything).
					Run(func(args mock.Arguments) {
						b := args[0].([]byte)
						copy(b[:], responseBuffer.Bytes())
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
				rwm.On("Write", mock.Anything).
					Return(1000, errors.New("example error"))

				return rwm
			},
			wantError: true,
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			err := handle(tc.args())
			if tc.wantError {
				assert.Error(t, err)

				return
			}

			assert.NoError(t, err)
		})
	}
}
