package httpserver

import (
	"context"
	"github.com/gorilla/mux"
	"my_projects/crypto/pkg/models"
	"my_projects/crypto/tools"
	"net/http"
)

type service interface {
	Alive(ctx context.Context) (output models.AliveResponse, err error)
	Sign(ctx context.Context, input *models.RegisterRequest) (output models.RegisterResponse, err error)
	LogIn(ctx context.Context, input *models.LogInRequest) (output models.RegisterResponse, err error)
	GetWallets(ctx context.Context) (output []*models.WalletsResponse, err error)
	Transaction(ctx context.Context, input models.TransactionRequest) (err error)
	GetTransactions(ctx context.Context, perPage int, pageNum int) (response models.GetTransactionResponse, err error)
}

//================================================
// AliveServer
//================================================
type aliveServer struct {
	transport AliveTransport
	service   service
}

// ServeHTTP implements http.Handler.
func (s *aliveServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.transport.DecodeRequest(r.Context(), r)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	response, err := s.service.Alive(r.Context())
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	if err := s.transport.EncodeResponse(r.Context(), w, &response); err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}
}

// NewAliveServer the server creator
func NewAliveServer(transport AliveTransport, service service) http.HandlerFunc {
	ls := aliveServer{
		transport: transport,
		service:   service,
	}
	return ls.ServeHTTP
}

//================================================
// SignIn
//================================================
type signServer struct {
	transport SignInTransport
	service   service
}

// ServeHTTP implements http.Handler.
func (s *signServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := s.transport.DecodeRequest(r.Context(), r)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	response, err := s.service.Sign(r.Context(), &req)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	if err := s.transport.EncodeResponse(r.Context(), w, &response); err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}
}

// NewSignServer the server creator
func NewSignServer(transport SignInTransport, service service) http.HandlerFunc {
	ls := signServer{
		transport: transport,
		service:   service,
	}
	return ls.ServeHTTP
}

//================================================
// LogInServer
//================================================
type logInServer struct {
	transport LogInTransport
	service   service
}

// ServeHTTP implements http.Handler.
func (s *logInServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := s.transport.DecodeRequest(r.Context(), r)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	response, err := s.service.LogIn(r.Context(), &req)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	if err := s.transport.EncodeResponse(r.Context(), w, &response); err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}
}

// NewLogInServer the server creator
func NewLogInServer(transport LogInTransport, service service) http.HandlerFunc {
	ls := logInServer{
		transport: transport,
		service:   service,
	}
	return ls.ServeHTTP
}

//================================================
// GetWalletsServer
//================================================
type getWalletsServer struct {
	transport GetWalletsTransport
	service   service
}

// ServeHTTP implements http.Handler.
func (s *getWalletsServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := s.transport.DecodeRequest(r.Context(), r)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	response, err := s.service.GetWallets(r.Context())
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	if err := s.transport.EncodeResponse(r.Context(), w, response); err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}
}

// NewGetWalletsServer the server creator
func NewGetWalletsServer(transport GetWalletsTransport, service service) http.HandlerFunc {
	ls := getWalletsServer{
		transport: transport,
		service:   service,
	}
	return ls.ServeHTTP
}

//================================================
// TransactionServer
//================================================
type transactionServer struct {
	transport TransactionTransport
	service   service
}

// ServeHTTP implements http.Handler.
func (s *transactionServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resp, err := s.transport.DecodeRequest(r.Context(), r)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	err = s.service.Transaction(r.Context(), resp)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	if err := s.transport.EncodeResponse(r.Context(), w); err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}
}

// NewTransactionServer the server creator
func NewTransactionServer(transport TransactionTransport, service service) http.HandlerFunc {
	ls := transactionServer{
		transport: transport,
		service:   service,
	}
	return ls.ServeHTTP
}

//================================================
// GetTransactionServer
//================================================
type getTransactionServer struct {
	transport GetTransactionsTransport
	service   service
}

// ServeHTTP implements http.Handler.
func (s *getTransactionServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	perPage, pageNum, err := s.transport.DecodeRequest(r.Context(), r)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	resp, err := s.service.GetTransactions(r.Context(), perPage, pageNum)
	if err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}

	if err := s.transport.EncodeResponse(r.Context(), w, resp); err != nil {
		tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
		return
	}
}

// NewGetTransactionsServer the server creator
func NewGetTransactionsServer(transport GetTransactionsTransport, service service) http.HandlerFunc {
	ls := getTransactionServer{
		transport: transport,
		service:   service,
	}
	return ls.ServeHTTP
}

// NewPreparedServer ...
func NewPreparedServer(svc service) *mux.Router {
	aliveTransport := NewAliveTransport()
	signInTransport := NewSignInTransport()
	logInTransport := NewLogInTransport()
	getWalletsTransport := NewGetWalletsTransport()
	transactionTransport := NewTransactionTransport()
	getTransactionsTransport := NewGetTransactionsTransport()
	return MakeRouter(
		[]*HandlerSettings{
			{
				Path:    URIPathGetAlive,
				Method:  http.MethodGet,
				Handler: NewAliveServer(aliveTransport, svc),
			},
			{
				Path:    URIPathSignIn,
				Method:  http.MethodPost,
				Handler: NewSignServer(signInTransport, svc),
			},
			{
				Path:    URIPathLogIn,
				Method:  http.MethodPost,
				Handler: NewLogInServer(logInTransport, svc),
			},
			{
				Path:    URIPathGetWallets,
				Method:  http.MethodGet,
				Handler: NewGetWalletsServer(getWalletsTransport, svc),
			},
			{
				Path:    URIPathTransaction,
				Method:  http.MethodPost,
				Handler: NewTransactionServer(transactionTransport, svc),
			},
			{
				Path:    URIPathGetTransactions,
				Method:  http.MethodGet,
				Handler: NewGetTransactionsServer(getTransactionsTransport, svc),
			},
		},
	)
}
