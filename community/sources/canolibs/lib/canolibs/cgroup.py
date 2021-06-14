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

#import logging
from crecord import crecord
try:
	from caccount import caccount
except:
	pass

class cgroup(crecord):
	def __init__(self, record=None, account_ids=[], description=None, *args, **kargs):
		crecord.__init__(self, *args, **kargs)
		self.type = 'group'
		self._id = '%s.%s' % (self.type,str(self.name))
		self.account_ids = account_ids
		self.description = description
		
		#HACK
		self.admin_group = 'group.CPS_account_admin'
		
		if isinstance(record, crecord):
			crecord.__init__(self, _id=self._id, record=record, type=self.type, *args, **kargs)

	def dump(self,json=False):
		self.data['account_ids'] = self.account_ids
		self.data['description'] = self.description
		return crecord.dump(self,json=json)
		
	def load(self, dump):
		crecord.load(self, dump)
		if 'account_ids' in self.data:
			self.account_ids = self.data['account_ids']
		else:
			self.account_ids = []
		
		if 'description' in  self.data:
			self.description = self.data['description']
		
	def add_accounts(self, accounts,storage=None):
		if not storage:
			storage = self.storage	
			
		if not isinstance(accounts,list):
			accounts = [accounts]
			
		#string _id to account
		account_list = []
		for account in accounts:
			if isinstance(account,crecord):
				account_list.append(account)
			elif isinstance(account, str):
				if storage:
					try:
						record = storage.get(account)
						account_list.append(caccount(record,storage=storage))
					except Exception,err:
						raise Exception('Account not found: %s', err)

		#add accounts
		for account in account_list:				
				if account._id not in self.account_ids:
					self.account_ids.append(unicode(account._id))
					if self.storage:
						self.save()
					
				if self._id not in account.groups:
					account.groups.append(unicode(self._id))
					if account.storage:
						account.save()

				
	def remove_accounts(self,accounts,storage=None):
		if not storage:
			storage = self.storage
		
		if not isinstance(accounts,list):
			accounts = [accounts]
			
		#string _id to account
		account_list = []
		for account in accounts:
			if isinstance(account,crecord):
				account_list.append(account)
			elif isinstance(account, str):
				if storage:
					try:
						record = storage.get(account)
						account_list.append(caccount(record,storage=storage))
					except Exception,err:
						raise Exception('Account not found: %s', err)
						
		#remove accounts
		for account in account_list:			
			if unicode(account._id) in self.account_ids:
				self.account_ids.remove(account._id)
				if self.storage:
					self.save()
				
			if unicode(self._id) in account.groups:
				account.groups.remove(self._id)
				if account.storage:
					account.save()

