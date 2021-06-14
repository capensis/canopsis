//need:app/lib/view/cgrid.js
/*
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
*/
Ext.define('canopsis.view.Schedule.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.ScheduleGrid',

	model: 'Schedule',
	store: 'Schedules',

	opt_paging: false,

	opt_db_namespace: 'object',

	opt_menu_delete: true,
	opt_menu_run_item: true,
	opt_bar_duplicate: true,
	opt_menu_rights: true,

	columns: [
		{
			header: _('Loaded'),
			width: 55,
			dataIndex: 'loaded',
			renderer: rdr_boolean
		},{
			header: _('Success'),
			width: 55,
			sortable: true,
			dataIndex: 'log',
			renderer: function(value) {
				return rdr_boolean(value.success);
			}
		},{
			header: _('Last execution'),
			flex: 1,
			sortable: true,
			dataIndex: 'log',
			renderer: function(value) {
				return rdr_tstodate(value.timestamp);
			}
		},{
			header: _('Next execution'),
			flex: 1,
			sortable: true,
			dataIndex: 'next_run_time',
			renderer: rdr_utcToLocal
		},{
			header: _('Schedule'),
			flex: 1,
			sortable: true,
			dataIndex: 'cron',
			renderer: rdr_task_crontab
		},{
			header: _('Name'),
			flex: 2,
			sortable: true,
			dataIndex: 'crecord_name'
		},{
			header: _('Message'),
			flex: 5,
			sortable: true,
			dataIndex: 'log',
			renderer: function(value) {
				var celery = value.celery_output;
				var duration = value.duration;

				if (celery !== undefined && duration !== undefined) {
					return celery + ' (in ' + duration + 's)';
				}
			}
		},{
			header: _('Mailing'),
			width: 55,
			dataIndex: 'mail',
			renderer: rdr_boolean
		}
	],

	initComponent: function() {
		this.callParent(arguments);
	}
});
