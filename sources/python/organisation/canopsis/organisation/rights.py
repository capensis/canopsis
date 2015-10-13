# -*- coding: utf-8 -*-
# --------------------------------
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
# ---------------------------------

from logging import ERROR
from uuid import uuid4
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category)
from canopsis.middleware.registry import MiddlewareRegistry

CATEGORY = 'RIGHTS'


@conf_paths('organisation/rights.conf')
@add_category(CATEGORY)
class Rights(MiddlewareRegistry):

    DATA_SCOPE = 'rights'

    def __init__(
        self, data_scope=DATA_SCOPE,
        logging_level=ERROR,
        *args, **kwargs
    ):

        super(Rights, self).__init__(data_scope=data_scope, *args, **kwargs)

    # Generic getter
    def get_from_storage(self, s_type):
        def get_from_storage_(elem):
            return self[s_type + '_storage'].get_elements(
                ids=elem, query={'crecord_type': s_type})
        return get_from_storage_

    def get_users(self, projection={'_id': 1}):
        return self['user_storage'].get_elements(
            query={'crecord_type': 'user'},
            projection=projection
        )

    def _configure(self, unified_conf, *args, **kwargs):

        super(Rights, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        self.profile_storage = self['profile_storage']
        self.group_storage = self['group_storage']
        self.role_storage = self['role_storage']
        self.action_storage = self['action_storage']
        self.user_storage = self['user_storage']

        self.get_profile = self.get_from_storage('profile')
        self.get_action = self.get_from_storage('action')
        self.get_group = self.get_from_storage('group')
        self.get_role = self.get_from_storage('role')
        self.get_user = self.get_from_storage('user')

        self.actions = {
            'remove': {
                'profile': self.remove_profile,
                'group': self.remove_group,
                },
            'add': {
                'profile': self.add_profile,
                'group': self.add_group,
                },
            'delete': {
                'profile': self.delete_profile,
                'group': self.delete_group,
                'user': self.delete_user,
                'role': self.delete_role,
                'action': self.delete_action
                }
            }

    # Add an action to the referenced action
    def add(self, a_id, a_desc):
        """
        Args:
            a_id: id of the action to reference
            a_desc: description of the action to reference
        Returns:
            A document describing the effect of the put_elements
            if the action was created
            ``None`` otherwise
        """

        return self['action_storage'].put_element(
            element={
                'crecord_name': a_id,
                'crecord_type': 'action',
                'desc': a_desc
            }, _id=a_id
        )

    # Delete an action from the reference list
    def delete_action(self, a_id):
        """
        Args:
            a_id: id of the action to be deleted
        Returns:
            See MongoStorage
        """

        return self['action_storage'].remove_elements(a_id)

    # Check if an entity has the flags for a specific rigjt
    # The entity must have a rights field with a rights maps within
    def check(self, entity, right_id, checksum):
        """
        Args:
            entity: entity to be checked
            right_id: right to be checked
            checksum: minimum flags needed
        Returns:
            ``True`` if the entity has enough permissions on the right
            ``False`` otherwise
        """

        if not entity or not entity.get('rights', None):
            self.logger.info(
                'Entity empty or has no rights field, ' +
                'can not check right {}'.format(right_id)
                )
            return False

        found = entity['rights'].get(right_id, None)
        if (found and found.get('checksum', 0) & checksum >= checksum):
            return True

        return False

    # Check if an user has the flags for a specific right
    # Each of the user's entities (Role, Profile, and Groups)
    # will be checked
    def check_rights(self, u_name, right_id, checksum):
        """
        Args:
            u_name: user to be checked
            right_id: right to be checked
            checksum: minimum flags needed
        Returns:
            ``True`` if the user has enough permissions
            ``False`` otherwise
        """

        user = self.get_user(u_name)
        role = None
        if user:
            role = self.get_role(user.setdefault('role', None))
        profiles = self.get_profile(role['profile'])

        # Do not edit the following for a double for loop
        # list grouprehensions are much faster
        groups = [self['group_storage'][x]
                  for y in profiles
                  for x in y['group']]

        if 'group' in role:
            groups += [self['group_storage'][x]
                       for x in role['group']]
        if 'group' in user:
            groups += [self['group_storage'][x]
                       for x in user['group']]

        # check in the role's comsposite
        if ((user and self.check(user, right_id, checksum)) or
            (role and self.check(role, right_id, checksum)) or
            # check in the profile's group
            (len(profiles) and any(self.check(x, right_id, checksum)
                                   for x in profiles)) or
            # check in the profile's groups groups
            (len(groups) and any(self.check(x, right_id, checksum)
                                 for x in groups))):
            return True

        return False

    # Add a right to the entity linked
    # If the right already exists, the checksum will be summed accordingly
    # checksum |= old_checksum
    # entity can be a role, a profile, or a group
    def add_right(self, e_name, e_type, right_id, checksum, **kwargs):
        """
        Args:
            e_name: name of the entity to add the right to
            e_type: type of the entity
            right_id: right to be modified
            checksum: flags to add
        Returns:
            The checksum of the right if the flags were added
            ``0`` otherwise
        """

        # Action not referenced, can't create a right
        if not self.get_action(right_id):
            self.logger.error(
                ('Can not create right, the action {0} ' +
                    'is not referenced').format(right_id)
                )
            return 0

        entity = None

        e_type += '_storage'

        if e_type in self:
            entity = self[e_type].get_elements(ids=e_name)

        if not entity:
            self.logger.error(
                '{0} {1} {2}.'.format(
                    'Can not create right, entity',
                    e_name,
                    'is empty or does not exist.'
                )
            )
            return 0

        if not entity.get('rights', None):
            entity['rights'] = {}

        # If it does not exist, create it
        if not self.check(entity, right_id, 0):
            entity['rights'].update({
                right_id: {
                    'crecord_type': 'right',
                    'checksum': checksum
                }
            })
        else:
            entity['rights'][right_id]['checksum'] = checksum

        # Add the new context and other fields, if any
        for key in kwargs:
            if kwargs[key]:
                entity['rights'][right_id][key] = context

        self[e_type].put_element(element=entity, _id=e_name)
        result = entity['rights'][right_id]['checksum']
        return result if result else True

    # Delete the checksum right of the entity linked
    # new_checksum ^= checksum
    def remove_right(self, entity, e_type, right_id, checksum):
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

        entity = self[e_type + '_storage'].get_elements(ids=entity)

        if (
                entity['rights']
                and entity['rights'][right_id]
                and entity['rights'][right_id]['checksum'] >= checksum
        ):

            # remove the permissions passed in checksum
            ochecksum = int(entity['rights'][right_id]['checksum'])
            entity['rights'][right_id]['checksum'] = ochecksum ^ int(checksum)

            # If all the permissions were removed from the right, delete it
            if not entity['rights'][right_id]['checksum']:
                del entity['rights'][right_id]
                self[e_type + "_storage"].put_element(
                    element=entity, _id=entity['_id']
                )
                return True

            self[e_type + "_storage"].put_element(
                element=entity, _id=entity['_id']
            )
            result = entity['rights'][right_id]['checksum']
            return result if result else True

        return 0

    # Create a new rights group grouposed of the rights passed in group_rights
    # group_rights should be a map of rights referenced in the action catalog
    def create_group(self, group_name, group_rights):
        """
        Args:
            group_name: id of the group to create
            group_rights: map of rights to init the group with
        Returns:
            The name of the group if it was created
            ``None`` otherwise
        """

        # Do nothing if it already exists
        if self.get_group(group_name):
            self.logger.error(
                'Can not create group, group {0} already exists'.format(
                    group_name)
                )
            return None

        new_group = {
            'crecord_type': 'group',
            'crecord_name': group_name,
            'rights': {}
        }

        self.group_storage.put_element(element=new_group, _id=group_name)

        if not group_rights:
            return group_name

        # Use add_right to check if the action is referenced
        for right_id in group_rights:
            self.add_right(group_name,
                           'group',
                           right_id,
                           group_rights[right_id]['checksum'])

        return group_name

    # Create a new profile grouposed of the groups p_groups
    #   and which name will be p_name
    # If the profile already exists, groups from p_groups
    #   that are not already in the profile's groups will be added
    def create_profile(self, p_name, p_groups):
        """
        Args:
            p_name: id of the profile to be created
            p_groupsites: list of groups to init the Profile with
        Returns:
            The name of the profile if it was created
            ``None`` otherwise
        """

        # Do nothing if it already exists
        if self.get_profile(p_name):
            self.logger.error(
                'Can not create group, group {0} already exists'
                .format(group_name)
            )
            return None

        new_profile = {
            'crecord_type': 'profile',
            'crecord_name': p_name
        }
        if isinstance(p_groups, list):
            new_profile['group'] = p_groups
        else:
            new_profile.setdefault('group', []).append(p_groups)

        self.profile_storage.put_element(element=new_profile, _id=p_name)

        return p_name

    # Delete entity of id e_name
    # t_type is the storage to check for relations
    # entity can be a profile, or group
    def delete_entity(self, e_name, e_type):
        """
        Args:
            e_name: id of the entity to be deleted
            e_type: type of the entity
        Returns:
            ``True`` if the entity was deleted
            ``False`` otherwise
        """

        from_storage = e_type + '_storage'
        t_type = 'profile' if e_type == 'group' else 'role'
        to_storage = t_type + '_storage'

        if self[from_storage].get_elements(ids=e_name):
            self[from_storage].remove_elements(e_name)

            # remove the entity from every other entities that use it
            for entity in self[to_storage].get_elements(
                    query={'crecord_type': t_type}):
                if e_type in entity and e_name in entity[e_type]:
                    entity[e_type].remove(e_name)
                    self[to_storage].put_element(
                        _id=entity['_id'], element=entity
                    )

            return True

        self.logger.error(
            'Can not delete entity, entity {0} does not exist'.format(e_name)
        )
        return False

    def delete_role(self, r_name):
        """
        Args:
            r_name: id of the role to be deleted
        Returns:
            ``True`` if the role was deleted
            ``False`` otherwise
        """

        if self.get_role(r_name):
            for user in self['user_storage'].get_elements(
                    query={'crecord_type': 'user'}
            ):
                if user and 'role' in user and r_name == user['role']:
                    user.pop('role', None)
                    self['user_storage'].put_element(
                        _id=user['_id'], element=user
                    )

            self['role_storage'].remove_elements(r_name)
            return True

        self.logger.error(
            'Can not delete role, role {0} does not exist'.format(r_name)
        )
        return False

    def delete_user(self, u_name):
        """
        Args:
            u_name: id of the name to be deleted
        Returns:
            ``True`` if the role was deleted
            ``False`` otherwise
        """

        return self['user_storage'].remove_elements(u_name)

    # delete_entity wrapper
    def delete_profile(self, p_name):
        """
        Args:
            p_name: id of the profile to be deleted
        Returns:
            ``True`` if the profile was deleted
            ``False`` otherwise
        """

        return self.delete_entity(p_name, 'profile')

    # delete_entity wrapper
    def delete_group(self, c_name):
        """
        Args:
            c_name: id of the group to be deleted
        Returns:
            ``True`` if the group was deleted
            ``False`` otherwise
        """

        return self.delete_entity(c_name, 'group')

    # Add the group named group_name to the entity
    # If the group does not exist and
    #   group_rights is specified it will be created first
    # entity can be a profile or a role
    def add_group(self, e_name, e_type, group_name, group_rights=None):
        """
        Args:
            e_name: name of the entity to be modified
            e_type: type of the entity
            group_name: id of the group to add to the entity
            group_rights: specified if the group has to be created beforehand
        Returns:
            ``True`` if the group was added to the entity
        """

        e_type += '_storage'

        if not self.get_group(group_name):
            self.create_group(group_rights, group_name)

        entity = self[e_type].get_elements(ids=e_name)
        if 'group' not in entity or group_name not in entity['group']:
            entity.setdefault('group', []).append(group_name)
            self[e_type].put_element(_id=e_name, element=entity)

        return True

    # add_group wrapper
    def add_group_profile(self, e_name, group_name, group_rights=None):
        """
        Args:
            e_name: profile id to add the group to
            group_name: group to be added
            group_rights: specified if the group has to be created beforehand
        Returns:
            ``True`` if the group was added to the profile
            ``False`` otherwise
        """

        return self.add_group(e_name, 'profile', group_name, group_rights)

    # add_group wrapper
    def add_group_role(self, e_name, group_name, group_rights=None):
        """
        Args:
            e_name: role id to add the group to
            group_name: group to be added
            group_rights: specified if the group has to be created beforehand
        Returns:
            ``True`` if the group was added to the role
            ``False`` otherwise
        """

        return self.add_group(e_name, 'role', group_name, group_rights)

    # add_group wrapper
    def add_group_user(self, e_name, group_name, group_rights=None):
        """
        Args:
            e_name: user id to add the group to
            group_name: group to be added
            group_rights: specified if the group has to be created beforehand
        Returns:
            ``True`` if the group was added to the user
            ``False`` otherwise
        """

        return self.add_group(e_name, 'user', group_name, group_rights)

    # Add the profile of name p_name to the role
    # If the profile does not exists and p_groups is specified
    #    it will be created first
    def add_profile(self, role, e_type, p_name, p_groups=None):
        """
        Args:
            role: id of the role to add the Profile to
            p_name: name of the Profile to be added
            p_groups: specified if the profile has to be created beforehand
        Returns:
            ``True`` if the profile was created
            ``False`` otherwise
        """

        profile = self.get_profile(p_name)
        if not profile:
            self.create_profile(p_name, p_groups)

        # retrieve the profile
        if profile:
            s_role = self.get_role(role)

            if 'profile' not in s_role or p_name not in s_role['profile']:
                s_role.setdefault('profile', []).append(p_name)
                self.role_storage.put_element(_id=role, element=s_role)

            return True

    # Add the profile of name p_name to the role
    # If the profile does not exists and p_groups is specified
    #    it will be created first
    def add_role(self, u_name, r_name, r_profile=None):
        """
        Args:
            u_name: id of the user to add the role to
            r_name: name of the role to be added
            r_groups: specified if the role has to be created beforehand
        Returns:
            ``True`` if the profile was created
            ``False`` otherwise
        """

        role = self.get_role(r_name)
        if not role:
            self.create_role(r_name, r_profile)

        # retrieve the profile
        if role:
            s_user = self.get_user(u_name)
            s_user['role'] = r_name
            self.user_storage.put_element(_id=u_name, element=s_user)

            return True

    # Remove the entity e_name from from_name
    # from_name can be a profile or a role
    # e_name can be a profile or a group
    def remove_entity(self, from_name, from_type, e_name, e_type):
        entity = self[from_type + '_storage'].get_elements(
            query={'crecord_type': from_type}, ids=from_name
        )

        if e_type in entity and e_name in entity[e_type]:
            entity[e_type].remove(e_name)
            self[from_type + '_storage'].put_element(
                _id=from_name, element=entity
            )
            return True

        return False

    # remove_entity wrapper
    def remove_group(self, e_name, e_type, group_name):
        """
        Args:
            e_name: name of the entity to be modified
            e_type: type of the entity
            group_name: id of the group to remove from the entity
        Returns:
            ``True`` if the group was removed from the entity
            ``False`` otherwise
        """

        return self.remove_entity(e_name, e_type, group_name, 'group')

    # remove_group wrapper
    def remove_group_role(self, r_name, c_name):
        """
        Args:
            r_name: role to removed the group from
            c_name: group to remove
        Return:
            ``True`` if the group was removed from the role
            ``False`` otherwise
        """

        return self.remove_group(r_name, 'role', c_name)

    # remove_group wrapper
    def remove_group_profile(self, p_name, c_name):
        """
        Args:
            p_name: profile to removed the group from
            c_name: group to remove
        Return:
            ``True`` if the group was removed from the profile
            ``False`` otherwise
        """

        return self.remove_group(p_name, 'profile', c_name)

    # remove_group wrapper
    def remove_group_user(self, u_name, c_name):
        """
        Args:
            u_name: user to removed the group from
            c_name: group to remove
        Return:
            ``True`` if the group was removed from the profile
            ``False`` otherwise
        """

        return self.remove_group(u_name, 'user', c_name)

    # remove_entity wrapper
    def remove_profile(self, r_name, e_type, p_name):
        """
        Args:
            r_name: id of the role to remove the Profile from
            p_name: name of the Profile to be removed
        Returns:
            ``True`` if the profile was removed from the entity
            ``False`` otehrwise
        """

        return self.remove_entity(r_name, 'role', p_name, 'profile')

    # remove_entity wrapper
    def remove_role(self, u_name, r_name):
        """
        Args:
            u_name: id of the user to remove the role from
            r_name: name of the role to be removed
        Returns:
            ``True`` if the role was removed from the entity
            ``False`` otehrwise
        """

        return self.remove_entity(u_name, 'user', r_name, 'role')

    # Create a new role grouposed of the profile r_profile
    #   and which name will be r_name
    # Any extra field can be specified in the kwargs
    # If the role already exists, the profile will be changed for r_profile
    def create_role(self, r_name, r_profile):
        """
        Args:
            r_name: id of the Role to be created
            r_profile: id of the Profile to init the Role with
        Returns:
            ``Name`` of the role if it was created
        """

        if self.get_role(r_name):
            return r_name

        new_role = {
            'crecord_type': 'role',
            'crecord_name': r_name,
        }
        if isinstance(r_profile, list):
            new_role['profile'] = r_profile
        else:
            new_role.setdefault('profile', []).append(r_profile)

        self.role_storage.put_element(_id=r_name, element=new_role)

        return r_name

    def create_user(
            self, u_id, u_role,
            contact=None, rights=None,
            groups=None
    ):
        """
        Args:
            u_id: nick of the user to create, usually first letter of first
                name and last name (i.e.: jdoe for John Doe).

            u_role: role to init the user with.
            contact: map containing full name, email, adress, and/or phone
                number of the user.

            rights: map containing specific rights.
            groups: list of specific groups.
        Returns:
            Map of the newly created user.
        """

        user = self.get_user(u_id)

        if user:
            return user

        user = {
            'crecord_type': 'user',
            'enable': True,
            'crecord_name': u_id,
            'role': u_role,
            'authkey': str(uuid4())
        }

        if contact and isinstance(contact, dict):
            user['contact'] = contact

        if rights and isinstance(rights, dict):
            user['rights'] = rights

        if groups and isinstance(groups, list):
            user['groups'] = groups

        self.user_storage.put_element(_id=u_id, element=user)
        return user

    def set_user_fields(self, u_id, fields):
        """
        Args:
            u_id: id of the user which fields to change.
            fields: map of fields to change and their new values.
        Returns:
            Map of the modified user.
        """

        user = self.get_user(u_id)

        supported_fields = {'name', 'email', 'address', 'phone'}

        for key in fields:
            if key in supported_fields:
                user.setdefault('contact', {})[key] = fields[key]

        self.user_storage.put_element(_id=u_id, element=user)
        return user

    def get_user_rights(self, u_id):
        """
        Args:
            u_uid: id of the user to get the rights from.
        Returns:
            dict of user's rights.
        """

        profiles = []
        n_groups = []
        user = self.get_user(u_id)

        if not user:
            return {}

        role = user.setdefault('role', None)
        if role:
            role = self.get_role(role)

            if role and 'profile' in role:
                profiles = self.get_profile(role['profile'])
                n_groups = [x for y in profiles for x in y['group']]

            if role and 'group' in role:
                n_groups += role['group']

        if 'group' in user:
            n_groups += user['group']

        specific_rights = [self['group_storage'][x]['rights']
                           for x in set(n_groups)]

        specific_rights.append(user.setdefault('rights', {}))
        if role:
            specific_rights.append(role.setdefault('rights', {}))
        if profiles:
            (specific_rights.append(x.setdefault('rights', {}))
             for x in profiles)

        rights = {}
        for e_rights in specific_rights:
            for r_id in e_rights:
                if r_id in rights:
                    rights[r_id]['checksum'] |= e_rights[r_id]['checksum']
                else:
                    rights[r_id] = e_rights[r_id]

        return rights

    def get_entity_field(self, e_id, e_type, field):
        """
        Args:
            e_id: entity to get the field from.
            e_type: type of the entity.
            field: field to get.
        Returns:
            value of the field if the field exists if the entity e_id.
            ``None`` otherwise.
        """

        if not field or not e_id or not e_type:
            return None

        entity = self[e_type + '_storage'].get_elements(
            ids=e_id, query={'crecord_type': e_type}
            )

        return entity.setdefault(field, None)

    # Update entity name
    def update_entity_name(self, e_id, e_type, new_name):
        """
        Args:
            e_id: id of the entity to update.
            e_type: type of the entity to update.
            new_name: new name of the entity.
        Returns:
            ``True`` if the name was updated.
            ``False`` otherwise.
        """

        entity = self[e_type + '_storage'].get_elements(
            query={'crecord_type': e_type}, ids=e_id)

        if entity:
            entity['crecord_name'] = new_name
            self[e_type + '_storage'].put_element(_id=e_id, element=entity)
            return True

        return False

    def update_field(self, e_id, e_type, new_elems, elem_type, entity):
        """
        Args:
            e_id: id of the entity to update.
            e_type: type of the entity to update.
            new_elems: elements to update.
            entity: entity to be updated.
        Returns:
            ``True`` if the entity was thoroughly updated.
            ``False`` otherwise.
        """

        if entity and elem_type in entity:
            to_remove = [] if new_elems is None else entity[elem_type]
            if new_elems:
                to_remove = set(entity[elem_type]) - set(new_elems)
            for elem in to_remove:
                if not self.actions['remove'][elem_type](e_id, e_type, elem):
                    return False
        if new_elems:
            for elem in new_elems:
                if not self.actions['add'][elem_type](e_id, e_type, elem):
                    return False
        return True

    def update_rights(self, e_id, e_type, e_rights, entity):
        """
        Args:
            e_id id of the entity to update
            e_type type of the entity to update
            e_rights rights to update
            entity entity to be updated
        Returns:
            ``True`` if the entity was thoroughly updated
            ``False`` otherwise
        """

        if entity and 'rights' in entity:
            to_remove = [] if e_rights is None else entity['rights']
            if e_rights:
                to_remove = set(entity['rights']) - set(e_rights)
            for right in to_remove:
                if not self.remove_right(
                        e_id,
                        e_type,
                        right,
                        entity['rights'][right]['checksum']
                ):
                    return False
        if e_rights:
            for right in e_rights:
                if not self.add_right(
                        e_id,
                        e_type,
                        right,
                        e_rights[right]['checksum']
                ):
                    return False
        return True

    def update_profile(self, e_id, e_type, profiles, entity):
        """
        Args:
            e_id id of the entity to update
            e_type type of the entity to update
            profiles profiles to update
            entity entity to be updated
        Returns:
            ``True`` if the entity was thoroughly updated
            ``False`` otherwise
        """

        return self.update_field(e_id, e_type, profiles, 'profile', entity)

    def update_group(self, e_id, e_type, groups, entity):
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

        return self.update_field(e_id, e_type, groups, 'group', entity)

    def update_fields(self, e_id, e_type, fields):
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

        entity = self[e_type + '_storage'].get_elements(ids=e_id)

        if entity and not isinstance(entity, list):
            for key in fields:
                entity[key] = fields[key]
            return self[e_type + '_storage'].put_element(
                _id=e_id, element=entity
            )
        else:
            return False

    def delete(self, e_type, e_id):
        """
        Args:
            e_type type of the entity to delete
            e_id id of the entity to delete
        Returns:
            ``True`` if the entity was deleted
            ``False`` otherwise
        """

        return not not (None if e_type not in self.actions['delete']
                        else self.actions['delete'][e_type](e_id))
