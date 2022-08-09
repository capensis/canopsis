# Args defined here are only available to FROM commands
# See https://docs.docker.com/engine/reference/builder/#understand-how-arg-and-from-interact
ARG ALPINE_VERSION
ARG ALPINE_BUILD_VERSION
ARG GOLANG_VERSION

#
# Build the engines
#

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_BUILD_VERSION} AS builder
RUN apk add --no-cache make git gcc binutils libc-dev

ARG CANOPSIS_EDITION
WORKDIR /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/

# Do a first pass that only handles Go dependencies, so that we can cache this part
COPY community/go-engines-community/go.mod /monorepo/community/go-engines-community/go.mod
COPY community/go-engines-community/go.sum /monorepo/community/go-engines-community/go.sum
COPY community/go-engines-community/Makefile /monorepo/community/go-engines-community/Makefile
COPY community/go-engines-community/.env /monorepo/community/go-engines-community/.env
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

COPY community/go-engines-community/config /opt/canopsis/share/config

ARG CANOPSIS_EDITION
ARG CMD
ENV CMD ${CMD}

RUN ln -sf /${CMD} /opt/canopsis/bin/${CMD}
RUN case "${CMD}" in canopsis-api*) ln -sf /canopsis-api-${CANOPSIS_EDITION} /canopsis-api && ( cd /opt/canopsis/bin && ln -sf canopsis-api-${CANOPSIS_EDITION} canopsis-api ) ;; esac
COPY --from=builder /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/build/cmd/${CMD} /${CMD}

USER canopsis:canopsis

CMD ["/bin/sh", "-c", "/${CMD}"]

#
# canopsis-reconfigure override
#

FROM final-container AS canopsis-reconfigure

USER root
COPY community/go-engines-community/cmd/canopsis-reconfigure/canopsis-community.toml /canopsis-community.toml
COPY community/go-engines-community/cmd/canopsis-reconfigure/canopsis-pro.toml /canopsis-pro.toml
RUN mkdir -p /opt/canopsis/share/migrations/mongodb /opt/canopsis/share/migrations/postgres
COPY community/go-engines-community/database/migrations/ /opt/canopsis/share/migrations/mongodb/
COPY community/go-engines-community/database/postgres_migrations/ /opt/canopsis/share/migrations/postgres/
USER canopsis:canopsis

ARG CMD
ENV CMD ${CMD}

# note: CPS_EDITION is defined at Docker Compose runtime in the .env file
CMD ["/bin/sh", "-c", "/${CMD} -conf /canopsis-${CPS_EDITION:-community}.toml -postgres-migration-directory=/opt/canopsis/share/migrations/postgres -migrate-postgres=true -postgres-migration-mode=up"]

#
# engine-che override (Community and Pro)
#

FROM final-container AS engine-che

USER root
RUN mkdir -p /opt/canopsis/lib/go/plugins
# XXX: horrible work around the fact that COPY can't be conditional, while mongo.so
# only exists when CANOPSIS_EDITION is 'pro', so we use a pattern that will always
# copy something AND mongo.so by accident if we're on Pro (sigh). The best solution
# would be to move the condition to the 'make install' target, and just use that.
# But conflicting Community/Pro binaries were added even though we warned that the
# build system wasn't ready for this yet...
COPY --from=builder /monorepo/${CANOPSIS_EDITION}/go-engines-${CANOPSIS_EDITION}/build/*/*ce* /opt/canopsis/lib/go/plugins/
RUN find /opt/canopsis/lib/go/plugins \( \! -name mongo.so \) -type f -delete
USER canopsis:canopsis

#
# engine-webhook override
#

FROM final-container AS engine-webhook

USER root
COPY pro/go-engines-pro/cmd/engine-webhook/webhook.conf.toml /opt/canopsis/etc/webhook.conf.toml
USER canopsis:canopsis

#
# import-context-graph override
#

FROM final-container AS import-context-graph

USER root
COPY pro/go-engines-pro/config/import-context-graph /opt/canopsis/share/config/import-context-graph
USER canopsis:canopsis

#
# canopsis-api-* overrides
#

FROM final-container AS canopsis-api-community
EXPOSE 8082

FROM final-container AS canopsis-api-pro
EXPOSE 8082
