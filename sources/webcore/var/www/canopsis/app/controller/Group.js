//need:app/lib/controller/cgrid.js,app/view/Group/Grid.js,app/view/Group/Form.js,app/store/Groups.js
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
Ext.define('canopsis.controller.Group', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Group.Grid', 'Group.Form'],
	stores: ['Groups'],
	models: ['Group'],

	iconCls: 'icon-crecord_type-group',

	logAuthor: '[controller][Group]',

	checkInternal: true,

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.formXtype = 'GroupForm';
		this.listXtype = 'GroupGrid';

		this.modelId = 'Group';

		this.callParent(arguments);
	},

	preSave: function(record) {
		record.data.id = 'group.' + record.data.crecord_name;
		return record;
	},

	beforeload_EditForm: function(form) {
		var field = Ext.ComponentQuery.query('#' + form.id + ' textfield[name=crecord_name]')[0];

		if(field) {
			field.hide();
		}
	},

	validateForm: function(store, data, form) {
		var already_exist = false;

		// in creation mode
		if(!form.editing && store.findExact('crecord_name', data['crecord_name']) >= 0) {
			already_exist = true;
		}

		var field = form.findField('crecord_name');
		if(field) {
			field.markInvalid(_('Invalid field'));
		}

		if(already_exist) {
			log.debug('Group already exist', this.logAuthor + '[validateForm]');
			global.notify.notify(data['crecord_name'] + ' already exist', 'you can\'t add the same group twice', 'error');
			return false;
		}
		else {
			return true;
		}
	}
});
