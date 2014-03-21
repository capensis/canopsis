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
			switch(data.frequency) {
				case 'year':
					if (month !== undefined && month !== null) {
						result['month'] = month;
					} else {
						result['month'] = date.getMonth();
					}
				case 'month':
					if (day !== undefined && day !== null) {
						result['day'] = day;
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
						result['minute'] = d.getMinutes();
						result['hour'] = d.getHours();
					}
					break;
				default:
					console.log('Wrong frequency: ' + frequency);
					break;
			}
			return result;
		}

		var crontab = time2struct(data.crontab_hours, data.cron_tab_day, data.cron_tab_day_of_week, data.crontab_month);
		var from = time2struct(data.from_hours, data.from_day, data.from_day_of_week, data.from_month);
		var to = time2struct(data.to_hours, data.to_day, data.to_day_of_week, data.to_month);

		record.set('cron', crontab);

		from.before = data.from_before === "on";
		to.enable = data.to_enable === "on";

		to.before = data.to_before === "on";

		record.set('from', from);
		record.set('to', to);

		kwargs['_from'] = from;
		kwargs['_to'] = to;

		var before = {
			unit: data.frequency + 's',
			count: data.exporting_before === "on"? 1 : 0
		};

		record.set('before', before);
		kwargs['before'] = before;

		record.set('exporting_type', data.exporting_advanced === 'on' ? 'fixed' : 'duration');

		record.set('exporting_intervalLength', data.exporting_intervalLength);
		record.set('exporting_intervalUnit', data.exporting_intervalUnit);

		kwargs['exporting_type'] = record.get('exporting_type');
		kwargs['exporting_intervalUnit'] = record.get('exporting_intervalUnit');
		kwargs['exporting_intervalLength'] = record.get('exporting_intervalLength');

		record.set('kwargs',kwargs);

		return record;
	},

	beforeload_EditForm: function(form, item) {
		var crontab = item.get('cron');

		if(crontab && crontab.hour !== undefined && crontab.minute !== undefined) {
			var d = new Date();
			d.setUTCHours(crontab.hour, crontab.minute);
			var minutes = d.getMinutes();

			if(minutes < 10) {
				minutes = '0' + minutes;
			}

			form.down('textfield[name=crontab_hours]').setValue(d.getHours() + ':' + minutes);
		}

		var from = item.get('from');

		if (from !== undefined) {
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
				form.down('*[name=from_day]').setValue(from.day);
			}
			if (from.before !== undefined) {
				form.down('*[name=from_before]').setValue(from.before);
			}
			if (from.day_of_week !== undefined) {
				form.down('*[name=from_day_of_week]').setValue(from.day_of_week);
			}
		}

		var to = item.get('to');

		if (to !== undefined) {
			if (to.hour !== undefined) {
				var d = new Date();
				d.setUTCHours(to.hour, to.minute);
				var minutes = d.getMinutes();
				if (minutes < 10) {
					minutes = '0' + minutes;
				}
				hours = d.getHours();
				if (hours < 10) {
					hours = '0' + hours;
				}
				form.down('*[name=to_hours]').setValue(hours + ':' + minutes);
			}
			if (to.month !== undefined) {
				form.down('*[name=to_month]').setValue(to.month);
			}
			if (to.day !== undefined) {
				form.down('*[name=to_day]').setValue(to.day);
			}
			if (to.day_of_week !== undefined) {
				form.down('*[name=to_day_of_week]').setValue(to.day_of_week);
			}
			if (to.before !== undefined) {
				form.down('*[name=to_before]').setValue(to.before);
			}
			form.down('*[name=to_enable]').setValue(to.enable);
		}

		var before = item.get('before');

		if (before !== undefined && before !== null) {
			form.down('*[name=exporting_before]').setValue(before.count > 0);
		}

		var exporting_intervalLength = item.get('exporting_intervalLength');
		var exporting_intervalUnit = item.get('exporting_intervalUnit');

		var exporting_advanced = item.get('exporting_type') === 'fixed' ? true : false;

		form.down('*[name=exporting_intervalLength]').setValue(exporting_intervalLength);
		form.down('*[name=exporting_intervalUnit]').setValue(exporting_intervalUnit);
		form.down('*[name=exporting_advanced]').setValue(exporting_advanced);

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

		var now = new Date();

		var exporting_type = item.get('exporting_type');

		switch(exporting_type) {
			case 'fixed':
				function update_date_with_before(date) {
					if (options.before !== undefined) {
						switch(options.before.unit) {
							case 'days':
								date.setDate(date.getDate() - 1);
								break;
							case 'weeks':
								date.setDate(date.getDate() - 7);
								break;
							case 'months':
								date.setMonth(date.getMonth() - 1);
								break;
							case 'years':
								date.setYear(date.getYear() - 1);
								break;
							default:
								console.log('Wrong before unit : ' + options.before.unit);
								break;
						}
					}
				}

				function convert_struct_to_timestamp(struct) {
					var result = undefined;
					var date = new Date(now.getTime());

					if (struct.before === true) {
						update_date_with_before(date);
					}

					var frequency = item.get('frequency');

					switch(frequency) {
						case "year":
							if (struct.month !== undefined && struct.month !== null) {
								date.setMonth(struct.month);
							}
						case "month":
							if (struct.day !== undefined && struct.day !== null) {
								date.setDate(struct.day);
							}
						case "week":
							if (item.data.frequency === "week" && struct.day_of_week !== undefined && struct.day_of_week !== null) {
								date.setDate(date.getDate() + date.getDay() - struct.day_of_week);
							}
						case "day":
							if (struct.hour !== undefined && struct.hour !== null) {
								date.setHours(struct.hour);
							}
							if (struct.minute !== undefined && struct.minute !== null) {
								date.setMinutes(struct.minute);
							}
							break;
						default:
							console.logger("Wrong frequency: " + frequency);
					}

					result = date.getTime();
					return result;
				}

				if (options._from !== undefined && options._from !== null) {
					start_time = convert_struct_to_timestamp(options._from);
					stop_time = (options._to === undefined || options._to === null || ! options._to.enable) ? (now.getTime()) : convert_struct_to_timestamp(options._to);
				}
				break;

			case 'duration':
				var exporting_intervalUnit = item.data.exporting_intervalUnit;
				var exporting_intervalLength = parseInt(item.data.exporting_intervalLength, 10);
				var date = new Date(now.getTime());
				switch(exporting_intervalUnit) {
					case "hours":
						date.setHours(date.getHours() - exporting_intervalLength);
						break;
					case "days":
						date.setDate(date.getDate() - exporting_intervalLength);
						break;
					case "weeks":
						date.setDate(date.getDate() - 7 * exporting_intervalLength);
						break;
					case "months":
						date.setMonth(date.getMonth() - exporting_intervalLength);
						break;
					case "years":
						date.setYear(date.getYear() - exporting_intervalLength);
						break;
					default:
						console.error('Wrong exporting unit : ' + exporting_intervalUnit);
				}
				start_time = date.getTime();
				stop_time = now.getTime();
				break;
		}

		var mail = options.mail;

		if(mail) {
			mail = Ext.encode(mail);
		}

		this.getController('Reporting').launchReport(view_name, start_time, stop_time, mail, undefined, undefined);
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
