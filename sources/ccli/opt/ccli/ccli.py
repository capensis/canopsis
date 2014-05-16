#!/usr/bin/env python2.7
#--------------------------------
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

import readline
import socket
import imp
import sys
import os

from ccmd import ccmd

modules_path = os.path.expanduser('~/opt/ccli/libexec')
sys.path.append(modules_path)


class Application(object):
	class cli(ccmd):
		def __init__(self, app):
			ccmd.__init__(self, '{0}/'.format(app.prompt))

			self.app = app

	def __init__(self, *args, **kwargs):
		super(Application, self).__init__(*args, **kwargs)

		self.hostname = socket.gethostname()
		self.user = 'root'

		self.prompt = '{0}@{1}:'.format(self.user, self.hostname)

		self.modules = {}

		for module in os.listdir(modules_path):
			if module[0] != '.':
				modname, ext = os.path.splitext(module)

				if ext == '.py':
					abspath = os.path.join(modules_path, module)
					self.modules[modname] = imp.load_source(modname, abspath)

		for modname in self.modules:
			module = self.modules[modname]

			def caller(s, line):
				module.start_cli(s.app.prompt)

			fname = 'do_{0}'.format(modname)
			caller.__name__ = fname
			Application.cli.__dict__[fname] = caller

	def do_help(self):
		print ("Usage: ccli [module]\n")
		print ("Modules:")

		for modname in self.modules:
			print (" - {0}".format(modname))

	def __call__(self):
		print ('Welcome to Canopsis CLI')

		try:
			if len(sys.argv) == 2:
				if sys.argv[1] == 'help':
					self.do_help()

				else:
					module = sys.argv[1]

					self.modules[module].start_cli(self.prompt)

			else:
				Application.cli(self).cmdloop()

		except KeyboardInterrupt:
			print ('Received KeyboardInterrupt, exiting...')
			sys.exit(0)


if __name__ == '__main__':
	app = Application()
	app()