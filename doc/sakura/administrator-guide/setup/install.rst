.. _admin-setup-install:

Install from sources
====================

Requirements
------------

Install requirements with ``root`` user.

Systems:
^^^^^^^

Canopsis can be installed on the following systems :

* Debian 7, 8
* Ubuntu 12, 14 (can be installed but not supported by team)
* RedHat / CentOS 7

Debian Like:
^^^^^^^^^^^^

.. code-block:: bash

    apt-get update
    apt-get install sudo git-core libcurl4-gnutls-dev libncurses5-dev

Redhat Like:
^^^^^^^^^^^^

Disable some services   
 We don't provide any SELinux context so it's better to disable it. Feel free to help us writing one !  
 You can see :ref:`admin-setup-firewall` to configure Iptables

.. code-block:: bash

    ## Disable SELinux and Firewall
    setenforce 0
    chkconfig iptables off
    chkconfig ip6tables off
    chkconfig qpidd off
    service iptables stop
    service ip6tables stop
    service qpidd stop

Iptables an qpidd may not be available on RedHat/CentOS 7. Take a look at firewalld

Install some packages

.. code-block:: bash

    yum install wget make redhat-lsb gcc gcc-c++ zlib-devel ncurses-devel git


Download sources
----------------

Clone git repository:

.. code-block:: bash

    git clone https://git.canopsis.net/canopsis/canopsis.git
    cd canopsis


Available branches are:

* master: the more stable version of canopsis Sakura
* rc: the future of the master version, contains new features. This branch should be almost stable
* develop: nightly build with latests complete feature. This branch can carry bugs that have to be fixed for new rc release.

.. code-block:: bash

    git submodule init
    git submodule update


Build and install
-----------------

.. code-block:: bash

    sudo ./build-install.sh

If build failed, you can see logs in ``log/`` directory.

Note that install dir will be /opt/canopsis by default.
You can change it by editing SOURCE_PATH/sources/canohome/lib/common.sh

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

And you can also check needed services :(in ``canopsis`` environment)

.. code-block:: bash

    hypcontrol status

Troubleshooting
---------------

During some occasions, you could encounter some funny errors.  
Please have a look at :ref:`Troubleshooting page <admin-troubleshooting>`.

