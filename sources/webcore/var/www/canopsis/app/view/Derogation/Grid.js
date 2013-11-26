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
Ext.define('canopsis.view.Derogation.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	controllerId: 'Derogation',

	alias: 'widget.DerogationGrid',

	model: 'Derogation',
	store: 'Derogations',

	opt_paging: true,
	opt_menu_delete: true,
	opt_bar_add: true,
	opt_menu_rights: false,
	opt_bar_search: true,
	opt_bar_enable: true,
	opt_tags_search: false,

	rdr_time_conditions: function(val) {
		var output = '';

		for(var i = 0; i < val.length; i++) {
			var condition = val[i];
			var type = condition['type'];

			if(type === 'time_interval') {
				output += '<b>Type</b>: ' + type + '<br>';
				output += '&nbsp; + <b>' + _('Start') + '</b>: ' + rdr_tstodate(condition['startTs']) + '<br>';
				output += '&nbsp; + <b>' + _('Stop') + '</b>: ' + rdr_tstodate(condition['stopTs']) + '<br>';
			}
		}

		return output;
	},

	rdr_override: function(val) {
		var output = '';
		for(var i = 0; i < val.length; i++) {
			var action = val[i];
			var type = action['type'];

			if(type === 'override') {
				output += '<b>' + action['field'] + "</b> -&gt; '" + action['value'] + "'<br>";
			}
		}

		return output;
	},

	rdr_statemap: function(val) {
		var found = false;

		for(var i = 0; i < val.length; ++i) {
			if (val[i].type == 'requalificate') {
				found = true;
				break;
			}
		}

		return rdr_boolean(found);
	},

	initComponent: function() {
		this.columns = [
			{
				header: '',
				width: 25,
				sortable: false,
				renderer: rdr_crecord_type,
				dataIndex: 'crecord_type'
			},{
				header: _('Enabled'),
				align: 'center',
				width: 55,
				dataIndex: 'enable',
				renderer: rdr_boolean
			},{
				header: _('Active'),
				align: 'center',
				width: 55,
				dataIndex: 'active',
				renderer: rdr_boolean
			},{
				header: _('Name'),
				flex: 1,
				dataIndex: 'crecord_name'
			},{
				header: _('Description'),
				flex: 2,
				dataIndex: 'description'
			},{
				header: _('Condition'),
				flex: 1,
				dataIndex: 'time_conditions',
				renderer: this.rdr_time_conditions
			},{
				header: _('Override'),
				flex: 1,
				dataIndex: 'actions',
				renderer: this.rdr_override
			},{
				header: _('Statemap'),
				flex: 1,
				dataIndex: 'actions',
				renderer: this.rdr_statemap,
			}
		];

		this.callParent();
	}
});
