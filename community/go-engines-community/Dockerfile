# Args defined here are only available to from commands
# See https://docs.docker.com/engine/reference/builder/#understand-how-arg-and-from-interact
ARG ALPINE_VERSION
ARG GOLANG_VERSION

# Upgrade alpine
FROM alpine:${ALPINE_VERSION} AS alpine_updated
RUN apk upgrade --no-cache

# Build the engines
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS builder
RUN apk add --no-cache make git gcc build-base g++ coreutils php

WORKDIR /go-engines

COPY go.mod /go-engines/go.mod
COPY go.sum /go-engines/go.sum
RUN go mod download

COPY . /go-engines

RUN make canopsis.toml
RUN make build

# Create the final container
FROM alpine_updated AS final_container
RUN adduser -D canopsis
USER canopsis:canopsis
ENTRYPOINT /engine
COPY --from=builder /go-engines/canopsis.toml /
ARG CMD
COPY --from=builder /go-engines/build/cmd/${CMD} /engine


FROM final_container AS canopsis-api
EXPOSE 8082
