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
Ext.define('canopsis.view.LiveSearch.View', {
    extend: 'Ext.form.Panel',
    alias: 'widget.LiveSearch',

	store: ['Inventory'],

	layout: {
			type: 'vbox',
			align: 'stretch'
		},
		border: false,
		bodyPadding: 10,

		fieldDefaults: {
			labelAlign: 'top',
			labelWidth: 100,
			labelStyle: 'font-weight:bold'
		},
		defaults: {
			margins: '0 0 10 0'
		},

		items: [{
			xtype: 'fieldcontainer',
			fieldLabel: _('Search Option'),
			labelStyle: 'font-weight:bold;padding:0',
			layout: 'hbox',
			defaultType: 'textfield',

			fieldDefaults: {
				labelAlign: 'top'
			},

			items: [{
				flex: 1,
				//name: 'firstName',
				itemId: 'source_name',
				fieldLabel: _('Source Name')
			},{
				flex: 1,
				//name: 'lastName',
				itemId: 'type',
				fieldLabel: _('Type'),
				margins: '0 0 0 5'
			},{
				flex: 1,
				//name: 'lastName',
				itemId: 'source_type',
				fieldLabel: _('Source type'),
				margins: '0 0 0 5'
			},{
				flex: 1,
				//name: 'lastName',
				itemId: 'component',
				fieldLabel: _('Component'),
				margins: '0 0 0 5'
			},{
				xtype: 'button',
				flex: 1,
				text: 'search',
				itemId: 'LiveSearchButton',
				margins: '0 0 0 5'
			}]
		}, {
			xtype: 'LiveGrid'
			/*
			xtype: 'grid',
			flex: 1,
			margins: '0',
			store : ['Inventory'],
			columns : [
				{header : 'name', dataIndex : '_id', flex : 1},
				{header : 'type', dataIndex : 'source_type', flex : 1},
			],
			/*store : Ext.create('Ext.data.Store', {
				fields: ['name']
			})*/

		}]
/*
		buttons: [{
			text: _('Cancel'),
		}, {
			text: _('Send'),
			}
		}]
		*/
	});
