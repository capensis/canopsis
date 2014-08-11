==========================================
Storage: Library for managing data storage
==========================================

.. module:: canopsis.storage
   :synopsis: storing data-oriented approach library.

.. moduleauthor:: jonathan labejof
.. sectionauthor:: jonathan labejof

Objective
=========

Storage provides a way to store data in a data-oriented approach (aka, technology and data-model agnostic).

Also, developers can choose to store data depending of kind of data to store, instaed of choosing the way to store data and to think about data-model.

Indices and tables
==================

* :ref:`genindex`
* :ref:`search`

Functional description
======================

Storage provides the Storage and Manager classes in order to manage data storage in a data_oriented approach.

Storage
-------

A storage works like a data base dedicated to specific data type. Each storage dedicated to data types is specialized and optimized for such data types. Five types of storages exists: Storage, PeriodicStorage, TimedStorage, TypedStorage and TimedTypedStorage.

Manager
-------

A manager is the interface between business code and storages. It specialized and optimized of set of storages for business components which needs to use several types of storages.

Technical description
=====================

.. data:: __version__

   Current package version : 0.1

.. class:: DataBase(canopsis.middleware.Middleware)

   Abstract base class for all storages. In charge of connecting any storage to a data base, and provide methods such as size(data) and drop(db).

   .. data:: db

      data base name

   .. data:: journaling

      journaling mode enabling.

   .. method:: drop(table=None)

      Drop related all tables or one table if given.

      :param table: table to drop
      :type table: str

      :return: True if dropped
      :rtype: bool

   .. method:: size(table=None, criteria=None)

.. class:: Storage(DataBase)

   .. data:: __storage_type__ = 'storage'

      storage type name

   .. data:: DATA_ID = 'id'

      data id field name

   .. data:: ASC = 1

      ascending order in search operations

   .. data:: DESC = -1

      descending order in search operations

   .. method:: bool_compare_and_swap(_id, oldvalue, newvalue)

      Performs an atomic compare_and_swap operation on database related to \
        input _id.

      :remarks: this method is not atomic

      :returns: True if the swamp succeed

   .. method:: val_compare_and_swap(_id, oldvalue, newvalue)

      Performs an atomic val_compare_and_swap operation on database related \
        to input _id, oldvalue and newvalue.

      :remarks: this method is not atomic

      :returns: True if the comparison succeed

   .. method:: get_elements(ids=None, limit=0, skip=0, sort=None)

      Get a list of elements where id are input ids

      :param ids: element ids or an element id to get if not None
      :type ids: list of str

      :param limit: max number of elements to get
      :type limit: int

      :param skip: first element index among searched list
      :type skip: int

      :param sort: contains a list of couples of field (name, ASC/DESC)
         or field name which denots an implicitelly ASC order
      :type sort: list of {(str, {ASC, DESC}}), or str}

      :return: input id elements, or one element if ids is an element
         (None if this element does not exist)
      :rtype: iterable of dict or dict or NoneType

   .. method:: find_elements(request, limit=0, skip=0, sort=None):

      Find elements corresponding to input request and in taking care of
      limit, skip and sort find parameters.

      :param request: set of couple of (field name, field value)
      :type request: dict(str, object)

      :param limit: max number of elements to get
      :type limit: int

      :param skip: first element index among searched list
      :type skip: int

      :param sort: contains a list of couples of field (name, ASC/DESC)
         or field name which denots an implicitelly ASC order
      :type sort: list of {(str, {ASC, DESC}}), or str}

      :return: input request elements
      :rtype: list of objects

   .. method:: remove_elements(ids)

      Remove elements identified by the unique input ids

      :param ids: ids of elements to delete
      :type ids: list of str

   .. method:: put_element(_id, element)

      Put an element identified by input id

      :param id: element id to update
      :type id: str

      :param element: element to put (couples of field (name,value))
      :type element: dict

      :return: True if updated
      :rtype: bool

   .. method:: count_elements(request)

      Count elements corresponding to the input request

      :param id: request which contain set of couples (key, value)
      :type id: dict

      :return: Number of elements corresponding to the input request
      :rtype: int

   .. method _find(*args, **kwargs)

      Find operation dedicated to technology implementation.

   .. method _update(*args, **kwargs)

      Update operation dedicated to technology implementation.

   .. method _remove(*args, **kwargs)

      Remove operation dedicated to technology implementation.

   .. method _insert(*args, **kwargs)

      Insert operation dedicated to technology implementation.

   .. method _count(*args, **kwargs)

      Count operation dedicated to technology implementation.

   .. method:: get_table()

      Table name related to elf type and data_type.

      :return: table name
      :rtype: str

   .. method:: copy(target)

      Copy self content into target storage.
      target type must implement the same class in cstorage packege as self.
      If self implements directly cstorage.Storage, we don't care about
      target type

      :param target: target storage where copy content
      :type target: same as self or any storage if type(self) is Storage

   .. method:: _copy(target)

      Called by Storage.copy(self, target) in order to ensure than target type is the same as self

   .. method:: _element_id(element)

      Get element id related to self behavior

   .. method:: _get_category()

      Get configuration category for self storage

   .. method:: _get_storage_type()

      Get storage type (last_value and timed are two storage types)

      :return: storage type name
      :rtype: str

   .. staticmethod:: _update_sort(sort)

      Add ASC values by default if not specified in input sort.

      :param sort: sort configuration
      :type sort: list of {tuple(str, int), str}

