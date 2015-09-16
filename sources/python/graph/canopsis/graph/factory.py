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

"""Graph loader module.
"""

from canopsis.common.utils import singleton_per_scope
from canopsis.graph.manager import GraphManager


class GraphFactory(object):
    """In charge of creating graphs from a data format.

    The data format is a dictionary which contains:

    - uid: elt id.

    and other elements which are defined in the graph.elements module.
    """

    def load(self, serializedelts, graphmgr=None):
        """Load serialized graphs.

        If serialized elts correspond to existing graph elements, the graph
        element is updated with serialized information.

        :param dict(s) serializedelts: serialized elements to load.
        :param GraphManager graphmgr: graph manager to use.
        :return: list of loaded GraphElements.
        :rtype:
        """

        # ensure graphs is a list of graphs.
        if isinstance(serializedelts, dict):
            serializedelts = [serializedelts]

        if graphmgr is None:
            graphmgr = singleton_per_scope(GraphManager)

        result = graphmgr.put_elts(serializedelts)

        return result


class GraphParser(object):
    """In charge of parsing a graph data format to the graph factory expected
    format.
    """

    def parse(self, data):
        """Parse input data and return graph element serialization format.

        :param data: data to parse.
        :return: serialized format.
        :rtype: dict(s)
        """

        raise NotImplementedError()
