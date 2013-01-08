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
Ext.define('canopsis.view.Consolidation.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.ConsolidationGrid',

	model: 'Consolidation',
	store: 'Consolidations',

	opt_db_namespace: 'object',

	opt_menu_delete: true,
	opt_bar_duplicate: true,
	opt_menu_rights: false,
	opt_bar_enable: true,

	columns: [
		/*{
			header: _('State'),
			align: 'center',
			width: 50,
			dataIndex: 'state',
			renderer: rdr_status
		}/*,{
			header: _('Enabled'),
			align: 'center',
			width: 55,
			dataIndex: 'enable',
			renderer: rdr_boolean
		},*/{
			header: _('Loaded'),
			align: 'center',
			width: 55,
			dataIndex: 'loaded',
			renderer: rdr_boolean
		},{
			header: _('Enable'),
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
			header: _('Component'),
			flex: 1,
			sortable: true,
			dataIndex: 'component'
		},{
			header: _('Resource'),
			flex: 1,
			sortable: true,
			dataIndex: 'resource'
		}, {
			header: _('Aggregation type'),
			flex: 3,
			dataIndex: 'type'
		},{
			header: _('Output engine'),
			//align: 'center',
			flex: 6,
			dataIndex: 'output_engine',
			//renderer : function(){return 'Ressource = Regex(cpu*) && Metric = system '}
		},{
			header: _('Nb elements'),
			align: 'center',
			width: 70,
			dataIndex: 'nb_items',
			//renderer : function(){return 34}
		}
	],

	initComponent: function() {
		this.callParent(arguments);
	}

});
