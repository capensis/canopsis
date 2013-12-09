//need:app/lib/form/cfield.js
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

Ext.define('canopsis.lib.form.field.cdate', {
	extend: 'Ext.container.Container',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.cdate',

	layout: {
		type: 'hbox',
		align: 'stretch'
	},

	date_label_width: 40,
	date_width: 150,
	hour_width: 75,
	date_value: undefined,
	max_value: undefined,
	label_text: undefined,
	value: undefined,

	now: false,

	initComponent: function() {
		this.logAuthor = '[' + this.id + ']';
		log.debug('Initialize ...', this.logAuthor);

		if(!this.date_value) {
			this.date_value = new Date();
		}

		var config = {
			//upper form does not retrieve this element
			isFormField: false,

			labelWidth: this.date_label_width,
			value: this.date_value,
			editable: false,
			width: this.date_width,
			maxValue: this.max_value,
			allowBlank: false
		};

		if(this.label_text) {
			config.fieldLabel = this.label_text;
		}

		this.date = Ext.widget('datefield', config);

		config = {
			isFormField: false,
			margin: '0 0 0 2',
			width: this.hour_width,
			allowBlank: false,
			regex: getTimeRegex()
		};

		if(this.now) {
			if(is12Clock()) {
				config.value = Ext.Date.format(new Date, 'g:i a');
			}
			else {
				config.value = Ext.Date.format(new Date, 'G:i');
			}
		}
		else {
			if(is12Clock()) {
				config.value = Ext.Date.format(this.date_value, 'g:i a');
			}
			else {
				config.value = Ext.Date.format(this.date_value, 'G:i');
			}
		}

		this.hour = Ext.widget('textfield', config);

		this.items = [this.date, this.hour];

		this.callParent(arguments);

		this.relayEvents(this.date, ['select']);

		if(this.value) {
			this.setValue(this.value);
		}
	},

	getValue: function() {
		var date = parseInt(Ext.Date.format(this.date.getValue(), 'U'));
		var hour = stringTo24h(this.hour.getValue());

		var timestamp = date + (hour.hour * 60 * 60) + (hour.minute * 60);

		return parseInt(timestamp, 10);
	},

	setValue: function(value) {
		log.debug('cdate ' + this.name + ' is setValue with ' + value, this.logAuthor);
		this.date.setValue(new Date(value * 1000));

		if(is12Clock()) {
			this.hour.setValue(Ext.Date.format(new Date(value * 1000), 'g:i a'));
		}
		else {
			this.hour.setValue(Ext.Date.format(new Date(value * 1000), 'G:i'));
		}
	},

	setDisabled: function(bool) {
		this.callParent(arguments);
		this.date.setDisabled(bool);
		this.hour.setDisabled(bool);
	},

	isValid: function() {
		return (this.date.isValid() && this.hour.isValid());
	},

	setMaxDate: function(value) {
		this.date.setMaxValue(value);
	},

	setMinDate: function(value) {
		this.date.setMinValue(value);
	}
});
