How to manage data load on canopsis build
=========================================

When canopsis is build, some documents are inserted into database. These files are the canopsis basis that make canopsis work with default parameters.
It is possible to customize how Canopsis documents are loaded, and finally, how Canopsis behaves.

1. filldb
---------

The filldb command makes canopsis to run some python scripts located in ``/opt/canopsis/opt/mongodb/load.d`` in sequential order. These files contains instructions in python format on how to load default configuration into canopsis.

This script can be run with two parameters (logged as canopsis user):

   python ~/opt/mongodb/filldb.py **init** | **update**

Init parameter will erase all previous information from database, whereas update will only update the database.


2. json loader
--------------

The file module ``11_jsonloader.py`` enables to tell canopsis to load custom files into the database depending in upsert mode following below considerations:

The json loader script searches each folder located in the ``/opt/canopsis/opt/mongodb/load.d`` folder that are prefixed by **json_**.  Theses folders name have to be followed by a collection name that will tell where to upsert json documents. If a json document have to be upsert in the object collection, the json folder have to be called **json_object**.


Json folders to contain **<filename>.json** files that contains either a document or a list of document. These document must be identified by a special key that allow upsert. This key is named ``loader_id`` and must be uniq over the collection folder documents. Below a sample of a json document being inserted in the object collection because the file path is ``/opt/canopsis/opt/mongodb/load.d/json_object/mydocument.json``


in a single document:

.. code-block:: javascript

   {
      "loader_id":"000",
      "document_key_1": "document_value_1"
   }

or in a document list:

.. code-block:: javascript

   [
      {
         "loader_id":"000",
         "document_key_2": "document_value_2"
      },
      {
         "loader_id":"000",
         "document_key_3": "document_value_3"
      }
   ]


it is possible to prevent a document for beiing updated by the json loader by adding the ``loader_no_update`` key equals to **true** in the json document.
