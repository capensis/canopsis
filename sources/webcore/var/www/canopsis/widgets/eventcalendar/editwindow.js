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

	modal: true,
	closeAction: 'hide',

	addMode: true,

	currentEditedEvent: {},
	currentEditedEvent: null,

	initComponent: function() {
		this.calendar = this.initialConfig.calendar;

		this.callParent(arguments);
	},

	_buildForm: function() {
		// The data store containing the list of states

		this.sources = Ext.create('Ext.data.Store', {
			fields: ['name'],
		});

		this.calendar.computeIcsSources();

		for (var i = this.calendar.ics_sources_array.length - 1; i >= 0; i--) {
			this.sources.add({'name' : this.calendar.ics_sources_array[i]});
		};

		//Title
		this._form.add({
			fieldLabel: _('Event title'),
			xtype: 'textfield',
			name: 'event_title',
			itemId: 'event_title',
			anchor: '100%',
			emptyText: _('Type here the event title')
		});

		this._form.add({
			xtype: 'combobox',
			fieldLabel: _('Event source'),
			store: this.sources,
			itemId: 'event_source',
			name: 'event_source',
			editable: false,
			queryMode: 'local',
			displayField: 'name',
			valueField: 'name',
			emptyText: _('Type here the event source'),
			allowBlank: false
		});

		this._form.add({
			"xtype": "cfieldset",
			"title": _('Start'),
			"margin": 0,
			"padding": 0,
			"items": [{
					fieldLabel: _('Date'),
					xtype: 'datefield',
					name: 'start_date',
					itemId: 'start_date'
				},{
					fieldLabel: _('Time'),
					xtype: 'timefield',
					name: 'start_time',
					itemId: 'start_time'
				}]
		});

		this._form.add({
			"xtype": "cfieldset",
			"title": _('End'),
			"margin": 0,
			"padding": 0,
			"items": [{
					fieldLabel: _('Date'),
					xtype: 'datefield',
					name: 'end_date',
					itemId: 'end_date'
				},{
					fieldLabel: _('Time'),
					xtype: 'timefield',
					name: 'end_time',
					itemId: 'end_time'
				}]
		});
	},

	afterRender: function() {
		this.callParent(arguments);
	},

	showNewEvent : function(start, end, allDay){
		this.currentEditedEvent = {};
		this.currentEditedEventHtml = null;
		this.addMode = true;

		this._form.down("#event_title").setValue("");

		this._form.down("#event_source").setValue("");
		this._form.down("#event_source").setDisabled(false);

		this._form.down("#start_date").setValue(start);
		this._form.down("#end_date").setValue(end);

		if(!allDay)
		{
			this._form.down("#start_time").setValue(start);
			this._form.down("#end_time").setValue(end);
		}
		else
		{
			this._form.down("#start_time").setValue(null);
			this._form.down("#end_time").setValue(null);
		}

		this.show();
	},

	showEditEvent: function(event, eventHtml){
		this.currentEditedEvent = event;
		this.currentEditedEventHtml = eventHtml;
		this.addMode = false;

		this._form.down("#event_title").setValue(event.title);

		this._form.down("#event_source").setValue(event.component);
		this._form.down("#event_source").setDisabled(true);

		this._form.down("#start_date").setValue(new Date(event.start));

		if(event.end === null)
			event.end = event.start;

		this._form.down("#end_date").setValue(new Date(event.end));

		if(!event.allDay)
		{
			this._form.down("#start_time").setValue(event.start);
			this._form.down("#end_time").setValue(event.end);
		}
		else
		{
			this._form.down("#start_time").setValue(null);
			this._form.down("#end_time").setValue(null);
		}

		this.show();
	},

	ok_button_function: function(){
		var newEvent = {}; //TODO set this as a property to save hidden props

		newEvent.title = this._form.down("#event_title").getValue();

		if(this.currentEditedEvent.id)
			newEvent.id = this.currentEditedEvent.id;
		else
		{
			var now = new Date().toString('yyyy-MM-dd-hh-mm-ss');
			newEvent.id = newEvent.title + "-" + now + "@" + "widget-calendar";
		}

		combine = function(me, date, time, all_day) {
				if(Ext.isString(date)) date = me.parseDate(date);
				if(!date) date = new Date();
				if(Ext.isString(time)) time = me.parseDate(time);
				if(!time) time = new Date();
				var rv = new Date(date);
				if(!all_day)
				{
					rv.setHours(time.getHours());
					rv.setMinutes(time.getMinutes());
					rv.setSeconds(time.getSeconds());
				}
				return rv;
			}

		var startDateWidget = this._form.down("#start_date");
		var startTimeWidget = this._form.down("#start_time");
		var endDateWidget = this._form.down("#end_date");
		var endTimeWidget = this._form.down("#end_time");

		var start_datetime = combine(this, startDateWidget.getValue(), startTimeWidget.getValue());
		var end_datetime = combine(this, endDateWidget.getValue(), endTimeWidget.getValue());
		var all_day = startTimeWidget.getValue() === null || endTimeWidget.getValue() === null;

		newEvent.component = this._form.down("#event_source").getValue();

		newEvent.start = start_datetime;
		newEvent.end = end_datetime;

		var isHourfilledOnlyOnStartOrEnd = (startTimeWidget.getValue() === null && endTimeWidget.getValue() !== null)
											|| (startTimeWidget.getValue() !== null && endTimeWidget.getValue() === null)

		if(isHourfilledOnlyOnStartOrEnd)
			global.notify.notify(_('Form problem'), _('"You must fill start and end times, or nothing for an all day event"'), 'info');
		else if(newEvent.component === undefined || newEvent.component === null || newEvent.component === "")
			global.notify.notify(_('Form problem'), _('"You must specify a source for the event"'), 'info');
		else
		{
			newEvent.allDay = all_day;
			//add the new event to the calendar

			this.calendar.send_events([newEvent]);

			this.hide();
		}
	},

	hide: function() {
		this.calendar.resetEventStyle(this.currentEditedEventHtml);
		this.callParent();
	}

});