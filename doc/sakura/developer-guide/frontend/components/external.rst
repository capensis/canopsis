.. _dev-frontend-cmp-external:

External components by Canopsis
===============================

Some components may be useful for other uses than in Canopsis.
Those components will be written in ``externals/webcore-libs``, for example :

.. code-block:: bash

   $ cd externals/webcore-libs
   $ mkdir -p ember-myawesomecomponent/lib
   $ touch ember-myawesomecomponent/lib/component.js

And here is the skeleton :

.. code-block:: javascript

   define(['ember'], function(Ember) {
       var get = Ember.get,
           set = Ember.set;
   
       var component = Ember.Component.extend({
           template: Ember.Handlebars.compile('<h1>{{myData}}</h1>'),
           myData: 'Hello World!'
       });
   
       Ember.Application.initializer({
           name: 'component-myawesomecomponent',
           initialize: function(container, application) {
               application.register('component:component-myawesomecomponent', component);
           }
       });
   
       return component;
   });

As you can see, there is no ``template.html`` for external components, it must be
compiled inside the component.
