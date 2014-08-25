Install from sources
====================

Requirements
------------

Install requirements with ``root`` user.

Debian Like:
^^^^^^^^^^^^

.. code-block:: bash

    apt-get update
    apt-get install sudo git-core libcurl4-gnutls-dev libncurses5-dev

Redhat Like:
^^^^^^^^^^^^

Disable some services

.. code-block:: bash

    ## Disable SELinux and Firewall
    setenforce 0
    chkconfig iptables off
    chkconfig ip6tables off
    chkconfig qpidd off
    service iptables stop
    service ip6tables stop
    service qpidd stop

Install some packages

.. code-block:: bash

    yum install wget make redhat-lsb gcc gcc-c++ zlib-devel ncurses-devel git

Install xorg-x11-server-Xvfb

.. code-block:: bash

    yum install xorg-x11-server-Xvfb

If package ``xorg-x11-server-Xvfb 1.10.4-6`` not found, you can try:

.. code-block:: bash

    wget http://vault.centos.org/6.2/os/x86_64/Packages/xorg-x11-server-Xvfb-1.10.4-6.el6.x86_64.rpm
    yum localinstall xorg-x11-server-Xvfb-1.10.4-6.el6.x86_64.rpm

Download sources
----------------

Clone git repository:

.. code-block:: bash

    git clone https://github.com/capensis/canopsis.git
    cd canopsis
    git submodule init
    git submodule update

You can switch to development branch:

.. code-block:: bash

    git checkout develop

Build and install
-----------------

.. code-block:: bash

    cd sources
    sudo ./build-install.sh

If build failed, you can see logs in ``log/`` directory.

Start Canopsis
--------------

Log in ``canopsis`` and start it:

.. code-block:: bash

    sudo su - canopsis
    hypcontrol start

Check installation
------------------

You can verify installation: (in ``canopsis`` environment)

.. code-block:: bash

    python opt/canotools/functional-test.py

Trouble shooting
----------------

During some occasions, you could encounter some funny error messages such as :

* Supervisord still running

.. code-block:: bash

    unix:///opt/canopsis/tmp/supervisor.sock no such file

This error occurs when ``supervisord`` failed to start during the installation. Simply start it in a ``canopsis`` environement.

* Erlang refuses to work and crashes

.. code-block:: bash

    Crash dump was written to: erl_crash.dump
    Kernel pid terminated (application_controller) ({application_start_failure,kernel,{shutdown,{kernel,start,[normal,[]]}}})
    + Declare Admin user ...
    {error_logger,{{2014,4,28},{9,20,0}},"Error when reading /opt/canopsis/.erlang.cookie: eacces",[]}
    [...]

This error occurs when rabbit-ms configuration is not properly set next to a system crash or equivalent. It is possible to fix this issue by removing the erlang cookie in canopsis root folder ``rm /opt/canopsis/.erlang.cookie``. this may have no side effect when canopsis is in single instance mode (no HA)

