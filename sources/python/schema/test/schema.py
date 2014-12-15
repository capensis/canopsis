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
from os.path import join, isfile, abspath, dirname
from os import remove
from mock import patch

from canopsis.schema.schema import Schema


def mock_get_schema_path(*args):
    this_file_dir = dirname(abspath(__file__))
    return join(this_file_dir, 'xml/schema', *args)


def get_data_path(*args):
    this_file_dir = dirname(abspath(__file__))
    return join(this_file_dir, 'xml/data', *args)


class TestSchema(TestCase):

    def setUp(self):
        self.s = Schema()

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_get_schema(self, get_schema_path):
        schema_query1 = self.s.get_schema('profile_to_name.xsl')
        schema_query2 = self.s.get_schema('profile>name:1.0')
        schema_query3 = self.s.get_schema('profile>name')

        self.assertIsInstance(schema_query1, str)
        self.assertIsInstance(schema_query2, str)
        self.assertIsInstance(schema_query3, str)

        self.assertEqual(schema_query1, schema_query2)
        self.assertEqual(schema_query2, schema_query3)

        self.assertRaises(ValueError,
                          self.s.get_schema,
                          ('a_wrong_query_token'))

    def test_get_data_type_schemas(self):
        self.assertRaises(
            NotImplementedError,
            self.s.get_data_type_schemas,
            'data_type')

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    @patch('canopsis.schema.schema.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_push_schema(self, schema_get_schema_path, utils_get_schema_path):
        schema = ('<?xml version="1.0" encoding="UTF-8" ?>'
                  '<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"'
                  '           xmlns="new:1.0" '
                  '           targetNamespace="new:1.0"'
                  '           elementFormDefault="qualified">'
                  ' <xs:complexType name="profiletype">'
                  '  <xs:sequence>'
                  '   <xs:element name="name" type="xs:string"/>'
                  '   <xs:element name="age" type="xs:integer"/>'
                  '  </xs:sequence>'
                  ' </xs:complexType>'
                  ' <xs:element name="profile" type="profiletype"/>'
                  '</xs:schema>')

        no_key_schema = (
            '<?xml version="1.0" encoding="UTF-8" ?>'
            '<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema">'
            ' <xs:complexType name="profiletype">'
            '  <xs:sequence>'
            '   <xs:element name="name" type="xs:string"/>'
            '   <xs:element name="age" type="xs:integer"/>'
            '  </xs:sequence>'
            ' </xs:complexType>'
            ' <xs:element name="profile" type="profiletype"/>'
            '</xs:schema>')

        existing_key_schema = (
            '<?xml version="1.0" encoding="UTF-8" ?>'
            '<xs:schema xmlns:xs="http://www.w3.org/2001/XMLSchema"'
            '           xmlns="profile:1.0"'
            '           targetNamespace="profile:1.0"'
            '           elementFormDefault="qualified">'
            ' <xs:complexType name="profiletype">'
            '  <xs:sequence>'
            '   <xs:element name="name" type="xs:string"/>'
            '   <xs:element name="age" type="xs:integer"/>'
            '  </xs:sequence>'
            ' </xs:complexType>'
            ' <xs:element name="profile" type="profiletype"/>'
            '</xs:schema>')

        self.assertRaises(ValueError,
                          self.s.push_schema,
                          'profile.xsd', '<not matter what>')
        self.assertRaises(SyntaxError,
                          self.s.push_schema,
                          'new.xsd', '<syntax error>')
        self.assertRaises(AttributeError,
                          self.s.push_schema,
                          'new.xsd', no_key_schema)
        self.assertRaises(ValueError,
                          self.s.push_schema,
                          'new.xsd', existing_key_schema)

        self.s.push_schema('new.xsd', schema)
        self.assertTrue(isfile(mock_get_schema_path('new.xsd')))

        remove(mock_get_schema_path('new.xsd'))

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_validate_schema(self, get_schema_path):
        self.assertEqual(self.s.validate_schema('profile.xsd'),
                         [True, 'XMLSchema'])
        # keyless schema is not recommended, but it should pass though
        self.assertEqual(self.s.validate_schema('keyless_profile.xsd'),
                         [True, 'XMLSchema'])

        self.assertEqual(self.s.validate_schema('profile_to_name.xsl'),
                         [True, 'XSLT'])

        self.assertEqual(self.s.validate_schema('<not a schema>'),
                         [False, None])
        self.assertEqual(self.s.validate_schema('syntaxerror_schema.xsd'),
                         [False, None])

    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_validate_data(self, get_schema_path):

        with open(get_data_path('profile.xml')) as data_file:
            data = data_file.read()

        self.assertTrue(self.s.validate_data(data, 'profile.xsd'))
        self.assertTrue(self.s.validate_data(data, 'profile:1.0'))

        self.assertTrue(self.s.validate_data(data, 'profile.xsd',
                                             validate_schemas=True))
        self.assertTrue(self.s.validate_data(data, 'profile:1.0',
                                             validate_schemas=True))

        with self.assertRaises(AssertionError):
            self.s.validate_data(data,
                                 'syntaxerror.xsd',
                                 validate_schemas=True)

    def test_transform(self):
        pass

if __name__ == '__main__':
    main()

