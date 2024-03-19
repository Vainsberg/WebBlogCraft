package service

import "github.com/stretchr/testify/mock"

type MockRepositoryUsers struct {
	mock.Mock
}

func (m *MockRepositoryUsers) AddPasswordAndUserName(userName, userPassword string) error {
	args := m.Called(userName, userPassword)
	return args.Error(0)
}
