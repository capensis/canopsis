.. _dev-frontend-cmp-internal:

Internal components for Canopsis plugins
========================================

Internal components exists within the ``components`` folder of the plugin.

Create your first component
---------------------------

A Canopsis component is basically an Ember component :

.. code-block:: bash

   $ cd var/www/canopsis/core/components
   $ mkdir myawesomecomponent
   $ cd myawesomecomponent
   $ touch component.js
   $ touch template.html

Now, you can edit your ``component.js`` :

.. code-block:: javascript

   define(['ember'], function(Ember) {
   
       var get = Ember.get,
           set = Ember.set;
   
       var component = Ember.Component.extend({
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

And your ``template.html`` :

.. code-block:: html

   <h1>{{myData}}</h1>

Now you have your skeleton, ready to implement some magic.

You will be able to add mixins, actions, or even use other components within
your template.
