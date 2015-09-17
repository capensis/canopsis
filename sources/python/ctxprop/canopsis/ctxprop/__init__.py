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

"""This project is dedicated to enrich a context with dynamic information such
as perfdata, periodic behavior, etc.

A ctxprop manager is dedicated to get/put/update/delete context information
from a single point (... of failure ?! Stop to criticize my idea :p) in order
to apply a same logic whatever concerns.

Such context information manager implement the CTXInfoFunder interface which
is able to execute the ctxprop manager methods.
"""

__version__ = '0.1'
