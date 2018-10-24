.. _dev-backend-storage-file:

File Storage
============

This type of storage is used to store binary documents into the database.

Such file storage is quite different than a file system in an operating system:

- it manages document versions.
   + only new documents can be written.
   + only stored documents can be read.
- when accessing to documents, a dedicated I/O object is given in order to read or write document data.

How To
------

Let ``FS`` an instance of a file storage.

Write 'test' in a new document named ``my file``
################################################

.. code-block:: python

   # put 'test' directly in a new version of 'my file'
   filestream = FS.put(name='my file', data='test')

Get data from the second version of ``my file``
###############################################

.. code-block:: python

   # get a specific document version
   filestream = FS.get(name='my file', version=2)
   # get document data
   data = filestream.get()

You can list/delete documents with methods ``list`` and ``delete`` as well.
