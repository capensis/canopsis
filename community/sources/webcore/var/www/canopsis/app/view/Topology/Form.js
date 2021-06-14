//need:app/lib/view/cform.js,app/lib/form/field/ctopo.js,app/lib/form/field/cinventory.js
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
Ext.define('canopsis.view.Topology.Form', {
	extend: 'canopsis.lib.view.cform',
	alias: 'widget.TopologyForm',

	requires: [
		'canopsis.lib.form.field.ctopo',
		'canopsis.lib.form.field.cinventory'
	],

	layout: {
		type: 'hbox',
		align: 'stretch'
	},

	initComponent: function() {
		var me = this;
		this.rightPanel = Ext.create('canopsis.lib.form.field.ctopo', {
			flex: 1,
			margin: '0 0 5 0'
		});

		this.checksGrid = Ext.create('canopsis.lib.form.field.cinventory', {
			select: false,
			isFormField: false,
			padding: 0,
			search_grid_border: false,
			height: '100%',
			width: '100%',
			base_filter: {event_type: 'check'}
		});

		this.selsGrid = Ext.create('canopsis.lib.form.field.cinventory', {
			select: false,
			isFormField: false,
			padding: 0,
			search_grid_border: false,
			height: '100%',
			width: '100%',
			showResource: false,
			base_filter: {event_type: 'selector'}
		});

		this.toposGrid = Ext.create('canopsis.lib.form.field.cinventory', {
			select: false,
			isFormField: false,
			padding: 0,
			search_grid_border: false,
			height: '100%',
			width: '100%',
			showResource: false,
			base_filter: {event_type: 'topology'}
		});

		this.optionsForm = Ext.create('Ext.form.Panel', {
			border: false,
			defaultType: 'textfield',
			layout: 'anchor',
			defaults: {
				anchor: '80%',
				margin: '5 5 0 5'
			},

			items: [{
				fieldLabel: _('Name'),
				name: 'crecord_name',
				listeners: {
					change: function(value) {
						me.rightPanel.setRootNodeName(value.value);
					}
				},
				allowBlank: false
			},{
				fieldLabel: _('Display name'),
				name: 'display_name'
			},{
				xtype: 'textarea',
				fieldLabel: _('Description'),
				name: 'description',
				listeners: {
					change: function(value) {
						me.rightPanel.setRootNodeDescription(value.value);
					}
				}
			}]
		});

		this.leftPanel = Ext.widget('tabpanel', {
			title: _('Toolbox'),
			width: 400,
			margin: '0 5 5 0',
			border: true,
			collapsible: true,
			collapseDirection: 'left',
			animCollapse: false,
			items: [{
				title: _('Options'),
				border: false,
				layout: 'fit',
				items: [this.optionsForm]
			},{
				title: _('Checks'),
				border: false,
				layout: 'fit',
				items: [this.checksGrid]
			},{
				title: _('Selectors'),
				border: false,
				layout: 'fit',
				items: [this.selsGrid]
			},{
				title: _('Topologies'),
				border: false,
				layout: 'fit',
				items: [this.toposGrid]
			},{
				title: _('Operators'),
				border: false,
				layout: 'fit',
				items: [this.build_object_grid()]
			}]
		});

		this.items = [this.leftPanel, this.rightPanel];

		this.callParent();
	},

	build_object_grid: function() {
		var store = Ext.create('Ext.data.Store', {
			fields: [
				{name: '_id'},
				{name: 'component'},
				{name: 'source_type'},
				{name: 'event_type'},
				{name: 'description', defaultValue: ''},
				{name: 'nodeMaxOutConnexion', defaultValue: 1},
				{name: 'nodeMaxInConnexion', defaultValue: 2},
				{name: 'form', defaultValue: undefined}

			],
			proxy: {
				type: 'ajax',
				url: '/topology/getOperators',
				reader: {
					type: 'json',
					root: 'data'
				}
			},
			autoLoad: true
		});

		// Translate Operators
		store.on('load', function(store, records) {
			void(store);

			for(var i = 0; i < records.length; i++) {
				var record = records[i];
				record.set('description', _(record.get('description')));
			}
		}, this, {single: true});

		this.object_grid = Ext.create('Ext.grid.Panel', {
			store: store,
			border: false,
			viewConfig: {
				markDirty: false,
				plugins: {
					ptype: 'gridviewdragdrop',
					enableDrop: false,
					dragGroup: 'search_grid_DNDGroup',
					dropGroup: 'search_grid_DNDGroup'
				}
			},
			columns: [{
				header: _('Id'),
				width: 100,
				dataIndex: 'component'
			},{
				header: _('Description'),
				flex: 1,
				dataIndex: 'description'
			}]
		});

		return this.object_grid;
	}
});
