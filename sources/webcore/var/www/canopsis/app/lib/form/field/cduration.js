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

Ext.define('canopsis.lib.form.field.cduration' , {
	extend: 'Ext.form.FieldContainer',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.cduration',

	layout: 'hbox',

	number_step: 1,
	number_min_val: 0,
	value: global.commonTs.day,
	add_none_value: false,

	fieldLabel: undefined,

	initComponent: function() {
		this.logAuthor = '[' + this.id + ']';
		log.debug('Initialize ...', this.logAuthor);

		this.ts_window = Ext.widget('numberfield', {
			isFormField: false,
			name: 'ts_window',
			width: 50,
			value:this.value,
			minValue: this.number_min_val,
			step: this.number_step
		});

		var store_data = [
			{'name': _('Minute'), 'value': global.commonTs.minute},
			{'name': _('Hour'), 'value': global.commonTs.hours},
			{'name': _('Day'), 'value': global.commonTs.day},
			{'name': _('Week'), 'value': global.commonTs.week},
			{'name': _('Month'), 'value': global.commonTs.month},
			{'name': _('Year'), 'value': global.commonTs.year}
		];

		if(this.add_none_value) {
			store_data.push({'name': _('None'), 'value': undefined});
		}

		this.ts_unit = Ext.widget('combobox', {
			isFormField: false,
			editable: false,
			width: 97,
			name: 'ts_unit',
			queryMode: 'local',
			displayField: 'name',
			valueField: 'value',
			store: {
				xtype: 'store',
				fields: ['value', 'name'],
				data: store_data
			}
		});

		this.items = [this.ts_window, this.ts_unit];

		this.callParent(arguments);

		this.setValue(this.value);

	},

	checkDisable: function(combo, value) {
		void(combo);

		if(value[0].raw.value) {
			this.ts_window.setDisabled(false);
		}
		else {
			this.ts_window.setDisabled(true);
		}
	},

	show: function() {
		this.callParent(arguments);
		this.ts_window.show();
		this.ts_window.setDisabled(false);
		this.ts_unit.show();
		this.ts_unit.setDisabled(false);
	},

	hide: function() {
		this.callParent(arguments);
		this.ts_window.hide();
		this.ts_window.setDisabled(true);
		this.ts_unit.hide();
		this.ts_unit.setDisabled(true);
	},

	getValue: function() {
		if(!this.ts_unit.getValue()) {
			return undefined;
		}

		return this.ts_window.getValue() * this.ts_unit.getValue();
	},

	setValue: function(value) {
		if(value) {
			if(value >= global.commonTs.year) {
				this.ts_unit.setValue(global.commonTs.year);
				this.ts_window.setValue(Math.floor(value / global.commonTs.year));
			}
			else if(value >= global.commonTs.month) {
				this.ts_unit.setValue(global.commonTs.month);
				this.ts_window.setValue(Math.floor(value / global.commonTs.month));
			}
			else if(value >= global.commonTs.week) {
				this.ts_unit.setValue(global.commonTs.week);
				this.ts_window.setValue(Math.floor(value / global.commonTs.week));
			}
			else if(value >= global.commonTs.day) {
				this.ts_unit.setValue(global.commonTs.day);
				this.ts_window.setValue(Math.floor(value / global.commonTs.day));
			}
			else if(value >= global.commonTs.hours) {
				this.ts_unit.setValue(global.commonTs.hours);
				this.ts_window.setValue(value / global.commonTs.hours);
			}
			else {
				this.ts_unit.setValue(global.commonTs.minute);
				this.ts_window.setValue(value / global.commonTs.minute);
			}
		}
		else {
			this.ts_unit.select(this.ts_unit.getStore().data.items[0]);
		}
	}
});
