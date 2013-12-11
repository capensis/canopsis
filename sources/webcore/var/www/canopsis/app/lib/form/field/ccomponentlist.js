//need:app/lib/form/cfield.js,app/lib/view/cgrid.js,app/lib/store/cstore.js,app/lib/controller/cgrid.js,app/lib/menu/cclear.js
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

Ext.define('canopsis.lib.form.field.ccomponentlist' , {
	extend: 'Ext.panel.Panel',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.ccomponentlist',

	requires: [
		'canopsis.lib.view.cgrid',
		'canopsis.lib.store.cstore',
		'canopsis.lib.controller.cgrid',
		'canopsis.lib.menu.cclear'
	],

	border: false,
	search_grid_border: true,

	metrics: true,

	select: true,
	multiSelect: true,

	dropGroup: 'search_grid_DNDGroup',
	dragGroup: 'search_grid_DNDGroup',

	inventory_url: '/rest/events/event',

	vertical_multiselect: false,
	padding: 5,
	base_filter: undefined,

	showResource: true,

	loaded_value: [],

	idKey: "component",

	fields: ["component", "ressource"],

	//switch variable for mask, because asynchronous is hard to handle
	value_already_fetched: false,

	getName: function() {
		return this.name;
	},

	isValid: function() {
		return true;
	},

	validate: function() {
		return this.isValid();
	},

	getSubmitData: function() {
		var data = {};
		data[this.name] = this.getValue();
		return data;
	},

	initComponent: function() {
		this.logAuthor = '[' + this.id + ']';
		log.debug('Initialize ...', this.logAuthor);

		var default_layout = {
			type: 'hbox',
			align: 'stretch'
		};

		var default_defaults = {
			padding: this.padding
		};

		this.padding = 0;

		this.items = this.buildView();

		if(!this.multiSelect || this.vertical_multiselect) {
			default_layout.type = 'vbox';
		}

		//HACK: for test in widget, wrap in container
		if(this.height) {
			log.debug(' + Widget mode', this.logAuthor);
			this.items = {xtype: 'container', items: this.items, layout: default_layout, defaults: default_defaults};
			this.layout = 'fit';
		}
		else {
			this.layout = default_layout;
			this.defaults = default_defaults;
		}

		this.callParent(arguments);
	},

	buildView: function() {
		log.debug(' + Build view ...', this.logAuthor);

		var me = this;
		var items = [];

		// building model on the fly with the additionnal field if needed
		var fields = this.fields;

		if(this.additional_field) {
			fields.push(this.additional_field.name);
		}

		Ext.define('simplified_event', {
			extend: 'Ext.data.Model',
			fields: fields,
			idProperty: this.idKey
		});

		var model = Ext.ModelManager.getModel('simplified_event');

		this.columns = [
			{
				header: '',
				width: 25,
				sortable: false,
				dataIndex: 'source_type',
				renderer: rdr_source_type
			},{
				header: _('Component'),
				flex: 1,
				dataIndex: 'component'
			}
		];

		if(this.showResource) {
			this.columns.push({
					header: _('Resource'),
					flex: 2,
					dataIndex: 'resource'
			});
		}


		// Selection GRID

		if(this.select) {
			log.debug(' + Selection grid', this.logAuthor);

			this.selection_render = function(value, p, record) {
				void(value, p);

				var node = '';

				if(record.data.resource) {
					node = Ext.String.format('<b>{0}</b><br>&nbsp;&nbsp;{1}', record.data.component, record.data.resource);
				}
				else {
					node = Ext.String.format('<b>{0}</b>', record.data.component);
				}

				return node;
			};

			this.selection_store = Ext.create('Ext.data.Store', {
				model: model
			});

			var selection_height = undefined;

			if(!this.multiSelect) {
				selection_height = 130;
			}

			var selection_grid_config = {
				title: _('Selection'),
				border: true,
				multiSelect: this.multiSelect,
				opt_bar: false,
				opt_allow_edit: false,
				opt_paging: false,
				flex: 1,
				height: selection_height,
				store: this.selection_store,
				hideHeaders: true,
				autoScroll: true,
				columns: [
					{
						header: '',
						width: 25,
						sortable: false,
						dataIndex: 'source_type',
						renderer: rdr_source_type
					},{
						sortable: false,
						dataIndex: this.idKey,
						flex: 2,
						renderer: this.selection_render
					}
				],

				viewConfig: {
					markDirty: false,
					plugins: {
						ptype: 'gridviewdragdrop',
						dragGroup: this.dragGroup,
						dropGroup: this.dropGroup
					}
				}
			};

			// additional field (if specified only)
			if(this.additional_field) {
				selection_grid_config.plugins = [
					Ext.create('Ext.grid.plugin.CellEditing', {
						clicksToEdit: 1,
						autoCancel: true,
						listeners: {
							beforeedit: function( editor, e, eOpts )
							{
								//Dirty hack to make ccolorfield work with cellediting
								this.record = e.record;
							}
						}
					})
				];

				var editor_config = {
					sortable: false,
					dataIndex: this.additional_field.name,
					editor: this.additional_field,
					maskOnDisable: false,
					flex: 3
				};

				selection_grid_config.emptyText = this.additional_field.emptyText;

				if(this.additional_field.name === 'link' || this.additional_field.name === 'display_name') {
					editor_config.renderer = function(val) {
						if(!val) {
							return Ext.String.format('<span style="color:grey">{0}</span>', this.emptyText);
						}
						else {
							return val;
						}
					};
				}

				selection_grid_config.columns.push(editor_config);
				selection_grid_config.flex = 2;
			}

			this.selection_grid = Ext.create('canopsis.lib.view.cgrid', selection_grid_config);
		}

		// Search GRID
		log.debug(' + Search grid', this.logAuthor);

		this.search_store = Ext.create('canopsis.lib.store.cstore', {
			model: model,
			pageSize: global.pageSize,
			proxy: {
				type: 'rest',
				url: this.inventory_url,
				reader: {
					type: 'json',
					root: 'data',
					totalProperty: 'total',
					successProperty: 'success'
				}
			},

			autoLoad: false
		});

		//set base filter if given
		if(this.base_filter !== undefined) {
			this.search_store.addFilter(this.base_filter);
		}

		this.search_store.load();

		this.search_grid = Ext.create('canopsis.lib.view.cgrid', {
			multiSelect: this.multiSelect,
			opt_bar: true,
			opt_bar_search: true,
			opt_bar_add: false,
			opt_allow_edit: false,
			opt_bar_duplicate: false,
			opt_bar_reload: true,
			opt_bar_delete: false,
			opt_bar_search_field: [this.idKey],
			border: this.search_grid_border,
			opt_paging: true,
			flex: 2,
			store: this.search_store,
			columns: this.columns,
			viewConfig: {
				copy: true,
				plugins: {
					ptype: 'gridviewdragdrop',
					enableDrop: false,
					dragGroup: this.dragGroup
				}
			}
		});

		// Bind cgrid controller on search grid
		this.search_ctrl = Ext.create('canopsis.lib.controller.cgrid');

		this.on('afterrender', function() {
			this.search_ctrl._bindGridEvents(this.search_grid);
		}, this);

		// Bind events
		log.debug(' + Bind events', this.logAuthor);

		if(this.selection_grid) {
			this.contextMenu = Ext.create('canopsis.lib.menu.cclear', {
				grid: this.selection_grid
			});

			this.selection_grid.on('itemdblclick', function(grid, record, item, index) {
				void(grid, item);

				this.selection_store.removeAt(index);
			}, this);

			this.selection_grid.getView().on('beforedrop', function(event, data) {
				var records = data.records;

				if(data.view.id !== this.selection_grid.getView().id) {
					for(var i = 0; i < records.length; i++) {
						this.addRecord(records[i]);
					}

					event.cancel = true;
					event.dropStatus = true;
					return false;
				}
			}, this);
		}

		this.search_grid.on('itemdblclick', function(grid, record) {
			void(grid);

			this.addRecord(record);
		}, this);

		// Push items
		log.debug(' + Set items', this.logAuthor);
		items.push(this.search_grid);

		var rightPanel = Ext.create('Ext.Panel', {
			align: 'right',
			layout: {
				type: 'vbox',
				align: 'stretch'
			},
			flex: 2
    	});

		var addButton = Ext.create('Ext.Button', {
			text: 'Add',
			handler: function() {
				var formRecord = me.getAddFormValues();

        		me.addRecord(formRecord);
    		}
		});

		if(this.addForm)
		{
			this.addForm = Ext.create("canopsis.lib.form.field.cfieldset", {
				items: this.addForm,
				itemId: "fieldsetAddForm"
			});

			this.addForm.add(addButton);
			rightPanel.add(this.addForm);
		}

		rightPanel.add(this.selection_grid);

		if(this.selection_grid) {
			items.push(rightPanel);
		}

		return items;

		console.log("AFTERREENNNNNNDERRRRR");
		if(this.additional_field.name === 'color') {
			//Dirty hack to make ccolofield work with cellediting
			console.log("this.additional_field");
			console.log(this.additional_field);
			this.additional_field.menuListeners.select = function() {
				this.setValue(d);
				var rec =  this.ownerCt.editingPlugin.record;
				rec.set('color', this.getValue());
				this.ownerCt.editingPlugin.ccomponentlist.addRecord(rec);
			}
		}
	},

	getAddFormValues: function() {
		var result = { data : {} };

		var formItems = this.addForm.items;

		for (var i = formItems.items.length - 1; i >= 0; i--) {
			var formItem = formItems.items[i];
			if(formItem.getValue !== undefined)
			{
				result.data[formItem.name] = formItem.getValue();
			}
		};

		return result;
	},

	addRecord: function(record, index) {
		if(this.selection_grid) {
			if(this.selection_store.findExact(this.idKey, record.data[this.idKey]) === -1) {
				var record_data = record.data;

				if(!this.multiSelect) {
					this.selection_store.removeAll();
				}

				//set additionnal value if needed
				if(this.additional_field && this.loaded_value[record.data[this.idKey]] !== undefined) {
					record_data[this.additional_field.name] = this.loaded_value[record.data[this.idKey]];
				}

				if(index !== undefined) {
					this.selection_store.insert(index, record_data);
				}
				else {
					this.selection_store.add(record_data);
				}
			}
			else {
				log.debug(record.data.id + ' already selected', this.logAuthor);
			}
		}
	},

	// GetValue for wizard ...
	getValue: function() {
		var dump = [];

		if(this.selection_grid) {
			this.selection_store.each(function(selectedElement) {
				var obj = Ext.clone(selectedElement.data);

				if(this.additional_field) {
					var additional_value = selectedElement.data[this.additional_field.name];
					obj[this.additional_field.name] = additional_value;
					dump.push(obj);
				}
				else {
					dump.push(obj);
				}
			}, this);
		}
		return dump;
	},

	setValue_record: function(records) {
		for(var i = 0; i < records.length; i++) {
			this.addRecord({
				data: records[i]
			});
		}
	},

	setValue: function(data) {
		if(this.selection_grid) {
			this.loaded_value = data;
			var ids = [];

			//Get only id
			if(this.additional_field && Ext.isObject(data[0])) {
				//push id
				for(var i = 0; i < data.length; i++) {
					ids.push(data[i].id);
				}

				//list to dict
				var dict = {};
				for (i = 0; i < this.loaded_value.length; i++) {
					dict[this.loaded_value[i][this.idKey]] = this.loaded_value[i][this.additional_field.name];
				}

				this.loaded_value = dict;
			}
			else {
				ids = data;
			}

			this.setValue_record(data);

			if(this.selection_grid.rendered && this.selection_grid.getEl().isMasked()) {
				this.loading_mask = this.selection_grid.getEl().unmask();
			}

			this.value_already_fetched = true;

		}
	},

	beforeDestroy: function() {
		this.search_ctrl.destroy();

		if(this.selection_grid) {
			this.selection_grid.destroy();
		}

		this.search_grid.destroy();

		Ext.grid.Panel.superclass.beforeDestroy.call(this);
		log.debug('Destroyed.', this.logAuthor);
	}
});
