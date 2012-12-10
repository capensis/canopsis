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
Ext.define('canopsis.view.Aggregation.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.AggregationGrid',

	model: 'Perfdata',
	store: 'Perfdatas',

	opt_db_namespace: 'object',

	opt_menu_delete: true,
	opt_bar_duplicate: true,
	opt_menu_rights: false,
	opt_bar_enable: true,

	columns: [
		{
			header: '',
			width: 25,
			sortable: false,
			//renderer: rdr_crecord_type,
			dataIndex: 'crecord_type'
        },{
			header: _('State'),
			align: 'center',
			width: 50,
			dataIndex: 'state',
			renderer: rdr_status
		},{
			header: _('Enabled'),
			align: 'center',
			width: 55,
			dataIndex: 'enable',
			renderer: rdr_boolean
		},{
			header: _('Loaded'),
			align: 'center',
			width: 55,
			dataIndex: 'loaded',
			renderer: rdr_boolean
		},{
			header: _('Name'),
			flex: 1,
			sortable: true,
			dataIndex: 'crecord_name'
		},{
			header: _('Aggregation type'),
			flex: 1,
			dataIndex: 'description'
		},{
			header: _('Selection regex'),
			align: 'center',
			width: 2,
			dataIndex: 'sla_state',

		}
	],

	initComponent: function() {
		this.callParent(arguments);
	}

});
