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
function init_REST_Store(collection, selector, groupField) {
	var options = {};
	log.debug("Init REST Store, Collection: '" + collection + "', selector: '" + selector + "', groupField: '" + groupField + "'");

	options['storeId'] = collection + selector;
	options['id'] = collection + selector;
	//options['model'] = Ext.create('canopsis.model.'+collection)
	//options['model'] = 'canopsis.model.'+collection
	options['model'] = Ext.ModelMgr.getModel('canopsis.model.' + collection);
	if (groupField) {
		options['groupField'] = groupField;
	}

	var store = Ext.create('canopsis.store.Mongo-REST', options);
	store.proxy.url = '/webservices/rest/' + collection + '/' + selector;

	return store;
}

//Ajax action
var ajaxAction = function(url, params, cb, scope, method) {
	if (!method)
		method = 'GET';

	var options = {
		method: method,
		url: url,
		scope: scope,
		success: cb,
		params: params,
		failure: function(result, request) {
			log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
		}
	};
	Ext.Ajax.request(options);
};

// Create Global "extend" method
var extend = function(obj, extObj) {
    if (arguments.length > 2) {
        for (var a = 1; a < arguments.length; a++) {
            extend(obj, arguments[a]);
        }
    } else {
        for (var i in extObj) {
            obj[i] = extObj[i];
        }
    }
    return obj;
};

var random_id = function() { return Math.floor(Math.random() * 11)};

//find the greatest common divisor
function find_gcd(nums)
{
        if (!nums.length)
                return 0;
        for (var r, a, i = nums.length - 1, GCDNum = nums[i]; i;)
                for (a = nums[--i]; r = a % GCDNum; a = GCDNum, GCDNum = r);
        return GCDNum;
}

// Split AMQP Routing key
function split_amqp_rk(rk) {
    var srk = rk.split('.');

    if (srk[2] == 'check') {
        var component;
        var resource;
        if (srk[3] == 'resource') {
            var expr = /^(\w*)\.(\w*)\.(\w*)\.(\w*)\.(.*)\.([\w\-]*)$/g;
            var result = expr.exec(rk);
            if (result) {
                component = result[5];
                resource = result[6];
            }
        }else {
            var expr = /^(\w*)\.(\w*)\.(\w*)\.(\w*)\.(.*)$/g;
            var result = expr.exec(rk);
            if (result)
                component = result[5];
        }

        return {source_type: srk[3] , component: component, resource: resource};
    }
    return {};
}

function get_timestamp_utc(date) {
	if (! date)
		date = new Date();

	var localTime = parseInt(date.getTime() / 1000);
	var localOffset = parseInt(date.getTimezoneOffset() * 60);

	return localTime - localOffset;
}

function isEmpty(obj) {
	for (var prop in obj) {
		if (obj.hasOwnProperty(prop))
			return false;
	}
	return true;
}

function getPct(value, max, decimal) {
	if (! decimal)
		decimal = 2;

	if (max == 0)
		return 100;

	var div = Math.pow(10, decimal);

	return Math.round(((100 * value) / max) * div) / div;
}

function getMidnight(timestamp) {
	var time = new Date(timestamp);
	var new_time = timestamp - (time.getHours() * global.commonTs.hours * 1000);
	//floor to hour, time / hour * hour
	new_time = parseInt(new_time / (global.commonTs.hours * 1000)) * (global.commonTs.hours * 1000);
	return new_time;
}

function check_color(color) {
	if (! color)
		return color;

	if (color.charAt(0) == '#' && color.charAt(1) != '#')
		return color;

	//Clean color
	while(color.charAt(0) == '#')
		color = color.slice(1)

	return '#' + color;
}

function strip_blanks(val) {
	return val.replace(/\n/g, '').replace(/ /g, '');
}

function stringTo24h(src_time) {
	var time = src_time.split(' ');

	if (time.length > 1) {
		//---------Format 12h
		var hour_type = time[1];
		var clock = time[0];

		clock = clock.split(':');
		var minute = parseInt(clock[1], 10);
		var hour = parseInt(clock[0], 10);

		if (hour_type == 'am' && hour == 12)
			hour = 0;

		if (hour_type == 'pm' && hour != 12)
			hour = hour + 12;

	} else {
		//--------Format 24h
		var time = src_time.split(':');

		var minute = time[1];
		var hour = time[0];
	}

	return {minute: parseInt(minute, 10), hour: parseInt(hour, 10)};
}

