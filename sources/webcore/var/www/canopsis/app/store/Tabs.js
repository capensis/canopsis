//need:app/lib/store/cstore.js
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
Ext.define('canopsis.store.Tabs.model', {
	extend: 'Ext.data.Model',
	fields: ['id', 'view_id', 'options', 'title', 'closable']
});

if(Ext.isIE) {
	Ext.define('canopsis.store.Tabs', {
		extend: 'canopsis.lib.store.cstore',
		model: 'canopsis.store.Tabs.model',
		id: 'Tabs',

		proxy: {
			type: 'memory',
			reader: {
				type: 'json'
			}
		},

		//HACK :we don't want this store to talk everytime
		listeners: {},

		autoLoad: false,
		autoSync: false
	});
}
else {
	Ext.define('canopsis.store.Tabs', {
		extend: 'canopsis.lib.store.cstore',
		model: 'canopsis.store.Tabs.model',
		id: 'Tabs',

		proxy: {
			type: 'localstorage',
			id: 'canopsis.store.tabs'
		},

		//HACK :we don't want this store to talk everytime
		listeners: {},

		autoLoad: false,
		autoSync: true
	});
}

