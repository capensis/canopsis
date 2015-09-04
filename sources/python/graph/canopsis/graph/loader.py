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

from canopsis.common.utils import singleton_per_scope, lookup
from canopsis.graph.manager import GraphManager
from canopsis.graph.elements import Vertice


def load(serializedelts, graphmgr=None):
    """Load serialized graphs.

    :param serializedelts: serialized elements to load.
    :type serializedelts: list or dict
    :param GraphManager graphmgr: graph manager to use.
    """

    # ensure graphs is a list of graphs.
    if isinstance(serializedelts, dict):
        serializedelts = [serializedelts]

    if graphmgr is None:
        graphmgr = singleton_per_scope(GraphManager)

    for serializedelt in serializedelts:
        # remove contents from the serialized element
        contents = serializedelt.pop('contents', [])
        # call recursively load on the contents if not empty
        if contents:
            load(contents, graphmgr=graphmgr)

        # finally, create the right element

        # if a dedicated cls exists, use it
        cls = serializedelt.get('_cls')
        if cls is None:  # else get default class
            # if base type is not precised, use vertice
            _type = serializedelt.get('_type', Vertice.BASE_TYPE)
            # get the right class
            cls = lookup(
                "canopsis.graph.elements.{0}".format(_type.capitalize())
            )
        # instantiate it
        elt = cls(**serializedelt)
        # and save it with the manager
        elt.save(manager=graphmgr)
