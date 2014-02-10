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
Ext.define('canopsis.lib.store.cstore', {
	extend: 'Ext.data.Store',

	pageSize: global.pageSize,
	remoteSort: true,

	logAuthor: '[cstore]',

	loaded: false,

	listeners: {
		update: function() {
			if(this.storeId !== 'Tabs' && global.websocketCtrl) {
				global.websocketCtrl.publish_event('store', this.storeId, 'update');
			}
		},
		remove: function() {
			if(this.storeId !== 'Tabs' && global.websocketCtrl) {
				global.websocketCtrl.publish_event('store', this.storeId, 'remove');
			}
		},
		write: function(store, operation) {
			void(store);

			this.displaySuccess(operation);
		}
	},

	displaySuccess: function(operation) {
		if(operation.success) {
			if(operation.action === 'create') {
				global.notify.notify(_('Success'), _('Record saved'), 'success');
			}
			else if(operation.action === 'destroy') {
				global.notify.notify(_('Success'), _('Record deleted'), 'success');
			}
			else if(operation.action === 'update') {
				global.notify.notify(_('Success'), _('Record updated'), 'success');
			}
		}
	},

	constructor: function(config) {
		//now this.filter_list is local
		this.filter_list = {};
		this.baseFilter = false;

		this.callParent([config]);

		this.on('load', function() {
			this.loaded = true;
		}, this, {single: true});
	},


	//function for search and filter
	setFilter: function(filter) {
		log.debug('Setting base store filter', this.logAuthor);

		if(typeof(filter) === 'object') {
			this.baseFilter = filter;
		}
		else {
			this.baseFilter = Ext.JSON.decode(filter);
		}
	},

	addFilter: function(filter) {
		var md5 = $.md5(Ext.encode(filter));
		this.filter_list[md5] = filter;
		return md5;
	},

	deleteFilter: function(filter_id) {
		if(this.filter_list[filter_id]) {
			delete this.filter_list[filter_id];
		}
	},

	clearFilter: function() {
		if(this.baseFilter) {
			this.proxy.extraParams.filter = Ext.JSON.encode(this.baseFilter);
		}
		else {
			delete this.proxy.extraParams['filter'];
		}
	},

	getFilter: function() {
		return this.proxy.extraParams.filter;
	},

	getOrFilter: function(filter) {
		if (filter.length === 1) {
			return filter[0];
		}
		return {'$or': filter};
	},

	getAndFilter: function(filter) {
		if (filter.length === 1) {
			return filter[0];
		}
		return {'$and': filter};
	},

	getInFilter: function(filter) {
		if(Ext.isArray(filter)) {
			return {'$in': filter};
		}
		else {
			return {'$in': [filter]};
		}
	},

	getNinFilter: function(filter) {
		if(Ext.isArray(filter)) {
			return {'$nin': filter};
		}
		else {
			return {'$nin': [filter]};
		}
	},

	search: function(filter, autoLoad) {
		if(autoLoad === undefined) {
			autoLoad = true;
		}

		if(this.search_filter_id) {
			this.deleteFilter(this.search_filter_id);
		}

		if(filter !== undefined && filter !== '') {
			this.search_filter_id = this.addFilter(filter);
		}

		if (autoLoad) {
			this.load();
		}
	},

	getFilterList: function() {
		var filters = [];
		var filter_list = this.filter_list;

		Ext.Object.each(filter_list, function(key) {
			Ext.Object.each(filter_list[key], function(key, value) {
				var filter = {};
				filter[key] = value;
				filters.push(filter);
			});
		});

		if(this.baseFilter) {
			filters.push(this.baseFilter);
		}

		return Ext.JSON.encode(this.getAndFilter(filters));
	},

	load: function() {
		if(Ext.Object.getSize(this.filter_list) !== 0 || this.baseFilter) {
			this.proxy.extraParams.filter = this.getFilterList();
		}
		else {
			if(!this.proxy.extraParams) {
				this.proxy.extraParams = {};
			}

			if(this.proxy.extraParams.filter) {
				delete this.proxy.extraParams.filter;
			}
		}

		this.callParent(arguments);
	}
});
