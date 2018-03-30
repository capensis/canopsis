.. _FR__CalendarWidget:

===========================
Data queries to the backend
===========================

This document describes the way data queries should be handled in the frontend.

.. contents::
   :depth: 3


References
==========

Updates
=======

.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Gwenael Pluchon", "2016/02/10", "0.1", "Document creation", ""

Contents
========

.. _FR__Title__Desc:

Description
-----------

Adapter creation
^^^^^^^^^^^^^^^^

Adapters must extend DS.Adapter or a subclass of DS.Adapter, and be registered into the main Ember Application Container.

Here is an example of an adapter:

.. code-block:: javascript

   /**
    * @adapter event
    */
   var adapter = ApplicationAdapter.extend({

     buildURL: function(type, id) {
       void(id);
       return '/event';
     },

     findQuery: function(store, type, query) {

       var url = '/rest/events';

       if (query.skip !== undefined){
         query.start = query.skip;
         delete query.skip;
       }

       return this.ajax(url, 'GET', { data: query });
     }
   });

Fetching data
^^^^^^^^^^^^^

Fetching data requires three main steps :

- Retreive the correct adapter
- Query for data
- Handle the results

Retreive the correct adapter
''''''''''''''''''''''''''''

This step tends to be supported by Ember Data if you manage your data via stores. If it's the case, this step is completely transparent for the developer.

However, if you have particular needs and do not want to use Ember Data stores, it is possible to borrow the main application store to find the correct adapter.


.. code-block:: javascript

   var mainStore = container.lookup('store:main');
   var adapter = mainStore.adapterFor('datatype');


This is usually not something you want to do manually.

Query for data
''''''''''''''

If you have a dedicated store for your data, it is possible to ask for data directly through it.

.. code-block:: javascript

   store.findQuery(datatype, options).then(function(result) {
      //handle results
   })

If you have found your adapter manually, you can directly query it for data.

.. code-block:: javascript

   adapter.findQuery(undefined, datatype, options).then(function(results) {
      //handle results
   });


Handle the results
''''''''''''''''''

Adapters methods usually returns promises, that can be used to handle data correctly when they are found.
