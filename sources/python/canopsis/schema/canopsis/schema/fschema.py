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

from StringIO import StringIO
from os import listdir
from lxml.etree import parse, tostring, XMLSchema, XSLT, XML, XMLParser

from canopsis.schema.utils import get_unique_key, get_schema_path, \
    is_name_available
from canopsis.schema.xslt2py import xslt2py


class FSchema(object):

    def __init__(self, *args, **kwargs):
        super(FSchema, self).__init__(*args, **kwargs)

        #: Store both metamorphic and structural schemas in cache. They
        #: can be retrieved by their unique_key.
        self.cache = {}
        self.load_cache()

        #: Factorisable schemas (xslt) are stored in a special cache,
        #: which is first checked
        self.factorisation_cache = {}
        self.load_factorisation_cache()

    def load_cache(self):
        """
        Put each available schema in cache
        """

        # For the moment only one location is allowed for schemas : in
        # sys.prefix/share/canopsis/schema/xml
        for schema_file in listdir(get_schema_path()):
            try:
                xschema = parse(get_schema_path(schema_file))
            except:
                # If we don't manage to parse it, it may be a non-xml
                # file that we ignore
                continue

            unique_key = get_unique_key(xschema)
            # unique_key is None when an xml parsable ressource does not
            # fit the unique_key requirement
            if unique_key is not None and unique_key not in self.cache:
                self.cache[unique_key] = xschema

    def load_factorisation_cache(self):
        """
        Put factorisable schemas in cache
        """
        cache_temp = self.cache.copy()

        for uk, xslt in cache_temp.items():
            if self.validate_schema(uk) == [True, 'XSLT']:
                # Factorisable schemas are cached only if they are
                # factorisable. If it's not, it's still in the standard
                # cache.
                try:
                    parser = XMLParser(target=xslt2py())
                    factorised_schema = XML(tostring(xslt.getroot()), parser)
                    self.factorisation_cache[uk] = factorised_schema
                except:
                    pass

    def cache_schema(self, xschema):
        """
        (Over)write a cache entry for supplied schema

        :param _ElementTree xschema: schema to cache
        :return: False is xschema has no unique_key, True otherwise
        """
        unique_key = get_unique_key(xschema)
        if unique_key is None:
            return False
        else:
            self.cache[unique_key] = xschema
            return True

    def get_existing_unique_keys(self):
        """
        Return all unique keys from all available schemas regarding to
        the API

        :return: {unique_key: xschema+}
        :rtype: dict
        """
        return self.cache.keys()

    def is_unique_key_existing(self, unique_key):
        """
        Assert if a unique_key already exists

        :param str unique_key: unique key token
        :return: True if a match is found, False otherwise
        :rtype: bool
        """
        return unique_key in self.get_existing_unique_keys()

    def get_cached_schema(self, unique_key):
        """
        Checks if schema is cached and returns it

        :param str unique_key: schema unique_key
        :return: schema
        :rtype: _ElementTree or None
        """

        if self.is_unique_key_existing(unique_key):
            return self.cache[unique_key]
        else:
            return None

    def get_schema(self, unique_key):
        """
        Return an available schema

        :param str unique_key: schema unique_key
        :return: schema in xml or None
        :rtype: str or None
        """

        xschema = self.get_cached_schema(unique_key)
        if xschema is not None:
            return tostring(xschema.getroot())
        else:
            return None

    def push_schema(self, name, schema):
        """
        Store a schema on server and cache it

        :param str name: name of schema file
        :param str schema: xml content in a string
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

        if self.is_unique_key_existing(unique_key):
            raise ValueError('Unique key (\'{}\') already exists'.
                             format(unique_key))

        # if no exception has been raised until here, we can write the
        # schema and add it to cache
        self.cache_schema(xschema)
        with open(get_schema_path(name), 'w') as schema_file:
            schema_file.write(schema)

    def validate_schema(self, schema):
        """
        Make sure provided schema's syntax/grammar are correct

        :param str schema: xml (schema itself) or unique_key
        :return: [True, <schema_type>] if schema is correct and
          [False, None] otherwise
        :rtype: list

        .. note:: <schema_type> can either be 'XMLSchema' or 'XSLT'
        """

        if schema in self.get_existing_unique_keys():
            xschema = self.get_cached_schema(schema)
        else:
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

        :param str schema: xml (schema itself) or unique_key
        :raises AssertionError: if schema is not a structural one
        """

        assert(self.validate_schema(schema) == [True, 'XMLSchema']), \
            'Schema \'{}\' has not been validated as a structural schema.'. \
            format(schema)

    def assert_metamorphic_schema(self, schema):
        """
        Raise an error if a schema is not a metamorphic one

        :param str schema: xml (schema itself) or unique_key
        :raises AssertionError: if schema is not a metamorphic one
        """

        assert(self.validate_schema(schema) == [True, 'XSLT']), \
            'Schema \'{}\' has not been validated as a metamorphic schema.'. \
            format(schema)

    def validate_data(self, data, structural_schema, metamorphic_schema=None,
                      validate_schemas=False):
        """
        Ensure that a data structure matches a schema (xml schema)

        :param str data: data to check
        :param str structural_schema: unique_key
        :param str metamorphic_schema: unique_key
        :param bool validate_schemas: Any provided schema will be
          valided before use if set to True. They are used as is
          otherwise. This option is suited for better perfs. Use it at
          your own risk.
        :return: True if data is valid, False otherwise
        :rtype: bool
        :raises
        """

        if validate_schemas is True:
            self.assert_structural_schema(structural_schema)

            if metamorphic_schema is not None:
                self.assert_metamorphic_schema(metamorphic_schema)

        if metamorphic_schema is not None:
            self.transform(data, metamorphic_schema)

        xsl_xml = self.get_cached_schema(structural_schema)
        xmlschema = XMLSchema(xsl_xml)

        try:
            xml = parse(StringIO(data))
        except:
            raise SyntaxError('Wrong xml syntax : {}'.format(data))

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
        :param str metamorphic_schema: unique_key
        :param from_structure: Structural schema unique_key to check
          data. Check is not performed if None.
        :type from_structure: str or None
        :param to_structure: Structural schema unique_key
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

        if metamorphic_schema in self.factorisation_cache:
            xin = parse(StringIO(data))
            factorisation_namespace = {'xin': xin}

            factorised_py = self.factorisation_cache[metamorphic_schema]

            exec(factorised_py, factorisation_namespace)

            transformed_data = factorisation_namespace['xout']
        else:
            xsl_xml = self.get_cached_schema(metamorphic_schema)
            xslt = XSLT(xsl_xml)

            xml = parse(StringIO(data))

            transformed_data = tostring(xslt(xml).getroot())

        if validate_data is True and to_structure is not None:
            if self.validate_data(transformed_data, to_structure) is False:
                raise ValueError('Transformed data has not been validated '
                                 'according to \'{}\''.format(to_structure))

        return transformed_data

