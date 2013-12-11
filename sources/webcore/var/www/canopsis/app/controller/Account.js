//need:app/lib/controller/cgrid.js,app/view/Account/Grid.js,app/view/Account/Form.js,app/store/Accounts.js
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

	bindGridEvents: function(grid){
		var ldapCtrl = this.getController('Ldap');

		// Bind ldap button
		var btns = Ext.ComponentQuery.query('#' + grid.id + ' button[action=ldap]');

		for(var i = 0; i < btns.length; i++) {
			btns[i].on('click', ldapCtrl.accountButton, ldapCtrl);
		}
	},

	getConfig: function(id, default_value) {
		log.debug(' + getConfig ' + id, this.logAuthor);
		if(!global.account[id]) {
			if(global[id]) {
				global.account[id] = global[id];
			}
			else {
				global.account[id] = default_value;
			}
		}

		return global.account[id];
	},

	setConfig: function(id, value, cb) {
		log.debug(' + setConfig ' + id + ' => ' + value, this.logAuthor);
		global.account[id] = value;

		var url = '/account/setConfig/' + id;

		if(!cb) {
			cb = function() {
				log.debug(' + setConfig Ok', this.logAuthor);
			};
		}

		ajaxAction(url, {value: value}, cb, this, 'POST');

		return global.account[id];
	},

	_deleteButton: function(button) {
		log.debug('Clicked deleteButton', this.logAuthor);
		var grid = this.grid;
		var me = this;

		var selection = grid.getSelectionModel().getSelection();

		if(selection) {
			//check right
			var ctrlAccount = this.getController('Account');
			var authorized = true;

			for(var i = 0; i < selection.length; i++) {
				if(!ctrlAccount.check_record_right(selection[i], 'w')) {
					authorized = false;
				}

				if(this.checkInternal && selection[i].get('internal')) {
					authorized = false;
				}

				if(!authorized) {
					break;
				}
			}

			if(authorized === true) {
				Ext.Msg.show({
					title:_('Are you sure ?'),
					inputSelection: selection,
					msg: _('When you delete an account ALL its object (view, curve...) is deleted with it. You can choose to just disable the account in order to keep its objects.'),
					buttons: Ext.Msg.YESNOCANCEL,
					buttonText: {
						yes: _('DELETE'),
						no: _('Disable')
					},
					scope: me,
					fn: function(buttonId, text, opt) {
						void(text);

						if(buttonId === 'yes') {
							this.grid.store.remove(opt.inputSelection);
						}
						if(buttonId === 'no') {
							this._enabledisable();
						}
					}
				});
			}
			else {
				global.notify.notify(_('Access denied'), _('You don\'t have the rights to modify this object'), 'error');
			}
		}

		if(this.deleteButton) {
			this.deleteButton(button, grid, selection);
		}

	},

	beforeload_DuplicateForm: function(form, copy) {
		void(form);

		/* remove unique data */
		copy.set('user', undefined);
		copy.set('firstname', undefined);
		copy.set('lastname', undefined);
		copy.set('mail', undefined);
	},

	afterload_DuplicateForm: function(form, copy) {
		this.form_loadGroups(form, copy);
	},

	logout: function() {
		Ext.Ajax.request({
			url: '/logout',
			scope: this,
			success: function() {
				log.debug(' + Success.', this.logAuthor);
				window.location.reload();
			},
			failure: function() {
				log.error("Logout impossible, maybe you're already logout");
			}
		});
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
					if(btn === 'ok') {
						if(global.minimified) {
							window.location.href = '/' + locale + '/';
						}
						else {
							window.location.href = '/' + locale + '/static/canopsis/index.debug.html';
						}
					}
				}
			});
		};

		Ext.util.Cookies.set('locale', locale);

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
					if(btn === 'ok') {
						window.location.reload();
					}
				}
			});
		};

		this.setConfig('clock_type', clock_type, cb);
	},

	beforeload_EditForm: function(form,item) {
		var pass_textfield = Ext.ComponentQuery.query('#' + form.id + ' textfield[name=passwd]')[0];

		if(pass_textfield) {
			pass_textfield.allowBlank = true;
		}

		var user_textfield = Ext.ComponentQuery.query('#' + form.id + ' textfield[name=user]')[0];

		if(user_textfield) {
			user_textfield.hide();
		}

		// groups
		this.form_loadGroups(form, item);
	},

	form_loadGroups: function(form, item) {
		var store = Ext.getStore('Groups');
		var groups = item.get('groups');

		if(groups.length > 0) {
			var groups_records = [];

			for(var i = 0; i < groups.length; i++) {
				record = store.findExact('_id', groups[i]);

				if(record !== -1) {
					groups_records.push(store.getAt(record));
				}
			}

			var groups_grid = form.down('grid');
			groups_grid.getSelectionModel().select(groups_records, false, true);
		}
	},

	preSave: function(record, data, form) {
		//don't update password if it's empty
		if(form.editing && (record.get('passwd') === '')) {
			delete record.data.passwd;
		}

		if(record.get('external') === true) {
			log.debug('Impossible to change password on external account', this.logAuthor);
			delete record.data.passwd;
		}

		//add groups
		var checkGrid = form.down('grid');
		var record_list = checkGrid.getSelectionModel().getSelection();
		var groups = [];

		for(var i = 0; i < record_list.length; i++) {
			groups.push(record_list[i].get('crecord_name'));
		}

		record.set('groups', groups);

		//hack, otherwise webservice go wild (temporary)
		record.modified.aaa_group = data['aaa_group'];

		return record;
	},

	validateForm: function(store, data, form) {
		var already_exist = false;

		if ( !store ) { 
			var store = Ext.getStore('Groups');
		}

		// in creation mode
		if(!form.editing && store.findExact('user', data['user']) >= 0) {
			already_exist = true;
		}

		var field = form.findField('user');

		if(field) {
			field.markInvalid(_('Invalid field'));
		}

		if(already_exist) {
			log.debug('User already exist', this.logAuthor + '[validateForm]');
			global.notify.notify(data['user'] + ' already exist', 'you can\'t add the same user twice', 'error');

			return false;
		}
		else {
			return true;
		}
	},

	//check if user have right on this record
	check_record_right: function(record, option) {
		var user = global.account._id;
		var group = global.account.aaa_group;
		var groups = global.account.groups;

		//root can do everything
		if(this.checkRoot()) {
			return true;
		}

		if((option === 'r') || (option === 'w')) {
			if((user === record.get('aaa_owner')) && (Ext.Array.contains(record.data.aaa_access_owner, option))) {
				return true;
			}
			else if((group === record.get('aaa_group')) && (Ext.Array.contains(record.data.aaa_access_group, option))) {
				return true;
			}
			else if((Ext.Array.contains(groups, record.get('aaa_group'))) && (Ext.Array.contains(record.data.aaa_access_group, option))) {
				return true;
			}
			else if((Ext.Array.contains(groups, record.get('aaa_admin_group'))) || (group === record.get('aaa_admin_group'))) {
				return true;
			}
			else {
				return false;
			}
		}
		else {
			log.error(_('Incorrect right option given'), this.logAuthor);
		}
	},

	//check if user have right on this obj
	check_right: function(obj, option) {
		var user = global.account._id;
		var group = global.account.aaa_group;
		var groups = global.account.groups;

		//root can do everything
		if(this.checkRoot()) {
			return true;
		}

		if((option === 'r') || (option === 'w')) {
			if((user === obj.aaa_owner) && (Ext.Array.contains(obj.aaa_access_owner, option))) {
				return true;
			}
			else if((group === obj.aaa_group) && (Ext.Array.contains(obj.aaa_access_group, option))) {
				return true;
			}
			else if((Ext.Array.contains(groups, obj.aaa_group)) && (Ext.Array.contains(obj.aaa_access_group, option))) {
				return true;
			}
			else if((Ext.Array.contains(groups, obj.aaa_admin_group)) || group === obj.aaa_admin_group) {
				return true;
			}
			else {
				return false;
			}
		}
		else {
			log.error(_('Incorrect right option given'), this.logAuthor);
		}

	},

	// if the callback function is not null and ajax request is successful,
	// then the callback is called in passed scope with new key as argument.
	new_authkey: function(account, callback_func, scope) {
		if(account) {
			// ajax request
			log.debug('Ask webserver for new authentification key', this.logAuthor);

			Ext.Ajax.request({
				url: '/account/getNewAuthKey/' + account,
				method: 'GET',
				scope: scope,
				success: function(response) {
					var object_response = Ext.decode(response.responseText);

					if(object_response.success === true) {
						global.notify.notify(_('Success'), _('Your authentification key is updated'), 'success');

						var authkey = object_response.data.authkey;
						global.account.authkey = authkey;

						if(callback_func) {
							callback_func.call(this, authkey);
						}
					}
					else {
						log.error('Ajax output incorrect', this.logAuthor);
					}
				},
				failure: function() {
					global.notify.notify(_('Error'), _('An error have occured during the updating process'), 'error');
					log.error('Error while fetching new Authkey', this.logAuthor);
				}
			});
		}
		else {
			log.error('No account provided for Authkey');
		}
	},

	get_authkey: function(account, callback_func, scope) {
		log.debug('Ask webserver for authentification key', this.logAuthor);

		Ext.Ajax.request({
			url: '/account/getAuthKey/' + account,
			method: 'GET',
			scope: scope,
			success: function(response) {
				var object_response = Ext.decode(response.responseText);

				if(object_response.success === true) {
					var authkey = object_response.data.authkey;

					if(callback_func) {
						callback_func.call(this, authkey);
					}
				}
				else {
					log.error('Ajax output incorrect', this.logAuthor);
				}
			},
			failure: function() {
				global.notify.notify(_('Error'), _('An error have occured during the process'), 'error');
				log.error('Error while fetching new Authkey', this.logAuthor);
			}
		});
	},

	add_to_group: function(group, account) {
		log.debug('Ask webserver adding ' + account + ' to ' + 'group', this.logAuthor);

		if(group.search('group.') === -1) {
			group = 'group.' + group;
		}

		if (account.search('account.') === -1) {
			account = 'account.' + account;
		}

		Ext.Ajax.request({
			url: '/account/addToGroup/' + group + '/' + account,
			method: 'POST',
			success: function(response) {
				var object_response = Ext.decode(response.responseText);

				if(object_response.success === true) {
					global.notify.notify(_('Group added'), _('Group successfuly added to secondary groups'), 'success');
				}
				else {
					log.error('Ajax output incorrect', this.logAuthor);
				}
			},
			failure: function() {
				global.notify.notify(_('Error'), _('An error have occured during the process'), 'error');
				log.error('Error while fetching new Authkey', this.logAuthor);
			}
		});
	},

	remove_from_group: function(group, account) {
		log.debug('Ask webserver removing ' + account + ' from ' + 'group', this.logAuthor);

		Ext.Ajax.request({
			url: '/account/removeFromGroup/' + group + '/' + account,
			method: 'POST',
			success: function(response) {
				var object_response = Ext.decode(response.responseText);

				if(object_response.success === true) {
					global.notify.notify(_('Group removed'), _('Group successfuly removed from secondary groups'), 'success');
				}
				else {
					log.error('Ajax output incorrect', this.logAuthor);
				}
			},
			failure: function() {
				global.notify.notify(_('Error'), _('An error have occured during the process'), 'error');
				log.error('Error while fetching new Authkey', this.logAuthor);
			}
		});
	},

	checkGroup: function(group) {
		if(global.account.aaa_group === group || (Ext.Array.contains(global.account.groups, group))) {
			return true;
		}
		else {
			return false;
		}
	},

	checkRoot: function() {
		if(global.account.user === 'root' || global.account.aaa_group === 'group.CPS_root' || (Ext.Array.contains(global.account.groups, 'group.CPS_root'))) {
			return true;
		}
		else {
			return false;
		}
	},

	setPassword: function(password){
		var passwd = $.encoding.digests.hexSha1Str(password);

		this.setConfig("shadowpasswd", passwd, function() {
			global.notify.notify(_('Success'), _('Your password is updated'), 'success');
		});
	},

	setAvatar: function(file_id, filename) {
		var me = this;

		var extension = filename.split('.').pop();

		if($.inArray(extension, ['png', 'jpeg', 'jpg', 'gif']) === -1) {
			global.notify.notify(_('Failed'), "File extension not valid", 'error');
			return;
		}

		var avatar_id = file_id;

		log.debug('Set avatar_id in backend', this.logAuthor);

		Ext.Ajax.request({
			method: 'POST',
			url: '/account/setConfig/avatar_id',
			params: {
				value: avatar_id
			},
			success: function() {
				global.notify.notify(_('Success'), _('Avatar setted'), 'success');
				log.debug(' + Done', me.logAuthor);

				global.account.avatar_id = avatar_id;

				// Update icon in main bar
				log.debug('Refresh Mainbar', me.logAuthor);
				Ext.ComponentQuery.query('button[iconCls="icon-mainbar icon-avatar-bar"]')[0].setIcon('/account/getAvatar');
			},
			failure: function() {
				global.notify.notify(_('Failed'), _('Impossible to set this file as avatar'), 'error');
				log.error('Impossible to set this file as avatar', me.logAuthor);
			}
		});
	}
});
