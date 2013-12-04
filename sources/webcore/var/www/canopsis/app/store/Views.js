//need:app/lib/store/cstore.js,app/model/View.js
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
Ext.define('canopsis.store.Views', {
	extend: 'canopsis.lib.store.cstore',
	model: 'canopsis.model.View',

	storeId: 'store.View',

	autoLoad: true,
	autoSync: true,

	listeners: {
		remove: function() {
			if(this.storeId !== 'Tabs' && global.websocketCtrl) {
				global.websocketCtrl.publish_event('store', this.storeId, 'remove');
			}
		},
		update: function() {
			if(this.storeId !== 'Tabs' && global.websocketCtrl) {
				global.websocketCtrl.publish_event('store', this.storeId, 'update');
			}
		}
	},

	proxy: {
		type: 'rest',
		url: '/rest/object/view',
		extraParams: {
			noInternal: true,
			limit: 0
		},
		reader: {
			type: 'json',
			root: 'data',
			totalProperty: 'total',
			successProperty: 'success'
		},
		writer: {
			type: 'json',
			writeAllFields: false
		}
	}
});

