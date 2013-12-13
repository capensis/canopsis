//need:app/lib/view/cform.js,app/lib/form/cfield.js,app/lib/form/field/cdate.js,app/lib/form/field/cfieldset.js,app/lib/form/field/cduration.js
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
Ext.define('canopsis.view.Derogation.Form', {
	extend: 'canopsis.lib.view.cform',

	requires: [
		'canopsis.lib.form.field.cdate',
		'canopsis.lib.form.field.cduration'
	],

	alias: 'widget.DerogationForm',

	width: 700,
	minHeight: 560,
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

							if(value === 'for') {
								durationDate.show();
								stopDate.hide();
								stopDate.setDisabled(true);
							}
							else if(value === 'to') {
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
				border: false,
				items: [],
			}],
		},{
			title: _('Requalificate'),
			xtype: 'DerogationStatemapField',
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
		this.eventFilter     = this.down('cfilter[name=evfilter]');

		if(!this.editing) {
			this.addNewField();
		}
	},

	addNewField: function(variable, value) {
		log.debug(' + Adding a new field', this.logAuthor);

		var actions = this.down('cfieldset[title="' + _('Actions') + '"]');

		actions.add(Ext.create('derogation.override', {
			variable: variable,
			value: value
		}));
	},

	setRequalification: function(statemap_id) {
		var statemapfield = this.down('[name=statemap]');

		statemapfield.setValue(statemap_id);
		statemapfield.statemap_id = statemap_id;
	}
});

Ext.define('derogation.override', {
	extend: 'Ext.form.Panel',
	mixins: ['canopsis.lib.form.cfield'],
	alias: 'widget.derogationfield',

	border: false,
	plain: true,
	bodyPadding: 10,
	layout: 'hbox',
	name: 'actions',

	state_icon_path: 'widgets/weather/icons/set1/',
	icon_weather1: 'state_0.png',
	icon_weather2: 'state_1.png',
	icon_weather3: 'state_2.png',

	alert_icon_path: 'widgets/weather/icons/alert/',
	icon_alert1: 'workman.png',
	icon_alert2: 'slippery.png',
	icon_alert3: 'alert.png',

	icon_class: 'widget-weather-form-icon',

	initComponent: function() {
		var me = this;

		this.callParent();

		this.add({
			xtype: 'combobox',

			name: 'key_field',
			isFormField: false,
			editable: false,

			disabled: (!!this.variable),
			flex: 1,

			queryMode: 'local',
			displayField: 'text',
			valueField: 'value',
			value: this.variable || 'output',

			store: {
				xtype: 'store',
				fields: ['value', 'text'],
				data: [
					{value: 'state', text: _('State')},
					{value: 'output', text: _('Comment')},
					{value: 'alert_msg', text: _('Alert message')},
					{value: 'alert_icon', text: _('Alert icon')},
				],
			},

			listeners: {
				select: function(combo, records) {
					void(combo);

					var value = records[0].get('value');

					var fields = {
						state:      me.down('[name=state]'),
						alert_icon: me.down('[name=alert_icon]'),
						output:     me.down('[name=output]'),
						alert_msg:  me.down('[name=alert_msg]')
					};

					for(key in fields) {
						fields[key].hide();
						fields[key].setDisabled(true);
					}

					fields[value].show();
					fields[value].setDisabled(false);

					me.variable = value;
					me.value    = fields[value].getValue();
				}
			}
		},{
			xtype: 'combobox',

			name: 'state',
			isFormField: false,
			editable: false,

			disabled: (this.variable !== 'state'),
			hidden: (this.variable !== 'state'),
			autoRender: true,
			flex: 2,

			queryMode: 'local',
			displayField: 'text',
			valueField: 'value',
			value: this.value || 0,

			store: {
				xtype: 'store',
				fields: ['value', 'text', 'icon'],
				data: [
					{value: 0, text: _('Ok'), icon: this.icon_weather1},
					{value: 1, text: _('Warning'), icon: this.icon_weather2},
					{value: 2, text: _('Critical'), icon: this.icon_weather3},
				],
			},

			listConfig: {
				getInnerTpl: function() {
					return '<div><img src="' + me.state_icon_path + '{icon}" class="' + me.icon_class + '" /> {text}</div>';
				},
			},

			listeners: {
				change: function(value) {
					if(me.variable === 'state') {
						me.value = value;
					}
				},
			},
		},{
			xtype: 'combobox',

			name: 'alert_icon',
			isFormField: false,
			editable: false,

			disabled: (this.variable !== 'alert_icon'),
			hidden: (this.variable !== 'alert_icon'),
			autoRender: true,
			flex: 2,

			queryMode: 'local',
			displayField: 'text',
			valueField: 'value',
			value: this.value || 0,

			store: {
				xtype: 'store',
				fields: ['value', 'text', 'icon'],
				data: [
					{value: 0, text: _('Unavailable'), icon: this.icon_alert1},
					{value: 1, text: _('Be carefull'), icon: this.icon_alert2},
					{value: 2, text: _('Simple alert'), icon: this.icon_alert3},
				],
			},

			listConfig: {
				getInnerTpl: function() {
					return '<div><img src="' + me.alert_icon_path + '{icon}" class="' + me.icon_class + '" /> {text}</div>';
				},
			},

			listeners: {
				change: function(value) {
					if(me.variable === 'alert_icon') {
						me.value = value;
					}
				},
			},
		},{
			xtype: 'textfield',

			name: 'output',
			isFormField: false,

			disabled: (this.variable !== 'output'),
			hidden: (this.variable !== 'output'),
			autoRender: true,
			flex: 2,

			emptyText: _('Type here new comment...'),
			value: this.value || '',

			listeners: {
				change: function(value) {
					if(me.variable === 'output') {
						me.value = value;
					}
				},
			},
		},{
			xtype: 'textfield',

			name: 'alert_msg',
			isFormField: false,

			disabled: (this.variable !== 'alert_msg'),
			hidden: (this.variable !== 'alert_msg'),
			autoRender: true,
			flex: 2,

			emptyText: _('Type here alert message...'),
			value: this.value || '',

			listeners: {
				change: function(value) {
					if(me.variable === 'alert_msg') {
						me.value = value;
					}
				},
			},
		},{
			xtype: 'button',
			iconCls: 'icon-cancel',

			listeners: {
				click: function() {
					Ext.destroy(me);
				},
			},
		});

		if(!this.variable) {
			var outputfield = this.down('[name=output]');

			outputfield.setDisabled(false);
			outputfield.hidden = false;
		}
	},

	getValue: function() {
		var key = this.down('[name=key_field]').getValue();

		var fields = {
			state:      this.down('[name=state]').getValue(),
			alert_icon: this.down('[name=alert_icon]').getValue(),
			output:     this.down('[name=output]').getValue(),
			alert_msg:  this.down('[name=alert_msg]').getValue()
		};

		return {type: 'override', field: key, value: fields[key]};
	}
});

