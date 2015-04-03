How to develop a cat module for canopsis
========================================

In this how to, we want to create a new cat module for canopsis that is made of various component types. We'll call it **testmodule** and it will have both a backend and frontend elements.

This documents is about how to develop a canopsis CAT module that is integrable to an canopsis install. This will covers topics such as packaging, setup, backend module, frontend module and lifecycle.


packaging
---------

First create a ``canopsis-cat/sources/python/testmodule`` folder. this is the place where source code will live in the repository. In this folder, respecting the canopsis file arborescence to create custom files will make the **testmodule** setup merging folders between the cat sources and the canopsis install folder.

Thus, when a ``canopsis-cat/sources/python/testmodule/etc/myconffile.ini`` is created, on module install, this file will be copied to the folder (by default) ``/opt/canopsis/etc/myconffile.ini``. This structure works for folders such as ``etc, opt, var ...``.

frontend
--------

for frontend files to be integrated to the UI environment, the setup just copy ``canopsis-cat/sources/python/testmodule/var/www/canopsis/testmodule/init.js`` to ``/opt/canopsis/var/www/canopsis/connectors/init.js``. Any other file in the ``canopsis-cat/sources/python/testmodule/var/www/canopsis/testmodule/`` folder will be copied logically in the Canopsis install web folder.

The init.js file is the root file called by the UI environment and is in charge to load recursively the whole frontend **testmodule**. actually, all concerned files are in the ``canopsis-cat/sources/python/testmodule/var/www/canopsis/testmodule/`` folder that may contain various files type such as controllers, adapters, views, templates and such. All those files are part of a frontend module and will work the same was as they are structured the same way.

Finally, in order to enable the frontend module in the canopsis UI, it has to be registered in the enable frontend modules list. This must be done with the following command:

.. code-block:: bash

	#Available in the canopsis environment.
	webmodulemanager add testmodule

This operation may be executed on module setup, operation described in the setup section of this document.

backend
-------

The backend structure for a CAT module works with the same ides than the frontend one. With this in mind, creating backend elements can be done with the following rules:

 - write engine files in the ``canopsis-cat/sources/python/testmodule/canopsis/engines/testmodule.py`` file.
 - write webservices api in the ``canopsis-cat/sources/python/testmodule/canopsis/webcore/services/testmodule.py`` file.
 - put specific module's files in the ``canopsis-cat/sources/python/testmodule/canopsis/testmodule`` folder

This way, creating a backend module will keep the a clean backend on **testmodule** install with respect of canopsis conventions.

.. note:: do not forget **__init__.py** files in folders that contains python files to tell python the folder is a module.

setup
-----

When the ``build-install.sh -o testmodule -d path_to_canopsis_sources`` is executed, this will run setup scripts that are located in the 
Then create the following files ````


.. code-block:: bash


refer to the `administrator guide <../administration/amqp2engines.html>`_ for more information about this kind of configuration.

 