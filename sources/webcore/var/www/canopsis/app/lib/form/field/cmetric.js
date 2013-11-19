//need:app/lib/form/cfield.js,app/lib/store/cstore.js,app/lib/view/cgrid.js,app/lib/controller/cgrid.js
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

Ext.define('canopsis.lib.form.field.cmetric', {
	extend: 'Ext.panel.Panel',
	mixins: ['canopsis.lib.form.cfield'],

	requires: [
		'canopsis.lib.store.cstore',
		'canopsis.lib.view.cgrid',
		'canopsis.lib.controller.cgrid'
	],

	alias: 'widget.cmetric',

	logAuthor: '[form][field][cmetric]',

	border: false,
	layout: {
		type: 'vbox',
		align: 'stretch'
	},

	show_internals: false,

	multiSelect: true,

	sharedStore : undefined,

	initComponent: function() {
		this.callParent(arguments);

		log.debug('Initialize ...', this.logAuthor);

		this.extra_field = [];

		this.build_stores();
		this.build_grids();

		this.on('afterrender', function() {
			this.bind_event();
		}, this);

		var config = {
			layout: {
				type: 'hbox',
				align: 'stretch'
			},
			flex: 2
		};

		var container = Ext.create('Ext.container.Container', config);

		container.add(this.meta_grid);

		this.items.add(container);
		this.items.add(this.selected_grid);
	},

	afterRender: function(){
		this.callParent(arguments);

		if(this.sharedStore) {
			this.parentWizard = this.findParentByType('cwizard');
			this.parentWizard.childStores[this.sharedStore] = this.selected_store;
		}
	},

	build_stores: function() {
		log.debug('Build stores', this.logAuthor);

		// create model
		var fields = [
			{name: '_id'},
			{name: 'id', mapping: '_id'},
			{name: 'co'},
			{name: 're', defaultValue: undefined},
			{name: 'me'},
			{name: 't'},
			{name: 'u'},

			// optional configuration for customize metrics tab in widget's wizard

			{name: 'label', defaultValue: undefined},
			{name: 'curve_color', defaultValue: undefined},

			// widget: bar_graph
			{name: 'trend', defaultValue: undefined},

			// widget: diagram
			{name: 'category', defaultValue: undefined},

			// widget: gauge
			{name: 'mi', defaultValue: undefined},
			{name: 'tw', defaultValue: undefined},
			{name: 'tc', defaultValue: undefined},
			{name: 'ma', defaultValue: undefined},

			// widget: line_graph
			{name: 'curve_type', defaultValue: undefined},
			{name: 'area_color', defaultValue: undefined},
			{name: 'trend_curve', defaultValue: undefined},
			{name: 'yAxis', defaultValue: undefined},

			{name: 'threshold_warn', defaultValue: undefined},
			{name: 'threshold_crit', defaultValue: undefined},

			// widget: mini_chart
			{name: 'printed_value', defaultValue: undefined},
			{name: 'display_pct', defaultValue: undefined},

			// widget: trends
			{name: 'show_sparkline', defaultValue: undefined},
			{name: 'chart_type', defaultValue: undefined}
		];


		if(this.additional_field) {
			for(var i = 0; i < this.additional_field.length; i++) {
				fields.push({name: this.additional_field[i]});
			}
		}

		Ext.define('Meta', {
			extend: 'Ext.data.Model',
			fields: fields
		});

		// store
		this.meta_store = Ext.create('canopsis.lib.store.cstore', {
			model: 'Meta',
			remoteSort: true,
			sorters: [{
				property: 'co',
				direction: 'ASC'
			},{
				property: 're',
				direction: 'ASC'
			}],
			proxy: {
				type: 'ajax',
				url: '/perfstore/get_all_metrics',
				extraParams: {
					'show_internals': this.show_internals
				},
				reader: {
					type: 'json',
					root: 'data'
				}
			},
			autoLoad: true
		});

		if(this.sharedStore) {
			this.selected_store = Ext.create('canopsis.lib.store.cstore',
				Ext.Object.merge({model: 'Meta'}, this.sharedStore)
			);
		}
		else {
			this.selected_store = Ext.create('canopsis.lib.store.cstore', {
				model: 'Meta'
			});
		}
	},

	build_grids: function() {
		log.debug('Build grids', this.logAuthor);

		var bar_search = [{
			xtype: 'button',
			iconCls: 'icon-internal-metrics',
			pack: 'end',
			tooltip: _('Display internal metrics'),
			enableToggle: true,
			toggleHandler: function(button, state) {
				void(button);

				this.show_internals = state;
				this.meta_store.getProxy().extraParams.show_internals = this.show_internals;
				this.meta_store.load();
			},
			scope: this
		}];

		// first grid
		this.meta_grid = Ext.create('canopsis.lib.view.cgrid', {
			store: this.meta_store,
			flex: 2,
			margin: 3,

			opt_menu_rights: false,
			opt_bar: true,
			opt_bar_search: true,
			opt_bar_add: false,
			opt_allow_edit: false,
			opt_bar_duplicate: false,
			opt_bar_reload: true,
			opt_bar_delete: false,
			opt_multiSelect: this.multiSelect,
			opt_paging: true,
			opt_simple_search: true,

			bar_search: bar_search,

			border: true,

			columns: [
				{
					header: _('Component'),
					sortable: false,
					dataIndex: 'co',
					flex: 1
				},{
					header: _('Resource'),
					sortable: false,
					dataIndex: 're',
					flex: 1
				},{
					header: _('Metric'),
					sortable: false,
					dataIndex: 'me',
					flex: 1
				}
			],
			viewConfig: {
				copy: true,
				plugins: {
					ptype: 'gridviewdragdrop',
					enableDrop: false,
					dragGroup: 'search_grid_DNDGroup'
				}
			}
		});

		// Create controller and bind with meta_grid
		this.meta_grid_ctrl = Ext.create('canopsis.lib.controller.cgrid');
		this.meta_grid.on('afterrender', function() {
			this.meta_grid_ctrl._bindGridEvents(this.meta_grid);
		}, this);


		// Selection grid
		var _columns = [
			{
				xtype: 'actioncolumn',
				width: 25,
				align: 'center',
				tooltip: _('Delete'),
				icon: './themes/canopsis/resources/images/icons/bin_closed.png',
				handler: function(view, rowIndex) {
					view.getStore().removeAt(rowIndex);
				}
			},{
				header: _('Component'),
				sortable: false,
				dataIndex: 'co',
				flex: 1
			},{
				header: _('Resource'),
				sortable: false,
				dataIndex: 're',
				flex: 1
			},{
				header: _('Metric'),
				sortable: false,
				dataIndex: 'me',
				flex: 1
			}
		];

		//additionnal columns
		var _plugins = [];

		//create grid
		this.selected_grid = Ext.widget('grid', {
			store: this.selected_store,
			flex: 1,
			margin: 3,
			multiSelect: true,
			border: true,
			scroll: true,
			columns: {
				defaults: {
					menuDisabled: true,
					sortable: false
				},
				items: _columns
			},
			plugins: _plugins,
			viewConfig: {
				markDirty: false,
				plugins: {
					ptype: 'gridviewdragdrop',
					copy: false,
					dragGroup: 'search_grid_DNDGroup',
					dropGroup: 'search_grid_DNDGroup'
				}
			}
		});

		// build menu
		this.clearAllButton = Ext.create('Ext.Action', {
			iconCls: 'icon-delete',
			text: _('Clear all'),
			action: 'clear'
		});

		this.deleteButton = Ext.create('Ext.Action', {
			iconCls: 'icon-delete',
			text: _('Delete selected'),
			action: 'delete'
		});

		this.contextMenu = Ext.create('Ext.menu.Menu', {
			items: [this.deleteButton, this.clearAllButton]
		});
	},

	bind_event: function() {
		log.debug('Binding events', this.logAuthor);

		// Meta inventory
		this.meta_grid.on('itemdblclick', function(view, record) {
			void(view);

			this.select_meta(record);
		}, this);

		// drop function

		this.selected_grid.getView().on('beforedrop', function(html_node, data) {
			void(html_node);

			//only do action if is not reorder
			if(data.view.id !== this.selected_grid.getView().id) {
				var records = data.records;

				for(var i = 0; i < records.length; i++) {
					this.select_meta(records[i]);
				}

				return false;
			}
		}, this);

		// Menu option
		this.selected_grid.on('itemcontextmenu', this.open_menu, this);

		this.clearAllButton.setHandler(function() {
			this.selected_store.removeAll();
		}, this);

		this.deleteButton.setHandler(this.deleteSelected, this);
	},

	fetch_metrics: function(record) {
		log.debug('Fetch metrics', this.logAuthor);

		var metric_array = [];
		var metrics = record.get('metrics');
		var node = record.get('node');
		var dn = record.get('dn');

		for(var i = 0; i < metrics.length; i++) {
			metric_array.push({
				'node': node,
				'metric': metrics[i].dn,
				'dn': dn
			});
		}

		return metric_array.sort(this.sort_by_metric);
	},

	sort_by_metric: function(a, b) {
		a = a.metric;
		b = b.metric;

		if(a === b) {
			return 0;
		}
		if (a > b) {
			return 1;
		}
		else {
			return -1;
		}
	},

	select_meta: function(record) {
		var _id = record.get('_id');
		log.debug('Select Meta ' + _id, this.logAuthor);

		if(!this.selected_store.getById(_id)) {
			if(!this.multiSelect) {
				this.selected_store.removeAll();
			}

			this.selected_store.add(record.copy());
		}
		else {
			log.debug(' + Already selected' , this.logAuthor);
		}
	},

	open_menu: function(view, rec, node, index, e) {
		void(node, index);

		e.preventDefault();
		//don't auto select if multi selecting
		var selection = this.selected_grid.getSelectionModel().getSelection();

		if (selection.length < 2) {
			view.select(rec);
		}

		this.contextMenu.showAt(e.getXY());
		return false;
	},

	deleteSelected: function() {
		log.debug('delete selected metrics', this.logAuthor);
		var selection = this.selected_grid.getSelectionModel().getSelection();

		for(var i = 0; i < selection.length; i++) {
			this.selected_store.remove(selection[i]);
		}
	},

	getValue: function() {
		log.debug('Write values', this.logAuthor);

		var output = [];
		var order = 0;

		this.selected_store.each(function(record) {
			record.data.order = order++;
			var data = Ext.clone(record.data);

			//clean that
			if(data.me) {
				data.metrics = [data.me];
			}

			output.push(data);
		}, this);

		return output;
	},

	setValue: function(data) {
		log.debug('Load values', this.logAuthor);
		log.dump(data);

		var metricList = [];

		//retrocompatibility
		if(Ext.isArray(data)) {
			for(var i = 0; i < data.length; i++) {
				var item = Ext.clone(data[i]);

				if(!item.co) {
					item.co = item.component;
				}

				if(!item.re) {
					item.re = item.resource;
				}

				if(!item.t) {
					item.t = item.type;
				}

				if(!item.me) {
					item.me = item.metrics[0];
				}

				var config = {
					id: item.id,
					co: item.co,
					re: item.re,
					me: item.me,
					t: item.type
				};

				metricList.push(Ext.create('Meta', config));
			}
		}

		if(Ext.isObject(data)){
			Ext.Object.each(data, function(key, value) {
				void(key);

				var item = Ext.clone(value);

				if(!item.co) {
					item.co = item.component;
				}

				if(!item.re) {
					item.re = item.resource;
				}

				if(!item.t) {
					item.t = item.type;
				}

				if(!item.me) {
					item.me = item.metrics[0];
				}

				// TODO: check if next line is important
				if(Ext.isArray) {
					metricList.push(Ext.create('Meta', item));
				}
			}, this);
		}

		if(metricList.length > 0 && data[metricList[0].get('id')]['order'] !== undefined) {
			metricList.sort(function(a, b) {
				return data[a.get('id')]['order'] - data[b.get('id')]['order'];
			});
		}

		this.selected_store.add(metricList);
	},

	beforeDestroy: function(){
		//deference store
		if(this.sharedStore && this.parentWizard) {
			delete this.parentWizard.childStores[this.sharedStore];
		}
	}
});
