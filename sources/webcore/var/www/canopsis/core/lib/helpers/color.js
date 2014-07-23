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

define(['ember', 'app/application'], function(Ember, Application) {
	Ember.Handlebars.helper('menu', Application.MenuView);
	Ember.Handlebars.helper('color', function(color) {

		var style = '';

		if (color.toLowerCase() !== '#null') {
			style = 'style="background-color:' + color + '"';
		} else {
			color = 'no color';
		}

		return new Ember.Handlebars.SafeString('<div class="color" ' + style + '>' + color + '</div>');

	});

});
