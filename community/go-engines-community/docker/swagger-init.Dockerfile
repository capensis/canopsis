# Note: as a special exception, this can run on its own variable Golang tag
FROM golang:1.16-alpine

RUN apk add --no-cache gcc binutils binutils-gold libc-dev

ENV GO111MODULE on

WORKDIR /go/src/community/go-engines-community
COPY . .

RUN \
    go get -u github.com/swaggo/swag/cmd/swag@v1.6.7 && \
    go get -u github.com/swaggo/http-swagger && \
    go get -u github.com/alecthomas/template

CMD swag init -g ./cmd/canopsis-api/main.go
