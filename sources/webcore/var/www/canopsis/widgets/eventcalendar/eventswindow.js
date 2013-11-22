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
	extend: 'Ext.window.Window',

	alias: 'widget.eventcalendar.eventswindow',

	height: 500,
	width: 800,
	layout: 'fit',

	closeAction: 'hide',

	addMode: true,
	hasFooter: false,

	initComponent: function() {
		this.calendar = this.initialConfig.calendar;

		this.setTitle( _("Event browsing"));
		var form = this.buildForm();
		this.items = [form];
		this.callParent(arguments);
	},

	buildForm: function() {
		this._form = Ext.create('Ext.form.Panel', {
			layout: 'anchor',
			bodyStyle: {
				background: '#ededed'
			},
			bodyPadding: 10,
			border: false
		});

		this._form.bodyPadding = 0;
		this.grid = Ext.create('canopsis.lib.view.cgrid_state', {
			exportMode: this.exportMode,
			opt_paging: this.paging,
			filter: this.filter,
			pageSize: this.pageSize,
			remoteSort: true,
			height:490,
			opt_bar_bottom:true,
			opt_paging:true,
			opt_show_component: this.calendar.opt_show_component,
			opt_show_resource: this.calendar.opt_show_resource,
			opt_show_state: this.calendar.opt_show_state,
			opt_show_state_type: this.calendar.opt_show_state_type,
			opt_show_source_type: this.calendar.opt_show_source_type,
			opt_show_last_check: this.calendar.opt_show_last_check,
			opt_show_output: this.calendar.opt_show_output,
			opt_show_tags: this.calendar.opt_show_tags
		});
		this._form.add(this.grid);

		return this._form;
	},

	afterRender: function() {
		this.callParent(arguments);
	},

	showEvents : function(calEvent, tags){
		var d = calEvent.start;

		d.setHours(0);
		d.setMinutes(0);
		d.setSeconds(0);

		var startOfDayTimestamp = d / 1000;

		//increment by one day
		d = new Date(d.getTime() + (24 * 60 * 60 * 1000));

		var endOfDayTimestamp = d / 1000;

		var filter = this.calendar.computeTagsFilter(startOfDayTimestamp, endOfDayTimestamp);

		this.grid.store.setFilter(filter);
		this.grid.store.load();

		this.show();
	}
});