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

import unittest

from cwebservices import cwebservices

WS = None

get_auth_uri = [
	# account
	'/account/me',
	'/account/',
	'/account/account.root',
	# rest
	'/rest/object',
	'/rest/object/account',
	'/rest/object/account/account.root',
]

class KnownValues(unittest.TestCase): 
	def setUp(self):
		pass

	def test_01_Init(self):
		global WS
		WS = cwebservices()

	def test_02_TestURL(self):
		for uri in get_auth_uri:
			success=False
			try:
				WS.get(uri)
				success=True
			except:
				pass
		
			if success:
				raise Exception("%s is open !" % uri)

	def test_03_Login(self):
		WS.login('root', 'root')

	def test_04_TestURLAuthed(self):
		WS.login('root', 'root')

		for uri in get_auth_uri:
			WS.get(uri)

	def test_99_Logout(self):
		WS.logout()
	
		
if __name__ == "__main__":
	unittest.main(verbosity=2)
