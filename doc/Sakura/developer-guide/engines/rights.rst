Rights
_______


Data Structures
================

User 
-----

.. code-block:: javascript

    User = {
        'rights': ...,               // Map of type Rights, every user-specific rights goes here
        'groups': ...,               // List of group names, every user-specific groups goes here
        'role': ...,                 // List of role names that defines the User's profile, groups, and rights
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

Example:

.. code-block:: javascript

    User = {
        'role': 'manager',
        'contact': {
            'mail': 'jharris@scdp.com',
            'phone_number': '+33678695041',
            'adress': '1271 6th Avenue, Rockefeller Center, NYC, New York'
            }
        'name': 'Joan Harris',
        '_id': '1407160264.joan.harris.manager'
        }


Role
-------

A Role is specific to a small number of users

.. code-block:: javascript

    'name': {
        'rights': ...               // Map of type Rights, every role-specific rights goes here
        'profile': ...              // ID of the profile (string)
        FIELD: ...                  // You can add any number of fields that can be used with data-specific rules
        ...
        }
        
        
Example:

.. code-block:: javascript

    Roles = {
        'manager': {
            'profile': 'Manager',
            'rights': {},
            'list_of_directors': ['Ted Chaough', 'Peggy Olson', 'Don Draper']
            }
        }

    
Profile
---------

A profile is generic and global to all users

.. code-block:: javascript
 
    'name': {                            // String of profile's name
        'composites': ...                // List of the groups the profile belongs to
        }



Example:

.. code-block:: javascript

    An Administrator profile exists, it has all rights and belongs to the Group Management as well as the root Group
    Profiles = {
        'Manager': {
            'composites': ['managements', 'supervizion']
        }
        
    

Composite
-------

A composite is generic and global to all users

.. code-block:: javascript

    'name': {                        // String of group's name
        'members': ...,              // List of members ids
        'rights': ...                // Map of type Rights
        }
        
        
Example:

.. code-block:: javascript

    Groups = {
        'management': {
            'members': ['1407160264.joan.harris.manager'],
            'rights': {
                userconf_view_id: {
                    'checksum': 1,
                    'desc': ['Access user configuration']
                    },
                role_specific_id: {
                    'checksum': 15,
                    'field': 'list_of_directors',
                    'desc': ['Access and change directors' configuration']
                }
            }
        }
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

    if Rights[object_idXYZ]['right'] & (READ | CREATE | UPDATE | DELETE) == (READ | CREATE | UPDATE | DELETE):
        #the user has all rights on the object identified with object_idXYZ
        
    if not Rights[object_idXYZ]['right'] & (CREATE | DELETE):
        #the user has none of the rights on the object identified with object_idXYZ

User-specific and role-specific rights
.......................................

By default, the users have their groups rights, if a user needs or wants specific rights, they are added to its own ``Rights`` field.

Example::

    Group_1 = Alice, Bob
    Group_2 = Alice, Mark, Tom
    Group_3 = Jerry, Tom

    Alice creates a widget and sets the visibility to her groups; We add the right to the Group_1's and Group_2's rights

    Alice, Bob, Mark, and Tom will be able to access the widget. 

    Alice creates a Widget and sets the visibility to only her; We add the right to Alice's rights

    Only Alice can access the Widget, 
