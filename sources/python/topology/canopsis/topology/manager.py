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

"""
Description
===========

Functional
----------

A topology manager aims to get topology elements from DB.

Technical
---------

A topology manager inherits from a GraphManager in order to play with topology
elements such as graph elements.
"""

from canopsis.common.init import basestring
from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.graph.manager import GraphManager
from canopsis.graph.element import GraphElement, Vertice
from canopsis.topology.element import TopoVertice

CONF_PATH = 'topology/topology.conf'
CATEGORY = 'TOPOLOGY'


@add_category(CATEGORY)
@conf_paths(CONF_PATH)
class TopologyManager(GraphManager):
    """
    Manage topological graph.
    """

    DATA_SCOPE = 'topology'  #: default data scope

    def __init__(self, data_scope=DATA_SCOPE, *args, **kwargs):

        super(TopologyManager, self).__init__(
            data_scope=data_scope, *args, **kwargs
        )

    def causals(self, vertice, depth=-1, errstate=None):
        """Get error causal vertices related to input vertices.

        Such vertices are the source vertices from input vertice where state
        matches with input state.

        :param Vertice vertice: starting vertice from where find much more
            possible causal vertices.
        :param int depth: maximal iteration of depth search. If negative,
            search without limit.
        :param int errstate: minimal error state which allow recursive depth
            search.
        """

        info = {TopoVertice.STATE: {'$gte': errstate}}

        # init ids
        if isinstance(vertice, basestring):
            ids = [vertice]
        elif isinstance(vertice, dict):
            ids = [vertice[GraphElement.ID]]
        elif isinstance(vertice, Vertice):
            ids = [vertice.id]

        result = self.get_neighbourhood(
            ids=ids, info=info, orientation=GraphElement.SOURCES, depth=depth
        )

        return result

    def consequences(self, vertice, depth=-1, errstate=None):
        """Get error consequence vertices related to input vertices.

        Such vertices are the target vertices from input vertice where state
        matches with input state.

        :param Vertice vertice: starting vertice from where find much more
            possible causal vertices.
        :param int depth: maximal iteration of depth search. If negative,
            search without limit.
        :param int errstate: minimal error state which allow recursive depth
            search.
        """

        info = {TopoVertice.STATE: {'$gte': errstate}}

        # init ids
        if isinstance(vertice, basestring):
            ids = [vertice]
        elif isinstance(vertice, dict):
            ids = [vertice[GraphElement.ID]]
        elif isinstance(vertice, Vertice):
            ids = [vertice.id]

        result = self.get_neighbourhood(
            ids=ids, info=info, orientation=GraphElement.TARGETS, depth=depth
        )

        return result

