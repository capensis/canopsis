#!/bin/sed -f
s|{{BASE_DOCKER_IMAGE}}|debian:jessie-slim|g
s|{{BUILD_SETUP}}|apt-get update; \\ \
	apt-get -y install apt-utils; \\ \
	apt-get -y --no-install-recommends install \\ \
		build-essential \\ \
		curl \\ \
		libcurl4-openssl-dev \\ \
		libsasl2-dev \\ \
		libxml2-dev \\ \
		libxslt1-dev \\ \
		libssl-dev \\ \
		libffi-dev \\ \
		libxmlsec1-dev \\ \
		libxmlsec1-openssl \\ \
		libldap2-dev \\ \
		pkg-config \\ \
		python2.7-dev \\ \
		python-pip \\ \
		python-pkg-resources \\ \
		python-virtualenv \\ \
		python-wheel \\ \
		virtualenv; \\ \
	apt-get clean|g
s|{{FINAL_IMAGE_SETUP}}|echo "deb http://ftp.fr.debian.org/debian/ jessie main contrib non-free" > /etc/apt/sources.list ; \\ \
	echo "deb http://security.debian.org/ jessie/updates main" >> /etc/apt/sources.list; \\ \
	# \
	# Keep the blank comments. It's used to make the file more readable and \
	# it prevent docker to complain about an empty line. \
	# \
	rm -f /etc/localtime; \\ \
	ln -s /usr/share/zoneinfo/UTC /etc/localtime; \\ \
	# \
	apt-get update; \\ \
	apt-get dist-upgrade -y; \\ \
	apt-get -y --no-install-recommends install locales; \\ \
	# \
	export LANG="en_US.UTF-8"; \\ \
	echo "LANG=${LANG}" > /etc/locale.conf; \\ \
	echo "${LANG} UTF-8" > /etc/locale.gen; \\ \
	locale-gen; \\ \
	# \
	apt-get -y --no-install-recommends install \\ \
		apt-transport-https \\ \
		base-files \\ \
		bash \\ \
		ca-certificates \\ \
		curl \\ \
		dnsutils \\ \
		iputils-ping \\ \
		libsasl2-2 \\ \
		libxml2 \\ \
		libxslt1.1 \\ \
		lsb-base \\ \
		libffi6 \\ \
		libgmp10 \\ \
		libgnutlsxx28 \\ \
		libgnutls-openssl27 \\ \
		libhogweed2 \\ \
		libicu52 \\ \
		libidn11 \\ \
		libnettle4 \\ \
		libp11-kit0 \\ \
		libpsl0 \\ \
		libssl1.0.0 \\ \
		libtasn1-6 \\ \
		libxmlsec1 \\ \
		libxmlsec1-openssl \\ \
		libldap-2.4-2 \\ \
		python2.7 \\ \
		rsync \\ \
		smitools \\ \
		sudo \\ \
		vim; \\ \
		apt-get clean \
		# python \\ \
		# tmux \\ \
		# snmp|g
