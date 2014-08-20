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

from bottle import get, delete, put

from canopsis.topology.manager import TopologyManager

manager = TopologyManager()

ROOT_ROUTE = '/rest/topology'


def _route(function, *parameters):
    """
    Get the right route related to input function
    """

    result = '%s/%s' % (ROOT_ROUTE, function)

    for parameter in parameters:
        result = '%s/:%s' % (result, parameter)

    return result


def apply_rest(op, *parameters, mandatory=True):
    """
    Decorator which apply input op on a function with parameters

    :param mandatory: specify if parameters are optional or mandatory
    """

    def apply_rest_on_function(function):

        function_name = function.__name__

        if mandatory:
            route = _route(function_name, *parameters)
            function = op(route)(function)

        else:
            route = _route(function_name)
            function = op(route)(function)
            for i in range(1, len(parameters)):
                _parameters = parameters[:i]
                route = _route(function_name, *_parameters)
                function = op(route)(function)

        return function

    return apply_rest_on_function


@apply_rest(get, 'ids', 'add_nodes', mandatory=False)
def get(ids=None, add_nodes=True):

    if not ids:
        ids = None

    result = manager.get(ids=ids, add_nodes=add_nodes)

    return result


@apply_rest(get, 'regex', 'add_nodes', mandatory=False)
def find(regex=None, add_nodes=False):

    result = manager.find(regex=regex, add_nodes=add_nodes)

    return result


@apply_rest(put, 'topology')
def put(topology=None):

    manager.put(topology=topology)


@apply_rest(delete, 'ids', mandatory=False)
def remove(ids=None):

    manager.remove(ids=ids)


@apply_rest(get, 'ids', mandatory=False)
def get_nodes(ids=None):

    result = manager.get_nodes(ids=ids)

    return result


@apply_rest(get, 'entity_id')
def find_nodes_by_entity_id(entity_id):

    result = manager.find_nodes_by_entity_id(entity_id=entity_id)

    return result
