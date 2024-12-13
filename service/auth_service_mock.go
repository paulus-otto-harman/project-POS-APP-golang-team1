package service

import (
	"github.com/stretchr/testify/mock"
	"project/domain"
)

type AuthServiceMock struct {
	mock.Mock
}

func (serviceMock *AuthServiceMock) Login(user domain.User) (string, bool, error) {
	args := serviceMock.Called(user)
	if sessionResult := args.Get(1); sessionResult != nil {
		return "", sessionResult.(bool), args.Error(2)
	}
	return "", false, args.Error(2)
}
