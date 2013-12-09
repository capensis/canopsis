//need:app/lib/view/cform.js
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
Ext.define('canopsis.view.Ldap.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.LdapForm',

	iconCls: 'icon-book',

	logAuthor: '[Controller][Ldap][Form]',

	layout: 'fit',

	width: 500,
	height: 350,

	items: [{
		xtype: 'fieldset',
		defaultType: 'textfield',
		border: false,
		defaults: {
			width: 450,
			labelWidth: 150
		},
		items: [{
			xtype: 'checkboxfield',
			fieldLabel: _('Enable'),
			name: 'enable',
			inputValue: true,
			uncheckedValue: false
		},{
			fieldLabel: _('URI'),
			name: 'uri',
			allowBlank: false
		},{
			fieldLabel: _('base_dn'),
			name: 'base_dn',
			allowBlank: false
		},{
			fieldLabel: _('user_dn'),
			name: 'user_dn',
			allowBlank: true
		},{
			fieldLabel: _('domain'),
			name: 'domain',
			allowBlank: true
		},{
			fieldLabel: _('user_filter'),
			name: 'user_filter',
			allowBlank: false
		},{
			fieldLabel: _('field_lastname'),
			name: 'lastname',
			allowBlank: false
		},{
			fieldLabel: _('field_firstname'),
			name: 'firstname',
			allowBlank: false
		},{
			fieldLabel: _('field_mail'),
			name: 'mail',
			allowBlank: false
		},{
			fieldLabel: _('field_group'),
			name: 'aaa_group',
			allowBlank: false
		},{
			xtype: 'label',
			text: _('Warning, you should restart the webserver to apply changes!')
		}]
	}]
});
