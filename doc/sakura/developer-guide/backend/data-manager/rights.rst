.. _dev-backend-mgr-rights:

Rights
_______

Table of contents
-------------------

1. What_

   1. SetUp_

2. How_

   1. Rights

   2. Groups

   3. Profile

   4. Role

   5. Update_

3. Library_

   1. Data Structures_

      1. User

      2. Role

      3. Profile

      4. Group

      5. Rights


.. _what:

What
-----

Rights are defined within the groups (groups) upon their creations or added unitarily in a specific profile, role, or user.

The id of a right must be the id of the action it is acting upon.

Each profile must belongs to at least one group (group), each role must be paired with an existing profile and each user shall have a role.

.. _setup:

SetUp
......

A list of actions must be referenced in the Rights' storages in order for it to find the actions ID and let users create new rights.

To reference a new action, simply use :

.. code-block:: python

    from canopsis.organisation.rights import Rights

    right_module = Rights()

    #                ACTION_ID       DESCRIPTION
    right_module.add('1234.ack', 'Acknowledge events')

*See the unit tests for more thorough examples.*

.. _how:

How
----

Here will be documented the methods needed to develop a proper API to interact with the Rights module.

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
    right_module.add('1234.ack', 'Acknowledge events')


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
    right_module.check(right_module.get_group('manager',
                                                '1234.ack',
                                                8)

Check if an user has the flags for a specific right
Each of the user's entities (Role, Profile, and Groups) will be checked
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
    right_module.check_rights(right_module.get_role('DirectorsManager',
                                                  'management.5412',
                                                  8)


Groups
.......

Creation

.. code-block:: python

    def create_group(group_name, group_rights)
    """
    Args:
        group_name: id of the group to create
        group_rights: map of rights to init the group with
    Returns:
        The name of the group if it was created
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

    right_module.create_group('manager', rights)


Deletion

.. code-block:: python

    def delete(e_type, e_id)
    """
    Args:
        e_type type of the entity to delete
        e_id id of the entity to delete
    Returns:
        ``True`` if the entity was deleted
        ``False`` otherwise
    """

    # Example
    right_module.delete('group', 'manager')


Profiles
.........

Create a Profile

.. code-block:: python

    def create_profile(p_name, p_groups)
    """
    Args:
        p_name: id of the profile to be created
        p_groups: list of groups to init the Profile with
    Returns:
        The name of the profile if it was created
        ``None`` otherwise
    """

    # Example
    right_module.create_profile('Manager', ['manager'])


Deletion

.. code-block:: python

    def delete(e_type, e_id)
    """
    Args:
        e_type type of the entity to delete
        e_id id of the entity to delete
    Returns:
        ``True`` if the entity was deleted
        ``False`` otherwise
    """

    # Example
    right_module.delete('profile', 'Manager')

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
    right_module.create_role('DirectorsManager', 'Manager')


Delete a Role

.. code-block:: python

    def delete(e_type, e_id)
    """
    Args:
        e_type type of the entity to delete
        e_id id of the entity to delete
    Returns:
        ``True`` if the entity was deleted
        ``False`` otherwise
    """

    # Example
    right_module.delete('DirectorsManager')


.. _update:

Update
.......

Groups
,,,,,,,

Update the groups of an entity
Every groups present in the groups list and not in the entity will be added
Every groups present in the entity and not in the groups list will be deleted

If the groups list is None, nothing will be deleted

.. code-block:: python

    def update_comp(self, e_id, e_type, groups, entity):
        """
        Args:
            e_id id of the entity to update
            e_type type of the entity to update
            groups groups to update
            entity entity to be updated
        Returns:
            ``True`` if the entity was thoroughly updated
            ``False`` otherwise
        """

    # Example, add groups to the 'Manager' profile
    right_module.update_comp('Manager', 'profile',
                             ['supervision', 'visualisation'],
                             right_module.get_profile('Manager'))

Rights
,,,,,,

Update the rights of an entity
Every rights present in the rights list and not in the entity will be added
Every rights present in the entity and not in the rights list will be deleted

If the rights list is None, nothing will be deleted

.. code-block:: python

    def update_rights(self, e_id, e_type, rights, entity):
        """
        Args:
            e_id id of the entity to update
            e_type type of the entity to update
            rights rights to update
            entity entity to be updated
        Returns:
            ``True`` if the entity was thoroughly updated
            ``False`` otherwise
        """

    # Example, add rights to the 'Manager' profile
    right_module.update_rights('Manager', 'profile',
                               {'eventview' : {'checksum': 15}},
                               right_module.get_profile('Manager'))

General
,,,,,,,,

Update several fields by overriding previous values

.. code-block:: python

    def update_fields(e_id, e_type, fields):
        """
        Args:
            e_id id of the entity to update
            e_type type of the entity to update
            fields map of the fields to update
        Returns:
            A document describing the effect of the put_elements
            if the entity was thoroughly updated
            ``False`` otherwise
        """

    # Example, change name of user and its role
    right_module.update_fields('Joan Harris', 'user',
                               {'crecord_name': 'Peggy Olson',
                                'role': 'CreativeDirector'})


Update a single ``list`` or ``dictionary`` field by adding or deleting differences in values

.. code-block:: python

    def update_field(self, e_id, e_type, new_elems, elem_type, entity):
        """
        Args:
            e_id id of the entity to update
            e_type type of the entity to update
            new_elems elements to update
            entity entity to be updated
        Returns:
            ``True`` if the entity was thoroughly updated
            ``False`` otherwise
        """

    # Example, add groups to the 'Manager' profile
    right_module.update_comp('Manager', 'profile',
                             ['supervision', 'visualisation'],
                             'group',
                             right_module.get_profile('Manager'))
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

        'crecord_type': 'user',
        'crecord_name': ...,                 // String of user's name

        'role': ...,                 // Role of the user that defines his profiles and groups
        'contact': {                 // Map of contact informations
            'mail': ...,
            'phone_number': ...,
            ...
            }
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

        'crecord_type': 'user',
        'crecord_name': 'Joan Harris',
        'role': 'manager',
        'contact': {
            'mail': 'jharris@scdp.com',
            'phone_number': '+33678695041',
            'adress': '1271 6th Avenue, Rockefeller Center, NYC, New York'
            }
        '_id': '1407160264.joan.harris.manager'

        }


Role
,,,,

A Role is specific to a small number of users

.. code-block:: javascript

    'name': {

        'crecord_type': 'role',
        'crecord_name':             // Name of the role

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

        'crecord_type': 'profile',
        'crecord_name':             // Name of the profile

        'groups': ...                // List of the groups the profile belongs to

        // Empty by default
        'rights': ...               // Map of type Rights, every profile-specific rights goes here

        }



Example:

.. code-block:: javascript

    An Administrator profile exists, it has all rights and belongs to the Group Management as well as the root Group
    Profiles = {
        'Manager': {
            'groups': ['managements', 'supervizion'],
            'crecord_type': 'profile',
            'crecord_name': 'Manager'
        }



Group (aka Groups)
,,,,,,,,,,,,,,,,,,,,,,

A group is generic and global to all users

.. code-block:: javascript

    'name': {                        // String of group's name

        'crecord_type': 'group',
        'crecord_name':             // Name of the group

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
            'crecord_type': 'group',
            'crecord_name': 'management'
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

            'crecord_type': 'right'

            }
        }

The keys of a map of type ``Rights`` are the ids of the objects accessible from the web application.
The ``right`` field is a 4-bit integer that goes from 1 to 15 and that describes the available action on the object.


.. code-block:: python

    if Rights[object_idXYZ]['right'] & (READ | CREATE | UPDATE | DELETE) == (READ | CREATE | UPDATE | DELETE):
        #the user has all rights on the object identified with object_idXYZ

    if not Rights[object_idXYZ]['right'] & (CREATE | DELETE):
        #the user has none of the rights on the object identified with object_idXYZ


Functions
..........

Rights
,,,,,,,


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
    right_module.delete_right('manager', 'group', '1234.ack', 4)

Groups
,,,,,,,


Add a group to an existing entity (Profile or Role)

.. code-block:: python

    def add_group(e_name, e_type, group_name, group_rights=None)
    """
    Args:
        e_name: name of the entity to be modified
        e_type: type of the entity
        group_name: id of the group to add to the entity
        group_rights: specified if the group has to be created beforehand
    Returns:
        ``True`` if the group was added to the entity
        ``False`` otherwise
    """

    # Example
    right_module.add_group('Manager', 'profile', 'manager')
    # or
    right_module.add_group('DirectorsManager', 'role', 'manager')

    # This also works, it is merely a wrapper of add_group to make it more user-friendly
    right_module.add_group_to_profile('Manager', 'manager')
    # or
    right_module.add_group_to_role('DirectorsManager', 'manager')

Remove a group from an existing entity (Profile or Role)

.. code-block:: python

    def remove_group(e_name, e_type, group_name)
    """
    Args:
        e_name: name of the entity to be modified
        e_type: type of the entity
        group_name: id of the group to remove from the entity
    Returns:
        ``True`` if the group was removed from the entity
        ``False`` otherwise
    """

    # Example
    right_module.remove_group('Manager', profile', 'manager')
    # or
    right_module.remove_group('DirectorsManager', 'role', 'manager')

    # This also works, it is merely a wrapper of remove_Group to make it more user-friendly
    right_module.rm_group_profile('Manager', 'manager')
    # or
    right_module.rm_group_role('DirectorsManager', 'manager')

profiles
,,,,,,,,,,


Add a Profile to an existing Role

.. code-block:: python

    def add_profile(role, p_name, p_groups=None)
    """
    Args:
        role: id of the role to add the Profile to
        p_name: name of the Profile to be added
        p_groups: specified if the profile has to be created beforehand
    Returns:
        ``True`` if the profile was created
        ``False`` otherwise
    """

    # Example
    right_module.add_profile('DirectorsManager', 'manager')

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
    right_module.remove_profile('DirectorsManager', 'Manager')
