package service

import (
	"context"
	"errors"
	"github.com/crypto_app/pkg/models"
	"github.com/crypto_app/tools"
	"net/http"
)

type crypto interface {
	Alive(ctx context.Context) (output models.AliveResponse, err error)
	Sign(ctx context.Context, input *models.RegisterRequest) (output models.RegisterResponse, err error)
	LogIn(ctx context.Context, input *models.LogInRequest) (output models.RegisterResponse, err error)
	GetWallets(ctx context.Context) (output []*models.WalletsResponse, err error)
	Transaction(ctx context.Context, input models.TransactionRequest) (success bool, err error)
	GetTransactions(ctx context.Context, perPage int, pageNum int) (response models.GetTransactionResponse, err error)
}

type Service interface {
	Alive(ctx context.Context) (output models.AliveResponse, err error)
	Sign(ctx context.Context, input *models.RegisterRequest) (output models.RegisterResponse, err error)
	LogIn(ctx context.Context, input *models.LogInRequest) (output models.RegisterResponse, err error)
	GetWallets(ctx context.Context) (output []*models.WalletsResponse, err error)
	Transaction(ctx context.Context, input models.TransactionRequest) (err error)
	GetTransactions(ctx context.Context, perPage int, pageNum int) (response models.GetTransactionResponse, err error)
}

type service struct {
	crypto crypto
}

func (s *service) Alive(ctx context.Context) (output models.AliveResponse, err error) {
	output, err = s.crypto.Alive(ctx)
	return
}

func (s *service) Sign(ctx context.Context, input *models.RegisterRequest) (output models.RegisterResponse, err error) {
	output, err = s.crypto.Sign(ctx, input)
	return
}

func (s *service) LogIn(ctx context.Context, input *models.LogInRequest) (output models.RegisterResponse, err error) {
	output, err = s.crypto.LogIn(ctx, input)
	return
}

func (s *service) GetWallets(ctx context.Context) (output []*models.WalletsResponse, err error) {
	output, err = s.crypto.GetWallets(ctx)
	return
}

func (s *service) Transaction(ctx context.Context, input models.TransactionRequest) (err error) {
	success, err := s.crypto.Transaction(ctx, input)
	if !success {
		err = tools.NewErrorMessage(errors.New("this transaction was aborted"),
			"Данная транзакция не завершилась успешно", http.StatusInternalServerError)
	}
	return
}

func (s *service) GetTransactions(ctx context.Context, perPage int, pageNum int) (response models.GetTransactionResponse, err error) {
	response, err = s.crypto.GetTransactions(ctx, perPage, pageNum)
	return
}

// NewService ...
func NewService(crypto crypto) Service {
	return &service{
		crypto: crypto,
	}
}
