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
	extend: 'Ext.panel.Panel',

	requires: [
		'canopsis.lib.form.field.cdate'
	],

	alias: 'widget.ReportingBar',

	layout: {
		type: 'vbox',
		align: 'stretch'
	},

	//false to prevent reloading after choosing date/duration
	reloadAfterAction: false,

	initComponent: function() {
		this.callParent(arguments);

		this.advancedMode = false;
		this.toolbar = this.add({
			xtype:'toolbar',
			dock: 'top'
		});
		// Create items

		var today = new Date();
		var tommorow = new Date(today.getTime() + (global.commonTs.day * 1000));

		this.previousButton = this.toolbar.add({
			xtype: 'button',
			cls: 'x-btn-icon x-tbar-page-prev',
			action: 'previous'
		});

		this.textFor = this.toolbar.add({
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

		this.periodNumber = this.toolbar.add({
			xtype: 'numberfield',
			width: 55,
			value: 1,
			minValue: 1
		});

		this.combo = this.toolbar.add({
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

		this.buttonExpandAdvancedMode = this.toolbar.add({
			xtype: 'button',
			enableToggle: true,
			cls: 'x-btn-icon x-group-by-icon',
			action: 'toggleAdvancedFilters',
			hidden: true
		});

		this.textFrom = this.toolbar.add({xtype: 'tbtext', text: _('From') + ': ', hidden: true});
		this.fromTs = this.toolbar.add({
			xtype: 'cdate',
			date_width: 130,
			hour_width: 70,
			date_value: today,
			max_value: tommorow,
			hidden: true
		});


		this.nextButton = this.toolbar.add({
			xtype: 'button',
			cls: 'x-btn-icon x-tbar-page-next',
			action: 'next'
		});

		this.textTo = this.toolbar.add({
			xtype: 'tbtext',
			text: _('To') + ': ',
			hidden: true
		});

		this.textBefore = this.toolbar.add({xtype: 'tbtext', text: _('Before') + ': '});
		this.toTs = this.toolbar.add({
			xtype: 'cdate',
			date_width: 130,
			hour_width: 70,
			now: true,
			max_value: tommorow
		});


		this.toolbar.add('->');

		// Buttons

		this.toolbar.add('-');

		this.toggleButton = this.toolbar.add({
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

		// this.advancedFilters.cfilter = Ext.create("widget.cfilter", {});

		this.advancedFilters = this.add({
			xtype: 'ccard',
			hidden: true,
			height:300,
			wizardSteps:[{
				title:"Component/Resource",
				items: [{
					xtype:"cgrid",
					itemId:"componentResourceGrid",
					store: {
						xtype: "store",
						fields: ["component", "resource"],
						reader: {
							type: 'json'
						}
					},
					columns:[{
						header: "Component",
						sortable: false,
						dataIndex: "component",
						editor: "field",
						// renderer: rdr_tstodate,
						flex: 3
					},{
						header: "Resource",
						sortable: false,
						dataIndex: "resource",
						editor: "field",
						// renderer: rdr_tstodate,
						flex: 3
					}],
					height:100,
					opt_bar_reload: false,
			        queryMode: 'local',
			        opt_keynav_del : true,
					opt_menu_delete: true,
			        opt_bar_delete : true,
					opt_bar_enable: false,
					opt_confirmation_delete: false
				}],
			},{
				title:"Hostgroups",
				items: [{
					xtype:"cgrid",
					itemId:"hostgroupsGrid",
					store: {
						xtype: "store",
						fields: ["hostgroup"],
						reader: {
							type: 'json'
						}
					},
					columns:[{
						header: "Hostgroup",
						sortable: false,
						dataIndex: "hostgroup",
						editor: "field",
						// renderer: rdr_tstodate,
						flex: 3
					}],
					height:100,
					opt_bar_reload: false,
			        queryMode: 'local',
			        opt_keynav_del : true,
					opt_menu_delete: true,
			        opt_bar_delete : true,
					opt_bar_enable: false,
					opt_confirmation_delete: false
				}],
			},{
				title:"Exclusion Intervals",
				items: [{
					xtype:"cgrid",
					itemId:"exclusionIntervalGrid",
					store: {
						xtype: "store",
						fields: ["from", "to"],
						reader: {
							type: 'json'
						}
					},
					columns:[{
						header: "from",
						sortable: false,
						dataIndex: "from",
						editor: "field",
						renderer: rdr_tstodate,
						flex: 3
					},{
						header: "to",
						sortable: false,
						dataIndex: "to",
						editor: "field",
						renderer: rdr_tstodate,
						flex: 3
					}],
					height:100,
					opt_bar_reload: false,
			        queryMode: 'local',
			        opt_keynav_del : true,
					opt_menu_delete: true,
			        opt_bar_delete : true,
					opt_bar_enable: false,
					opt_confirmation_delete: false
				}]
			},{
				title:"Downtimes",
				items: [{

					xtype:"cinventory",
					name : "Downtimes",
					multiSelect: true,
					vertical_multiselect: true,
				}],
			}]
		});

		this.addExclusionIntervalWindow = Ext.create('Ext.window.Window', {
			closeAction:'hide',
			cls: 'addExclusionIntervalWindow',
			modal:true,
			items:[{
				xtype: "panel",
				items:[{
					xtype: "cdate",
					itemId: "newExclusionInterval_from",
					label_text: "From"
				},{
					xtype: "cdate",
					itemId: "newExclusionInterval_to",
					label_text: "To"
				},{
					xtype: "button",
					text: "Save",
					action: "addExclusionInterval"
				}]
			}]
		});

		this.addComponentResourceWindow = Ext.create('Ext.window.Window', {
			closeAction:'hide',
			cls: 'addComponentResourceWindow',
			modal:true,
			items:[{
				xtype: "panel",
				items:[{
					xtype: "textfield",
					itemId: "component",
					fieldLabel: "Component"
				},{
					xtype: "textfield",
					itemId: "resource",
					fieldLabel: "Resource"
				},{
					xtype: "button",
					text: "Save",
					action: "addComponentResource"
				}]
			}]
		});

		this.addHostgroupWindow = Ext.create('Ext.window.Window', {
			closeAction:'hide',
			cls: 'addHostgroupWindow',
			modal:true,
			items:[{
				xtype: "panel",
				items:[{
					xtype: "textfield",
					itemId: "hostgroup",
					fieldLabel: "Hostgroup"
				},{
					xtype: "button",
					text: "Save",
					action: "addHostgroup"
				}]
			}]
		});

		// this.advancedFilters.cfilterTab.cfilter = this.advancedFilters.cfilterTab.add({xtype:cfilter});
	}
});
