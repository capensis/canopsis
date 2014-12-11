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

from StringIO import StringIO
from lxml.etree import parse, tostring, XMLSchema, XSLT

from canopsis.schema.utils import is_name_available, get_unique_key, get_xml, \
    is_unique_key_existing, get_schema_path


class Schema(object):

    def __init__(self, *args, **kwargs):
        super(Schema, self).__init__(*args, **kwargs)
        pass

    def get_schema(self, schema):
        """
        Return an available schema

        :param str schema: unique_key or filename
        :return: schema in xml
        :rtype: str
        :raises ValueError: if schema param is neither a unique_key nor
          a file
        """

        # Xml Schema. It will be set to a filename if a match is found.
        xschema = get_xml(schema)

        # xschema not set <=> no matches found
        if xschema is None:
            raise ValueError(
                '{} neither matches with a schema filename nor with a '
                'schema unique_key (targetNamespace)'.format(schema))

        return tostring(xschema.getroot())

    def get_data_type_schemas(self, data_type):
        raise NotImplementedError

    def push_schema(self, name, schema):
        """
        Store a schema on server

        :param str name: name of schema file
        :param str schema: xml content in a string
        :return: True if insertion has been performed
        :rtype: bool
        :raises ValueError: if schema name already exists
        :raises SyntaxError: if xml cannot be read by lxml
        :raises AttributeError: if unique_key does not exist
        :raises ValueError: if unique_key exists in an other schema
        """

        if is_name_available(name):
            raise ValueError(
                'Provided schema name : \'{}\' already exists'.format(name))

        if self.validate_schema(schema)[0] is False:
            raise SyntaxError('Provided xml can neither be interpreted as '
                              'XMLSchema nor as XSLT')

        xschema = parse(StringIO(schema))
        unique_key = get_unique_key(xschema)
        if unique_key is None:
            raise AttributeError('Unique key (targetNamespace attribute of '
                                 'root element) does not exist')

        if is_unique_key_existing(unique_key):
            raise ValueError('Unique key (\'{}\') already exists'.
                             format(unique_key))

        # if no exception has been raised until here, we can write the
        # schema
        with open(get_schema_path(name), 'w') as schema_file:
            schema_file.write(schema)

        return True

    def validate_schema(self, schema):
        """
        Make sure provided schema's syntax/grammar are correct

        :param str schema: xml (schema itself) or unique_key or
          filename
        :return: [True, <schema_type>] if schema is correct and
          [False, None] otherwise
        :rtype: list

        .. info:: <schema_type> can be either 'XMLSchema' or 'XSLT'
        """

        xschema = get_xml(schema)

        if xschema is None:
            try:
                xschema = parse(StringIO(schema))
            except:
                return [False, None]

        try:
            XMLSchema(xschema)
            return [True, 'XMLSchema']
        except:
            pass

        try:
            XSLT(xschema)
            return [True, 'XSLT']
        except:
            pass

        return [False, None]

    def assert_structural_schema(self, schema):
        """
        Raise an error if a schema is not a structural one

        :param str schema: xml (schema itself) or unique_key or
          filename
        :raises AssertionError: if schema is not a structural one
        """

        assert(self.validate_schema(schema) == [True, 'XMLSchema']), \
            'Schema \'{}\' has not been validated as a structural schema.'

    def assert_metamorphic_schema(self, schema):
        """
        Raise an error if a schema is not a metamorphic one

        :param str schema: xml (schema itself) or unique_key or
          filename
        :raises AssertionError: if schema is not a metamorphic one
        """

        assert(self.validate_schema(schema) == [True, 'XLST']), \
            'Schema \'{}\' has not been validated as a metamorphic schema.'

    def validate_data(self, data, structural_schema, metamorphic_schema=None,
                      validate_schemas=False):
        """
        Ensure that a data structure matches a schema (xml schema)

        :param str data: data to check
        :param str structural_schema: unique_key or filename
        :param str metamorphic_schema: unique_key or filename
        :param bool validate_schemas: Any provided schema will be
          valided before use if set to True. They are used as is
          otherwise. This option is suited for better perfs. Use it at
          your own risk.
        :return: True if data is valid, False otherwise
        :rtype: bool
        """

        if validate_schemas is True:
            self.assert_structural_schema(structural_schema)

            if metamorphic_schema is not None:
                self.assert_metamorphic_schema(metamorphic_schema)

        if metamorphic_schema is not None:
            self.transform(data, metamorphic_schema)

        xmlschema = XMLSchema(get_xml(structural_schema))
        xml = parse(StringIO(data))

        return xmlschema.validate(xml)

    def transform(self,
                  data,
                  metamorphic_schema,
                  from_structure=None,
                  to_structure=None,
                  validate_data=False,
                  validate_schemas=False,
                  ):
        """
        Transform provided data with a metamorphic schema (xslt)

        :param str data: data to transform
        :param str metamorphic_schema: unique_key or filename
        :param from_structure: Structural schema (unique_key or
          filename) to check data. Check is not performed if None.
        :type from_structure: str or None
        :param to_structure: Structural schema (unique_key or filename)
          to check transformed data. Check is not performed if None.
        :type to_structure: str or None
        :param bool validate_data: Structural schemas are used if True,
          they are ignored otherwise. This option is suited for better
          perfs. Use it at your own risk.
        :param bool validate_schemas: Any provided schema will be
          valided before use if set to True. They are used as is
          otherwise. This option is suited for better perfs. Use it at
          your own risk.
        :return: transformed data
        :rtype: str (xml)
        """

        if validate_schemas is True:
            if validate_data is True:

                if from_structure is not None:
                    self.assert_structural_schema(from_structure)

                if to_structure is not None:
                    self.assert_structural_schema(to_structure)

            self.assert_metamorphic_schema(metamorphic_schema)

        if validate_data is True and from_structure is not None:
            if self.validate_data(data, from_structure) is False:
                raise ValueError('Original data has not been validated '
                                 'according to \'{}\''.format(from_structure))

        xslt = XSLT(get_xml(metamorphic_schema))
        xml = parse(StringIO(data))

        transformed_data = xslt(xml)

        if validate_data is True and to_structure is not None:
            if self.validate_data(transformed_data, to_structure) is False:
                raise ValueError('Transformed data has not been validated '
                                 'according to \'{}\''.format(to_structure))

        return transformed_data

