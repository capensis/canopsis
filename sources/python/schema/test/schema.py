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
from os.path import join, isfile, abspath, dirname
from os import remove
from mock import patch
from lxml.etree import parse, tostring
from StringIO import StringIO

from canopsis.schema.schema import Schema


def mock_get_schema_path(*args):
    this_file_dir = dirname(abspath(__file__))
    return join(this_file_dir, 'xml/schema', *args)


def get_data_path(*args):
    this_file_dir = dirname(abspath(__file__))
    return join(this_file_dir, 'xml/data', *args)


class TestSchema(TestCase):

    @patch('canopsis.schema.schema.get_schema_path',
           side_effect=mock_get_schema_path)
    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def setUp(self, utils_get_schema_path, schema_get_schema_path):
        self.s = Schema()

    def test_load_cache(self):
        # cache is supposed to be loaded in constructor
        self.assertEqual(len(self.s.cache), 3)
        self.assertIn('profile:1.0', self.s.cache)
        self.assertIn('profile:2.0', self.s.cache)
        self.assertIn('profile>name:1.0', self.s.cache)

    def test_cache_schema(self):
        keyless_xschema = parse(mock_get_schema_path('keyless_profile.xsd'))
        self.assertFalse(self.s.cache_schema(keyless_xschema))

        existing_cached_xschema = self.s.cache['profile:1.0']
        existing_xschema = parse(mock_get_schema_path('profile.xsd'))
        self.assertTrue(self.s.cache_schema(existing_xschema))
        self.assertIn('profile:1.0', self.s.cache)
        self.assertEqual(tostring(existing_cached_xschema),
                         tostring(existing_xschema))
        self.assertEqual(self.s.cache['profile:1.0'],
                         existing_xschema)

        new_schema = ('<?xml version="1.0" encoding="UTF-8" ?>'
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
        new_xschema = parse(StringIO(new_schema))
        self.assertTrue(self.s.cache_schema(new_xschema))
        self.assertIn('new:1.0', self.s.cache)
        self.assertEqual(self.s.cache['new:1.0'], new_xschema)

    def test_get_existing_unique_keys(self):
        unique_keys = self.s.get_existing_unique_keys()

        self.assertEqual(len(unique_keys), 3)
        self.assertIn('profile:1.0', unique_keys)
        self.assertIn('profile:2.0', unique_keys)
        self.assertIn('profile>name:1.0', unique_keys)

    def test_is_unique_key_existing(self):
        self.assertTrue(self.s.is_unique_key_existing('profile:1.0'))
        self.assertTrue(self.s.is_unique_key_existing('profile>name:1.0'))
        self.assertFalse(self.s.is_unique_key_existing('profile'))
        self.assertFalse(self.s.is_unique_key_existing('wrong_key'))

    def test_get_cached_schema(self):
        existing_xschema = parse(mock_get_schema_path('profile.xsd'))
        self.assertEqual(tostring(existing_xschema),
                         tostring(self.s.get_cached_schema('profile:1.0')))
        self.assertIs(None, self.s.get_cached_schema('wrong_key'))

    def test_get_schema(self):
        supposed_schema = tostring(parse(
            mock_get_schema_path('profile_to_name.xsl')))
        schema = self.s.get_schema('profile>name:1.0')

        self.assertIsInstance(schema, str)
        self.assertEqual(supposed_schema, schema)

        self.assertIs(None, self.s.get_schema('wrong_key'))

    @patch('canopsis.schema.schema.get_schema_path',
           side_effect=mock_get_schema_path)
    @patch('canopsis.schema.utils.get_schema_path',
           side_effect=mock_get_schema_path)
    def test_push_schema(self, utils_get_schema_path, schema_get_schema_path):
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

    def test_validate_schema(self):
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

        self.assertEqual(self.s.validate_schema('profile:1.0'),
                         [True, 'XMLSchema'])
        self.assertEqual(self.s.validate_schema('profile>name:1.0'),
                         [True, 'XSLT'])
        self.assertEqual(self.s.validate_schema(schema), [True, 'XMLSchema'])
        # well, schema with no key should not be allowed, but it
        # passes the test
        self.assertEqual(self.s.validate_schema(no_key_schema),
                         [True, 'XMLSchema'])

        self.assertEqual(self.s.validate_schema('<not a schema>'),
                         [False, None])
        self.assertEqual(self.s.validate_schema('wrong_key:1.0'),
                         [False, None])

    def test_validate_data(self):

        with open(get_data_path('profile.xml')) as data_file:
            data = data_file.read()

        self.assertTrue(self.s.validate_data(data, 'profile:1.0'))

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

