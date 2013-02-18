/*
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
*/
Ext.define('canopsis.model.Schedule', {
	extend: 'Ext.data.Model',
	fields: [
		{name: '_id'},
		{name: 'id', mapping: '_id'},
		{name: 'crecord_type', defaultValue: 'schedule'},
		{name: 'func_ref'},
		{name: 'loaded', defaultValue: false},
		{name: 'crecord_name'},
		{name: 'args', defaultsValue: []},
		{name: 'kwargs' , defaultsValue: {}},
		{name: 'next_run_time'},
		{name: 'cron', defaultValue: undefined},
		{name: 'log'},
		{
			name: 'log_success',
			convert: function(value, record) {return record.get('log').success}
		},
		{
			name: 'log_output',
			convert: function(value, record) {
				var celery = record.get('log').celery_output;
				var duration = record.get('log').duration;
				if (celery != undefined && duration != undefined) {
					return celery + ' (in ' + duration + 's)';
				}
			}
		},
		{
			name: 'log_last_execution',
			convert: function(value, record) {return record.get('log').timestamp}
		},
		{
		 name: 'mail',
		 convert: function(value, record) {
					var kwargs = record.get('kwargs');
					if (kwargs['mail'] != undefined && kwargs['mail'].sendMail)
						return true;
					return false;
				}
		},

		{name: 'aaa_access_owner', defaultValue: ['r', 'w']},
		{name: 'aaa_access_group', defaultValue: ['r']},
		{name: 'aaa_access_other', defaultValue: []},
		{name: 'aaa_admin_group'},
		{name: 'aaa_group'},
		{name: 'aaa_owner'},

		{name: 'exporting_intervalLength'},
		{name: 'exporting_intervalUnit'},	
		{name: 'frequency',defaultValue:'day'},


		{name: 'exporting_interval'},
		{name: 'exporting_account'},
		{name: 'exporting_task'},
		{name: 'exporting_method'},
		{name: 'exporting_owner'},
		{name: 'exporting_viewName'},

		{name: 'exporting_mail'},
		{name: 'exporting_recipients'},
		{name: 'exporting_subject'},
		
		{name: 'crontab_hours'},
		{name: 'crontab_month'},
		{name: 'crontab_day_of_week'},
		{name: 'crontab_day'},


	]
});
