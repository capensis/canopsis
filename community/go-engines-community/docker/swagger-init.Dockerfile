# sync with Makefile.var:GOLANG_IMAGE_TAG
FROM golang:1.16.4

ADD . ${GOPATH}/src/git.canopsis.net/canopsis/canopsis-community/community/go-engines-community

WORKDIR ${GOPATH}/src/git.canopsis.net/canopsis/canopsis-community/community/go-engines-community

RUN \
    go get -u github.com/swaggo/swag/cmd/swag@v1.6.7 && \
    go get -u github.com/swaggo/http-swagger && \
    go get -u github.com/alecthomas/template

CMD swag init -g ./cmd/canopsis-api/main.go
