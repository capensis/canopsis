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

Ext.define('canopsis.auth' , {
	extend: 'Ext.form.Panel',

	url:'/auth',
	frame:true,

	logAuthor: '[auth]',

	title: _('Authentification'),

	renderTo: 'auth_form',

	fieldDefaults: {
		msgTarget: 'side',
		labelWidth: 75
	},
	defaultType: 'textfield',
	defaults: {
		anchor: '100%'
	},

	items: [{
		fieldLabel: _('Username'),
		name: 'login',
		allowBlank:false
	},{
		fieldLabel: _('Password'),
		name: 'password',
		id: 'password',
		inputType: 'password',
		allowBlank:false
	}],

	on_authed: undefined,

	initComponent: function() {
		log.debug("Bind enter key", this.logAuthor)
		this.navkeys = Ext.create('Ext.util.KeyNav', Ext.getDoc(), {
			scope: this,
			enter : this.submit
		});

		this.buttons = [{
			text: _('Connect'),
			id: 'submitbutton',
			scope: this,
			submitOnEnter : false,
			handler : function(){this.submit();}
		}],

		this.callParent(arguments);
	},

	submit: function() {
		log.debug("Sublit form", this.logAuthor)

		var form = this.getForm()
		if (form.isValid()){
			var login = form.getFieldValues().login
			var passwd = form.getFieldValues().password;
			var passwd_sha1 = $.encoding.digests.hexSha1Str(passwd);

			var timestamp = Math.floor(Ext.Date.format(new Date(), 'U') / 10) * 10
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
				success: function(response){
					response = Ext.JSON.decode(response.responseText)
					log.debug(" + Auth Ok", this.logAuthor)
					global.account = response.data[0]
					if (this.on_authed)
						this.on_authed(global.account)
				},
				failure: function(form, action) {
					log.debug(" + Auth Failed", this.logAuthor)
					Ext.MessageBox.alert(_('authentification failed'), _('Login or password invalid'))
				}
			});

		}else{
			log.debug("+ Form is invalid", this.logAuthor)
			Ext.Msg.alert(_('Invalid Data'), _('Please correct form errors.'))
		}	
	}

});