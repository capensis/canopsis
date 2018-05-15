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
    name: 'MonitoringApplicationRouteReopen',
    after: 'ApplicationRoute',
    initialize: function(container, application) {
        var ApplicationRoute = container.lookupFactory('route:application');

        var get = Ember.get,
            set = Ember.set;


        var applicationRoute = ApplicationRoute.reopen({
            /**
             * @method buildBeforeModelPromises
             * @param {Transition} transition
             * @return {Promise}
             *
             * Feed the ApplicationController with ticket options.
             */
            buildBeforeModelPromises: function() {
                var store = get(this, 'store');
                var ticketPromise = store.find('ticket', 'cservice.ticket');
                var appController = this.controllerFor('application');

                ticketPromise.then(function(queryResults) {
                    appController.ticketConfig = queryResults;
                });

                get(this, 'promiseArray').pushObject(ticketPromise);

                return this._super();
            }
        });

        application.register('route:application', applicationRoute);
    }
});
