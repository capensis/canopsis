Install from packages
=====================

Requirements
------------

**Debian 6, Ubuntu 12.04 LTS and CentOS 6 64 bits only.**

Debian Like
^^^^^^^^^^^

.. code-block:: bash

    apt-get update
    apt-get install sudo uuid-runtime curl xvfb bzip2 libbz2-dev bc libevent-dev libxrender1 libfontconfig1 libltdl7 python-pip git-core python-gevent


.. caution:: Debian 6 and Ubuntu 10.04 Only, upgrade Pip

.. code-block:: bash

    easy_install pip
    rm /usr/bin/pip
    ln -sv /usr/local/bin/pip-2.6 /usr/bin/pip
    pip install pip --upgrade

Redhat / CentOS
^^^^^^^^^^^^^^^
Disable some services

.. code-block:: bash

    #Disable SELinux and Firewall
    sed -i "s#enforcing#disabled#" /etc/selinux/config
    setenforce 0
    chkconfig iptables off
    chkconfig ip6tables off
    chkconfig qpidd off
    service iptables stop
    service ip6tables stop
    service qpidd stop

Install some packages

.. code-block:: bash

    yum install wget make redhat-lsb gcc gcc-c++ zlib-devel ncurses-devel git python-setuptools libevent
    easy_install pip

Install xorg-x11-server-Xvfb

.. code-block:: bash

    yum install xorg-x11-server-Xvfb

If package ``xorg-x11-server-Xvfb 1.10.4-6`` is not found, you may try:

.. code-block:: bash

    wget http://vault.centos.org/6.2/os/x86_64/Packages/xorg-x11-server-Xvfb-1.10.4-6.el6.x86_64.rpm
    yum localinstall xorg-x11-server-Xvfb-1.10.4-6.el6.x86_64.rpm

Environment
-----------

.. code-block:: bash

    pip install --upgrade git+https://github.com/socketubs/ubik.git@0.1
    useradd -m -d /opt/canopsis -s /bin/bash canopsis


Installer
---------

Login with Canopsis
^^^^^^^^^^^^^^^^^^^

.. code-block:: bash

    sudo su - canopsis

    # You can export your HTTP proxy configuration
    export http_proxy="http://<USER>:<PASS>@<SERVER>:<PORT>"
    export https_proxy=$http_proxy

Stable
^^^^^^

Download and install the installer (as ``canopsis`` user):

.. code-block:: bash

    mkdir tmp && cd tmp
    wget http://repo.canopsis.org/stable/canopsis_installer.tgz
    tar xfz canopsis_installer.tgz
    cd canopsis_installer

If you have a Proxy, you must edit `ubik.conf` before installation.
After installation `ubik.conf` is copied into `~/etc/ubik.conf`.

.. code-block:: bash

    ./install.sh
    exit

Daily
^^^^^

Download and install the installer (as ``canopsis`` user):

.. code-block:: bash

    mkdir tmp && cd tmp
    wget http://repo.canopsis.org/daily/canopsis_installer.tgz
    tar xfz canopsis_installer.tgz
    cd canopsis_installer

If you have a Proxy, you must edit `ubik.conf` before installation.
After installation `ubik.conf` is copied into `~/etc/ubik.conf`.

.. code-block:: bash

    ./install.sh daily
    exit

Install Main packages
---------------------

Install package ``cmaster`` (as ``canopsis`` user):

.. code-block:: bash

    sudo su - canopsis

    ubik update
    ubik list

    ubik install cmaster

Start Canopsis
--------------

.. code-block:: bash

    hypcontrol start

Now you can login on WebUI: `http://Your_IP:8082` (Login: `root`, Password: `root`)

Enjoy ;)