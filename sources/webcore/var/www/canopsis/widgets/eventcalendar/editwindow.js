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

Ext.define('widgets.eventcalendar.editwindow' , {
	extend: 'canopsis.lib.view.cpopup',

	alias: 'widget.eventcalendar.editwindow',

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
			xtype: 'displayfield',
			value: _('Event title') + ':'
		});

		this._form.add({
			xtype: 'textfield',
			name: 'event_title',
			itemId: 'event_title',
			anchor: '100%',
			emptyText: _('Type here the event title')
		});

		//StartDate
		this._form.add({
			xtype: 'displayfield',
			value: _('Start date') + ':'
		});

		this._form.add({
			xtype: 'datefield',
			name: 'start_date',
			itemId: 'start_date'
		});

		this._form.add({
			xtype: 'timefield',
			name: 'start_time',
			itemId: 'start_time'
		});


		//EndDate
		this._form.add({
			xtype: 'displayfield',
			value: _('End date') + ':'
		});

		this._form.add({
			xtype: 'datefield',
			name: 'end_date',
			itemId: 'end_date'
		});

		this._form.add({
			xtype: 'timefield',
			name: 'end_time',
			itemId: 'end_time'
		});

		//AllDay
		this._form.add({
			xtype: 'displayfield',
			value: _('All day event') + ':'
		});

		this._form.add({
			xtype: 'checkbox',
			name: 'all_day'
		});
	},

	afterRender: function() {
		this.callParent(arguments);
	},

	showNewEvent : function(start, end, allDay){
		this.addMode = true;

		this._form.down("#event_title").setValue("");
		this._form.down("#start_date").setValue(start);
		this._form.down("#end_date").setValue(end);

		this.show();

		return { title: "aaaaaa"};
	},

	showEditEvent: function(event){
		this.addMode = false;

		this._form.down("#event_title").setValue(event.title);
		this._form.down("#start_date").setValue(new Date(event.start));
		this._form.down("#end_date").setValue(new Date(event.end));

		this.show();
	},

	ok_button_function: function(){
		var newEvent = {};
		newEvent.title = this._form.down("#event_title").getValue();
		newEvent.start = this._form.down("#start_date").getValue();
		newEvent.end = this._form.down("#end_date").getValue();
		newEvent.allDay = this._form.all_day;

		if(this.addMode == true)
		{
			//add the new event to the calendar
			this.calendar.add_events([newEvent]);
		}

		this.hide();
	}
});