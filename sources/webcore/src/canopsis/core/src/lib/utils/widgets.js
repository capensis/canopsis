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
    name: 'WidgetsUtils',
    after: ['UtilityClass'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');

        var get = Ember.get;

        var widgetsUtils = Utility.create({

            name: 'widgetsUtils',

            getParentViewForWidget: function(widget) {
                var currentItem = widget;

                while (get(currentItem, 'crecord_type') !== 'view') {
                    currentItem = get(currentItem, 'target');
                }

                return currentItem;
            }
        });

        application.register('utility:widgets', widgetsUtils);
    }
});
