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
Ext.define('canopsis.view.Curves.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	controllerId: 'Curves',

	alias: 'widget.CurvesGrid',

	model: 'Curve',
	store: 'Curves',

	opt_paging: true,
	opt_menu_delete: true,
	opt_bar_duplicate: true,
	opt_menu_rights: true,
	opt_bar_search: true,
	opt_bar_search_field: ['metric'],

	columns: [
		{
			header: _('Line Color'),
			sortable: false,
			align: 'center',
			dataIndex: 'line_color',
			renderer: rdr_color
		},{
			header: _('Area color'),
			sortable: false,
			align: 'center',
			dataIndex: 'area_color',
			renderer: rdr_color
		},{
			header: _('Line style'),
			sortable: false,
			align: 'center',
			dataIndex: 'dashStyle'
		},{
			header: _('Area opacity'),
			sortable: false,
			align: 'center',
			dataIndex: 'area_opacity'
		},{
			header: _('zIndex'),
			sortable: false,
			align: 'center',
			dataIndex: 'zIndex'
		},{
			header: _('Invert'),
			sortable: false,
			align: 'center',
			dataIndex: 'invert',
			renderer: rdr_boolean
		},{
			header: _('Metric name'),
			flex: 7,
			sortable: true,
			dataIndex: 'metric'
		},{
			header: _('Label'),
			flex: 7,
			sortable: true,
			dataIndex: 'label'
		}
	]
});
