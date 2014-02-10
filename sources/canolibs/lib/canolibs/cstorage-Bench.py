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

from cstorage import cstorage
from crecord import crecord
from caccount import caccount
from ctimer import ctimer

import random

def go(account, nb):
	storage.account=account
	## Insert 1000 records
	insert_nb = nb
	timer.start()
	for i in range(0, insert_nb):
		record = crecord({'number': i})
		storage.put(record)
	timer.stop()
	insert_speed = int(insert_nb / timer.elapsed)
	
	## Read all records
	timer.start()
	records = storage.find()
	timer.stop()
	read_nb = len(records)
	read_speed = int(read_nb / timer.elapsed)

	## Update records
	new_records = []
	for record in records:
		record.data = {'check': 'update'}
		new_records.append(record)

	update_nb = len(new_records)
	timer.start()
	records = storage.put(new_records)
	timer.stop()
	update_speed = int(update_nb / timer.elapsed)

	## Remove all records
	timer.start()
	storage.remove(records)
	timer.stop()
	remove_nb = len(records)
	remove_speed = int(remove_nb / timer.elapsed)
	
	print " + Insert Speed:",insert_speed,"records/s (%s records)" % insert_nb
	print " + Read Speed:",read_speed,"records/s (%s records)" % read_nb
	print " + Update Speed:",update_speed,"records/s (%s records)" % update_nb
	print " + Remove Speed:",remove_speed,"records/s (%s records)" % remove_nb



namespace = "bench-"+str(random.randint(0,1000))
account = caccount()
storage = cstorage(account=account, namespace=namespace)
timer = ctimer()

print "Bench with 'anonymous' account ..."
account = caccount()
go(account, 5000)

print "Bench with 'root' account ..."
account = caccount(user="root", group="root")
go(account, 5000)

storage.drop_namespace(namespace)
