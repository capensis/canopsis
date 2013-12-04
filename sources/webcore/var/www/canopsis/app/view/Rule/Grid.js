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
Ext.define('canopsis.view.Rule.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.RuleGrid',

	model: 'Rule',
	store: 'Rules',

	opt_db_namespace: 'object',

	opt_menu_delete: true,
	opt_bar_duplicate: true,
	opt_menu_rights: true,
	opt_bar_enable: true,
	opt_paging: false,

	opt_bar_search: true,
	opt_bar_search_field: ['crecord_name', 'priority', 'rule', 'action'],

	opt_bar_customs: [{
		alias: 'widget.defaultAction',
		fieldLabel: _('Default action'),
		xtype: 'combobox',
		iconCls: 'icon-book',
		editable: false,
		store: {
			xtype: 'store',
			fields: ['value', 'text'],
			data: [
				{value: 'pass', text: _('Pass')},
				{value: 'drop', text: _('Drop')}
			]
		}
	}],

	columns: [
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
		sortable: false,
		renderer: rdr_boolean
	},{
		header: 'Priority',
		width: 55,
		sortable: true,
		dataIndex: 'priority'
	},{
		header: _('Name'),
		width: 200,
		sortable: false,
		dataIndex: 'crecord_name'
	},{
		header: _('Rule'),
		flex: 1,
		sortable: false,
		dataIndex: 'mfilter'
	},{
		header: _('Action'),
		align: 'center',
		width: 55,
		sortable: false,
		dataIndex: 'action',
		renderer: rdr_rule_action
	}
	],

	initComponent: function() {
		this.callParent(arguments);
	}

});
