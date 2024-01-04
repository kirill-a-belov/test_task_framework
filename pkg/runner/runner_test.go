package runner

import (
	"context"
	"errors"
	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"os"
	"syscall"
	"testing"
)

func TestRunner_Run(t *testing.T) {
	testCasesList := []struct {
		name        string
		args        func() (appMock appController, logMock logger.Logger)
		wantSigTerm bool
	}{
		{
			name: "Success",
			args: func() (appController, logger.Logger) {
				appMock := &appMock{}
				appMock.On("Start", mock.Anything).Return(nil)
				appMock.On("Stop", mock.Anything)

				logMock := &logMock{}

				return appMock, logMock
			},
		},
		{
			name: "Error",
			args: func() (appController, logger.Logger) {
				appMock := &appMock{}
				appMock.On("Start", mock.Anything).Return(errors.New("example error"))

				logMock := &logMock{}
				appMock.On("Error", mock.Anything)

				return appMock, logMock
			},
		},
		{
			name: "Panic",
			args: func() (appController, logger.Logger) {
				appMock := &appMock{}
				appMock.On("Start", mock.Anything).Panic("example panic")

				logMock := &logMock{}
				appMock.On("Error", mock.Anything)

				return appMock, logMock
			},
		},
	}

	for _, tc := range testCasesList {
		t.Run(tc.name, func(t *testing.T) {
			testRunner := New(tc.args())
			go func() {
				require.NotPanics(t, func() {
					testRunner.Run(context.Background())
				})
			}()

			if tc.wantSigTerm {
				testRunner.sigChan = make(chan os.Signal, 1)
				testRunner.sigChan <- syscall.SIGTERM
			}
		})
	}
}

type appMock struct {
	mock.Mock
}

func (am *appMock) Start(context.Context) error {
	args := am.Called()

	return args.Error(0)
}
func (am *appMock) Stop(context.Context) {
	am.Called()
}

type logMock struct {
	mock.Mock
}

func (lm *logMock) Error(err error, details ...interface{}) {
	lm.Called(err, details)
}
