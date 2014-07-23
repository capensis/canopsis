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

	Ember.Handlebars.helper('duration', function(timestamp) {
		if (timestamp) {
			timestamp = parseInt(timestamp);

			var dt = {
				days: 0,
				hours: 0,
				minutes: 0,
				seconds: timestamp
			};

			dt.minutes = parseInt(dt.seconds / 60);
			dt.seconds -= (dt.minutes * 60);

			dt.hours = parseInt(dt.minutes / 60);
			dt.minutes -= (dt.hours * 60);

			dt.days = parseInt(dt.hours / 24);
			dt.hours -= (dt.days * 24);

			var str = "";

			if (dt.days > 0) {
				str += dt.days + ' day' + (dt.days > 1 ? 's ' : ' ');
			}

			if (dt.hours > 0) {
				str += dt.hours + ' h ';
			}

			if (dt.minutes > 0) {
				str += dt.minutes + ' min ';
			}

			if (dt.seconds > 0) {
				str += dt.seconds + ' s';
			}

			return new Ember.Handlebars.SafeString(str);
		}
		else {
			return new Ember.Handlebars.SafeString('');
		}
	});

});
