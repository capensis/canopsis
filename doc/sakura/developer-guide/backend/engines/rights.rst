Rights
_______

Table of contents
-------------------

1. What_

   1. SetUp_

2. How_

   1. Rights

   2. Composites

   3. Profile

   4. Role

3. Library_

   1. Data Structures_

      1. User

      2. Role

      3. Profile

      4. Composite

      5. Rights


.. _what:

What
-----

Rights are defined within the composites (groups) upon their creations or added unitarily in a specific profile, role, or user.

The id of a right must be the id of the action it is acting upon.

Each profile must belongs to at least one composite (group), each role must be paired with an existing profile and each user shall have a role.

.. _setup:

SetUp
......

A list of actions must be referenced in the Rights' storages in order for it to find the actions ID and let users create new rights.

To reference a new action, simply use :

.. code-block:: python

    from canopsis.organisation.rights import Rights

    self.Rights = Rights()

    #                ACTION_ID       DESCRIPTION
    self.Rights.add('1234.ack', 'Acknowledge events')

*See the unit tests for more throgouh examples.*

.. _how:

How
----

Rights
.......

Add an action to the referenced actions list

.. code-block:: python

    def add(self, a_id, a_desc)
    """
    Args:
        a_id: id of the action to reference
        a_desc: description of the action to reference
    Returns:
        A document describing the effect of the put_elements
        if the action was created
        ``None`` otherwise
    """

    # Example
    self.Rights.add('1234.ack', 'Acknowledge events')


Check if an entity has the flags for a specific right
The entity must have a ``rights`` field with a Rights map within

