package test_helper

import "github.com/stretchr/testify/mock"

type ReadWriterMock struct {
	mock.Mock
}

func (rwm *ReadWriterMock) Write(b []byte) (int, error) {
	args := rwm.Called(b)

	return args.Int(0), args.Error(1)
}
func (cm *ReadWriterMock) Read(b []byte) (int, error) {
	args := cm.Called(b)

	return args.Int(0), args.Error(1)
}
