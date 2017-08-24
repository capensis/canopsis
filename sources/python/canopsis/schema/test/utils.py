# /usr/bin/env python
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

from unittest import TestCase, main
from mock import patch
from os.path import join, abspath, dirname
from lxml.etree import parse, _ElementTree

from canopsis.schema.utils import get_unique_key, get_xml, is_name_available


def mock_get_schema_path(*args):
    this_file_dir = dirname(abspath(__file__))
    return join(this_file_dir, 'xml/schema', *args)


class TestUtils(TestCase):

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
    def test_get_xml(self, get_schema_path):
        xschema = get_xml('profile>name:1.0')
        wrong_xschema = get_xml('wrong_key')

        self.assertIsInstance(xschema, _ElementTree)
        self.assertIs(wrong_xschema, None)

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_is_name_available(self, get_schema_path):
        self.assertTrue(is_name_available('profile.xsd'))
        self.assertTrue(is_name_available('profile_to_name.xsl'))
        self.assertFalse(is_name_available('a_wrong_query_token'))


if __name__ == '__main__':
    main()

