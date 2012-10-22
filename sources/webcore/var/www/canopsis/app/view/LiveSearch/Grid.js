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
Ext.define('canopsis.view.LiveSearch.Grid' , {
	extend: 'Ext.grid.Panel',
	alias: 'widget.LiveGrid',

	store: 'Inventory',
	id: 'LiveGrid',

	requires: [
		'Ext.grid.plugin.CellEditing',
		'Ext.form.field.Text',
		'Ext.toolbar.TextItem'
	],

	title: '',
	//iconCls: 'icon-grid',
	//frame: true,
	features: [Ext.create('Ext.grid.feature.Grouping', {
        groupHeaderTpl: 'Source type: {source_type}'
    })],

	border: false,

	columns: [
				{header: 'type', dataIndex: 'source_type', flex: 1},
				{header: 'name', dataIndex: '_id', flex: 1}
	],

	initComponent: function() {
		this.callParent(arguments);
	}

});
