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
			xtype: 'cfieldset',
			title: _('General options'),
			layout: 'anchor',
			defaults:{flex:1},
			items: [
				{
					xtype: 'textfield',
					fieldLabel: _('Schedule name'),
					name: 'crecord_name',
					allowBlank: false
				},{
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
				},{
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
				},{
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
		},{
			xtype: 'cfieldset',
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
					value: _('January'),
					disabled: true,
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
				},{
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
				},{
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
				},{
					xtype: 'textfield',
					name: 'crontab_hours',
					fieldLabel: _('Hours (local time)'),
					allowBlank: false,
					regex: getTimeRegex()
				}
			]
		},{
			xtype: 'cfieldset',
			title: _('Exporting Interval'),
			layout: 'hbox',
			checkboxName: 'exporting_interval',
			items: [
				{
					xtype: 'combobox',
					name: 'exporting_intervalLength',
					queryMode: 'local',
					displayField: 'text',
					width: 90,
					valueField: 'value',
					padding: '0 0 5 5',
					value: global.commonTs.day,
					store: {
						xtype: 'store',
						fields: ['value', 'text'],
						data: [
							{value: global.commonTs.day, text: _('Day')},
							{value: global.commonTs.week, text: _('Week')},
							{value: global.commonTs.month, text: _('Month')},
							{value: global.commonTs.year, text: _('Year')}
						]
					}
				},{
					xtype: 'numberfield',
					name: 'exporting_intervalUnit',
					fieldLabel: _('The last'),
					minValue: 1,
					value: 1,
					width: 160
				}
			]
		},{
			xtype: 'cfieldset',
			checkboxName: 'exporting_mail',
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
	],

	initComponent: function() {
		//IE hack
		if(Ext.isIE) {
			this.height = 500;
			this.width = 300;
		}

		this.callParent(arguments);

		var durationCombo = this.down('*[name="frequency"]');
		dayCombo = this.down('*[name="crontab_day"]');
		dayWeekCombo = this.down('*[name="crontab_day_of_week"]');
		monthCombo = this.down('*[name="crontab_month"]');

		durationCombo.on('change', function(combo, value) {
			void(combo);

			switch (value) {
				case 'day':
					dayCombo.hide().setDisabled(true);
					dayWeekCombo.hide().setDisabled(true);
					monthCombo.setDisabled(true);
					break;

				case 'week':
					dayCombo.hide().setDisabled(true);
					dayWeekCombo.show().setDisabled(false);
					monthCombo.setDisabled(true);
					break;

				case 'month':
					dayCombo.show().setDisabled(false);
					dayWeekCombo.hide().setDisabled(true);
					monthCombo.setDisabled(true);
					break;

				case 'year':
					dayCombo.show().setDisabled(false);
					dayWeekCombo.hide().setDisabled(true);
					monthCombo.setDisabled(false);
					break;

				default:
					log.debug('Wrong value');
					break;
			}
		}, this);

		this.down('*[name="exporting_owner"]').setValue(global.account.user);
	}
});
