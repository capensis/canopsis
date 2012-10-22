#!/usr/bin/env python
# --------------------------------
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

import re, os, sys
import subprocess

from pprint import pprint

def usage(message=None):
	if message:
		sys.stderr.write("%s\n" % message)
	else:
		sys.stderr.write("Usage: %s MIB_DIR > PY_FILE\n" % sys.argv[0])

try:
	mib = sys.argv[1]
except:
	usage()
	sys.exit(1)

#mib_path=os.path.expanduser("~/var/snmp/%s" % mib)
mib_path=os.path.expanduser("%s" % mib)
mib_file = "%s.mib" % mib

if not os.path.exists(mib_path):
	usage("'%s' dir not exist ..." % mib_path)
	sys.exit(1)	

os.chdir(mib_path)

if not os.path.exists(mib_file):
	usage("'%s/%s' file not exist ..." % (mib_path, mib_file))
	sys.exit(1)

## Get all oids
oids = subprocess.check_output("smidump -q -k -l1 -f identifiers -p ../models/RFC1155-SMI.mib -p ../models/RFC* -p models/* %s | tr -s ' '" % mib_file, shell=True)
oids = oids.split("\n")

notifications_oid = {}

for line in oids:
	if line.find('notification') > 0:
		line = line.split(' ')
		notification_name = line[1]
		oid = line[3]
		notifications_oid[oid] = notification_name

## Parse Novell NMS informations
fd = open(mib_file, "r")
notification_bloc=False
notification_name=None

notifications = {}
for line in fd:
	if line.find('TRAP-TYPE') >= 0:
		m = re.search('^\s*(\w*)\s*TRAP-TYPE', line)
		notification_name = m.group(1)
		if notification_name:
			notification_bloc=True
			notifications[notification_name] = {}

	if line.find('::=') > 0 and notification_bloc:
		notification_bloc=False
		notification_name=None
	
	if notification_bloc:
		if line.find('--#') >= 0:
			m = re.search('--#(\w*)\s*"?([a-zA-Z0-9_,\.}{ %:]*)"?\s*', line)
			index = m.group(1)
			value = m.group(2)
			notifications[notification_name][index] = value
fd.close()

sys.stdout.write("notifications_oid = ")
pprint(notifications_oid)

sys.stdout.write("notifications = ")
pprint(notifications)
