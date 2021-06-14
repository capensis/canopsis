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
Ext.define('canopsis.view.Account.Password', {
	extend: 'Ext.window.Window',

	alias: 'widget.Password',

	title: _('Change your password'),
	height: 120,
	width: 300,
	layout: 'fit',

	logAuthor: '[view][Password]',

	items: [{
		xtype: 'form',
		bodyPadding: '5 5 5 5',
		defaultType: 'textfield',
		defaults: {
			allowBlank: false,
			listeners: {
				specialkey: function(field, event) {
					if(event.getKey() === event.ENTER) {
						field.up('form').getForm().submit();
					}
				}
			}
		},

		items: [{
			fieldLabel: _('Password'),
			inputType: 'password',
			name: 'password'
		},{
			fieldLabel: _('Confirm'),
			inputType: 'password',
			name: 'confirm'
		}],

		submit: function() {
			var win = this.owner.up('window');
			var form = this;
			var values = form.getValues();

			if(values["password"] !== values["confirm"]) {
				form.markInvalid({"confirm": _("Pasword doesn't match")});
				return;
			}

			if(!form.isValid()) {
				return;
			}

			global.accountCtrl.setPassword(values["password"]);
			win.close();
		},

		buttons: [{
			text: 'Submit',
			handler: function() {
				this.up('form').getForm().submit();
			}
		}]
	}]
});