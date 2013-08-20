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
Ext.define('canopsis.controller.Ldap', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Ldap.Form'],
	stores: ['Ldaps'],
	models: ['Ldap'],

	logAuthor: '[controller][ldap]',

	//ditMethod: 'tab',

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.formXtype = 'LdapForm';
		this.listXtype = undefined;

		this.modelId = 'Ldap';

		this.callParent(arguments);
	},

	accountButton: function(){
		var store = this.getStore('Ldaps');

		this.grid = this.getController('Account').grid;
		
		store.on('load', function() {
			var record = store.getAt(0);

			log.debug('Edit:', this.logAuthor);
			log.dump(record);

			//if (record)
				this._editRecord(undefined, record, 0, store)

		}, this, {single: true});

		store.load();
	},


	preSave: function(record, data, form) {
		record.set('id', "ldap.config");
		//record.set('groups', groups);
		/*
		{name: 'aaa_access_group', defaultValue: undefined},
		{name: 'aaa_access_other', defaultValue: undefined},
		{name: 'aaa_access_owner', defaultValue: undefined},
		{name: 'aaa_admin_group', defaultValue: undefined},
		{name: 'aaa_group', defaultValue: undefined},
		{name: 'aaa_owner', defaultValue: undefined},
		*/
		return record;
	}

});
