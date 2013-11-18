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

Ext.define('canopsis.auth', {
	extend: 'Ext.form.Panel',

	url: '/auth',
	frame: true,

	title: _('Authentication'),

	renderTo: 'auth_form',

	fieldDefaults: {
		msgTarget: 'side',
		labelWidth: 100
	},
	defaultType: 'textfield',
	defaults: {
		anchor: '100%'
	},

	items: [],

	on_authed: undefined,

	initComponent: function() {
		var default_login = '';
		var default_password = '';

		if(window.location.host === "demo-devel.canopsis.org" || window.location.host === "demo.canopsis.org") {
			default_login    = 'root';
			default_password = 'root';
		}

		this.items = [
			{
				fieldLabel: _('Username'),
				name: 'login',
				allowBlank: false,
				value: default_login
			},{
				fieldLabel: _('Password'),
				name: 'password',
				id: 'password',
				inputType: 'password',
				allowBlank: false,
				value: default_password
			},{
				xtype: 'combo',
				name: 'locale',
				queryMode: 'local',
				displayField: 'text',
				valueField: 'value',
				fieldLabel: _('Language'),
				value: ENV['locale'],
				store: {
					xtype: 'store',
					fields: ['value', 'text'],
					data: [
							{'value': 'fr', 'text': 'Français'},
							{'value': 'en', 'text': 'English'}
					]
				},
				iconCls: 'no-icon',
				listeners: {
					select: function(combo, records) {
						void(combo);

						if(records.length) {
							locale = records[0].get('value');

							if(locale !== ENV['locale']) {
								Ext.util.Cookies.set('locale', locale);

								if(global.minimified) {
									window.location.href = '/' + locale + '/';
								}
								else {
									window.location.href = '/' + locale + '/static/canopsis/index.debug.html';
								}
							}
						}
					}
				}
			}
		];

		log.debug("Bind enter key", this.logAuthor);
		this.navkeys = Ext.create('Ext.util.KeyNav', Ext.getDoc(), {
			scope: this,
			enter : this.submit
		});

		this.buttons = [{
			text: _('Connect'),
			id: 'submitbutton',
			scope: this,
			submitOnEnter : false,
			handler : function() {
				this.submit();
			}
		}];

		this.callParent(arguments);
	},

	onSuccess: function() {
		this.navkeys.destroy();

		if(this.on_authed) {
			this.on_authed(global.account);
		}
	},

	onFailure: function() {
		log.debug(" + Auth Failed", this.logAuthor);

		Ext.MessageBox.alert(_('Authentification failed'), _('Password invalid or account disabled'));
	},

	auth_m1: function(login, passwd, passwd_sha1) {
		var timestamp = Math.floor(Ext.Date.format(new Date(), 'U') / 10) * 10;
		var authkey = $.encoding.digests.hexSha1Str(passwd_sha1 + timestamp.toString());

		Ext.Ajax.request({
			method: 'GET',
			url: this.url,
			scope: this,
			params: {
				cryptedKey: 'True',
				password: authkey,
				login: login
			},
			success: function(response) {
				response = Ext.JSON.decode(response.responseText);
				log.debug(" + M1 Auth Ok", this.logAuthor);
				global.account = response.data[0];
				this.onSuccess();
			},
			failure: function() {
				log.debug(" + M1 Failed", this.logAuthor);
				this.auth_m2(login, passwd, passwd_sha1);
			}
		});
	},

	auth_m2: function(login, passwd, passwd_sha1) {
		//relaunch auth with sha 1 mdp
		Ext.Ajax.request({
			method: 'GET',
			url: this.url,
			scope: this,
			params: {
				shadow: 'True',
				password: passwd_sha1,
				login: login
			},
			success: function(response) {
				response = Ext.JSON.decode(response.responseText);
				log.debug(" + M2 Auth Ok", this.logAuthor);
				global.account = response.data[0];
				this.onSuccess();
			},
			failure: function() {
				log.debug(" + M2 Failed", this.logAuthor);

				if(global['auth_plain']) {
					this.auth_m3(login, passwd, passwd_sha1);
				}
				else {
					this.onFailure();
				}
			}
		});
	},

	// WARNING: Plain method /!\
	auth_m3: function(login, passwd) {
		Ext.Ajax.request({
			method: 'GET',
			url: this.url,
			scope: this,
			params: {
				password: passwd,
				login: login
			},
			success: function(response) {
				response = Ext.JSON.decode(response.responseText);
				log.debug(" + M3 Auth Ok", this.logAuthor);
				global.account = response.data[0];
				this.onSuccess();
			},
			failure: function() {
				log.debug(" + M3 Failed", this.logAuthor);
				this.onFailure();
			}
		});
	},

	submit: function() {
		log.debug("Submit form", this.logAuthor);

		var form = this.getForm();

		if(form.isValid()) {
			var FieldValues = form.getFieldValues();
			var login       = FieldValues.login;
			var passwd      = FieldValues.password;
			var passwd_sha1 = $.encoding.digests.hexSha1Str(passwd);

			this.auth_m1(login, passwd, passwd_sha1);
		}
		else {
			log.debug("+ Form is invalid", this.logAuthor);
			Ext.Msg.alert(_('Invalid Data'), _('Please correct form errors.'));
		}
	}
});

function checkAuth(callback) {
	//check Auth
	log.debug('Check auth ...', "[index]");

	Ext.Ajax.request({
		type: 'rest',
		url: '/account/me',
		reader: {
			type: 'json',
			root: 'data',
			totalProperty  : 'total',
			successProperty: 'success'
		},
		success: function(response) {
			request_state = Ext.JSON.decode(response.responseText).success;

			if(request_state) {
				global.account = Ext.JSON.decode(response.responseText).data[0];

				log.debug(' + Logged', "[app]");
				callback(global.account);
			}
			else {
				log.debug(' + Please loggin', "[index]");
				createAuthForm(callback);
			}
		},
		failure: function(response) {
			if(response.status === 403) {
				log.debug(' + Please loggin', "[index]");
			}
			else {
				log.debug(' + Error in request', "[index]");
			}

			createAuthForm(callback);
		}
	});
}

function createAuthForm(callback) {
	Ext.get('auth').setVisible(true, true);
	Ext.create('canopsis.auth', {
		on_authed: callback
	});
}