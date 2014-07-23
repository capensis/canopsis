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

	//TODO check if it used or not
	Ember.Handlebars.helper('json2html', function(json) {
		if (typeof json === 'string') {

			try {
				json = JSON.parse(json);
			} catch (exception) {
				console.error('unable to deserialize json');
				json = {'invalid': 'json string'};
			}
		}

		function parseJson(object) {
			var html = '';
			if (Ember.isArray(object)) {
				html += '<ul class="jsonUl">';
				for (var element in  object) {
					html += '<li class="jsonLi">' + parseJson(object[element]) + '</li>';
				}
				html += '</ul>';
			} else if (typeof object === 'object') {
				for (var key in object) {
					html += '<ul class="jsonUl"><li class="jsonLi">';
					html += '<span class="label label-primary">'+ key +'</span>';
					if (typeof object[key] === 'object') {
						html += parseJson(object[key]);
					} else {
						html +=  '&nbsp;<span class="glyphicon glyphicon-arrow-right" style="display:inline"></span><span class="label label-warning">'+ object[key] +'</span>' ;
					}
					html += '</li></ul>';
				}
			}

			return html;
		}

		var html = parseJson(json);

		return new Ember.Handlebars.SafeString(html);
	});

});
