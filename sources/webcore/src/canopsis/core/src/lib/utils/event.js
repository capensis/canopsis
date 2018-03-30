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
    name: 'EventUtils',
    after: ['UtilityClass'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');

        var get = Ember.get,
            set = Ember.set,
            __ = Ember.String.loc,
            isNone = Ember.isNone;

        //TODO delete this, as it looks more like a registry than an util
        var eventUtil= Utility.create({

            name: 'event',

            getFields: function() {
                return [
                    'connector',
                    'connector_name',
                    'component',
                    'resource',
                    'perimeter',
                    'domain',
                    'state',
                    'status',
                    'timestamp',
                    'output'
                ];
            },

            /**
             * apply a filter to an element before to transmit it to an other view.
             * see goToInfo method of brick-calendar for example
             * @method applyVolatileFilter
             * @param transition: promise to transmit the element to an other view
             * @param element: concerned event
             * @param filter: a mongoDB style query as a string
             */
            applyVolatileFilter: function (transition, element, filter) {
                transition.promise.then(function(routeInfos){
                    var list = get(routeInfos, 'controller.content.containerwidget');
                    list = get(list, '_data.items')[0];
                    list = get(list, 'widget');
                    console.log(list);
                    if(filter !== '') {
                        if(!get(list, 'volatile')) {
                            set(list, 'volatile', {});
                        }

                        set(list, 'volatile.forced_filter', filter);
                        set(list, 'rollbackable', true);
                        set(list, 'title', element.title);
                    }
                });
            }
        });

        application.register('utility:event', eventUtil);
    }
});