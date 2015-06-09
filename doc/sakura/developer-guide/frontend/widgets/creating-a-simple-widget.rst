.. _dev-frontend-widgets-howto:

Creating a simple widget
========================

Introduction
------------

This guide will help you to master the key concepts of creating widgets for Canopsis.


Requirements
------------

This chapter requires some basic Javascript and Object Oriented Programming (OOP) knowledge.


Creating a simple widget
------------------------

First of all, a widget needs at least 3 elements :

 - A **schema**, to provide a data model to the widget and to make it editable and storable
 - A **controller**, that is a javascript file that contains an canopsis factory witch role is to instanciate a widget on demand. This factory is parametrized with the widget controller.
 - A **template**, that contains the layout of the widget (it is basically html with a template language)

Creating boilerplate files
--------------------------

Let's start with a simple case, and create a widget that displays only a label which can be configured in the widget edit form.

First, define a **schema** that contains one editable attribute, the text to display :

``etc/schema.d/widget.label.json``

.. code-block:: javascript

   {
      "type": "object",
      "categories": [{
         "title": "General",
         "keys":["texttodisplay"]
      }],
      "properties": {
         "texttodisplay": {
            "type": "string"
         }
      }
   }


Then, define a **controller** that is required but will not contain any particular method or property:

``core/widgets/label/controller.js``

.. code-block:: javascript

   define([
      'app/lib/factories/widget',
      'app/lib/loaders/schemas'
   ], function(WidgetFactory) {

      //below the instanciation of the controller that only contains an empty object
      var widget = WidgetFactory('label', {});

      return widget;
   });


And finally, define a **template** that show the field setted by the view manager:

``core/widgets/label/template.html``

.. code-block:: html

   <h1>
      {{texttodisplay}}
   </h1>


Test the widget
---------------

Now you have a working widget, but it is not registered nor loaded by the web application.

To register the schema, make sure the file is in the schema directory, and execute the ``schema2db``. The new schema will be sent to the database.

To register the widget JS and HTML, open the widget loader file, (usually in ``<js plugin file>/lib/loaders/widgets``), and add ensure the widget is referenced in the widget list.


You should now be able to display a view, enter the edit mode, place a widget, display it and reconfigure it !


Deeper in the widget creation
-----------------------------

As the widget factory makes your widget inherit from Canopsis Widget class, some behaviors are available from this super class. This inheritance layer brings for exemple the layout placement management in the GUI, the parametrable refreshing system for your widget depending on the **refreshableWidget** boolean property in the schema and the **refreshInterval** value in seconds.


Imagine you want to make your widget reresh a new label each 30 seconds, what you have to do is to add the following property to your widget schema

.. code-block:: javascript

   {
      "properties": {
         "refreshableWidget":  {
            "type": "boolean",
            "default": true
         }
      }
   }

Then, by inheritance, the **refreshInterval** property will give your widget the resfreshInterval property set to 60 seconds by default. We will change it in a first time in a hardcoded way for demonstration purposes.

Now let update our controller definition with the hardcoded param from:

.. code-block:: javascript

   //header code ...

   var widget = WidgetFactory('label', {});

   //end widget code ...


.. code-block:: javascript

   //header code ...

   //good practice in canopsis is to define and use shortcuts to Ember.get and Ember.set
   var get = Ember.get,
       set = Ember.set;

   var widget = WidgetFactory('label', {

      init: function () {
         //The hardcoded value set
         set(this, 'refreshInterval', 30);

         //Calling the super call is required when overriding the constructor
         this._super();

      }

   });

   //end widget code ...


Using Ember js set and get methods will trigger databinding recomputation and this way, all the widget remains up to date with the lastest information.
Now we have updated our widget with a custom value the widget should refresh sooner than by default.

Beyond the simple widget
------------------------

A widget will now become whatever you want as the given widget basis upper let you create what you need. Thus, in widgets it is possible to use components (see `canopsis components <#components>`_) in a way as simple as the following code

.. code-block:: html

   {{component-mycustomcomponent content=dataBindingVariable}}

where the dataBindingVariable will be updated by the sub component and be reachable in your widget with the following code in the controller:

.. code-block:: javascript

   // controller header...

   methodUsingComponentValue: function () {
      var componentValue = get(this, 'dataBindingVariable');
      //process value from component
   }


- Manipulating remote data can be done preferably with records and adapters or can be acheived with jquery ajax queries.
- It is possible to display custom template information mixed with components, loop controls, loginc controls, helpers and any other facilities canopsis environement provides within the widget template (see `architecture </developer-guide/frontend/architecture.html>`_).
- Don't forget what Ember framework brings to you, it is possible to run code on dom element ready in a widget because **didInsertElement** is called in the widget view when dom element is rendered (don't forget to use this._super()).
- The current dom element can be reached with **this.$()** in the widget view.
- Uderstanding how Ember js works will help you writting your own widget `Ember js <http://emberjs.com/api/>`_.
