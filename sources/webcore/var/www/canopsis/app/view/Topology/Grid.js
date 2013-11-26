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
Ext.define('canopsis.view.Topology.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.TopologyGrid',

	model: 'Topology',
	store: 'Topologies',

	opt_db_namespace: 'object',

	opt_menu_delete: true,
	opt_bar_duplicate: false,
	opt_menu_rights: false,
	opt_bar_enable: true,

    opt_export_import: true,

	columns: [
		{
			header: '',
			width: 25,
			sortable: false,
			renderer: rdr_crecord_type,
			dataIndex: 'crecord_type'
		},{
			header: _('State'),
			align: 'center',
			width: 50,
			dataIndex: 'state',
			renderer: rdr_status
		},{
			header: _('Loaded'),
			align: 'center',
			width: 55,
			dataIndex: 'loaded',
			renderer: rdr_boolean
		},{
			header: _('Enabled'),
			align: 'center',
			width: 55,
			dataIndex: 'enable',
			renderer: rdr_boolean
		},{
			header: _('Name'),
			flex: 1,
			sortable: true,
			dataIndex: 'crecord_name'
		},{
			header: _('Display name'),
			flex: 1,
			sortable: true,
			dataIndex: 'display_name'
		},{
			header: _('Description'),
			flex: 2,
			dataIndex: 'description'
		},{
			flex: 1,
			dataIndex: 'aaa_owner',
			renderer: rdr_clean_id,
			text: _('Owner')
		},{
			flex: 1,
			dataIndex: 'aaa_group',
			renderer: rdr_clean_id,
			text: _('Group')
		}
	],

	initComponent: function() {
		this.callParent(arguments);
	}

});
