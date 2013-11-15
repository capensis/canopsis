//need:app/lib/form/field/cdate.js
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
Ext.define('canopsis.view.ReportingBar.ReportingBar', {
	extend: 'Ext.toolbar.Toolbar',

	requires: [
		'canopsis.lib.form.field.cdate'
	],

	alias: 'widget.ReportingBar',

	dock: 'top',

	//false to prevent reloading after choosing date/duration
	reloadAfterAction: false,

	initComponent: function() {
		this.callParent(arguments);

		this.advancedMode = false;

		// Create items

		var today = new Date();
		var tommorow = new Date(today.getTime() + (global.commonTs.day * 1000));

		this.previousButton = this.add({
			xtype: 'button',
			cls: 'x-btn-icon x-tbar-page-prev',
			action: 'previous'
		});

		this.textFor = this.add({
			xtype: 'tbtext',
			text: _('For') + ':'
		});

		var comboStore = Ext.create('Ext.data.Store', {
			fields: ['name', 'value'],
			data: [
				{'name': _('Day'), 'value': global.commonTs.day},
				{'name': _('Week'), 'value': global.commonTs.week},
				{'name': _('Month'), 'value': global.commonTs.month},
				{'name': _('Year'), 'value': global.commonTs.year}
			]
		});

		comboStore.load();

		this.periodNumber = this.add({
			xtype: 'numberfield',
			width: 55,
			value: 1,
			minValue: 1
		});

		this.combo = this.add({
			xtype: 'combobox',
			store: comboStore,
			queryMode: 'local',
			editable: false,
			displayField: 'name',
			width: 85,
			valueField: 'value',
			forceSelection: true,
			value: _('Day')
		});

		this.combo.setValue(86400);

		this.textFrom = this.add({xtype: 'tbtext', text: _('From') + ': ', hidden: true});
		this.fromTs = this.add({
			xtype: 'cdate',
			date_width: 130,
			hour_width: 70,
			date_value: today,
			max_value: tommorow,
			hidden: true
		});


		this.nextButton = this.add({
			xtype: 'button',
			cls: 'x-btn-icon x-tbar-page-next',
			action: 'next'
		});

		this.textTo = this.add({
			xtype: 'tbtext',
			text: _('To') + ': ',
			hidden: true
		});

		this.textBefore = this.add({xtype: 'tbtext', text: _('Before') + ': '});
		this.toTs = this.add({
			xtype: 'cdate',
			date_width: 130,
			hour_width: 70,
			now: true,
			max_value: tommorow
		});


		this.add('->');

		// Buttons

		this.add('-');

		this.toggleButton = this.add({
			xtype: 'button',
			iconCls: 'icon-calendar',
			action: 'toggleMode',
			tooltip: _('Toggle to advanced/simple mode')
		},{
			xtype: 'button',
			iconCls: 'icon-run',
			action: 'search',
			tooltip: _('Display data of the selected time')
		},{
			xtype: 'button',
			iconCls: 'icon-save',
			action: 'save',
			tooltip: _('Export this view to pdf')
		},{
			xtype: 'button',
			iconCls: 'icon-page-html',
			action: 'link',
			tooltip: _('View page in html')
		},{
			xtype: 'button',
			iconCls: 'icon-close',
			action: 'exit',
			tooltip: _('Leave reporting mode')
		});
	}
});
