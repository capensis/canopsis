 * **POST** ``/perfstore/values`` parameters

.. code-block:: javascript

     {
          'nodes': [
               {
                    'id':          // metric's node identifier
                    'from':        // Interval's beginning
                    'to':          // Interval's end
               },
               // ...
          ],

          'consolidation_method': 'mean' or 'min' or 'max' or 'sum' or 'delta',
          'aggregate_method': 'first' or 'last' or 'mean' or 'min' or 'max' or 'sum' or 'delta',

          'aggregate_interval':    // Aggregation interval in seconds
          'aggregate_max_points':  // Maximum of points to returns
          'aggregate_round_time':  // if True and the aggregation interval is 15m, and start of data are at 4:03, the interval will be [4:00, 4:15], otherwise it will be [4:03, 4:18]
          'timezone':              // Time Zone for intervals calculation
     }
