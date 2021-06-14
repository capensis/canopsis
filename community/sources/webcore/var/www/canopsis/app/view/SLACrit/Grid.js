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
Ext.define('canopsis.view.SLACrit.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.SLACritGrid',

	model: 'SLACrit',
	store: 'SLACrit',

	opt_db_namespace: 'object',

	opt_menu_delete: true,
	opt_bar_duplicate: true,
	opt_menu_rights: true,
	opt_bar_enable: true,
	opt_paging: false,

	opt_bar_search: false,

	// opt_bar_customs: [{
	// 	alias: 'widget.defaultAction',
	// 	fieldLabel: _('Default action'),
	// 	xtype: 'combobox',
	// 	iconCls: 'icon-book',
	// 	editable: false,
	// 	store: {
	// 		xtype: 'store',
	// 		fields: ['value', 'text'],
	// 		data: [
	// 			{value: 'pass', text: _('Pass')},
	// 			{value: 'drop', text: _('Drop')}
	// 		]
	// 	}
	// }],

	opt_bar_customs: [{
		text: 'Macros',
		xtype: 'button',
		iconCls: 'icon-mainbar-edit-task',
		action: 'edit-macros-button',
	}],

	columns: [
		{
			header: '',
			width: 25,
			sortable: false,
			renderer: rdr_crecord_type,
			dataIndex: 'crecord_type'
		},{
			header: 'Criticity',
			sortable: true,
			// renderer: rdr_crecord_type,
			dataIndex: 'crit'
		},{
			header: 'Delay',
			sortable: true,
			renderer: rdr_time_interval,
			dataIndex: 'delay'
		}
	],

	initComponent: function() {
		this.callParent(arguments);
	}
});
