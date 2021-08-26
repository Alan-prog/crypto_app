package main

import (
	"context"
	"github.com/crypto_app/middlewhare"
	"github.com/crypto_app/pkg/crypto_app"
	"github.com/crypto_app/service"
	"github.com/crypto_app/service/httpserver"
	"github.com/crypto_app/tools/db"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	serverPort = "8080"
	login      = "postgres"
	pass       = "somepass"
	name       = "postgres"
	host       = "127.0.0.1"
	dbPort     = uint16(5432)
)

func main() {
	ctx := context.Background()

	dbAdp, err := db.NewDbConnector(ctx, login, pass, host, name, dbPort)
	if err != nil {
		log.Fatalf("error while connecting to db: %v", err)
	}
	defer dbAdp.Close()

	crypto := crypto_app.NewCrypto(dbAdp)
	svc := service.NewService(crypto)

	router := httpserver.NewPreparedServer(svc)
	http.Handle("/", router)

	log.Printf("server starting on port: %s", serverPort)
	log.Fatal(http.ListenAndServe(":"+serverPort, middlewhare.ExampleMiddleware(router)))
}
