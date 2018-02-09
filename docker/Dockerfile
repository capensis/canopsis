FROM debian:jessie-slim as canopsis-base

RUN echo "deb http://ftp.fr.debian.org/debian/ jessie non-free" >> /etc/apt/sources.list

RUN apt-get update \
  && apt-get -y install python \
    python-pip \
    tmux \
    build-essential \
    bash \
    bash-completion \
    apt-transport-https \
    python2.7-dev \
    python-virtualenv \
    libcurl4-openssl-dev \
    libsasl2-dev \
    libxml2-dev \
    libxslt1-dev \
    base-files \
    lsb-core \
    lsb \
    libssl-dev \
    patch \
    ca-certificates \
    libffi6 \
    libffi-dev \
    libgmp10 \
    libgnutls-deb0-28 \
    libhogweed2 \
    libicu52 \
    libidn11 \
    libnettle4 \
    libp11-kit0 \
    libpsl0 \
    libssl1.0.0 \
    libtasn1-6 \
    libxmlsec1 \
    libxmlsec1-dev \
    libldap2-dev \
  && apt-get clean

# CANOPSIS ENGINES
FROM canopsis-base as canopsis-engines

ENV CPS_HOME /opt/canopsis
ENV CPS_USER canopsis
ENV CPS_GROUP canopsis

RUN groupadd ${CPS_GROUP} && useradd -d ${CPS_HOME} -m -g ${CPS_GROUP} -s /bin/bash ${CPS_USER}

RUN mkdir -p ${CPS_HOME}/var/log/engines \
  && mkdir -p ${CPS_HOME}/var/cache/canopsis \
  && mkdir -p ${CPS_HOME}/var/lib/canopsis/unittest \
  && mkdir -p ${CPS_HOME}/var/run \
  && mkdir -p ${CPS_HOME}/tmp \
  && mkdir -p ${CPS_HOME}/repo/canopsis-externals/python-libs \
  && mkdir -p ${CPS_HOME}/repo/canopsis

COPY docker/files/sudoers /etc/sudoers.d/canopsis
COPY docker/files/bashrc ${CPS_HOME}/.bashrc
COPY docker/files/vimrc ${CPS_HOME}/.vimrc
COPY sources/db-conf/opt/ ${CPS_HOME}/opt/
COPY sources/externals/python-libs/ ${CPS_HOME}/repo/canopsis-externals/python-libs/
COPY sources/canopsis/ ${CPS_HOME}/repo/canopsis/

RUN virtualenv ${CPS_HOME}/
RUN echo "[easy_install]\nallow_hosts = ''\nfind_links = file://${CPS_HOME}/repo/canopsis-externals/python-libs/" > /root/.pydistutils.cfg

RUN . ${CPS_HOME}/bin/activate \
  && pip install --no-index --find-links=file://${CPS_HOME}/repo/canopsis-externals/python-libs/ --upgrade setuptools distribute

# TODO modifier
RUN . ${CPS_HOME}/bin/activate \
  && pip install --no-index --find-links=file://${CPS_HOME}/repo/canopsis-externals/python-libs/ -r ${CPS_HOME}/repo/canopsis/requirements.txt ${CPS_HOME}/repo/canopsis \
  && rm -rf ${CPS_HOME}/repo/canopsis-externals


COPY sources/canopsis/etc ${CPS_HOME}/etc
#RUN rsync -avKSH  ${CPS_HOME}/repo/canopsis/etc/ ${CPS_HOME}/etc/

COPY docker/files/amqp.conf ${CPS_HOME}/etc
RUN  mkdir -p ${CPS_HOME}/etc/init.d

COPY docker/files/cstorage.conf ${CPS_HOME}/etc
COPY docker/files/mongo/storage.conf ${CPS_HOME}/etc/mongo
COPY docker/files/mongo/mongo_store.conf ${CPS_HOME}/etc/common
#COPY sources/canopsis/etc/influx/storage.conf ${CPS_HOME}/etc/influx

#RUN . ${CPS_HOME}/bin/activate \
# && cd ${CPS_HOME}/repo/canopsis-cat/sources/canopsis_cat/ && pip install --no-index --find-links=file://${CPS_HOME}/repo/canopsis-externals/python-libs/ .

#RUN rsync -avKSH  ${CPS_HOME}/repo/canopsis-cat/sources/canopsis_cat/etc/ ${CPS_HOME}/etc/


RUN chown -R ${CPS_USER}:${CPS_GROUP} ${CPS_HOME}

ADD docker/files/entrypoint.sh /
ENTRYPOINT ["/entrypoint.sh"]
