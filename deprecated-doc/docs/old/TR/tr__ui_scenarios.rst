.. _TR__Title:

============
UI Scenarios
============

UI Scenarios provides a way to automatically trigger a sequence of actions in the frontend.

Applications:

 - Schedule a sequence of actions from a widget (a button, a weather tile, ...)
 - Provide functionnal testing ability to Canopsis core developers, product owner, integrators and continuous integration bot

.. contents::
   :depth: 2

References
==========

Updates
=======


.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Gwenael Pluchon", "2015/11/30", "1.0", "First draft", ""

Contents
========

Scenario schema
---------------

.. code-block:: javascript

   {
      "type": "array",
      "items": {
         "type": "object",
         "properties": {
            "xtype": {
               "type": "string"
               "enum":["action"],
               "description" : "type of document. Fixed to \"action\""
            },
            "actionType": {
               "type": "string",
               "description" : "type of action"
            },
            "parameters": {
               "type": "object",
               "description" : "Parameters passed as-is to to action implementation"
            },
         }
      }
   }


Using scenarios in widgets
--------------------------

Scenarios can be used in widgets. To put a scenario in widgets properties, the role "scenario" have to be put on the corresponding property.

Implementing actions
--------------------

Actions are Ember Objects that derivates from the Action Class. Actions classes follows the "command" design pattern, and have the following methods :

 - do(source, parameters)
 - undo(source, parameters)

These methods should return promises to be able to have a control over actions sequences (to be able to know when actions have been executed).

With this kind of implementation, actions can be stacked easily on a data structure.

Action playback
---------------

At first, actions are executed sequentially one after another.

Example
-------

Here is the example of actions that could be implemented :

 - changeView: transitions to a defined view
 - applyFilterToLists: apply a defined filter to all list widgets in the view
 - changeWidgetsProperty: change a defined property to defined widget(s)
