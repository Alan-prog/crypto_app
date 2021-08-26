package httpserver

import (
	"context"
	"encoding/json"
	"github.com/crypto_app/pkg/models"
	"github.com/crypto_app/tools"
	"net/http"
	"strconv"
)

// AliveTransport ...
//================================================
// AliveTransport
//================================================
type AliveTransport interface {
	DecodeRequest(ctx context.Context, r *http.Request) (err error)
	EncodeResponse(ctx context.Context, w http.ResponseWriter, response *models.AliveResponse) (err error)
}

type aliveTransport struct {
}

// DecodeRequest method for decoding requests on server side
func (t *aliveTransport) DecodeRequest(ctx context.Context, r *http.Request) (err error) {
	return
}

// EncodeResponse method for encoding response on server side
func (t *aliveTransport) EncodeResponse(ctx context.Context, w http.ResponseWriter, response *models.AliveResponse) (err error) {
	w.Header().Set("Content-Type", "application/json")
	byteResp, er := json.Marshal(response)
	if er != nil {
		err = tools.NewErrorMessage(err, "error while marshal Alive response", http.StatusInternalServerError)
		return
	}

	_, er = w.Write(byteResp)
	if er != nil {
		err = tools.NewErrorMessage(err, "error while marshal Alive response", http.StatusInternalServerError)
		return
	}
	return
}

// NewAliveTransport the transport creator for http requests
func NewAliveTransport() AliveTransport {
	return &aliveTransport{}
}

// SignInTransport ...
//================================================
// SignInTransport
//================================================
type SignInTransport interface {
	DecodeRequest(ctx context.Context, r *http.Request) (response models.RegisterRequest, err error)
	EncodeResponse(ctx context.Context, w http.ResponseWriter, response *models.RegisterResponse) (err error)
}

type signInTransport struct {
}

// DecodeRequest method for decoding requests on server side
func (t *signInTransport) DecodeRequest(ctx context.Context, r *http.Request) (response models.RegisterRequest, err error) {
	er := json.NewDecoder(r.Body).Decode(&response)
	if er != nil {
		err = tools.NewErrorMessage(er, "Error while unmarshal Sign request", http.StatusInternalServerError)
	}
	return
}

