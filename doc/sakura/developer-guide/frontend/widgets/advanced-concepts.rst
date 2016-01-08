.. _dev-frontend-widgets-advanced:

Advanced concepts
=================

Widget factory options
----------------------

.. NOTE::

   TODO: document this part

View mixins
-----------

The widget factory uses the ``viewMixins`` array to apply listed mixins to the
widget's view.

For example :

.. code-block:: javascript

   define([
      'jquery',
      'ember',
      'app/lib/factories/widget'
   ], function($, Ember, WidgetFactory) {
      var get = Ember.get,
          set = Ember.set;
   
      var widgetoptions = {};
   
      var SimpleViewMixin = Ember.Mixin.create({
         didInsertElement: function() {
            var ctrl = get(this, 'controller');
            var domElement = this.$();

            console.log('Widget inserted in DOM:', ctrl, domElement);

            this._super.apply(this, arguments);
         }
      });
   
      var widget = WidgetFactory('simplewidget', {
         viewMixins: [SimpleViewMixin],
   
         init: function() {
            this._super.apply(this, arguments);
         },
   
         findItems: function() {
            // callback for data fetching at refresh
         }
      }, widgetoptions);
   
      return widget;
   });

Make your widget interactive
----------------------------

Let's take the example of the Timegraph widget.
You can zoom on the graph by selecting an range of time on it. If you want to
reset the default time window, you have to press the button made for it.

This button sends an action to the view, it's this simple template :

.. code-block:: html

   <button type="button" class="btn btn-info btn-sm pull-right" {{action resetZoom target=view}}>
      <i class="fa fa-undo"></i>
   </button>

And now implement the handler in the view mixin :

.. code-block:: javascript

   var SimpleViewMixin = Ember.Mixin.create({
      didInsertElement: function() {
         var ctrl = get(this, 'controller');
         var domElement = this.$();
   
         console.log('Widget inserted in DOM:', ctrl, domElement);
   
         this._super.apply(this, arguments);
      },
   
      actions: {
         resetZoom: function() {
            // do something
         }
      }
   });

Now, if you want to send the action to the controller, just remove the ``target=view``
part of the template, and move the action handler to the widget controller :

.. code-block:: javascript

   var widget = WidgetFactory('simplewidget', {
      viewMixins: [SimpleViewMixin],
   
      init: function() {
         this._super.apply(this, arguments);
      },
   
      findItems: function() {
         // callback for data fetching at refresh
      },
   
      actions: {
         resetZoom: function() {
            // do something
         }
      }
   }, widgetoptions);

Now, if your template is made of Ember components, it makes them able to interact
with the widget's view and/or controller.
