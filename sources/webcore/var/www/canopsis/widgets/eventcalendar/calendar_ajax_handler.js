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

Ext.define('widgets.eventcalendar.calendar_ajax_handler' , {
	alias: 'widget.calendar_ajax_handler',

	load: function(calendar) {
		this.calendar = calendar;
	},

	computeStackedUrl: function(start, end){
		if(!!this.calendar.stacked_events_filter)
		{
			var url = "/rest/events?_dc=1383151536066&limit=2000";

			var filter = this.computeTagsFilter(start, end);
			url += "&filter=";
			url += JSON.stringify(filter);

			return encodeURI(url);
		}
		return null;
	},

	computeTagsFilter: function(from, to) {
		if(!!this.calendar.stacked_events_filter)
		{
			return JSON.parse(this.calendar.stacked_events_filter);
		}
	},

	computeIcsUrl: function(start, end){
		if(this.calendar.sources)
		{
			var urls = [];
			for (var i = this.calendar.sources.length - 1; i >= 0; i--) {
				var url = "/cal/" + this.calendar.sources[i].component + "/" + start + "/" + end;
				url = encodeURI(url);
				urls.push(url);
			};

			return urls;
		}
		else
			return null;
	},

	getStackedEvents: function(start, end, callback)
	{
		var calendarRoot = this.calendar;
		var ajaxHandler = this;

		var result = [];

		var url = this.computeStackedUrl(start, end);
		if(url === null)
		{
			sendEventsToFullCalendar = function(calEvents){
				result = result.concat(calEvents);
				callback(result);
			}

			ajaxHandler.getCalendarEvents(start, end, sendEventsToFullCalendar);
		}
		else
		{
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

					ajaxHandler.getCalendarEvents(start,end, sendEventsToFullCalendar);
				}
			});
		}
	},

	getCalendarEvents: function(start, end, callback, currentSource, calEventsStack)
	{
		var calendarRoot = this.calendar;
		var ajaxHandler = this;

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
				ajaxHandler.getCalendarEvents(start, end, callback, currentSource + 1, calEventsStack);
			}
		});
		};
	},
});