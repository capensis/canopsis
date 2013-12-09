//need:app/lib/store/ctreeStore.js,app/model/View.js
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

Ext.define('canopsis.store.TreeStoreViews', {
	extend: 'canopsis.lib.store.ctreeStore',
	model: 'canopsis.model.View',

	storeId: 'store.TreeStoreViews',

	autoLoad: false,
	autoSync: false,

	clearOnLoad: true,

	defaultRootId: 'directory.root',

	proxy: {
		batchActions: true,
		type: 'rest',
		url: '/ui/view',
		reader: {
			type: 'json'
		},
		writer: {
			type: 'json',
			method: 'POST'
		}
	},

	afterCorrectWrite : function() {
		Ext.getStore('Views').load();
	}
});
