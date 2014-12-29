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

"""
Description
===========

Functional
----------

A topology manager aims to get topology elements from DB.

Technical
---------

Three types of nodes exist in topology::

- cluster: operation between node states.
- node: node bound to an entity, like components, resources, and the root.
- root: root node.

A topology data contains several thinks among::

- weight: node weight in the graph.
- state: node state which change at runtime depending on bound entity state and
    event propagation.
"""

from canopsis.configuration.configurable.decorator import (
    conf_paths, add_category
)
from canopsis.graph.manager import GraphManager

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
