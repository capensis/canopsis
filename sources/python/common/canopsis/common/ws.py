#!/usr/bin/env python
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

from inspect import getargspec

from canopsis.common.utils import force_iterable


def response(data):
    """
    Construct a REST response from input data.
    """

    result_data = force_iterable(data)

    result = {
        'total': 0 if data is None else len(result_data),
        'data': result_data,
        'success': True
    }

    return result


def route_name(operation_name, *parameters):
    """
    Get the right route related to input operation_name
    """

    result = '/%s' % operation_name.replace('_', '-')

    for parameter in parameters:
        result = '%s/:%s' % (result, parameter)

    return result


def route(op, name=None):
    """
    Decorator which add ws routes to a callable object.

    :param str name: specify ws name if different than function name

    Example::

        @route(get)
        def entities(a, b, c=None, d=None):
            ...
        fill a, b, c, d parameters in entities function and provide
        the three urls:

            - '/entities/a/b'
            - '/entities/a/b/c'
            - '/entities/a/b/c/d'
    """

    def apply_route_on_function(function):

        function_name = function.__name__ if name is None else name

        argspec = getargspec(function)

        args, defaults = argspec.args, argspec.defaults
        len_args = len(args)
        len_defaults = 0 if defaults is None else len(defaults)

        # add routes with optional parameters
        for n in range(len_defaults):
            route = route_name(function_name, *args[:len_args - n])
            function = op(route)(function)

        # add route with mandatory parameters
        if len_args >= len_defaults:
            route = route_name(function_name, *args[:len_args - len_defaults])
            function = op(route)(function)

        return function

    return apply_route_on_function
