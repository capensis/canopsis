FROM debian:jessie-slim as canopsis-base

RUN echo "deb http://ftp.fr.debian.org/debian/ jessie non-free" >> /etc/apt/sources.list

RUN apt-get update \
  && apt-get -y install python \
    python-pip \
    build-essential \
    ssh \
    sudo \
    bash \
    bash-completion \
    rsync \
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
    libldap-dev \
    rsync \
    snmp \
    smitools \
    smistrip \
    wget \
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
    openssl \
    vim \
  && apt-get -y install -d snmp-mibs-downloader \
  && dpkg -x /var/cache/apt/archives/snmp-mibs-downloader_1.1_all.deb / \
  && download-mibs \
  && apt-get clean

FROM canopsis-base as canopsis-engines

# define environment variables
ENV canopsis_install_dir /opt/canopsis
ENV canopsis_user canopsis
ENV canopsis_group canopsis

# create the canopsis environment
COPY files/sudoers /etc/sudoers.d/canopsis
RUN groupadd ${canopsis_group} && useradd -d ${canopsis_install_dir} -m -g ${canopsis_group} -s /bin/bash ${canopsis_user}

COPY files/bashrc ${canopsis_install_dir}/.bashrc

RUN mkdir -p ${canopsis_install_dir}/var/log/engines \
  && mkdir -p ${canopsis_install_dir}/var/cache/canopsis \
  && mkdir -p ${canopsis_install_dir}/var/lib/canopsis/unittest \
  && mkdir -p ${canopsis_install_dir}/var/run \
  && mkdir -p ${canopsis_install_dir}/tmp \
  && mkdir -p ${canopsis_install_dir}/etc/init.d \
  && mkdir -p ${canopsis_install_dir}/repo \
  && mkdir -p ${canopsis_install_dir}/repo/canopsis-externals/

COPY sources/db-conf/opt/ ${canopsis_install_dir}/opt/
COPY sources/externals/python-libs ${canopsis_install_dir}/repo/canopsis-externals/python-libs

COPY sources/extra/conf/supervisord.conf ${canopsis_install_dir}/etc/supervisord.conf
RUN virtualenv ${canopsis_install_dir}/
RUN echo "[easy_install]\nallow_hosts = ''\nfind_links = file://${canopsis_install_dir}/repo/canopsis-externals/python-libs/" > /root/.pydistutils.cfg

RUN . ${canopsis_install_dir}/bin/activate \
  && pip install --no-index --find-links=file://${canopsis_install_dir}/repo/canopsis-externals/python-libs/ --upgrade setuptools distribute

COPY sources/canopsis ${canopsis_install_dir}/repo/canopsis

# TODO modifier
RUN . ${PREFIX}/bin/activate \
  && pip install --no-index --find-links=file://${PREFIX}/repo/canopsis-externals/python-libs/ ${PREFIX}/repo/canopsis

RUN rsync -avKSH  ${PREFIX}/repo/canopsis/etc/ ${PREFIX}/etc/

ADD files/amqp.conf ${PREFIX}/etc
ADD files/cstorage.conf ${PREFIX}/etc
ADD files/mongo/storage.conf ${PREFIX}/etc/mongo
ADD files/mongo/mongo_store.conf ${PREFIX}/etc/common
ADD files/influx/storage.conf ${PREFIX}/etc/influx

RUN . ${PREFIX}/bin/activate \
  && cd ${PREFIX}/repo/canopsis-cat/sources/canopsis_cat/ && pip install --no-index --find-links=file://${PREFIX}/repo/canopsis-externals/python-libs/ .

RUN rsync -avKSH  ${PREFIX}/repo/canopsis-cat/sources/canopsis_cat/etc/ ${PREFIX}/etc/

RUN rm -rf ${PREFIX}/repo/canopsis-externals/python-libs

RUN chown -R ${CANOPSIS_USER}:${CANOPSIS_GROUP} ${PREFIX}
