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
    name:'LoadingindicatorMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get,
            set = Ember.set;
            isNone = Ember.isNone;


        /**
         * Mixin showing a loading indicator when there is an action that takes time. Used on the Application controller.
         *
         * @class LoadingindicatorMixin
         * @extensionfor ApplicationController
         * @static
         */
        var mixin = Mixin('loadingindicator', {
            /**
             * @method init
             */
            init: function() {
                if(isNone(get(this, 'partials.indicators'))) {
                    set(this, 'partials.indicators', []);
                }
                this.partials.indicators.pushObject('loadingindicator');

                this._super();
            },

            /**
             * @property isLoading
             * @type Number
             * @description the number of concurrent loadings (usually requests) pending
             */
            isLoading: 0,

            /**
             * @method addConcurrentLoading
             * @param {string} name the name of the loading process that is starting
             */
            addConcurrentLoading: function(name) {
                if(isNone(get(concurrentLoadingsPending, name))) {
                    set(concurrentLoadingsPending, name, { count: 1 });
                } else {
                    concurrentLoadingsPending.incrementProperty(name + '.count');
                }

                recomputeIsLoading(this);
            },

            /**
             * @method removeConcurrentLoading
             * @param {string} name the name of the loading process that is ending
             */
            removeConcurrentLoading: function(name) {
                concurrentLoadingsPending.decrementProperty(name + '.count');

                if(get(concurrentLoadingsPending, name + '.count') <= 0) {
                    delete concurrentLoadingsPending[name];
                }

                recomputeIsLoading(this);
            }
        });


        /**
         * @method recomputeIsLoading
         * @private
         * @description reassign a value to isLoading depending on if there are loading processes engaged or not
         */
        recomputeIsLoading = function(controller) {
            var totalCount = 0;

            var concurrentLoadingsKeys = Ember.keys(concurrentLoadingsPending);
            for (var i = 0; i < concurrentLoadingsKeys.length; i++) {
                var selectedConcurrentLoading = concurrentLoadingsPending[concurrentLoadingsKeys[i]];
                totalCount += get(selectedConcurrentLoading, 'count') || 0;
            }

            set(controller, 'isLoading', totalCount > 0);
        };


        /**
         * @property concurrentLoadingsPending
         * @type array
         * @private
         * @static
         * @description the list of concurrent Loadings
         */
        concurrentLoadingsPending = Ember.Object.create({
            userview: {count : 0},
            perfdata: {count : 0}
        });

        application.register('mixin:loadingindicator', mixin);
    }
});
