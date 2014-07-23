/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

define([
	'ember',
	'app/application'
], function(Ember, Application) {
	var get = Ember.get,
		set = Ember.set;

	/**
	  Implements pagination in ArrayControllers

	  You should define on the ArrayController:
		  - the `findOptions` property
		  - the `findItems()` method

	*/
	Application.PaginationMixin = Ember.Mixin.create({
		paginationMixinContent: function() {
			console.warn("paginationMixinContent should be defined on the concrete class");
		},
		paginationMixinFindOptions: function() {
			console.warn("paginationMixinFindOptions should be defined on the concrete class");
		},

		actions: {
			prevPage: function() {
				if (this.get("currentPage") > 1) {
					this.set("currentPage", get(this, "currentPage") - 1);
				}
			},
			nextPage: function() {
				if (get(this, "currentPage") < get(this, "totalPages")) {
					this.set("currentPage", get(this, "currentPage") + 1);
				}
			}
		},

		itemsTotal: 1,
		itemsPerPage: 5,
		currentPage: 1,
		totalPages: 1,
		paginationFirstItemIndex: 1,
		paginationLastItemIndex: 1,

		onCurrentPageChanges: function() {
			this.refreshContent();
		}.observes('currentPage'),

		refreshContent: function() {
			console.group("paginationMixin refreshContent");

			if (get(this, 'paginationMixinFindOptions') === undefined) {
				set(this, 'paginationMixinFindOptions', {});
			}

			var itemsPerPage = get(this, 'itemsPerPage');

			//HACK when widget is saved and the app is not refreshed, itemsPerPage is a string!
			if (typeof itemsPerPage === 'string') {
				itemsPerPage = parseInt(itemsPerPage, 10);
			}
			if (itemsPerPage === 0) {
				console.warn("itemsPerPage is 0 in widget", this);
				console.warn("assuming itemsPerPage is 5");
				itemsPerPage = 5;
			}
			if (typeof itemsPerPage !== 'number' || itemsPerPage % 1 !== 0) {
				itemsPerPage = 5;
			}

			var start = itemsPerPage * (this.currentPage - 1);
			console.log("start", start);

			set(this, 'paginationMixinFindOptions.start', start);
			set(this, 'paginationFirstItemIndex', start + 1);
			set(this, 'paginationLastItemIndex', start + itemsPerPage);
			set(this, 'paginationMixinFindOptions.limit', itemsPerPage);
			this._super.apply(this, arguments);

			console.groupEnd();
		},

		extractItems: function(queryResult) {
			get(this, 'paginationMixinContent');
			this.set('itemsTotal', get(this, 'widgetDataMetas').total);

			var itemsPerPage = get(this, 'itemsPerPage');
			if (itemsPerPage === 0) {
				console.warn("itemsPerPage is 0 in widget", this);
				console.warn("assuming itemsPerPage is 5");
				itemsPerPage = 5;
			}

			if (get(this, 'itemsTotal') === 0) {
				set(this, 'totalPages', 0);
			} else {
				set(this, 'totalPages', Math.ceil(get(this, 'itemsTotal') / itemsPerPage));
			}

			this._super(queryResult);
		}
	});

	return Application.PaginationMixin;
});
