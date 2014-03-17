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
Ext.define('canopsis.controller.ReportingBar', {
	extend: 'Ext.app.Controller',

	views: ['ReportingBar.ReportingBar'],
	logAuthor: '[controller][ReportingBar]',

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.control({
			'ReportingBar': {
				afterrender: this._bindBarEvents
			},
			'ReportingBar button[action="toggleMode"]': {
				click: this.toggle_mode
			},
			'ReportingBar button[action="search"]': {
				click: this.launchReport
			},
			'ReportingBar button[action="save"]': {
				click: this.saveButton
			},
			'ReportingBar button[action="link"]': {
				click: this.htmlReport
			},
			'ReportingBar button[action="exit"]': {
				click: this.exitButton
			},
			'ReportingBar button[action="next"]': {
				click: this.nextButton
			},
			'ReportingBar button[action="previous"]': {
				click: this.previousButton
			},
			'ReportingBar button[action="toggleAdvancedFilters"]': {
				click: this.toggleAdvancedFilters
			},
			'ReportingBar button[action="computeAdvancedFilters"]': {
				click: this.getAdvancedFilters
			},
			'ReportingBar cgrid#hostgroupsGrid button[action="add"]': {
				click: this.showAddHostgroupWindow
			},
			'ReportingBar cgrid#exclusionIntervalGrid button[action="add"]': {
				click: this.showAddExclusionIntervalWindow
			},
			'window[cls=addExclusionIntervalWindow] button[action="addExclusionInterval"]': {
				click: this.addExclusionInterval
			},
			'window[cls=addHostgroupWindow] button[action="addHostgroup"]': {
				click: this.addHostgroup
			}
		});

		this.ctrl = Ext.create('canopsis.lib.controller.cgrid');

		this.callParent(arguments);
	},

	_bindBarEvents: function(bar) {
		log.debug('Bind events...', this.logAuthor);
		this.bar = bar;

		bar.toTs.on('select', this.setMaxDate, this);
		bar.fromTs.on('select', this.setMinDate, this);
	},

	launchReport: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();

		var timestamps = this.getReportTime();
		var startTimestamp = timestamps.start;
		var stopTimestamp = timestamps.stop;

		var advancedFilters = this.getAdvancedFilters();

		if(startTimestamp && stopTimestamp) {
			log.debug('------------------------Asked Report date-----------------------');
			log.debug('from : ' + startTimestamp + ' To : ' + stopTimestamp, this.logAuthor);
			log.debug('startReport date is : ' + Ext.Date.format(new Date(startTimestamp * 1000), 'Y-m-d H:i:s'), this.logAuthor);
			log.debug('endReport date is : ' + Ext.Date.format(new Date(stopTimestamp * 1000), 'Y-m-d H:i:s'), this.logAuthor);
			log.dump(advancedFilters);
			log.debug('----------------------------------------------------------------');
			tab.setReportDate(startTimestamp * 1000, stopTimestamp * 1000, advancedFilters);
		}
		else {
			log.debug('Timestamps are, start: ' + startTimestamp + ' stop: ' + stopTimestamp, this.logAuthor);
			global.notify.notify(_('Invalid date'), _('The selected date is invalid'));
		}
	},

	nextButton: function() {
		log.debug('Next button pressed', this.logAuthor);
		var dateField = this.bar.toTs;

		var selectedTime = dateField.getValue();
		var timeUnit = this.bar.combo.getValue();

		log.debug('selected time : ' + selectedTime);
		log.debug('time unit : ' + timeUnit);

		var timestamp = selectedTime + (timeUnit * this.bar.periodNumber.getValue());
		dateField.setValue(timestamp);

		if(dateField.isValid()) {
			this.launchReport();
		}
		else {
			global.notify.notify(_('Invalid date'), _('The selected date is invalid'));
		}
	},

	previousButton: function() {
		log.debug('Previous button pressed', this.logAuthor);
		var dateField = this.bar.toTs;

		var selectedTime = dateField.getValue();
		var timeUnit = this.bar.combo.getValue();

		var timestamp = selectedTime - (timeUnit * this.bar.periodNumber.getValue());
		dateField.setValue(timestamp);

		if (dateField.isValid()) {
			this.launchReport();
		}
		else {
			global.notify.notify(_('Invalid date'), _('The selected date is invalid'));
		}
	},

	saveButton: function() {
		log.debug('launching pdf reporting', this.logAuthor);

		var timestamps = this.getReportTime();
		var startTimestamp = timestamps.start;
		var stopTimestamp = timestamps.stop;

		if(startTimestamp && stopTimestamp) {
			var view_id = Ext.getCmp('main-tabs').getActiveTab().view_id;
			var ctrl = this.getController('Reporting');

			log.debug('view_id : ' + view_id, this.logAuthor);
			log.debug('startReport : ' + startTimestamp * 1000, this.logAuthor);
			log.debug('stopReport : ' + stopTimestamp * 1000, this.logAuthor);

			ctrl.launchReport(view_id, startTimestamp * 1000, stopTimestamp * 1000);
		}
		else {
			log.debug('Timestamps are, start: ' + startTimestamp + ' stop: ' + stopTimestamp, this.logAuthor);
			global.notify.notify(_('Invalid date'), _('The selected date is in futur'));
		}
	},

	htmlReport: function() {
		log.debug('launching html window reporting', this.logAuthor);

		var timestamps = this.getReportTime();
		var startTimestamp = timestamps.start;
		var stopTimestamp = timestamps.stop;

		if(startTimestamp && stopTimestamp) {
			var ctrl = this.getController('Reporting');
			var view = Ext.getCmp('main-tabs').getActiveTab().view_id;
			ctrl.openHtmlReport(view, startTimestamp * 1000, stopTimestamp * 1000);
		}
		else {
			log.debug('Timestamps are, start: ' + startTimestamp + ' stop: ' + stopTimestamp, this.logAuthor);
			global.notify.notify(_('Invalid date'), _('The selected date is in futur'));
		}
	},

	getReportTime: function() {
		var startTimestamp = undefined;
		var stopTimestamp  = undefined;

		if(this.bar.advancedMode) {
			startTimestamp = this.bar.fromTs.getValue();
			stopTimestamp  = this.bar.toTs.getValue();
		}
		else {
			var timeUnit     = this.bar.combo.getValue();
			var periodLength = this.bar.periodNumber.getValue();

			stopTimestamp  = this.bar.toTs.getValue();
			startTimestamp = stopTimestamp - (timeUnit * periodLength);
		}

		return {
			start: startTimestamp,
			stop: stopTimestamp
		};
	},

	getTimestamp: function(date_element,hour_element) {
		var date = date_element;
		var hour = hour_element;

		if(date.isValid() && hour.isValid()) {
			var tsDate = parseInt(Ext.Date.format(date.getValue(), 'U'));
			var hourObject = stringTo24h(hour.getValue());

			//date + hour in seconds + minute in second
			var timestamp = tsDate + (hourObject.hour * 60 * 60) + (hourObject.minute * 60);
		}
		else {
			log.debug('getTimestamp Invalid', this.logAuthor);
			return undefined;
		}

		return parseInt(timestamp, 10);
	},

	exitButton: function() {
		log.debug('Exit reporting mode', this.logAuthor);
		var tab = Ext.getCmp('main-tabs').getActiveTab();
		tab.report_window.destroy();
		tab.report_window = undefined;
		this.getController('Tabs').reload_active_view();
	},

	enable_reporting_mode: function() {
		log.debug('Enable reporting mode', this.logAuthor);
		Ext.getCmp('main-tabs').getActiveTab().addReportingBar();
	},

	toggle_mode: function() {
		if(this.bar.advancedMode) {
			this.bar.fromTs.hide();
			this.bar.textFrom.hide();
			this.bar.textTo.hide();
			this.bar.textFor.show();
			this.bar.textBefore.show();
			this.bar.previousButton.show();
			this.bar.nextButton.show();
			this.bar.periodNumber.show();
			this.bar.combo.show();
			this.bar.buttonExpandAdvancedMode.hide();
			this.bar.advancedMode = false;
		}
		else {
			this.bar.fromTs.show();
			this.bar.textFrom.show();
			this.bar.textTo.show();
			this.bar.textFor.hide();
			this.bar.textBefore.hide();
			this.bar.previousButton.hide();
			this.bar.nextButton.hide();
			this.bar.periodNumber.hide();
			this.bar.combo.hide();
			this.bar.buttonExpandAdvancedMode.show();
			this.bar.advancedMode = true;
		}

	},

	toggleAdvancedFilters: function() {
		if(this.advancedFiltersShown === true)
		{
			this.bar.advancedFilters.hide();
			this.advancedFiltersShown = false;
			this.bar.setHeight(this.bar.toolbar.getHeight());
		}
		else
		{
			this.bar.advancedFilters.show();
			this.advancedFiltersShown = true;
			this.bar.setHeight(300);
		}
	},

	showAddExclusionIntervalWindow: function() {
		log.debug("showAddExclusionIntervalWindow");
		this.bar.addExclusionIntervalWindow.show();
	},

	showAddHostgroupWindow: function() {
		console.log("showAddHostgroupWindow");
		this.bar.addHostgroupWindow.show();
	},

	addExclusionInterval: function() {
		log.debug("new exclusion interval");
		var from = this.bar.addExclusionIntervalWindow.down("#newExclusionInterval_from").getValue();
		var to = this.bar.addExclusionIntervalWindow.down("#newExclusionInterval_to").getValue();

		this.bar.addExclusionIntervalWindow.hide();

		var grid = this.bar.down("#exclusionIntervalGrid");
		grid.store.add({from: from, to: to});
	},

	addHostgroup: function() {
		console.log("new hostgroup");
		var hostgroup = this.bar.addHostgroupWindow.down("#hostgroup").getValue();

		this.bar.addHostgroupWindow.hide();

		var grid = this.bar.down("#hostgroupsGrid");
		grid.store.add({hostgroup: hostgroup});
	},

	setMinDate: function(cdate, date) {
		void(cdate);

		this.bar.toTs.setMinDate(date);
	},

	setMaxDate: function(cdate, date) {
		void(cdate);

		this.bar.fromTs.setMaxDate(date);
	},

	getAdvancedFilters: function() {
		return {
			component_resources: this.computeComponentResource(),
			exclusions: this.computeExclusionFilter(),
			hostgroups: this.computeHostgroups(),
			downtimes: 	this.computeDowntimes(),
		};
	},

	computeExclusionFilter: function() {
		var result = [];
		var grid = this.bar.down("#exclusionIntervalGrid");

		for (var i = grid.store.data.items.length - 1; i >= 0; i--) {
			var from = grid.store.data.items[i].data.from;
			var to = grid.store.data.items[i].data.to;
			result.push({"from": from, "to": to});
		};

		return result;
	},

	computeComponentResource: function() {
		var result = [];

		var inventory_components_resources = this.bar.down("#component_resource");
		var components_resources = inventory_components_resources.getValue();

		for (var i = inventory_components_resources.selection_store.data.items.length - 1; i >= 0; i--) {
			var component = inventory_components_resources.selection_store.data.items[i].data.component;
			var resource = inventory_components_resources.selection_store.data.items[i].data.resource;

			result.push({"component": component, "resource": resource});
		};
		return result;
	},

	computeHostgroups: function() {
		var result = [];

		grid = this.bar.down("#hostgroupsGrid");

		for (var i = grid.store.data.items.length - 1; i >= 0; i--) {
			var hostgroup = grid.store.data.items[i].data.hostgroup;
			result.push(hostgroup);
		};

		return result;
	},

	computeDowntimes: function() {
		var result = [];

		var inventory_downtimes = this.bar.down("#Downtimes");

		for (var i = inventory_downtimes.selection_store.data.items.length - 1; i >= 0; i--) {
			var component = inventory_downtimes.selection_store.data.items[i].data.component;
			var resource = inventory_downtimes.selection_store.data.items[i].data.resource;

			result.push({"component": component, "resource": resource});
		};

		return result;
	}
});
