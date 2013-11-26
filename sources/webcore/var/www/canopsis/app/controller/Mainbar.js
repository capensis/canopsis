//need:app/view/Mainbar/Bar.js,app/lib/view/cauthkey.js,app/model/Schedule.js
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
Ext.define('canopsis.controller.Mainbar', {
	extend: 'Ext.app.Controller',

	requires: [
		'canopsis.lib.view.cauthkey',
		'canopsis.model.Schedule'
	],

	views: ['Mainbar.Bar'],

	logAuthor: '[controller][Mainbar]',

	init: function() {
		this.control({
			'#region-north' : {
				collapse: this.onCollapseMainbar,
				expand: this.onExpandMainbar
			},
			'Mainbar [action="logout"]': {
				click: this.logout
			},
			'Mainbar menuitem[action="cleartabscache"]': {
				click: this.cleartabscache
			},
			'Mainbar menuitem[action="authkey"]': {
				click: this.authkey
			},
			'Mainbar combobox[action="viewSelector"]': {
				select: this.openViewSelector
			},
			'Mainbar combobox[action="dashboardSelector"]': {
				select: this.setDashboard
			},
			'Mainbar combobox[action="localeSelector"]': {
				select: this.setLocale
			},
			'Mainbar combobox[action="clockTypeSelector"]': {
				select: this.setClockType
			},
			'Mainbar combobox[action="avatarSelector"]': {
				select: this.setAvatar
			},
			'Mainbar menuitem[action="openDashboard"]': {
				click: this.openDashboard
			},
			'Mainbar menuitem[action="editView"]': {
				click: this.editView
			},
			'Mainbar menuitem[action="editAccount"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="editSelector"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="openViewsManager"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="editGroup"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="editSchedule"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="openBriefcase"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="newView"]': {
				click: this.newView
			},
			'Mainbar menuitem[action="exportView"]': {
				click: this.exportView
			},
			'Mainbar menuitem[action="reportingMode"]': {
				click: this.reportingMode
			},
			'Mainbar menuitem[action="eventLog_navigation"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="ScheduleExportView"]': {
				click: this.ScheduleExportView
			},
			'Mainbar [name="clock"]': {
				afterrender: this.setClock
			},
			'Mainbar menuitem[action="openViewMenu"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="openDerogationManager"]': {
				click: this.openViewMenu
			},
			'Mainbar menuitem[action="openStatemapManager"]' : {
				click: this.openViewMenu
			}
		});

		this.callParent(arguments);

		// Bind Websocket Events
		global.websocketCtrl.on('transport_up', function() {
			var button = Ext.getCmp('Mainbar-menu-Websocket');

			if(button) {
				button.setIconCls('icon-bullet-green');
			}
		}, this);

		global.websocketCtrl.on('transport_down', function() {
			var button = Ext.getCmp('Mainbar-menu-Websocket');

			if(button) {
				button.setIconCls('icon-bullet-red');
			}
		}, this);
	},

	onCollapseMainbar: function() {
	},

	onExpandMainbar: function() {
	},

	logout: function() {
		log.debug('Logout', this.logAuthor);
		this.getController('Account').logout();
	},

	cleartabscache: function() {
		log.debug('Clear tabs localstore', this.logAuthor);
		this.getController('Tabs').clearTabsCache();
	},

	authkey: function() {
		log.debug('Show authkey', this.logAuthor);
		var authkey = Ext.create('canopsis.lib.view.cauthkey');
		authkey.show();
	},

	openViewMenu: function(item) {
		this.getController('Tabs').open_view({
			view_id: item.viewId,
			title: _(item.text)
		});
	},

	setClock: function(item) {
		log.debug('Set Clock', this.logAuthor);

		var refreshClock = function() {
			var thisTime = new Date();

			if(is12Clock()) {
				item.update("<div class='cps-account' >" + Ext.Date.format(thisTime, 'g:i a - l d F Y') + '</div>');
			}
			else {
				item.update("<div class='cps-account' >" + Ext.Date.format(thisTime, 'G:i - l d F Y') + '</div>');
			}
		};

		Ext.TaskManager.start({
			run: refreshClock,
			interval: 60000
		});
	},

	setLocale: function(combo, records) {
		void(combo);

		var locale = records[0].get('value');
		log.debug('Set language to ' + locale, this.logAuthor);
		this.getController('Account').setLocale(locale);
	},

	setClockType: function(combo, records) {
		void(combo);

		var clock_type = records[0].get('value');
		log.debug('Set clock to ' + clock_type, this.logAuthor);
		this.getController('Account').setClock(clock_type);
	},

	setAvatar: function(combo, records) {
		var avatarId = records[0].data.id;
		var avatarFilename = records[0].data.file_name;

		log.debug('Set avatar to ' + avatarId + '(' + avatarFilename + ')', this.logAuthor);

		combo.up('button').hideMenu();

		global.accountCtrl.setAvatar(avatarId, avatarFilename);
	},

	setDashboard: function(combo, records) {
		void(combo);

		var view_id = records[0].get('id');
		log.debug('Set dashboard to ' + view_id, this.logAuthor);

		//set new dashboard
		this.getController('Account').setConfig('dashboard', view_id);

		var maintabs = Ext.getCmp('main-tabs');

		//close view selected if open
		var tab = Ext.getCmp(view_id + '.tab');
		if(tab) {
			maintabs.remove(tab.id);
		}

		var current_dashboard_id = maintabs.getComponent(0).id;

		this.getController('Tabs').open_dashboard();

		//close current dashboard
		maintabs.remove(current_dashboard_id);
	},

	openDashboard: function() {
		log.debug('Open dashboard', this.logAuthor);
		var maintabs = Ext.getCmp('main-tabs');
		maintabs.setActiveTab(0);
	},

	openViewSelector: function(combo, records) {
		var view_id = records[0].get('id');
		var view_name = records[0].get('crecord_name');

		log.debug('Open view "' + view_name + '" (' + view_id + ')', this.logAuthor);
		combo.clearValue();

		this.getController('Tabs').open_view({
			view_id: view_id,
			title: view_name
		});
	},

	editView: function() {
		log.debug('Edit view', this.logAuthor);
		var ctrl = this.getController('Tabs');
		ctrl.edit_active_view();
	},

	reportingMode: function() {
		log.debug('Live reporting mode activated', this.logAuthor);
		var ctrl = this.getController('ReportingBar');
		ctrl.enable_reporting_mode();
	},

	newView: function() {
		log.debug('New view', this.logAuthor);
		var ctrl = this.getController('Tabs');
		ctrl.create_new_view();
	},

	exportView: function() {
		var view_id = Ext.getCmp('main-tabs').getActiveTab().view_id;
		this.getController('Reporting').launchReport(view_id);
	},

	ScheduleExportView: function() {
		log.debug('Schedule active view export', this.logAuthor);
		var activetabs = Ext.getCmp('main-tabs').getActiveTab();
		var record = Ext.create('canopsis.model.Schedule', activetabs.view);
		this.getController('Schedule').scheduleWizard(record, activetabs.id);
	}
});
