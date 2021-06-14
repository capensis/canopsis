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
Ext.define('canopsis.view.Statemap.Form', {
	extend: 'canopsis.lib.view.cform',
	alias: 'widget.StatemapForm',

	width: 500,
	minHeight: 460,
	border: false,
	bodyPadding: 10,

	items: [{
		xtype: 'hiddenfield',
		name: '_id',
		value: undefined,
	},{
		xtype: 'hiddenfield',
		name: 'tags',
		value: undefined,
	},{
		xtype: 'textfield',
		name: 'crecord_name',
		fieldLabel: _('Name'),
		allowBlank: false,
		width: 295,
	},{
		xtype: 'statemapfield',
		name: 'statemap',
		fieldLabel: _('States'),
	}],

	initComponent: function() {
		this.callParent();
	},
});

Ext.define('statemap.field', {
	extend: 'Ext.form.Panel',
	mixins: ['canopsis.lib.form.cfield'],
	alias: 'widget.statemapfield',

	border: false,
	plain: true,
	bodyPadding: 10,
	layout: 'vbox',

	state_icon_path: 'widgets/weather/icons/set1/',
	icon_weather1: 'state_0.png',
	icon_weather2: 'state_1.png',
	icon_weather3: 'state_2.png',
	icon_weather4: 'state_3.png',

	icon_class: 'widget-weather-form-icon',

	initComponent: function() {
		this.tbar = [{
			xtype: 'button',
			text: _('Add'),
			iconCls: 'icon-add',

			listeners: {
				click: function() {
					this.addStateAssoc(0, 3);
				},
				scope: this,
			}
		}];

		this.callParent();
	},

	addStateAssoc: function(state, canostate) {
		var me = this;

		this.add({
			xtype: 'fieldcontainer',
			isFormField: false,

			layout: 'hbox',
			margin: '5 0 0 0',

			defaults: {
				labelWidth: 50,
				border: false,
			},

			items: [{
				xtype: 'numberfield',
				name: 'state',
				isFormField: false,
				editable: false,
				fieldLabel: _('State'),
				width: 100,
				margin: '0 0 0 5',
				value: state,
			},{
				xtype: 'combobox',
				name: 'canostate',
				isFormField: false,
				editable: false,
				width: 100,
				margin: '0 0 0 5',

				queryMode: 'local',
				displayField: 'text',
				valueField: 'value',
				value: canostate,

				listConfig: {
					getInnerTpl: function() {
						return '<div><img src="' + me.state_icon_path + '{icon}" class="' + me.icon_class + '"/>{text}</div>';
					}
				},

				store: {
					xtype: 'store',
					fields: ['value', 'text', 'icon'],
					data: [
						{value: 0, text: _('Ok'), icon: this.icon_weather1 },
						{value: 1, text: _('Warning'), icon: this.icon_weather2 },
						{value: 2, text: _('Critical'), icon: this.icon_weather3 },
						{value: 3, text: _('Unknown'), icon: this.icon_weather4 },
					]
				},
			},{
				xtype: 'button',
				iconCls: 'icon-cancel',
				margin: '0 0 0 5',

				listeners: {
					click: function() {
						var fcontainer = this.up('fieldcontainer');
						Ext.destroy(fcontainer);
					},
				},
			}]
		});
	},

	getValue: function() {
		var value = [];
		var me = this;

		this.items.each(function() {
			var state     = this.down('[name=state]').getValue();
			var canostate = this.down('[name=canostate]').getValue();

			value[state] = canostate;
		});

		/* if the user didn't specified some necessary states, assume it's UNKNOWN */
		for(var idx = 0; idx < value.length; ++idx) {
			if(value[idx] === undefined) {
				value[idx] = 3;
			}
		}

		return value;
	},

	setValue: function(val) {
		for(var i = 0; i < val.length; ++i) {
			this.addStateAssoc(i, val[i]);
		}
	},
});
