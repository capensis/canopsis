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

(function(){

    var folderPath = '/static/canopsis/media/images/';

    var images = {
        nagios : folderPath + 'nagioslogo.png',
        shinken : folderPath + 'shinkenlogo.png',
        schneider : folderPath + 'schneiderlogo.png',
        collectd : folderPath + 'collectd.jpg',
        sikuli: folderPath + 'sikulilogo.png',
        cucumber: folderPath + 'cucumberlogo.png',
        watir: folderPath + 'watirlogo.png',
        jmeter: folderPath + 'jmeterlogo.jpg',
        centreon: folderPath + 'centreonlogo.png',
        Engine: folderPath + 'engine.png'

    };

    Ember.Handlebars.helper('logo', function(value) {

        var logoPath = images[value];

        if(logoPath !== undefined) {
            return '<img alt="Source" src="'+ logoPath + '"/>';
        } else {
            return value;
        }
    });

    Ember.Handlebars.helper('logofromstring', function(imageName) {
        var  logoPath = images[imageName];

        return logoPath;
    });

    Ember.Handlebars.helper('logofromstring2', function(imageName) {
        var  logoPath = images[imageName];

        return '<img alt="Source" src="'+ logoPath + '"/>';
    });

})();
