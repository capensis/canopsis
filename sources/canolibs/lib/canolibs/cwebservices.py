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

import urllib2, cookielib, json
import logging

logging.basicConfig(level=logging.DEBUG,
                    format='%(asctime)s %(name)s %(levelname)s %(message)s',
                    )

class cwebservices(object):
	def __init__(self, host="127.0.0.1", port=8082, logging_level=logging.DEBUG):

		self.logger = logging.getLogger('cwebservice')
		self.logger.setLevel(logging_level)

		self.logger.debug('Init urlib object ...')
		self.base_url = 'http://' + host + ':' + str(port)

		self.jar = cookielib.CookieJar()
		self.jar.clear_session_cookies()

		self.handler = urllib2.HTTPCookieProcessor(self.jar)
		self.opener = urllib2.build_opener(self.handler)

		urllib2.install_opener(self.opener)

		self.is_login = False

	def get(self, uri, parsing=True):
		url = self.base_url + uri
		self.logger.debug(' + GET '+url)
		data = urllib2.urlopen(url).read()

		if parsing:
			self.logger.debug('   + Try to parse json Data ...')
			try:
				data_json = json.loads(str(data))
				data = data_json['data']
				state = data_json['success']
			except:
				self.logger.debug('     + Failed')
				raise Exception("Failed to parse response ...")

			if not state:
				raise Exception("Request marked failed by server ...")
			
			self.logger.debug('     + Success')
	
		return data

	def login(self, login, password):
		self.logger.debug("Login with '%s'" % login)
		self.get("/auth/" + login + "/" + password, False)

		#print self.get("/online")

		self.is_login = True

	def logout(self):
		self.logger.debug("Logout.")
		self.get("/logout", False)

	def __del__(self):
		if self.is_login:
			self.logout()

