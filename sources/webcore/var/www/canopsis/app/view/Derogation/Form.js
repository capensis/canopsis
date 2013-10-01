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
Ext.define('canopsis.view.Derogation.Form', {
	extend: 'canopsis.lib.view.cform',
	alias: 'widget.DerogationForm',

	width: 500,
	minHeight: 460,
	bodyPadding: 10,
	border: false,
	now: false,

	items: [{
		xtype: 'tabpanel',
		deferredRender: false,
		border: false,
		plain: true,

		defaults: {
			border: false,
			autoScroll: true,
			layout: 'anchor',
		},

		items: [{
			title: _('General'),
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
				xtype: 'textarea',
				name: 'description',
				fieldLabel: 'Description',
				width: 295,
			},{
				xtype: 'fieldcontainer',
				fieldLabel: _('Begin'),
				layout: 'hbox',
				items: [{
					xtype: 'cdate',
					name: 'startTs',
					date_width: 110,
					now: this.now,
				}],
			},{
				xtype: 'fieldcontainer',
				fieldLabel: _('Ending'),
				layout: 'hbox',
				items: [{
					xtype: 'combobox',
					name: 'periodType',

					isFormField: false,
					editable: false,

					width: 60,
					margin: '0 5 0 0',

					queryMode: 'local',
					displayField: 'text',
					valueField: 'value',
					value: 'for',

					store: {
						xtype: 'store',
						fields: ['value', 'text'],
						data: [
							{value: 'for', text: _('For')},
							{value: 'to', text: _('To')},
						]
					},

					listeners: {
						change: function(combo, value) {
							var me = combo.up('DerogationForm');
							var durationDate = me.down('cduration[name=forTs]');
							var stopDate     = me.down('cdate[name=stopTs]');

							if (value == 'for') {
								durationDate.show();
								stopDate.hide();
								stopDate.setDisabled(true);
							}
							else if (value == 'to') {
								durationDate.hide();
								stopDate.show();
								stopDate.setDisabled(false);
							}
						},
					},
				},{
					xtype: 'cduration',
					name: 'forTs',
					value: global.commonTs.day,
				},{
					xtype: 'cdate',
					name: 'stopTs',
					hidden: true,
					disabled: true,
				}]
			}],
		},{
			title: _('Override'),
			items: [{
				xtype: 'button',
				text: _('Add'),
				iconCls: 'icon-add',
				margin: '5 5 5 5',

				listeners: {
					click: function() {
						var me = this.up('DerogationForm');

						me.addNewField();
					}
				}
			},{
				xtype: 'cfieldset',
				title: _('Actions'),
				border: 'false',
				items: [],
			}],
		},{
			title: _('Requalificate'),
		},{
			title: _('Filter'),
			xtype: 'cfilter',
			name: 'evfilter',
		}]
	}],

	initComponent: function() {
		this.callParent();

		this.periodTypeCombo = this.down('combobox[name=periodType]');
		this.durationDate    = this.down('cduration[name=forTs]');
		this.startDate       = this.down('cdate[name=startTs]');
		this.stopDate        = this.down('cdate[name=stopTs]');

		if(!this.editing) {
			this.addNewField();
		}
	},

	afterRender: function() {
		this.callParent();

		var tabpanel = this.down('tabpanel');

		for(var i = 0; i < tabpanel.items.getCount(); ++i) {
			tabpanel.setActiveTab(i);
		}

		tabpanel.setActiveTab(0);
	},

	addNewField: function(variable, value) {
		log.debug(' + Adding a new field', this.logAuthor);

		var actions = this.down('cfieldset[title="' + _('Actions') + '"]');

		actions.add(Ext.create('derogation.override', {_variable: variable, _value: value}));
	}
});

