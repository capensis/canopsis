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

import time
import hashlib

from crecord import crecord
from random import getrandbits

try:
	from cgroup import cgroup
except:
	pass

class caccount(crecord):
	def __init__(self, record=None, user=None, group=None, lastname=None, firstname=None, mail=None, groups=[], authkey=None, *args, **kargs):
		self.user = user or "anonymous"
		self.groups = groups
		self.group = group or "group.anonymous"
		self.shadowpasswd = None
		
		self.authkey = authkey or self.generate_new_authkey()

		self.lastname = lastname
		self.firstname = firstname
		self.mail = mail

		self.type = "account"

		self._id = self.type + "." + self.user
		
		self.access_owner = ['r', 'w']
		self.access_group = []
		self.access_other = []
		self.access_unauth = []

		self.external = False

		if isinstance(record, crecord):
			crecord.__init__(self, _id=self._id, record=record, type=self.type, *args, **kargs)
		else:
			crecord.__init__(self, _id=self._id, owner="account.%s" % self.user, group=self.group, type=self.type, *args, **kargs)

	def get_full_mail(self):
		return "\"%s %s\" <%s>" % (self.firstname, self.lastname, self.mail)

	def passwd(self, passwd):
		self.shadowpasswd = hashlib.sha1(str(passwd)).hexdigest()

	def check_passwd(self, passwd):
		return self.check_shadowpasswd(hashlib.sha1(str(passwd)).hexdigest())

	def check_shadowpasswd(self, shadowpasswd):
		shadowpasswd = str(shadowpasswd).upper()
		if shadowpasswd == str(self.shadowpasswd).upper():
			return True

		return False
		
	def make_shadow(self, passwd):
		return hashlib.sha1(str(passwd)).hexdigest()
	
	
	def check_tmp_cryptedKey(self, authkey):
		authkey =  str(authkey).upper()
		if authkey == str(self.make_tmp_cryptedKey(self.shadowpasswd)).upper():
			return True
		
		return False
		
	def make_tmp_cryptedKey(self, shadow=None):
		if not shadow:
			shadow = self.shadowpasswd
			
		return hashlib.sha1(str(shadow).upper() + str( int( time.time() / 10)*10 )).hexdigest()

	def get_authkey(self):
		return self.authkey
		
	def check_authkey(self, authkey):
		if str(authkey).upper() == str(self.get_authkey()).upper():
			return True
		else:
			return False
		
	def get_mail_md5(self):
		m = hashlib.md5()
		m.update(self.mail)
		return m.hexdigest()

	def generate_new_authkey(self):
		return hashlib.sha224(str(getrandbits(512))).hexdigest()

	def dump(self):
		self.name = self.user
		self.data['user'] = self.user
		self.data['lastname'] = self.lastname
		self.data['firstname'] = self.firstname
		self.data['mail'] = self.mail
		self.data['groups'] = list(self.groups)
		self.data['external'] = self.external
		'''
		if self.group:
			self.data['groups'].insert(0, self.group)
		'''
		self.data['shadowpasswd'] = self.shadowpasswd
		self.data['authkey'] = self.authkey
		return crecord.dump(self)

	def load(self, dump):
		crecord.load(self, dump)
		self.user = self.data['user']
		self.lastname = self.data['lastname']
		self.firstname = self.data['firstname']
		self.mail = self.data['mail']
		self.groups = self.data['groups']
		self.external = self.data.get('external', self.external)
		'''
		if len(self.groups) > 0:
			if self.groups[0] == self.group:
				self.groups.pop(0)
		'''
		self.shadowpasswd = self.data['shadowpasswd']
		if 'authkey' in self.data:
			self.authkey = self.data['authkey']

	def cat(self):
		print "Id:\t", self._id
		print " + Fullname:\t", self.firstname, self.lastname
		print " + User:\t", self.user
		print " + Mail:\t", self.mail
		print " + Owner:\t", self.owner
		print " + Group:\t", self.group
		print " + Groups:\t", self.groups, "\n"
		
	def add_in_groups(self, groups, storage=None):
		if not storage:
			storage = self.storage
		
		if not isinstance(groups,list):
			groups = [groups]
			
		# String _id to cgroup
		group_list = []
		for group in groups:
			if isinstance(group,cgroup):
				group_list.append(group)
			elif isinstance(group, basestring):
				if storage:
					try:
						record = storage.get(group)
						group_list.append(cgroup(record,storage=storage))
					except Exception,err:
						raise Exception('Group not found: %s', err)
												
		# Add to groups
		for group in group_list:
				if unicode(group._id) not in self.groups:
					self.groups.append(unicode(group._id))
					if self.storage:
						self.save()
					
				if unicode(self._id) not in group.account_ids:
					group.account_ids.append(unicode(self._id))
					if group.storage:
						group.save()
	
	def remove_from_groups(self, groups, storage=None):
		if not storage:
			storage = self.storage
		
		if not isinstance(groups,list):
				groups = [groups]
				
		# String _id to cgroup
		group_list = []
		for group in groups:
			if isinstance(group,crecord):
				group_list.append(group)
			elif isinstance(group, str)  and isinstance(group, unicode):
				if storage:
					try:
						record = storage.get(group)
						group_list.append(cgroup(record, storage=storage))
					except Exception,err:
						raise Exception('Group not found: %s', err)
		# Remove groups
		for group in group_list:
				if unicode(group._id) in self.groups:
					self.groups.remove(group._id)
					if self.storage:
						self.save()
					
				if unicode(self._id) in group.account_ids:
					group.account_ids.remove(self._id)
					if group.storage:
						group.save()

#################

def caccount_getall(storage):
	accounts = []
	records = storage.find({'crecord_type': 'account'})
	for record in records:
		accounts.append(caccount(record))
	
	return accounts

def caccount_get(storage, user):
	record = storage.get('account.'+user)
	account = caccount(record)
	return account
