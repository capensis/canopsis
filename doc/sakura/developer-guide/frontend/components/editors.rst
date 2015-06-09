.. _dev-frontend-cmp-editors:

Editors
=======

An editor is a simple template, eventually made of Ember components.
It is used to define how data, in a JSON schema, will be edited.

For example :

.. code-block:: javascript

   {
       "type": "object",
       "categories": [{
           "title": "General",
           "keys": ["bgcolor"]
       }],
       "properties": {
           "bgcolor": {
               "type": "string",
               "role": "color"
           }
       }
   }

Here, the key ``role`` tells the ``modelform`` to use the ``color.html`` template (our editor), in order to render the field.

Now, we will create our template ``editors/color.html`` within the UI plugin folder (core, uibase, monitoring, ...), containing something like :

.. code-block:: html

   {{input type="color" value=attr.value}}

Or, if you have an awesome *color selection* component :

.. code-block:: html

   {{component-colorselect content=attr.value}}

