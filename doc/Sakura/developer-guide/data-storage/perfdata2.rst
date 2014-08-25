Perfdata2
=========

Collections
-----------

-  ``perfdata2``: Meta data + current series
-  ``perfdata2_bin``: Grid FS for compressed series

Structure
---------

.. code-block:: javascript

	

    {
      "_id": Internal id
      "co":  Component
      "re":  Resource
      "me":  Metric name
      "u":   Unit
      "t":   Type
      "tg":  Array of tags
      "mi":  Min
      "ma":  Max
      "lv":  Last value
      "fts": First timestamp
      "lts": Last timestamp
      "c":   Array of _id of compressed series
      "d":   Array of current values
      "r":   Retention time
      "tw":  Warning threshold
      "tc":  Critical threshold
    }