Ext.define('derogation.override', {
	extend: 'Ext.form.Panel',
	mixins: ['canopsis.lib.form.cfield'],

	border: false,
	layout: 'hbox',
	//bodyStyle:{'background': '#ededed'},

	state_icon_path: 'widgets/weather/icons/set1/',
	icon_weather1: 'state_0.png',
	icon_weather2: 'state_1.png',
	icon_weather3: 'state_2.png',

	alert_icon_path: 'widgets/weather/icons/alert/',
	icon_alert1: 'workman.png',
	icon_alert2: 'slippery.png',
	icon_alert3: 'alert.png',

	icon_class: 'widget-weather-form-icon',

	name: 'actions',

	initComponent: function() {

		//------------------config objects-----------------

		var config_key_field = {
			isFormField: false,
			editable: false,
			flex: 1,
			labelWidth: 40,
			margin: '5 0 0 0',
			fieldLabel: _('Field'),
			queryMode: 'local',
			displayField: 'text',
			valueField: 'value',
			value: 'output',
			name: 'key_field',
			store: {
				xtype: 'store',
				fields: ['value', 'text'],
				data: [
					{value: 'state', text: _('State')},
					{value: 'output', text: _('Comment')},
					{value: 'alert_msg', text: _('Alert message')},
					{value: 'alert_icon', text: _('Alert icon')}
				]
			}
		};

		if (this._variable)
			config_key_field.disabled = true;


		var config_list_state = {
			stateIconPath: this.state_icon_path,
			iconClass: this.icon_class,
			isFormField: false,
			xtype: 'combobox',
			editable: false,
			margin: '5 5 0 5',
			disabled: true,
			hidden: true,
			flex: 1,
			name: 'state',
			displayField: 'text',
			valueField: 'value',
			queryMode: 'local',
			value: 0,
			listConfig: {
				getInnerTpl: function() {
					var tpl = '<div>' +
							  '<img src="' + this.findParentByType('combobox').stateIconPath + '{icon}" class="' + this.findParentByType('combobox').iconClass + '"/>' +
							  '{text}</div>';
					return tpl;
				}
			},
			store: {
				xtype: 'store',
				fields: ['value', 'text', 'icon'],
				data: [
					{value: 0, text: _('Ok'), icon: this.icon_weather1 },
					{value: 1, text: _('Warning'), icon: this.icon_weather2 },
					{value: 2, text: _('Critical'), icon: this.icon_weather3 }
				]
			}
		};

		var config_alertIcon_radio = {
			alertIconPath: this.alert_icon_path,
			iconClass: this.icon_class,
			isFormField: false,
			border: false,
			editable: false,
			margin: '5 5 0 5',
			disabled: true,
			hidden: true,
			flex: 1,
			name: 'alert_icon',
			displayField: 'text',
			valueField: 'value',
			queryMode: 'local',
			value: 0,
			listConfig: {
				getInnerTpl: function() {
					var tpl = '<div>' +
							  '<img src="' + this.findParentByType('combobox').alertIconPath + '{icon}" class="' + this.findParentByType('combobox').iconClass + '"/>' +
							  '{text}</div>';
					return tpl;
				}
			},
			store: {
				xtype: 'store',
				fields: ['value', 'text', 'icon'],
				data: [
					{value: 0, text: _('Indisponible'), icon: this.icon_alert1 },
					{value: 1, text: _('Be carefull'), icon: this.icon_alert2 },
					{value: 2, text: _('Simple alert'), icon: this.icon_alert3 }
				]
			}
		};

		var config_output_textfield = {
			isFormField: false,
			flex: 1,
			name: 'output',
			emptyText: _('Type here new comment...'),
			margin: '5 5 0 5'
		};

		//if value, not display comment by default
		if (this._variable) {
			config_output_textfield.disabled = true;
			config_output_textfield.hidden = true;
		}

		var config_alert_textfield = {
			isFormField: false,
			flex: 1,
			disabled: true,
			hidden: true,
			name: 'alert_msg',
			emptyText: _('Type here alert message...'),
			margin: '5 5 0 5'
		};



		//---------------------build objects-------------
		this.key_field = Ext.widget('combobox',	config_key_field);

		this.items = [this.key_field];

		if (!this._variable || this._variable == 'state') {
			this.list_state = Ext.widget('combobox', config_list_state);
			this.items.push(this.list_state);
		}

		if (!this._variable || this._variable == 'alert_icon') {
			this.alertIcon_radio = Ext.widget('combobox', config_alertIcon_radio);
			this.items.push(this.alertIcon_radio);
		}

		if (!this._variable || this._variable == 'output') {
			this.output_textfield = Ext.widget('textfield', config_output_textfield);
			this.items.push(this.output_textfield);
		}

		if (!this._variable || this._variable == 'alert_msg') {
			this.alert_textfield = Ext.widget('textfield', config_alert_textfield);
			this.items.push(this.alert_textfield);
		}

		//----------------------other---------------------

		this.destroyButton = Ext.widget('button', {
			iconCls: 'icon-cancel',
			margin: '5 0 0 0'
		});
		this.items.push(this.destroyButton);

		this.callParent(arguments);

		//----------------------bind events--------------------
		this.key_field.on('select', this._onChange, this);
		this.destroyButton.on('click', this.selfDestruction, this);
	},

	afterRender: function() {
		this.callParent(arguments);
		if (this._variable) {
			this.key_field.setValue(this._variable);
			this.change(this._variable);
			var field = this.down('[name=' + this._variable + ']');
			if (field)
				field.setValue(this._value);
		}

	},

	selfDestruction: function() {
		//tweak, otherwise the textfield is not deleted
		Ext.destroy(this.output_textfield);
		Ext.destroy(this.alert_textfield);
		Ext.destroy(this);
	},

	_onChange: function(combo,records,options) {
		var value = records[0].get('value');
		this.change(value);
	},

	change: function(value) {
		var fields = Ext.ComponentQuery.query('#' + this.id + ' [name]');
		for (var i = 0; i < fields.length; i++) {
			var elem = fields[i];
			if (elem.name != 'key_field') {
				if (elem.name != value) {
					elem.hide();
					elem.setDisabled(true);
				}else {
					elem.show();
					elem.setDisabled(false);
				}
			}
		}
	},

	getValue: function() {
		var field = this.key_field.getValue();
		var value = this.down('[name=' + field + ']').getValue();
		return {type: 'override', field: field, value: value};
	}
});
