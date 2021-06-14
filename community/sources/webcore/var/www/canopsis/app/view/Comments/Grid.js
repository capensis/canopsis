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
Ext.define('canopsis.view.Comments.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.CommentsGrid',

	model: 'Comment',
	store: 'Comments',

	opt_db_namespace: 'object',

	opt_menu_delete: true,
	opt_bar_duplicate: true,
	opt_menu_rights: true,
	opt_bar_enable: true,
	opt_paging: false,

	opt_bar_search: true,
	opt_bar_search_field: ['comment'],


	columns: [
	{
		header: _('Comment'),
		sortable: true,
		flex: 1,
		dataIndex: 'comment'
	}],

	initComponent: function() {
		this.callParent(arguments);
	}

});
