FROM node:10.6-stretch

VOLUME /dist/

COPY sources/webcore/src/canopsis-next /canopsis-next
COPY docker/build/canopsis-next.sh /build.sh

ENTRYPOINT ["/build.sh"]
