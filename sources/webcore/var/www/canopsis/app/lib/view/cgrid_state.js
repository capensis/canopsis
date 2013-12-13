//need:app/lib/view/cgrid.js,app/lib/store/cstore.js
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
Ext.define('canopsis.lib.view.cgrid_state' , {
	extend: 'canopsis.lib.view.cgrid',

	requires: [
		'canopsis.lib.store.cstore'
	],

	store: false,
	filter: false,
	autoload: false,
	remoteSort: false,

	opt_paging: false,
	opt_bar: false,

	opt_allow_edit: false,

	opt_column_sortable: false,

	opt_show_state_type: true,
	opt_show_component: false,
	opt_show_resource: true,
	opt_show_state: true,
	opt_show_state_type: true,
	opt_show_source_type: true,
	opt_show_last_check: true,
	opt_show_output: true,
	opt_show_tags: true,
	opt_show_ticket: true,
	
	opt_show_row_background: true,

	opt_bar_delete: false,
	opt_bar_add: false,

	border: false,

	namespace: 'events',

	pageSize: global.pageSize,

	sorters: [{
		property: 'state',
		direction: 'DESC'
	}],

	columns: [],

	fitler_buttons: false,

	initComponent: function() {
		this.columns = [];

		console.log( this );

		if ( this.opt_show_form_ack && ( this.opt_show_ack_state_solved || this.opt_show_ack_state_pendingsolved || this.opt_show_ack_state_pendingaction || this.opt_show_ack_state_pendingvalidation ) ) {
			this.selModel= Ext.create('Ext.selection.CheckboxModel');	
			this.opt_bar_ack = true;
		}
		//set columns
		if (this.opt_show_source_type) {
			this.columns.push({
				header: ' ',
				width: 25,
				sortable: this.opt_column_sortable,
				//menuDisabled: true,
				hideable: false,
				dataIndex: 'source_type',
				renderer: rdr_source_type
			});
		}
		
		if (this.opt_show_file_help && this.opt_file_help_url != null && this.opt_file_help_url != "" ) {
			this.columns.push({
				header: '<span class="icon icon-file_help"></span>',
				width: 25,
				sortable: this.opt_column_sortable,
				//menuDisabled: true,
				hideable: true,
				renderer: rdr_file_help
			});
		}
		
		if (this.opt_show_file_equipement && this.opt_file_equipement_url != null && this.opt_file_equipement_url != "" ) {
			this.columns.push({
				header: '<span class="icon icon-file_equipement"></span>',
				width: 25,
				sortable: this.opt_column_sortable,
				//menuDisabled: true,
				hideable: true,
				renderer: rdr_file_equipement
			});
		}
		
		if (this.opt_show_ack) {
			this.columns.push({
				header: '<span class="icon"></span>',
				width: 25,
				sortable: this.opt_column_sortable,
				//menuDisabled: true,
				hideable: true,
				renderer: rdr_ack
			});
		}

		if (this.opt_show_state_type) {
			this.columns.push({
				header: 'ST',
				sortable: this.opt_column_sortable,
				width: 25,
				dataIndex: 'state_type',
				renderer: rdr_state_type
			});
		}

		if (this.opt_show_state) {
			this.columns.push({
				header: _('S'),
				sortable: this.opt_column_sortable,
				width: 25,
				dataIndex: 'state',
				renderer: rdr_status
			});
		}

		if (this.opt_show_ticket) {
			this.columns.push({
				header: 'Ticket',
				sortable: this.opt_column_sortable,
				flex: 1,
				dataIndex: 'state_type',
				renderer: rdr_ticket
			});
		}

		if (this.opt_show_last_check) {
			this.columns.push({
				header: _('Last check'),
				sortable: this.opt_column_sortable,
				flex: 1,
				dataIndex: 'timestamp',
				renderer: rdr_tstodate
			});
		}

		if (this.opt_show_component) {
			this.columns.push({
				header: _('Component'),
				flex: 1,
				sortable: this.opt_column_sortable,
				dataIndex: 'component'
			});
		}

		if (this.opt_show_resource) {
			this.columns.push({
				header: _('Resource'),
				flex: 1,
				sortable: this.opt_column_sortable,
				dataIndex: 'resource'
			});
		}

		if (this.opt_show_output) {
			this.columns.push({
				header: _('Message'),
				flex: 4,
				sortable: this.opt_column_sortable,
				dataIndex: 'output'
			});
		}

		if (this.opt_show_tags) {
			this.columns.push({
				header: _('Tags'),
				flex: 4,
				sortable: this.opt_column_sortable,
				dataIndex: 'tags',
				renderer: rdr_tags
			});
		}


		//store
		if (! this.store) {
			this.store = Ext.create('canopsis.lib.store.cstore', {
				//extend: 'canopsis.lib.store.cstore',
				model: 'canopsis.model.Event',

				pageSize: this.pageSize,

				sorters: this.sorters,

				remoteSort: this.remoteSort,

				proxy: {
					type: 'rest',
					url: '/rest/' + this.namespace + '/event',
					reader: {
						type: 'json',
						root: 'data',
						totalProperty: 'total',
						successProperty: 'success'
					}
				}
			});

			if (this.filter) {
				this.store.setFilter(this.filter);
			}

			if (this.autoload) {
				this.store.load();
			}
		}

		this.viewConfig = {
			stripeRows: false
		};

		if (this.opt_show_row_background) {
			this.viewConfig.getRowClass = this.coloringRow;
		}

		// Look at the eventlog initComponent duplicate
		//-----------filter button-------------------------

		// if (this.fitler_buttons) {
		// 		this.bar_search = [{
		// 		xtype: 'button',
		// 		iconCls: 'icon-crecord_type-resource',
		// 		pack: 'end',
		// 		tooltip: _('Show resource'),
		// 		enableToggle: true,
		// 		pressed: true,
		// 		scope: this,
		// 		toggleHandler: function(button, state) {
		// 			if (!state) {
		// 				button.filter_id = this.store.addFilter(
		// 					{'source_type': {'$ne': 'resource'}}
		// 				);
		// 			}else {
		// 				if (button.filter_id)
		// 					this.store.deleteFilter(button.filter_id);
		// 			}
		// 			this.store.load();
		// 		}
		// 	},{
		// 		xtype: 'button',
		// 		iconCls: 'icon-crecord_type-component',
		// 		pack: 'end',
		// 		tooltip: _('Show component'),
		// 		enableToggle: true,
		// 		pressed: true,
		// 		scope: this,
		// 		toggleHandler: function(button, state) {
		// 			if (!state) {
		// 				button.filter_id = this.store.addFilter(
		// 					{'source_type': {'$ne': 'component'}}
		// 				);
		// 			}else {
		// 				if (button.filter_id)
		// 					this.store.deleteFilter(button.filter_id);
		// 			}
		// 			this.store.load();
		// 		}
		// 	},{
		// 		xtype: 'button',
		// 		iconCls: 'icon-state-0',
		// 		pack: 'end',
		// 		tooltip: _('Show state ok'),
		// 		enableToggle: true,
		// 		pressed: true,
		// 		scope: this,
		// 		toggleHandler: function(button, state) {
		// 			if (!state) {
		// 				button.filter_id = this.store.addFilter(
		// 					{'state': {'$ne': 0}}
		// 				);
		// 			}else {
		// 				if (button.filter_id)
		// 					this.store.deleteFilter(button.filter_id);
		// 			}
		// 			this.store.load();
		// 		}
		// 	},{
		// 		xtype: 'button',
		// 		iconCls: 'icon-state-1',
		// 		pack: 'end',
		// 		tooltip: _('Show state warning'),
		// 		enableToggle: true,
		// 		pressed: true,
		// 		scope: this,
		// 		toggleHandler: function(button, state) {
		// 			if (!state) {
		// 				button.filter_id = this.store.addFilter(
		// 					{'state': {'$ne': 1}}
		// 				);
		// 			}else {
		// 				if (button.filter_id)
		// 					this.store.deleteFilter(button.filter_id);
		// 			}
		// 			this.store.load();
		// 		}
		// 	},{
		// 		xtype: 'button',
		// 		iconCls: 'icon-state-2',
		// 		pack: 'end',
		// 		tooltip: _('Show state critical'),
		// 		enableToggle: true,
		// 		pressed: true,
		// 		scope: this,
		// 		toggleHandler: function(button, state) {
		// 			if (!state) {
		// 				button.filter_id = this.store.addFilter(
		// 					{'state': {'$ne': 2}}
		// 				);
		// 			}else {
		// 				if (button.filter_id)
		// 					this.store.deleteFilter(button.filter_id);
		// 			}
		// 			this.store.load();
		// 		}
		// 	},{
		// 		xtype: 'button',
		// 		iconCls: 'icon-state-3',
		// 		pack: 'end',
		// 		tooltip: _('Show state unknown'),
		// 		enableToggle: true,
		// 		pressed: true,
		// 		scope: this,
		// 		toggleHandler: function(button, state) {
		// 			if (!state) {
		// 				button.filter_id = this.store.addFilter(
		// 					{'state': {'$ne': 3}}
		// 				);
		// 			}else {
		// 				if (button.filter_id)
		// 					this.store.deleteFilter(button.filter_id);
		// 			}
		// 			this.store.load();
		// 		}
		// 	},{
		// 		xtype: 'button',
		// 		iconCls: 'icon-state-type-0',
		// 		pack: 'end',
		// 		tooltip: _('Show soft state'),
		// 		enableToggle: true,
		// 		pressed: true,
		// 		scope: this,
		// 		toggleHandler: function(button, state) {
		// 			if (!state) {
		// 				button.filter_id = this.store.addFilter(
		// 					{'state_type': {'$ne': 0}}
		// 				);
		// 			}else {
		// 				if (button.filter_id)
		// 					this.store.deleteFilter(button.filter_id);
		// 			}
		// 			this.store.load();
		// 		}
		// 	},{
		// 		xtype: 'button',
		// 		iconCls: 'icon-state-type-1',
		// 		pack: 'end',
		// 		tooltip: _('Show hard state'),
		// 		enableToggle: true,
		// 		pressed: true,
		// 		scope: this,
		// 		toggleHandler: function(button, state) {
		// 			if (!state) {
		// 				button.filter_id = this.store.addFilter(
		// 					{'state_type': {'$ne': 1}}
		// 				);
		// 			}else {
		// 				if (button.filter_id)
		// 					this.store.deleteFilter(button.filter_id);
		// 			}
		// 			this.store.load();
		// 		}
		// 	},'-'];
		// }

		this.callParent(arguments);
	},

	coloringRow: function(record,index,rowParams,store) {
		state = record.get('state');
		if (state == 0) {
			return 'row-background-ok';
		} else if (state == 1) {
			return 'row-background-warning';
		} else if (state == 2) {
			return 'row-background-error';
		} else {
			return 'row-background-unknown';
		}
	}

});
