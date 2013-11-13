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

Ext.require('widgets.eventcalendar.editwindow');

Ext.define('widgets.eventcalendar.eventcalendar' , {
	extend: 'canopsis.lib.view.cwebsocketWidget',

	alias: 'widget.eventcalendar',

	tags: '',

	limitEventNumber: 3,

	header_left: 'prev,next today',
	header_center: 'title',
	header_right: 'month,agendaWeek,agendaDay',
	defaultView: 'month',

    /**
     * @see cwebsocketWidget
     */
	amqp_queue: 'alerts',

	initComponent: function() {
		this.callParent(arguments);
		this.logAuthor = '[eventcalendar]';

		this.editwindow = Ext.create("widgets.eventcalendar.editwindow",
			{
				calendar: this
			});

		this.eventswindow = Ext.create("widgets.eventcalendar.eventswindow",
			{
				calendar: this
			});
	},

	afterContainerRender: function() {
		var calendarRoot = this;

		var eventSources = []

		var tags_url = this.computeTagsUrl()
		var ics_url = this.computeIcsUrl()
		if(tags_url)
			eventSources.push(tags_url);

		if(ics_url)
			eventSources.push(ics_url);

		$('#' + calendarRoot.id).fullCalendar({
			height: calendarRoot.height,
			events: calendarRoot.getEvents(),
			eventSources: eventSources,
			defaultView: this.defaultView,
			header: {
				left: this.header_left,
				center: this.header_center,
				right: this.header_right
			},
			selectable: calendarRoot.editable,
			selectHelper: calendarRoot.editable,
			editable: calendarRoot.editable,
			select: function(start, end, allDay) {
				calendarRoot.editwindow.showNewEvent(start, end, allDay);
			},
			eventDrop: function(calEvent, dayDelta, minuteDelta, allDay) {
				calEvent.allDay = allDay;
				calEvent.type = "calendar";
				calendarRoot.send_events([calEvent]);
			},
			eventResize: function(calEvent, dayDelta, minuteDelta, revertFunc, jsEvent) {
				calEvent.type = "calendar";
				calendarRoot.send_events([calEvent]);
			},
			eventClick: function(calEvent, jsEvent, view) {
				if(calEvent.type && calEvent.type === "non-calendar")
				{
					calendarRoot.eventswindow.showEvents(calEvent, calendarRoot.tags);
				}
				else
				{
					if(calendarRoot.editable)
					{
						$(this).css('border-color', 'red');
						calendarRoot.editwindow.showEditEvent(calEvent, this);
					}
				}

			}
		});

		this.subscribe();
		this.callParent(arguments);
	},


	onResize: function() {
		$('#'+ this.id).fullCalendar('option', 'height', this.height);
	},

	getEvents: function() {
		var date = new Date();
		var d = date.getDate();
		var m = date.getMonth();
		var y = date.getFullYear();

		return [];
	},

	send_events: function(calevents) {
		var calendarRoot = this;

		for (var i = calevents.length - 1; i >= 0; i--) {
			var isoStartDatetime = calevents[i].start.toISOString();
			var isoEndDatetime;
			//if end datetime is not define, define it to start + 2h
			if(calevents[i].end)
				isoEndDatetime = calevents[i].end.toISOString();
			else
			{
				var endDatetime = new Date(calevents[i].start);
				endDatetime.setHours(endDatetime.getHours() + 2);
				isoEndDatetime = endDatetime.toISOString();
			}

			var event_raw = {
				'connector': 'ics',
				'connector_name': 'ics2amqp',
				'source_type': 'resource',
				'event_type': 'calendar',
				'component': calevents[i].component,
				'resource': calevents[i].id,
				'output': calevents[i].title,
				'state': 0,
				'start': isoStartDatetime,
				'end': isoEndDatetime,
				'all_day': calevents[i].allDay
			};

			this.publishEvent('events', event_raw, false);
		};
	},

	computeTagsFilter: function(from, to) {
		var query = {
					"$and": [
						{ "timestamp": { "$gt": from } },
						{ "timestamp": { "$lt": to } }
					]
		};

		var tagQueryPart = {"$or":[]};
		var tags = this.tags.split(",");

		for (var i = tags.length - 1; i >= 0; i--) {
			tagQueryPart["$or"].push({ "tags" : tags[i] });
		};

		query["$and"].push(tagQueryPart);
		return query;
	},

	computeTagsUrl: function(from, to){
		//TODO limit should be dynamic
		//TODO manage several tags
		if(this.tags && this.tags != "")
		{
			var url = "/rest/events/event?_dc=1383151536066&limit=2000&filter=";

			var filter = this.computeTagsFilter("!start!", "!end!");

			url += JSON.stringify(filter);

			return encodeURI(url);
		}
		return null;
	},

	computeIcsUrl: function(from, to){

		this.computeIcsSources();

		//TODO limit should be dynamic
		//TODO manage several sources
		if(this.ics_sources && this.ics_sources != "")
		{
			var url = "/rest/events/event?_dc=1383151536066&limit=2000&filter=";

			var query = {
						"$and": [
							{ "timestamp": { "$gt": "!start!" } },
							{ "timestamp": { "$lt": "!end!" } }
						]
			};

			eventTypeQueryPart = {"event_type" : "calendar"}
			icsSourcesQueryPart = {"component" : {"$in": this.ics_sources_array}}
			query["$and"].push(eventTypeQueryPart);
			query["$and"].push(icsSourcesQueryPart);

			url += JSON.stringify(query);

			return encodeURI(url);
		}
		else
			return null;
	},

	computeIcsSources: function() {
		//TODO remove left and right spaces
		//TODO call this only on sources updates
		if(this.ics_sources && typeof this.ics_sources === "string" && this.ics_sources != "")
			this.ics_sources_array = this.ics_sources.split(",");
		else
			this.ics_sources_array = [];
	},

	/**
     * @see cwebsocketWidget
     */
	on_event: function(raw, rk) {
		this.computeIcsSources();
		var me = this;

		function in_sources_array(raw_component) {
			for (var i = me.ics_sources_array.length - 1; i >= 0; i--) {
				if(me.ics_sources_array[i] === raw_component)
					return true;
			};
			return false;
		};

		if(raw.event_type === "calendar" && in_sources_array(raw.component))
		{
			$('#'+ this.id).fullCalendar( 'removeEvents', raw.resource);
			var event = {
					id: raw.resource,
					title: raw.output,
					start: new Date(raw.start),
					end: new Date(raw.end),
					allDay: raw.all_day,
					type:raw.event_type,
					component:raw.component
				};

			$('#'+ this.id).fullCalendar('renderEvent',
				event
			);
		}
	},

	resetEventStyle: function(eventHtml)
	{
		if(eventHtml)
			$(eventHtml).css('border-color', '');

		$('#'+ this.id).fullCalendar('unselect');
	}
});
