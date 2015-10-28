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

from os import listdir

from os.path import isfile, join

from sys import prefix

from lxml.etree import parse

from canopsis.configuration.model import Parameter
from canopsis.configuration.configurable import Configurable
from canopsis.configuration.configurable.decorator import (
    add_category, conf_paths
)


@conf_paths('schema/schema.conf')
@add_category('SCHEMA', content=Parameter('schema_location'))
class SchemaManager(Configurable):

    def __init__(self, schema_location=None, *args, **kwargs):
        super(SchemaManager, self).__init__(*args, **kwargs)
        self._schema_location = schema_location

    @property
    def schema_location(self):
        return self._schema_location

    @schema_location.setter
    def schema_location(self, value):
        self._schema_location = value

try:
    _schema_manager = SchemaManager()
except IOError:
    _schema_manager = None


def get_schema_path(*args):
    if _schema_manager is not None:
        return join(prefix, _schema_manager.schema_location, *args)
    else:
        return join('.', *args)


def get_unique_key(xschema):
    """
    Standardize the way to determine a schema unique_key

    :param _ElementTree xschema: lxml _ElementTree object
    :return: unique key or None if not found
    :rtype: str or None
    """

    try:
        return xschema.xpath('/*/@targetNamespace')[0]
    except IndexError:
        return None


def get_xml(unique_key):
    """
    Converts a unique_key to an _ElementTree structure of the matching
    schema.

    :param schema: unique_key
    :return: schema or None if identifier did not lead to any schema
    :rtype: _ElementTree or None
    """

    for schema_file in listdir(get_schema_path()):
        try:
            xschema = parse(get_schema_path(schema_file))
        except:
            continue

        if unique_key == get_unique_key(xschema):
            return xschema
    else:
        return None


def is_name_available(name):
    """
    Assert if a filename exists at the schemas location

    :param str name: schema filename
    :return: True is the schema is present, False otherwise
    :rtype: bool
    """
    return isfile(get_schema_path(name))
