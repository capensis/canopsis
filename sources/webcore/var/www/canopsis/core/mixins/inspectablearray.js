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

	Application.InspectableArrayMixin = Ember.Mixin.create({
		attributesKeys: function() {
			var attributes = [];

			var attributesDict = this.get('inspectedDataArray.type.attributes.values');
			console.log("attributesDict", attributesDict);

			for (var key in attributesDict) {
				var attr = attributesDict[key];

				if (attr.options.hiddenInLists === false || attr.options.hiddenInLists === undefined) {
					attributes.push(Ember.Object.create({
						field: attr.name,
						type: attr.type,
						options: attr.options
					}));
					console.log("pushed attr", {
						field: attr.name,
						type: attr.type,
						options: attr.options
					});
				}
			}
			return attributes;
		}.property("inspectedProperty", "inspectedDataArray"),


		inspectedDataArray: function() { console.error("This must be defined on the base class"); }.property()
	});

	return Application.InspectableArrayMixin;
});
