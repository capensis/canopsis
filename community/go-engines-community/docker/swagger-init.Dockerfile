FROM golang:latest

ADD . /go/src/git.canopsis.net/canopsis/canopsis-community/community/go-engines-community

WORKDIR /go/src/git.canopsis.net/canopsis/canopsis-community/community/go-engines-community

RUN \
    go get -u github.com/swaggo/swag/cmd/swag@v1.6.7 && \
    go get -u github.com/swaggo/http-swagger && \
    go get -u github.com/alecthomas/template

CMD swag init -g ./cmd/canopsis-api/main.go
