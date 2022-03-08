ARG CANOPSIS_TAG
ARG CANOPSIS_DISTRIBUTION



## Canopsis Build Image
FROM canopsis/canopsis-sysbase:${CANOPSIS_DISTRIBUTION}-${CANOPSIS_TAG} as build

ARG CANOPSIS_TAG
ARG CANOPSIS_DISTRIBUTION
ARG PROXY

ENV http_proxy ${PROXY}
ENV https_proxy ${PROXY}

COPY docker/build/build-${CANOPSIS_DISTRIBUTION}.sh /build.sh
RUN /bin/bash /build.sh

WORKDIR ${CPS_HOME}

### Canopsis Core

#### Venv setup
COPY docker/build/pip-setup.sh /
RUN /bin/bash /pip-setup.sh

#### Dependencies
COPY sources/canopsis/requirements.txt /sources/canopsis/requirements.txt
COPY docker/build/pip-deps.sh /
RUN /bin/bash /pip-deps.sh

#### Ansible
COPY docker/build/pip-ansible.sh /
RUN /bin/bash /pip-ansible.sh

#### Canopsis only
COPY sources/canopsis /sources/canopsis
COPY docker/build/pip-canopsis.sh /
RUN /bin/bash /pip-canopsis.sh
RUN ln -s ${CPS_HOME}/bin/canoctl /usr/bin/canoctl

### Webcore
COPY sources/webcore /sources/webcore
RUN mkdir -p ${CPS_HOME}/var/www && rsync -a /sources/webcore/ ${CPS_HOME}/var/www/

COPY VERSION.txt /${CPS_HOME}/

## Canopsis Core Image
FROM canopsis/canopsis-sysbase:${CANOPSIS_DISTRIBUTION}-${CANOPSIS_TAG}

ARG PROXY

ENV http_proxy ${PROXY}
ENV https_proxy ${PROXY}
ENV CPS_OLD_API 0

COPY --from=build /${CPS_HOME} /${CPS_HOME}

WORKDIR ${CPS_HOME}

RUN mkdir -p ./share/conf/samples

COPY sources/canopsis/etc/*.example ./share/conf/samples/
COPY sources/canopsis/etc ./etc
COPY sources/canopsis/etc/amqp2engines-core.conf.example ./etc/amqp2engines.conf
COPY sources/db-conf/opt ./opt
COPY docker/files/sudoers /etc/sudoers.d/canopsis
COPY docker/files/bashrc .bashrc
COPY docker/files/bash_profile .bash_profile
COPY docker/files/etc/ ./etc/
COPY sources/canopsis/canopsis/webcore/services/swagger/ ./lib/python2.7/site-packages/canopsis/webcore/services/swagger/

RUN chmod +x /opt/canopsis/.bash_profile

ADD docker/files/entrypoint.sh /

# Do NOT chown the entire CPS_HOME directory:
#  * Security: the user must not be able to change runtime files.
#  * Image size: docker is dumb, until --squash is stable.
RUN mkdir -p ./etc/init.d && \
    mkdir -p var/log/engines && \
    mkdir -p var/cache/canopsis && \
    mkdir -p tmp && \
    sed -i -E -e "s@~@${CPS_HOME}@g" ./etc/oldapi.conf && \
    chown -R ${CPS_USER}:${CPS_GROUP} ${CPS_HOME}/var/cache ${CPS_HOME}/var/log ${CPS_HOME}/tmp /opt/canopsis/etc

# Ansible
COPY deploy-ansible/ ${CPS_HOME}/deploy-ansible

USER ${CPS_USER}:${CPS_GROUP}

EXPOSE 8081
ENTRYPOINT /entrypoint.sh
