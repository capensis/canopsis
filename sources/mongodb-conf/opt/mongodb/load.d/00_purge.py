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

from caccount import caccount
from cstorage import cstorage

##set root account
root = caccount(user="root", group="root")
storage = cstorage(account=root)

logger = None


def init():
	namespaces = ['cache', 'events', 'events_log', 'object']

	for namespace in namespaces:
		logger.info(" + Drop '%s' collection" % namespace)
		storage.drop_namespace(namespace)


def update():
	storage.drop_namespace('cache')
