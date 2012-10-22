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

logger = None

##set root account
root = caccount(user="root", group="root")
storage = get_storage(account=root, namespace='object')

def init():
	pass

def update():
	## Update 201111 -> 201201
	#search all user view and add them to root directory
	#	views = storage.find({'crecord_type':'view'}, account=account)
	#	for view in views:
	#		#if this account is the owner of the view
	#		if view.owner == user:
	#			#if view is not in administration views
	#			if view._id not in ['view._default_.dashboard','view.ComponentDetails','view.components','view.resources','view.group_manager','view.account_manager','view.view_manager']:
	#				if not view.parent:
	#					print(view.dump()['crecord_name'])
	#					view.data['items'] = []
	#					view.data['leaf'] = True
	#					userdir.add_children(view)
	#					storage.put(view)
	#	storage.put(userdir)
	pass
