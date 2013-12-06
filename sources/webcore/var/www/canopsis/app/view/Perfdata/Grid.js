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
Ext.define('canopsis.view.Perfdata.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.PerfdataGrid',

	model: 'Perfdata',
	store: 'Perfdatas',

	opt_paging: true,
	opt_menu_delete: true,
	opt_bar_add: false,
	opt_menu_rights: false,
	opt_bar_search: true,

	opt_tags_search: false,
	opt_simple_search: true,

	opt_cell_edit: false,


	columns: [
		{
			header: '',
			width: 25,
			sortable: false,
			renderer: function() {
				return "<span class='icon icon-mainbar-perfdata' />";
			}
		},{
			header: _('Component'),
			flex: 1,
			sortable: true,
			dataIndex: 'co'
		},{
			header: _('Resource'),
			flex: 1,
			sortable: true,
			dataIndex: 're'
		},{
			header: _('Metric'),
			flex: 2,
			sortable: true,
			dataIndex: 'me',
			editor: {xtype: 'textfield'}
		},{
			header: _('Retention') +' ' +  _('in seconds'),
			width: 150,
			sortable: true,
			dataIndex: 'r',
			align: 'center',
			renderer: rdr_time_interval,
			editor: {
				xtype: 'numberfield',
				minValue: 0,
				step: 60
			}
		},{
			header: _('First point'),
			width: 150,
			sortable: true,
			dataIndex: 'fts',
			align: 'center',
			renderer: rdr_tstodate
		},{
			header: _('Last point'),
			width: 150,
			sortable: true,
			dataIndex: 'lts',
			align: 'center',
			renderer: rdr_tstodate
		},{
			header: _('Type'),
			width: 70,
			sortable: true,
			dataIndex: 't',
			align: 'center',
			editor: {xtype: 'textfield'}
		},{
			header: _('Min'),
			width: 100,
			sortable: true,
			dataIndex: 'mi',
			align: 'right',
			renderer: function(value, metaData, record) {
				void(metaData);

				return rdr_humanreadable_value(value, record.get('u'));
			}
		},{
			header: _('Max'),
			width: 100,
			sortable: true,
			dataIndex: 'ma',
			align: 'right',
			renderer: function(value, metaData, record) {
				void(metaData);

				return rdr_humanreadable_value(value, record.get('u'));
			}
		},{
			header: _('Last value'),
			width: 100,
			sortable: true,
			dataIndex: 'lv',
			align: 'right',
			renderer: function(value, metaData, record) {
				void(metaData);

				return rdr_humanreadable_value(value, record.get('u'));
			}
		},{
			header: _('Unit'),
			width: 45,
			sortable: true,
			dataIndex: 'u',
			align: 'center',
			editor: {xtype: 'textfield'}
		},{
			header: _('Tags'),
			flex: 2,
			sortable: false,
			dataIndex: 'tg',
			renderer: rdr_tags
		}

	],

	bar_search: [{
		xtype: 'button',
		iconCls: 'icon-internal-metrics',
		pack: 'end',
		tooltip: _('Display internal metrics'),
		enableToggle: true,
		action: 'toggle_internal_metric'
	}],

	initComponent: function() {
		this.callParent(arguments);
	}
});
