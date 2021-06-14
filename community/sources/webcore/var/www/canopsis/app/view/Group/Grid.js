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
Ext.define('canopsis.view.Group.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.GroupGrid',

	model: 'Group',
	store: 'Groups',

	opt_menu_delete: true,
	opt_paging: false,
	opt_menu_rights: false,
	opt_allow_edit: false,
	opt_cell_edit: true,

	columns: [
		{
			header: '',
			width: 25,
			sortable: false,
			renderer: rdr_crecord_type,
			dataIndex: 'crecord_type'
		},{
			header: _('Name'),
			flex: 2,
			sortable: true,
			dataIndex: 'crecord_name'
		},{
			header: _('Description'),
			flex: 2,
			sortable: true,
			dataIndex: 'description',
			editor: {xtype: 'textfield'}
		}

	],

	initComponent: function() {
		this.callParent(arguments);
	}
});
