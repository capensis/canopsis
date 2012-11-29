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
Ext.define('canopsis.controller.Account', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Account.Grid', 'Account.Form'],
	stores: ['Accounts'],
	models: ['Account'],

	iconCls: 'icon-crecord_type-account',

	logAuthor: '[controller][Account]',

	init: function() {
		log.debug('[' + this.id + '] - Initialize ...');

		this.formXtype = 'AccountForm';
		this.listXtype = 'AccountGrid';

		this.modelId = 'Account';

		this.callParent(arguments);

		global.accountCtrl = this;
	},

	getConfig: function(id, default_value) {
		log.debug(' + getConfig ' + id, this.logAuthor);
		if (! global.account[id]) {
			if (global[id]) {
				global.account[id] = global[id];
			}else {
				global.account[id] = default_value;
			}
		}
		return global.account[id];
	},

	setConfig: function(id, value, cb) {
		log.debug(' + setConfig ' + id + ' => ' + value, this.logAuthor);
		global.account[id] = value;

		var url = '/account/setConfig/' + id;

		if (! cb) {
			cb = function() {
				log.debug(' + setConfig Ok', this.logAuthor);
			};
		}

		ajaxAction(url, {value: value}, cb, this, 'POST');

		return global.account[id];
	},

	setLocale: function(locale) {
		var cb = function() {
			log.debug(' + setLocale Ok', this.logAuthor);
			Ext.MessageBox.show({
				title: _('Configure language'),
				msg: _('Application must be reloaded, do you want to reload now ?'),
				icon: Ext.MessageBox.WARNING,
  				buttons: Ext.Msg.OKCANCEL,
  				fn: function(btn) {
					if (btn == 'ok') {
						window.location.reload();
					}
				}
			});
		};

		this.setConfig('locale', locale, cb);
	},

	setClock: function(clock_type) {
		var cb = function() {
			log.debug(' + setClock Ok', this.logAuthor);
			Ext.MessageBox.show({
				title: _('Configure clock display'),
				msg: _('Application must be reloaded, do you want to reload now ?'),
				icon: Ext.MessageBox.WARNING,
  				buttons: Ext.Msg.OKCANCEL,
  				fn: function(btn) {
					if (btn == 'ok') {
						window.location.reload();
					}
				}
			});
		};

		this.setConfig('clock_type', clock_type, cb);
	},

	beforeload_EditForm: function(form,item) {
		var pass_textfield = Ext.ComponentQuery.query('#' + form.id + ' textfield[name=passwd]')[0];
		if (pass_textfield) {
			pass_textfield.allowBlank = true;
		}

		var user_textfield = Ext.ComponentQuery.query('#' + form.id + ' textfield[name=user]')[0];
		if (user_textfield)
			user_textfield.hide();

		//------------groups-------------------
		var store = Ext.getStore('Groups');
		var groups = item.get('groups');
		if (groups.length > 0) {
			var groups_records = [];

			for (i in groups) {
				record = store.findExact('_id', groups[i]);
				if (record != -1)
					groups_records.push(store.getAt(record));
			}

			var selectMethod = form.checkGrid.getSelectionModel();
			selectMethod.select(groups_records);

		}
	},

	preSave: function(record,data,form) {
		//don't update password if it's empty
		if (form.editing && (record.get('passwd') == '')) {
			delete record.data.passwd;
		}

		//add groups
		var record_list = form.checkGrid.getSelectionModel().getSelection();
		var groups = [];
		for (i in record_list) {
			groups.push(record_list[i].get('crecord_name'));
		}

		record.set('groups', groups);

		//log.dump(record)

		//hack, otherwise webservice go wild (temporary)
		record.modified.aaa_group = data['aaa_group'];

		return record;
	},

	validateForm: function(store, data, form) {
		var already_exist = false;

		// in creation mode
		if (!form.editing)
			if (store.findExact('user', data['user']) >= 0)
				already_exist = true

			var field = form.findField('user')
			if (field)
				field.markInvalid(_("Invalid field"))

		if (already_exist) {
			log.debug('User already exist', this.logAuthor + '[validateForm]');
			global.notify.notify(data['user'] + ' already exist', 'you can\'t add the same user twice', 'error');
			return false;
		}else {
			return true;
		}
	},

	//check if user have right on this record
	check_record_right: function(record,option) {
		var user = global.account._id;
		var group = global.account.aaa_group;
		var groups = global.account.groups;

		//root can do everything
		if (this.checkRoot()) {
			return true;
		}

		if ((option == 'r') || (option == 'w')) {
			if ((user == record.get('aaa_owner')) && (record.data.aaa_access_owner.indexOf(option) > -1)) {
				//log.debug('owner')
				return true;
			} else if ((group == record.get('aaa_group')) && (record.data.aaa_access_group.indexOf(option) > -1)) {
				//log.debug('group')
				return true;
			} else if ((groups.indexOf(record.get('aaa_group')) != -1) && (record.data.aaa_access_group.indexOf(option) > -1)) {
				//log.debug('group')
				return true;
			} else if ((groups.indexOf(record.get('aaa_admin_group')) != -1) || group == record.get('aaa_admin_group')) {
				return true;
			} else {
				//log.debug('nothing')
				return false;
			}
		} else {
			log.error(_('Incorrect right option given'), this.logAuthor);
		}
	},

	//check if user have right on this obj
	check_right: function(obj,option) {
		var user = global.account._id;
		var group = global.account.aaa_group;
		var groups = global.account.groups;

		//root can do everything
		if (this.checkRoot()) {
			return true;
		}

		if ((option == 'r') || (option == 'w')) {
			if ((user == obj.aaa_owner) && (obj.aaa_access_owner.indexOf(option) > -1)) {
				//log.debug('owner')
				return true;
			} else if ((group == obj.aaa_group) && (obj.aaa_access_group.indexOf(option) > -1)) {
				//log.debug('group')
				return true;
			} else if ((groups.indexOf(obj.aaa_group) != -1) && (obj.aaa_access_group.indexOf(option) > -1)) {
				//log.debug('group')
				return true;
			} else if ((groups.indexOf(obj.aaa_admin_group) != -1) || group == obj.aaa_admin_group) {
				return true;
			} else {
				//log.debug('nothing')
				return false;
			}
		} else {
			log.error(_('Incorrect right option given'), this.logAuthor);
		}

	},

	//if callback_func != null and ajax success -> callback is call in passed scope with
	//new key as argument
	new_authkey: function(account,callback_func,scope) {
		if (account) {
			//------------------------------ajax request----------------------
			log.debug('Ask webserver for new authentification key', this.logAuthor);
			Ext.Ajax.request({
				url: '/account/getNewAuthKey/' + account,
				method: 'GET',
				scope: scope,
				success: function(response) {
					var object_response = Ext.decode(response.responseText);
					if (object_response.success == true) {
						global.notify.notify(_('Success'), _('Your authentification key is updated'), 'success');
						var authkey = object_response.data.authkey;
						global.account.authkey = authkey;
						if (callback_func)
							callback_func.call(this, authkey);
					}else {
						log.error('Ajax output incorrect', this.logAuthor);
					}
				},
				failure: function(response) {
					global.notify.notify(_('Error'), _('An error have occured during the updating process'), 'error');
					log.error('Error while fetching new Authkey', this.logAuthor);
				}
			});
		}else {
			log.error('No account provided for Authkey');
		}
	},

	get_authkey: function(account,callback_func,scope) {
		log.debug('Ask webserver for authentification key', this.logAuthor);
		Ext.Ajax.request({
			url: '/account/getAuthKey/' + account,
			method: 'GET',
			scope: scope,
			success: function(response) {
				var object_response = Ext.decode(response.responseText);
				if (object_response.success == true) {
					var authkey = object_response.data.authkey;
					if (callback_func)
						callback_func.call(this, authkey);
				}else {
					log.error('Ajax output incorrect', this.logAuthor);
				}
			},
			failure: function(response) {
				global.notify.notify(_('Error'), _('An error have occured during the process'), 'error');
				log.error('Error while fetching new Authkey', this.logAuthor);
			}
		});
	},

	add_to_group: function(group,account)	{
		log.debug('Ask webserver adding ' + account + ' to ' + 'group', this.logAuthor);
		if (group.search('group.') == -1)
			group = 'group.' + group;
		if (account.search('account.') == -1)
			account = 'account.' + account;

		Ext.Ajax.request({
			url: '/account/addToGroup/' + group + '/' + account,
			method: 'POST',
			success: function(response) {
				var object_response = Ext.decode(response.responseText);
				if (object_response.success == true) {
					global.notify.notify(_('Group added'), _('Group successfuly added to secondary groups'), 'success');
				}else {
					log.error('Ajax output incorrect', this.logAuthor);
				}
			},
			failure: function(response) {
				global.notify.notify(_('Error'), _('An error have occured during the process'), 'error');
				log.error('Error while fetching new Authkey', this.logAuthor);
			}
		});
	},

	remove_from_group: function(group,account) {
		log.debug('Ask webserver removing ' + account + ' from ' + 'group', this.logAuthor);
		Ext.Ajax.request({
			url: '/account/removeFromGroup/' + group + '/' + account,
			method: 'POST',
			success: function(response) {
				var object_response = Ext.decode(response.responseText);
				if (object_response.success == true) {
					global.notify.notify(_('Group removed'), _('Group successfuly removed from secondary groups'), 'success');
				}else {
					log.error('Ajax output incorrect', this.logAuthor);
				}
			},
			failure: function(response) {
				global.notify.notify(_('Error'), _('An error have occured during the process'), 'error');
				log.error('Error while fetching new Authkey', this.logAuthor);
			}
		});
	},

	checkGroup: function(group) {
		if (global.account.aaa_group == group || (global.account.groups.indexOf(group) != -1))
			return true;
		else
			return false;
	},

	checkRoot: function() {
		if (global.account.user == 'root' || global.account.aaa_group == 'group.CPS_root' || (global.account.groups.indexOf('group.CPS_root') != -1))
			return true;
		else
			return false;
	}

});
