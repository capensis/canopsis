/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'LinksDetailsHelper',
    initialize: function(container, application) {
        void(application);

        var get = Ember.get,
            isNone = Ember.isNone;


        Ember.Handlebars.helper('linksdetails', function(info, value) {

            var extensions = ['jpg', 'JPG', 'jpeg', 'JPEG', 'png', 'PNG'];

            var details = '<ul>';
            for (var prop in info) {
                console.error('prop :',prop);
                console.error('info[prop] :',info[prop]);
                var reg = new RegExp("^.*\.(" + extensions.join('|') + ")$")
                for (i=0; i < info[prop].length; i++){
                    console.error('toto', info[prop][i]);
                    if (info[prop][i].match(reg)) {
                        details = details + '<li><a href=' + info[prop][i] + ' target="_blank"><img src=' + info[prop][i] + ' alt=' + prop + ' class="image-center" style="max-width: 100px;max-height: 100px;"></a></li>';
                    } else {
                        details = details + '<li><a href=' + info[prop][i] + ' target="_blank">' + prop + '</a></li>';
                    }
                }
            }
            details = details + '</ul>'
            return new Ember.Handlebars.SafeString(details);
        });
    }
});
