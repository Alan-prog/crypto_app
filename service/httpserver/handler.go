package httpserver

import (
	"github.com/gorilla/mux"
	"net/http"
)

// HandlerSettings ...
type HandlerSettings struct {
	Path    string
	Method  string
	Handler http.HandlerFunc
}

// MakeRouter ...
func MakeRouter(handlerSettings []*HandlerSettings) *mux.Router {
	router := mux.NewRouter()

	for _, settings := range handlerSettings {
		router.HandleFunc(settings.Path, settings.Handler).Methods(settings.Method)
	}

	return router
}
