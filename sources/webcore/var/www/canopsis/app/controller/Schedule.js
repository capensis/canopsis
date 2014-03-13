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
		else {
			kwargs['mail'] = null;
		}

		record.set('loaded', false);

		function time2struct(hours, day, day_of_week, month, type, intervalLength, intervalUnit){
			var result = {};
			var date = new Date();
			switch(type === undefined? 'Date and Time' : type) {
				case 'Date and Time':
					switch(data.frequency) {
						case 'year':
							if (month !== undefined && month !== null) {
								result['month'] = month;
							} else {
								result['month'] = date.getMonth();
							}
						case 'month':
							if (day !== undefined && day !== null) {
								result['day'] = month;
							} else {
								result['day'] = date.getDate();
							}
						case 'week':
							if (data.frequency === 'week'){
								if(day_of_week !== undefined && day_of_week !== null) {
									result['day_of_week'] = day_of_week;
								} else {
									result['day_of_week'] = date.getDay();
								}
							}
						case 'day':
							if (hours !== undefined && hours !== null) {
								var time = stringTo24h(hours);
								var d = new Date();
								d.setHours(time.hour);
								d.setMinutes(time.minute, 10);
								result['minute'] = d.getUTCMinutes();
								result['hour'] = d.getUTCHours();
							}
							break;
						default:
							console.log('Wrong frequency: ' + frequency);
							break;
					}
					break;
				case 'Duration':
					if (intervalLength !== undefined) {
						result['intervalLength'] = intervalLength;
					}
					if (intervalUnit !== undefined) {
						result['intervalUnit'] = intervalUnit;
					}
					break;
				default:
					console.log('Wrong type: ' + type);
					break;
			}
			if (type !== undefined) {
				result['type'] = type;
			}
			return result;
		}

		var crontab = time2struct(data.crontab_hours, data.cron_tab_day, data.cron_tab_day_of_week, data.crontab_month);
		var from = time2struct(data.from_hours, data.from_day, data.from_day_of_week, data.from_month, data.from_type, data.from_intervalLength, data.from_intervalUnit);
		var to = time2struct(data.to_hours, data.to_day, data.to_day_of_week, data.to_month, data.to_type, data.to_intervalLength, data.to_intervalUnit);

		record.set('cron', crontab);
		record.set('from', from);
		record.set('to', to);

		kwargs['_to'] = to;
		kwargs['_from'] = from;

		kwargs['timezone'] = data.timezone * -3600;

		record.set('kwargs',kwargs);

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

		timezone = item.data.timezone;

		if (timezone === undefined) {
			timezone = new Date().getTimezoneOffset() / -60;
		} else {
			timezone /= 3600;
		}

		form.down('*[name=timezone]').setValue(timezone);

		from = item.data.from;

		if (from !== undefined) {
			form.down('*[name=from_type]').setValue(from.type);
			if (from.intervalLength !== undefined) {
				form.down('*[name=from_intervalLength]').setValue(from.intervalLength);
			}
			if (from.intervalUnit !== undefined) {
				form.down('*[name=from_intervalUnit]').setValue(from.intervalUnit);
			}
			if (from.hour !== undefined) {
				var d = new Date();
				d.setUTCHours(from.hour, from.minute);
				var minutes = d.getMinutes();
				if (minutes < 10) {
					minutes = '0' + minutes;
				}
				form.down('*[name=from_hours]').setValue(d.getHours() + ':' + minutes);
			}
			if (from.month !== undefined) {
				form.down('*[name=from_month]').setValue(from.month);
			}
			if (from.day !== undefined) {
				form.down('*[name=from_month]').setValue(from.day);
			}
			if (from.day_of_week !== undefined) {
				form.down('*[name=from_day_of_week]').setValue(from.day_of_week);
			}
		}

		to = item.data.to;

		if (to !== undefined) {
			form.down('*[name=to_type]').setValue(to.type);
			if (to.intervalLength !== undefined) {
				form.down('*[name=to_intervalLength]').setValue(to.intervalLength);
			}
			if (to.intervalUnit !== undefined) {
				form.down('*[name=to_intervalUnit]').setValue(to.intervalUnit);
			}
			if (to.hour !== undefined) {
				var d = new Date();
				d.setUTCHours(to.hour, to.minute);
				var minutes = d.getMinutes();
				if (minutes < 10) {
					minutes = '0' + minutes;
				}
				form.down('*[name=to_hours]').setValue(d.getHours() + ':' + minutes);
			}
			if (to.month !== undefined) {
				form.down('*[name=to_month]').setValue(to.month);
			}
			if (to.day !== undefined) {
				form.down('*[name=to_month]').setValue(to.day);
			}
			if (to.day_of_week !== undefined) {
				form.down('*[name=to_day_of_week]').setValue(to.day_of_week);
			}
		}
	},

	validateForm: function(store, data, form) {
		var field = undefined;

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
		var stop_time = undefined;

		var now = new Date().getTime();

		function convert_struct_to_timestamp(struct) {
			var result = undefined;
			var date = undefined;
			switch(struct.type) {
				case 'Duration':
					var intervalUnit = parseFloat(struct.intervalUnit);
					date = start_time === undefined ? new Date() : new Date(start_time);
					switch(struct.intervalLength) {
						case 'hours':
							if (start_time === undefined) {
								date.setHours(date.getHours() - intervalUnit);
							} else {
								date.setHours(date.getHours() + intervalUnit);
							}
							break;
						case 'days':
							if (start_time === undefined) {
								date.setDate(date.getDate() - intervalUnit);
							} else {
								date.setDate(date.getDate() + intervalUnit);
							}
							break;
						case 'weeks':
							if (start_time === undefined) {
								date.setDate(date.getDate() - 7 * intervalUnit);
							} else {
								date.setDate(date.getDate() + 7 * intervalUnit);
							}
							break;
						case 'months':
							if (start_time === undefined) {
								date.setMonth(date.getMonth() - intervalUnit);
							} else {
								date.setMonth(date.getMonth() + intervalUnit);
							}
							break;
						case 'years':
							if (start_time === undefined) {
								date.setYear(date.getYear() - intervalUnit);
							} else {
								date.setYear(date.getYear() + intervalUnit);
							}
							break;
						default:
							console.log("Wrong interval length: " + struct.intervalLength);
					}
					break;
				case 'Date and Time':
					date = new Date();
					switch(item.frequency) {
						case 'year':
							if (struct.month !== undefined && struct.month !== null) {
								date.setMonth(struct.month);
							}
						case 'month':
							if (struct.day !== undefined && struct.day !== null) {
								date.setDate(struct.day);
							}
						case 'week':
							if (item.frequency === 'week' && struct.day_of_week !== undefined && struct.day_of_week !== null) {
								day.setDate(date.getDate() + date.getDay() - struct.day_of_week);
							}
						case 'day':
							if (struct.hour !== undefined && struct.hour !== null) {
								day.setHours(struct.hour);
							}
							if (struct.minute !== undefined && struct.minute !== null) {
								day.setMinutes(struct.minute);
							}
							break;
						default:
						console.log('Error of frequency: ' + item.frequency);
						break;
					}
					break;
				default:
					console.log('Error of type: ' + struct.type);
					break;
			}
			result = date.getTime();
			return result;
		}

		if (options._from !== undefined && options._from !== null) {
			start_time = convert_struct_to_timestamp(options._from);
			stop_time = options._to === undefined || options._to === null ? (now) : convert_struct_to_timestamp(options._to);
		}

		var mail = options.mail;

		if(mail) {
			mail = Ext.encode(mail);
		}

		this.getController('Reporting').launchReport(view_name, start_time, stop_time, mail, undefined, undefined, options.timezone);
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
