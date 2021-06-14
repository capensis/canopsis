//need:app/lib/store/cstore.js,app/model/Perfdata.js
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
Ext.define('canopsis.store.Perfdatas', {
	extend: 'canopsis.lib.store.cstore',
	model: 'canopsis.model.Perfdata',

	storeId: 'store.Perfdatas',

	logAuthor: '[store][Perfdatas]',

	autoLoad: false,
	autoSync: true,

	proxy: {
		type: 'rest',
		url: '/perfstore',
		appendId: false,
		batchActions: true,
		extraParams: {show_internals: false},
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
	},

	toggle_internal_metric: function() {
		log.debug('Toggle internal metric view/hide', this.logAuthor);

		if(this.proxy.extraParams.show_internals) {
			this.proxy.extraParams.show_internals = false;
			this.load();
		}
		else {
			this.proxy.extraParams.show_internals = true;
		}
	}
});
