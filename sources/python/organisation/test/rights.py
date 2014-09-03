#!/usr/bin/env python
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

from logging import getLogger
from unittest import main, TestCase
from canopsis.organisation.rights import Rights
from pprint import pprint

class RightsTest(TestCase):

    def setUp(self):
        self.logger = getLogger()
        self.rights = Rights()
        self.data_types = ['profile', 'composite', 'role']

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
        self.rights.delete_composite('composite_test1')
        self.rights.delete_composite('composite_test2')

        self.rights.delete_user('jharris')

        self.rights.delete_profile('profile_test1')
        self.rights.delete_profile('profile_test2')

        self.rights.delete_role('role_test1bis')


    def tests(self):
        # Test creation of composites
        rights = {
            '1234': {'desc': 'test right in comp', 'checksum': 15},
            '1235': {'desc': 'test right in comp', 'checksum': 8},
            '1236': {'desc': 'test right in comp', 'checksum': 12},
            '1237': {'desc': 'test right in comp', 'checksum': 1},
            '1238': {'desc': 'test right in comp', 'checksum': 15},
            '1239': {'desc': 'test right in comp', 'checksum': 15},
            '1240': {'desc': 'test right in comp', 'checksum': 8},
            '1241': {'desc': 'test right in comp', 'checksum': 8}
            }
        rights_scnd = {
            '2344': {'desc': 'test right in comp', 'checksum': 15},
            '2345': {'desc': 'test right in comp', 'checksum': 8},
            '2346': {'desc': 'test right in comp', 'checksum': 12},
            '2347': {'desc': 'test right in comp', 'checksum': 1},
            '2348': {'desc': 'test right in comp', 'checksum': 15},
            '2349': {'desc': 'test right in comp', 'checksum': 15},
            '4210': {'desc': 'test right in comp', 'checksum': 8},
            '4211': {'desc': 'test right in comp', 'checksum': 8}
            }

        # basic composite creation
        self.rights.create_composite('composite_test1', rights)
        self.assertEqual(
            self.rights['composite_storage'].get_elements(
                query={'type':'composite'})[0]['_id'],
            'composite_test1'
            )

        self.rights.create_composite('composite_test2', rights_scnd)
        self.assertEqual(
            self.rights['composite_storage'].get_elements(
                query={'type':'composite'})[1]['_id'],
            'composite_test2'
            )

        # basic profile creation
        self.rights.create_profile('profile_test1', ['composite_test1'])
        self.assertEqual(
            self.rights['profile_storage'].get_elements(
                query={'type':'profile'})[0]['_id'],
            'profile_test1'
            )

        self.rights.create_profile('profile_test2', ['composite_test2'])
        self.assertEqual(
            self.rights['profile_storage'].get_elements(
                query={'type':'profile'})[1]['_id'],
            'profile_test2'
            )

        # create new role and assign it to the user
        self.rights.create_role('role_test1bis', 'profile_test1')
        self.assertEqual(
            self.rights['role_storage'].get_elements(
                query={'type':'role'})[0]['_id'],
            'role_test1bis'
            )

        self.rights.create_user('jharris', 'role_test1bis')
        self.assertEqual(
            self.rights['user_storage'].get_elements(
                query={'type':'user'})[0]['_id'],
            'jharris'
            )

        # Basic check
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 1), True)

        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add permissions to the rights
        self.rights.add_right('composite_test1', 'composite', '1237', 12)
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Test right deletion
        self.rights.remove_right('composite_test1', 'composite', '1237', 8)
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add it back
        self.rights.add_right('composite_test1', 'composite', '1237', 12)

        # Test remove_entity
        self.rights.remove_comp_profile('profile_test1', 'composite_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add it back
        self.rights.add_comp_profile('profile_test1', 'composite_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Test removing the profile
        self.rights.remove_profile('role_test1bis', 'profile_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add it back
        self.rights.add_profile('role_test1bis', 'profile_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Remove the profile to add it to the role
        self.rights.remove_profile('role_test1bis', 'profile_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), False)

        # Add the composite to the role
        self.rights.add_comp_role('role_test1bis', 'composite_test1')
        self.assertEqual(
            self.rights.check_rights(
                'jharris', '1237', 12), True)

        # Remove it
        self.rights.remove_comp_role('role_test1bis', 'composite_test1')
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


if __name__ == '__main__':
    main()
