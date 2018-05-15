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
    name: 'EventConsumer',
    after: ['MixinFactory'],
    initialize: function(container, application) {
        var MixinFactory = container.lookupFactory('factory:mixin');

        var get = Ember.get;

        /**
         * @mixin EventConsumerMixin
         * @augments Mixin
         * Provide Event fetching mecanism to widgets.
         */
        var mixin = MixinFactory('eventconsumer', {
            /**
             * @method fetchEvents
             * @memberof EventConsumerMixin
             * @param {array} events - Array of labelled events
             * @returns promise
             */
            fetchEvents: function(events) {
                var store = get(this, 'widgetDataStore'),
                    rks = [],
                    labelsByRk = {},
                    me = this;

                events.forEach(function(evt) {
                    var rk = get(evt, 'rk'),
                        label = get(evt, 'label');

                    rks.push(rk);
                    labelsByRk[rk] = label;
                });

                var promise = store.findQuery('event', {
                    filter: JSON.stringify({_id: {'$in': rks}})
                });

                return promise.then(function(result) {
                    me.onEvents(get(result, 'content'), labelsByRk);
                });
            },

            /**
             * @abstract
             * @method onEvents
             * @memberOf EventConsumerMixin
             * @param {array} events - Events fetched by ``fetchEvents()`` method
             * @param {object} labelsByRk - Labels for each events
             * Called by ``fetchEvents()`` method.
             */
            onEvents: function(events, labelsByRk) {
                void(events);
                void(labelsByRk);
            }
        });

        application.register('mixin:eventconsumer', mixin);
    }
});
