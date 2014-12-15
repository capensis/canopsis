# /usr/bin/env python
# -*- coding: utf-8 -*-
# --------------------------------
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

from unittest import TestCase, main
from mock import patch
from os.path import join, abspath, dirname
from lxml.etree import parse, tostring, _ElementTree

from canopsis.schema.utils import get_unique_key, get_existing_unique_keys, \
    get_xml, get_xml_from_name, is_name_available, is_unique_key_existing


def mock_get_schema_path(*args):
    this_file_dir = dirname(abspath(__file__))
    return join(this_file_dir, 'xml/schema', *args)


class TestUtils(TestCase):

    def setUp(self):
        pass

    def test_get_unique_key(self):
        supposed_schemas = [
            ('profile.xsd', 'profile:1.0'),
            ('profile_to_name.xsl', 'profile>name:1.0'),
            ('keyless_profile.xsd', None),
        ]

        for schema_file, supposed_key in supposed_schemas:
            xschema = parse(mock_get_schema_path(schema_file))
            returned_key = get_unique_key(xschema)

            self.assertEqual(supposed_key, returned_key)

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_get_existing_unique_keys(self, get_schema_path):
        supposed_keys = {
            'profile:1.0': 'profile.xsd',
            'profile:2.0': 'profile2.xsd',
            'profile>name:1.0': 'profile_to_name.xsl',
        }

        returned_keys = get_existing_unique_keys()

        self.assertEqual(supposed_keys, returned_keys)

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_get_xml(self, get_schema_path):
        xschema_query1 = get_xml('profile_to_name.xsl')
        xschema_query2 = get_xml('profile>name:1.0')
        xschema_query3 = get_xml('profile>name')
        xschema_wrong_query = get_xml('a_wrong_query_token')

        self.assertIsInstance(xschema_query1, _ElementTree)
        self.assertIsInstance(xschema_query2, _ElementTree)
        self.assertIsInstance(xschema_query3, _ElementTree)
        self.assertIs(xschema_wrong_query, None)

        self.assertEqual(tostring(xschema_query1), tostring(xschema_query2))
        self.assertEqual(tostring(xschema_query2), tostring(xschema_query3))

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_get_xml_from_name(self, get_schema_path):
        xschema_query = get_xml_from_name('profile_to_name.xsl')
        self.assertIsInstance(xschema_query, _ElementTree)

        self.assertRaises(IOError, get_xml_from_name, ('profile>name:1.0'))
        self.assertRaises(IOError, get_xml_from_name, ('a_wrong_query_token'))

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_is_name_available(self, get_schema_path):
        self.assertTrue(is_name_available('profile.xsd'))
        self.assertTrue(is_name_available('profile_to_name.xsl'))
        self.assertFalse(is_name_available('a_wrong_query_token'))

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_is_unique_key_existing(self, get_schema_path):
        for key in get_existing_unique_keys():
            self.assertTrue(is_unique_key_existing(key))

        self.assertFalse(is_unique_key_existing('profile.xsd'))
        self.assertFalse(is_unique_key_existing('a_wrong_key'))


if __name__ == '__main__':
    main()

