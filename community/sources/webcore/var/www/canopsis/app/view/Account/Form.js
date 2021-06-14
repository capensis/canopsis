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
Ext.define('canopsis.view.Account.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.AccountForm',

	iconCls: 'icon-crecord_type-account',

	logAuthor: '[Controller][Account][Form]',

	layout: 'fit',

	width: 350,
	height: 350,

	items: [{
		xtype: 'tabpanel',
		plain: true,
		border: false,
		deferredRender: false,
		defaults: {
			border: false,
			autoScroll: true
		},
		items: [{
			title: _('General options'),
			defaultType: 'textfield',
			bodyStyle: 'padding:5px 5px 0',
			layout: 'anchor',
			border: false,
			items: [{
				xtype: 'fieldset',
				defaultType: 'textfield',
				border: false,
				items: [{
					fieldLabel: _('Login'),
					name: 'user',
					allowBlank: false,
					regex: /^[A-Za-z0-9_]+$/,
					regexText: _('Invalid login') + ', ' + _('use alphanumeric characters only') + '<br/>([A-Za-z0-9_])'
				},{
					fieldLabel: _('First Name'),
					name: 'firstname',
					allowBlank: false
				}, {
					fieldLabel: _('Last Name'),
					name: 'lastname',
					allowBlank: false
				},{
					fieldLabel: _('E-mail'),
					name: 'mail',
					vtype: 'email',
					allowBlank: true
				},{
					fieldLabel: _('Group'),
					name: 'aaa_group',
					store: 'Groups',
					value: 'Canopsis',
					displayField: 'crecord_name',
					valueField: '_id',
					editable: false,
					xtype: 'combobox',
					allowBlank: false
				},{
					fieldLabel: _('Password'),
					inputType: 'password',
					name: 'passwd'
				}]

			}]
		},{
			title: _('Groups'),
			xtype: 'grid',
			store: 'Groups',
			autoScroll: true,
			selType: 'checkboxmodel',
			selModel: {
				mode: 'MULTI',
				showHeaderCheckbox: false
			},
			columns: [
				{text: _('Name'), dataIndex: 'crecord_name', flex: 1}
			],
			columnLines: true,
			hideHeaders: true,
			name: 'groups'
		}]
	}]
});
