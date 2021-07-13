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

ARG CANOPSIS_EDITION
WORKDIR /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/

COPY community/go-engines-community/go.mod /monorepo/community/go-engines-community/go.mod
COPY community/go-engines-community/go.sum /monorepo/community/go-engines-community/go.sum
COPY ${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.mod /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.mod
COPY ${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.sum /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.sum
RUN go mod download

COPY . /monorepo
RUN make canopsis.toml
RUN make build

# Create the final container
FROM alpine_updated AS final_container
RUN adduser -D canopsis
USER canopsis:canopsis
ENTRYPOINT /engine

ARG CANOPSIS_EDITION
COPY --from=builder /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/build/canopsis.toml /
ARG CMD
COPY --from=builder /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/build/cmd/${CMD} /engine

FROM final_container AS canopsis-api
EXPOSE 8082
