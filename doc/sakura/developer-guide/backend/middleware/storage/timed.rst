TimedStorage
============

.. module:: canopsis.storage.timed

.. class:: TimedStorage(canopsis.storage.Storage)

   Store dedicated to manage timed data. It saves one value at one timestamp. Two consecutives timestamp values can not be same values.

   .. data:: __datatype__ = 'timed'

      storage data type name

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
