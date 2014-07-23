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
var DS;

define([
	'ember',
	'app/lib/factories/wrapper',
	'webcore-libs/dev/ember-data'
], function(Ember, Wrapper) {

	DS.ArrayTransform = DS.Transform.extend({
		deserialize: function(serialized) {
			if (Ember.typeOf(serialized) === 'array') {
				return serialized;
			}

			return [];
		},

		serialize: function(deserialized) {
			var type = Ember.typeOf(deserialized);

			if (type === 'array') {
				return deserialized;
			}
			else if (type === 'string') {
				return deserialized.split(',').map(function(item) {
					return jQuery.trim(item);
				});
			}

			return [];
		}
	});

	DS.IntegerTransform = DS.Transform.extend({
		deserialize: function(serialized) {
			console.log("deserialize integer: ", typeof serialized);

			if (typeof serialized === "number") {
				return serialized;
			} else {
				console.warn("deserialized value is not a number as it is supposed to be", arguments);
				return 0;
			}
		},

		serialize: function(deserialized) {
			console.log("serialize : ",deserialized);

			return Ember.isEmpty(deserialized) ? null : Number(deserialized);

		}
	});

	DS.ObjectTransform = DS.Transform.extend({
		deserialize: function(serialized) {
			if (Ember.typeOf(serialized) === 'object') {
				return serialized;
			}

			return {};
		},

		serialize: function(deserialized) {
			var type = Ember.typeOf(deserialized);

			if (type === 'object') {
				return deserialized;

			} else if (type === 'string') {
				console.log("bad format");
			}

			return {};
		}
	});

	return Wrapper("ember-data", DS, arguments, DS.VERSION);
});