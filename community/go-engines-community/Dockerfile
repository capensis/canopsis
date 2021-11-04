# Args defined here are only available to FROM commands
# See https://docs.docker.com/engine/reference/builder/#understand-how-arg-and-from-interact
ARG ALPINE_VERSION
ARG ALPINE_BUILD_VERSION
ARG GOLANG_VERSION

#
# Build the engines
#

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_BUILD_VERSION} AS builder
RUN apk add --no-cache make git gcc binutils binutils-gold libc-dev

ARG CANOPSIS_EDITION
WORKDIR /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/

# Do a first pass that only handles Go dependencies, so that we can cache this part
COPY community/go-engines-community/go.mod /monorepo/community/go-engines-community/go.mod
COPY community/go-engines-community/go.sum /monorepo/community/go-engines-community/go.sum
COPY community/go-engines-community/Makefile /monorepo/community/go-engines-community/Makefile
COPY ${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.mod /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.mod
COPY ${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.sum /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/go.sum
COPY ${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/Makefile /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/Makefile
RUN make deps

# Always build everything in the current $CANOPSIS_EDITION
COPY . /monorepo
RUN make

#
# Prepare the common final container
#

FROM alpine:${ALPINE_VERSION} AS final-container

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup --system canopsis && adduser --system --disabled-password --shell /sbin/nologin --no-create-home --home /opt/canopsis --ingroup canopsis canopsis
RUN mkdir -p /opt/canopsis/bin /opt/canopsis/etc /opt/canopsis/share
RUN ln -sf /canopsis.toml /opt/canopsis/etc/canopsis.toml

COPY community/go-engines-community/config /opt/canopsis/share/config/

ARG CANOPSIS_EDITION
ARG CMD

RUN ln -sf /${CMD} /opt/canopsis/bin/${CMD}
COPY community/go-engines-community/cmd/canopsis-reconfigure/canopsis-${CANOPSIS_EDITION}.toml /canopsis.toml
COPY --from=builder /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/build/cmd/${CMD} /${CMD}

USER canopsis:canopsis

CMD ["/bin/sh", "-c", "/${CMD}"]

#
# engine-webhook override
#

FROM final-container AS engine-webhook

USER root
COPY pro/go-engines-pro/cmd/engine-webhook/webhook.conf.toml /opt/canopsis/etc/webhook.conf.toml
USER canopsis:canopsis

#
# connector-junit override
#

FROM final-container AS connector-junit

USER root
RUN mkdir -p /opt/canopsis/var/lib/junit-files && chown canopsis:canopsis /opt/canopsis/var/lib/junit-files
RUN mkdir -p /tmp/canopsis/junit && chown canopsis:canopsis /tmp/canopsis /tmp/canopsis/junit
USER canopsis:canopsis

#
# canopsis-api override
#

FROM final-container AS canopsis-api

EXPOSE 8082
