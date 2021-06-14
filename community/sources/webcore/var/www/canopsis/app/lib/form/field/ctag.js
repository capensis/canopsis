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

Ext.define('canopsis.lib.form.field.ctag' , {
	extend: 'Ext.panel.Panel',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.ctag',

	border: false,

	layout: 'hbox',

	initComponent: function() {
		this.logAuthor = '[' + this.id + ']';
		log.debug('Initialize ...', this.logAuthor);

		// store
		this.operator_store = Ext.create('Ext.data.Store', {
			fields: ['operator', 'text', 'type'],
			data: [
				{'operator': '$in', 'text': _('In')},
				{'operator': '$all', 'text': _('Match all')}
			]
		});

		// create operator combo
		this.operator_combo = Ext.widget('combobox', {
			flex: 1,
			queryMode: 'local',
			displayField: 'text',
			editable: false,
			value: '$in',
			valueField: 'operator',
			store: this.operator_store
		});

		// textArea
		this.textArea = Ext.widget('textfield', {
			margin: '0 0 0 5',
			flex: 3,
			emptyText: _('Type your tags here')
		});

		this.items = [this.operator_combo, this.textArea];

		this.callParent(arguments);
	},

	getValue: function() {
		log.debug('Get value', this.logAuthor);
		var rawString = this.textArea.getValue();
		var separator = undefined;

		if(Ext.Array.contains(rawString, ';')) {
			separator = ',';
		}
		else if(Ext.Array.contains(rawString, ',')) {
			separator = ';';
		}
		else if(Ext.Array.contains(rawString, ' ')) {
			separator = ' ';
			rawString = rawString.replace(/  +/g, ' ');
		}

		if(!separator) {
			return undefined;
		}

		log.debug('String separator is: "' + separator + '"', this.logAuthor);

		if(separator !== ' ') {
			rawString = strip_blanks(rawString);
		}

		var tag_array = rawString.split(separator);

		var filter = {'tags': {}};
		filter['tags'][this.operator_combo.getValue()] = tag_array;

		var output = Ext.encode(filter);
		log.debug('Generated filter is: ' + output, this.logAuthor);

		return Ext.encode(filter);
	},

	setValue: function(value) {
		value = Ext.decode(value);

		if(value['tags']) {
			value = value['tags'];
			var operator = Ext.Object.getKeys(value)[0];
			var value_array = value[operator];
			var tags_string = '';

			for(var i = 0; i < value_array.length; i++) {
				tags_string = tags_string + ' ' + value_array[i];
			}

			this.operator_combo.setValue(operator);
			this.textArea.setValue(tags_string);
		}
	}
});
