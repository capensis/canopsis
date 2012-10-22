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
Ext.define('canopsis.controller.LiveSearch', {
	extend: 'Ext.app.Controller',

	views: ['LiveSearch.View', 'LiveSearch.Grid'],
	stores: ['Inventory'],

	logAuthor: '[controller][LiveSearch]',
	//models: [''],

	//iconCls: 'icon-crecord_type-account',

	refs: [{

		ref: 'grid',
		selector: 'LiveGrid'

	},{
		ref: 'liveSearch',
		selector: 'LiveSearch'
	}],

	init: function() {
		this.callParent(arguments);

		this.control({
			'LiveSearch #LiveSearchButton' : {
				click: this.addFilter
			}
		});

	},

	addFilter: function() {
		log.debug('Search button pushed', this.logAuthor);
		var store = this.getGrid().getStore();
		store.clearFilter();

		search = {};

		var searchValue = this.getLiveSearch().down('#source_name').value;
		if (searchValue) {
			store.load().filter('source_name', searchValue);
		}

		searchValue = this.getLiveSearch().down('#type').value;
		if (searchValue) {
			store.load().filter('type', searchValue);
		}

		searchValue = this.getLiveSearch().down('#source_type').value;
		if (searchValue) {
			store.load().filter('source_type', searchValue);
		}

		searchValue = this.getLiveSearch().down('#component').value;
		if (searchValue) {
			store.load().filter('component', searchValue);
		}

		store.load();
	}

});
