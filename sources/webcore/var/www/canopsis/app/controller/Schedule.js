//need:app/lib/controller/cgrid.js,app/view/Schedule/Grid.js,app/view/Schedule/Form.js,app/store/Schedules.js
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
Ext.define('canopsis.controller.Schedule', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Schedule.Grid', 'Schedule.Form'],
	stores: ['Schedules'],
	models: ['Schedule'],

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.formXtype = 'ScheduleForm';
		this.listXtype = 'ScheduleGrid';

		this.modelId = 'Schedule';

		this.callParent(arguments);
	},

	preSave: function(record, data) {
		var interval = null;

		if(data.exporting_interval) {
			interval = data.exporting_intervalLength * data.exporting_intervalUnit;
		}

		record.set('exporting_interval', data.exporting_interval);
		record.set('exporting_account', global.account.user);
		record.set('exporting_task', 'task_reporting');
		record.set('exporting_method', 'render_pdf');

		var kwargs = {
			viewName: record.get('exporting_viewName'),
			account: record.get('exporting_account'),
			task: record.get('exporting_task'),
			method: record.get('exporting_method'),
			interval: interval,
			_scheduled: record.get('crecord_name'),
			owner: record.get('exporting_owner')
		};

		//check if a mail must be send
		if(data.exporting_mail) {
			if(data.exporting_recipients !== '' && data.exporting_recipients !== undefined) {
				log.debug('sendMail is true', this.logAuthor);

				var stripped_recipients = data.exporting_recipients.replace(/ /g, '');
				var recipients = stripped_recipients.split(',');

				if(recipients.length === 1) {
					recipients = stripped_recipients.split(';');
				}

				var mail = {
					'recipients': recipients,
					'subject': record.get('exporting_subject'),
					'body': 'Scheduled task reporting'
				};

				kwargs['mail'] = mail;
			}
		}
		else {
			kwargs['mail'] = null;
		}

		record.set('kwargs',kwargs);
		record.set('loaded', false);

		// formating crontab
		var time = stringTo24h(data.crontab_hours);

		//apply offset to get utc
		var d = new Date();
		d.setHours(time.hour);
		d.setMinutes(time.minute, 10);

		//set crontab
		var crontab = {
			minute: d.getUTCMinutes(),
			hour: d.getUTCHours()
		};

		if(data.crontab_month) {
			crontab['month'] = data.crontab_month;
		}

		if(data.crontab_day_of_week) {
			crontab['day_of_week'] = data.crontab_day_of_week;
		}

		if(data.crontab_day) {
			crontab['day'] = data.crontab_day;
		}

		record.set('cron', crontab);

		return record;
	},

	beforeload_EditForm: function(form, item) {
		crontab = item.data.cron;

		if(crontab && crontab.hour !== undefined && crontab.minute !== undefined) {
			var d = new Date();
			d.setUTCHours(crontab.hour, crontab.minute);
			var minutes = d.getMinutes();

			if(minutes < 10) {
				minutes = '0' + minutes;
			}

			form.down('textfield[name=crontab_hours]').setValue(d.getHours() + ':' + minutes);
		}
	},

	validateForm: function(store, data, form) {
		var field = undefined;

		//check mail options
		if(data['exporting_mail'] && !data['exporting_subject'] || !data['exporting_subject']) {
			log.debug('Invalid mail options', this.logAuthor + '[validateForm]');
			global.notify.notify(' Invalid mail options', '', 'error');

			field = form.findField('exporting_subject');

			if(!data['exporting_subject'] && field) {
				field.markInvalid(_('Invalid field'));
			}

			field = form.findField('exporting_recipients');

			if(!data['exporting_recipients'] && field) {
				field.markInvalid(_('Invalid field'));
			}

			return false;
		}

		//Check duplicate
		var already_exist = false;

		if(!form.editing && store.findExact('crecord_name', data['crecord_name']) >= 0) {
			already_exist = true;
		}

		field = form.findField('crecord_name');

		if(field) {
			field.markInvalid(_('Invalid field'));
		}

		if (already_exist) {
			log.debug('Schedule already exist exist', this.logAuthor + '[validateForm]');
			global.notify.notify(data['crecord_name'] + ' already exist', 'you can\'t add the same Schedule twice', 'error');
			return false;
		}

		return true;
	},

	runItem: function(item) {
		log.debug('Clicked on run item', this.logAuthor);

		var options    = item.get('kwargs');
		var view_name  = options.viewName;
		var start_time = undefined;

		if(options.interval) {
			start_time = Ext.Date.now() - (options.interval * 1000);
		}

		var mail = options.mail;

		if(mail) {
			mail = Ext.encode(mail);
		}

		this.getController('Reporting').launchReport(view_name, start_time, undefined, mail);
	},

	//call a window wizard to schedule Schedule with passed argument
	scheduleWizard: function(item, renderTo) {
		//temporary hack, check if called by cgrid or ctree
		var form = Ext.create('canopsis.view.Schedule.Form', {
			EditMethod: 'window',
			editing: false
		});

		var store = Ext.getStore('Schedules');

		if(item !== undefined) {
			var viewName = item.get('_id');
			var combo = form.down('combobox[name=view]');

			if(combo !== null) {
				combo.setValue(viewName);
			}
		}

		var window_wizard = Ext.widget('window', {
			title: _('Scheduling'),
			items: [form],
			constrain: true,
			renderTo: renderTo
		});

		form.win = window_wizard;

		window_wizard.show();

		// binding events
		var btns = form.down('button[action=save]');

		btns.on('click', function() {
			this._saveForm(form, store);
		}, this);

		btns = form.down('button[action=cancel]');

		btns.on('click', function() {
			window_wizard.close();
		}, this);
	},

	format_time: function(cron) {
		//format time
		var d = new Date();
		d.setUTCHours(parseInt(cron.hour, 10));
		d.setUTCMinutes(parseInt(cron.minute, 10));

		var minute = d.getMinutes();
		var hour = d.getHours();

		//cosmetic
		if(minute < 10) {
			minute = '0' + minute;
		}

		if(hour < 10) {
			hour = '0' + hour;
		}

		//check 12h / 24h clock
		var hours = undefined;

		if(!is12Clock()) {
			hours = hour + ':' + minute;
		}
		else {
			if(hour > 12) {
				hours = (hour - 12) + ':' + minute + ' pm';
			}
			else {
				hours = hour + ':' + minute + ' am';
			}
		}

		return hours;
	}
});
