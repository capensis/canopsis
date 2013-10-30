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
Ext.define('widgets.eventcalendar.eventcalendar' , {
	extend: 'canopsis.lib.view.cwebsocketWidget',

	alias: 'widget.eventcalendar',

	logAuthor: '[eventcalendar]',

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

	},

	afterRender: function() {
		var calendarRoot = this;
		$('#' + calendarRoot.id).fullCalendar({
			height: calendarRoot.height,
			events: calendarRoot.getEvents(),
			header: {
				left: 'prev,next today',
				center: 'title',
				right: 'month,agendaWeek,agendaDay'
			},
			selectable: true,
			selectHelper: true,
			select: function(start, end, allDay) {
				var title = prompt('Event Title:');
				if (title) {
					$('#'+ calendarRoot.id).fullCalendar('renderEvent',
						{
							title: title,
							start: start,
							end: end,
							allDay: allDay
						},
						true // make the event "stick"
					);
				}
				$('#'+ calendarRoot.id).fullCalendar('unselect');
			}
		});

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
		if (events.length >= this.max)
			this.wcontainer.removeAll(true);

		this.wcontainer.insert(0, events);

		//Remove last components
		while (this.wcontainer.items.length > this.max) {
			var item = this.wcontainer.getComponent(this.wcontainer.items.length - 1);
			this.wcontainer.remove(item.id, true);
		}
	}

});