.. module:: canopsis.storage.periodic

.. class:: PeriodicStorage(Storage)

   Storage dedicated to manage periodic data which are a set of values in an interval of timestamp. The minimal timestamp corresponds to a round time with a saved period.

   .. data:: __storage_type__ = 'periodic'

      storage type name

   .. data:: TIMESTAMP = 'timestamp'

      timestamp field name

   .. data:: VALUES = 'values'

      values field name

   .. data:: PERIOD = 'period'

      period field name

   .. data:: LAST_UPDATE = 'last_update'

      last update timestamp field name

   .. method:: count(data_id, period, timewindow=None)

        Get number of periodic documents for input data_id.

   .. method:: size(data_id=None, period=None, timewindow=None)

      Get size occupied by research filter data_id

   .. method:: get(data_id, period, timewindow=None, limit=0, skip=0, sort=None)

      Get a list of points.

   .. method:: put(data_id, period, points)

      Put periodic points in periodic collection with specific period values.

      points is an iterable of (timestamp, value)

   .. method:: remove(data_id, period=None, timewindow=None)

      Remove periodic data related to data_id, timewindow and period.
      If timewindow is None, remove all periodic_data with input period.
      If period is None

.. module:: canopsis.storage.timed

.. class:: TimedStorage(Storage)

   Store dedicated to manage timed data. It saves one value at one timestamp. Two consecutives timestamp values can not be same values.

   .. data:: __storage_type__ = 'timed'

      storage type name

   .. class:: Index:

      result field index

      .. data:: TIMESTAMP = 0
      .. data:: VALUE = 1
      .. data:: DATA_ID = 2

   .. data:: VALUE = 'value'

      value field name

   .. data:: TIMESTAMP = 'timestamp'

      timestamp field name

   .. method:: get(data_ids, timewindow=None, limit=0, skip=0, sort=None)

      Get a dictionary of sorted list of triplet of dictionaries such as :

      dict(
         tuple(
            timestamp,
            dict(data_type, data_value), document id))

      If timewindow is None, result is all timed document.

      :return:
      :rtype: dict of tuple(float, dict, str)

   .. method:: count(data_id)

      Get number of timed documents for input data_id.

   .. method:: put(data_id, value, timestamp)

      Put a dictionary of value by name in collection.

   .. method:: remove(data_ids, timewindow=None)

      Remove timed_data existing on input timewindow.

.. module:: canopsis.storage.typed

