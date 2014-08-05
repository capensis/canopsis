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


class Rights(object):

    def __init__(self):
        pass

    # Entity can be a right_composite, a profile or a role since
    # all 3 of them have a rights field
    def check(self, entity, right_id):
        """
        Check the from the rights of entity for a right of id right_id
        """

        raise NotImplementedError()

    # Check in the rights of the user (in the role), then the profile's ones,
    # then in the rights_composite
    # Return the value as soon as it's found
    # True if found else False
    def check_user_rights(self, user, right_id):
        """
        Check if user has the right of id right_id
        """

        raise NotImplementedError()


    def get(self, right_id):
        """
        Return the map of the right of id right_id
        """

        raise NotImplementedError()


    # Add a right to the entity linked
    # If the right already exists, the checksum will be summed accordingly
    # checksum |= old_checksum
    # entity can be a role, a profile, or a composite
    def add(self, entity, right_id, checksum,
            context=None, element_id="", desc=""):

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
    def add_composite(self, entity, comp_name, comp_rights=None):
        """
        Add the composite comp_name to the entity
        """

        raise NotImplementedError()


    # Remove the composite named comp_name from the entity
    # entity can be a profile or a user
    def remove_composite(self, entity, comp_name);
        """
        Remove the composite comp_name from the entity
        """

        raise NotImplementedError()


    # Create a new profile composed of the composites p_composites
    #   and which name will be p_name
    # If the profile already exists, composites from p_composites
    #   that are not already in the profile's composites will be added
    def create_profile(self, p_name, p_composites):
        """
        Create profile p_name composed of p_composites
        """
        raise NotImplementedError()


    # Delete profile of name p_name
    def delete_profile(self, p_name):
        """
        Delete profile p_name
        """

        raise NotImplementedError()


    # Add the profile of name p_name to the role
    def add_profile(self, role, profile):
        """
        Add profile p_name to role['profile']
        """

        raise NotImplementedError()




