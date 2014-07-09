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

"""
Python utility library.
"""


def resolve_element(path):
    """
    Get element reference from input full path element.

    :limitations: does not resolve class method.

    :param path: full path to a python element.
        Examples:
            - __builtin__.open
            - ccommon.utils.resolve_element
    :type path: str

    :return: python object which is accessible thourgh input path.
    :rtype: object
    """

    components = path.split('.')

    # if mod does not exist
    result = None

    module_name = components[0]

    # try to import the first component name
    try:
        result = __import__(module_name)
    except ImportError:
        pass

    # try to import all sub-modules/packages
    if result is not None:

        try:  # check if name is defined from an external module
            # find the right module

            for index in range(1, len(components)):
                module_name = '{0}.{1}'.format(module_name, components[index])
                result = __import__(module_name)

        except ImportError:
            pass

    # if result exist
    if result is not None:
        # path its content
        for comp in components[1:]:
            result = getattr(result, comp)

    return result


def path(element):
    """
    Get full path of a given element.

    Do the inverse of resolve_element

    :param element: must be directly defined into a module or a package
    :type element: object
    """

    result = '{0}.{1}'.format(element.__module__, element.__name__)

    return result
