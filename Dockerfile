FROM golang:1.13.10-alpine AS build
ARG GOOS
ENV CGO_ENABLED=0 \
    GOOS=$GOOS \
    GOARCH=amd64 \
    CGO_CPPFLAGS="-I/usr/include" \
    UID=0 GID=0 \
    CGO_CFLAGS="-I/usr/include" \
    CGO_LDFLAGS="-L/usr/lib -lpthread -lrt -lstdc++ -lm -lc -lgcc -lz " \
    PKG_CONFIG_PATH="/usr/lib/pkgconfig"
RUN apk add --no-cache git make
WORKDIR /go/src/my_projects/crypto
COPY ./cmd ./cmd
COPY ./pkg ./pkg
COPY ./service ./service
COPY ./middlewhare ./middlewhare
COPY ./tools ./tools
COPY ./vendor ./vendor
RUN go test ./...
RUN mkdir bin
RUN go build  -o bin/crypto cmd/main.go

FROM postgres
COPY migrations.sql /docker-entrypoint-initdb.d/
COPY --from=build /go/src/my_projects/crypto/bin/crypto .
EXPOSE 8080