.. class:: TypedStorage(Storage)

   Storage dedicated to manage typed data identified by the couple (type, id).

   .. data:: __storage_type__ = 'typed'

      storage type name

   .. method:: get(_ids=None, data_type=None, limit=0, skip=0, sort=None)

      Get a list of data identified among data_ids or a type

      :param data_ids: data ids to get
      :type data_id: list of str

      :param data_type: data_id type to get if not None
      :type data_type: str

      :param limit: max number of data to get
      :type limit: int

      :param skip: starting index of research if multi data to get
      :type skip: int

      :param sort: couples of field (name, value) to sort with ASC/DESC
         Storage fields
      :type sort: dict

      :return: a list of couples of field (name, value) or None respectivelly
         if such data exist or not
      :rtype: list of dict of field (name, value)

   .. method:: put(_id, data, data_type=None)

      Put a data related to an id

      :param _id: data id
      :type _id: str

      :param data_type: data type to update
      :type data_type: str

      :param data: data to update
      :type data: dict

   .. method:: remove(_ids=None, data_type=None):

      Remove data from ids or type

      :param _ids: list of data id
      :type _ids: list

      :param data_type: data type to remove if not None
      :type data_type: str

.. module:: canopsis.storage.timedtyped

.. class:: TimedTypedStorage(Storage)

   Store dedicated to manage timed typed data identified by the triplet (type, timestamp, id).

   .. data:: __storage_type__ = 'timedtyped'

      storage type name

   .. class:: Index

      result values field indexes

        TIMESTAMP = 0
        VALUE = 1
        TYPE = 2
        DOCUMENT_ID = 3

   .. data:: VALUE = 'value'

      value field name

   .. data:: TIMESTAMP = 'timestamp'

      timestamp field name

   .. data:: TYPE = 'type'

      type field name

   .. method:: get(data_ids, data_type=None, timewindow=None, limit=0, skip=0, sort=None)

      Get a dictionary of sorted list of triplet of dictionaries such as :

      dict(
         tuple(
            timestamp,
            dict(data_type, data_value), document id))

      If timewindow is None, result is all timed document.

      :param data_ids: list of data_id to find
      :type data_ids: list of str

      :param data_type: data type to find if not None
      :type data_type: str

      :param timewindow: timewindow
      :type timewindow: ctimeserie.timewindow.TimeWindow

      :param limit: max number of data to get
      :type limit: int

      :param skip: starting index of research if multi data to get
      :type skip: int

      :param sort: couples of field (name, value) to sort with ASC/DESC
         Storage fields
      :type sort: dict

      :return:
      :rtype: dict of tuple(float, dict, str)

   .. method:: count(data_id=None, data_type=None, *args, **kwargs)

      Get number of timed documents for input data_id.

   .. method:: put(data_id, data_type, value, timestamp)

      Put a dictionary of value by name in collection.

   .. method:: remove(data_ids=None, data_type=None, timewindow=None)

      Remove timed_data existing on input timewindow.

.. module:: canopsis.storage.manager

.. class:: Manager(canopsis.configuration.Configurable)

   Manages storages by name.

   .. data:: CONF_RESOURCE = 'manager/manager.conf'

      Configuration resource

   .. data:: TIMED_STORAGE = 'timed_storage'
   .. data:: PERIODIC_STORAGE = 'periodic_storage'
   .. data:: STORAGE = 'storage'
   .. data:: TYPED_STORAGE = 'typed_storage'
   .. data:: TIMED_TYPED_STORAGE = 'timed_typed_storage'

      configuration name for storage types

   .. data:: AUTO_CONNECT = 'auto_connect'

      configuration parameter which auto-connect storages

   .. data:: SHARED = 'shared'

      share storage by name in the same processus

   .. data:: CATEGORY = 'MANAGER'

      Configuration category

   .. data:: STORAGE_SUFFIX = '_storage'

      Configuration storage name suffix

   .. method:: get_storage(data_type=None, storage_type=None, shared=None, auto_connect=None, *args, **kwargs)

      Load a storage related to input data type and storage type.

      If shared, the result instance is shared among same storage type and data type.

      :param data_type: storage data type
      :type data_type: str

      :param storage_type: storage type (among timed, last_value ,etc.)
      :type storage_type: Storage or str

      :param shared: if True, the result is a shared storage instance among
         managers. If None, use self.shared
      :type shared: bool

      :return: storage instance corresponding to input storage_type
      :rtype: Storage
