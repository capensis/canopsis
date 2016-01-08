.. _admin-manage-engines-datacleaner:

Data cleaner
============

The data cleaner engine is an asynchronous engine. It is ran by default in canopsis configuration. It aims to perform database purge depending on some parameters that are explained below.

How does it works
-----------------

* Database cleaned collection are ``events`` and ``events_log``. These values are hard coded.
* By default The engine is ran each day thanks to it's beat_interval value set to `3600 * 24` seconds.
* This engine can be parametrized from the UI via the **settings menu > data cleaner**


Engine parameters
-----------------

* The engine is parametrized thanks to a database document that contains two fields ``retention_duration`` and ``use_secure_delay``.
* ``retention_duration`` is the duration in seconds from now to perform data deletion. For instance, if this value is set to 3600, all events and events log will be cleaned when these document's timestamp reach a timestamp lower or equal to ``'now timestamp' - 3600``. This parameter is modulated by the ``use_secure_delay`` parameter explained below
* ``use_secure_delay`` is a boolean value that will affect the retention delay computed from `now`. When ``True`` and if the retention duration is lower than one year, then the retention duration will automatically be set to ``now - one_year`` (one year is computed as 3600 * 24 * 365). Otherwise when this parameter is set to false, secutity is removed and the user can input a retention duration that will be directly applyed and so it is possible to remove data from now minus one hour.


The engine parameter is stored in the database with following informations and is loaded by the json loader at canopsis initialisation. more details about json loader at `json loader <../../../administrator-guide/setup/filldb.html>`_


.. code-block:: javascript

   {
       "_id": "datacleaner",
       "loader_id": "datacleaner",
       "crecord_type": "datacleaner",
       "retention_duration" : 31536000,
       "use_secure_delay" : true,
       "loader_no_update": true
   }
