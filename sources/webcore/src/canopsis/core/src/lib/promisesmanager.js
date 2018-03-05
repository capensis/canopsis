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
    name: 'PromisesRegistry',
    after: 'AbstractClassRegistry',
    initialize: function(container, application) {
        var Abstractclassregistry = container.lookupFactory('registry:abstractclass');

        var get = Ember.get,
            set = Ember.set;

        var registry = Abstractclassregistry.create({
            name: 'promises',

            all: [],

            pending: Ember.A(),
            errors: Ember.A(),

            pendingCount: 0,
            errorsCount: 0,

            byClass: {},

            handlePromise: function(promise) {
                var me = this;
                // Ember.run.schedule('sync', this, function() {
                    me.all.pushObject(promise);
                    me.pending.pushObject(promise);
                    set(me, 'pendingCount', me.pendingCount + 1);
                // });
            },

            promiseSuccess: function(promise) {
                console.info('promise success', promise);
            },

            promiseFail: function(promise) {
                if(promise._detail !== undefined && promise._detail.status === 200) {
                    console.warn('promise failed with error code 200, assuming it\'s a success');
                    this.promiseSuccess(promise);
                } else {
                    console.error('promise failed', promise, new Error().stack);
                    var me = this;
                    // Ember.run.schedule('sync', this, function() {
                        console.error('promise failed', promise);
                        get(me, 'errors').pushObject(promise);
                        set(me, 'errorsCount', me.errorsCount + 1);
                    // });
                }
            },

            promiseFinally: function (promise) {
                var me = this;
                Ember.run.schedule('sync', this, function() {
                    console.log('finally');
                    get(me, 'pending').removeObject(promise);
                    set(me, 'pendingCount', me.pendingCount - 1);
                });
            }
        });

        application.register('registry:promises', registry);
    }
});
