package logger

import (
	"errors"
	"io"
	stdLog "log"
	"os"
	"testing"

	"github.com/stretchr/testify/mock"
)

func Test_log_Error(t *testing.T) {
	logTest(t, "Error")
}

type stdOutMock struct {
	mock.Mock
}

func (som *stdOutMock) Write(p []byte) (n int, err error) {
	som.Called(p)

	return
}

func Test_log_Info(t *testing.T) {
	logTest(t, "Info")
}

func logTest(t *testing.T, method string) {
	var (
		testLog          log
		mockCalledResult = make([]byte, 0)
		testFuncError    func(error, ...interface{})
		testFuncInfo     func(...interface{})
	)

	switch method {
	case "Error":
		mockCalledResult = []byte("example_prefix logger.go:28:  [example_str 1]  example error \n")
		testFuncError = testLog.Error

	case "Info":
		mockCalledResult = []byte("example_prefix logger.go:28:  [example_str 1]  <nil> \n")
		testFuncInfo = testLog.Info

	default:
		t.Fatal("invalid method name")
	}

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
				som.On("Write", mockCalledResult).
					Return(mock.Anything)

				return testArgs{
					prefix:     "example_prefix ",
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

			testLog.l = l

			if testFuncError != nil {
				testFuncError(ta.err, ta.payload)
			}
			if testFuncInfo != nil {
				testFuncInfo(ta.payload)
			}
		})
	}
}
