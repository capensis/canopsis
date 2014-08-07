================================================
Storage: Library for managing data storage
================================================

.. module:: canopsis.storage
    :synopsis: storage library for storing data specifically to generic properties.

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Storage provides (abstract) classes in order to manage data storage depending of generic properties and a Manager which is in charge of separating storage implementation from business code from the developer view point.

USed storage implementation of technologies are (re-)configured by thus Managers.

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Package contents
================

.. data:: __version__

    Current package version : 0.1

.. class:: DataBase(canopsis.configuration.Configurable)

    Abstract base class for all storages. In charge of connecting any storage to a data base.

.. data:: auto_connect

    If True, connect the database as soon as possible (at the end of the initialization, or at the end of any reconfiguration processing)

.. data:: host

.. data:: port

.. data:: db

.. data:: journaling

.. data:: safe

.. data:: wtimeout

.. data:: user

.. data:: pwd

.. data:: ssl

.. data:: ssl_key

.. data:: ssl_cert

.. function:: connect(self)

.. function:: disconnect(self)

.. function:: reconnect(self)

.. function:: connected(self)

.. function:: drop(self, table=None)

.. function:: size(self, table=None, criteria=None)

.. class:: Storage(DataBase)

.. data:: ASC

.. data:: DESC

.. data:: data_type

.. function:: bool_compare_and_swap(self, _id, oldvalue, newvalue)

.. function:: val_compare_and_swap(self, _id, oldvalue, newvalue)

.. function:: get_elements(self, ids=None, limit=0, skip=0, sort=None)

.. function:: find_elements(self, request, limit=0, skip=0, sort=None)

.. function:: remove_elements(self, ids)

.. function:: put_element(self, _id, element)

.. function:: get_table(self)

.. function:: copy(self, target)

.. function:: _element_id(self, element)

.. function:: _get_category(self)

.. function:: _get_storage_type(self)

.. module:: canopsis.storage.periodic

.. class:: PeriodicStorage(Storage)

.. module:: canopsis.storage.timed

.. class:: TimedStorage(Storage)

.. module:: canopsis.storage.typed

.. class:: TypedStorage(Storage)

.. module:: canopsis.storage.timedtyped

.. class:: TimedTypedStorage(Storage)

.. module:: canopsis.storage.manager

.. class:: Manager(canopsis.configuration.Configurable)
