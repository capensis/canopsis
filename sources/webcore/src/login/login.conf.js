/*
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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

require.config({
    waitSeconds: 40,
    baseUrl: '/static/',
    paths: {
        'app': 'canopsis/core',
        'jquery': 'canopsis/core/lib/wrappers/jquery',
        'bootstrap': 'canopsis/uibase/lib/externals/bootstrap/dist/js/bootstrap.min',
        'handlebars': 'canopsis/core/lib/externals/handlebars/handlebars',
        'ember-template-compiler': 'canopsis/core/lib/externals/min/ember-template-compiler',
        'ember-lib': 'canopsis/core/lib/externals/min/ember.debug',
        'ember-data-lib': 'canopsis/core/lib/externals/min/ember-data'
    },
    shim: {
        'bootstrap': {
            deps: ['jquery']
        },
        'ember-lib': {
            deps: ['jquery', 'ember-template-compiler', 'handlebars']
        },
        'ember-data-lib': {
            deps: ['ember-lib']
        }
    }
});

define([
    'canopsis/canopsisConfiguration',
    'bootstrap'
], function(canopsisConfiguration) {
    console.log(canopsisConfiguration);
    //Set page title
    var title = canopsisConfiguration.TITLE;
    if (title !== undefined) {
        $('title').html(title);
    }
});
