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

		function time2struct(hours, day, day_of_week, month, before){
			var result = {};
			var date = new Date();

			if (before === 'on') {
				var frequency_unit = data.frequency + 's';
				result['before'] = {};
				result['before'][frequency_unit] = 1;
			}

			switch(data.frequency) {
				case 'year':
					if (month !== undefined && month !== null) {
						result['month'] = month;
					} else {
						result['month'] = date.getUTCMonth();
					}
				case 'month':
					if (day !== undefined && day !== null) {
						result['day'] = day;
					} else {
						result['day'] = date.getUTCDate();
					}
				case 'week':
					if (data.frequency === 'week'){
						if(day_of_week !== undefined && day_of_week !== null) {
							result['day_of_week'] = day_of_week;
						} else {
							result['day_of_week'] = date.getUTCDay();
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
			return result;
		}

		var crontab = time2struct(data.crontab_hours, data.crontab_day, data.crontab_day_of_week, data.crontab_month);
		var from = time2struct(data.from_hours, data.from_day, data.from_day_of_week, data.from_month, data.from_before);
		var to = time2struct(data.to_hours, data.to_day, data.to_day_of_week, data.to_month, data.to_before);

		record.set('cron', crontab);

		to.enable = data.to_enable === "on";

		var timezone_type = data.timezone_type;
		if (timezone_type === undefined || timezone_type === null) {
			timezone_type = "local";
		}
		var timezone_value = data.timezone_value;
		if (timezone_value === null || timezone_value === undefined) {
			timezone_value = new Date().getTimezoneOffset() * 60;
		}

		var timezone = {
			type: timezone_type,
			value: timezone_value
		}

		var exporting = {
			enable: data.exporting_enable === 'on',
			length: data.exporting_intervalLength,
			unit: data.exporting_intervalUnit,
			type: data.exporting_advanced === 'on' ? "fixed" : "duration",
			from: from,
			to: to,
			timezone: timezone
		};

		record.set("exporting", exporting);

		kwargs['exporting'] = exporting;

		record.set('kwargs', kwargs);

		return record;
	},

	beforeload_EditForm: function(form, item) {
		var crontab = item.get('cron');

		var exporting = item.get('exporting');

		function update_forms_from_timestruct(timestruct, name) {
			if (timestruct !== undefined) {
				if (timestruct.hour !== undefined) {
					var d = new Date();
					d.setUTCHours(timestruct.hour, timestruct.minute);
					var minutes = d.getMinutes();
					if (minutes < 10) {
						minutes = '0' + minutes;
					}
					form.down('*[name='+name+'_hours]').setValue(d.getHours() + ':' + minutes);
				}
				if (timestruct.month !== undefined) {
					form.down('*[name='+name+'_month]').setValue(timestruct.month);
				}
				if (timestruct.day !== undefined) {
					form.down('*[name='+name+'_day]').setValue(timestruct.day);
				}
				if (timestruct.day_of_week !== undefined) {
					form.down('*[name='+name+'_day_of_week]').setValue(timestruct.day_of_week);
				}
				if (timestruct.before !== undefined) {
					for (before_type in timestruct.before) {
						if (timestruct.before[before_type] === 1) {
							form.down('*[name='+name+'_before]').setValue(true);
							break;
						}
					}
				}
			}
		}

		var from = exporting.from;
		var to = exporting.to;

		update_forms_from_timestruct(crontab, 'crontab');
		update_forms_from_timestruct(from, 'from');
		update_forms_from_timestruct(to, 'to');

		if (to !== null && to !== undefined && to.enable) {
			form.down('*[name=to_enable]').setValue(to.enable);
		}

		var exporting_intervalLength = exporting.length;
		var exporting_intervalUnit = exporting.unit;
		var exporting_advanced = exporting.type === 'fixed';

		var exporting_enable = exporting.enable === true;

		form.down('*[name=exporting_enable]').setValue(exporting_enable);

		if (exporting_intervalLength !== undefined && exporting_intervalLength !== null) {
			form.down('*[name=exporting_intervalLength]').setValue(exporting_intervalLength);
		}
		if (exporting_intervalUnit !== undefined && exporting_intervalUnit !== null) {
			form.down('*[name=exporting_intervalUnit]').setValue(exporting_intervalUnit);
		}
		form.down('*[name=exporting_advanced]').setValue(exporting_advanced);

		var timezone = exporting.timezone;
		if (timezone !== null && timezone !== undefined) {
			var timezone_type = timezone.type;
			var timezone_value = timezone.value;

			form.down('*[name=timezone_type]').setValue(timezone_type);
			if (timezone_value !== undefined && timezone_value !== null) {
				form.down('*[name=timezone_value]').setValue(timezone_value);
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

		var now = new Date();

		var exporting = item.get('exporting');

		var exporting_type = exporting.type;

		if (exporting.enable) {

			switch(exporting_type) {
				case 'fixed':
					function update_date_with_before(date, before) {
						if (before !== undefined) {
							for (unit in before) {
								switch(unit) {
									case 'days':
										date.setDate(date.getUTCDate() - before[unit]);
										break;
									case 'weeks':
										date.setDate(date.getUTCDate() - before[unit] * 7);
										break;
									case 'months':
										date.setMonth(date.getUTCMonth() - before[unit]);
										break;
									case 'years':
										date.setYear(date.getYear() - before[unit]);
										break;
									default:
										console.log('Wrong before unit : ' + unit);
										break;
								}
							}
						}
					}

					function convert_struct_to_timestamp(struct) {
						var result = undefined;
						var date = new Date(now.getTime());

						update_date_with_before(date, struct.before);

						var frequency = item.get('frequency');

						switch(frequency) {
							case "year":
								if (struct.month !== undefined && struct.month !== null) {
									date.setUTCMonth(struct.month);
								}
							case "month":
								if (struct.day !== undefined && struct.day !== null) {
									date.setUTCDate(struct.day);
								}
							case "week":
								if (item.data.frequency === "week" && struct.day_of_week !== undefined && struct.day_of_week !== null) {
									date.setUTCDate(date.getUTCDate() + date.getUTCDay() - struct.day_of_week);
								}
							case "day":
								if (struct.hour !== undefined && struct.hour !== null) {
									date.setUTCHours(struct.hour);
								}
								if (struct.minute !== undefined && struct.minute !== null) {
									date.setUTCMinutes(struct.minute);
								}
								break;
							default:
								console.logger("Wrong frequency: " + frequency);
						}

						result = date.getTime();
						return result;
					}

					var from = exporting.from;

					if (from !== undefined && from !== null) {
						start_time = convert_struct_to_timestamp(from);
						var to = exporting.to;
						stop_time = (to === undefined || to === null || ! to.enable) ?
							(now.getTime())
							:
							convert_struct_to_timestamp(to);
					}
					break;

				case 'duration':
					var exporting_intervalUnit = exporting.unit;
					var exporting_intervalLength = parseInt(exporting.length, 10);
					var date = new Date(now.getTime());
					switch(exporting_intervalUnit) {
						case "hours":
							date.setHours(date.getUTCHours() - exporting_intervalLength);
							break;
						case "days":
							date.setDate(date.getUTCDate() - exporting_intervalLength);
							break;
						case "weeks":
							date.setDate(date.getUTCDate() - 7 * exporting_intervalLength);
							break;
						case "months":
							date.setMonth(date.getUTCMonth() - exporting_intervalLength);
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
		}

		var mail = options.mail;

		if(mail) {
			mail = Ext.encode(mail);
		}

		var timezone = exporting['timezone'];
		if (timezone !== undefined) {
			timezone = timezone['value'];
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

		var minute = d.getUTCMinutes();
		var hour = d.getUTCHours();

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