.. code-block:: python

    def check(entity, right_id, checksum)
    """
    Args:
        entity: entity to be checked
        right_id: right to be checked
        checksum: minimum flags needed
    Returns:
        ``True`` if the entity has enough permissions on the right
        ``False`` otherwise
    """

    # Example
    self.Rights.check(self.Rights.get_composite('manager',
                                                '1234.ack',
                                                8)

Check if an user has the flags for a specific right
Each of the user's entities (Role, Profile, and Composites) will be checked
For now, you must specify the user's role

.. code-block:: python

    def check_rights(role, right_id, checksum)
    """
    Args:
        role: user's role to be checked
        right_id: right to be checked
        checksum: minimum flags needed
    Returns:
        ``True`` if the user's role has enough permissions
        ``False`` otherwise
    """

    # Example
    self.Rights.check_rights(self.Rights.get_role('DirectorsManager',
                                                  'management.5412',
                                                  8)


Delete the checksum of a Right from an entity

.. code-block:: python

    def delete_right(entity, e_type, right_id, checksum)
    """
    Args:
        entity: entity to delete the right from
        e_type: type of the entity
        right_id: right to be modified
        checksum: flags to remove
     Returns:
        The checksum of the right if it was modified
        ``0`` otherwise
     """

    # Example
    self.Rights.delete_right('manager', 'composite', '1234.ack', 4)




Composites
...........

Creation

.. code-block:: python

    def create_composite(comp_name, comp_rights)
    """
    Args:
        comp_name: id of the composite to create
        comp_rights: map of rights to init the composite with
    Returns:
        The name of the composite if it was created
        ``None`` otherwise
    """

    # Example
    rights = {
        '1234.ack': {
                'desc': 'create and manage ACKs',
                'checksum': 15
                },
        'management.5412': {
                'desc': 'manage list of directors',
                'checksum': '12',
                'context': 'field',
                'field': 'list_of_directors'
                }
        }

    self.Rights.create_composite('manager', rights)


Deletion

.. code-block:: python

    def delete_composite(c_name)
    """
    Args:
        c_name: id of the composite to be deleted
    Returns:
        ``True`` if the composite was deleted
        ``False`` otherwise
    """

    # Example
    self.Rights.delete_composite('manager')

Add a composite to an existing entity (Profile or Role)

.. code-block:: python

    def add_composite(e_name, e_type, comp_name, comp_rights=None)
    """
    Args:
        e_name: name of the entity to be modified
        e_type: type of the entity
        comp_name: id of the composite to add to the entity
        comp_rights: specified if the composite has to be created beforehand
    Returns:
        ``True`` if the composite was added to the entity
        ``False`` otherwise
    """

    # Example
    self.Rights.add_composite('Manager', 'profile', 'manager')
    # or
    self.Rights.add_composite('DirectorsManager', 'role', 'manager')

    # This also works, it is merely a wrapper of add_composite to make it more user-friendly
    self.Rights.add_comp_to_profile('Manager', 'manager')
    # or
    self.Rights.add_comp_to_role('DirectorsManager', 'manager')

Remove a composite from an existing entity (Profile or Role)

.. code-block:: python

    def remove_composite(e_name, e_type, comp_name)
    """
    Args:
        e_name: name of the entity to be modified
        e_type: type of the entity
        comp_name: id of the composite to remove from the entity
    Returns:
        ``True`` if the composite was removed from the entity
        ``False`` otherwise
    """

    # Example
    self.Rights.remove_composite('Manager', profile', 'manager')
    # or
    self.Rights.remove_composite('DirectorsManager', 'role', 'manager')

    # This also works, it is merely a wrapper of remove_Composite to make it more user-friendly
    self.Rights.rm_comp_profile('Manager', 'manager')
    # or
    self.Rights.rm_comp_role('DirectorsManager', 'manager')

Profiles
.........

Create a Profile

.. code-block:: python

    def create_profile(p_name, p_compites)
    """
    Args:
        p_name: id of the profile to be created
        p_compsites: list of composites to init the Profile with
    Returns:
        The name of the profile if it was created
        ``None`` otherwise
    """

    # Example
    self.Rights.create_profile('Manager', ['manager'])


Delete a Profile

.. code-block:: python

    def delete_profile(p_name)
    """
    Args:
        p_name: id of the profile to be deleted
    Returns:
        ``True`` if the profile was deleted
        ``False`` otherwise
    """

    # Example
    self.Rights.delete_profile('Manager')

Add a Profile to an existing Role

.. code-block:: python

    def add_profile(role, p_name, p_composites=None)
    """
    Args:
        role: id of the role to add the Profile to
        p_name: name of the Profile to be added
        p_composites: specified if the profile has to be created beforehand
    Returns:
        ``True`` if the profile was created
        ``False`` otherwise
    """

    # Example
    self.Rights.add_profile('DirectorsManager', 'manager')

Remove a Profile from an existing Role

.. code-block:: python

    def remove_profile(role, p_name)
    """
    Args:
        role: id of the role to remove the Profile from
        p_name: name of the Profile to be removed
    Returns:
        ``True`` if the profile was removed from the entity
        ``False`` otehrwise
    """

    # Example
    self.Rights.remove_profile('DirectorsManager', 'Manager')


Role
.....

Create a Role

.. code-block:: python

    def create_role(r_name, r_profile)
    """
    Args:
        r_name: id of the Role to be created
        r_profile: id of the Profile to init the Role with
    Returns:
        ``Name`` of the role if it was created
    """

    # Example
    self.Rights.create_role('DirectorsManager', 'Manager')


Delete a Role

.. code-block:: python

    def delete_role(r_name)
    """
    Args:
        r_name: id of the role to be deleted
    Returns:
        ``True`` if the role was deleted
        ``False`` otherwise
    """

    # Example
    self.Rights.delete_role('DirectorsManager')



.. _library:

Library
-------

.. _structures:

Data Structures
................

User
,,,,,

.. code-block:: javascript

    User = {

        'role': ...,                 // List of role names that defines the User's profile, groups, and rights
        'contact': {                 // Map of contact informations
            'mail': ...,
            'phone_number': ...,
            ...
            }
        'name': ...,                 // String of user's name
        '_id': ...                   // uniq id

        // Empty by default
        'rights': ...,               // Map of type Rights, every user-specific rights goes here
        'groups': ...,               // List of group names, every user-specific groups goes here
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
,,,,

A Role is specific to a small number of users

.. code-block:: javascript

    'name': {

        'profile': ...              // ID of the profile (string)

        // Empty by default
        'rights': ...               // Map of type Rights, every role-specific rights goes here
        FIELD: ...                  // You can add any number of fields that can be used with data-specific rules
        ...

        }


Example:

.. code-block:: javascript

    Roles = {
        'manager': {
            'profile': 'DirectorsManager',
            'list_of_directors': ['Ted Chaough', 'Peggy Olson', 'Don Draper']
            }
        }


Profile
,,,,,,,,

A profile is generic and global to all users

.. code-block:: javascript

    'name': {                            // String of profile's name

        'composites': ...                // List of the groups the profile belongs to

        // Empty by default
        'rights': ...               // Map of type Rights, every profile-specific rights goes here

        }



Example:

.. code-block:: javascript

    An Administrator profile exists, it has all rights and belongs to the Group Management as well as the root Group
    Profiles = {
        'Manager': {
            'composites': ['managements', 'supervizion']
        }



Composite (aka Groups)
,,,,,,,,,,,,,,,,,,,,,,

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
                    'desc': ['Access and change directors configuration']
                }
            }
        }
    }


Rights
,,,,,,

.. code-block:: javascript

    Rghts = {
        object_id...: {             // Right on the object with the identifier id

            'checksum': ...,        // 1 == Read, 2 == Update, 4 == Create, 8 == Delete

            // Additional Field
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

