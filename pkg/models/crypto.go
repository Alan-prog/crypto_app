package models

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
)

type CtxKey string

var JwtSigningKey = []byte("secret")

const (
	SqlNoRows = "no rows in result set"
)

type AliveResponse struct {
	Text   string `json:"text"`
	UserID int32  `json:"user_id"`
}

type RegisterRequest struct {
	Name     string `json:"name"`
	LastName string `json:"last_name"`
	Email    string `json:"email"`
	Pass     string `json:"pass"`
}

type LogInRequest struct {
	Email string `json:"email"`
	Pass  string `json:"pass"`
}

type TransactionRequest struct {
	FromAddress int32   `json:"from_address"`
	ToAddress   int32   `json:"to_address"`
	Amount      float64 `json:"amount"`
}

type WalletsResponse struct {
	Salary  string  `json:"salary"`
	Balance float64 `json:"balance"`
	Address string  `json:"address"`
}

type RegisterResponse struct {
	AccessToken string `json:"access_token"`
}

type SingleUserDataDbResponse struct {
	ID          int
	Name        string
	DOB         sql.NullTime
	Address     string
	Description string
	CreateAt    sql.NullTime
	UpdatedAt   sql.NullTime
}

type SingleUserData struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	DOB         *string `json:"dob"`
	Address     string  `json:"address"`
	Description string  `json:"description"`
	CreateAt    string  `json:"createAt"`
	UpdatedAt   *string `json:"updatedAt"`
}

type AllUsersData []*SingleUserData

type UpdateUserData struct {
	ID          int
	Name        *string
	DOB         *string
	Address     *string
	Description *string
}

type ClaimWithID struct {
	ID string `json:"custom_id"`
	jwt.StandardClaims
}

type SingleTransaction struct {
	FromAddress string  `json:"from_address"`
	ToAddress   string  `json:"to_address"`
	Sum         float64 `json:"sum"`
	Commission  float64 `json:"commission"`
	Date        string  `json:"date"`
	Success     bool    `json:"success"`
}

type Meta struct {
	Total      int32 `json:"total"`
	PageNum    int32 `json:"page_num"`
	PagesCount int32 `json:"pages_count"`
}

type GetTransactionResponse struct {
	Items []*SingleTransaction `json:"items"`
	Meta  Meta                 `json:"meta"`
}
