package service

import (
	"context"
	"github.com/stretchr/testify/mock"
	"my_projects/crypto/pkg/models"
	"my_projects/crypto/tools"
)

// MockService ...
type MockService struct {
	mock.Mock
}

// Alive ...
func (m *MockService) Alive(ctx context.Context) (output models.AliveResponse, err tools.ErrorMessage) {
	output.Text = "service is okay"
	return
}

// CreateUser ...
func (m *MockService) CreateUser(ctx context.Context, input *models.RegisterRequest) (output models.SingleUserData, err tools.ErrorMessage) {
	args := m.Called(context.Background(), input)
	if a, ok := args.Get(0).(models.SingleUserData); ok {
		if args.Get(1) == nil {
			return a, nil
		} else {
			return output, args.Get(1).(tools.ErrorMessage)
		}
	}
	return output, args.Get(1).(tools.ErrorMessage)
}

func (m *MockService) GetUser(ctx context.Context, input int) (output models.SingleUserData, err tools.ErrorMessage) {
	args := m.Called(context.Background(), input)
	if a, ok := args.Get(0).(models.SingleUserData); ok {
		if args.Get(1) == nil {
			return a, nil
		} else {
			return output, args.Get(1).(tools.ErrorMessage)
		}
	}
	return output, args.Get(1).(tools.ErrorMessage)
}

func (m *MockService) DeleteUser(ctx context.Context, input int) (err tools.ErrorMessage) {
	args := m.Called(context.Background(), input)
	if _, ok := args.Get(0).(models.SingleUserData); ok {
		if args.Get(1) == nil {
			return nil
		} else {
			return args.Get(1).(tools.ErrorMessage)
		}
	}
	return args.Get(1).(tools.ErrorMessage)
}

func (m *MockService) UpdateUser(ctx context.Context, input *models.UpdateUserData) (output models.SingleUserData, err tools.ErrorMessage) {
	return
}

func (m *MockService) GetAllUser(ctx context.Context) (output models.AllUsersData, err tools.ErrorMessage) {
	return
}
