package logger

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"io"
	stdLog "log"
	"os"
	"testing"
)

func Test_log_Error(t *testing.T) {
	type testArgs struct {
		prefix     string
		err        error
		payload    []interface{}
		stdOutMock io.Writer
	}

	testCaseList := []struct {
		name string
		args func() testArgs
	}{
		{
			name: "Success",
			args: func() testArgs {
				som := &stdOutMock{}
				som.On("Write", []byte("example_prefixlogger.go:27:  [example_str 1]  example error \n")).
					Return(mock.Anything)

				return testArgs{
					prefix:     "example_prefix",
					err:        errors.New("example error"),
					payload:    []interface{}{"example_str", 1},
					stdOutMock: som,
				}
			},
		},
	}

	for _, tc := range testCaseList {
		t.Run(tc.name, func(t *testing.T) {
			ta := tc.args()

			l := stdLog.New(os.Stdout, ta.prefix, stdLog.Lshortfile)
			l.SetOutput(ta.stdOutMock)

			testLog := log{
				l: l,
			}
			testLog.Error(ta.err, ta.payload)
		})
	}
}

type stdOutMock struct {
	mock.Mock
}

func (som *stdOutMock) Write(p []byte) (n int, err error) {
	som.Called(p)

	return
}
