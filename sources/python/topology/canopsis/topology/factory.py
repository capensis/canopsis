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


from canopsis.context.manager import Context
from canopsis.topology.format import formatter

class Factory(object):
    """docstring for Factory"""
    def __init__(self, arg):
        super(Factory, self).__init__()
        self.arg = arg

    def create_topology(self):
        '''
            TODO
        '''
        pass

    def create_component(self):
        '''
            TODO
        '''
        pass

    def create_connections(self, topology):
        '''
            TODO
        '''
        pass

    def get_topo_id(self, top_ctx):
        '''
            TODO
        '''
        pass
