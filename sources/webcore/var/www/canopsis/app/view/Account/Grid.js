//need:app/lib/view/cgrid.js
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

Ext.define('canopsis.view.Account.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	controllerId: 'Account',

	alias: 'widget.AccountGrid',

	model: 'Account',
	store: 'Accounts',

	opt_grouping: true,
	opt_paging: false,
	opt_bar_duplicate: true,
	opt_menu_delete: true,
	opt_menu_rights: false,
	opt_menu_authKey: true,
	opt_bar_enable: true,

	opt_bar_customs: [{
		text: 'Ldap',
		xtype: 'button',
		iconCls: 'icon-book',
		action: 'ldap'
	}],

	columns: [
		{
			header: '',
			width: 25,
			sortable: false,
			renderer: rdr_crecord_type,
			dataIndex: 'crecord_type'
		},{
			header: '',
			width: 25,
			sortable: false,
			renderer: function(val) {
				if(val === true) {
					return "<span class='icon icon-book_link' />";
				}
			},
			dataIndex: 'external'
		},{
			header: _('Enabled'),
			align: 'center',
			width: 55,
			dataIndex: 'enable',
			renderer: rdr_boolean
		},{
			header: _('Login'),
			flex: 2,
			sortable: true,
			dataIndex: 'user'
		},{
			header: _('First name'),
			flex: 2,
			sortable: false,
			dataIndex: 'firstname'
		},{
			header: _('Last name'),
			flex: 2,
			sortable: false,
			dataIndex: 'lastname'
		},{
			header: _('Email'),
			flex: 2,
			sortable: false,
			dataIndex: 'mail'
		},{
			header: _('Group'),
			flex: 2,
			sortable: false,
			dataIndex: 'aaa_group',
			renderer: rdr_clean_id
		},{
			header: _('Groups'),
			flex: 2,
			sortable: false,
			dataIndex: 'groups',
			renderer: rdr_display_groups
		}
	]
});
