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
	logAuthor: '[eventcalendar]',

	tags: '',
	editable: true,

	initComponent: function() {
		// var url = "widgets/eventcalendar/fullcalendar.js";
		// var onload = function () {
		// 	console.log("sucess");
		// }
		// var onerror = function () {
		// }
		// var scope = this;

		// Ext.Loader.injectScriptElement(url, onload, onerror, scope);

		// console.log("=========================");
		// console.log(this.id);

		this.callParent(arguments);

		this.editwindow = Ext.create("widgets.eventcalendar.editwindow",
			{
				calendar: this
			});
	},

	afterRender: function() {
		var calendarRoot = this;



		$('#' + calendarRoot.id).fullCalendar({
			height: calendarRoot.height,
			events: calendarRoot.getEvents(),
			eventSources: [
        		this.computeUrl()
        	],
			header: {
				left: 'prev,next today',
				center: 'title',
				right: 'month,agendaWeek,agendaDay'
			},
			selectable: calendarRoot.editable,
			selectHelper: calendarRoot.editable,
			editable: calendarRoot.editable,
			select: function(start, end, allDay) {
				console.log("select");
				console.log(start);
				console.log(end);
				calendarRoot.editwindow.showNewEvent(start, end, allDay);
			},
			// eventDrop: function() {
			// }
			eventClick: function(calEvent, jsEvent, view) {
				console.log("Click");
				console.log(calEvent);
		        $(this).css('border-color', 'red');
				calendarRoot.editwindow.showEditEvent(calEvent);

    		}
		}).limitEvents(3);

		calendarRoot.add_events([]);
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

	add_events: function(events) {
		var calendarRoot = this;

		for (var i = events.length - 1; i >= 0; i--) {
			$('#'+ calendarRoot.id).fullCalendar('renderEvent',
				{
					title: events[i].title,
					start: events[i].start,
					end: events[i].end,
					allDay: events[i].allDay
				},
				true // make the event "stick"
			);
		};
	},

	computeUrl: function(from, to){
		var url = "/rest/events/event?_dc=1383151536066&filter=";

		var query = {
						"$and": [
        					{ "timestamp": { "$gt": "!start!" } },
        					{ "timestamp": { "$lt": "!end!" } }
    					]
		};

		if(this.tags && this.tags != "")
		{
			console.log("tags found");
			tagQueryPart = {"tags" : this.tags}
			query["$and"].push(tagQueryPart);
		}
		else
			console.log("no tags found");


		url += JSON.stringify(query);

		return encodeURI(url);
	}
});
