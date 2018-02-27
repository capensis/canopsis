ARG TAG

FROM canopsis/canopsis-sysbase:${TAG} as build

ARG TAG
ARG PROXY

ENV http_proxy ${PROXY}
ENV https_proxy ${PROXY}

COPY docker/build/debian-build.sh /
RUN /bin/bash /debian-build.sh

WORKDIR ${CPS_HOME}

### Canopsis Core

#### Venv setup
COPY sources/externals/python-libs /sources/python-libs
COPY docker/build/pip-setup.sh /
RUN /bin/bash /pip-setup.sh

#### Dependencies
COPY sources/canopsis/requirements.txt /sources/canopsis/requirements.txt
COPY docker/build/pip-deps.sh /
RUN /bin/bash /pip-deps.sh

#### Canopsis only
COPY sources/canopsis /sources/canopsis
COPY docker/build/pip-canopsis.sh /
RUN /bin/bash /pip-canopsis.sh

### Webserver
COPY sources/webcore /sources/webcore
COPY tools/brickmanager /sources/brickmanager
RUN mkdir -p ${CPS_HOME}/var/www && rsync -a --exclude=/sources/webcore/doc /sources/webcore/ ${CPS_HOME}/var/www/



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
COPY docker/files/etc/ ./etc/
COPY tools/brickmanager /brickmanager

ADD docker/files/entrypoint.sh /

RUN mkdir -p ./etc/init.d && \
    mkdir -p var/log/engines && \
    mkdir -p var/cache/canopsis && \
    mkdir -p tmp && \
    sed -r "s@~@${CPS_HOME}@g" -i ./etc/webserver.conf && \
    chown -R ${CPS_USER}:${CPS_GROUP} ${CPS_HOME}/

# Copy UI bricks as root
COPY docker/bricks/ ${CPS_HOME}/var/www/src/canopsis/
RUN . ${CPS_HOME}/bin/activate && \
    python /brickmanager enable brick-querybuilder

USER ${CPS_USER}:${CPS_GROUP}

EXPOSE 8082
ENTRYPOINT /entrypoint.sh
