package runner

import (
	"context"
	"errors"
	"os"
	"syscall"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/kirill-a-belov/test_task_framework/pkg/logger"
	"github.com/kirill-a-belov/test_task_framework/pkg/test_helper"
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

				logMock := &test_helper.LoggerMock{}

				return appMock, logMock
			},
		},
		{
			name: "Error",
			args: func() (appController, logger.Logger) {
				appMock := &appMock{}
				appMock.On("Start", mock.Anything).Return(errors.New("example error"))

				logMock := &test_helper.LoggerMock{}
				logMock.On("Error", mock.Anything)

				return appMock, logMock
			},
		},
		{
			name: "Panic",
			args: func() (appController, logger.Logger) {
				appMock := &appMock{}
				appMock.On("Start", mock.Anything).Panic("example panic")

				logMock := &test_helper.LoggerMock{}
				logMock.On("Error", mock.Anything)

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
