Administration
==============

.. toctree::
   :maxdepth: 2

   services-management/index
   mongo-collections/index
   engines-management/index
   amqp2engines


Once canopsis installed, it is possible to get informations about current build. When the webserver is started, go to the link

``http://your.canopsis.dns:8082/en/static/canopsis/canopsis-meta.json``

.. code-block:: javascript

   {
      "build-date": "Thu Feb 19 08:28:48 UTC 2015",
      "build-timestamp": 1424334528,
      "git-commit": "202396bd8b2d200938cc353dccf590f3d6c2422f"
   }
