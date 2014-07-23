Access Control List
===================

User rights in Canopsis are managed via 3 entities :

Users
-----

Groups
------

Profiles
--------

A profile is the entity that owns rights records.

Rights
^^^^^^

Generic Rights specification
''''''''''''''''''''''''''''

one namespace = one record in db

.. code-block:: javascript

   {
      id: "rk",
      allowed_profiles: []
   }

the id field, which is a string separated by dots, allows to group rights within categories. The principle is to set descriptive ids into rights, such as ``canopsis.ui.widgets.timegraph``, going from the more general tag to the more specific one. The last part of the rk is the action name. Examples :

 - ``canopsis.ui.widgets.refresh``
 - ``canopsis.ui.widgets.timegraph.refresh``


Specific Rights specification
'''''''''''''''''''''''''''''

Specific rights takes place when users override rights on a per-widget basis. Rights are then stored within the widget record, with a syntax similar as :

.. code-block:: javascript

   {
      id: "widget1",
      title: "mywidget",
      // [...]
      rights: [
         refresh: ["profile_1", "profile_2"]
      ]
   }


Referential integrity
'''''''''''''''''''''

Referential integrity is maintained between profiles and views