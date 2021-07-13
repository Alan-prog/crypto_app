package main

import (
	"context"
	"log"
	"math/rand"
	"my_projects/crypto/middlewhare"
	"my_projects/crypto/pkg/crypto_app"
	"my_projects/crypto/service"
	"my_projects/crypto/service/httpserver"
	"my_projects/crypto/tools/db"
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
