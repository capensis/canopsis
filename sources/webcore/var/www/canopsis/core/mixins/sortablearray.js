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
	  Implements sorting in arraycontrollers

	  You should define on the ArrayController:
		  - the `findOptions` property
		  - the `refreshContent()` method

	*/
	Application.SortableArrayMixin = Ember.Mixin.create({

		sort_direction: false,

		actions: {
			sort: function(attribute) {
				var direction;

				direction = get(this, 'sort_direction') ? 'ASC' : 'DESC';
				set(this, 'sort_direction', !get(this, 'sort_direction'));

				if (get(this, 'sortedAttribute') !== undefined) {
					set(this, 'sortedAttribute.headerClassName', "sorting");
				}

				console.log('attribute', attribute);
				set(attribute, 'headerClassName', 'sorting_' + direction.toLowerCase());

				set(this, 'sortedAttribute', attribute);

				console.log("sortBy", arguments);
				if (this.findOptions === undefined) {
					this.findOptions = {};
				}

				this.findOptions.sort = JSON.stringify([{"property": attribute.field,"direction": direction}]);

				this.refreshContent();
			}
		},

		attributesKeys: function() {
			console.log("attributesKeys from sortableArray");
			var keys = this._super.apply(this);
			var sortedAttribute = get(this, 'sortedAttribute');

			console.log("sortedAttribute", sortedAttribute);
			if(sortedAttribute !== undefined)
			{
				for (var i = 0; i < keys.length; i++) {
					var currentKey = keys[i];
					var sortedAttributeField = get(sortedAttribute, 'field');
					var sortedAttributeHeaderClassName = get(sortedAttribute, 'headerClassName');

					if(get(currentKey, 'field') === sortedAttributeField) {
						set(currentKey, 'headerClassName', sortedAttributeHeaderClassName);
					} else {
						set(currentKey, 'headerClassName', 'sorting');
					}
				}
			} else {
				for (var i = 0; i < keys.length; i++) {
					set(keys[i], 'headerClassName', 'sorting');
				}
			}
			return keys;
		}.property("inspectedProperty", "inspectedDataArray"),
	});

	return Application.SortableArrayMixin;
});
