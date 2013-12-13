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
Ext.define('canopsis.view.Statemap.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	controllerId: 'Statemap',

	alias: 'widget.StatemapGrid',

	model: 'Statemap',
	store: 'Statemaps',

	opt_paging: true,
	opt_menu_delete: true,
	opt_bar_add: true,
	opt_menu_rights: false,
	opt_bar_search: true,
	opt_bar_enable: true,
	opt_tags_search: false,

	rdr_statemap: function(val) {
		var output = '<p>';

		for(var i = 0; i < val.length; ++i) {
			output += '<b>' + i + '</b> -> ';

			switch(val[i]) {
				case 0:
					output += 'OK';
					break;

				case 1:
					output += 'WARNING';
					break;

				case 2:
					output += 'CRITICAL';
					break;

				default:
					output += 'UNKNOWN';
					break;
			}

			if (i != val.length - 1) {
				output += ', ';
			}
		}

		return output;
	},

	initComponent: function() {
		this.columns = [{
			header: '',
			width: 25,
			sortable: false,
			renderer: rdr_crecord_type,
			dataIndex: 'crecord_type',
		},{
			header: _('Enabled'),
			align: 'center',
			width: 55,
			dataIndex: 'enable',
			renderer: rdr_boolean,
		},{
			header: _('Name'),
			flex: 1,
			dataIndex: 'crecord_name',
		},{
			header: _('Statemap'),
			flex: 1,
			dataIndex: 'statemap',
			renderer: this.rdr_statemap,
		}];

		this.callParent();
	}
});
