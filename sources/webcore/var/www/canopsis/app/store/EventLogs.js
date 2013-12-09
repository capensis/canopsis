//need:app/lib/store/cstore.js,app/model/Event.js
/*
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
*/
Ext.define('canopsis.store.EventLogs', {
	extend: 'canopsis.lib.store.cstore',
	model: 'canopsis.model.Event',

	storeId: 'store.Derogations',

	logAuthor: '[store][curve]',

	autoLoad: false,
	autoSync: true,

	sortOnLoad: true,

	eventFilter: {},

	sorters: [
		{
			property: 'timestamp',
			direction: 'DESC'
		}
	],

	proxy: {
		type: 'rest',
		url: '/rest/events_log/event',
		batchActions: true,
		reader: {
			type: 'json',
			root: 'data',
			totalProperty: 'total',
			successProperty: 'success'
		}
	},

	toggleEventFilter: function(field, value) {
		if(this.eventFilter[field] === undefined) {
			this.eventFilter[field] = [value];
		}
		else {
			if(Ext.Array.contains(this.eventFilter[field], value)) {
				var index = this.eventFilter[field].indexOf(value);
				this.eventFilter[field].splice(index, 1);
			}
			else {
				this.eventFilter[field].push(value);
			}
		}

		this.buildEventFilter();
	},

	buildEventFilter: function() {
		var cleaned_filter = {};

		for(var i = 0; i < this.eventFilter.length; i++) {
			cleaned_filter[i] = this.getNinFilter(this.eventFilter[i]);
		}

		this.setFilter(cleaned_filter);
	}
});
