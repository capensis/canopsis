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
Ext.define('canopsis.view.Account.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.AccountForm',

	iconCls: 'icon-crecord_type-account',

	logAuthor: '[Controller][Account][Form]',

	//layout : 'hbox',

	initComponent: function() {
		log.debug('Initializing...', this.logAuthor);


		//---------------------------General options-------------------------
		var g_options = [{
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
				store: Ext.create('canopsis.store.Groups'),
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
				//allowBlank : false
			}];



		var g_options_panel = Ext.widget('fieldset', {
				title: _('General options'),
				defaultType: 'textfield',
				items: g_options
			});

		//----------------------- drag and drop-------------------
		var checkboxModel = Ext.create('Ext.selection.CheckboxModel');
		this.checkGrid = Ext.create('Ext.grid.Panel', {
			store: 'Groups',
			autoScroll: true,
			height: 200,
			selModel: checkboxModel,
			columns: [
				{text: _('Name'), dataIndex: 'crecord_name', flex: 1}
			],
			columnLines: true,
			title: _('Groups'),
			hideHeaders: true,
			collapsible: true,
			//collapsed : true,
			name: 'groups',
			scroll: 'vertical',
			margin: '4 0 6 0'
		});

		var secondary_group = Ext.widget('fieldset', {
				title: _('Secondary groups'),
				items: [this.checkGrid]
			});

		this.callParent(arguments);

		this.add([g_options_panel, secondary_group]);
	}

});
