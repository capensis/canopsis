Serializers and Adapters
************************

Mock data reception
-------------------

You can easily add a mock on the RestAdapter to help debugging.

.. code-block:: javascript

   //
   // Mock the HTTP requests
   //
   DS.RESTAdapter.reopen({
      find: function(store, type, id){
         return new Ember.RSVP.Promise(function(resolve, reject) {
            var json = {
               "child": {
                  "name": "Herbert",
                  "toys": [{
                     "kind": "Robot",
                     "size": {
                       "height": 5,
                        "width": 5,
                        "depth": 10
                     },
                     "features": [{
                        "name": "Beeps"
                     },{
                        "name": "Remote Control"
                     }]
                  }]
               }
            };
            json.child.id = id; // Set ID that we're looking for on the object
            Ember.run(null, resolve, json);
         });
      }
   });