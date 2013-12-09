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

Ext.define('canopsis.lib.form.field.cfieldset' , {
	extend: 'Ext.form.FieldSet',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.cfieldset',

	checkboxToggle: false,
	inputValue: true,
	collapsible: false,
	collapsed: false,
	border:false,
	defaults: { labelWidth: 200 },

	getName: function() {
		return this.checkboxName;
	},

	getSubmitData: function() {
		var data = {};
		data[this.checkboxName] = this.getValue();
		return data;
	},

	initComponent: function() {
		if(!this.name) {
			this.name = this.checkboxName;
		}

		if(this.name !== undefined) {
			this.collapsed = true;
			this.checkboxToggle = true;
		}

		if(this.value === true) {
			this.collapsed = false;
		}

		//don't move this, otherwise it won't work
		this.style = {'border-width': "1px 0px 0px 0px"};

		this.callParent(arguments);
	},

	getValue: function() {
		if(this.checkboxCmp) {
			var value = this.checkboxCmp.getValue();

			if(value) {
				return true;
			}
			else {
				return false;
			}
		}
		else {
			return this.value;
		}
	},

	setValue: function(value) {
		if(value === undefined && this.value) {
			value = this.value;
		}

		if(value === undefined) {
			value = false;
		}

		this.value = value;

		if(this.checkboxCmp) {
			this.checkboxCmp.setValue(value);
		}
		else {
			this.collapsed = !value;

			if(!value) {
				this.collapse();
			}
			else {
				this.expand();
			}
		}
	},

	createCheckboxCmp: function() {
		var checkbox = this.callParent(arguments);
		checkbox.isFormField = false;
		checkbox.uncheckedValue = false;

		return checkbox;
	}
});
