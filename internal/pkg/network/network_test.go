package network

import (
	"context"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	testCases := []struct {
		name      string
		args      func() (conn, interface{}, error)
		wantError bool
	}{
		{
			name: "Success",
			args: func() (conn, interface{}, error) {
				resultMsg := exampleMessage{
					ExampleFieldOne: 1,
					ExampleFieldTwo: "example",
				}

				mockBytesPartOne := []byte{0x43, 0x7f, 0x3, 0x1, 0x1, 0xe, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1, 0xff, 0x80, 0x0, 0x1, 0x2, 0x1, 0xf, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x6e, 0x65, 0x1, 0x4, 0x0, 0x1, 0xf, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x77, 0x6f, 0x1, 0xc, 0x0, 0x0, 0x0}
				mockBytesPartTwo := []byte{0xe, 0xff, 0x80, 0x1, 0x2, 0x1, 0x7, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x0}
				connMock := &connMock{}
				connMock.On("Write", mockBytesPartOne).Return(len(mockBytesPartOne), nil)
				connMock.On("Write", mockBytesPartTwo).Return(len(mockBytesPartTwo), nil)

				return connMock, resultMsg, nil
			},
			wantError: false,
		},
		{
			name: "Bad conn",
			args: func() (conn, interface{}, error) {
				connMock := &connMock{}
				connMock.On("Write", mock.Anything).Return(0, errors.New("example error"))

				return connMock, nil, nil
			},
			wantError: true,
		},
		{
			name: "Bad msg",
			args: func() (conn, interface{}, error) {
				return nil, make(chan int), nil
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputConn, inputMsg, err := tc.args()
			require.NoError(t, err)

			err = Send(context.Background(), inputConn, inputMsg)
			if tc.wantError {
				assert.Error(t, err)

				return
			}
			assert.NoError(t, err)
		})
	}
}

type exampleMessage struct {
	ExampleFieldOne int
	ExampleFieldTwo string
}

type connMock struct {
	mock.Mock
}

func (cm *connMock) Write(b []byte) (int, error) {
	args := cm.Called(b)

	return args.Int(0), args.Error(1)
}
func (cm *connMock) Read(b []byte) (int, error) {
	args := cm.Called(b)

	return args.Int(0), args.Error(1)
}

func TestReceive(t *testing.T) {
	testCases := []struct {
		name      string
		args      func() (conn, interface{}, error)
		wantError bool
	}{
		{
			name: "Success",
			args: func() (conn, interface{}, error) {
				resultMsg := exampleMessage{
					ExampleFieldOne: 1,
					ExampleFieldTwo: "example",
				}

				mockBytesPartOne := []byte{0x43, 0x7f, 0x3, 0x1, 0x1, 0xe, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x1, 0xff, 0x80, 0x0, 0x1, 0x2, 0x1, 0xf, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x4f, 0x6e, 0x65, 0x1, 0x4, 0x0, 0x1, 0xf, 0x45, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x77, 0x6f, 0x1, 0xc, 0x0, 0x0, 0x0, 0xe, 0xff, 0x80, 0x1, 0x2, 0x1, 0x7, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x0}
				connMock := &connMock{}
				connMock.On("Read", mock.Anything).
					Run(func(args mock.Arguments) {
						bytes := args[0].([]byte)
						copy(bytes[:], mockBytesPartOne)
					}).
					Return(len(mockBytesPartOne), nil)

				return connMock, resultMsg, nil
			},
			wantError: false,
		},
		{
			name: "Bad conn",
			args: func() (conn, interface{}, error) {
				connMock := &connMock{}
				connMock.On("Read", mock.Anything).Return(1, errors.New("example error"))

				return connMock, nil, nil
			},
			wantError: true,
		},
		{
			name: "Bad msg",
			args: func() (conn, interface{}, error) {
				connMock := &connMock{}
				connMock.On("Read", mock.Anything).
					Run(func(args mock.Arguments) {
						bytes := args[0].([]byte)
						copy(bytes[:], "example bad msg")
					}).
					Return(1, nil)
				return connMock, nil, nil
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			conn, msg, err := tc.args()
			require.NoError(t, err)

			result, err := Receive[exampleMessage](context.Background(), conn)
			if tc.wantError {
				assert.Error(t, err)

				return
			}
			assert.NoError(t, err)

			assert.Equal(t, msg, result)
		})
	}
}
