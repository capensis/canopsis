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
Ext.define('canopsis.view.Aggregation.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.AggregationForm',

	defaultType: undefined,

	requires: [
		'canopsis.lib.form.field.cfieldset',
		'canopsis.lib.form.field.cinventory',
		'canopsis.lib.form.field.cfilter'
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
						/*defaults: {
							anchor: '100%'
						},*/
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
								items: [
									{
										fieldLabel: _('Name'),
										name: 'crecord_name',
										allowBlank: false
									}
								]
							},{
								xtype: 'fieldset',
								checkboxToggle: true,
								title: _('Aggregation option'),
								//checkboxName: 'dostate',
								//defaultType: 'textfield',
								value: true,
								//defaults: { labelWidth: labelWidth },
								items: [{
										xtype: 'combobox',
										name: 'state_algorithm',
										fieldLabel: _('Function'),
										queryMode: 'local',
										displayField: 'text',
										valueField: 'value',
										forceSelection: true,
										value: 0,
										store: {
											xtype: 'store',
											fields: ['value', 'text'],
											data: [
												{value: 'mean', text: _('Mean')}
											]
										}
									}
								]
							}
						]
					}/*,{
						title: _('Include'),
						name: 'include_ids',
						xtype: 'cinventory'
					},{
						title: _('Exclude'),
						name: 'exclude_ids',
						xtype: 'cinventory'
					}*/,{
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
