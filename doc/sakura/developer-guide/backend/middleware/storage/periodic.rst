Periodic Storage
================

.. module:: canopsis.storage.periodic

.. class:: PeriodicStorage(canopsis.storage.Storage)

   Storage dedicated to manage periodic data which are a set of values in an interval of timestamp. The minimal timestamp corresponds to a round time with a saved period.

   .. data:: __datatype__ = 'periodic'

      storage data type name

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
