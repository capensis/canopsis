ARG CANOPSIS_TAG
ARG CANOPSIS_DISTRIBUTION

FROM canopsis/canopsis-core:${CANOPSIS_DISTRIBUTION}-${CANOPSIS_TAG}
ADD docker/files/entrypoint-prov-sync.sh /
ADD docker/files/entrypoint-prov.sh /

VOLUME /provisioning

ENTRYPOINT /entrypoint-prov.sh
