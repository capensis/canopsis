//need:app/lib/form/field/cinventory.js
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

Ext.define('canopsis.view.MetricNavigation.MetricNavigation', {
	extend: 'Ext.panel.Panel',

	alias: 'widget.MetricNavigation',

	requires: [
		'canopsis.lib.form.field.cinventory'
	],

	border: false,

	nodes: [],

	initComponent: function() {
		//create cinventory
		this.metricTab = Ext.create('canopsis.lib.form.field.cinventory', {
			title: _('select metrics'),
			vertical_multiselect: true,
			default_padding: '3 0 3 0'
		});

		//create toolbar button
		var toolbar = Ext.create('Ext.toolbar.Toolbar');

		this.buttonCancel = toolbar.add({
			xtype: 'button',
			text: _('Cancel'),
			action: 'cancel'
		});

		toolbar.add('->');

		this.buttonDisplay = toolbar.add({
			xtype: 'button',
			text: _('Display'),
			action: 'display'
		});

		//time selection panel
		this.createTimePanel();

		//create config panel
		var config_tabPanel = {
			region: 'west',
			width: 550,
			border: false,
			title: _('Configuration panel'),
			collapsible: true,
			collapseDirection: 'left',
			bbar: toolbar,
			items: [this.metricTab, this.timePanel]
		};

		this.tabPanel = Ext.create('Ext.tab.Panel', config_tabPanel);

		//create render panel
		this.renderContent = Ext.create('Ext.panel.Panel', {
			layout: 'column',
			border: false,
			autoScroll: true
		});

		var config_renderPanel = {
			region: 'center',
			layout: 'fit',
			margin: '0 0 0 1',
			title: 'renderPanel',
			border: false,
			items: [this.renderContent]
		};

		this.renderPanel = Ext.create('Ext.panel.Panel', config_renderPanel);

		// building layout

		var masterpanel = Ext.create('Ext.panel.Panel', {
			layout: 'border',
			items: [this.tabPanel, this.renderPanel]
		});

		this.items = [masterpanel];
		this.callParent(arguments);
	},

	createTimePanel: function() {
		log.debug('Creating time panel', this.logAuthor);

		// fieldSet creation

		var fromField = Ext.widget('fieldset', {
			xtype: 'fieldset',
			title: _('From'),
			collapsible: false,
			layout: 'hbox'
		});

		var toField = Ext.widget('fieldset', {
			xtype: 'fieldset',
			title: _('To'),
			collapsible: false,
			layout: 'hbox'
		});

		// some default var
		var today = new Date();
		var tommorow = new Date(today.getTime() + (global.commonTs.day * 1000));
		var yesterday = new Date(today.getTime() - (global.commonTs.day * 1000));

		// from field
		this.fromHour = fromField.add({
			xtype: 'timefield',
			name: 'fromHour',
			fieldLabel: _('Hour'),
			increment: 15,
			labelWidth: 50,
			flex: 1,
			submitFormat: 'G:i',
			value: '0:00 AM'
		});

		this.fromDate = fromField.add({
			xtype: 'datefield',
			name: 'fromDate',
			editable: false,
			fieldLabel: _('Date'),
			labelWidth: 50,
			width: 110,
			margin: '0 0 0 10',
			value: yesterday,
			maxValue: tommorow,
			flex: 1
		});

		// to field
		this.toHour = toField.add({
			xtype: 'timefield',
			name: 'toHour',
			fieldLabel: _('Hour'),
			labelWidth: 50,
			increment: 15,
			flex: 1,
			submitFormat: 'G:i',
			value: '0:00 AM'
		});

		this.toDate = toField.add({
			xtype: 'datefield',
			name: 'toDate',
			editable: false,
			width: 110,
			value: today,
			fieldLabel: _('Date'),
			labelWidth: 50,
			margin: '0 0 0 10',
			maxValue: tommorow,
			flex: 1
		});

		// Button

		// Building
		var config = {
			title: _('Period selection'),
			items: [fromField, toField],
			border: false,
			padding: 10
		};

		this.timePanel = Ext.create('Ext.form.Panel', config);

		return this.timePanel;
	}
});