var updateRecord = function(namespace, crecord_type, model, _id, data, on_success, on_error) {
	var logAuthor = '[tools][updateRecord]';

	if (! data) {
		log.error('You must specify data to write', logAuthor);
		return;
	}

	var base_url = '/rest/' + namespace + '/' + crecord_type + '/' + _id;

	log.debug('Update ' + _id, logAuthor);
	Ext.Ajax.request({
		url: base_url,
		jsonData: data,
		method: 'PUT',
		success: function(operation) {
			log.debug(' + Success', logAuthor);
			global.notify.notify(_('Saved'), _('Successfully'), 'success');
			if (on_success)
				on_success(operation);
		},
		failure: function(response) {
			log.error(' + Impossible to deal with webservice', logAuthor);
			global.notify.notify(_('Error'), _('Imposible to deal with webservice, record not saved.'), 'error');
			if (on_error)
				on_error();
		}

	});
};

var demultiplex_cps_state = function(cps_state) {
	var state = cps_state.toString();
	if (state.length == 2)
		return {state: 0, state_type: state[0], state_extra: state[1]};
	else if (state.length == 3)
		return {state: state[0], state_type: state[1], state_extra: state[2]};
	else
		return undefined;
};

function getMetaId(component, resource, metric ) {
	if (!resource || resource == null)
		var name = component + metric;
	else
		var name = component + resource + metric;

	//console.log('-------------------'+name+'----------------------')
	return $.md5(name);
}

function split_search_box(raw) {
	// Split search string by space
	var search_value_array = [];

	var tmp = raw.split('"');
	if (tmp.length > 1) {
		for (var i = 0; i < tmp.length; i++) {
			var w = tmp[i];
			if (w.length > 1) {
				if (w[0] == ' ')
					w = w.slice(1);
				if (w[w.length - 1] == ' ')
					w = w.slice(0, w.length - 1);

				search_value_array.push(w);
			}
		}
	}else {
		search_value_array = raw.split(' ');
	}

	return search_value_array;
}

function is12Clock() {
	if (global.is12Clock != undefined)
		return global.is12Clock;

	if (global.account != undefined) {

		if (global.account.clock_type && global.account.clock_type != 'auto')
			global.is12Clock = global.account.clock_type == '12h';
		else
			global.is12Clock = Ext.Array.contains(global.am_pm_lang, global.locale);

		return global.is12Clock;

	}else {
		return global.default_is12Clock;
	}
}

function getTimeRegex() {
	if (is12Clock() == true)
		return /([01]?\d)(:\d{2})(\s)?(am|pm)?$/;
	else
		return /^([01]?\d|2[0-3]):?([0-5]\d)$/;

}

// Check if record exist
function isRecordExist(namespace, crecord_type, field, record, callback, scope) {
		var filter = {};
		filter[field] = record.get(field);
		filter = {
			filter: Ext.encode(filter),
			limit: 1
		};

		Ext.Ajax.request({
			method: 'GET',
			scope: scope,
			params: filter,
			url: '/rest/' + namespace + '/' + crecord_type,
			success: function(response, opts) {
				var data = Ext.decode(response.responseText).data;
				callback(this, record, data.length == 0);
			},
			failure: function(response) {
				log.error(' + Impossible to deal with webservice', '[tools][isExist]');
				callback(this, record, false);
			}
		});
}

function parseBool(arg) {
	if (arg == 0 || arg == null || arg == undefined || arg == false)
		return false;
	else
		return true;
}

function roundSignifiantDigit(value, sig) {
	var mult = Math.pow(10, sig) 
	value = Math.round(value * mult)
	value = value / mult
    return value;
}

function sciToDec(number) {
	val = number;
	if (Ext.isNumber(number))
		val = number.toString();

	if (val.match(/^[-+]?[1-9]\.[0-9]+e[-]?[1-9][0-9]*$/)) {
		var arr = new Array();
		arr = scinum.split('e');
		var exponent = Math.abs(arr[1]);
		var precision = new Number(exponent);
		arr = arr[0].split('.');
		precision += arr[1].length;
		val = (+val).toFixed(precision);
	}
	return val;
}

function cleanTimestamp(number){
	if(Ext.isNumber(number))
		number = parseInt(number,10).toString()
	if(number.length > 12){
		var cleaned_timestamp = parseInt(number,10)
		return parseInt(cleaned_timestamp/1000,10)
	}else{
		return parseInt(number,10)
	}
}