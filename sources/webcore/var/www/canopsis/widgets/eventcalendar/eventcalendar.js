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
		$('#' + this.id).fullCalendar({
			height: this.height,
			events: this.getEvents()
		});

		this.initializeStore();
	},

	initializeStore: function(){
		this.store = Ext.create('canopsis.store.EventLogs', {
			model: 'canopsis.model.EventLogs',

			autoLoad: false,

			proxy: {
				type: 'rest',
				url: '/event/perftop',
				extraParams: {
					'limit': this.limit,
					'sort': this.sort,
					'mfilter': this.mfilter,
					'threshold': this.threshold,
					'threshold_direction': this.threshold_direction,
					'expand': this.expand,
					'percent': this.show_percent,
					'threshold_on_pct': this.threshold_on_pct,
					'report': this.reportMode || this.exportMode
				},
				reader: {
					type: 'json',
					root: 'data',
					totalProperty: 'total',
					successProperty: 'success'
				}
			}
		});
	},

	onResize: function() {
		$('#'+ this.id).fullCalendar('option', 'height', this.height);
	},

	getEvents: function() {
		var date = new Date();
		var d = date.getDate();
		var m = date.getMonth();
		var y = date.getFullYear();

		return [{
					title: 'All Day Event',
					start: new Date(y, m, 1)
				},
				{
					title: 'Long Event',
					start: new Date(y, m, d-5),
					end: new Date(y, m, d-2)
				},
				{
					id: 999,
					title: 'Repeating Event',
					start: new Date(y, m, d-3, 16, 0),
					allDay: false
				},
				{
					id: 999,
					title: 'Repeating Event',
					start: new Date(y, m, d+4, 16, 0),
					allDay: false
				},
				{
					title: 'Meeting',
					start: new Date(y, m, d, 10, 30),
					allDay: false
				},
				{
					title: 'Lunch',
					start: new Date(y, m, d, 12, 0),
					end: new Date(y, m, d, 14, 0),
					allDay: false
				},
				{
					title: 'Birthday Party',
					start: new Date(y, m, d+1, 19, 0),
					end: new Date(y, m, d+1, 22, 30),
					allDay: false
				},
				{
					title: 'Click for Google',
					start: new Date(y, m, 28),
					end: new Date(y, m, 29),
					url: 'http://google.com/'
				}];
	}
});

