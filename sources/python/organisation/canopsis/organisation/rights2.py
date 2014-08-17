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

from canopsis.configuration import conf_paths, add_category
from canopsis.middleware.manager import Manager

CATEGORY = 'RIGHTS'


@conf_paths('organisation/rights.conf')
@add_category(CATEGORY)
class Rights(Manager):

    DATA_SCOPE = 'rights'

    def __init__(
        self, data_scope=DATA_SCOPE,
        *args, **kwargs
    ):

        super(Rights, self).__init__(data_scope=data_scope, *args, **kwargs)

        self.default_profile = 'vizualisation'

    def _configure(self, unified_conf, *args, **kwargs):

        super(Rights, self)._configure(
            unified_conf=unified_conf, *args, **kwargs)

        self.profile_storage = self['profiles_storage']
        self.composite_storage = self['composite_storage']
        self.role_storage = self['role_storage']
        self.default_profile = 'vizualisation'

    # Entity can be a right_composite, a profile or a role since
    # all 3 of them have a rights field
    def check(self, entity, right_id, checksum):
        """
        Check the from the rights of entity for a right of id right_id
        """

        # check the fields
        if not entity or not entity.get('rights', None):
            return False

        found = entity['rights'].get(right_id, None)

        if (found and found.get(right_id, None)
                and found[right_id].get('checksum', 0) & checksum >= checksum):
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

        profile = self.profiles_storage.get_elements(
            ids=[role.get('profile', None)])

        p_composites = profile.get('composites', None)

        composites = [self.composite_storage.get_elements(ids=[x])
                      for x in p_composites]

        # check in the role's comsposite
        if ((role and self.check(role, right_id, checksum)) or
            # check in the profile's composite
            (profile and self.check(profile, right_id, checksum)) or
            # check in the profile's groups composites
            (p_composites and any(self.check(x, right_id, checksum)
                                  for x in composites))):
            return True

        return False

    # Add a right to the entity linked
    # If the right already exists, the checksum will be summed accordingly
    # checksum |= old_checksum
    # entity can be a role, a profile, or a composite
    def add(self, e_name, e_type, right_id, checksum,
            **kwargs):

        entity = None

        if e_type in self:
            entity = self[e_type].get_elements(ids=[e_name])

        if not entity:
            return False

        if not entity.get('rights', None):
            entity['rights'] = {}

        # If it does not exist, create it
        if not self.check(entity, right_id, 0):
            entity['rights'].update({right_id: {'checksum': checksum}})

        else:
            entity['rights'][right_id]['checksum'] |= checksum

        # Add the new desc and/or context
        for key in kwargs:
            if kwargs[key]:
                entity['rights'][right_id][key] = context

        self[e_type].put_element(entity)
        return True

    # Delete the checksum right of the entity linked
    # new_checksum ^= checksum
    # entity can be a role, a profile, or a composite
    # return 0 if the new right is deleted or not found
    # else return the new checksum
    def delete(self, entity, right_id, checksum):

        if (entity['rights']
                and entity['rights'][right_id]
                and entity['rights'][right_id]['checksum'] >= checksum):

            # remove the permissions passed in checksum
            entity['rights'][right_id]['checksum'] ^= checksum

            # If all the permissions were removed from the right, delete it
            if not entity['rights'][right_id]['checksum']:
                del entity['rights'][right_id]
                return 0

            return entity['rights'][right_id]['checksum']

        return 0

    # Create a new rights composite composed of the rights passed in comp_rights
    # comp_rights should be a map of right maps
    def create_composite(self, comp_name, comp_rights):
        """
        Create a new rights composite
        and add it in the appropriate Mongo collection
        """

        new_comp = {}

        for right_id in comp_rights:
            new_comp[right_id] = comp_rights[right_id].copy()

        self.composite_storage.put_element(comp_name, new_comp)

    # Delete composite named comp_name
    # Return True if the composite was successfully deleted
    def delete_composite(self, comp_name):
        """
        Delete the composite comp_name
        and update the appropriate Mongo collection
        """

        if self.composite_storage.get_elements(ids=[comp_name]):

            self.composite_storage.remove_elements(comp_name)
            # remove the composite from every profiles
            for p in self.profile_storage:
                p['composites'].discard(comp_name)
                # if the profile only had that composite, remove the profile
                if not len(p['composites']):
                    self.delete_profile(p)
            return True

        return False

    # Add the composite named comp_name to the entity
    # If the composite does not exist and
    #   comp_rights is specified it will be created first
    # entity can be a profile or a role
    # Return True if the composite was added, False otherwise
    def add_composite(self, entity, comp_name, comp_rights=None):
        """
        Add the composite comp_name to the entity
        """

        if not self.composite.get(comp_name, None):
            if comp_rights:
                self.create_composite(comp_rights, comp_name)
            else:
                return False

        entity['composites'].add(
            self.composite_storage.get_elements(ids=[comp_name]))
        return True

    # Remove the composite named comp_name from the entity
    # entity can be a profile or a role
    # Return True if the composite was removed, False otherwise
    def remove_composite(self, entity, comp_name):
        """
        Remove the composite comp_name from the entity
        """

        return entity['rights'].pop(comp_name, None)

    # Create a new profile composed of the composites p_composites
    #   and which name will be p_name
    # If the profile already exists, composites from p_composites
    #   that are not already in the profile's composites will be added
    # Return True if the profile was created, False otherwise
    def create_profile(self, p_name, p_composites):
        """
        Create profile p_name composed of p_composites
        """

        if self.profile_storage.get_elements(ids=[p_name]):
            return False

        new_profile = {}
        new_profile['composites'] = p_composites
        self.profile_storage.put_element(p_name, new_profile)
        return True

    # Delete profile of name p_name
    # Return True if the prodile was deleted, False otherwise
    def delete_profile(self, p_name):
        """
        Delete profile p_name
        """

        if self.profile_storage.get_elements(ids=[p_name]):
            self.profile_storage.remove_element(p_name)
            # remove the profile from every roles
            for r in self.role_storage:
                r['profile'].discard(comp_name)
                # if the role only had that profile, remove it
                if not len(r['profile']):
                    self.add_profile(r, self.default_profile)
            return True

        return False

    # Remove the profile p_name from the role
    # Return True if it was removed, False otherwise
    def remove_profile(self, role, p_name):
        """
        Remove profile p_name from the role
        """

        if not role.pop(p_name, None):
            return False

        # replace self.default_profile by a set if you want
        #    to enable multi-profiles
        role['profile'] = self.default_profile

        return True

    # Add the profile of name p_name to the role
    # If the profile does not exists and p_composites is specified
    #    it will be created first
    # Return True if the profile was added, False otherwise
    def add_profile(self, role, p_name, p_composites=None):
        """
        Add profile p_name to role['profile']
        """

        if not self.profile_storage.get_elements(ids=[p_name]):
            if p_composites:
                self.create_profile(p_name, p_composites)
            else:
                return False

        # retrieve the profile
        if self.profile_storage.get_elements(ids=[p_name]):
            # change role['profiles'] to a set of strings
            #   if you want to allow several profiles on
            #   the same role
            role['profiles'] = p_name

    # Create a new role composed of the profile r_profile
    #   and which name will be r_name
    # Any extra field can be specified in the kwargs
    # If the role already exists, the profile will be changed for r_profile
    # Return True if the role was created, False otherwise
    def create_role(self, r_name, r_profile):
        """
        Create role p_name composed of p_composites
        """

        if self.role_storage.get_elements(ids=[r_name]):
            return False

        new_role = {}
        new_role['profiles'] = r_profile
        self.role_storage.put_element(r_name, new_role)
        return True

    def delete_role(self, r_name):
        if self.role_storage.get_elements(ids=[r_name]):
            self.role_storage.remove_element(r_name)
            return True

        return False

    # def create_entity(self, entity_name, entity_type, **kwargs):
    #     """
    #     Create a new entity of type entity_type
    #     The fields specified in the kwargs will be added to the entity
    #     """
    #     new_entity = {}
    #     for key in kwargs:
    #         if kwargs[key]:
    #             new_entity[key] = kwargs[key].copy()

    #     if entity_type in self.storages:
    #         self.storages[entity_type].put_element(entity_name, new_entity)
