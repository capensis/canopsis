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
function rdr_tstodate(val) {
	if(val) {
		var dval = new Date(parseInt(val) * 1000);

		if(is12Clock()) {
			return Ext.Date.format(dval, 'Y-m-d h:i:s a');
		}
		else {
			return Ext.Date.format(dval, 'Y-m-d H:i:s');
		}
	}
	else {
		return '';
	}
}

function rdr_utcToLocal(val) {
	if(val !== undefined && val !== '') {
		//format date
		var array_split = val.split(' ');
		var date = array_split[0].split('-');
		var hour = array_split[1].split(':');

		//create date
		var dval = new Date(date[0], date[1] - 1, date[2], hour[0], hour[1], hour[2]);

		return rdr_tstodate(get_timestamp_utc(dval));
	}
}

function rdr_boolean(val) {
	if(val) {
		return "<span class='icon icon-true' />";
	}
	else {
		return "<span class='icon icon-false' />";
	}
}

function rdr_status(val) {
	if(typeof(val) === 'number') {
		return "<span class='icon icon-state-" + val + "' />";
	}

	return val;
}

function rdr_color(val) {
	return "<span class='icon' style='background-color: #" + val + ';color: #' + val + ";'/>";
}

function rdr_state_type(val) {
	return "<span class='icon icon-state-type-" + val + "' />";
}

function rdr_source_type(val) {
	return "<span class='icon icon-crecord_type-" + val + "' />";
}

function rdr_crecord_type(val) {
	if(val !== '') {
		return "<span class='icon icon-crecord_type-" + val + "' />";
	}
}

function rdr_file_type(val) {
	if(!val) {
		return "<span class='icon icon-unknown' />";
	}

	var split = val.split('/');
	var extension = undefined;

	if(split.length > 0) {
		extension = split[split.length - 1];
	}
	else {
		extension = "unknown";
	}

	if($.inArray(extension, ['jpg', 'jpeg', 'gif', 'png']) !== -1) {
		return "<span class='icon icon-mimetype-png' />";
	}

	if ($.inArray(extension, ['ogg']) !== -1) {
		return "<span class='icon icon-mimetype-video' />";
	}

	if (split.length > 0) {
		return "<span class='icon icon-mimetype-" + extension + "' />";
	}
}

function rdr_havePerfdata(val) {
	if (val !== '') {
		return "<span class='icon icon-perfdata'/>";
	}
}

function rdr_widget_preview(val, metadata, record, rowIndex) {
	void(val, metadata, record);

	return "<span style='background-color:" + global.default_colors[rowIndex] + ';color:' + global.default_colors[rowIndex] + ";'>__</span>";
}

function rdr_task_crontab(val) {
	var output = '';

	if(val !== undefined) {
		//second condition is if minutes are str and not int
		if(val.hour !== undefined && val.minute !== undefined) {
			var d = new Date();
			d.setUTCHours(parseInt(val.hour, 10));
			d.setUTCMinutes(parseInt(val.minute, 10));

			var utc_minutes = d.getUTCMinutes();
			var local_minutes = d.getMinutes();
			var local_hours = d.getHours();

			//cosmetic
			if(utc_minutes < 10) {
				utc_minutes = '0' + utc_minutes;
			}

			if(local_minutes < 10) {
				local_minutes = '0' + local_minutes;
			}

			//12h translate
			if(!is12Clock()) {
				output += local_hours + ':' + local_minutes;
			}
			else {
				//utc AM/PM check
				if(local_hours > 12) {
					output += (local_hours - 12) + ':' + local_minutes + ' pm';
				}
				else {
					output += local_hours + ':' + local_minutes + ' am';
				}
			}

		}

		if(val.month !== undefined && val.day !== undefined) {
			output += '   |    ' + _('month') + ' : ' + global.numberToMonth[val.month] + ' |  day : ' + val.day;
		}

		if(val.day_of_week !== undefined) {
			output += '   |   ' + _('day') + ' : ' + _(val.day_of_week);
		}
	}

	return output;
}

