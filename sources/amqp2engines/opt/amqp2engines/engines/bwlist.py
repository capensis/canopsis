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

from cengine import cengine
from caccount import caccount
from cstorage import get_storage
import cevent
import logging

import time
from datetime import datetime

NAME="bw-list"

DROP = -1


class engine(cengine):

	def pre_run(self):
		account = caccount.caccount(user="root", group="root")
		self.storage = cstorage.get_storage(logging_level=logging.DEBUG, account=account)		
#		self.configuration = self.storage.find({'configuration':'bw-list'}, namespace='object')
		self.configuration = {'rules': [
			{'filter':'connector', 'pattern': 'nagios', 'action': 'pass'},
			{'filter':'connector', 'pattern': 'collectd', 'action': 'drop'},
		]}

	def work(self, event, *xargs, **kwargs):		
		
		if not self.configuration:
			return event

		#When list configuration then check black and white lists
		for filterItem in self.configuration['rules']:
			if self.filterMatch(event, filterItem):
				if filterItem['action'] == 'pass':
					return event
				elif filterItem['action'] == 'drop':
					return DROP
				else:
					self.logger('Unknown action')
		return event

	def filterMatch(self, event, filterItem):
		return event[filterItem['filter']] == filterItem['pattern']
		


