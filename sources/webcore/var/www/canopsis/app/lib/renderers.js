/*
#--------------------------------
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
# ---------------------------------
*/
var rdr_tstodate = function(val) {
	if (val) {
		var dval = new Date(parseInt(val) * 1000);

		var dval_day = Ext.Date.format(dval, 'Y-m-d');
		var now_day = Ext.Date.format(new Date(), 'Y-m-d');

		if (dval_day == now_day) {
			if (is12Clock())
				return Ext.Date.format(dval, 'h:i:s a');
			else
				return Ext.Date.format(dval, 'H:i:s');
		}else {
			if (is12Clock())
				return Ext.Date.format(dval, 'Y-m-d h:i:s a');
			else
				return Ext.Date.format(dval, 'Y-m-d H:i:s');
		}
	}else {
		return '';
	}
};

var rdr_utcToLocal = function(val) {
	if (val != undefined && val != '') {
		//format date
		var array_split = val.split(' ');
		var date = array_split[0].split('-');
		var hour = array_split[1].split(':');

		//create date
		var dval = new Date(date[0], date[1] - 1, date[2], hour[0], hour[1], hour[2]);

		return rdr_tstodate(get_timestamp_utc(dval));
	}
};

var rdr_boolean = function(val, metadata, record, rowIndex, colIndex, store) {
	if (val)
		return "<span class='icon icon-true' />";
	else
		return "<span class='icon icon-false' />";
};

var rdr_status = function(val, metadata, record, rowIndex, colIndex, store) {
	if (typeof(val) == 'number')
		return "<span class='icon icon-state-" + val + "' />";
	return val;
};

var rdr_color = function(val, metadata, record, rowIndex, colIndex, store) {
	return "<span class='icon' style='background-color: #" + val + ';color: #' + val + ";'/>";
};

var rdr_state_type = function(val, metadata, record, rowIndex, colIndex, store) {
	return "<span class='icon icon-state-type-" + val + "' />";
};

var rdr_source_type = function(val, metadata, record, rowIndex, colIndex, store) {
	return "<span class='icon icon-crecord_type-" + val + "' />";
};

var rdr_crecord_type = function(val, metadata, record, rowIndex, colIndex, store) {
	if (val != '') {
		return "<span class='icon icon-crecord_type-" + val + "' />";
	}
};

var rdr_file_type = function(val, metadata, record, rowIndex, colIndex, store) {
	var split = val.split('/');
	if (split.length > 0) {
		return "<span class='icon icon-mimetype-" + split[split.length - 1] + "' />";
	}
};

var rdr_havePerfdata = function(val, metadata, record, rowIndex, colIndex, store) {
	if (val != '') {
		return "<span class='icon icon-perfdata'/>";
	}
};

var rdr_widget_preview = function(val, metadata, record, rowIndex, colIndex, store) {
	return "<span style='background-color:" + global.default_colors[rowIndex] + ';color:' + global.default_colors[rowIndex] + ";'>__</span>";
};

var rdr_task_crontab = function(val, metadata, record, rowIndex, colIndex, store) {

	var output = '';

	if (val != undefined) {
		//second condition is if minutes are str and not int
		if (val.hour != undefined && val.minute != undefined) {
			var d = new Date();
			d.setUTCHours(parseInt(val.hour, 10));
			d.setUTCMinutes(parseInt(val.minute, 10));

			var utc_minutes = d.getUTCMinutes();
			var utc_hours = d.getUTCHours();
			var local_minutes = d.getMinutes();
			var local_hours = d.getHours();

			//cosmetic
			if (utc_minutes < 10)
				utc_minutes = '0' + utc_minutes;
			if (local_minutes < 10)
				local_minutes = '0' + local_minutes;

			//12h translate
			if (!is12Clock()) {
				output += local_hours + ':' + local_minutes;
			} else {
				//utc AM/PM check
				if (local_hours > 12)
					output += (local_hours - 12) + ':' + local_minutes + ' pm';
				else
					output += local_hours + ':' + local_minutes + ' am';
			}

		}

		if (val.month != undefined && val.day != undefined)
			output += '   |    ' + _('month') + ' : ' + global.numberToMonth[val.month] + ' |  day : ' + val.day;

		if (val.day_of_week != undefined)
			output += '   |   ' + _('day') + ' : ' + _(val.day_of_week);
	}

	return output;
};

//Function for rendering export to pdf button, we haven't find another solution
var rdr_export_button = function(val, metadata, record, rowIndex, colIndex, store,view) {
	var id = Ext.id();
	if (record.get('leaf')) {
		var output = '';
		output += Ext.String.format(
			'<div style="{0}" class="{1}" onclick="Ext.getCmp(\'{2}\').ownerCt.export_pdf(\'{3}\')"></div>',
			'height:16px;width:16px;',
			'icon-mimetype-pdf',
			view.id,
			record.get('_id')
		);
		return output;
	}
};

var rdr_mail_information = function(val, metadata, record, rowIndex, colIndex, store,view) {
	if (val == false)
		return _('This task is not send by mail');

	output = '';
	if (val.recipients != undefined)
		output += _('Recipients :') + ' ' + val.recipients;

	return output;
};

var rdr_clean_id = function(val) {
	if (val.search('.') != -1) {
		tmp = val.split('.');
		val = tmp[1];
	}

	return val;
};

