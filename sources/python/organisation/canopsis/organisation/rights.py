# -*- coding: utf-8 -*-
#--------------------------------
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

from logging import DEBUG, ERROR, getLogger

from canopsis.configuration import conf_paths, add_category
from canopsis.middleware.manager import Manager

CATEGORY = 'RIGHTS'

@conf_paths('organisation/rights.conf')
@add_category(CATEGORY)
class Rights(Manager):

    DATA_SCOPE = 'rights'

    def __init__(
        self, data_scope=DATA_SCOPE,
        logging_level=ERROR,
        *args, **kwargs
    ):

        super(Rights, self).__init__(data_scope=data_scope, *args, **kwargs)


    def _configure(self, unified_conf, *args, **kwargs):

        super(Rights, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        self.profile_storage = self['profile_storage']
        self.composite_storage = self['composite_storage']
        self.role_storage = self['role_storage']
        self.action_storage = self['action_storage']
        self.default_profile = 'vizualisation'

        self.get_profile = self.get_from_storage('profile_storage')
        self.get_action = self.get_from_storage('action_storage')
        self.get_composite = self.get_from_storage('composite_storage')
        self.get_role = self.get_from_storage('role_storage')


    # Add an action to the referenced action
    def add(self, a_id, a_desc):
        action = {'type': 'action',
                 'desc': a_desc}
        self['action_storage'].put_element(a_id, action)

    # Entity can be a right_composite, a profile or a role since
    # all 3 of them have a rights field
    # Called by check_user_rights
    def check(self, entity, right_id, checksum):
        """
        Check the from the rights of entity for a right of id right_id
        """

        # check the fields
        if not entity or not entity.get('rights', None):
            return False


        found = entity['rights'].get(right_id, None)
        if (found and found.get('checksum', 0) & checksum >= checksum):
            return True

        return False


    # Check in the rights of the user (in the role), then in the profile,
    # then in the rights_composite
    # Return the value as soon as it's found
    # True if found and user has the right else False
    def check_user_rights(self, role, right_id, checksum):
        """
        Check if user has the right of id right_id
        """

        role = self.profile_storage.get_elements(ids=role)
        profiles = self.profile_storage.get_elements(ids=role['profile'])

        # Do not edit the following for a double for loop
        # list comprehensions are much faster
        composites = [self['composite_storage'][x]
                      for y in profiles
                      for x in y['composite']]

        # check in the role's comsposite
        if ((role and self.check(role, right_id, checksum)) or
            # check in the profile's composite
            (len(profiles) and any(self.check(x, right_id, checksum)
                                   for x in profiles)) or
            # check in the profile's groups composites
            (len(composites) and any(self.check(x, right_id, checksum)
                                     for x in composites))):
            return True

        return False

    # Add a right to the entity linked
    # If the right already exists, the checksum will be summed accordingly
    # checksum |= old_checksum
    # entity can be a role, a profile, or a composite
    def add_right(self, e_name, e_type, right_id, checksum,
            **kwargs):

        # Action not referenced, can't create a right
        if not self['action_storage'].get_elements(
            ids=right_id, query={'type':'action'}):
            print (
                'Can not add right {0} to entity {1}: action is not referenced'.format(right_id, e_name))
            return False

        entity = None

        e_type += '_storage'

        if e_type in self:
            entity = self[e_type].get_elements(ids=e_name)

        if not entity:
            return False

        if not entity.get('rights', None):
            entity['rights'] = {}

        # If it does not exist, create it
        if not self.check(entity, right_id, 0):
            entity['rights'].update({right_id: {'type': 'right',
                                                'checksum': checksum
                                                }
                                     })
        else:
            entity['rights'][right_id]['checksum'] |= checksum

        # Add the new context and other fields, if any
        for key in kwargs:
            if kwargs[key]:
                entity['rights'][right_id][key] = context

        self[e_type].put_element(e_name, entity)
        return True

    # Delete the checksum right of the entity linked
    # new_checksum ^= checksum
    # entity can be a role, a profile, or a composite
    # return 0 if the new right is deleted or not found
    # else return the new checksum
    def delete_right(self, entity, e_type, right_id, checksum):

        entity = self[e_type + '_storage'].get_elements(ids=entity)

        if (entity['rights']
            and entity['rights'][right_id]
            and entity['rights'][right_id]['checksum'] >= checksum):

            # remove the permissions passed in checksum
            entity['rights'][right_id]['checksum'] ^= checksum

            # If all the permissions were removed from the right, delete it
            if not entity['rights'][right_id]['checksum']:
                del entity['rights'][right_id]
                self[e_type + "_storage"].put_element(entity['_id'], entity)
                return 0

            self[e_type + "_storage"].put_element(entity['_id'], entity)
            return entity['rights'][right_id]['checksum']

        return 0


    # Create a new rights composite composed of the rights passed in comp_rights
    # comp_rights should be a map of rights referenced in the action catalog
    def create_composite(self, comp_name, comp_rights):
        """
        Create a new rights composite
        and add it in the appropriate Mongo collection
        """

        # Do nothing if it already exists
        if self['composite_storage'].get_elements(
            ids=comp_name, query={'type': 'composite'}
            ):
            return False

        new_comp = {'type': 'composite',
                    'rights': {}
                    }

        self.composite_storage.put_element(comp_name, new_comp)

        # Use add_right to check if the action is referenced
        for right_id in comp_rights:
            self.add_right(comp_name,
                           'composite',
                           right_id,
                           comp_rights[right_id]['checksum'])

        return True


    # Create a new profile composed of the composites p_composites
    #   and which name will be p_name
    # If the profile already exists, composites from p_composites
    #   that are not already in the profile's composites will be added
    # Return True if the profile was created, False otherwise
    def create_profile(self, p_name, p_composites):
        """
        Create profile p_name composed of p_composites
        """

        # Do nothing if it already exists
        if self['profile_storage'].get_elements(
            ids=p_name, query={'type': 'profile'}
            ):
            return False

        new_profile = {'type':'profile',
                       'composites': []
                       }

        self.profile_storage.put_element(p_name, new_profile)

        for comp in p_composites:
            self.add_composite(p_name, 'profile', comp)

        return True


    # Delete entity of id e_name
    # Return True if the entity was deleted, False otherwise
    # f_type is the storage to remove it from
    # t_type is the storage to check for relations
    # entity can be a profile, or composite
    def delete_entity(self, e_name, e_type):
        """
        Delete entity
        """

        from_storage = e_type + '_storage'
        t_type = 'profile' if e_type == 'composite' else 'role'
        to_storage = t_type + '_storage'

        if self[from_storage].get_elements(ids=e_name):
            self[from_storage].remove_elements(e_name)

            # remove the entity from every other entities that use it
            for entity in self[to_storage].get_elements(query={'type':t_type}):
                if e_type in entity and e_name in entity[e_type]:
                    entity[e_type].remove(e_name)
                    self[to_storage].put_element(entity['_id'], entity)

            return True

        return False


    # to be removed when user module is created
    def delete_role(self, r_name):
        if self.role_storage.get_elements(ids=r_name):
            self.role_storage.remove_elements(r_name)
            return True

        return False


    # delete_entity wrapper
    def delete_profile(self, p_name):
        self.delete_entity(p_name, 'profile')


    # delete_entity wrapper
    def delete_composite(self, c_name):
        self.delete_entity(c_name, 'composite')


    # Add the composite named comp_name to the entity
    # If the composite does not exist and
    #   comp_rights is specified it will be created first
    # entity can be a profile or a role
    # Return True if the composite was added, False otherwise
    def add_composite(self, e_name, e_type, comp_name, comp_rights=None):
        """
        Add the composite comp_name to the entity
        """

        e_type += '_storage'

        if not self.composite_storage.get_elements(ids=comp_name):
            if comp_rights:
                self.create_composite(comp_rights, comp_name)
            else:
                return False

        entity = self[e_type].get_elements(ids=e_name)
        if not 'composite' in entity:
            entity['composite'] = []
        if not comp_name in entity['composite']:
            entity['composite'].append(comp_name)
            self[e_type].put_element(e_name, entity)

        return True


    # Add the profile of name p_name to the role
    # If the profile does not exists and p_composites is specified
    #    it will be created first
    # Return True if the profile was added, False otherwise
    def add_profile(self, role, p_name, p_composites=None):
        """
        Add profile p_name to role['profile']
        """

        profile = self.profile_storage.get_elements(ids=p_name)
        if not profile:
            if p_composites:
                self.create_profile(p_name, p_composites)
            else:
                return None

        # retrieve the profile
        if profile:
            s_role = self.role_storage.get_elements(ids=role)
            if not 'profile' in s_role or not len(s_role['profile']):
                s_role['profile'] = []

            if not p_name in s_role['profile']:
                s_role['profile'].append(p_name)
                self.role_storage.put_element(role, s_role)

            return p_name


    # Remove the entity e_name from from_name
    # from_name can be a profile or a role
    # e_name can be a profile or a composite
    # Return True if e_name was removed, False otherwise
    def remove_entity(self, from_name, from_type, e_name, e_type):
        entity = self[from_type + '_storage'].get_elements(
            query={'type': from_type}, ids=from_name)

        if e_type in entity and e_name in entity[e_type]:
            entity[e_type].remove(e_name)
            self[from_type + '_storage'].put_element(from_name, entity)
            return True

        return False


    # remove_entity wrapper
    def remove_composite(self, e_name, e_type, comp_name):
        """
        Remove the composite comp_name from the entity
        """

        return self.remove_entity(e_name, e_type, comp_name, 'composite')


    # remove_entity wrapper
    def remove_profile(self, role, p_name):
        """
        Remove profile p_name from the role
        """

        return self.remove_entity(role, 'role', p_name, 'profile')


    # Create a new role composed of the profile r_profile
    #   and which name will be r_name
    # Any extra field can be specified in the kwargs
    # If the role already exists, the profile will be changed for r_profile
    # Return the newly created role's name or False it the creation failed
    def create_role(self, r_name, r_profile):
        """
        Create role p_name composed of p_composites
        """

        if self.role_storage.get_elements(ids=r_name):
            return r_name

        new_role = {'type': 'role'}
        new_role['profile'] = []
        if isinstance(r_profile, list):
            new_role['profile'] = r_profile
        else:
            new_role['profile'].append(r_profile)
        self.role_storage.put_element(r_name, new_role)
        return r_name

    # Generic getter
    def get_from_storage(self, s_name):
        def get_from_storage_(elem, default=None):
            if not elem in self[s_name]:
                return default
            return self[s_name][elem]
        return get_from_storage_

