#!/usr/bin/env python
#--------------------------------
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
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
from pyperfstore.node import make_metric_id

from pyperfstore.node import node
from pyperfstore.metric import metric
from pyperfstore.dca import dca

# Stores
from pyperfstore.filestore import filestore
from pyperfstore.memstore import memstore
from pyperfstore.mongostore import mongostore


# Common functions

def node_exist(store, node_id):
	raw = store.get_raw(node_id)
	if raw:
		return True
	else:
		return False

