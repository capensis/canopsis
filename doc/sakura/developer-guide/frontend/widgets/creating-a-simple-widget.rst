Creating a simple widget
************************

Introduction
------------

This guide will help you to master the key concepts of creating widgets.


Requirements
------------

This chapter requires some basic Javascript and Object Oriented Programming (OOP) knowledge.

Creating a simple widget
------------------------

First of all, a widget needs at least 3 elements :

 - a **schema**, to provide a data model to the widget and to make it editable and storable
 - a **controller**, that contains the algorithm of the widget
 - a **template**, that contains the layout of the widget (it is basically html with a template language)

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

      var widget = WidgetFactory('label', {});

      return widget;
   });


And finally, define a **template** that show the field setted by the view manager:

``core/widgets/label/template.html``

.. code-block:: html

   <h1>
      {{texttodisplay}}
   </h1>


Now you have a working widget, but it is not registered nor loaded by the web application.

To register the schema, make sure the file is in the schema directory, and execute the ``schema2db``. The new schema will be sent to the database.

To register the widget JS and HTML, open the widget loader file, (usually in ``<js plugin file>/lib/loaders/widgets``), and add ensure the widget is referenced in the widget list.


You should now be able to display a view, enter the edit mode, place a widget, display it and reconfigure it !
