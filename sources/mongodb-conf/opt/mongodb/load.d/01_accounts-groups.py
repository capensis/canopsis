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

from caccount import caccount
from cstorage import get_storage
from crecord import crecord
from cgroup import cgroup

logger = None

##set root account
root = caccount(user="root", group="root")

#need this in two functions, key:group_name , value : group_description
groups =  {
	'CPS_root':'Have all rights.',
	'Canopsis':'Base canopsis group.',
	'CPS_curve_admin':'Create and modify curves parameters for UI.',
	'CPS_view_admin':'Manage all view in canopsis, add, remove or edit.',
	'CPS_view':'Create and manage his own view.',
	'CPS_schedule_admin':'View and create his own reporting schedules.',
	'CPS_reporting_admin':'Launch and export report, use live reporting.',
	'CPS_account_admin':'Manage all account and groups in canopsis.',
	'CPS_event_admin':'Send event by the webservice or websocket.',
	'CPS_selector_admin':'Manage all selectors in canopsis, add, remove or edit.',
	'CPS_perfdata_admin':'Manage all perfdata in canopsis, remove or edit',
	'CPS_derogation_admin':'Manage and create derogation on canopsis',
	'CPS_authkey':'Give the ability to renew the account authkey',
	'CPS_consolidation_admin':'Manage and create consolidations',
	'CPS_rule_admin':'Manage and create rules (black and white lists for incoming events)'
}

def init():
	base_init()
	update_for_new_rights()

def update():
	base_init()
	check_user_canopsis_groups()
	check_and_create_authkey()
	add_description_to_group()

def base_init():
	storage = get_storage(account=root, namespace='object')
	
	# (0'login', 1'pass', 2'group', 3'lastname', 4'firstname', 5'groups' ,6'email')
	accounts = [
		('root','root', 'CPS_root', 'Lastname', 'Firstname', [] ,''),
		('canopsis','canopsis', 'Canopsis', 'Psis', 'Cano', ['group.CPS_view'],'')
	]

	for name in groups:
		try:
			# Check if exist
			record = storage.get('group.%s' % name)
			record.data['internal'] = True
			storage.put(record)
		except:
			logger.info(" + Create group '%s'" % name)
			record = crecord({'_id': 'group.%s' % name }, type='group', name=name, group='group.CPS_account_admin')
			record.admin_group = 'group.CPS_account_admin'
			record.data['description'] = groups[name]
			record.data['internal'] = True
			record.chmod('o+r')
			storage.put(record)
		
	for account in accounts:
		user = account[0]
		try:
			# Check if exist
			record = storage.get('account.%s' % user)
		except:
			logger.info(" + Create account '%s'" % user)
			
			record = caccount(user=user, group=account[2])
			record.firstname = account[4]
			record.lastname = account[3]
			record.groups = account[5]
			record.chown(record._id)
			record.chgrp(record.group)
			record.admin_group = 'group.CPS_account_admin'
			record.chmod('g+r')
			record.passwd(account[1])
			record.generate_new_authkey()
			storage.put(record)
		

	###Root directory
	try:
		# Check if exist
		rootdir = storage.get('directory.root')
	except:
		logger.info(" + Create root directory")
		rootdir = crecord({'_id': 'directory.root','id': 'directory.root','expanded':'true'},type='view_directory', name="root directory")
		rootdir.chmod('o+r')
		storage.put(rootdir)
	
	records = storage.find({'crecord_type': 'account'}, namespace='object', account=root)
	for record in records:
		user = record.data['user']
		
		try:
			# Check if exist
			record = storage.get('directory.root.%s' % user)
		except:
			logger.info(" + Create '%s' directory" % user)
			userdir = crecord({'_id': 'directory.root.%s' % user,'id': 'directory.root.%s' % user ,'expanded':'true'}, type='view_directory', name=user)
			userdir.chown('account.%s' % user)
			userdir.chgrp('group.%s' % user)
			userdir.admin_group = 'group.CPS_view_admin'
			userdir.chmod('g-w')
			userdir.chmod('g-r')

			storage.put(userdir)
			rootdir.add_children(userdir)

			storage.put(rootdir)
			storage.put(userdir)


