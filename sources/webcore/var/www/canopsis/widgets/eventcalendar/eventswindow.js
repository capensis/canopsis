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

Ext.define('widgets.eventcalendar.eventswindow' , {
	extend: 'canopsis.lib.view.cpopup',

	alias: 'widget.eventcalendar.eventswindow',

	height: 400,
	width: 400,
	layout: 'fit',
	items: {
		border: false,
	},
	closeAction: 'hide',

	addMode: true,

	initComponent: function() {
		this.calendar = this.initialConfig.calendar;

		this.callParent(arguments);
	},

	_buildForm: function() {
		//Title
		this._form.add({
			xtype: 'grid',
    		title: 'Events',
    		columns: [
        		{ text: 'Component',  dataIndex: 'component' },
        		{ text: 'Ressource', dataIndex: 'ressource', flex: 1 },
        		{ text: 'Output', dataIndex: 'output' }
    		],
		});

	},

	afterRender: function() {
		this.callParent(arguments);
	},

	showEvents : function(calEvent){
		this.show();
	}
});