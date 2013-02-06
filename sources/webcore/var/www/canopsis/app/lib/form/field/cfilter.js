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

Ext.define('canopsis.lib.form.field.cfilter' , {
	extend: 'Ext.panel.Panel',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.cfilter',

	border: false,

	url: undefined,
	namespace: 'events',
	ctype: 'event',
	model: 'canopsis.model.Event',
	autoScroll: true,
	params: {},

	columns: [
		{
			header: '',
			width: 25,
			sortable: false,
			dataIndex: 'source_type',
			renderer: rdr_source_type
   		},{
   			header: _("Component"),
			sortable: false,
			dataIndex: 'component',
			flex: 2
 		},{
 			header: _("Resource"),
			sortable: false,
			dataIndex: 'resource',
			flex: 2
 	}],

	filter: undefined,

	operator_fields: [
		{'operator': 'connector_name',	'text': _('Connector name'),	'type': 'all' },
		{'operator': 'event_type',	'text': _('Event type'),	'type': 'all'},
		{'operator': 'state',	'text': _('State'),	'type': 'all'},
		{'operator': 'state_type',	'text': _('State type'),	'type': 'all'},
		{'operator': 'resource',	'text': _('Resource'),	'type': 'all'},
		{'operator': 'component',	'text': _('Component'),	'type': 'all'},
		{'operator': 'tags', 'text': _('Tags'),	'type': 'all'}
	],

	layout: {
        type: 'vbox',
        align: 'stretch'
    },

    checkObjectValidity: true,

	initComponent: function() {
		this.logAuthor = '[' + this.id + ']';
		log.debug('Initialize ...', this.logAuthor);

		this.define_object();
		this.build_store();

		var url = this.url;

		if (! url) {
			url = '/rest/' + this.namespace;
			if (this.ctype)
				url += '/' + this.ctype;
		}

		//-----------------preview windows----------------
		this.preview_store = Ext.create('canopsis.lib.store.cstore', {
			model: this.model,
			proxy: {
				type: 'rest',
				url: url,
				extraParams: this.params,
				reader: {
					type: 'json',
					root: 'data',
					totalProperty: 'total',
					successProperty: 'success'
				}
			 },
			autoLoad: false
		});

		this.preview_grid = Ext.widget('grid', {
			store: this.preview_store,
			border: true,
			hidden: true,
			hideHeaders: false,
			columns: this.columns
		});

		//-------------cfilter (wizard part)---------------
		this.cfilter = Ext.create('cfilter.object', {
			operator_store: this.operator_store,
			sub_operator_store: this.sub_operator_store,
			opt_remove_button: false,
			start_with_and: false
		});

		//--------------edit area (hand writing part)--------

		this.edit_area = Ext.widget('textarea', {
			isFormField: false,
			hidden: true,
			validator: this.check_json_validity,
			flex: 1
		});

		//---------------------TBAR--------------------------
		this.wizard_button = Ext.widget('button', {handler: this.show_wizard,
			iconCls: 'icon-wizard',
			tooltip: _('Wizard'),
			scope: this,
			disabled: true,
			margin: 5
		});
		this.edit_button = Ext.widget('button', {
			handler: this.show_edit_area,
			tooltip: _('Edit'),
			iconCls: 'icon-edit',
			margin: 5,
			scope: this
		});
		this.preview_button = Ext.widget('button', {
			handler: this.show_preview,
			tooltip: _('Preview'),
			iconCls: 'icon-preview',
			margin: 5,
			scope: this
		});
		this.clean_button = Ext.widget('button', {
			handler: this.reset_cfilter,
			tooltip: _('Clean'),
			iconCls: 'icon-clean',
			margin: 5,
			scope: this
		});

		var button_panel = Ext.widget('panel', {
			border: false,
			items: [this.wizard_button, this.edit_button, this.preview_button, this.clean_button ]
		});

		this.items = [button_panel, this.cfilter, this.edit_area, this.preview_grid];
		this.callParent(arguments);
	},

	reset_cfilter: function(){
		this.cfilter.destroy()
		this.cfilter = Ext.create('cfilter.object', {
			operator_store: this.operator_store,
			sub_operator_store: this.sub_operator_store,
			opt_remove_button: false,
			start_with_and: false
		});
		this.add(this.cfilter)
	},

	check_json_validity: function(value) {
		if (value == '')
			return true;
		try {
			var obj = Ext.decode(value);
		}catch (err) {
			return 'Error: invalid JSON';
		}
		return true;
	},

	check_object_validity: function(obj) {
		if (obj && this.checkObjectValidity) {
			var output = true;
			for (var i = 0; i < obj.length; i++) {
				if (Ext.isArray(obj[i])) {
					if (obj[i].length == 0)
						return false;
					else
						output = this.check_object_validity(obj[i]);
				}else if (Ext.isObject(obj[i])) {
					output = this.check_object_validity(obj[i]);
				}
			}
			return output;
		}else {
			return true;
		}
	},

	isValid: function() {
		log.debug('Execute isValid function', this.logAuthor);
		var value = this.getRawValue();
		if (this.check_json_validity(value) == true) {
			if (this.check_object_validity(Ext.decode(value))) {
				return true;
			}else {
				global.notify.notify('Invalid filter', "You can't let an array empty (and / or / in ...)", 'warning');
				return false;
			}
		}else {
			return false;
		}
	},


	switch_elements_visibility: function(cfilter,edit_area,preview_grid) {
		(edit_area) ? this.edit_area.show() : this.edit_area.hide();
		(preview_grid) ? this.preview_grid.show() : this.preview_grid.hide();
		(cfilter) ? this.cfilter.show() : this.cfilter.hide();
	},

	switch_button_state: function(wizard,edit,preview) {
		(wizard) ? this.wizard_button.setDisabled(false) : this.wizard_button.setDisabled(true);
		(edit) ? this.edit_button.setDisabled(false) : this.edit_button.setDisabled(true);
		(preview) ? this.preview_button.setDisabled(false) : this.preview_button.setDisabled(true);

		(wizard) ? this.clean_button.setDisabled(true) : this.clean_button.setDisabled(false);
	},


	show_wizard: function() {
		if (!this.edit_area.isHidden()) {
			if (this.edit_area.validate()) {
				var filter = this.edit_area.getValue();
				filter = strip_return(filter);
				if (filter && filter != '') {
					this.cfilter.remove_all_cfilter();
					this.setValue(filter);
				}

				this.switch_elements_visibility(true, false, false);
				this.switch_button_state(false, true, true);
			}else {
				log.debug('Incorrect JSON given', this.logAuthor);
			}
		}else {
			this.switch_elements_visibility(true, false, false);
			this.switch_button_state(false, true, true);
		}
	},

	show_edit_area: function() {
		var filter = Ext.decode(this.getValue());
		if (filter) {
			filter = JSON.stringify(filter, undefined, 8);
			this.edit_area.setValue(filter);
		}
		this.switch_elements_visibility(false, true, false);
		this.switch_button_state(true, false, true);
	},

	show_preview: function() {
		var filter = this.getValue();

		if (filter) {
			if (this.check_object_validity(Ext.decode(filter))) {
				this.preview_store.clearFilter();
				log.debug('Showing preview with filter: ' + filter, this.logAuthor);
				this.preview_store.setFilter(filter);
				this.preview_store.load();

				this.switch_elements_visibility(false, false, true);
				this.switch_button_state(true, true, false);
			}else {
				global.notify.notify('Invalid filter', "You can't let an array empty (and / or / in ...)", 'warning');

			}
		}
	},

	define_object: function() {

		//for array input
		Ext.define('cfilter.array_field', {
			extend: 'Ext.panel.Panel',

			border: false,
			value: undefined,

			margin: '0 0 0 5',
			layout: 'hbox',

			initComponent: function() {
				this.textfield_panel = Ext.widget('panel', {
					border: false,
					margin: '0 0 0 5'
				});

				if (!this.value) {
					this.add_textfield();
				}
				//--------buttons--------
				this.add_button = Ext.widget('button', {
					iconCls: 'icon-add',
					//margin: '0 0 0 5',
					tooltip: _('Add new value to this list')
				});
				//--------build object----
				this.items = [this.add_button, this.textfield_panel];
				this.callParent(arguments);
				//--------bindings-------
				this.add_button.on('click', function() {this.add_textfield()},this);
			},

			add_textfield: function(value) {
				var config = {
					emptyText: _('Type value here'),
					isFormField: false
				};

				if (value)
					config.value = value;

				var textfield = Ext.widget('textfield', config);
				var remove_button = Ext.widget('button', {
					iconCls: 'icon-cancel',
					margin: '0 0 0 5',
					width: 24,
					tooltip: _('Remove this from list of value')
				});

				var item_array = [textfield];

				//if it's not first elem, add remove button
				if (this.textfield_panel.items.length >= 1)
					item_array.push(remove_button);

				var panel = Ext.widget('panel', {
					border: false,
					margin: '0 0 5 0',
					layout: 'hbox',
					items: item_array
				});
				remove_button.on('click', function(button) {button.up().destroy()});

				return this.textfield_panel.add(panel);
			},

			getValue: function() {
				var output = [];
				for (var i = 0; i < this.textfield_panel.items.items.length; i++) {
					var panel = this.textfield_panel.items.items[i];
					var textfield = panel.down('.textfield');
					output.push(textfield.getValue());
				}
				return output;
			},

			setValue: function(array) {
				this.textfield_panel.removeAll();
				for (var i = 0; i < array.length; i++)
					this.add_textfield(array[i]);
			}
		});


		//this object is made of two component, upper panel with combobox
		//and the bottom panel with object (itself) and add button
		Ext.define('cfilter.object' , {
			extend: 'Ext.panel.Panel',
			border: false,

			operator_store: undefined,
			sub_operator_store: undefined,
			filter: undefined,
			start_with_and: false,

			opt_remove_button: true,

			contain_other_cfilter: false,

			margin: 5,

			initComponent: function() {
				this.logAuthor = '[' + this.id + ']';
				log.debug('init sub object', this.logAuthor);
				//------------------create operator combo----------------
				this.operator_combo = Ext.widget('combobox', {
								queryMode: 'local',
								displayField: 'text',
								isFormField: false,
								//Hack: don't search in store
								minChars: 50,
								valueField: 'operator',
								emptyText: _('Type value or choose operator'),
								store: this.operator_store
							});


				//-------------sub operator combo ($in etc...)-----
				this.sub_operator_combo = Ext.widget('combobox', {
								queryMode: 'local',
								displayField: 'text',
								isFormField: false,
								valueField: 'operator',
								value: '$eq',
								editable: false,
								margin: '0 0 0 5',
								store: this.sub_operator_store
							});

				//-----------------is/isnot combo------------------
				this.is_isnot_combo = Ext.widget('combobox', {
								queryMode: 'local',
								displayField: 'text',
								isFormField: false,
								valueField: 'operator',
								value: '$is',
								editable: false,
								margin: '0 0 0 5',
								width: 80,
								store: {
									fields: ['operator', 'text', 'type'],
									data: [
										{'operator': '$is', 'text': _('Is'), 'type': 'value'},
										{'operator': '$not', 'text': _('Is Not'), 'type': 'value' }
									]}
							});

				//--------------------panel-------------------------
				this.add_button = Ext.widget('button', {
					iconCls: 'icon-add',
					margin: '0 0 0 5',
					hidden: true,
					tooltip: _('Add new field/condition')
				});

				if (this.opt_remove_button)
					this.remove_button = Ext.widget('button', {
						iconCls: 'icon-cancel',
						margin: '0 5 0 0',
						tooltip: _('Remove this condition')
					});

				this.string_value = Ext.widget('textfield', {
					margin: '0 0 0 5',
					emptyText: 'Type value here',
					isFormField: false,
					getValue: function() {
						var string = Ext.form.field.Text.superclass.getValue.call(this);
						var number = parseInt(string);
						if (! isNaN(number))
							return number;
						return string;
					}
					});
				this.array_field = Ext.create('cfilter.array_field', {hidden: true});

				var items_array = [];
				if (this.opt_remove_button)
					items_array.push(this.remove_button);
				items_array.push(this.operator_combo, this.is_isnot_combo, this.sub_operator_combo, this.string_value, this.array_field, this.add_button);

				//upper panel
				var config = {
					items: items_array,
					layout: 'hbox',
					border: false
				};

				this.upperPanel = Ext.widget('panel', config);

				//bottom panel
				var config = {
					margin: '0 0 0 20',
					bodyStyle: 'border-top:none;border-bottom:none;border-right:none;'
				};
				this.bottomPanel = Ext.widget('panel', config);

				//----------------------bind events-------------------
				//combo binding
				this.operator_combo.on('change', function(combo,value,oldvalue) {
					this.operator_combo_change(combo, value, oldvalue);
				},this);

				this.sub_operator_combo.on('change', function(combo,value,oldvalue) {
					this.sub_operator_combo_change(combo, value, oldvalue);
				},this);

				this.sub_operator_combo.on('beforeselect',
					this.sub_operator_combo_check_validity
, this);

				//button binding
				this.add_button.on('click', function() {this.add_cfilter()},this);

				if (this.opt_remove_button)
					this.remove_button.on('click', this.remove_button_func, this);
				//-------------------building cfilter-----------------
				this.items = [this.upperPanel, this.bottomPanel];
				this.callParent(arguments);

				//--------------load filter if there is filter--------
				if (this.filter)
					this.setValue(this.filter);
				else;
					if (this.start_with_and) {
						this.operator_combo.setValue('$and');
						this.add_cfilter();
					}
			},


			//launched when value selected in combobox
			operator_combo_change: function(combo,value,oldvalue) {
				log.debug(' + Catch changes on operator combobox, value : ' + value, this.logAuthor);
				var operator_type = this.get_type_from_operator(value, this.operator_store);
				log.debug(' + The type of the operator is: ' + operator_type, this.logAuthor);
				if (operator_type == 'object') {
					//log.debug('   + Field is a known operator',this.logAuthor)
					this.contain_other_cfilter = true;
					this.add_button.show();
					this.string_value.hide();
					this.array_field.hide();
					this.sub_operator_combo.hide();
					this.bottomPanel.show();
					this.is_isnot_combo.hide();
				} else {
					//log.debug('   + Unknown operator',this.logAuthor)
					this.contain_other_cfilter = false;
					this.add_button.hide();
					this.sub_operator_combo.show();
					if (operator_type == 'array')
						this.sub_operator_combo.setValue('$in');
					else
						this.sub_operator_combo_change();
					this.bottomPanel.hide();
					this.is_isnot_combo.show();
				}
			},

			sub_operator_combo_check_validity: function(combo,record) {
				var valid_operator_type = this.get_type_from_operator(this.operator_combo.getValue(), this.operator_store);
				var operator_type = this.get_type_from_operator(record.get('operator'), this.sub_operator_store);

				if (valid_operator_type && valid_operator_type != 'all' && operator_type != valid_operator_type) {
					global.notify.notify('Wrong operator', "You can't use this operator", 'info');
					return false;
				}

			},

			sub_operator_combo_change: function(combo,value,oldvalue) {
				//log.debug(' + Catch changes on sub operator combobox, value : ' + value, this.logAuthor)
				if (!value)
					var value = this.sub_operator_combo.getValue();

				var valid_operator_type = this.get_type_from_operator(this.operator_combo.getValue(), this.operator_store);
				var operator_type = this.get_type_from_operator(value, this.sub_operator_store);
				/*
				if (valid_operator_type && valid_operator_type != 'all' && operator_type != valid_operator_type) {
					global.notify.notify('Wrong operator', "You can't use this operator", 'info');
					this.sub_operator_combo.setValue(oldvalue);
					return;
				}
				*/
				switch (operator_type) {
					case 'value':
						this.string_value.show();
						this.array_field.hide();
						break;
					case 'array':
						this.string_value.hide();
						this.array_field.show();
						break;
					default:
						//log.debug('   + Unrecognized field type',this.logAuthor)
						break;
				}
			},

			//give operator and store, return associated type
			get_type_from_operator: function(operator,store) {
				if (operator && operator.length >= 2) {
					var index_search = store.findExact('operator', operator);
					if (index_search != -1) {
						var operator_record = store.getAt(index_search);
						var operator_type = operator_record.get('type');
						return operator_type;
					}else {
						return null;
					}
				} else {
					return null;
				}
			},

			add_cfilter: function(filter) {
				return this.bottomPanel.add(this.build_field_panel(filter));
			},

			//return an ready to add cfilter
			build_field_panel: function(filter) {
				//Hack: clean store filters (otherwise combo are empty)
				this.operator_store.clearFilter();
				return Ext.create('cfilter.object', {
							operator_store: this.operator_store,
							sub_operator_store: this.sub_operator_store,
							filter: filter
						});
			},

			remove_button_func: function() {
				this.destroy();
			},

			remove_all_cfilter: function() {
				this.bottomPanel.removeAll();
			},

			//------------get / set value--------------------
			getValue: function() {
				var items = this.bottomPanel.items.items;

				var field = undefined;
				if (this.operator_combo.validate())
					field = this.operator_combo.getValue();
				if (!field || field == '')
					return undefined;

				var value = this.string_value.getValue();
				var sub_operator = undefined;

				var output = {};

				if (this.contain_other_cfilter) {
					//get into cfilter
					var values = [];
					//get all cfilter values
					for (var i = 0; i < items.length; i++) {
						var cfilter = items[i];
						var value = cfilter.getValue()
						if (value)
							values.push(value);
					}
				}else {
					//just simple value (no inner cfilter)
					var values = {};
					var sub_operator = this.sub_operator_combo.getValue();
					var sub_operator_type = this.get_type_from_operator(sub_operator, this.sub_operator_store);

					//choose between array or value
					if (sub_operator_type == 'value') {
						if (sub_operator != '$eq')
							values[sub_operator] = this.string_value.getValue();
						else
							values = this.string_value.getValue();
					}else if (sub_operator_type == 'array') {
						values[sub_operator] = this.array_field.getValue();
					}
				}

				//is/isnot combo process
				if (this.is_isnot_combo.getValue() == '$not') {
					if (sub_operator == '$eq') {
						output[field] = {'$ne': values};
					/*}else if(sub_operator == '$regex'){
						console.log('----------------------------------')
						var regex = new RegExp(values[sub_operator])
						console.log(regex)
						output[field] = {"$not": /toto/ }
					*/
					}else {
						output[field] = {'$not': values};
					}
				}else {
					output[field] = values;
				}

				return output;
			},

			setValue: function(filter) {
				log.debug('Set value', this.logAuthor);

				this.remove_all_cfilter();

				if (typeof(filter) == 'string')
					filter = Ext.decode(filter);

				var key = Ext.Object.getKeys(filter)[0];
				var value = filter[key];

				//$not case processing
				if (Ext.isObject(value)) {
					var sub_key = Ext.Object.getKeys(value)[0];
					if (sub_key == '$not') {
						this.is_isnot_combo.setValue('$not');
						if (Ext.isObject(value[sub_key]))
							value = value[sub_key];
					}
				}

				//Hack: clear filter before research, otherwise -> search always = -1
				this.operator_store.clearFilter();
				log.debug('Search for the operator "' + key + '" in store', this.logAuthor);
				var search = this.operator_store.findExact('operator', key);

				this.operator_combo.setValue(key);

				if (search == -1) {
					if (typeof(value) == 'object') {
						log.debug('  + "' + key + '" have a sub operator', this.logAuthor);
						var object_key = Ext.Object.getKeys(value)[0];
						var object_value = value[object_key];
						this.sub_operator_combo.setValue(object_key);

						//check sub operator type
						var sub_operator_type = this.get_type_from_operator(object_key, this.sub_operator_store);
						if (sub_operator_type == 'array') {
							log.debug('   + The sub operator is an array', this.logAuthor);
							this.array_field.setValue(object_value);
						}else {
							log.debug('   + The sub operator is a value', this.logAuthor);
							this.string_value.setValue(object_value);
						}
					}else {
						log.debug('  + "' + key + '" is a simple value', this.logAuthor);
						this.string_value.setValue(value);
					}
				}else {
					log.debug('  + "' + key + '" is a registred operator', this.logAuthor);
					var operator_type = this.get_type_from_operator(key, this.operator_store);
					log.debug('      + Type: ' + operator_type, this.logAuthor);
					if (operator_type == 'array' || operator_type == 'value' || operator_type == 'all' || operator_type == null) {
						log.debug('  + "' + key + '" contain value', this.logAuthor);
						var object_value = value;
						try {
							var object_key = Ext.Object.getKeys(value)[0];
							this.sub_operator_combo.setValue(object_key);
							object_value = value[object_key];
						}catch (err) {
							log.debug('there is no sub operator -> $eq', this.logAuthor);
						}

						if (Ext.isArray(object_value))
							this.array_field.setValue(object_value);
						else
							this.string_value.setValue(object_value);
					}else {
						for (var i = 0; i < value.length; i++) {
							log.debug('  + "' + key + '" contain another cfilter object', this.logAuthor);
							this.add_cfilter(value[i]);
						}
					}
				}
			}
		});
	},

	build_store: function() {
		log.debug('Build stores', this.logAuthor);

		var operator_fields = [
			{'operator': '$nor', 'text': _('Nor'), 'type': 'object'},
			{'operator': '$or', 'text': _('Or'), 'type': 'object'},
			{'operator': '$and', 'text': _('And'), 'type': 'object'}
			//{'operator': '$not', 'text': _('Not'), 'type': 'object'}
		];

		operator_fields = Ext.Array.union(operator_fields, this.operator_fields);

		//---------------------operator store----------------
		this.operator_store = Ext.create('Ext.data.Store', {
			fields: ['operator', 'text', 'type'],
			data: operator_fields
		});

		this.sub_operator_store = Ext.create('Ext.data.Store', {
			fields: ['operator', 'text', 'type'],
			data: [
				{'operator': '$eq', 'text': _('Equal'), 'type': 'value'},
				{'operator': '$lt', 'text': _('Less'), 'type': 'value' },
				{'operator': '$lte', 'text': _('Less or equal'), 'type': 'value' },
				{'operator': '$gt', 'text': _('Greater'), 'type': 'value' },
				{'operator': '$gte', 'text': _('Greater or equal'), 'type': 'value' },
				{'operator': '$all', 'text': _('Match all'), 'type': 'array' },
				{'operator': '$exists', 'text': _('Exists'), 'type': 'value' },
				{'operator': '$ne', 'text': _('Not equal'), 'type': 'value' },
				{'operator': '$in', 'text': _('In'), 'type': 'array'},
				{'operator': '$nin', 'text': _('Not in'), 'type': 'array'},
				{'operator': '$regex', 'text': _('Regex'), 'type': 'value'}
			]
		});

	},

	getRawValue: function() {
		var value = undefined;

		if (!this.edit_area.isHidden()) {
			if (this.edit_area.validate())
				value = strip_return(this.edit_area.getValue());
		} else {
			value = this.cfilter.getValue();
		}

		if (Ext.isObject(value))
			value = Ext.encode(value);

		return value;
	},

	getValue: function() {
		value = this.getRawValue();

		if (value) {
			if (typeof(value) != 'string')
				value = Ext.encode(value);

			log.debug('The filter is : ' + value, this.logAuthor);
			return value;
		}else {
			log.debug('Invalid JSON value', this.logAuthor);
			return undefined;
		}
	},

	setValue: function(value) {
		log.debug('Set value', this.logAuthor);
		if (value != null && value != undefined && value != '') {
			if (typeof(value) == 'string')
				value = Ext.decode(value);

			log.debug('The filter to set is : ' + Ext.encode(value), this.logAuthor);
			this.cfilter.setValue(value);
		}
	},

	beforeDestroy: function() {
		this.checkObjectValidity = false;
		this.callParent(arguments);
	}

});
