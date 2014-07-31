Rights
_______


Data Structures
================

User 
-----

.. code-block:: javascript

    User = {
        'rights': ...,               // Map of type Rights
        'groups': ...,               // List of strings (groups names)
        'profile': ...,              // String of profile name (Admin, Root, Manager, ...)
        'contact': {                 // Map of contact informations
            'mail': ...,
            'phone_number': ...,
            ...
            }
        'name': ...,                 // String of user's name
        '_id': ...                   // uniq id
    }

When an action is triggered, the ``object_id`` of the target of the action is sent and we check if one of the user's groups has the rights needed to perform the action.
If no groups among the user's has the right, we then check the user's own rights if he has any.

Group
-------

.. code-block:: javascript

    Group = {
        'name': ...,                 // String of group's name
        'members': ...,              // List of strings (members names)
        'rights': ...                // Map of type Rights
    }

Rights
----------

.. code-block:: javascript

    Rghts = {
        object_id...: {             // Right on the object with the identifier id
            'right': ...,           // 1 == Read, 2 == Update, 4 == Create, 8 == Delete
            'desc': ...,            // Short desc of the right
            'context': ...          // Time period
            }
    }

The keys of a map of type ``Rights`` are the ids of the objects accessible from the web application.
The ``right`` field is a 4-bit integer that goes from 1 to 15 and that describes the available action on the object.


.. code-block:: python

    if Rights[object_idXYZ]['right'] & (READ | CREATE | UPDATE | DELETE):
        #the user has all rights on the object identified with object_idXYZ
        
        
User-specific rights
=====================

By default, the users have their groups rights, if a user needs or wants specific rights, they are added to its own ``Rights`` field.

Example :

Group_1 = Alice, Bob
Group_2 = Alice, Mark, Tom
Group_3 = Jerry, Tom

Alice creates a widget and sets the visibility to her groups; We add the right to the Group_1's and Group_2's rights

Alice, Bob, Mark, and Tom will be able to access the widget. 


Alice creates a widhet and sets the visibility to only her; We add the right to Alice's rights

Only Alice can access the Widget, 
