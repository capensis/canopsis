Packaging
=========

Build Process
-------------

The Canopsis build process is split in 3 parts :

 * host dependencies
 * environment dependencies
 * Canopsis packages

This process is done by the ``build-install.sh`` script.

Firstly, it detects the actual system and execute the corresponding script located
at ``sources/extra/dependencies`` in the source tree.

Then, it creates (if needed) the Canopsis environment and execute all scripts located
in ``sources/build.d`` in the source tree.
Those scripts are executed ordered by alphabetical order.

Finally, the build process list the available packages in the ``sources/packages``
folder, and execute the ``control`` file located in the package's folder.

Package control file
--------------------

Here is an example of a *control* file (without functions implementations) :

.. code-block:: bash

	#!/bin/bash

	# Package's name
	NAME="amqp2engines"

	# Distributed version
	VERSION=0.7

	# Release of the package (each time you update the package, increment the counter)
	RELEASE=0

	# Description and dependencies
	DESCRIPTION=""
	# Each names refers to another Canopsis package
	REQUIRES="canohome python canolibs supervisord-conf pyperfstore2"

	# Platform informations
	# If not defined, or set to false, ARCH, DIST, and DISTVERS will be automatically
	# detected and set to the build platform values
	NO_ARCH=true
	NO_DIST=true
	NO_DISTVERS=true

	function pre_install() {
		# ...
	}

	function post_install() {
		# ...
	}

	function pre_remove() {
		# ...
	}

	function post_remove() {
		# ...
	}

	function pre_update() {
		# ...
	}

	function post_update() {
		# ...
	}

	function purge() {
		# ...
	}

Generating packages
-------------------

Similar to the ``build-install.sh`` script, there is a ``build-packages.sh``.

This script first removes the existing Canopsis environment, then clone the Canopsis
repository to ``/usr/local/src/canopsis``.

The default cloned branch is ``freeze``, pass the name of the branch to clone as
first parameter of the ``build-packages.sh`` script.

After the source tree is fetched, the script runs the build process and track all
installations to generate packages in ``/usr/local/src/canopsis/binaries``.

Update the repository
---------------------

The ``/usr/local/src/canopsis/binaries`` acts as a single-release repository.

Here is an example of repository generated on a **Debian 7 - x86_64** :

 * Packages.json
 * canopsis_installer.tgz
 * noarch
    * nodist
       * novers
          * amqp2engines.tar
          * ...
 * x86_64
    * debian
       * 7
          * collectd.tar
          * ...

Now, let's explore the ``Packages.json`` file :

.. code-block:: javascript

	[
		{
			"arch":              // package's architecture or noarch
			"dist":              // package's distribution or nodist
			"vers":              // package's distribution version or novers
			"name":              // package's name
			"description":       // package's description
			"requires":          // list of required packages
			"version":           // distributed version
			"release":           // package's version
			"md5":               // archive's md5
		},
		// ...
	]

If you want generate your repository for many distributions, be sure to update the
``Packages.json`` file correctly.