package middlewhare

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"my_projects/crypto/pkg/models"
	"my_projects/crypto/tools"
	"net/http"
	"strings"
)

const (
	auth = "Authorization"
)

var (
	publicAPI = [][]string{
		{"POST", "/crypto/register"},
		{"POST", "/crypto/log_in"},
	}
)

func ExampleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Our middleware logic goes here...
		if isRoutePublic(r.Method, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		p := strings.Split(r.Header.Get(auth), " ")
		if len(p) != 2 || p[0] != "Bearer" {
			err := tools.NewErrorMessage(errors.New("bad token format"), "Неправильный формат токена",
				http.StatusUnauthorized)
			tools.EncodeIntoResponseWriter(w, err)
			return
		}

		userID, err := checkTheTokenValidness(p[1])
		if err != nil {
			tools.EncodeIntoResponseWriter(w, err.(tools.ErrorMessage))
			return
		}
		ctx := context.WithValue(r.Context(), models.CtxKey("id"), userID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func isRoutePublic(method, url string) (response bool) {
	for i := range publicAPI {
		if len(publicAPI[i]) == 2 {
			if publicAPI[i][0] == method {
				if publicAPI[i][1] == url {
					return true
				}
			}
			continue
		}
	}
	return false
}

func checkTheTokenValidness(tokenString string) (userID string, err error) {
	const (
		tokenInvalidErr = "token is invalid"
	)

	token, err := jwt.ParseWithClaims(tokenString, &models.ClaimWithID{}, func(token *jwt.Token) (interface{}, error) {
		return models.JwtSigningKey, nil
	})

	if err != nil {
		err = tools.NewErrorMessage(err, "Некорректный токен", http.StatusUnauthorized)
		return
	}

	if claims, ok := token.Claims.(*models.ClaimWithID); ok && token.Valid {
		userID = claims.ID
		return
	} else {
		err = tools.NewErrorMessage(errors.New(tokenInvalidErr), "Некорректный токен", http.StatusUnauthorized)
	}

	return

}