Ext.define('derogation.statemap', {
	extend: 'canopsis.lib.view.cgrid',
	mixins: ['canopsis.lib.form.cfield'],
	alias: 'widget.DerogationStatemapField',

	requires: [
		'canopsis.store.Statemaps',
	],

	opt_paging: true,
	opt_menu_delete: false,
	opt_bar_add: false,
	opt_menu_rights: false,
	opt_bar_search: true,
	opt_bar_enable: false,
	opt_tags_search: false,

	selType: 'rowmodel',
	model: 'Statemap',
	store: 'Statemaps',

	minHeight: 480,

	name: 'statemap',

	listeners: {
		itemdblclick: function() {
			/* ignore double-clicks usually handled to edit items */
			return false;
		},

		select: function(grid, record) {
			this.statemap_id = record.data._id;
		},

		selectionchange: function(grid, selected) {
			if(selected.length === 0) {
				this.statemap_id = undefined;
			}
		}
	},

	rdr_statemap: function(val) {
		var output = '<p>';

		for(var i = 0; i < val.length; ++i) {
			output += '<b>' + i + '</b> -> ';

			switch(val[i]) {
				case 0:
					output += 'OK';
					break;

				case 1:
					output += 'WARNING';
					break;

				case 2:
					output += 'CRITICAL';
					break;

				default:
					output += 'UNKNOWN';
					break;
			}

			if (i !== val.length - 1) {
				output += ', ';
			}
		}

		return output;
	},

	initComponent: function() {
		this.columns = [{
			header: '',
			width: 25,
			sortable: false,
			renderer: rdr_crecord_type,
			dataIndex: 'crecord_type',
		},{
			header: _('Enabled'),
			align: 'center',
			width: 55,
			dataIndex: 'enable',
			renderer: rdr_boolean,
		},{
			header: _('Name'),
			flex: 1,
			dataIndex: 'crecord_name',
		},{
			header: _('Statemap'),
			flex: 1,
			dataIndex: 'statemap',
			renderer: this.rdr_statemap,
		}];

		this.callParent();
	},

	getValue: function() {
		if(this.statemap_id) {
			return {type: 'requalificate', statemap: this.statemap_id};
		}
	},

	setValue: function(val) {
		if(val) {
			this.statemap_id = val;

			var selectmodel = this.getSelectionModel();
			var record = this.store.find('_id', val);
			selectmodel.select(record);
		}
	}
});
