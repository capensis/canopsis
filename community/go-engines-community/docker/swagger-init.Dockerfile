# sync with Makefile.var:GOLANG_IMAGE_TAG
FROM golang:1.16.4

ENV GO111MODULE on

WORKDIR /go/src/community/go-engines-community
COPY . .

RUN \
    go get -u github.com/swaggo/swag/cmd/swag@v1.6.7 && \
    go get -u github.com/swaggo/http-swagger && \
    go get -u github.com/alecthomas/template

CMD swag init -g ./cmd/canopsis-api/main.go
