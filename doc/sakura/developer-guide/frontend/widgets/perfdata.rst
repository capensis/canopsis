.. _dev-frontend-widgets-perfdata:

Using Performance Data
======================

This document describe how to use series and metrics in your widget.

Updating the model
------------------

You will need to choose the series/metrics to use when creating your widget.
For that purpose, add this to your widget's model:

.. code-block:: javascript

   {
      "type": "object",
      "categories": [{
         "title": "Series",
         "keys": ["series"]
      },{
         "title": "Metrics",
         "keys": ["metrics"]
      }],
      "properties": {
         "series": {
            "title": "Select computed metrics",
            "type": "array",
            "role": "arrayclassifiedcrecordselector",
            "model": "serie",
            "crecordtype": "serie",
            "valueKey": "id",
            "multiselect": true
         },
         "metrics": {
            "title": "Data source for series",
            "type": "array",
            "role": "cmetric",
            "items": {
               "type": "string"
            }
         }
      }
   }

Updating the controller
-----------------------

You will need to include the controllers ``series`` and ``perfdata`` in order to
fetch data.

Let start with the following skeleton:

.. code-block:: javascript

   define([
      'jquery',
      'ember',
      'app/lib/factories/widget',
      'app/controllers/serie',
      'app/controllers/perfdata'
   ], function($, Ember, WidgetFactory) {
      var get = Ember.get,
          set = Ember.set,
          isNone = Ember.isNone;
   
      var widgetoptions = {};
   
      var widget = WidgetFactory('simplewidget', {
         needs:  ['serie', 'perfdata']
   
         init: function() {
            this._super.apply(this, arguments);
         },
   
         findItems: function() {
            var now = new Date().getTime();
            var to = now, from = get(this, 'lastRefresh');

            if (isNone(from)) {
               // for example:
               from = to - get(this, 'time_window');
            }

            this.fetchSeries(from, to);
            this.fetchMetrics(from, to);
         },

         fetchSeries: function(from, to) {
            // find series
         },

         fetchMetrics: function(from, to) {
            // find metrics
         }
      }, widgetoptions);
   
      return widget;
   });

Retrieving series
-----------------

The serie selector will fill the ``series`` array with all selected serie ids.
You can use this array to retrieve, from the Ember store, the serie records.
They will be passed to the ``serie`` controller.

.. code-block:: javascript

   fetchSeries: function(from, to) {
      var series = get(this, 'series'),
          ctrl = get(this, 'controllers.serie'),
          store = get(this, 'widgetDataStore');

      // fetch series
      store.findQuery('serie', {ids: series}).then(function(result) {
         var promises = [], records = get(result, 'content');

         for(var i = 0, l = get(result, 'meta.total'); i < l; i++) {
            var serie = records[i];

            // fetch points from serie
            var promise = ctrl.fetch(serie, from, to);
            promises.push(promise);
         }

         // resolve when all promises resolved
         Ember.RSVP.all(promises).then(function(results) {
            // the serie controller returns an array of points
            // so results is an array of arrays of points

            for(var i = 0, l = results.length; i < l; i++) {
               var points = results[i];

               // do something with points
            }
         });
      });
   }

Retrieving metrics
------------------

Just like the serie selector, the metric selector will fill the ``metrics`` array
with metric ids.
The difference is that this time, you don't need to call the store, just pass the
array to the ``perfdata`` controller:

.. code-block:: javascript

   fetchMetrics: function(from, to) {
      var metrics = get(this, 'metrics'),
          ctrl = get(this, 'controllers.perfdata');

      ctrl.fetchMany(metrics, from, to).then(function(result) {
         // result = {total: X, data: [...], success: True/False}

         for(var i = 0, l = result.total; i < l; i++) {
            var metric = result.data[i];
            var meta = metric.meta,
                points = metric.points;

            // do something with points
         }
      });
   }

Aggregating metrics
-------------------

To aggregate metrics, you need a method of aggregation, and an interval.
The aggregation method will be used to calculate a single value from all points
within the interval.

.. code-block:: javascript

   var ctrl = get(this, 'controllers.perfdata');

   // will return one point each 5 minutes
   // this point is the average of all points in the previous 5 minutes
   ctrl.aggregateMany(metrics, from, to, 'average', 300 * 1000);

**NB:** The controller ``perfdata`` also provides ``fetch`` and ``aggregate``.
They take a string, representing the metric id, instead of an array of strings.
