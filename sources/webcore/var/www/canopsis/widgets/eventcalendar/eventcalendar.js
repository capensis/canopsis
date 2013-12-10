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

function CalendarException (message, event) {
	this.message = message;
	this.event = event;
	this.name = "CalendarException";
}

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

	event_display_size : "normal_size",

	opt_show_component : true,
	opt_show_resource : true,
	opt_show_state : true,
	opt_show_state_type : true,
	opt_show_source_type : true,
	opt_show_last_check : true,
	opt_show_output : true,
	opt_show_tags : true,

	weekends : true,

	defaultEventColor :"#3a87ad",
	/**
	 * @see cwebsocketWidget
	 */
	amqp_queue: 'alerts',
	sources: [],
	sources_byComponent:{},

	initComponent: function() {
		this.callParent(arguments);
		this.logAuthor = '[eventcalendar]';

		this.editwindow = Ext.create("widgets.eventcalendar.editwindow", {
			calendar: this
		});

		this.eventswindow = Ext.create("widgets.eventcalendar.eventswindow", {
			calendar: this
		});

		if(this.defaultView !== "month" && this.defaultView !== "agendaWeek" && this.defaultView !== "agendaDay" && this.defaultView !== "basicWeek" && this.defaultView !== "basicDay")
			this.defaultView = "month";
	},

	afterContainerRender: function() {
		var calendarRoot = this;

		var eventSources = [];

		var tags_url = this.computeStackedUrl("!start!", "!end!");
		var ics_url = this.computeIcsUrl("\"!start!\"", "\"!end!\"");
		if(tags_url)
			eventSources.push(tags_url);

		if(ics_url)
			eventSources = eventSources.concat(ics_url);

		$('#' + calendarRoot.wcontainer.id).fullCalendar({
			firstDay:1,
			height: calendarRoot.wcontainer.height,
			// eventSources: eventSources,
			events: function(start, end, callback){
				var events = [];
				var start_unixTimestamp = new Date(start).getTime() / 1000;
				var end_unixTimestamp = new Date(end).getTime() / 1000;
				// events = events.concat(calendarRoot.getCalendarEvents(start_unixTimestamp, end_unixTimestamp, callback));
				calendarRoot.getStackedEvents(start_unixTimestamp, end_unixTimestamp, callback);

				// callback(events);
			},
			defaultView: this.defaultView,
			weekends : this.show_weekends,
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
			},
			eventRender: function(event, element) {
				if(!event.component && event.type === "calendar")
					throw new CalendarException("Event of type calendar does not have a component property", event);

				for (var i = calendarRoot.sources.length - 1; i >= 0; i--) {
					var currentSource = calendarRoot.sources[i];
					calendarRoot.sources_byComponent[currentSource.component] = currentSource;
				};

				if(event.type === "non-calendar")
				{
					log.debug('no component (ics source) for event, assuming the event is stacked regular events', calendarRoot.logAuthor);
					element.css({"background-color" : calendarRoot.defaultEventColor, "border-color" : calendarRoot.defaultEventColor});
				}
				else
				{
					if(!calendarRoot.sources_byComponent[event.component])
					{
						log.debug('component not found on calendar sources', calendarRoot.logAuthor);
						return;
					}

					var sourceColor = calendarRoot.sources_byComponent[event.component].color;

					if(!!sourceColor)
					{
						element.css({ "background-color" : sourceColor, "border-color" : sourceColor});
					}
				}
				return true;
    		},

		});

		$('#' + calendarRoot.wcontainer.id).addClass(calendarRoot.event_display_size);

		this.subscribe();
		this.callParent(arguments);
	},

	getStackedEvents: function(start, end, callback)
	{
		var calendarRoot = this;

		var url = this.computeStackedUrl(start, end);

		function groupBy(array, field){
			var associativeArray = {};
			for (var i = 0; i < array.length; i++) {
				var item = array[i];
				var associativeArrayItem = associativeArray[item[field]];
				if(!associativeArray[item[field]])
				{
					associativeArray[item[field]] = {};
					associativeArray[item[field]]["items"] = [item];
					associativeArray[item[field]]["count"] = 1;
				}
				else
				{
					associativeArray[item[field]]["items"].push(item);
					associativeArray[item[field]]["count"] ++;
				}
			};

			return associativeArray;
		};

		$.ajax({
            url: url,
            dataType: 'json',
			success: function(request_result) {
				var events = request_result["data"] || [];
				var result = [];

				for (var i = 0; i < events.length; i++) {
					var d = new Date(events[i].timestamp * 1000);
					d.setHours(0);
					d.setMinutes(0);
					d.setSeconds(0);
					events[i].day =  d / 1000;
				};

				events = groupBy(events, "day");

				for(day in events)
				{
					var newEvent = {};
					newEvent.start = day;

					var evCount = events[day].count;
					var evCountText = evCount === 1 ? " event" : " events";
					newEvent.title = events[day].count.toString() + evCountText;
					newEvent.type = "non-calendar";
					newEvent.editable = false;
					result.push(newEvent);
				}

				sendEventsToFullCalendar = function(calEvents){
					result = result.concat(calEvents);
					callback(result);
				}

				calendarRoot.getCalendarEvents(start,end, sendEventsToFullCalendar);

			}
		});
	},

	getCalendarEvents: function(start, end, callback, currentSource, calEventsStack)
	{
		var calendarRoot = this;
		var urls = this.computeIcsUrl(start, end);

		//this happens when recursion begins
		if(currentSource === undefined)
			currentSource = 0;
		if(calEventsStack === undefined)
			calEventsStack = [];

		if(urls[currentSource] === undefined) {
			callback(calEventsStack);
		} else {
			var url = urls[currentSource];

			$.ajax({
            url: url,
            dataType: 'json',
			success: function(request_result) {
				var events = request_result["data"] || [];
				var result = [];

				for (var i = 0; i < events.length; i++) {
					var startDate = new Date(events[i].start * 1000);
					var endDate = new Date(events[i].end * 1000);

					var newEvent = {};
					newEvent.start = startDate;
					newEvent.end = endDate;
					newEvent.title = events[i].output;
					newEvent.type = events[i].event_type;
					newEvent.component = events[i].component;
					newEvent.id = events[i].resource;
					newEvent.rrule = events[i].rrule;

					if(events[i].all_day && events[i].all_day === true)
						newEvent.allDay = true;
					else
						newEvent.allDay = false;

					result.push(newEvent);
				}

				calEventsStack = calEventsStack.concat(result);
				calendarRoot.getCalendarEvents(start, end, callback, currentSource + 1, calEventsStack);
			}
		});
		};
	},

	onResize: function() {
		$('#'+ this.wcontainer.id).fullCalendar('option', 'height', this.getHeight());
	},

	send_events: function(calevents) {
		var calendarRoot = this;

		for (var i = calevents.length - 1; i >= 0; i--) {
			var tsStartDatetime = calevents[i].start.getTime();
			var tsEndDatetime;
			//if end datetime is not define, define it to start + 2h
			if(calevents[i].end)
				tsEndDatetime = calevents[i].end.getTime();
			else
			{
				var endDatetime = new Date(calevents[i].start);
				endDatetime.setHours(endDatetime.getHours() + 2);
				tsEndDatetime = endDatetime.getTime();
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
				'start': tsStartDatetime / 1000,
				'end': tsEndDatetime / 1000,
				'all_day': calevents[i].allDay
			};

			if(calevents[i].rrule !== null && calevents[i].rrule !== undefined)
				event_raw.rrule = calevents[i].rrule;

			this.publishEvent('events', event_raw, false);
		};
	},

	computeTagsFilter: function(from, to) {
		if(!!this.stacked_events_filter)
		{
			var query = {
						"$and": [
							{ "timestamp": { "$gt": from } },
							{ "timestamp": { "$lt": to } }
						]
			};

			query["$and"].push(JSON.parse(this.stacked_events_filter));
			return query;
		}
	},

	computeStackedUrl: function(start, end){
		if(!!this.stacked_events_filter)
		{
			var url = "/events?_dc=1383151536066&limit=2000";

			var filter = this.computeTagsFilter(start, end);
			url += "&filter=";
			url += JSON.stringify(filter);

			var perfdata_history = { "start": start, "end": end};
			url += "&perfdata_history=";
			url += JSON.stringify(perfdata_history);

			return encodeURI(url);
		}
		return null;
	},

	computeIcsUrl: function(start, end){

		//TODO limit should be dynamic
		if(this.sources)
		{
			var urls = [];
			for (var i = this.sources.length - 1; i >= 0; i--) {
				var url = "/cal/" + this.sources[i].component + "/" + start + "/" + end;
				url = encodeURI(url);
				urls.push(url);
			};

			return urls;
		}
		else
			return null;
	},

	/**
	 * @see cwebsocketWidget
	 */
	on_event: function(raw, rk) {
		var me = this;

		function in_sources_array(raw_component) {
			for (var i = me.sources.length - 1; i >= 0; i--) {
				if(me.sources[i].component === raw_component)
					return true;
			};
			return false;
		};

		if(raw.event_type === "calendar" && in_sources_array(raw.component))
		{
			$('#'+ this.wcontainer.id).fullCalendar( 'removeEvents', raw.resource);
			var event = {
					id: raw.resource,
					title: raw.output,
					start: new Date(raw.start * 1000),
					end: new Date(raw.end * 1000),
					allDay: raw.all_day,
					type: raw.event_type,
					component: raw.component,
					rrule: raw.rrule
				};

			$('#'+ this.wcontainer.id).fullCalendar('renderEvent',
				event
			);

			if(!!event.rrule)
			{
				log.debug("event has a rrule. Refetch all calendar events of the displayed period.", this.logAuthor);
				$('#'+ this.wcontainer.id).fullCalendar('refetchEvents');
			}
		}
	},

	resetEventStyle: function(eventHtml, event)
	{
		for (var i = this.sources.length - 1; i >= 0; i--) {
			var currentSource = this.sources[i];
			this.sources_byComponent[currentSource.component] = currentSource;
		};

		if(eventHtml)
			$(eventHtml).css("border-color", this.sources_byComponent[event.component].color);

		$('#'+ this.wcontainer.id).fullCalendar('unselect');
	}
});
