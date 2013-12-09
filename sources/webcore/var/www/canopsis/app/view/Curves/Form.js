//need:app/lib/view/cform.js
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
Ext.define('canopsis.view.Curves.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.CurvesForm',

	logAuthor: '[view][Curves][form]',
	layout: 'vbox',

	border: false,
	items: [
		{
			name: 'id',
			hidden: true
		},{
			fieldLabel: _('Metric name'),
			name: 'metric',
			allowBlank: false
		},{
			fieldLabel: _('Label'),
			name: 'label'
		},{
			xtype: 'combobox',
			name: 'dashStyle',
			fieldLabel: _('Dash style'),
			queryMode: 'local',
			displayField: 'text',
			editable: false,
			valueField: 'value',
			store: {
				xtype: 'store',
				fields: ['value', 'text'],
				data: [
					{value: 'Solid', text: _('Solid')},
					{value: 'ShortDash', text: _('ShortDash')},
					{value: 'ShortDot', text: _('ShortDot')},
					{value: 'ShortDashDot', text: _('ShortDashDot')},
					{value: 'ShortDashDotDot', text: _('ShortDashDotDot')},
					{value: 'Dot', text: _('Dot')},
					{value: 'Dash', text: _('Dash')},
					{value: 'LongDash', text: _('LongDash')},
					{value: 'DashDot', text: _('DashDot')},
					{value: 'LongDashDot', text: _('LongDashDot')},
					{value: 'LongDashDotDot', text: _('LongDashDotDot')}
				]
			}
		},{
			xtype: 'colorfield',
			colors: global.default_colors,
			fieldLabel: _('Line color'),
			name: 'line_color',
			allowBlank: false
		},{
			xtype: 'colorfield',
			colors: global.default_colors,
			fieldLabel: _('Area color'),
			name: 'area_color',
			allowBlank: true
		},{
			xtype: 'numberfield',
			name: 'area_opacity',
			fieldLabel: _('Area opacity') + ' (%)',
			minValue: 1,
			maxValue: 100,
			value: 75
		},{
			xtype: 'numberfield',
			name: 'zIndex',
			fieldLabel: _('zIndex'),
			value: 0
		},{
			xtype: 'checkboxfield',
			fieldLabel: _('Invert values'),
			name: 'invert',
			inputValue: true,
			uncheckedValue: false
		}
	]
});
