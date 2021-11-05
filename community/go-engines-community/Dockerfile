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
RUN mkdir -p /opt/canopsis/var/lib/junit-files /opt/canopsis/var/lib/remediation-files /tmp/canopsis/junit && chown canopsis:canopsis /opt/canopsis/var/lib/junit-files /opt/canopsis/var/lib/remediation-files /tmp/canopsis/junit
RUN ln -sf /canopsis.toml /opt/canopsis/etc/canopsis.toml

COPY community/go-engines-community/config /opt/canopsis/share/config/

ARG CANOPSIS_EDITION
ARG CMD

RUN ln -sf /canopsis-${CANOPSIS_EDITION}.toml /canopsis.toml
COPY community/go-engines-community/cmd/canopsis-reconfigure/canopsis-${CANOPSIS_EDITION}.toml /canopsis-${CANOPSIS_EDITION}.toml
RUN ln -sf /${CMD} /opt/canopsis/bin/${CMD}
COPY --from=builder /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/build/cmd/${CMD} /${CMD}
RUN case "${CMD}" in canopsis-api*) ln -sf /canopsis-api-${CANOPSIS_EDITION} /canopsis-api && ( cd /opt/canopsis/bin && ln -sf canopsis-api-${CANOPSIS_EDITION} canopsis-api ) ;; esac

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
# canopsis-api-* overrides
#

FROM final-container AS canopsis-api-community
EXPOSE 8082

FROM final-container AS canopsis-api-pro
EXPOSE 8082