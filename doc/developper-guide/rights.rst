Rights
_______


Data Structures
================

User 
-----

.. code-block:: javascript

    User = {
        'rights': ...,               // Map of type Rights
        'groups': ...,               // Map of groups the User belongs to
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
        'members': ...,              // Map of members ids (members[id]==true if the id is in the group)
        'rights': ...                // Map of type Rights
        }
    
Profile
---------

.. code-bloc:: javascript

    Profile = {
        'name': ...,                 // String of profile's name
        'rights': ...,               // Map of type Rights
        'groups': ...                // Map of groups the profile belongs to
        }


Profiles are used as a template to create new users, it sets the User's fields with the Profile's one and allows a better management of the users

Example::

    An Administrator profile exists, it has all rights and belongs to the Group Management
    Profile = {
        'name': 'Administrator',
        'rights': ...,
        'groups': {'Management': true}
        }
        
    If you now create a User and specifiy the profile Administrator during the creation,
    you will get :
    
    User = {
        'rights': ...,                   // Rights specified in the profile, here, everything
        'groups': {'Management': true},
        'profile': 'Administrator',  
        '_id': ...                       // uniq id
        }  
        
    You can now set user-specific values (such as the name, contact information, etc..)

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

    if Rights[object_idXYZ]['right'] & (READ | CREATE | UPDATE | DELETE) == Rights[object_idXYZ]['right']:
        #the user has all rights on the object identified with object_idXYZ
        
    if not Rights[object_idXYZ]['right'] & (CREATE | DELETE):
        #the user has none of the rights on the object identified with object_idXYZ

User-specific rights
......................

By default, the users have their groups rights, if a user needs or wants specific rights, they are added to its own ``Rights`` field.

Example::

    Group_1 = Alice, Bob
    Group_2 = Alice, Mark, Tom
    Group_3 = Jerry, Tom

    Alice creates a widget and sets the visibility to her groups; We add the right to the Group_1's and Group_2's rights

    Alice, Bob, Mark, and Tom will be able to access the widget. 

    Alice creates a Widget and sets the visibility to only her; We add the right to Alice's rights

    Only Alice can access the Widget, 
