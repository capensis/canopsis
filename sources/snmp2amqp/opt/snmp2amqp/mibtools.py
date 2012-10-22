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

import sys, os, logging

mib_path="~/var/snmp/"

severity_to_state = {
	'INFORMATIONAL': 0,
	'MINOR': 1,
	'MAJOR': 2,
	'CRITICAL': 2,
	'FATAL': 2,
	'WARNING': 1,
	'NORMAL': 0,
	'UNKNOWN': 3
}

state_to_state = {
	'FAILED': 2,
	'DEGRADED': 1,
	'NONOPERATIONAL': 2,
	'OPERATIONAL': 0,
	'UNKNOWN': 3
}

class mib(object):
	def __init__(self, name):
		self.name = name
		self.logger = logging.getLogger('mib')

		self.logger.info("Load mib %s ..." % name)

		
		sys.path.append(os.path.expanduser(mib_path+name))
		mod_name = name.replace('-', '_')
		exec("from %s import notifications_oid" % mod_name )
		exec("from %s import notifications" % mod_name )

		self.notifications =     notifications
		self.notifications_oid = notifications_oid

		self.logger.debug(" + Loaded with %s Notifications" % len(self.notifications_oid) )

	def get_notification(self, oid):
		try:
			name = self.notifications_oid[oid]
			notification = self.notifications[name]
			notification['name'] = name
			try:
				notification['SEVERITY']
			except:
				notification['SEVERITY'] = "UNKNOWN"

			try:
				notification['STATE']
			except:
				notification['STATE'] = "UNKNOWN"
				
			try:
				notification['ARGUMENTS']
			except:
				notification['ARGUMENTS'] = []
				
			try:
				notification['SUMMARY']
			except:
				notification['SUMMARY'] = ''

			return notification
		except:
			self.logger.debug("OID is not in this mib")
			return None

		
