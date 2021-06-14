#!/usr/bin/env python
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

from logging import getLogger
from unittest import main, TestCase
from canopsis.organisation.rights import Rights


class RightsTest(TestCase):

    def setUp(self):
        self.logger = getLogger()
        self.rights = Rights()
        self.data_types = ['profile', 'group', 'role']

        # This should be in a filldb script
        referenced_rights = {
            '1234',
            '1235',
            '1236',
            '1237',
            '1238',
            '1239',
            '1240',
            '1241',
            '2344',
            '2345',
            '2346',
            '2347',
            '2348',
            '2349',
            '4210',
            '4211'
        }
        for x in referenced_rights:
            self.rights.add(x, 'desc test rule')

        # delete everything before starting
        self.rights.delete_group('group_test1')
        self.rights.delete_group('group_test2')

        self.rights.delete_user('jharris')

        self.rights.delete_profile('profile_test1')
        self.rights.delete_profile('profile_test2')

        self.rights.delete_role('role_test1bis')

    def tests(self):
        # Test creation of groups
        rights = {
            '1234': {'desc': 'test right in group', 'checksum': 15},
            '1235': {'desc': 'test right in group', 'checksum': 8},
            '1236': {'desc': 'test right in group', 'checksum': 12},
            '1237': {'desc': 'test right in group', 'checksum': 1},
            '1238': {'desc': 'test right in group', 'checksum': 15},
            '1239': {'desc': 'test right in group', 'checksum': 15},
            '1240': {'desc': 'test right in group', 'checksum': 8},
            '1241': {'desc': 'test right in group', 'checksum': 8}
        }
        rights_scnd = {
            '2344': {'desc': 'test right in group', 'checksum': 15},
            '2345': {'desc': 'test right in group', 'checksum': 8},
            '2346': {'desc': 'test right in group', 'checksum': 12},
            '2347': {'desc': 'test right in group', 'checksum': 1},
            '2348': {'desc': 'test right in group', 'checksum': 15},
            '2349': {'desc': 'test right in group', 'checksum': 15},
            '4210': {'desc': 'test right in group', 'checksum': 8},
            '4211': {'desc': 'test right in group', 'checksum': 8}
        }

        # basic group creation
        self.rights.create_group('group_test1', rights)
        self.rights.create_group('group_test2', rights_scnd)

        # basic profile creation
        self.rights.create_profile('profile_test1', ['group_test1'])
        self.rights.create_profile('profile_test2', ['group_test2'])

        # create new role and assign it to the user
        self.rights.create_role('role_test1bis', 'profile_test1')
        self.rights.create_user('jharris', 'role_test1bis')

        # Basic check
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 1), True)

        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add permissions to the rights
        self.rights.add_right('group_test1', 'group', '1237', 12)
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Test right deletion
        self.rights.remove_right('group_test1', 'group', '1237', 8)
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add it back
        self.rights.add_right('group_test1', 'group', '1237', 12)

        # Test remove_entity
        self.rights.remove_group_profile('profile_test1', 'group_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add it back
        self.rights.add_group_profile('profile_test1', 'group_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Test removing the profile
        self.rights.remove_profile('role_test1bis', None, 'profile_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add it back
        self.rights.add_profile('role_test1bis', None, 'profile_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Remove the profile to add it to the role
        self.rights.remove_profile('role_test1bis', None, 'profile_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add the group to the role
        self.rights.add_group_role('role_test1bis', 'group_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Remove it
        self.rights.remove_group_role('role_test1bis', 'group_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add the specific right to the user
        self.rights.add_right('jharris', 'user', '1237', 12)
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Change the checksum
        self.rights.remove_right('jharris', 'user', '1237', 4)
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 8), True)

        # Add a right on the same action to different fields
        # and check that it is summed correctly
        self.rights.remove_right('jharris', 'user', '1237', 8)
        self.rights.remove_right('group_test1', 'user', '1237', 12)
        self.rights.add_right('jharris', 'user', '1237', 2)
        self.rights.add_right('role_test1bis', 'user', '1237', 4)
        self.rights.add_right('group_test2', 'user', '1237', 8)
        self.rights.add_group_user('jharris', 'group_test2')
        self.rights.add_group_role('role_test1bis', 'group_test1')
        self.rights.add_profile('role_test1bis', None, 'profile_test1')
        self.rights.add_group_profile('role_test1bis', 'group_test1')
        ##TODO4-01-2017
		#self.assertEqual(
        #    self.rights.get_user_rights('jharris')['1237']['checksum'],
        #    15)

        # Change entity name
        self.assertTrue('group_test2' in self.rights.get_user('jharris')['group'])
        self.assertTrue(
            'group_test2' == self.rights.get_group('group_test2')['crecord_name']
            )
        self.assertTrue(
            self.rights.update_entity_name('group_test2', 'group', 'name_changed')
            )
        self.assertTrue('group_test2' in self.rights.get_user('jharris')['group'])
        self.assertTrue(
            'name_changed' == self.rights.get_group('group_test2')['crecord_name']
            )

if __name__ == '__main__':
    main()
