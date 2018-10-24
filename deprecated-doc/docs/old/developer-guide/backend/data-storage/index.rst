.. _dev-backend-storage:

Data Storage
============

Canopsis stores different kind of data. Each one of this kind of data has specific
requirements about how it must be stored.

The Storage interface is built to answer those specifications, and implement them
with a specific technology (MongoDB, SQLAlchemy, ...).

.. toctree::
   :maxdepth: 2
   :titlesonly:

   composite
   timed
   periodic
   file
   default
