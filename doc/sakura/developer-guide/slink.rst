.. _dev-slink:

Slink
=====

When developping Canopsis features, It is necessary to test the code to build again the system. A tool named **slink** helps the developper to test quicker it's code.

This tool just replace some installation (/opt/canopsis) folders by some **symbolic links** pointing directly in the canopsis sources.

This way the developer just have to edit source code and he just need to restart the appropriate service (engine, webserver) to reload python code. The frontend side with slinks just require to reload the canopsis UI in the browser instead of build installing again.

The slink system can be installed by going to the canopsis source folder and run.

.. code-block:: bash

   sudo ./tools/slink

This command will replace folders in ``/opt/canopsis`` with symbolic links. To remove these links in the canopsis install just build-install again and ask yes to the question about removing symbolic links.


.. note::
   When using slink some folder may not be slinked properly (error messages on slink command) just for now create these folder in the canopsis installation folder before issuing again the slink command that will work properly then.

   .. code-block:: bash

      sudo mkdir /opt/canopis/folder-in-error1 /opt/canopis/folder-in-error2 ...
      sudo ./tool/slink

   This hack have to be fixed properly for the future install process.
