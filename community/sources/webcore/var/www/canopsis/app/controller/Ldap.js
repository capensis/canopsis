//need:app/lib/controller/cgrid.js,app/view/Ldap/Form.js,app/store/Ldaps.js
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
Ext.define('canopsis.controller.Ldap', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Ldap.Form'],
	stores: ['Ldaps'],
	models: ['Ldap'],

	logAuthor: '[controller][ldap]',

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

			if(record) {
				this._editRecord(undefined, record, 0, store);
			}
			else {
				this._showForm(undefined, store);
			}

		}, this, {single: true});

		store.load();
	},


	preSave: function(record) {
		record.set('id', "ldap.config");

		return record;
	},

	_save: function(record, edit, unused_store, form) {
		void(unused_store);

		//override the parent _save method, to manually choose the store to save data in
		var store = this.getStore('Ldaps');
		console.log(store);

		var batch = undefined;

		if(edit) {
			batch = store.proxy.batch({update: [record]});
		}
		else {
			batch = store.proxy.batch({create: [record]});
		}

		batch.on('complete', function(batch, operation, opts) {
			this.displaySuccess(batch, operation, opts);
			log.debug('Reload store', this.logAuthor);
			this.load();
		}, store);

		this._postSave(record);
		this._cancelForm(form);
	}
});
