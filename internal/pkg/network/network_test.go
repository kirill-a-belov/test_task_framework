package network

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestSend(t *testing.T) {
	testCases := []struct {
		name      string
		inputData func() (conn, interface{}, error)
		wantError bool
	}{
		{
			name: "Success",
			inputData: func() (conn, interface{}, error) {
				resultMsg := exampleMessage{
					exampleFieldOne: 1,
					exampleFieldTwo: "example",
				}

				mockBytes, err := json.Marshal(resultMsg)
				if err != nil {
					return nil, exampleMessage{}, err
				}

				connMock := &connMock{}
				connMock.On("Write", mockBytes).Return(len(mockBytes), nil)

				return connMock, resultMsg, nil
			},
			wantError: false,
		},
		{
			name: "Bad conn",
			inputData: func() (conn, interface{}, error) {
				resultMsg := exampleMessage{
					exampleFieldOne: 1,
					exampleFieldTwo: "example",
				}

				mockBytes, err := json.Marshal(resultMsg)
				if err != nil {
					return nil, exampleMessage{}, err
				}

				connMock := &connMock{}
				connMock.On("Write", mockBytes).Return(0, errors.New("example error"))

				return connMock, resultMsg, nil
			},
			wantError: true,
		},
		{
			name: "Bad msg",
			inputData: func() (conn, interface{}, error) {
				return nil, make(chan int), nil
			},
			wantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			inputConn, inputMsg, err := tc.inputData()
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
	exampleFieldOne int
	exampleFieldTwo string
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
