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

import readline, os, sys
from ccmd import ccmd
import socket

## Config

hostname = socket.gethostname()
user = "root"
welcome = "Welcome to Canopsis CLI, have fun !"
prompt = user+'@'+hostname+":"

libexec_dir=os.path.expanduser("~/opt/ccli/libexec")
sys.path.append(libexec_dir)

## Import libexec
libexec_list = []
for libexec in os.listdir(libexec_dir):
	if libexec[0] != "." :
		ext = libexec.split(".")[1]
		module = libexec.split(".")[0]
		if ext == "py":
			try:
				exec "import %s" % module
				libexec_list.append(module)
			except Exception, err:
				print("Impossible to load %s:\n\t%s" % (module, err))


## Main object
class cli(ccmd):
	def __init__(self):
		ccmd.__init__(self, prompt + '/')

	for libexec in libexec_list:
		exec "def do_%s(self, line): %s.start_cli(prompt)" % (libexec, libexec)



## Launch CLI
try:
	print welcome

	try:
		if sys.argv[1]:
			try:
				exec "%s.start_cli(prompt)" % sys.argv[1]
			except KeyboardInterrupt:
				print
				sys.exit(1)
			except Exception, err:
				print err
	except:
		pass
	try:
		cli().cmdloop()
	except KeyboardInterrupt:
		print
		sys.exit(1)
	except Exception, err:
		print err
except:
	pass
