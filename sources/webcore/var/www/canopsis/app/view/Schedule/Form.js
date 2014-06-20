//need:app/lib/view/cform.js,app/lib/form/field/cfieldset.js
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
Ext.define('canopsis.view.Schedule.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.ScheduleForm',

	requires: [
		'canopsis.lib.form.field.cfieldset'
	],

	items: [
		{
			xtype: 'ccard',
			height: 350,
			width: 600,
			resizable: true,
			wizardSteps: [
				{
					title: _('General options'),
					items: [
						{
							xtype: 'textfield',
							fieldLabel: _('Schedule name'),
							name: 'crecord_name',
							allowBlank: false
						}, {
							xtype: 'combobox',
							fieldLabel: _('Action'),
							name: 'task',
							queryMode: 'local',
							displayField: 'text',
							valueField: 'value',
							value: 'task_reporting.render_pdf',
							disabled: true,
							store: {
								xtype: 'store',
								fields: ['value', 'text'],
								data: [
									{value: 'task_reporting.render_pdf', text: _('Reporting')}
								]
							}
						}, {
							xtype: 'combobox',
							fieldLabel: _('View'),
							name: 'exporting_viewName',
							store: 'Views',
							displayField: 'crecord_name',
							valueField: 'id',
							typeAhead: false,
							allowBlank: false,
							minChars: 2,
							queryMode: 'remote',
							emptyText: _('Select a view') + ' ...'
						}, {
							xtype: 'combobox',
							fieldLabel: _('Owner'),
							name: 'exporting_owner',
							store: 'Accounts',
							displayField: 'user',
							valueField: 'user',
							typeAhead: false,
							allowBlank: false,
							queryMode: 'remote'
						}
					]
				}, {
					title: _('Frequency'),
					items: [
						{
							xtype: 'combobox',
							name: 'frequency',
							fieldLabel: _('Every'),
							queryMode: 'local',
							displayField: 'text',
							valueField: 'value',
							value: 'day',
							store: {
								xtype: 'store',
								fields: ['value', 'text'],
								data: [
									{value: 'day', text: _('Day')},
									{value: 'week', text: _('Week')},
									{value: 'month', text: _('Month')},
									{value: 'year', text: _('Year')}
								]
							}
						},{
							xtype: 'combobox',
							name: 'crontab_month',
							fieldLabel: _('Month'),
							queryMode: 'local',
							displayField: 'text',
							valueField: 'value',
							value: 1,
							disabled: true,
							hidden: true,
							store: {
								xtype: 'store',
								fields: ['value', 'text'],
								data: [
									{value: 1, text: _('January')},
									{value: 2, text: _('February')},
									{value: 3, text: _('March')},
									{value: 4, text: _('April')},
									{value: 5, text: _('May')},
									{value: 6, text: _('June')},
									{value: 7, text: _('July')},
									{value: 8, text: _('August')},
									{value: 9, text: _('September')},
									{value: 10, text: _('October')},
									{value: 11, text: _('November')},
									{value: 12, text: _('December')}
								]
							}
						}, {
							xtype: 'combobox',
							name: 'crontab_day',
							fieldLabel: _('Day'),
							queryMode: 'local',
							displayField: 'text',
							valueField: 'value',
							value: 1,
							disabled: true,
							hidden: true,
							store: {
								xtype: 'store',
								fields: ['value', 'text'],
								data: (function() {
										var dayData = [];

										for(var i = 1; i < 32; i++) {
											dayData.push({value: i, text: i});
										}

										return dayData;
									})()
							}
						}, {
							xtype: 'combobox',
							name: 'crontab_day_of_week',
							fieldLabel: _('Day of week'),
							queryMode: 'local',
							displayField: 'text',
							valueField: 'value',
							value: 'mon',
							disabled: true,
							hidden: true,
							store: {
								xtype: 'store',
								fields: ['value', 'text'],
								data: [
									{value: 'mon', text: _('Monday')},
									{value: 'tue', text: _('Tuesday')},
									{value: 'wed', text: _('Wednesday')},
									{value: 'thu', text: _('Thursday')},
									{value: 'fri', text: _('Friday')},
									{value: 'sat', text: _('Satursday')},
									{value: 'sun', text: _('Sunday')}
								]
							}
						}, {
							xtype: 'textfield',
							name: 'crontab_hours',
							fieldLabel: _('Hours (local time)'),
							allowBlank: false,
							value: '00:00',
							regex: getTimeRegex()
						}
					]
				}, {
					title: _('Exporting interval'),
					items: [
						{
							xtype: 'checkbox',
							name: 'exporting_enable',
							checked: false,
							fieldLabel: _('Exporting interval enable')
						},
						{
							xtype: 'cfieldset',
							title: _('Time zone'),
							itemId: 'exporting_timezone',
							layout: "hbox",
							hidden: true,
							disabled: true,
							items: [
								{
									xtype: 'combobox',
									name: 'timezone_type',
									fieldLabel: _('Type'),
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									value: 'local',
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 'local', text: _('Local')},
											{value: 'utc', text: _('UTC')},
											{value: 'server', text: _('Server')},
											{value: 'custom', text: _('Custom (in minutes)')}
										]
									}
								}, {
									xtype: 'numberfield',
									hidden: true,
									disabled: true,
									name: 'timezone_value',
									value: new Date().getTimezoneOffset() * 60
								}
							]
						}, {
							xtype: "checkbox",
							fieldLabel: _('advanced'),
							checked: false,
							name: 'exporting_advanced',
							hidden: true,
							disabled: true,
						}, {
							xtype: 'cfieldset',
							layout: "hbox",
							title: _('Duration'),
							itemId: 'exporting_duration',
							hidden: true,
							disabled: true,
							items: [
								{
									xtype: "numberfield",
									name: "exporting_intervalLength",
									value: 1,
									minValue: 0
								},	{
									xtype: 'combobox',
									name: 'exporting_intervalUnit',
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									value: 'days',
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 'hours', text: _('Hours')},
											{value: 'days', text: _('Day')},
											{value: 'weeks', text: _('Week')},
											{value: 'months', text: _('Month')},
											{value: 'years', text: _('Year')}
										]
									}
								}
							]
						},  {
							xtype: 'cfieldset',
							title: _('From'),
							layout: 'vbox',
							itemId: 'from',
							disabled: true,
							hidden: true,
							items: [
								{
									xtype: 'combobox',
									name: 'from_month',
									fieldLabel: _('Month'),
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									value: 1,
									disabled: true,
									hidden: true,
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 1, text: _('January')},
											{value: 2, text: _('February')},
											{value: 3, text: _('March')},
											{value: 4, text: _('April')},
											{value: 5, text: _('May')},
											{value: 6, text: _('June')},
											{value: 7, text: _('July')},
											{value: 8, text: _('August')},
											{value: 9, text: _('September')},
											{value: 10, text: _('October')},
											{value: 11, text: _('November')},
											{value: 12, text: _('December')}
										]
									}
								}, {
									xtype: 'combobox',
									name: 'from_day',
									fieldLabel: _('Day'),
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									value: 1,
									disabled: true,
									hidden: true,
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: (function() {
												var dayData = [];

												for(var i = 1; i < 32; i++) {
													dayData.push({value: i, text: i});
												}

												return dayData;
											})()
									}
								}, {
									xtype: 'combobox',
									name: 'from_day_of_week',
									fieldLabel: _('Day of week'),
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									value: 0,
									disabled: true,
									hidden: true,
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 0, text: _('Monday')},
											{value: 1, text: _('Tuesday')},
											{value: 2, text: _('Wednesday')},
											{value: 3, text: _('Thursday')},
											{value: 4, text: _('Friday')},
											{value: 5, text: _('Satursday')},
											{value: 6, text: _('Sunday')}
										]
									}
								}, {
									xtype: 'textfield',
									name: 'from_hours',
									fieldLabel: _('Hours (local time)'),
									allowBlank: false,
									value: '00:00',
									regex: getTimeRegex()
								}, {
									xtype: 'checkbox',
									name: 'from_before',
									fieldLabel: _('The day before'),
									checked: false
								}
							]
						}, {
							xtype: 'cfieldset',
							title: _('To (scheduling date by default)'),
							layout: 'vbox',
							disabled: true,
							hidden: true,
							itemId: 'to',
							items: [
								{
									xtype: 'checkbox',
									name: 'to_enable',
									checked: false,
									fieldLabel: _('Enable')
								},	{
									xtype: 'combobox',
									name: 'to_month',
									fieldLabel: _('Month'),
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									value: 1,
									disabled: true,
									hidden: true,
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 1, text: _('January')},
											{value: 2, text: _('February')},
											{value: 3, text: _('March')},
											{value: 4, text: _('April')},
											{value: 5, text: _('May')},
											{value: 6, text: _('June')},
											{value: 7, text: _('July')},
											{value: 8, text: _('August')},
											{value: 9, text: _('September')},
											{value: 10, text: _('October')},
											{value: 11, text: _('November')},
											{value: 12, text: _('December')}
										]
									}
								}, {
									xtype: 'combobox',
									name: 'to_day',
									fieldLabel: _('Day'),
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									value: 1,
									disabled: true,
									hidden: true,
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: (function() {
												var dayData = [];

												for(var i = 1; i < 32; i++) {
													dayData.push({value: i, text: i});
												}

												return dayData;
											})()
									}
								}, {
									xtype: 'combobox',
									name: 'to_day_of_week',
									fieldLabel: _('Day of week'),
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									value: 0,
									disabled: true,
									hidden: true,
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 0, text: _('Monday')},
											{value: 1, text: _('Tuesday')},
											{value: 2, text: _('Wednesday')},
											{value: 3, text: _('Thursday')},
											{value: 4, text: _('Friday')},
											{value: 5, text: _('Satursday')},
											{value: 6, text: _('Sunday')}
										]
									}
								}, {
									xtype: 'textfield',
									name: 'to_hours',
									fieldLabel: _('Hours (local time)'),
									allowBlank: false,
									hidden: true,
									disabled: true,
									value: '00:00',
									regex: getTimeRegex()
								}, {
									xtype: 'checkbox',
									name: 'to_before',
									fieldLabel: _('The day before'),
									checked: false,
									hidden: true,
									disabled: true
								}
							]
						}
					]
				}, {
					title: _('Mailing Options'),
					items: [
						{
							xtype: 'textfield',
							fieldLabel: _('mailTo'),
							name: 'exporting_recipients'
						},{
							xtype: 'textfield',
							fieldLabel: _('Subject'),
							name: 'exporting_subject'
						}
					]
				}
			]
		}
	],

	initComponent: function() {
		//IE hack
		if(Ext.isIE) {
			this.height = 500;
			this.width = 300;
		}

		this.callParent(arguments);

		var frequencyCombo = this.down('*[name="frequency"]');
		var dayCombo = this.down('*[name="crontab_day"]');
		var dayWeekCombo = this.down('*[name="crontab_day_of_week"]');
		var monthCombo = this.down('*[name="crontab_month"]');

		var from_month = this.down('*[name="from_month"]');
		var from_day = this.down('*[name="from_day"]');
		var from_day_of_week = this.down('*[name="from_day_of_week"]');
		var from_hours = this.down('*[name="from_hours"]');

		var from_before = this.down('*[name="from_before"]');

		var to_month = this.down('*[name="to_month"]');
		var to_day = this.down('*[name="to_day"]');
		var to_day_of_week = this.down('*[name="to_day_of_week"]');
		var to_hours = this.down('*[name="to_hours"]');

		var to_before = this.down('*[name="to_before"]');

		var to_enable = this.down('*[name="to_enable"]');

		var exporting_advanced = this.down('*[name="exporting_advanced"]');

		var exporting_duration = this.down('#exporting_duration');

		var to = this.down('#to');
		var from = this.down('#from');

		var exporting_timezone = this.down("#exporting_timezone");
		var timezone_type = this.down('*[name=timezone_type]');
		var timezone_value = this.down('*[name=timezone_value]');

		var exporting_enable = this.down('*[name=exporting_enable]');

		var exporting_duration = this.down('#exporting_duration');

		exporting_enable.on('change', function(component, value) {
			switch(value) {
				case true:
					exporting_timezone.show().setDisabled(false);
					exporting_advanced.show().setDisabled(false);
					switch(exporting_advanced.getValue()) {
						case true:
							to.show().setDisabled(false);
							from.show().setDisabled(false);
							break;
						case false:
							exporting_duration.show().setDisabled(false);
							break;
					}
					break;
				case false:
					exporting_timezone.hide().setDisabled(true);
					exporting_advanced.hide().setDisabled(true);
					exporting_duration.hide().setDisabled(true);
					to.hide().setDisabled(true);
					from.hide().setDisabled(true);
					break;
			}
		});

		timezone_type.on('change', function(component, value) {
			switch(value) {
				case 'local':
				case 'server':
				case 'utc':
					timezone_value.hide().setDisabled(true);
					break;
				case 'custom':
					timezone_value.show().setDisabled(false);
					break;
				default:
					console.error('Wrong timezone type: ' + value);
			}
		});

		from_before.on('change', function(component, value) {
			switch(value) {
				case true:
					if (to_enable.getValue()) {
						to_before.show().setDisabled(false);
					}
					break;
				case false:
					to_before.hide().setDisabled(true);
					break;
			}
		});

		exporting_advanced.on('change', function(component, value) {
			switch(value) {
				case true:
					exporting_duration.hide().setDisabled(true);
					from.show().setDisabled(false);
					to.show().setDisabled(false);
					break;

				case false:
					exporting_duration.show().setDisabled(false);
					from.hide().setDisabled(true);
					to.hide().setDisabled(true);
					break;
			}
		});

		to_enable.on('change', function(component, value) {
			void(component);
			void(value);

			switch(value) {
				case true:
					if (from_before.getValue()) {
						to_before.show().setDisabled(false);
					} else {
						to_before.hide().setDisabled(true);
					}
					switch(frequencyCombo.getValue()) {
						case 'day':
							to_month.hide().setDisabled(true);
							to_day.hide().setDisabled(true);
							to_day_of_week.hide().setDisabled(true);
							to_hours.show().setDisabled(false);
							break;

						case 'week':
							to_month.hide().setDisabled(true);
							to_day.hide().setDisabled(true);
							to_day_of_week.show().setDisabled(false);
							to_hours.show().setDisabled(false);
							break;

						case 'month':
							to_month.hide().setDisabled(true);
							to_day.show().setDisabled(false);
							to_day_of_week.hide().setDisabled(true);
							to_hours.show().setDisabled(false);
							break;

						case 'year':
							to_month.show().setDisabled(false);
							to_day.show().setDisabled(false);
							to_day_of_week.hide().setDisabled(true);
							to_hours.show().setDisabled(false);
							break;

						default:
							log.debug('Wrong value');
							break;
					}
					break;

				case false:
					to_month.hide().setDisabled(true);
					to_day.hide().setDisabled(true);
					to_day_of_week.hide().setDisabled(true);
					to_hours.hide().setDisabled(true);
					to_before.hide().setDisabled(true);
					break;
			}
		});

		frequencyCombo.on('change', function(combo, value) {
			void(combo);

			from_before.setFieldLabel(_("The " + value + " before"));
			to_before.setFieldLabel(_("The " + value + " before"));

			switch (value) {
				case 'day':
					dayCombo.hide().setDisabled(true);
					dayWeekCombo.hide().setDisabled(true);
					monthCombo.hide().setDisabled(true);
					if (exporting_advanced.getValue()) {
						from_month.hide().setDisabled(true);
						from_day.hide().setDisabled(true);
						from_day_of_week.hide().setDisabled(true);
						if (to_enable.getValue()) {
							to_month.hide().setDisabled(true);
							to_day.hide().setDisabled(true);
							to_day_of_week.hide().setDisabled(true);
						}
					}
					break;

				case 'week':
					dayWeekCombo.show().setDisabled(false);
					dayCombo.hide().setDisabled(true);
					monthCombo.hide().setDisabled(true);
					if (exporting_advanced.getValue()) {
						from_month.hide().setDisabled(true);
						from_day.hide().setDisabled(true);
						from_day_of_week.show().setDisabled(false);
						if (to_enable.getValue()) {
							to_month.hide().setDisabled(true);
							to_day.hide().setDisabled(true);
							to_day_of_week.show().setDisabled(false);
						}
					}
					break;

				case 'month':
					dayWeekCombo.hide().setDisabled(true);
					dayCombo.show().setDisabled(false);
					monthCombo.hide().setDisabled(true);
					if (exporting_advanced.getValue()) {
						from_month.hide().setDisabled(true);
						from_day.show().setDisabled(false);
						from_day_of_week.hide().setDisabled(true);
						if (to_enable.getValue()) {
							to_month.hide().setDisabled(true);
							to_day.show().setDisabled(false);
							to_day_of_week.hide().setDisabled(true);
						}
					}
					break;

				case 'year':
					dayCombo.show().setDisabled(false);
					dayWeekCombo.hide().setDisabled(true);
					monthCombo.show().setDisabled(false);
					if (exporting_advanced.getValue()) {
						from_month.show().setDisabled(false);
						from_day.show().setDisabled(false);
						from_day_of_week.hide().setDisabled(true);
						if (to_enable.getValue()) {
							to_month.show().setDisabled(false);
							to_day.show().setDisabled(false);
							to_day_of_week.hide().setDisabled(true);
						}
					}
					break;

				default:
					log.debug('Wrong value');
					break;
			}
		}, this);

		this.down('*[name="exporting_owner"]').setValue(global.account.user);
	}
});
