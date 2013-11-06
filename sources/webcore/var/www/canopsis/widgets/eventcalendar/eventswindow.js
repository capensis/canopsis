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

Ext.define('widgets.eventcalendar.eventswindow' , {
	extend: 'canopsis.lib.view.cpopup',

	alias: 'widget.eventcalendar.eventswindow',

	height: 550,
	width: 800,
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
		this._form.bodyPadding = 0;
		this.grid = Ext.create('canopsis.lib.view.cgrid_state', {
			exportMode: this.exportMode,
			opt_paging: this.paging,
			filter: this.filter,
			pageSize: this.pageSize,
			remoteSort: true,
			height:490,
			opt_bar_bottom:true,
			opt_paging:true
		});
		this._form.add(this.grid);
	},

	afterRender: function() {
		this.callParent(arguments);
	},

	showEvents : function(calEvent, tags){
		console.log("showEvents");
		console.log(calEvent.start / 1000);
		console.log(tags);
		var d = calEvent.start;

		d.setHours(0);
		d.setMinutes(0);
		d.setSeconds(0);

		var startOfDayTimestamp = d / 1000;

		//increment by one day
		d = new Date(d.getTime() + (24 * 60 * 60 * 1000));

		var endOfDayTimestamp = d / 1000;

		var queryFilter = {
						"$and": [
							{ "timestamp": { "$gt": startOfDayTimestamp } },
							{ "timestamp": { "$lt": endOfDayTimestamp } }
						]
		};

		if(tags && tags != "")
		{
			console.log("tags found");
			tagQueryPart = {"tags" : tags}
			queryFilter["$and"].push(tagQueryPart);
		}
		else
			console.log("no tags found");


		console.log("showEvents::loadstore");
		this.grid.store.setFilter(queryFilter);
		console.log(queryFilter);

		this.grid.store.load();

		console.log("showEvents::end");

		this.show();
	}
});