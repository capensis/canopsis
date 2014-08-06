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


class User(Manager):

    def __init__(
        self,
        _id=None,
        rights=None, role=None, profile=None,  # permission informations
        contact=None, session=None,  # data information
        groups=None
    ):

        self.user_storage = self.get_storage()
        self.group_storage = self.get_typed_storage()

        self.id = _id

        self.rights = rights
        self.role = role

        if profile is not None:
            self.put_profile(profile=profile)

        self.contact = contact
        self.session = session

        self.load()

    def load(self):

        # update saved content of user
        user = self.user_storage(data_id=self.id)
        for key, value in user.iteritems():
            getattr(self, key).update(value)

        self.groups = self.group_storage.get(data_ids=self.groups)

    def save(self):

        self.user_storage.update(data_id=self.id)
        self.group_storage.put(data_ids=(group['id'] for group in self.groups))

    def add_group(self, group):
        pass

    def remove_group(self, group):
        pass

    def add_rights(self, rights):

        self.rights.add(rights)

    def remove_rights(self, rights):

        self.remove(rights)

    def put_profile(self, profile=None, concrete_relationships=None):
        """
        Set this profile in adding related relationships to related role.

        :param profile:
        :param concrete_relationships: relationships useful for the role
        """

        # update the input profile
        if self.role is None or self.role['profile'] != profile:
            self.role = dict()
            self.role['profile'] = profile

        # update role relationships with input concrete_relationships
        self.role['relationships'].update(concrete_relationships)

    def remove_profile(self, profile=None):

        if self.role is not None:
            self.role = None

    def check(self, right):
        """
        Check if this user has granted access to the input right.
        """

        result = False

        result = User.check_right(self.rights, right)

        if not result:
            role = self.user.role
            if role is not None:

                result = User.check_right(role.profile.rights, right)

                if not result:
                    for relationship in role.relationships:
                        rights = role.profile.relationships[relationship]
                        result = User.check_right(rights, right)
                        if result:
                            break

        return result

    @staticmethod
    def check_right(self, rights, right):

        """
        Check if at least one input rights right matches with input tight.
        """

        for self_right in self.rights:
            if self_right.element_id == right.element_id \
                and (self_right.checksum & right.checksum) \
                    == right.checksum:
                result = True
                break

        return result
