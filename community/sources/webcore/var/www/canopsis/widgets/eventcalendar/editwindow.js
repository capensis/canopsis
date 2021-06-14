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
	width: 500,
	layout: 'fit',

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
		var me = this;
		// The data store containing the list of states
		this.sources = Ext.create('Ext.data.Store', {
			fields: ['component'],
		});

		this.sources.add(this.calendar.sources);

		var tabs = Ext.create('Ext.tab.Panel', {
				itemId: "tabs",
				plain: true,
				border: false,
				bodyStyle: 'padding:5px 5px 0',
				border: 0,
				deferredRender: false,
				defaults: {
					border: false,
					autoScroll: true
				},
				items: [{
					title: 'Event',
					itemId: 'tabEvent',
					items:[{
						fieldLabel: _('Event title'),
						xtype: 'textfield',
						name: 'event_title',
						itemId: 'event_title',
						emptyText: _('Type here the event title'),
						allowBlank: false
					},{
						xtype: 'combobox',
						fieldLabel: _('Event source'),
						store: this.sources,
						itemId: 'event_source',
						name: 'event_source',
						editable: false,
						queryMode: 'local',
						displayField: 'component',
						valueField: 'component',
						emptyText: _('Type here the event source'),
						allowBlank: false
					},{
						xtype: "checkbox",
						fieldLabel: "All day event",
						itemId: 'all_day_event',
						inputValue: true,
						uncheckedValue: false,
						checked: true,
						listeners : {
								change: function(check, e, eOpts) {
									check.nextSibling().getComponent("start_time").setDisabled(check.getValue());
									check.nextSibling().nextSibling().getComponent("end_time").setDisabled(check.getValue());
								}
							}
					},{
						xtype: "cfieldset",
						title: _('Start'),
						margin: 0,
						padding: 0,
						items: [{
							fieldLabel: _('Date'),
							xtype: 'datefield',
							name: 'start_date',
							itemId: 'start_date'
						},{
							fieldLabel: _('Time'),
							xtype: 'timefield',
							name: 'start_time',
							itemId: 'start_time',
							disabled: true
						}]
					},{
						xtype: "cfieldset",
						title: _('End'),
						margin: 0,
						padding: 0,
						items: [{
								fieldLabel: _('Date'),
								xtype: 'datefield',
								name: 'end_date',
								itemId: 'end_date'
							},{
								fieldLabel: _('Time'),
								xtype: 'timefield',
								name: 'end_time',
								itemId: 'end_time',
								disabled: true
							}]
						}]
				}, {
					title: 'Recurrence',
					itemId: 'tabRecurrence',
					items:[{
							xtype: 'combobox',
							itemId: "rrule_preset",
							fieldLabel: "Preset",
							queryMode: "local",
							valueField: "value",
							store: {
								xtype: "store",
								fields: ["value", "text"],
								data : [
									{"value": "FREQ=DAILY", "text": "Daily"},
									{"value": "FREQ=WEEKLY", "text": "Weekly"},
									{"value": "FREQ=MONTHLY", "text": "Monthly"},
									{"value": "FREQ=WEEKLY;COUNT=10", "text": "Weekly, for 10 occurences"},
									{"value": "FREQ=WEEKLY;COUNT=10;BYDAY=TU,TH", "text": "Every Tuesday and Thursday for 10 occurences"},
									{"value": "FREQ=MONTHLY;COUNT=10;BYDAY=1FR", "text": "Monthly on the 1st Friday for ten occurrences"}
								]
							},
							listeners : {
								change: function(combo, e, eOpts) {
									combo.nextSibling().setValue(combo.getValue());

								}
							}
					},{
						xtype: 'textfield',
						itemId: "rrule",
						maskRe: /([0-9a-zA-Z,;=]+)$/,
						regex: /[0-9a-zA-Z,;=]/,
						fieldLabel: "Advanced rule"
					}]
				}]
		});

		this._form.add(tabs);
		this._form.bodyStyle = {
				background: '#ffffff'
			};
	},

	afterRender: function() {
		this.callParent(arguments);
	},

	showNewEvent : function(start, end, allDay){
		this.setTitle( _("New event"));

		this.currentEditedEvent = {};
		this.currentEditedEventHtml = null;
		this.addMode = true;

		this._form.down("#tabs").setActiveTab(0);

		var tabEvent = this._form.down("#tabEvent");
		var tabRecurrence = this._form.down("#tabRecurrence");
		tabRecurrence.down("#rrule").setValue("");

		tabEvent.down("#event_title").setValue("");

		tabEvent.down("#event_source").setValue("");
		tabEvent.down("#event_source").setDisabled(false);

		tabEvent.down("#start_date").setValue(start);
		tabEvent.down("#end_date").setValue(end);

		if(!allDay)
		{
			tabEvent.down("#start_time").setValue(start);
			tabEvent.down("#end_time").setValue(end);
			tabEvent.down("#all_day_event").setValue(false);
		}
		else
		{
			tabEvent.down("#start_time").setValue(null);
			tabEvent.down("#end_time").setValue(null);
			tabEvent.down("#all_day_event").setValue(true);
		}

		this.show();
	},

	showEditEvent: function(event, eventHtml){
		this.setTitle( _("Edit event"));

		this._form.down("#tabs").setActiveTab(0);

		var tabEvent = this._form.down("#tabEvent");

		var tabRecurrence = this._form.down("#tabRecurrence");

		this.currentEditedEvent = event;
		this.currentEditedEventHtml = eventHtml;
		this.addMode = false;

		tabEvent.down("#event_title").setValue(event.title);

		tabEvent.down("#event_source").setValue(event.component);
		tabEvent.down("#event_source").setDisabled(true);

		tabEvent.down("#start_date").setValue(new Date(event.start));

		if(event.end === null)
			event.end = event.start;

		tabEvent.down("#end_date").setValue(new Date(event.end));

		if(!event.allDay)
		{
			tabEvent.down("#start_time").setValue(event.start);
			tabEvent.down("#end_time").setValue(event.end);
			tabEvent.down("#all_day_event").setValue(false);
		}
		else
		{
			tabEvent.down("#start_time").setValue(null);
			tabEvent.down("#end_time").setValue(null);
			tabEvent.down("#all_day_event").setValue(true);
		}

		if(event.rrule !== null && event.rrule !== undefined)
			tabRecurrence.down("#rrule").setValue(event.rrule);
		else
			tabRecurrence.down("#rrule").setValue("");

		this.show();
	},

	ok_button_function: function(){
		var newEvent = {}; //TODO set this as a property to save hidden props

		var tabEvent = this._form.down("#tabEvent");
		var tabRecurrence = this._form.down("#tabRecurrence");

		newEvent.title = tabEvent.down("#event_title").getValue();

		if(this.currentEditedEvent.id)
			newEvent.id = this.currentEditedEvent.id;
		else
		{
			var now = new Date().toString('yyyy-MM-dd-hh-mm-ss');
			newEvent.id = newEvent.title + "-" + now + "@" + "widget-calendar";
		}

		combine = function(me, date, time, all_day) {
				if(Ext.isString(date)) date = new Date(date);
				if(!date) date = new Date();
				if(Ext.isString(time)) time = new Date(time);
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

		var startDateWidget = tabEvent.down("#start_date");
		var startTimeWidget = tabEvent.down("#start_time");
		var endDateWidget = tabEvent.down("#end_date");
		var endTimeWidget = tabEvent.down("#end_time");

		var start_datetime = combine(this, startDateWidget.getValue(), startTimeWidget.getValue());
		var end_datetime = combine(this, endDateWidget.getValue(), endTimeWidget.getValue());
		var all_day = startTimeWidget.getValue() === null || endTimeWidget.getValue() === null;

		newEvent.component = tabEvent.down("#event_source").getValue();

		newEvent.start = start_datetime;
		newEvent.end = end_datetime;

		var rrule = tabRecurrence.down("#rrule").getValue();

		if(rrule !== "")
			newEvent.rrule = rrule;

		var isHourfilledOnlyOnStartOrEnd = (startTimeWidget.getValue() === null && endTimeWidget.getValue() !== null)
											|| (startTimeWidget.getValue() !== null && endTimeWidget.getValue() === null)

		if(isHourfilledOnlyOnStartOrEnd)
			global.notify.notify(_('Form problem'), _('"You must fill start and end times, or nothing for an all day event"'), 'info');
		else if(newEvent.component === undefined || newEvent.component === null || newEvent.component === "")
			global.notify.notify(_('Form problem'), _('"You must specify a source for the event"'), 'info');
		else
		{
			newEvent.allDay = all_day;
			//add the new event to the calendar, by sending it to amqp
			this.calendar.send_events([newEvent]);

			this.hide();
		}
	},

	hide: function() {
		this.calendar.resetEventStyle(this.currentEditedEventHtml, this.currentEditedEvent);
		this.callParent();
	}
});