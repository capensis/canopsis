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
from camqp import camqp
import subprocess
import sys

import os

class cli(ccmd):
	def __init__(self, prompt):
		self.myprompt = prompt + 'rabbitmq'
		ccmd.__init__(self, self.myprompt)
		self.amqp = camqp()
		self.amqp.start()
		self.current_exchange = None

	def do_cd(self, queue):
		print "[TODO] rabbit cd to queue"

	def do_ls(self, queue=None):
		if self.current_exchange is None:
			exchanges_names = self.amqp.exchanges.keys()
			i = 0

			for name in exchanges_names:
				if i > 0:
					sys.stdout.write("\t")

				sys.stdout.write(name + " ")
				i += 1

			sys.stdout.write("\n")
			proc = subprocess.Popen("rabbitmqadmin list queues", shell=True, stdout=subprocess.PIPE)
			stdout_value = proc.communicate()[0]
			print stdout_value

		# else:

	def do_purge(self, queue=None):
		print "[TODO] rabbit purge queue (erase events)"

	def do_top(self, queue=None):
		print "[TODO] global status of queues"

def start_cli(prompt):
	try:
		mycli = cli(prompt)
		mycli.cmdloop()
	except Exception, err:
		print "Impossible to start module: %s" % err
