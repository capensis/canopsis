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
from canopsis.graph.elements import GraphElement, Vertice

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
            search. Default is vertice state.
        """

        vid, state = self._get_idstate(vertice=vertice)

        if errstate is None:

            errstate = 1 if state is None else state

        info = {'state': {'$gte': errstate}}

        result = self.get_neighbourhood(
            ids=vid, info=info, orientation=GraphElement.SOURCES, depth=depth
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
            search. Default is vertice state.
        """

        vid, state = self._get_idstate(vertice=vertice)

        if errstate is None:

            errstate = 1 if state is None else state

        info = {'state': {'$gte': errstate}}

        result = self.get_neighbourhood(
            ids=vid, info=info, orientation=GraphElement.TARGETS, depth=depth
        )

        return result

    def _get_idstate(self, vertice):
        """Get vertice id and state.

        :param vertice: vertice from where get state.
        :type vertice: str, dict or Vertice
        :return: vertice id and state.
        :rtype: str, int
        """

        vid, state = None, None

        if isinstance(vertice, basestring):
            vid = vertice
            vertice = self.get_vertices(ids=vid)
            if vertice is not None:
                state = vertice.info.get('state')

        elif isinstance(vertice, dict):
            vid = vertice[GraphElement.ID]
            if GraphElement.INFO in vertice:
                info = vertice[GraphElement.INFO]
                state = info.get('state')

        elif isinstance(vertice, Vertice):
            vid = vertice.id
            state = vertice.info.get('state')

        return vid, state

    def get_causalsconsequencespertopo(self, topoids, errstate=1):
        """Get a set of causals and consequences from input topos.

        Iterate on all input topos, and get vertices in errors. Once they are
        found, use a set of (vertice, consequences) where vertices are final
        causals (lower vertices in the graph of causals), and consequences are
        only last consequences (upper vertices in the graph of consequences).

        :param str(s) topoids: topology ids from where find causals and
            consequences.
        :param int errstate: minimal errstate for finding causals and
            consequences.
        :return: set of (maximal causals, maximal consequences).
        :rtype: dict
        """

        result = {}

        info = {'state': {'$gte': errstate}}

        # get vertices in error related to topoids
        vertices = self.get_vertices(
            graph_ids=topoids, info=info, serialize=False
        )

        # get causals
        for vertice in vertices:
            vstate = vertice['info']['state']
            causalset = self.causals(
                vertice=vertice, errstate=vstate
            )
            # find all vertice causals
            for causalitem in causalset:
                maxdepth = max(causalitem)
                causals = causalset[maxdepth]
                # find all causals
                for causal in causals:
                    consequenceset = self.consequences(
                        vertice=causal, errstate=vstate
                    )
                    # find all consequences
                    maxdepth = max(consequenceset)
                    result[causal] = consequenceset[maxdepth]

        return result

    def get_consequencescausalspertopo(self, topoids, errstate=1):
        """Get a set of consequences and causals from input topos.

        Iterate on all input topos, and get vertices in errors. Once they are
        found, use a set of (vertice, causals) where vertices are final
        consequences (upper vertices in the graph of consequences), and causals
        are only last causals (lower vertices in the graph of causals).

        :param str(s) topoids: topology ids from where find consequences and
            causals.
        :param int errstate: minimal errstate for finding consequences and
            causals.
        :return: set of (maximal consequences, maximal causals).
        :rtype: dict
        """

        result = {}

        info = {'state': {'$gte': errstate}}

        # get vertices in error related to topoids
        vertices = self.get_vertices(
            graph_ids=topoids, info=info, serialize=False
        )

        # get consequences
        for vertice in vertices:
            vstate = vertice['info']['state']
            consequenceset = self.consequences(
                vertice=vertice, errstate=vstate
            )
            # find all vertice consequences
            for consequenceitem in consequenceset:
                maxdepth = max(consequenceitem)
                consequences = consequenceset[maxdepth]
                # find all causals
                for consequence in consequences:
                    causalset = self.causals(
                        vertice=consequence, errstate=vstate
                    )
                    # find all consequences
                    maxdepth = max(causalset)
                    result[consequence] = causalset[maxdepth]

        return result
