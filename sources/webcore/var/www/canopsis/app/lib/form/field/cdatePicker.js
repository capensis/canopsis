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

Ext.define('canopsis.lib.form.field.cdatePicker', {
	extend: 'Ext.form.field.Date',
	mixins: ['canopsis.lib.form.cfield'],
	alias: 'widget.cdatePicker',

	getValue: function() {
		var inputDate = canopsis.lib.form.field.cdatePicker.superclass.getValue.call(this);
		inputDate = new Date(inputDate);

		return inputDate.getTime() / 1000;
	},

	setValue: function(_date) {
		if(Ext.isNumber(_date)) {
			_date = new Date(cleanTimestamp(_date) * 1000);
		}

		canopsis.lib.form.field.cdatePicker.superclass.setValue.call(this, _date);
	}
});
