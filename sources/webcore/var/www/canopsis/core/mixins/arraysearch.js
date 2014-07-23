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

	/**
	 * Implements search in arraycontrollers
	 *
	 * You should define on the ArrayController:
	 *	  - the `findOptions` property
	 *	  - the `refreshContent()` method
	 *	  - the searchableAttributes property
	 *
	 * @mixin
	 */
	Application.ArraySearchMixin = Ember.Mixin.create({
		actions: {
			searchItems: function(searchPhrase) {
				//TODO these checks should be asserts
				if (this.searchableAttributes === undefined) {
					console.warn("searchableAttributes not defined in controller, but searchItems still called. Doing nothing.", this);
					return;
				}
				if (typeof this.searchableAttributes !== "object") {
					console.warn("searchableAttributes should be an array.", this);
					return;
				}
				if (this.searchableAttributes.length === 0) {
					console.warn("Asking for a search on records with no searchableAttributes. Doing nothing.", this);
					return;
				}

				console.log("search", searchPhrase, this.get("attributesKeys"));

				if (this.findOptions === undefined) {
					this.findOptions = {};
				}

				var filter_orArray = [];
				for (var i = 0; i < this.searchableAttributes.length; i++) {
					var filter_orArrayItem = {};
					filter_orArrayItem[this.searchableAttributes[i]] = {"$regex": searchPhrase,"$options":"i"};
					filter_orArray.push(filter_orArrayItem);
				}

				this.findOptions.filter = JSON.stringify({"$and":[{"$or": filter_orArray }]});

				if (this.currentPage !== undefined) {
					this.set("currentPage", 1);
				}

				this.refreshContent();
			}
		}
	});

	return Application.ArraySearchMixin;
});
