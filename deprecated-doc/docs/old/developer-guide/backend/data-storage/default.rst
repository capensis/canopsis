.. _dev-backend-storage-default:

Default Storage
===============

A default storage is used to access store and retrieve documents from a database collection.
Searches, insertions, updates and deletions are made using the document's ID.

How to
------

It is possible to store and retrieve data from a database collection using the default storage.
Below is an example on how to manage data from a manager.

Here is the minimal code for the manager

.. code-block:: python

   # Allow retrieve configuration from configuration files
   from canopsis.configuration.configurable.decorator import (
       conf_paths, add_category)
   # Base class allowing to manage data storage in an easy way
   from canopsis.middleware.registry import MiddlewareRegistry

   # The file where the configuration is retrieved
   # here it is /opt/canopsis/etc/my_manager/my_manager.conf
   CONF_PATH = 'my_manager/my_manager.conf'

   #The category from the configuration file that is read and fetched at runtime
   CATEGORY = 'MANAGER'


   # Effective load of the configuration files and parameters
   @conf_paths(CONF_PATH)
   @add_category(CATEGORY)
   class Mymanager(MiddlewareRegistry):

      # The configuration key to read in the configuration file
      MANAGER_STORAGE = 'manager_storage'

      def __init__(self, *args, **kwargs):

        super(Mymanager, self).__init__(*args, **kwargs)



Here is a the Mymanager class manager's method to put data in a default storage:

.. code-block:: python

      def put(self):

         data_to_upsert = {
            'key0': 'value0',
            'key1': 'value1'
         }

         entity_id = 'entity_id'

         self[Connectors.MANAGER_STORAGE].put_element(
            _id=entity_id, element=entity
         )


This operation will lead to an upsert. If the document does not exists with the key entity_id, it is created and if this code is called twice with some different data_to updsert, the new keys are simply added to the document.


This will result in a database document (here with mongodb) looking like that:

.. code-block:: javascript

   {
      "_id" : "entity_id",
      "key1" : "value1",
      "key0" : "value0"
   }


Below the code to manage data deletion from the manager:

.. code-block:: python

   def delete(self):
      ids = ['entity_id']
      self[Connectors.MANAGER_STORAGE].remove_elements(ids=ids)


And in the end here is how to select data from the storage:

.. code-block:: python

   def get(self):
      ids = ['entity_id']
      state_documents = self[Connectors.MANAGER_STORAGE].get_elements(
         ids=ids
      )

The ``get_elements`` method also accepts a query parameter that is for instance a mongodb like filter taking a form of a dict.
