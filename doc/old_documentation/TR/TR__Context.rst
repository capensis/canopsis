.. _TR__Context:

=======
Context
=======

The new context engine manage a graph in MongoDB to have informations on entities inside canopsis.

.. contents::
   :depth: 2

References
==========

List of referenced functional and technical requirements...


Updates
=======


.. csv-table::
   :header: "Author(s)", "Date", "Version", "Summary", "Accepted by"

   "Thomas Gosselin, Arthur Dewarumez", "2017/02/06", "1.0", "context engine specifiation", ""

Contents
========

.. _TR__Context__Engine:

Engine
------

description:
>>>>>>>>>>>>

The context engine is in charge to keep the context's graph up to date.
He extracts connector, resource and component from events,
He builds ids and check if connectors, resources and components are already in the graph and if links exists.
If needed, the engine will update the graph.


Software architecture and costing
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

Data Model:

.. figure:: ../_static/images/context/data_model.png

Entities type relation models:

.. figure:: ../_static/images/context/entities_type_relation.png

Document in collection users:

.. code-block:: javascript

    {
        '_id': //user's id
        'name': //user's name
        'org': // user's organisation
        'access_to': //list of organisations the user can access
    }

Document in collection Organisations:

.. code-block:: javascript

    {
        '_id': //organisation's id
        'name': //organisation's name
        'parents': //list of organisation's parents _id
        'children': //list of organisation's children _id
        'views': //list of organisation's views _id
        'users': //list of users _id attached at the organisation
        'entities': //list of entities _id
    }

Document in collection Entities:

.. code-block:: javascript

    {
        '_id': //entity's' id
        'name': //entity's name
        'type': // component, resource, connector or application
        'depends': //list of entities _id
        'impact': //list of entities _id
        'measurements': //list of measurements _id
        'infos'://information about the entity
    }

Document in collection measurements:

.. code-block:: javascript

    {
        '_id': //measurements _id
        'tags': //lists of measurements tags
    }

Technical guide
>>>>>>>>>>>>>>>

UTs + costing
>>>>>>>>>>>>>

Description of Unit tests plus costing in day/man.

TFs + costing
>>>>>>>>>>>>>

Description of fonctional tests plus costing in day/man.

TPs + chiffrage
>>>>>>>>>>>>>>>

Description of performance tests plus costing in day/man.
