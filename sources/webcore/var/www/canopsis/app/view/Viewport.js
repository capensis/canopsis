//need:app/view/Mainbar/Bar.js,app/view/Tabs/View.js
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

Ext.define('canopsis.view.Viewport', {
	extend: 'Ext.container.Viewport',

	requires: [
		'Ext.layout.container.Border',
		'Ext.tab.Panel',
		'Ext.panel.Header',
		'canopsis.view.Mainbar.Bar',
		'canopsis.view.Tabs.View'
	],

	layout: 'border',

	logAuthor: '[Viewport]',

	items: [
		{
			region: 'north',
			id: 'region-north',
			border: false,
			height: 27,
			items: [{
				xtype: 'Mainbar',
				height: 26
			}],
			xtype: 'panel',
			collapsed: false,
			collapsible: true,
			preventHeader: true,
			placeholder: {
				xtype: 'header',
				height: 5,
				listeners: {
					mouseover: {
						element: 'el',
						fn: function() {
							Ext.getCmp('region-north').expand();
						}
					}
				}
			},
			listeners: {
				dblclick: {
					element: 'body',
					fn: function() {
						Ext.getCmp('region-north').collapse();
					}
				}
			}
		},{
			region: 'center',
			border: false,
			xtype: 'TabsView',
			id: 'main-tabs'
		}
	],

	initComponent: function() {
		log.debug('Render viewport ...', this.logAuthor);
		this.on('afterrender', this.afterrender, this);
		this.callParent(arguments);
	},

	afterrender: function() {
		log.debug('Viewport rendered', this.logAuthor);
	}

});
