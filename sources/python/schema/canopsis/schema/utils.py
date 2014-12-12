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

from os import listdir
from os.path import isfile
from lxml.etree import parse

from canopsis.configuration.parameters import Parameter
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import (
    add_category, conf_path
)


@conf_path('schema/schema.conf')
@add_category('SCHEMA', content=Parameter('schema_location'))
class SchemaManager(Configurable):

    @property
    def schema_location(self):
        return self._schema_location

    @schema_location.setter
    def schema_location(self, value):
        self._schema_location = value


_schema_manager = SchemaManager()


def get_schema_path(*args):
    return _schema_manager.schema_location
    #return join(prefix, 'etc/schema.d/xml', *args)


def get_unique_key(schema):
    """
    Standardize the way to determine a schema unique_key

    :param _ElementTree schema: lxml _ElementTree object
    :return: unique key or None if not found
    :rtype: str or None
    """

    try:
        return schema.xpath('/*/@targetNamespace')[0]
    except IndexError:
        return None


def get_existing_unique_keys():
    """
    Return all unique keys from all available schemas for the API

    :return: {unique_key: schema_filename+}
    :rtype: dict
    """

    unique_keys = {}

    # At the moment only one location is allowed for schemas : in
    # sys.prefix/etc/schema.d/xml
    for schema_file in listdir(get_schema_path()):
        try:
            schema = parse(get_schema_path(schema_file))
        except:
            # If we don't manage to parse it, it may be a non-xml
            # file that we ignore
            continue

        unique_key = get_unique_key(schema)
        # unique_key is None when an xml parsable ressource does not
        # fit the unique_key requirement
        if unique_key is not None:
            unique_keys[unique_key] = schema_file

    return unique_keys


def get_xml(schema):
    """
    Convert a schema identifier to a schema itself. A schema can be
    identified by :

       - its filename
       - its entire unique_key
       - its unique_key without version part

    .. info:: unique_keys have the following expression :
       <key>:<version>

    :param schema: unique_key or filename
    :return: schema or None if identifier did not lead to any schema
    :rtype: _ElementTree or None
    """

    # First, we check for the filename
    if is_name_available(schema):
        return get_xml_from_name(schema)

    # If no matches are found, we check for the unique_key, both with
    # and without version. A full match (with the version) is
    # prioritary on all partial matches.
    unique_keys = get_existing_unique_keys()
    match_without_version = None

    # The following algorithm anticipates the fact that the research
    # will not perform a full match. match_without_version records the
    # first partial match. In worst case we do not need to iterate over
    # unique_keys a second time.
    for key, schema_file in unique_keys.items():
        if key == schema:
            return get_xml_from_name(schema_file)
        if (key.split(':')[0] == schema and match_without_version is None):
            match_without_version = schema_file
    else:
        if match_without_version is not None:
            return get_xml_from_name(match_without_version)

    # Eventually, if nothing has been returned yet...
    return None


def get_xml_from_name(name):
    """
    Return an xml tree from a filename that can be used by lxml

    :param str name: filename
    :return: xml element tree
    :rtype: _ElementTree (lxml)
    """

    return parse(get_schema_path(name))


def is_name_available(name):
    """
    Assert if a filename exists at the schemas location

    :param str name: schema filename
    :return: True is the schema is present, False otherwise
    :rtype: bool
    """

    return isfile(get_schema_path(name))


def is_unique_key_existing(unique_key):
    """
    Assert if a unique_key already exists

    :param str unique_key: unique key token
    :return: True if a match is found, False otherwise
    :rtype: bool
    """

    unique_keys = get_existing_unique_keys()

    # Processing this function can be optimized depending on how it is
    # used. We assume here that calls will be made by giving entire
    # unique_keys (with version).
    uk_noversion = [key.split(':')[0] for key, schema in unique_keys.items()]

    return unique_key in (list(unique_keys.keys()) + uk_noversion)
