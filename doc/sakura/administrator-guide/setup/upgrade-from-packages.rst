.. include:: /includes/links.rst

Upgrade from packages
=====================

Introduction
------------

Before upgrade, it's recommended to make a snapshot or backup your Canopsis with this: |howToBackUp|

Upgrade Ubik
------------

As ``root``:

.. code-block:: bash

	pip install --upgrade git+https://github.com/socketubs/ubik.git@0.1

	rm /opt/canopsis/bin/ubik
	ln -s $(which ubik) /opt/canopsis/bin/ubik

Upgrade Canopsis
----------------

Log in as Canopsis:

.. code-block:: bash

	su - canopsis

Stop Canopsis
~~~~~~~~~~~~~

.. code-block:: bash

	hypcontrol stop


Upgrade Python and Python-libs
~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

.. code-block:: bash

	ubik update
	ubik upgrade python python-libs


Upgrade Canopsis
~~~~~~~~~~~~~~~~

.. code-block:: bash

	ubik update
	ubik upgrade

Start Canopsis
~~~~~~~~~~~~~~

.. code-block:: bash

	hypcontrol start
