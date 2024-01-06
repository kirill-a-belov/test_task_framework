package test_helper

import "github.com/stretchr/testify/mock"

type LoggerMock struct {
	mock.Mock
}

func (lm *LoggerMock) Error(err error, details ...interface{}) {
	lm.Called(err, details)
}

func (lm *LoggerMock) Info(details ...interface{}) {
	lm.Called()
}
