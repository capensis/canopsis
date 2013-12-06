//need:app/lib/form/cfield.js
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

Ext.define('canopsis.lib.form.field.cthreshold_metro' , {
	extend: 'Ext.panel.Panel',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.cthreshold_metro',

	border: false,

	layout: 'anchor',

	items: [
		{
			xtype: 'panel',
			anchor: '100% 10%',
			border: 0,
			autoScroll: true,
			items: [{
				xtype: 'label',
				text: 'If:'
			},{
				xtype: 'button',
				text: 'add',
				name: 'add_condition'
			}]
		},{
			xtype: 'panel',
			name: 'if_section',
			anchor: '100% 70%',
			margin: '0 0 0 40',
			border: 0,
			autoScroll: true
		},{
			xtype: 'panel',
			anchor: '100% 10%',
			border: 0,
			items: [{
				xtype: 'label',
				text: 'Else:'
			}]
		},{
			xtype: 'panel',
			name: 'else_section',
			anchor: '100% 10%',
			border: 0,
			items: [{
				xtype: 'combobox',
				displayField: 'text',
				isFormField: false,
				editable: false,
				margin: '0 0 0 40',
				valueField: 'value',
				width: 100,
				store: {
					fields: ['text', 'value'],
					data: [
						{text: 'Ok' , value: '0'},
						{text: 'Warn' , value: '1'},
						{text: 'Crit' , value: '2'},
						{text: 'Unknown' , value: '3'}
					]
				}
			}]
		}
	],

	initComponent: function() {
		this.callParent(arguments);
		this.condition_panel = this.down('panel[name=if_section]');
		this.down('button[name=add_condition]').on('click', this.add_condition, this);
		this.add_condition();
	},

	add_condition: function() {
		var panel = this.condition_panel.add(Ext.create('ctreshold_metro.field'));
		panel.on('moveUp', this.movePanelUp, this);
		panel.on('moveDown', this.movePanelDown, this);
	},

	movePanelUp: function(obj) {
		var element_index = this.condition_panel.items.keys.indexOf(obj.id);

		if(element_index - 1 >= 0) {
			var element = this.condition_panel.remove(Ext.getCmp(obj.id), false);
			this.condition_panel.insert(element_index - 1, element);
		}
	},

	movePanelDown: function(obj) {
		var element_index = this.condition_panel.items.keys.indexOf(obj.id);

		if(element_index + 1 <= this.condition_panel.items.length) {
			var element = this.condition_panel.remove(Ext.getCmp(obj.id), false);
			this.condition_panel.insert(element_index + 1, element);
		}
	},

	getValue: function() {
		var items = this.condition_panel.items.items;
		var output = [];

		for(var i = 0; i < items.length; i++) {
			output.push(items[i].getValue());
		}

		return undefined;
	}
});

Ext.define('ctreshold_metro.field', {
	extend: 'Ext.panel.Panel',

	border: false,
	value: undefined,

	defaults: {
		margin: '0 0 0 5'
	},

	margin: '5 0 0 0',

	layout: 'hbox',

	items: [{
		xtype: 'button',
		text: 'delete',
		name: 'delete'
	},{
		xtype: 'combobox',
		displayField: 'text',
		isFormField: false,
		name: 'metric',
		editable: false,
		valueField: 'value',
		emptyText: _('Type value or choose operator'),
		store: {
			fields: ['text', 'value'],
			data: [
				{text: 'Average' , value: 'average'},
				{text: 'Min' , value: 'min'},
				{text: 'Max' , value: 'max'},
				{text: 'Delta' , value: 'delta'},
				{text: 'Sum' , value: 'sum'}
			]
		}
	},{
		xtype: 'combobox',
		displayField: 'text',
		isFormField: false,
		name: 'operator',
		editable: false,
		valueField: 'value',
		width: 40,
		store: {
			fields: ['text', 'value'],
			data: [
				{'text': '>',	'value': '>' },
				{'text': '>=',	'value': '>=' },
				{'text': '<',	'value': '<' },
				{'text': '<=',	'value': '<=' },
				{'text': '=',	'value': '=' },
				{'text': '!=',	'value': '!=' }
			]
		}
	},{
		xtype: 'textfield',
		name: 'value',
		emptyText: 'Value',
		allowBlank: false
	},{
		xtype: 'label',
		text: '=>',
		margin: '0 5 0 5'
	},{
		xtype: 'combobox',
		displayField: 'text',
		isFormField: false,
		editable: false,
		valueField: 'value',
		name: 'state',
		width: 100,
		store: {
			fields: ['text', 'value'],
			data: [
				{text: 'Ok' , value: '0'},
				{text: 'Warn' , value: '1'},
				{text: 'Crit' , value: '2'},
				{text: 'Unknown' , value: '3'}
			]
		}
	},{
		xtype: 'button',
		text: 'up',
		name: 'up'
	},{
		xtype: 'button',
		text: 'down',
		name: 'down'
	}],

	initComponent: function() {
		this.callParent(arguments);

		this.down('button[name=delete]').on('click', this.destroy, this);

		this.down('button[name=up]').on('click', function() {
			this.fireEvent('moveUp', {id: this.id});
		}, this);

		this.down('button[name=down]').on('click', function() {
			this.fireEvent('moveDown', {
				id: this.id
			});
		}, this);
	},

	getValue: function() {
		var metric = this.down('combobox[name=metric]').getValue();
		var operator = this.down('combobox[name=operator]').getValue();
		var value = this.down('textfield[name=value]').getValue();
		var state = this.down('combobox[name=state]').getValue();

		console.log(metric + ' ' + operator + ' ' + value + ' ' + state);
	}
});