var rdr_time_interval = function(val) {
	if (val == 0 || val == null || val == undefined)
		return '';

	var tmp;

	tmp = Math.round((val / global.commonTs.year) * 100) / 100;
	if (tmp >= 1)
		return tmp + ' ' + _('Year') + '(s)';

	tmp = Math.round((val / global.commonTs.month) * 100) / 100;
	if (tmp >= 1)
		return tmp + ' ' + _('Month') + '(s)';

	tmp = Math.round((val / global.commonTs.week) * 100) / 100;
	if (tmp >= 1)
		return tmp + ' ' + _('Week') + '(s)';

	tmp = Math.round((val / global.commonTs.day) * 100) / 100;
	if (tmp >= 1)
		return tmp + ' ' + _('Day') + '(s)';

	tmp = Math.round((val / global.commonTs.hours) * 100) / 100;
	if (tmp >= 1)
		return tmp + ' ' + _('Hour') + '(s)';

	tmp = Math.round((val / global.commonTs.minute) * 100) / 100;
	if (tmp >= 1)
		return tmp + ' ' + _('Minute') + '(s)';

	if (val >= 1) {
		val = Math.round(val * 100) / 100;
		return val + ' ' + _('Second') + '(s)';
	}

	tmp = Math.round((val * 1000) * 100) / 100;
	if (tmp >= 1)
		return tmp + ' ' + _('Millisecond') + '(s)';

	tmp = Math.round((val * 1000000) * 100) / 100;
	if (tmp >= 1)
		return tmp + ' ' + _('Microsecond') + '(s)';

	return val;
};

rdr_elapsed_time = function(timestamp, full_length) {
	timestamp = parseInt(timestamp);

	var elapsed = parseInt(new Date().getTime() / 1000) - timestamp;

	var elapsed_text = elapsed + ' seconds ago';

	if (elapsed < 3)
		elapsed_text = 'just now';
	if (elapsed > 60)
		elapsed_text = parseInt(elapsed / 60) + ' mins ago';
	if (!full_length) {
		if (elapsed > 3600)
			elapsed_text = rdr_tstodate(timestamp);
	}else {
		if (elapsed > 3600)
			elapsed_text = parseInt(elapsed / 3600) + ' hours ago';
		if (elapsed > 86400)
			elapsed_text = parseInt(elapsed / 86400) + ' days ago';
	}

	return elapsed_text;
};

rdr_tags = function(tags) {
	var html = '';
	if (tags) {
		if (tags.length > 0) {
			html += "<ul class='tags'>";
			for (var i in tags)
				html += "<li><a href='#'>" + tags[i] + '</a></li>';
			html += '</ul>';
		}
	}
	return html;
};

rdr_display_groups = function(groups) {
	output = '';
	for (var i in groups) {
		var group = rdr_clean_id(groups[i]);

		output += group;

		if (i != (groups.length - 1))
			output += ',';
	}
	return output;
};


rdr_country = function(val ) {
	dicCountry = {
		'france': 'fr',
		'usa' : 'us',
		'espagne': 'es'
	};

	if (dicCountry[val.toLowerCase()])
		return '<span class=\"icon icon-country-' + dicCountry[val.toLowerCase()] + '\" />';
	else return val;
};

rdr_os = function(val ) {
	return '<span class=\"icon icon-os-' + val.toLowerCase() + '\" />';
};

rdr_browser = function(val ) {
	return '<span class=\"icon icon-browser-' + val.toLowerCase() + '\" />';
};

rdr_duration = function(timestamp, nb) {

	if (! nb && nb != 0)
		nb = 99;

	if (timestamp == 0)
		return 0;

	var times = [
		[global.commonTs.year,	'y',	0],
		[global.commonTs.month,	'M',	0],
		[global.commonTs.week,	'w',	0],
		[global.commonTs.day,	'd',	0],
		[global.commonTs.hours,	'h',	0],
		[global.commonTs.minute,	'm',	0],
		[1,	's',	0],
		[0.001,	'ms',	1],
		[0.000001,	'us',	1]
	];

	var output = '';
	var j = 0;
	for (var i in times) {
		if (timestamp >= times[i][0]) {
			var value = parseInt(timestamp / times[i][0]);
			timestamp -= value * times[i][0];
			if (output != '')
				output += ' ';

			output += value + times[i][1];
			j += 1;
		}

		if (timestamp == 0 || j >= nb || (times[i][2] && output))
			break;
	}

	return output;
};

rdr_yaxis = function(ori_value, multiple, decimal) {

	if (ori_value == 0)
		return 0;

	if (! multiple || parseInt(multiple) == NaN)
		multiple = 1000;

	if (! decimal || parseInt(decimal) == NaN)
		decimal = 2;


	var rounderer = Math.pow(10, decimal);
	var output = ori_value;

	value = Math.round(ori_value * rounderer) / rounderer;
	if (value < multiple)
		return value;

	value = Math.round((ori_value / multiple) * rounderer) / rounderer;
	if (value >= 1) {
		output = value + 'K';
		ori_value = value;
	}

	value = Math.round((ori_value / multiple) * rounderer) / rounderer;
	if (value >= 1) {
		output = value + 'M';
		ori_value = value;
	}

	value = Math.round((ori_value / multiple) * rounderer) / rounderer;
	if (value >= 1) {
		output = value + 'G';
		ori_value = value;
	}

	value = Math.round((ori_value / multiple) * rounderer) / rounderer;
	if (value >= 1) {
		output = value + 'T';
		ori_value = value;
	}

	return output;
};

rdr_humanreadable_value = function(value, unit) {
	if (value == 0)
		return 0;

	var multiple = 1000;

	if (! unit || unit == undefined) {
		unit = '';
	}else {

		if (unit == 'o' || unit == 'MB')
			multiple = 1024;

		if (unit == 'MB') {
			unit = 'B';
			value = value * 1024 * 1024;
		}

		if (unit == 'ms' || unit == 's') {
			if (unit == 'ms')
				value = value / 1000;
			return rdr_duration(value);
		}
	}

	value = rdr_yaxis(value, multiple);
	return value + unit;
};
