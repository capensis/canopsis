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

from canopsis.storage.manager import Manager


class Rights(Manager):

    class Error(Exception):
        """
        Raised when a Rights error occured
        """
        pass

    def __init__(self):

        self.right_storage = self.get_storage()
        self.profile_storage = self.get_storage()
        self.entity

    # Entity can be a right_composite, a profile or an user since
    # all 3 of them have a rights field
    def check(self, entity, right_id):
        """
        Check the from the rights of entity for a right of id right_id
        """

        right = self.get(right_id=right_id)

        result = Rights._check_right(entity.rights, right)

        return result

    # Check in the rights of the user (in the role), then the profile's ones,
    # then in the rights_composite
    # Return the value as soon as it's found
    # True if found else False
    def check_user_rights(self, user, right_id):
        """
        Check if user has the right of id right_id
        """

        raise NotImplementedError()

        right = self.get(right_id=right_id)

        result = False

        result = Rights._check_right(user.rights, right)

        if not result:
            role = user.user.role
            if role is not None:

                result = Rights.check_right(role.profile.rights, right)

                if not result:
                    for relationship in role.relationships:
                        rights = role.profile.relationships[relationship]
                        if Rights.check_right(rights, right):
                            result = True
                            break

        return result

    @staticmethod
    def _check_right(rights, right):
        """
        Check if at least one input rights right matches with input tight.
        """

        result = False

        for _right in rights:
            if _right.element_id == right.element_id \
                    and (_right.checksum & right.checksum) == right.checksum:
                result = True
                break

        return result

    def get(self, right_id):
        """
        Return the map of the right of id right_id
        """

        result = self.right_storage.get(data_ids=right_id)

        return result


    # Add a right to the entity linked
    # If the right already exists, the checksum will be summed accordingly
    # checksum |= old_checksum
    # entity can be a role, a profile, or a composite
    def add(self, entity, right_id, checksum,
            context=None, element_id=None, desc=None):

        if element_id is None:
            element_id = ""

        if desc is None:
            desc = ""

        if context is None:
            context = {}

        raise NotImplementedError()

    # Delete the checksum right of the entity linked
    # new_checksum ^= checksum
    # entity can be a role, a profile, or a composite
    def delete(self, entity, right_id, checksum):

        raise NotImplementedError()

    # Create a new rights composite composed of the rights passed in comp_rights
    # Which name will be name
    def create_composite(self, comp_rights, comp_name):
        """
        Create a new rights composite
        and add it in the appropriate Mongo collection
        """

        raise NotImplementedError()

    # Delete composite named comp_name
    def delete_composite(self, comp_name):
        """
        Delete the composite comp_name
        and update the appropriate Mongo collection
        """

        raise NotImplementedError()

    # Add the composite named comp_name to the entity
    # If the composite does not exist and
    #   comp_rights is specified it will be created first
    # entity can be a profile or a user
    def add_composite(self, entity, comp_name, comp_rights={}):
        """
        Add the composite comp_name to the entity
        """

        raise NotImplementedError()

    # Remove the composite named comp_name from the entity
    # entity can be a profile or a user
    def remove_composite(self, entity, comp_name):
        """
        Remove the composite comp_name from the entity
        """

        entity['rights'].pop(comp_name, None)

        self.right_composite.remove(data_id=comp_name)

    # Create a new profile composed of the composites p_composites
    #   and which name will be p_name
    # If the profile already exists, composites from p_composites
    #   that are not already in the profile's composites will be added
    def create_profile(self, p_name, p_composites, relationships):
        """
        Create profile p_name composed of p_composites
        """

        self.profile_storage.put(
            data_id=p_name,
            value={
                'composites': p_composites,
                'relationships': relationships})

    # Delete profile of name p_name
    def delete_profile(self, p_name):
        """
        Delete profile p_name
        """

        self.profile_storage.remove(data_ids=p_name)

    # Add the profile of name p_name to the role
    def add_profile(self, role, p_name):
        """
        Add profile p_name to role['profile']
        """

        # retrieve the profile
        profile = self.profile_storage.get(data_ids=p_name)

        # if profile exists
        if profile is not None:
            role['profile'] = p_name

            # add a map of concrete relationships in the role related to
            # profile relationships
            for name, rights in profile['relationships']:
                # last value contains concrete element_ids
                role['relationships'] = name, rights, None
        else:
            raise Rights.Error('Profile %s does not exist' % profile)
