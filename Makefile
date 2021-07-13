build:
	docker container rm --force crypto_app 2>/dev/null && docker build -t crypto_app . && docker run --name crypto_app -e POSTGRES_PASSWORD=somepass -e POSTGRES_USER=postgres -e POSTGRES_DB=postgres --rm -p 6001:5432 -p 8080:8080 -d crypto_app
run:
	docker exec -it crypto_app /crypto
lint:
	golangci-lint   run
