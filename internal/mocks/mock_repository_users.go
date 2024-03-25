package mocks

import "github.com/stretchr/testify/mock"

type MockRepositoryUsers struct {
	mock.Mock
}

func (m *MockRepositoryUsers) AddPasswordAndUserName(userName, userPassword string) error {
	args := m.Called(userName, userPassword)
	return args.Error(0)
}

func (m *MockRepositoryUsers) SearchPasswordAndUserName(userName string) (string, error) {
	args := m.Called(userName)
	return "", args.Error(0)
}

func (m *MockRepositoryUsers) CheckingPresenceUser(username string) (string, error) {
	args := m.Called(username)
	return "", args.Error(0)
}

func (m *MockRepositoryUsers) SearchUserName(UserID int) (string, error) {
	args := m.Called(UserID)
	return "", args.Error(0)
}
