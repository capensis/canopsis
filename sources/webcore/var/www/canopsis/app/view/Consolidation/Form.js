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
Ext.define('canopsis.view.Consolidation.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.ConsolidationForm',

	requires:['canopsis.lib.form.field.cthreshold_metro'],

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
						defaults: {anchor: '100%'},
						items: [
							{
								name: '_id',
								hidden: true
							},
							{
								xtype: 'fieldset',
								title: _('General'),
								defaultType: 'textfield',
								defaults: { labelWidth: labelWidth },
								items: [{
										fieldLabel: _('Display name'),
										name: 'crecord_name',
										allowBlank: false
									},{
										fieldLabel: _('Component'),
										name: 'component',
										allowBlank: false
									},{
										fieldLabel: _('Resource'),
										name: 'resource',
										allowBlank: false
									}]
							}
						]
					},{
						title: _('Consolidation'),
						bodyStyle: 'padding:5px 5px 0',
						layout: 'anchor',
						defaults: {anchor: '100%'},
						items:[{
							xtype: 'fieldset',
							title: _('Aggregation interval'),
							layout:'column',
							items: [{
								xtype: 'cduration',
								value: global.commonTs.minute,
								name: 'aggregation_interval',
								columnWidth: .50
							},{
								labelWidth: 70,
								xtype: "combobox",
								name: "aggregation_method",
								queryMode: "local",
								fieldLabel:'Operator',
								displayField: "text",
								valueField: "value",
								value: 'last',
								store: {
									xtype: "store",
									fields: ["value", "text"],
									data : [
										{"value": "last", "text": _("Last point")},
										{"value": "mean", "text": _("Mean")},
										{"value": "sum", "text": _("Sum")},
										{"value": "delta", "text": _("Delta")},
										{"value": "min", "text": _("Min")},
										{"value": "max", "text": _("Max")}
									]
								}
							}]
						},{
							xtype: 'fieldset',
							title: _('Consolidation operator'),
							items:[{
								xtype: 'checkboxgroup',
								name:'checkboxgroup',
								columns: 3,
	        					vertical: true,
								items:[{
									boxLabel  : _('Mean'),
									name      : 'consolidation_method',
									inputValue: 'mean'
								},{
									boxLabel  : _('Sum'),
									name 	  : 'consolidation_method',
									inputValue: 'sum'
								},{
									boxLabel  : _('Delta'),
									name      : 'consolidation_method',
									inputValue: 'delta'
								},{
									boxLabel  : _('Min'),
									name      : 'consolidation_method',
									inputValue: 'min'
								},{
									boxLabel  : _('Max'),
									name      : 'consolidation_method',
									inputValue: 'max'
								}]
							}]
						}]
						
					}/*,{
						title: _('Thresholds'),
						bodyStyle: 'padding:5px 5px 0',
						layout:'fit',
						items:[{
							xtype:'cthreshold_metro'
						}]
					}*/,{
						title: _('Filter'),
						xtype: 'cfilter',
						name: 'mfilter',
						url: "/perfstore",
						model:"canopsis.model.Perfdata",
						columns : [
							{
								header:"",
								sortable: false,
								flex: 2,
								dataIndex: "co",
							},{
								header:"",
								sortable: false,
								flex: 2,
								dataIndex:"re"
							}, {
								header:"",
								sortable:false,
								flex: 2,
								dataIndex:"me"
							}, {
								header:"",
								sortable: false,
								flex: 2,
								dataIndex:"u"
							}
						],
						operator_fields: [
							{ 'operator': 'co', 'text': _('Component'), 'type': 'all'},
							{ 'operator': 're', 'text': _('Resource'), 'type': 'all'},
							{ 'operator': 'me', 'text': _('Metric'), 'type': 'all'},
							{ 'operator': 'u', 'text': _('Unit'), 'type': 'all'},
							{ 'operator': 'tg', 'text':_('Tags'), 'type': 'all'}
						]
					}
				]
			}
		];

        this.callParent();
    }

});
