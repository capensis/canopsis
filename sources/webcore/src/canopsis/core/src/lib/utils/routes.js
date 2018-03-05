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
    name: 'RoutesUtils',
    after: ['UtilityClass', 'DataUtils'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');
        var dataUtils = container.lookupFactory('utility:data');

        /**
         * @class RoutesUtils
         * @extends Utility
         */
        var routesUtils = Utility.create({
            name: 'routes',

            /**
             * @method getCurrentRouteController
             * @return Ember.Controller
             */
            getCurrentRouteController: function() {
                var currentHandlers = dataUtils.getEmberApplicationSingleton().__container__.lookup("router:main").router.currentHandlerInfos;
                var currentRouteController = currentHandlers[currentHandlers.length - 1].handler.controller;

                console.log("currentHandlers", currentHandlers);
                console.log("currentRouteController", currentRouteController);

                return currentRouteController;
            },

            /**
             * @method getCurrentViewId
             * @return string
             */
            getCurrentViewId: function() {
                return dataUtils.getEmberApplicationSingleton().__container__.lookup("router:main").router.currentHandlerInfos[1].params.userview_id;
            }
        });

        application.register('utility:routes', routesUtils);
    }
});
