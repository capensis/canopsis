# Args defined here are only available to from commands
# See https://docs.docker.com/engine/reference/builder/#understand-how-arg-and-from-interact
ARG ALPINE_VERSION
ARG ALPINE_BUILD_VERSION
ARG GOLANG_VERSION

# Build the engines
FROM golang:${GOLANG_VERSION}-alpine${ALPINE_BUILD_VERSION} AS builder
RUN apk add --no-cache make git gcc build-base g++ coreutils php

ARG CANOPSIS_EDITION
WORKDIR /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/

COPY community/go-engines-community/go.mod /monorepo/community/go-engines-community/go.mod
COPY community/go-engines-community/go.sum /monorepo/community/go-engines-community/go.sum
COPY ${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.mod /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.mod
COPY ${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.sum /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.sum
COPY community/go-engines-community/Makefile /monorepo/community/go-engines-community/Makefile
COPY ${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/Makefile /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/Makefile
RUN make deps

COPY . /monorepo
RUN make canopsis.toml
RUN make

ARG CANOPSIS_EDITION
# Prepare the final container
FROM alpine:${ALPINE_VERSION} AS final-container

RUN adduser -D canopsis
RUN mkdir -p /opt/canopsis/share && chown canopsis:canopsis /opt/canopsis/share
COPY community/go-engines-community/config /opt/canopsis/share/config/
RUN cp /opt/canopsis/share/config/canopsis.${CANOPSIS_EDITION}.toml /opt/canopsis/share/config/canopsis.toml
USER canopsis:canopsis

ENTRYPOINT /engine

COPY --from=builder /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/build/canopsis.toml /
ARG CMD
COPY --from=builder /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/build/cmd/${CMD} /engine

FROM final-container AS canopsis-api
EXPOSE 8082