// EncodeResponse method for encoding response on server side
func (t *signInTransport) EncodeResponse(ctx context.Context, w http.ResponseWriter, response *models.RegisterResponse) (err error) {
	byteResp, err := json.Marshal(response)
	if err != nil {
		err = tools.NewErrorMessage(err, "Error while marshal Sign response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(byteResp)
	if err != nil {
		err = tools.NewErrorMessage(err, "Error while writing response to response writer in SignIn method",
			http.StatusInternalServerError)
	}
	return
}

// NewSignInTransport the transport creator for http requests
func NewSignInTransport() SignInTransport {
	return &signInTransport{}
}

// LogInTransport ...
//================================================
// LogInTransport
//================================================
type LogInTransport interface {
	DecodeRequest(ctx context.Context, r *http.Request) (response models.LogInRequest, err error)
	EncodeResponse(ctx context.Context, w http.ResponseWriter, response *models.RegisterResponse) (err error)
}

type logInTransport struct {
}

// DecodeRequest method for decoding requests on server side
func (t *logInTransport) DecodeRequest(ctx context.Context, r *http.Request) (response models.LogInRequest, err error) {
	er := json.NewDecoder(r.Body).Decode(&response)
	if er != nil {
		err = tools.NewErrorMessage(er, "Error while unmarshal LogIn request", http.StatusInternalServerError)
	}
	return
}

// EncodeResponse method for encoding response on server side
func (t *logInTransport) EncodeResponse(ctx context.Context, w http.ResponseWriter, response *models.RegisterResponse) (err error) {
	byteResp, err := json.Marshal(response)
	if err != nil {
		err = tools.NewErrorMessage(err, "Error while marshal LogIn response", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(byteResp)
	if err != nil {
		err = tools.NewErrorMessage(err, "Error while writing response to response writer in LogIn method",
			http.StatusInternalServerError)
	}
	return
}

// NewLogInTransport the transport creator for http requests
func NewLogInTransport() LogInTransport {
	return &logInTransport{}
}

// GetWallets ...
//================================================
// GetWallets
//================================================
type GetWalletsTransport interface {
	DecodeRequest(ctx context.Context, r *http.Request) (err error)
	EncodeResponse(ctx context.Context, w http.ResponseWriter, response []*models.WalletsResponse) (err error)
}

type getWalletsTransport struct {
}

// DecodeRequest method for decoding requests on server side
func (t *getWalletsTransport) DecodeRequest(ctx context.Context, r *http.Request) (err error) {
	return
}

// EncodeResponse method for encoding response on server side
func (t *getWalletsTransport) EncodeResponse(ctx context.Context, w http.ResponseWriter, response []*models.WalletsResponse) (err error) {
	byteResp, err := json.Marshal(response)
	if err != nil {
		err = tools.NewErrorMessage(err, "Error while marshal GetWallets response",
			http.StatusInternalServerError)
		return
	}

	_, err = w.Write(byteResp)
	if err != nil {
		err = tools.NewErrorMessage(err,
			"Error while writing response to response writer in GetWallets method",
			http.StatusInternalServerError)
	}
	return
}

// NewGetWalletsTransport the transport creator for http requests
func NewGetWalletsTransport() GetWalletsTransport {
	return &getWalletsTransport{}
}

// TransactionTransport ...
//================================================
// TransactionTransport
//================================================
type TransactionTransport interface {
	DecodeRequest(ctx context.Context, r *http.Request) (response models.TransactionRequest, err error)
	EncodeResponse(ctx context.Context, w http.ResponseWriter) (err error)
}

type transactionTransport struct {
}

// DecodeRequest method for decoding requests on server side
func (t *transactionTransport) DecodeRequest(ctx context.Context, r *http.Request) (response models.TransactionRequest, err error) {
	er := json.NewDecoder(r.Body).Decode(&response)
	if er != nil {
		err = tools.NewErrorMessage(er, "Error while unmarshal LogIn request", http.StatusInternalServerError)
	}
	return
}

// EncodeResponse method for encoding response on server side
func (t *transactionTransport) EncodeResponse(ctx context.Context, w http.ResponseWriter) (err error) {
	return
}

// NewTransactionTransport the transport creator for http requests
func NewTransactionTransport() TransactionTransport {
	return &transactionTransport{}
}

// GetTransactionsTransport ...
//================================================
// GetTransactionsTransport
//================================================
type GetTransactionsTransport interface {
	DecodeRequest(ctx context.Context, r *http.Request) (perPage, pageNum int, err error)
	EncodeResponse(ctx context.Context, w http.ResponseWriter, response models.GetTransactionResponse) (err error)
}

type getTransactionsTransport struct {
}

// DecodeRequest method for decoding requests on server side
func (t *getTransactionsTransport) DecodeRequest(ctx context.Context, r *http.Request) (perPage, pageNum int, err error) {
	perPage, err = strconv.Atoi(r.URL.Query().Get("per_page"))
	if err != nil {
		err = tools.NewErrorMessage(err, "Неправильно переданы query параметры", http.StatusBadRequest)
		return
	}
	pageNum, err = strconv.Atoi(r.URL.Query().Get("page_num"))
	if err != nil {
		err = tools.NewErrorMessage(err, "Неправильно переданы query параметры", http.StatusBadRequest)
	}
	return
}

// EncodeResponse method for encoding response on server side
func (t *getTransactionsTransport) EncodeResponse(ctx context.Context, w http.ResponseWriter, response models.GetTransactionResponse) (err error) {
	byteResp, err := json.Marshal(response)
	if err != nil {
		err = tools.NewErrorMessage(err, "Error while marshal GetWallets response",
			http.StatusInternalServerError)
		return
	}

	_, err = w.Write(byteResp)
	if err != nil {
		err = tools.NewErrorMessage(err,
			"Error while writing response to response writer in GetWallets method",
			http.StatusInternalServerError)
	}
	return
}

// NewGetTransactionsTransport the transport creator for http requests
func NewGetTransactionsTransport() GetTransactionsTransport {
	return &getTransactionsTransport{}
}
