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
import json
from crecord import crecord
from caccount import caccount
from cgroup import cgroup

#storage = cstorage.


class KnownValues(unittest.TestCase): 
	def setUp(self):
		self.anonymous_account = caccount()
		self.root_account = caccount(user="root", group="root")
		self.user_account = caccount(user="william", group="capensis")

		self.data = {'mydata1': 'data1', 'mydata2': 'data2', 'mydata3': 'data3'}

	def test_01_Init(self):
		record = crecord(self.data)
		if record.data != self.data:
			raise Exception('Data corruption ...')

	def test_02_InitFromRaw(self):
		raw = {'_id': None, 'parent': [], 'children': [], 'crecord_name': 'titi', 'aaa_access_group': ['r'], 'aaa_access_owner': ['r', 'w'], 'aaa_group': None, 'aaa_access_unauth': [], 'aaa_owner': None, 'aaa_access_other': [], 'mydata1': 'data1', 'mydata3': 'data3', 'mydata2': 'data2', 'crecord_type': 'raw', 'crecord_write_time': None, 'enable': True}

		record = crecord(raw_record=raw)

		dump = record.dump()
		print(' + _id: %s (%s)' % (dump['_id'], type(dump['_id'])))

		if not isinstance(dump['_id'], type(None)):
			raise Exception('Invalid _id type')

		del record.data['crecord_creation_time']

		if record.data != self.data:
			raise Exception('Data corruption ...')

	def test_03_InitFromRecord(self):
		record = crecord(self.data)

		record2 = crecord(record=record)
		
		del record2.data['crecord_creation_time']
		
		if record2.data != self.data:
			raise Exception('Data corruption ...')

	def test_04_ChOwnGrp(self):
		record = crecord(self.data)

		record.chown('toto')
		if record.owner != 'account.toto':
			raise Exception('chown dont work ...')

		record.chgrp('tata')
		if record.group != 'group.tata':
			raise Exception('chgrp dont work ...')

		#record.chown(self.user_account)
		#if record.owner != 'william' and record.group != 'capensis':
		#	raise Exception('chown with caccount dont work ...')
		
	def test_05_Chmod(self):
		record = crecord({'check': 'bidon'})

		record.chmod('u-w')
		record.chmod('u-r')
		record.chmod('u+w')

		if record.access_owner != ['w']:
			raise Exception('Chmod not work on "owner" ...')
		
		record.chmod('g-w')
		record.chmod('g-r')
		record.chmod('g+w')

		if record.access_group != ['w']:
			raise Exception('Chmod not work on "group" ...')

	def test_06_children(self):
		record1 = crecord(self.data)
		record2 = crecord(self.data)
		record3 = crecord(self.data)

		record1._id = 1
		record2._id = 2
		record3._id = 3
		
		record1.add_children(record2)
		record1.add_children(record3)

		if not record1.is_parent(record2):
			raise Exception('Invalid children association ...')
		if not record1.is_parent(record3):
			raise Exception('Invalid children association ...')

		record1.remove_children(record3)
			
		if record1.is_parent(record3):
			raise Exception('Invalid children supression ...')

	def test_07_enable(self):
		record = crecord(self.data)

		record.set_enable()
		if not record.is_enable():
			raise Exception('Impossible to enable ...')

		record.set_disable()
		if record.is_enable():
			raise Exception('Impossible to disable ...')
			
	def test_08_recursive_dump(self):
		record1 = crecord(self.data)
		record2 = crecord(self.data)
		record3 = crecord(self.data)
		record4 = crecord(self.data)
		
		record2.children.append(record3)
		
		record1.children.append(record2)
		record1.children.append(record4)

		json_output = record1.recursive_dump(json=True)
		json.dumps(json_output)
		
	def test_09_check_admin_rights(self):
		account = caccount(user='jean')
		group = cgroup(name='administrator')
		group.add_accounts(account)
		
		record = crecord(admin_group=group._id,group='nothing',owner='refrigerator')
		
		check = record.check_write(account)
		
		if not check:
			raise Exception('Admin group are not handle ...')
		

if __name__ == "__main__":
	unittest.main(verbosity=2)
	


