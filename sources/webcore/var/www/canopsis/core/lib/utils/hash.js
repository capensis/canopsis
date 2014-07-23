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

define(['ember'], function(Ember) {

	var hash = {
		generate_GUID: function() {
			//Generates a random GUID
			var s4 = function () {
				return Math.floor((1 + Math.random()) * 0x10000).toString(16).substring(1);
			};

			var token = s4() + s4() + '-' + s4() + '-' + s4() + '-' + s4() + '-' + s4() + s4() + s4();

			return token;
		},
		generateId: function(prefix) {

			var token = hash.generate_GUID();

			if(!Ember.isNone(prefix)) {
				token = prefix + '_' + token;
			}
			return token;
		},
	};

	return hash;
});