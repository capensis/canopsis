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

	getEvents: function(start, end, callback)
	{
		var calendarRoot = this.calendar;
		var ajaxHandler = this;

		var result = [];

		var stackedUrl = this.computeStackedUrl(start, end);
		var icsUrls = this.computeIcsUrl(start, end);

		var promises = [];

		if(stackedUrl !== undefined && stackedUrl !== null)
		{
			var ajaxRequest = $.ajax({
				url: stackedUrl,
				dataType: 'json',
			});

			ajaxRequest.eventsType = "stacked";
			promises.push(ajaxRequest);
		}

		for (var i = icsUrls.length - 1; i >= 0; i--) {
			var ajaxRequest = $.ajax({
				url: icsUrls[i],
				dataType: 'json',
			});

			ajaxRequest.eventsType = "calendar";

			promises.push(ajaxRequest);
		}

		if(this.calendar.downtimes !== undefined && this.calendar.downtimes !== []) {
			var ajaxRequest = $.ajax({
				url: "/rest/events",
				dataType: 'json',
				data: this.getDowntimesParams()
			});

			ajaxRequest.eventsType = "downtimes";
			promises.push(ajaxRequest);
		}

		$.when.apply($, promises).then(function(schemas) {
			var result = []

			for (var i = arguments.length - 1; i >= 0; i--) {
				var request = arguments[i];

				if(request !== undefined)
				{
					var data = request[0].data;

					events = ajaxHandler.format_events(data, request[2].eventsType);
					result = result.concat(events);
				}
				else
				{
					console.warning("request undefined");
				}
			};

			callback(result);
		});

	},

	format_events: function(events, eventsType) {
		calEvents = [];

		if(eventsType === "calendar")
		{
			for(key in events)
			{
				var event = events[key];

				var startDate = new Date(event.start * 1000);
				var endDate = new Date(event.end * 1000);

				var newEvent = {};
				newEvent.start = startDate;
				newEvent.end = endDate;
				newEvent.title = event.output;
				newEvent.type = event.event_type;
				newEvent.component = event.component;
				newEvent.id = event.resource;
				newEvent.rrule = event.rrule;

				if(event.all_day && event.all_day === true)
					newEvent.allDay = true;
				else
					newEvent.allDay = false;

				calEvents.push(newEvent);
			}
			return calEvents;
		}

		if(eventsType === "stacked")
		{
			function groupBy(array, field) {
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
				newEvent.editable = false;
				newEvent.type = "stacked";
				calEvents.push(newEvent);
			}

			return calEvents;
		}

		if(eventsType === "downtimes")
		{
			var result = [];

			for(key in events)
			{
				var event = events[key];
				var newEvent = {};
				newEvent.title = event.output;
				newEvent.start = event.start;
				newEvent.end = event.end;
				newEvent.allDay = false;
				newEvent.editable = false;
				newEvent.type = "downtime";

				result.push(newEvent);
			}
			return result;
		}

		return events;
	},

	getDowntimesParams: function() {
		var filter = {};
		filter["$or"] = [];

		// for (var i = cwidgetObject.nodes.length - 1; i >= 0; i--) {
		for(key in this.calendar.downtimes) {
			var downtime = this.calendar.downtimes[key];
			var currentFilterPart = { "$and" : [
						{component : downtime.component},
						{resource  : downtime.resource},
						{event_type  : "downtime"}
						// {"$or": [{end: { "$gt" : from}},
						// 		{start: { "$lt" : to}}]}
						// end: { "$gt" : from}
					]};
			filter["$or"].push(currentFilterPart);
		}

		var post_param = {
			filter: JSON.stringify(filter, undefined, 8)
		};

		return post_param;
	}
});