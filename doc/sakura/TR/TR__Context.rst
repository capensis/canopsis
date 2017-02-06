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

   "Thomas Gosselin", "2017/02/06", "1.0", "context engine specifiation", ""

Contents
========

.. _TR__Context__Engine:

Engine
------


Software architecture and costing
>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>

Data Model:

ICI ON MET LE SCHÉMA TOUT ÇA

Document in collection users:

.. code-block:: javascript

    {
        '_id': //user's id
        'name': //user's name
        'org': // user's organisation
        'access_to': //list of organisations 
    }

Document in collection Organisations:

.. code-block:: javascript

    {
        '_id': //organisation's id
        'name': //organisation's name
        'parents': //list of organisation's parents _id
        'children': //list of organisation's childre _id
        'views': //list of organisation's views
        'users': //list of users attached at the organisation
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
        '_id': //
        'tags': //
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
