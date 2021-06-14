//need:app/lib/view/cgrid_state.js,app/store/EventLogs.js
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
Ext.define('canopsis.view.Event.Log' , {
	extend: 'canopsis.lib.view.cgrid_state',

	alias: 'widget.EventLog',

	logAuthor: '[View][Event][Log]',

	requires: [
		'canopsis.lib.controller.cgrid',
		'canopsis.store.EventLogs'
	],

	store: undefined,

	opt_show_component: true,
	opt_allow_edit: false,
	opt_bar_delete: false,
	opt_bar_add: false,
	opt_bar_time_search: true,
	opt_bar_search: true,

	opt_paging: true,
	opt_bar: true,

	opt_bar_search_field: ['component', 'resource'],

	initComponent: function() {
		this.bar_search = [{
			xtype: 'button',
			iconCls: 'icon-crecord_type-resource',
			pack: 'end',
			tooltip: _('Show resource'),
			enableToggle: true,
			pressed: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(!state) {
					button.filter_id = this.store.addFilter(
						{'source_type': {'$ne': 'resource'}}
					);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-crecord_type-component',
			pack: 'end',
			tooltip: _('Show component'),
			enableToggle: true,
			pressed: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(!state) {
					button.filter_id = this.store.addFilter(
						{'source_type': {'$ne': 'component'}}
					);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-state-0',
			pack: 'end',
			tooltip: _('Show state ok'),
			enableToggle: true,
			pressed: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(!state) {
					button.filter_id = this.store.addFilter(
						{'state': {'$ne': 0}}
					);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-state-1',
			pack: 'end',
			tooltip: _('Show state warning'),
			enableToggle: true,
			pressed: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(!state) {
					button.filter_id = this.store.addFilter(
						{'state': {'$ne': 1}}
					);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-state-2',
			pack: 'end',
			tooltip: _('Show state critical'),
			enableToggle: true,
			pressed: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(!state) {
					button.filter_id = this.store.addFilter(
						{'state': {'$ne': 2}}
					);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-state-3',
			pack: 'end',
			tooltip: _('Show state unknown'),
			enableToggle: true,
			pressed: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(!state) {
					button.filter_id = this.store.addFilter(
						{'state': {'$ne': 3}}
					);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-state-type-0',
			pack: 'end',
			tooltip: _('Show soft state'),
			enableToggle: true,
			pressed: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(!state) {
					button.filter_id = this.store.addFilter(
						{'state_type': {'$ne': 0}}
					);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-state-type-1',
			pack: 'end',
			tooltip: _('Show hard state'),
			enableToggle: true,
			pressed: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(!state) {
					button.filter_id = this.store.addFilter(
						{'state_type': {'$ne': 1}}
					);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},'-'],

		// bind controller
		this.ctrl = Ext.create('canopsis.lib.controller.cgrid');

		var stop = (new Date().getTime()) / 1000;
		var start = stop - 86400;

		// add filter for timestamp
		var initialFilter = {'$and': [
			{
				timestamp: {
					'$gt': parseInt(start)
				}
			},{
				timestamp: {
					'$lt': parseInt(stop)
				}
			}
		]};

		this.store = Ext.create('canopsis.store.EventLogs');
		this.ctrl.filter_id = this.store.addFilter(initialFilter);

		this.callParent(arguments);

		this.on('afterrender', function() {
			this.ctrl._bindGridEvents(this);
		}, this);
	}
});
