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
Ext.define('canopsis.view.Selector.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.SelectorGrid',

	model: 'Selector',
	store: 'Selectors',

	opt_db_namespace: 'object',

	opt_menu_delete: true,
	opt_bar_duplicate: true,
	opt_menu_rights: true,
	opt_bar_enable: true,
	opt_paging: false,

	opt_bar_search: true,
	opt_bar_search_field: ['crecord_name', 'display_name', 'description'],

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
			header: _('Calcul state'),
			align: 'center',
			width: 80,
			dataIndex: 'dostate',
			renderer: rdr_boolean
		},{
			header: _('Calcul SLA'),
			align: 'center',
			width: 80,
			dataIndex: 'dosla',
			renderer: rdr_boolean
		},{
			header: _('SLA'),
			align: 'center',
			width: 50,
			dataIndex: 'sla_state',
			renderer: rdr_status
		},{
			header: _('SLA') + ': ' + _('time window'),
			align: 'center',
			dataIndex: 'sla_timewindow',
			width: 150,
			renderer: rdr_time_interval
		},{
			header: _('SLA Value'),
			align: 'center',
			width: 60,
			dataIndex: 'sla_timewindow_perfdata',
			renderer: function(val) {
				if(val) {
					perf = val[0];
					return perf.value + perf.unit;
				}
			}
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
		},{
			xtype: 'actioncolumn',
			width: 70,
			text: _('Derogation'),
			icon: './themes/canopsis/resources/images/icons/edit.png',
			iconCls: 'icon-clickable',
			handler: function(grid, rowIndex) {
				var rec = grid.getStore().getAt(rowIndex).raw;
				log.dump(rec);

				if(rec.rk) {
					global.derogationCtrl.derogate(rec.rk, rec.crecord_name);
				}
				else {
					global.notify.notify(_('Information not found'), _("Please wait a moment, some informations aren't availables"), 'info');
				}
			}
		}
	],

	initComponent: function() {
		this.callParent(arguments);
	}

});
