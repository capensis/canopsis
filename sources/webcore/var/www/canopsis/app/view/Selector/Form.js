//need:app/lib/view/cform.js,app/lib/form/field/cfieldset.js,app/lib/form/field/cinventory.js,app/lib/form/field/cfilter.js,app/lib/form/field/cduration.js
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
Ext.define('canopsis.view.Selector.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.SelectorForm',

	defaultType: undefined,

	requires: [
		'canopsis.lib.form.field.cfieldset',
		'canopsis.lib.form.field.cinventory',
		'canopsis.lib.form.field.cfilter',
		'canopsis.lib.form.field.cduration'
	],

	fieldDefaults: {
		labelWidth: 150
	},

	initComponent: function() {
		var labelWidth = 200;

		this.items = [
			{
				xtype: 'tabpanel',
				height: 400,
				width: 700,
				plain: true,
				border: false,
				defaults: {
					border: false,
					autoScroll: true
				},
				items: [
					{
						title: _('Options'),
						defaultType: 'textfield',
						bodyStyle: 'padding:5px 5px 0',
						layout: 'anchor',
						items: [
							{
								name: '_id',
								hidden: true
							},
							{
								xtype: 'cfieldset',
								title: _('General'),
								defaultType: 'textfield',
								defaults: {
									labelWidth: labelWidth
								},
								items: [
									{
										fieldLabel: _('Name'),
										name: 'crecord_name',
										allowBlank: false
									},{
										fieldLabel: _('Display name'),
										name: 'display_name'
									},{
										fieldLabel: _('Description'),
										xtype: 'textareafield',
										name: 'description'
									}
								]
							},{
								xtype: 'cfieldset',
								title: _('Calcul State'),
								checkboxName: 'dostate',
								defaultType: 'textfield',
								value: true,
								defaults: {
									labelWidth: labelWidth
								},
								items: [{
										xtype: 'combobox',
										name: 'state_algorithm',
										fieldLabel: _('Algorithm'),
										queryMode: 'local',
										displayField: 'text',
										valueField: 'value',
										forceSelection: true,
										value: 0,
										store: {
											xtype: 'store',
											fields: ['value', 'text'],
											data: [
												{value: 0, text: _('Worst state')}
											]
										}
									},{
										fieldLabel: _('Output Template'),
										xtype: 'textareafield',
										name: 'output_tpl',
										anchor: '100%',
										value: '{cps_sel_state_0} Ok, {cps_sel_state_1} Warning, {cps_sel_state_2} Critical'
									}
								]
							},{
								xtype: 'cfieldset',
								title: _('Calcul SLA (if state is calculated)'),
								checkboxName: 'dosla',
								defaultType: 'textfield',
								value: false,
								defaults: {
									labelWidth: labelWidth
								},
								items: [
									{
										xtype: 'fieldcontainer',
										fieldLabel: _('Time window'),
										layout: 'hbox',
										width: 500,
										items: [
											{
												xtype: 'numberfield',
												name: 'sla_timewindow_value',
												minValue: 1,
												value: 1,
												width: 60,
												allowBlank: false,
												padding: '0 5 0 0'
											},{
												xtype: 'combobox',
												name: 'sla_timewindow_unit',
												queryMode: 'local',
												displayField: 'text',
												width: 90,
												valueField: 'value',
												value: global.commonTs.day,
												store: {
													xtype: 'store',
													fields: ['value', 'text'],
													data: [
														{value: global.commonTs.day, text: _('Day')},
														{value: global.commonTs.week, text: _('Week')},
														{value: global.commonTs.month, text: _('Month')},
														{value: global.commonTs.year, text: _('Year')}
													]
												}
											}
										]
									},{
										fieldLabel: _('Consider unknown time'),
										xtype: 'checkboxfield',
										inputValue: true,
										uncheckedValue: false,
										name: 'sla_timewindow_doUnknown'
									},{
										xtype: 'numberfield',
										fieldLabel: _('Warning threshold'),
										name: 'thd_warn_sla_timewindow',
										minValue: 1,
										maxValue: 100,
										value: 98,
										allowBlank: true
									},{
										xtype: 'numberfield',
										fieldLabel: _('Critical threshold'),
										name: 'thd_crit_sla_timewindow',
										minValue: 1,
										maxValue: 100,
										value: 95,
										allowBlank: true
									},{
										fieldLabel: _('Output Template'),
										xtype: 'textareafield',
										name: 'sla_output_tpl',
										anchor: '100%',
										value: '{cps_pct_by_state_0}% Ok, {cps_pct_by_state_1}% Warning, {cps_pct_by_state_2}% Critical, {cps_pct_by_state_3}% Unknown'
									}
								]
							}
						]
					},{
						title: _('Include'),
						name: 'include_ids',
						xtype: 'cinventory'
					},{
						title: _('Exclude'),
						name: 'exclude_ids',
						xtype: 'cinventory'
					},{
						title: _('Filter'),
						xtype: 'cfilter',
						name: 'mfilter'
					}
				]
			}
		];

		this.callParent();
	}
});