#add CPS_view to canopsis if needed
def check_user_canopsis_groups() :
	try:
		storage = get_storage(account=root, namespace='object')
		record = storage.get('account.canopsis')
		account = caccount(record)
		if not 'group.CPS_view' in account.groups and not 'authkey' in record.data:
			account.groups.append('group.CPS_view')
			storage.put(account)
	except:
		pass

	
	
def add_description_to_group():
	storage = get_storage(account=root, namespace='object')
	for name in groups:
		try:
			record = storage.get('group.%s' % name)
			group_record = cgroup(record)
			if not group_record.description:
				group_record.description = groups[name]
				storage.put(group_record)
		except:
			pass

	
def check_and_create_authkey():
	storage = get_storage(account=root, namespace='object')
	records = storage.find({'crecord_type': 'account'}, namespace='object', account=root)
	accounts = []
	for record in records:
		if not 'authkey' in record.data:
			#caccount auto create authkey if not provided
			accounts.append(caccount(record))
	storage.put(accounts)

def update_for_new_rights():
	#Enable rights , update old record
	storage = get_storage(account=root, namespace='object')

	dump = storage.find({})

	for record in dump:
		if record.owner.find('account.') == -1:
			record.owner = 'account.%s' % record.owner
		if record.group.find('group.') == -1:
			record.group = 'group.%s' % record.group
		#for caccount
		if 'groups' in record.data:
			for group in record.data['groups']:
				if group.find('group.') == -1:
					group = 'group.%s' % group
		#for cgroup
		if 'account_ids' in record.data:
			for account in record.data['account_ids']:
				if account.find('account.') == -1:
					account = 'account.%s' % account
	storage.put(dump)

	
	#---------------rename canopsis group, root group and curve admin-------------
	
	try:
		storage.remove('group.root')
		records = storage.find({'aaa_group':'group.root'})
		for record in records:
			record.chgrp('CPS_root')
		storage.put(records)
	except:
		pass
		
	try:
		storage.remove('group.canopsis')
		records = storage.find({'aaa_group':'group.canopsis'})
		for record in records:
			record.chgrp('CPS_canopsis')
		storage.put(records)
	except:
		pass
		
	try:
		storage.remove('group.CPS_canopsis')
		records = storage.find({'aaa_group':'group.CPS_canopsis'})
		for record in records:
			record.chgrp('Canopsis')
		storage.put(records)
	except:
		pass
		
	try:
		storage.remove('group.curves_admin')
		records = storage.find({'aaa_group':'group.curves_admin'})
		for record in records:
			record.chgrp('CPS_curves_admin')
		storage.put(records)
	except:
		pass


	
		
	#clean all groups in account.groups
	try:
		group_list = ['group.canopsis','group.root','canopsis','root','curves_admin','group.curves_admin']
		records = storage.find({'crecord_type':'account','groups':{'$in':group_list}})
		if not isinstance(records,list):
			records = [records]
		
		for record in records:
			new_groups_array = []
			for group in record.data['groups']:
				if group == 'group.canopsis' or group == 'canopsis' or group == 'CPS_canopsis' or group == 'group.CPS_canopsis':
					group = 'Canopsis'
				if group == 'group.root' or group == 'root':
					group = 'CPS_root'
				if group == 'group.curves_admin' or group == 'curves_admin':
					group = 'CPS_curves_admin'
				new_groups_array.append(group)
			record.data['groups'] = new_groups_array

		storage.put(records)
	except:
		pass
	#---------------------update each record type--------------------
	#update view
	dump = storage.find({'$or': [{'crecord_type':'view'},{'crecord_type':'view_directory'}]})
	for record in dump:
		record.chgrp('group.CPS_view_admin')
		record.admin_group = 'group.CPS_view_admin'
		record.chmod('g+w')
		record.chmod('g+r')
	storage.put(dump, account=root)
	
	#update schedule
	dump = storage.find({'crecord_type':'schedule'})
	for record in dump:
		record.chgrp('group.CPS_schedule_admin')
		record.admin_group = 'group.CPS_schedule_admin'
		record.chmod('g+w')
		record.chmod('g+r')
	storage.put(dump)
	
	#update accounts
	dump = storage.find({'$or': [{'crecord_type':'account'},{'crecord_type':'group'}]})
	for record in dump:
		#record.chgrp('group.CPS_account_admin')
		record.admin_group = 'group.CPS_account_admin'
		record.chmod('g+w')
		record.chmod('g+r')
	storage.put(dump)
