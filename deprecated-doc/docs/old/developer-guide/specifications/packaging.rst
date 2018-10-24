.. _dev-spec-buildinstall:

Build System Specification
==========================

The Canopsis build-system is composed of a set of shell scripts, which builds,
installs and updates Canopsis packages inside the environment.

They also generates Ubik packages to make a binary repository of Canopsis.

Entry Point
-----------

All starts with the ``build-install.sh`` script :

.. code-block:: text

   Usage : ./build-install.sh [OPTIONS]
   
        Build and install Canopsis dependencies and packages
   
   Options:
       -c       -> Uninstall Canopsis
       -n       -> Don't build sources if possible
       -u       -> Run unittest at the end
       -p       -> Make packages
       -d       -> Don't check host dependencies
       -i       -> Build Canopsis Installer (for packages installation)
       -h, help -> Print this help

You can't run this script if you're not logged in as root.

Configure installation
----------------------

The first thing it will do, is trying to load the ``common.sh`` script from the
**canohome** package (``sources/canohome/lib/common.sh``).

This file contains functions and environment variables for the ``build-install.sh``
and the environment once installed :

.. code-block:: bash

   PREFIX="/opt/canopsis"
   HUSER="canopsis"
   HGROUP="canopsis"

By editing those variables, you are able to configure the installation folder of
Canopsis, for example :

.. code-block:: bash

   PREFIX="/opt/canopsis/ficus"
   HUSER="canopsis-ficus"
   HGROUP="canopsis"

And :

.. code-block:: bash

   PREFIX="/opt/canopsis/sakura"
   HUSER="canopsis-sakura"
   HGROUP="canopsis"

Host dependencies
-----------------

After initializing environment, it will try to install the build
dependencies on the host.

To do that, the function ``detect_os`` and the variable ``ARCH`` from the ``common.sh``
script will be used to detect the actual platform.

This will allow the ``build-install.sh`` to execute the correct script from ``sources/extra/dependencies``.

For example, if I'm on Debian Wheezy, the file will be ``debian_7``.
But if I'm on CentOS 6.5, the file will be ``centos_6``.

Building Canopsis packages
--------------------------

Canopsis packages are an extended model of an Ubik package. They are used to build
the Canopsis component, installing it in a Canopsis prefix, and generate Ubik
packages from it.

First, if a script ``sources/pre-build.sh`` exists, it will be executed.

Then each files, ending with ``.install`` suffix, in the ``sources/build.d`` folder,
will be loaded.

They provides the Canopsis package build instructions :

.. code-block:: bash

   NAME="canohome"
   
   function build(){
       true
   }
   
   function install(){
       install_basic_source $NAME
   }
   
   function update(){
       update_basic_source $NAME
   }

Each functions available in the ``build-install.sh`` script are available in those
files.

For each of those files, representing a single Canopsis packages, we will load
the file ``sources/packages/$NAME/control`` where ``$NAME`` is defined in the
``build.d`` script.

This ``control`` file respects the Ubik format, only the functions defined in the
``common.sh`` file are available :

.. code-block:: bash

   #!/bin/bash

   NAME="canohome"
   VERSION=0.6
   RELEASE=0
   DESCRIPTION=""
   REQUIRES=""
   
   NO_ARCH=true
   NO_DIST=true
   NO_DISTVERS=true
   
   function pre_install() {
       echo "Pre-install $NAME $VERSION-$RELEASE ..."
       # ...
   }
   
   function post_install() {
       echo "Post-install $NAME $VERSION-$RELEASE ..."
       # ...
   }
   
   function pre_remove() {
       echo "Pre-remove $NAME $VERSION-$RELEASE ..."
   }
   
   function post_remove() {
       echo "Post-remove $NAME $VERSION-$RELEASE ..."
       # ...
   }
   
   function pre_update() {
       echo "Pre-update $NAME $VERSION-$RELEASE ..."
   }
   
   function post_update() {
       echo "Post-update $NAME $VERSION-$RELEASE ..."
   }
   
   function purge() {
       echo "Purge $NAME $VERSION-$RELEASE ..."
       # ...
   }

If this file doesn't exist, a dummy one is created.

With the informations given with the ``control`` file, we're now able to define
``P_ARCH``, ``P_DIST`` and ``P_DISTVERS``.

If the file ``$PREFIX/var/lib/pkgmgr/packages/$NAME.info`` exists, then the
``build-install.sh`` script will perform an update (unless we are building MongoDB).

The process will be :

.. code-block:: bash

   pre_update
   update
   post_update
   echo "v${VERSION}-r${RELEASE}_${P_DIST}-${P_DISTVERS}_${P_ARCH}" > $VARLIB_PATH/$NAME.info

In the other case, we will check if we need to build the package.
If we need to, the function ``build`` will be called.
Then the process will be :

.. code-block:: bash

   pre_install
   install
   post_install
   echo "v${VERSION}-r${RELEASE}_${P_DIST}-${P_DISTVERS}_${P_ARCH}" > $VARLIB_PATH/$NAME.info

Now, the file ``$PREFIX/var/lib/pkgmgr/packages/$NAME.files`` will be generated.
It contains the listing of all the files created in the Canopsis prefix during
the installation.

This file will be used to generated the Ubik package in ``binaries`` folder.
After the Ubik package generation, the package's metadata are saved in ``sources/pkg.tmp``,
which will be converted into ``binaries/Packages.json`` at the end of the Canopsis
installation.

When the ``build-install.sh`` script is ended, if packages generation was asked,
there is an Ubik repository in the folder ``binaries``.

You can serve it via HTTP, and use Ubik to install Canopsis from it :

.. code-block:: text

   $ python -m SimpleHTTPServer

After the whole build process, if the script ``sources/post-build.sh`` exists,
it will be executed.

Packages configuration update
-----------------------------

In case of an update, the configuration files stored in ``etc`` are handled differently.

For each files, we will check two things :

 * if the file doesn't exist in the prefix: the file is copied
 * if the file exists, and haven't the same md5 checksum: we ask the user if we can replace it

Ubik package generation
-----------------------

The package is generated in the folder ``PPATH=binaries/$P_ARCH/$P_DIST/$P_DISTVERS/$NAME``.

The file ``$PREFIX/var/lib/pkgmgr/packages/$NAME.files`` is used to generate the
tar archive : ``$PPATH/files.bz2`` and the file ``$PPATH/files.lst``

Then the files ``sources/packages/$NAME/control`` and ``sources/packages/$NAME/blacklist``
are copied to the ``$PPATH`` folder, which is then archived under the name ``$NAME.tar``,
before being removed.
