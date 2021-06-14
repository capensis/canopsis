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

from ccmd import ccmd
from ccmd import cbrowser
from caccount import caccount

import os

class cli(ccmd):
	def __init__(self, prompt):
		self.myprompt = prompt + 'storage'
		ccmd.__init__(self, self.myprompt)

	def do_cd(self, namespace):
		cbrowser(self.myprompt + '/' +  namespace, caccount(user="root", group="root"), namespace).cmdloop()

	def do_mongo(self, line):
		os.system('mongo canopsis')

def start_cli(prompt):
	try:
		mycli = cli(prompt)
		mycli.cmdloop()
	except Exception, err:
		print "Impossible to start module: %s" % err
