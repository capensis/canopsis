ARG TAG

FROM canopsis/canopsis-sysbase:${TAG} as build

ARG TAG
ARG PROXY

ENV http_proxy ${PROXY}
ENV https_proxy ${PROXY}

COPY docker/build/debian-jessie-build.sh /
COPY docker/build/install.sh /
RUN /bin/bash /debian-jessie-build.sh

WORKDIR ${CPS_HOME}

### Sources

RUN mkdir -p /sources/bricks/

COPY sources/canopsis /sources/canopsis
COPY sources/externals/python-libs /sources/python-libs
COPY sources/webcore /sources/webcore
COPY tools/brickmanager /sources/brickmanager

### Canopsis Core and Engines
RUN /bin/bash /install.sh

### Webserver

RUN mkdir -p ${CPS_HOME}/var/www && rsync -a --exclude=/sources/webcore/doc /sources/webcore/ ${CPS_HOME}/var/www/
COPY docker/build/bricks.sh /
RUN /bin/bash /bricks.sh ${TAG} ${CPS_HOME}

## Canopsis Core Image
FROM canopsis/canopsis-sysbase:${TAG}

ARG PROXY

ENV http_proxy ${PROXY}
ENV https_proxy ${PROXY}
ENV CPS_WEBSERVER 0

COPY --from=build /${CPS_HOME} /${CPS_HOME}

WORKDIR ${CPS_HOME}

COPY sources/canopsis/etc ./etc
COPY sources/db-conf/opt ./opt
COPY docker/files/sudoers /etc/sudoers.d/canopsis
COPY docker/files/bashrc .bashrc
COPY docker/files/cstorage.conf ./etc
COPY docker/files/amqp.conf ./etc
COPY docker/files/mongo/* ./etc/mongo/
COPY docker/files/common/* ./etc/common/

RUN sed -r "s@~@${CPS_HOME}@g" -i ./etc/webserver.conf

ADD docker/files/entrypoint.sh /

RUN mkdir -p ./etc/init.d && \
    mkdir -p var/log/engines && \
    mkdir -p var/cache/canopsis && \
    mkdir -p tmp

RUN chown -R ${CPS_USER}:${CPS_GROUP} ${CPS_HOME}/var ${CPS_HOME}/tmp

USER ${CPS_USER}:${CPS_GROUP}
EXPOSE 8082
ENTRYPOINT /entrypoint.sh
