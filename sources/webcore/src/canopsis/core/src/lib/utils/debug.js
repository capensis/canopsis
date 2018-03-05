/**
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
    name: 'DebugUtils',
    after: ['UtilityClass'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');

        var set = Ember.set;

        var debugUtils = Utility.create({

            name: 'debug',

            inspectObject: function(object) {
                window.$E = object;

                set(this, 'inspectedObject', object);

                console.info('--- inspect object :', this.inspectedObject);
            },

            getViewFromJqueryElement: function($el, className) {
                if(className) {
                    return Ember.View.views[$el.closest('.ember-view .' + className).attr('id')];
                } else {
                    return Ember.View.views[$el.closest('.ember-view').attr('id')];
                }
            }
        });

        application.register('utility:debug', debugUtils);
    }
});
