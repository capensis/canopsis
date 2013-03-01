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

Ext.define('cfilter.object' , {
	extend: 'Ext.panel.Panel',
	border: false,
	margin: 5,
	autoScroll:true,

	initialCfilter: false,
	filter: false,

	items: [{
		name:'upperPanel',
		xtype: 'panel',
		layout: 'hbox',
		border: false,
		items:[{
			xtype:'button',
			name:'cfilterRemoveButton',
			iconCls: 'icon-cancel',
			margin: '0 5 0 0',
			width: 24,
			tooltip: _('Remove this from list of value')
		},{
			name:'cfilterField',
			xtype:'combobox',
			queryMode: 'local',
			displayField: 'text',
			isFormField: false,
			minChars: 50,
			valueField: 'operator',
			emptyText: _('Type value or choose operator'),
			store: {
					"xtype": "store",
					"fields": ["operator", "text","type"],
					"data" : []
				}
		},{
			name:'cfilterIsCombo',
			xtype:'combobox',
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
		},{
			name:'cfilterOperator',
			xtype:'combobox',
			queryMode: 'local',
			displayField: 'text',
			isFormField: false,
			valueField: 'operator',
			editable: false,
			margin: '0 0 0 5',
			store: {
				fields: ['operator', 'text', 'type'],
				data: []
			}
		},{
			cfilterField:true,
			cfilterType: 'string',
			xtype: 'textfield',
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
		},{
			xtype:'combobox',
			cfilterField:true,
			cfilterType: 'bool',
			margin: '0 0 0 5',
			value: true,
			hidden:true,
			isFormField: false,
			displayField: "text",
			valueField: "value",
			store: {
				xtype: 'store',
				fields: ["value", "text"],
				data : [
					{"value": true, "text": "True"},
					{"value": false, "text": "False"}
				]
			}
		},{
			cfilterField:true,
			cfilterType: 'date',
			xtype: 'datepicker',
			hidden:true,
			margin: '0 0 0 5',
			isFormField: false,
		},{
			cfilterField:true,
			cfilterType: 'array',
			extend:'cfilter.array_field',
			hidden:true,
		},{
			xtype:'button',
			name: 'cfilterAddButton',
			iconCls: 'icon-add',
			margin: '0 0 0 5',
			hidden: true,
			tooltip: _('Add new field/condition')
		}]
	},{
		cfilterField:true,
		cfilterType: 'object',
		xtype:'panel',
		name:'lowerPanel',
		margin: '0 0 0 20',
		bodyStyle: 'border-top:none;border-bottom:none;border-right:none;'
	}],
	
	initComponent: function() {
		this.logAuthor = '[' + this.id + ']';
		this.callParent(arguments);

		//stock cfilterField elements
		this.cfilterFieldElements = [this.down('panel[name=lowerPanel]')]
		var upperPanelId = this.down('panel[name=upperPanel]').id
		this.cfilterFieldElements = Ext.Array.union(
			this.cfilterFieldElements,
			Ext.ComponentQuery.query('#' + upperPanelId + ' > *[cfilterField]')
		)
		
		//stock frequently used element
		this.cfilterField = this.down('combobox[name=cfilterField]')
		this.cfilterOperator = this.down('combobox[name=cfilterOperator]')
		this.fieldStore = this.cfilterField.getStore()
		this.operatorStore = this.cfilterOperator.getStore()

		//prepare cfilter
		this.fieldStore.loadData(this.fields_array)
		this.operatorStore.loadData(this.operators_array)
		this.cfilterOperator.setValue('$eq')
		if(this.initialCfilter)
			this.down('button[name=cfilterRemoveButton]').hide()

		//bind events
		this.cfilterField.on('select',this.fieldChange,this)
		this.cfilterOperator.on('select',this.operatorChange,this)
		this.down('button[name=cfilterAddButton]').on('click',this.createInnerCfilter,this)
		this.down('button[name=cfilterRemoveButton]').on('click',function() {this.destroy()},this)

		//set data if existing
		if(this.filter)
			this.setValue(this.filter)
	},

	fieldChange: function(combo,records){
		log.debug('Field changed',this.logAuthor)
		var record = records[0]
		var allowed_type = record.get('type')

		if(allowed_type){
			if(allowed_type != 'all'){
				if(allowed_type == 'object'){
					if(!this.haveInnerCfilter)
						this.createInnerCfilter()
					this.showOnValueType('object')
				}else{
					this.operatorStore.clearFilter(true)
					this.operatorStore.filterBy(function(record){
							var record_types = record.get('type')
							if(Ext.Array.indexOf(record_types,allowed_type) == -1)
								return false
							else
								return true
					},this)
					this.showOnValueType(this.getValueType())
				}
			}else{
				this.operatorStore.clearFilter(false)
				this.showOnValueType(this.getValueType())
			}	
		}
	},

	createInnerCfilter: function(data){
		var cfilter = Ext.create('cfilter.object',{
			fields_array:this.fields_array,
			operators_array:this.operators_array,
			filter:data
		})

		this.down('panel[name=lowerPanel]').add(cfilter)

		if(!this.innerCfilter)
			this.innerCfilter = true
	},

	getFieldRecord: function(){
		var recordId = this.fieldStore.find('operator',this.cfilterField.getValue())
		if(recordId == -1)
			return undefined
		return this.fieldStore.getAt(recordId)
	},

	getOperatorRecord:function(){
		var recordId = this.operatorStore.find('operator',this.cfilterOperator.getValue())
		return this.operatorStore.getAt(recordId)
	},

	operatorChange: function(combo,record_or_records){
		log.debug('Operator changed',this.logAuthor)

		if(Ext.isArray(record_or_records))
			var operatorRecord = record_or_records[0]
		else
			var operatorRecord = record_or_records

		var operatorRecordType = operatorRecord.get('type')

		var fieldRecord_index = this.fieldStore.find('operator',this.cfilterField.getValue())
		if(fieldRecord_index != -1){
			var fieldRecord = this.fieldStore.getAt(fieldRecord_index)
			var fieldRecordType = fieldRecord.get('type')
		}else{
			//WARNING CLEAN THAT ----
			this.showOnValueType(operatorRecordType[0])
			return
		}

		if(fieldRecordType == 'all')
			//IF field doesn't require specific value (ex: "custom field")
			if(operatorRecord.get('array'))
				log.debug('array')
			else
				//WARNING CLEAN THAT ----
				this.showOnValueType(operatorRecordType[0])
		else
			//field require specif value, like "timestamp" who need date
			this.showOnValueType(fieldRecordType)

	},

	//this function aimed to determine final type of value (string/bool...)
	//first we check field type (timestamp need date, then operator type ($exist need bool)
	getValueType: function(){
		var operatorRecord = this.getOperatorRecord()
		var fieldRecord = this.getFieldRecord()

		var fieldType = undefined
		if(fieldRecord)
			var fieldType = fieldRecord.get('type')

		if(!fieldType || fieldType == 'all')
			return operatorRecord.get('type')[0]

		return fieldType
		
	},

	showOnValueType: function(type){
		var elements = this.cfilterFieldElements
		for(var i=0; i < elements.length; i++)
			if(elements[i].cfilterType == type)
				elements[i].show()
			else
				elements[i].hide()

		if(type =='object'){
			this.down('button[name=cfilterAddButton]').show()
			this.down('combobox[name=cfilterIsCombo]').hide()
			this.cfilterOperator.hide()
		}else{
			this.down('button[name=cfilterAddButton]').hide()
			this.down('combobox[name=cfilterIsCombo]').show()
			this.cfilterOperator.show()
		}
	},

	//return the value of the elements corresponding of given type (ex:"string/bool/date ...")
	getValueByElementType: function(type){
		var elements = this.cfilterFieldElements
		for(var i=0; i < elements.length; i++)
			if(elements[i].cfilterType == type)
				if(elements[i].getValue)
					return elements[i].getValue()
	},

	getValue: function(){
		var fieldRecord = this.getFieldRecord()
		var operatorRecord = this.getOperatorRecord()
		var operator = operatorRecord.get('operator')
		var isIsNotValue = this.down('combobox[name=cfilterIsCombo]').getValue()
		var inputValue = this.getValueByElementType(this.getValueType())
		var output = {}
		var values = {}

		//if contained another cfilter
		if(fieldRecord && fieldRecord.get('type') == 'object'){
			var listCfilterResult = []
			var panel = this.down('panel[name=lowerPanel]')
			for(var i=0; i < panel.items.items.length; i++)
				listCfilterResult.push(panel.items.items[i].getValue())
			output[this.cfilterField.getValue()] = listCfilterResult
			return output
		}

		//Get operator
		if(operatorRecord.get('operator') == '$eq')
			values = inputValue
		else
			values[operator] = inputValue

		//manage negation
		if(isIsNotValue == '$not')
			if(operator == '$eq')
				values = {'$ne': values};
			else
				values = {'$not': values};

		output[this.cfilterField.getValue()] = values

		return output
	},

	setValue: function(filter){
		var key = Ext.Object.getKeys(filter)[0];
		var value = filter[key];

		this.cfilterField.setValue(key)

		//if $and/$or
		if(Ext.isArray(value)){
			for(var i=0; i< value.length;i++)
				this.createInnerCfilter(value[i])
			this.showOnValueType('object')
			return
		}

		if(!Ext.isObject(value)){
			var type = this.getValueType()
			this.down('*[cfilterField=true][cfilterType='+type+']').setValue(value)
			return
		}

		//operator or not
		key = Ext.Object.getKeys(value)[0];
		value = value[key]

		if(key == '$not')
			this.down('combobox[name=cfilterIsCombo]').setValue('$not')

		if(Ext.isObject(value)){
			key = Ext.Object.getKeys(value)[0];
			this.cfilterOperator.setValue(key)
			value = value[key]
		}

		var type = this.getValueType()
		this.down('*[cfilterField=true][cfilterType='+type+']').setValue(value)


	}

})


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
		{'operator': 'tags', 'text': _('Tags'),	'type': 'all'},
		{'operator': 'timestamp', 'text': _('Timestamp'),	'type': 'date'}
	],

	layout: {
        type: 'vbox',
        align: 'stretch'
    },

    checkObjectValidity: true,

	initComponent: function() {
		this.logAuthor = '[' + this.id + ']';
		log.debug('Initialize ...', this.logAuthor);

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
			fields_array: this.operator_store,
			operators_array: this.sub_operator_store,
			opt_remove_button: false,
			initialCfilter: true,
			flex:1
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
		//log.debug('Execute isValid function', this.logAuthor);
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
		/*this.operator_store = Ext.create('Ext.data.Store', {
			fields: ['operator', 'text', 'type'],
			data: operator_fields
		});*/
		this.operator_store = operator_fields 
		this.sub_operator_store = [
				{'operator': '$eq', 'text': _('Equal'), 'type': ['string','date'], 'array':false},
				{'operator': '$lt', 'text': _('Less'), 'type': ['string','date'], 'array':false },
				{'operator': '$lte', 'text': _('Less or equal'), 'type': ['string','date'], 'array':false },
				{'operator': '$gt', 'text': _('Greater'), 'type': ['string','date'], 'array':false },
				{'operator': '$gte', 'text': _('Greater or equal'), 'type': ['string','date'], 'array':false },
				{'operator': '$all', 'text': _('Match all'), 'type': ['array'], 'array':true },
				{'operator': '$exists', 'text': _('Exists'), 'type': ['bool'], 'array':false },
				{'operator': '$ne', 'text': _('Not equal'), 'type': ['string','date'], 'array':false },
				{'operator': '$in', 'text': _('In'), 'type': ['string'], 'array':true},
				{'operator': '$nin', 'text': _('Not in'), 'type': ['string'],'array':true },
				{'operator': '$regex', 'text': _('Regex'), 'type': ['string'], 'array':false}
			]
		/*
		this.sub_operator_store = Ext.create('Ext.data.Store', {
			fields: ['operator', 'text', 'type'],
			data: [
				{'operator': '$eq', 'text': _('Equal'), 'type': ['value','date']},
				{'operator': '$lt', 'text': _('Less'), 'type': ['value','date'] },
				{'operator': '$lte', 'text': _('Less or equal'), 'type': ['value','date'] },
				{'operator': '$gt', 'text': _('Greater'), 'type': ['value','date'] },
				{'operator': '$gte', 'text': _('Greater or equal'), 'type': ['value','date'] },
				{'operator': '$all', 'text': _('Match all'), 'type': 'array' },
				{'operator': '$exists', 'text': _('Exists'), 'type': 'bool' },
				{'operator': '$ne', 'text': _('Not equal'), 'type': ['value','date'] },
				{'operator': '$in', 'text': _('In'), 'type': 'array'},
				{'operator': '$nin', 'text': _('Not in'), 'type': 'array'},
				{'operator': '$regex', 'text': _('Regex'), 'type': 'value'}
			]
		});*/

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