//Function for rendering export to pdf button, we haven't find another solution
function rdr_export_button(val, metadata, record, rowIndex, colIndex, store, view) {
	void(val, metadata, rowIndex, colIndex, store);

	if(record.get('leaf')) {
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
}

function rdr_mail_information(val) {
	if(val === false) {
		return _('This task is not send by mail');
	}

	var output = '';

	if(val.recipients !== undefined) {
		output += _('Recipients :') + ' ' + val.recipients;
	}

	return output;
}

function rdr_clean_id(val) {
	if(val.search('.') !== -1) {
		var tmp = val.split('.');
		val = tmp[1];
	}

	return val;
}

function rdr_time_interval(val) {
	if(!val) {
		return '';
	}

	var tmp = undefined;

	tmp = Math.round((val / global.commonTs.year) * 100) / 100;

	if(tmp >= 1) {
		return tmp + ' ' + _('Year') + '(s)';
	}

	tmp = Math.round((val / global.commonTs.month) * 100) / 100;

	if(tmp >= 1) {
		return tmp + ' ' + _('Month') + '(s)';
	}

	tmp = Math.round((val / global.commonTs.week) * 100) / 100;

	if(tmp >= 1) {
		return tmp + ' ' + _('Week') + '(s)';
	}

	tmp = Math.round((val / global.commonTs.day) * 100) / 100;

	if(tmp >= 1) {
		return tmp + ' ' + _('Day') + '(s)';
	}

	tmp = Math.round((val / global.commonTs.hours) * 100) / 100;

	if(tmp >= 1) {
		return tmp + ' ' + _('Hour') + '(s)';
	}

	tmp = Math.round((val / global.commonTs.minute) * 100) / 100;

	if(tmp >= 1) {
		return tmp + ' ' + _('Minute') + '(s)';
	}

	if(val >= 1) {
		val = Math.round(val * 100) / 100;
		return val + ' ' + _('Second') + '(s)';
	}

	tmp = Math.round((val * 1000) * 100) / 100;

	if(tmp >= 1) {
		return tmp + ' ' + _('Millisecond') + '(s)';
	}

	tmp = Math.round((val * 1000000) * 100) / 100;

	if(tmp >= 1) {
		return tmp + ' ' + _('Microsecond') + '(s)';
	}

	return val;
}

function rdr_elapsed_time(timestamp, full_length) {
	timestamp = parseInt(timestamp);

	var elapsed = parseInt(new Date().getTime() / 1000) - timestamp;

	var elapsed_text = elapsed + ' ' + _('seconds ago');

	if(elapsed < 3) {
		elapsed_text = 'just now';
	}

	if(elapsed > 60) {
		elapsed_text = parseInt(elapsed / 60) + ' ' + _('minutes ago');
	}

	if(!full_length) {
		if (elapsed > 3600) {
			elapsed_text = rdr_tstodate(timestamp);
		}
	}
	else {
		if(elapsed > 3600) {
			elapsed_text = parseInt(elapsed / 3600) + ' ' + _('hours ago');
		}

		if(elapsed > 86400) {
			elapsed_text = parseInt(elapsed / 86400) + ' ' + _('days ago');
		}
	}

	return elapsed_text;
}

function rdr_tags(tags) {
	var html = '';

	if(tags && tags.length > 0) {
		html += "<ul class='tags'>";

		for(var i = 0; i < tags.length; i++) {
			html += "<li><a href='#'>" + tags[i] + '</a></li>';
		}

		html += '</ul>';
	}

	return html;
}

function rdr_ack(ackStatus) {
	//default value
	var status = 'confirm';
	//set existing status value, status change is delegated to button click control

	if(ackStatus) {
		status = ackStatus;
	}

	if(status === 'cancelled') {
		return '<span>'+ status +'</span>';
	}
	else {
		return '<button class="ackChangeStatus">' + status + '</button>';
	}
}

function rdr_display_groups(groups) {
	var output = '';

	for(var i = 0; i < groups.length; i++) {
		var group = rdr_clean_id(groups[i]);

		output += group;

		if(i !== (groups.length - 1)) {
			output += ',';
		}
	}

	return output;
}

function rdr_country(val) {
	var dicCountry = {
		'france': 'fr',
		'usa' : 'us',
		'espagne': 'es'
	};

	if(dicCountry[val.toLowerCase()]) {
		return '<span class=\"icon icon-country-' + dicCountry[val.toLowerCase()] + '\" />';
	}
	else {
		return val;
	}
}

function rdr_os(val) {
	return '<span class=\"icon icon-os-' + val.toLowerCase() + '\" />';
}

function rdr_browser(val) {
	return '<span class=\"icon icon-browser-' + val.toLowerCase() + '\" />';
}

function rdr_duration(timestamp, nb) {
	if(!nb && nb !== 0) {
		nb = 99;
	}

	if(timestamp === 0) {
		return 0;
	}

	var is_neg = false;

	if(timestamp < 0) {
		is_neg = true;
	}

	timestamp = Math.abs(timestamp);

	var times = [
		[global.commonTs.year,   'y',  0],
		[global.commonTs.month,  'M',  0],
		[global.commonTs.week,   'w',  0],
		[global.commonTs.day,    'd',  0],
		[global.commonTs.hours,  'h',  0],
		[global.commonTs.minute, 'm',  0],
		[1,                      's',  0],
		[0.001,                  'ms', 1],
		[0.000001,               'us', 1]
	];

	var output = '';
	var j = 0;

	for(var i = 0; i < times.length; i++) {
		if(timestamp >= times[i][0]) {
			var value = parseInt(timestamp / times[i][0]);

			timestamp -= value * times[i][0];

			if(output !== '') {
				output += ' ';
			}

			output += value + times[i][1];
			j += 1;
		}

		if(timestamp === 0 || j >= nb || (times[i][2] && output)) {
			break;
		}
	}

	return output;
}

function rdr_yaxis(ori_value, multiple, decimal) {
	if(ori_value === 0) {
		return 0;
	}

	if(!multiple || parseInt(multiple) === NaN) {
		multiple = 1000;
	}

	if(!decimal || parseInt(decimal) === NaN) {
		decimal = 2;
	}

	var rounderer = Math.pow(10, decimal);
	var output = ori_value;

	value = Math.round(ori_value * rounderer) / rounderer;

	if(value < multiple) {
		return sciToDec(roundSignifiantDigit(ori_value, 2));
	}

	value = Math.round((ori_value / multiple) * rounderer) / rounderer;

	if(value >= 1) {
		output = value + 'K';
		ori_value = value;
	}

	value = Math.round((ori_value / multiple) * rounderer) / rounderer;

	if(value >= 1) {
		output = value + 'M';
		ori_value = value;
	}

	value = Math.round((ori_value / multiple) * rounderer) / rounderer;

	if(value >= 1) {
		output = value + 'G';
		ori_value = value;
	}

	value = Math.round((ori_value / multiple) * rounderer) / rounderer;

	if(value >= 1) {
		output = value + 'T';
		ori_value = value;
	}

	return output;
}

function rdr_humanreadable_value(value, unit) {
	var is_neg = false;

	if(value < 0) {
		is_neg = true;
		value = value * -1;
	}

	if(value === 0) {
		return 0;
	}

	var multiple = 1000;

	if(!unit || unit === undefined) {
		unit = '';
	}
	else {
		if(unit === 'S' || unit === 's') {
			value = rdr_duration(value);

			if(is_neg) {
				value = '-' + value;
			}

			return value;
		}

		try {
			var information = global.sizeTable[unit.toUpperCase()];

			if(information) {
				multiple = information['multiple'];
				unit = information['unit'];
				value = value * information['pow'];
			}
		}
		catch(err) {
			log.debug(err.message);
		}
	}

	value = rdr_yaxis(value, multiple);

	if(is_neg) {
		value = '-' + value;
	}

	return value + unit;
}

function rdr_access(val) {
	if(Ext.isArray(val)) {
		return val.sort();
	}
	else {
		return val;
	}
}

function rdr_rule_action(val) {
	return "<span class='icon icon-rule-" + val + "' />";
}

rdr_file_help = function(val, metadata, record, rowIndex, colIndex, store) {
	if ( typeof record.raw['component'] != null && typeof record.raw['resource'] != null && this.opt_file_help_url != null && this.opt_file_help_url != "" ) {
		var reg_com=new RegExp( "<component>", "g" );
		var reg_res=new RegExp( "<resource>", "g" );
		return '<a target="_blank" href="' + this.opt_file_help_url.replace( reg_com, record.raw['component'] ).replace( reg_res, record.raw['resource'] ) + '" title="Support Directive"><span class="icon icon-file_help"></span></a>';
	} else {
		return ''
	}
};

rdr_file_equipement = function(val, metadata, record, rowIndex, colIndex, store) {
	if ( typeof  record.raw['component'] != null && typeof record.raw['resource'] != null && this.opt_file_equipement_url != null && this.opt_file_equipement_url != "" ) {
			var reg_com=new RegExp( "<component>", "g" );
			var reg_res=new RegExp( "<resource>", "g" );
		return '<a target="_blank" href="' + this.opt_file_equipement_url.replace( reg_com, record.raw['component'] ).replace( reg_res, record.raw['resource'] ) + '" title="Equipement Information"><span class="icon icon-file_equipement"></span></a>';
	} else {
		return ''
	}
};

rdr_ticket = function(val, metadata, record, rowIndex, colIndex, store) {
	if ( typeof record.raw['ticket'] != null && record.raw['ticket'] != undefined ) {
		var reg=new RegExp( "<ticket>", "g" );
		return '<a href="' + this.opt_ticket_url.replace( reg, record.raw['ticket'] ) + '" target="_blank">' + record.raw['ticket'] + '</a>'
	} else {
		return ''
	}
};

rdr_ack = function(val, metadata, record, rowIndex, colIndex, store ) {
	if ( typeof record.raw['ack_state'] != null && record.raw['ack_state'] != undefined && record.raw['ack_state'] > 0 ) {
		var cssClass = ""
		switch( record.raw['ack_state'] ) {
			case 1:
				cssClass = "icon-ack-pendingsolved";
				break;
			case 2:
				cssClass = "icon-ack-pendingaction";
				break;
			case 3:
				cssClass = "icon-ack-pendingvalidation";
				break;
		}
		return '<span class="icon '+cssClass+'" title="'+  record.raw['ack_output'] +'"></span>'
	} else {
		return ''
	}
}
