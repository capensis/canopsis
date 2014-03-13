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
							xtype: "numberfield",
							value: - (new Date().getTimezoneOffset() / 60),
							fieldLabel: _('Time zone (hr)'),
							name: 'timezone'
						}, {
							xtype: 'cfieldset',
							title: _('From'),
							layout: 'vbox',
							items: [
								{
									xtype: 'combobox',
									name: 'from_type',
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									padding: '0 0 5 5',
									value: 'Duration',
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 'Duration', text: _('Duration')},
											{value: 'Date and Time', text: _('Date and Time')}
										]
									},
								},	{
									xtype: 'numberfield',
									name: 'from_intervalUnit',
									allowBlank: false,
									minValue: 1,
									value: 1
								},	{
									xtype: 'combobox',
									name: 'from_intervalLength',
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									padding: '0 0 5 5',
									value: 'days',
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 'hours', text: _('Hour')+'(s)'},
											{value: 'days', text: _('Day')+'(s)'},
											{value: 'weeks', text: _('Week')+'(s)'},
											{value: 'months', text: _('Month')+'(s)'},
											{value: 'years', text: _('Year')+'(s)'}
										]
									}
								}, {
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
									disabled: true,
									hidden: true,
									allowBlank: false,
									value: '00:00',
									regex: getTimeRegex()
								}
							]
						}, {
							xtype: 'cfieldset',
							title: _('To'),
							layout: 'vbox',
							items: [
								{
									xtype: 'combobox',
									name: 'to_type',
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									padding: '0 0 5 5',
									value: 'Duration',
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 'Duration', text: _('Duration')},
											{value: 'Date and Time', text: _('Date and Time')}
										]
									},
								}, {
									xtype: 'numberfield',
									name: 'to_intervalUnit',
									allowBlank: false,
									minValue: 1,
									value: 1
								}, {
									xtype: 'combobox',
									name: 'to_intervalLength',
									queryMode: 'local',
									displayField: 'text',
									valueField: 'value',
									padding: '0 0 5 5',
									value: 'days',
									store: {
										xtype: 'store',
										fields: ['value', 'text'],
										data: [
											{value: 'hours', text: _('Hour')+'(s)'},
											{value: 'days', text: _('Day')+'(s)'},
											{value: 'weeks', text: _('Week')+'(s)'},
											{value: 'months', text: _('Month')+'(s)'},
											{value: 'years', text: _('Year')+'(s)'}
										]
									}
								}, {
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
									disabled: true,
									hidden: true,
									value: '00:00',
									regex: getTimeRegex()
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

		var from_type = this.down('*[name="from_type"]');
		var from_intervalUnit = this.down('*[name="from_intervalUnit"]');
		var from_intervalLength = this.down('*[name="from_intervalLength"]');
		var from_month = this.down('*[name="from_month"]');
		var from_day = this.down('*[name="from_day"]');
		var from_day_of_week = this.down('*[name="from_day_of_week"]');
		var from_hours = this.down('*[name="from_hours"]');

		var to_type = this.down('*[name="to_type"]');
		var to_intervalUnit = this.down('*[name="to_intervalUnit"]');
		var to_intervalLength = this.down('*[name="to_intervalLength"]');
		var to_month = this.down('*[name="to_month"]');
		var to_day = this.down('*[name="to_day"]');
		var to_day_of_week = this.down('*[name="to_day_of_week"]');
		var to_hours = this.down('*[name="to_hours"]');

		from_type.on('change', function(elt, e, eOpts) {
			void(elt);
			void(eOpts);

			switch(e) {
				case 'Date and Time':
					from_intervalLength.hide().disable();
					from_intervalUnit.hide().disable();
					frequencyComboValue = frequencyCombo.getValue();
					switch (frequencyComboValue) {
						case 'day':
							from_day.hide().setDisabled(true);
							from_day_of_week.hide().setDisabled(true);
							from_month.hide().setDisabled(true);
							break;

						case 'week':
							from_day_of_week.show().setDisabled(false);
							from_day.hide().setDisabled(true);
							from_month.hide().setDisabled(true);
							break;

						case 'month':
							from_day_of_week.hide().setDisabled(true);
							from_day.show().setDisabled(false);
							from_month.hide().setDisabled(true);
							break;

						case 'year':
							from_day.show().setDisabled(false);
							from_day_of_week.hide().setDisabled(true);
							from_month.show().setDisabled(false);
							break;

						default:
							log.debug('Wrong value');
							break;
					}
					from_hours.show().enable();
					break;

				case 'Duration':
					from_intervalLength.show().enable();
					from_intervalUnit.show().enable();
					from_month.hide().disable();
					from_day.hide().disable();
					from_day_of_week.hide().disable();
					from_hours.hide().disable();
					break;
			}
		}, this);

		to_type.on('change', function(elt, e, eOpts) {
			void(elt);
			void(eOpts);

			switch(e) {
				case 'Date and Time':
					to_intervalLength.hide().disable();
					to_intervalUnit.hide().disable();
					frequencyComboValue = frequencyCombo.getValue();
					switch (frequencyComboValue) {
						case 'day':
							to_day.hide().setDisabled(true);
							to_day_of_week.hide().setDisabled(true);
							to_month.hide().setDisabled(true);
							break;

						case 'week':
							to_day_of_week.show().setDisabled(false);
							to_day.hide().setDisabled(true);
							to_month.hide().setDisabled(true);
							break;

						case 'month':
							to_day_of_week.hide().setDisabled(true);
							to_day.show().setDisabled(false);
							to_month.hide().setDisabled(true);
							break;

						case 'year':
							to_day.show().setDisabled(false);
							to_day_of_week.hide().setDisabled(true);
							to_month.show().setDisabled(false);
							break;

						default:
							log.debug('Wrong value');
							break;
					}
					to_hours.show().enable();
					break;

				case 'Duration':
					to_intervalLength.show().enable();
					to_intervalUnit.show().enable();
					to_month.hide().disable();
					to_day.hide().disable();
					to_day_of_week.hide().disable();
					to_hours.hide().disable();
					break;
			}
		}, this);

		function renewDateAndTimeDisplay() {
			if (from_type.getValue()==='Day and Time') {
				from_type.setValue('');
				from_type.setValue('Day and Time');
			}
			if (to_type.getValue()==='Day and Time') {
				to_type.setValue('');
				to_type.setValue('Day and Time');
			}
		}

		frequencyCombo.on('change', function(combo, value) {
			void(combo);

			switch (value) {
				case 'day':
					dayCombo.hide().setDisabled(true);
					dayWeekCombo.hide().setDisabled(true);
					monthCombo.hide().setDisabled(true);
					break;

				case 'week':
					dayWeekCombo.show().setDisabled(false);
					dayCombo.hide().setDisabled(true);
					monthCombo.hide().setDisabled(true);
					break;

				case 'month':
					dayWeekCombo.hide().setDisabled(true);
					dayCombo.show().setDisabled(false);
					monthCombo.hide().setDisabled(true);
					break;

				case 'year':
					dayCombo.show().setDisabled(false);
					dayWeekCombo.hide().setDisabled(true);
					monthCombo.show().setDisabled(false);
					break;

				default:
					log.debug('Wrong value');
					break;
			}
			renewDateAndTimeDisplay();
		}, this);

		this.down('*[name="exporting_owner"]').setValue(global.account.user);
	}
});
