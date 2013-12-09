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
	### Default Dasboard
	data = [{'position': {'width': 8, 'top': 2, 'left': 8, 'height': 7}, 'data': {'bar_search': False, 'show_last_check': True, 'xtype': 'list', 'pageSize': 100, 'title': 'Resource problems', 'show_source_type': True, 'border': True, 'default_sort_direction': 'DESC', 'scroll': True, 'filter': '{ "$and": [ {"source_type":"resource"}, {"state": { "$ne": 0 }} ]}', 'default_sort_column': 'state', 'paging': False, 'show_resource': True, 'reload': False, 'show_state': True, 'refreshInterval': 300, 'show_output': True, 'show_state_type': True, 'column_sort': True, 'hideHeaders': False, 'show_component': True}, 'id': '1336723949800-5'}, {'position': {'width': 8, 'top': 2, 'left': 0, 'height': 7}, 'data': {'bar_search': False, 'show_last_check': True, 'xtype': 'list', 'pageSize': 100, 'title': 'Component problems', 'show_source_type': True, 'border': True, 'default_sort_direction': 'DESC', 'scroll': True, 'filter': '{ "$and": [ {"source_type":"component"}, {"state": { "$ne" : 0 }} ]}', 'default_sort_column': 'state', 'paging': False, 'show_resource': False, 'reload': False, 'show_state': True, 'refreshInterval': 300, 'show_output': True, 'show_state_type': True, 'column_sort': True, 'hideHeaders': False, 'show_component': True}, 'id': '1336724023524-4'}, {'position': {'width': 4, 'top': 0, 'left': 0, 'height': 2}, 'data': {'refreshInterval': 0, 'title': '', 'border': False, 'xtype': 'text', 'text': '<img src="themes/canopsis/resources/images/logo_canopsis.png" height="100%">'}, 'id': '1336724801997-7'}]
	create_view('_default_.dashboard', 'Dashboard', data, autorm=False)
		
	### Account
	data = { 'xtype': 'AccountGrid'}
	create_view('account_manager', 'Accounts', data, internal=True)

	### Group
	data = { 'xtype': 'GroupGrid'}
	create_view('group_manager', 'Groups', data, internal=True)
	
	### Selector
	data = { 'xtype': 'SelectorGrid'}
	create_view('selector_manager', 'Selectors', data, internal=True)

	### Components
	data = { 'xtype': 'list', 'show_tags': True,'fitler_buttons': True, 'filter': '{"$and": [{"source_type":"component"}, {"event_type": {"$ne": "comment"}}, {"event_type": {"$ne": "user"}}]}', 'show_resource': False}
	create_view('components', 'Components', data, internal=True)

	### Resources
	data = { 'xtype': 'list', 'show_tags': True,'fitler_buttons': True, 'filter': '{"$and": [{"source_type":"resource"}, {"event_type": {"$ne": "comment"}}, {"event_type": {"$ne": "user"}}]}'}
	create_view('resources', 'Resources', data, internal=True)

	### View manager
	data = { 'xtype': 'ViewTreePanel'}
	create_view('view_manager', 'Views', data, internal=True)

	###task
	data = { 'xtype': 'ScheduleGrid'}
	create_view('schedule_manager', 'Schedules', data, internal=True)

	###briefcase
	data = { 'xtype': 'BriefcaseGrid'}
	create_view('briefcase', 'Briefcase', data, internal=True)
	
	###curves
	data = { 'xtype': 'CurvesGrid'}
	create_view('curves', 'Curves', data, internal=True)
	
	###derogation
	data = {'xtype':'DerogationGrid'}
	create_view('derogation_manager','Derogations',data, internal=True)
	
	###statemap
	data = {'xtype': 'StatemapGrid'}
	create_view('statemap_manager', 'Statemaps', data, internal=True)

	###perfdata
	data = {'xtype':'PerfdataGrid'}
	create_view('perfdata','Perfdata',data, internal=True)

	###Event log navigation
	data = { 'xtype': 'EventLog'}
	create_view('eventLog_navigation', 'Events log navigation', data, internal=True)
	
	### Topology
	data = { 'xtype': 'TopologyGrid'}
	create_view('topology_manager', 'Topologies', data, internal=True)

	### Consolidation
	data = { 'xtype': 'ConsolidationGrid'}
	create_view('consolidation_manager', 'Consolidation', data, internal=True)

	### Filter
	data = { 'xtype': 'RuleGrid' }
	create_view('rules_manager', 'Filter Rules', data, internal=True)

	###metric_navigator
	#data = {'xtype': 'MetricNavigation'}
	#create_view('metric_navigation', 'Metrics Navigation', data)

def update():
	init()
	update_view_for_new_metric_format()

def create_view(_id, name, data, position=None, mod='o+r', autorm=True, internal=False):
	#Delete old view
	try:
		record = storage.get('view.%s' % _id)
		if autorm:
			storage.remove(record)
		else:
			return record
	except:
		pass
		
	if not position:
		# fullscreen
		position = {'width': 1,'top': 0, 'left': 0, 'height': 1}
		
	logger.info(" + Create view '%s'" % name)
	record = crecord({'_id': 'view.%s' % _id, 'internal': internal }, type='view', name=name,group='group.CPS_view_admin')
	
	if  isinstance(data, list):
		record.data['items'] = data
	elif  isinstance(data, dict):
		record.data['items'] = [ {'position': position, 'data': data } ]
	else:
		raise("Invalide data ...")
		
	record.chmod(mod)
	storage.put(record)
	return record

def update_view_for_new_metric_format():
	records = storage.find({'crecord_type': 'view'}, namespace='object', account=root)
	for view in records:
		for item in view.data['items']:
			nodesObject = {}

			#check if old format
			if 'nodes' in item['data']:
				itemNodes = item['data']['nodes']

				if isinstance(itemNodes, list):
					itemXtype = item['data']['xtype']

					if itemXtype == 'weather':
						print('Ignore for weather widget')
						break

					#update for text widget
					if itemXtype == 'text' or itemXtype == 'topology_viewer':
						print('Update widget text/topology_viewer format')
						item['data']['inventory'] = item['data']['nodes']
						del item['data']['nodes']
						break

					for node in itemNodes:
						try:
							nodesObject[node['id']] = node

							# write extra_fields in node root
							if 'extra_field' in node:
								nodesObject[node['id']].update(node['extra_field'])

								#build ccustom in view
								del node['extra_field']
						except Exception as error:
							print('An error occured for the following widget: %s' % error)
							print(item)

					item['data']['nodes'] = nodesObject
					print(item['data']['nodes'])

				#check between commits
				if 'ccustom' in item['data']:
					if isinstance(item['data']['ccustom'], dict):
						for nodeId, customValue in item['data']['ccustom'].iteritems():
							if nodeId in itemNodes:
								itemNodes[nodeId].update(customValue)
						del item['data']['ccustom']
				
	storage.put(records)
				

