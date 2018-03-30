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
    name: 'PbehaviorAdapter',
    after: ['VeventAdapter', 'ModelsolveUtils'],
    initialize: function(container, application) {
        var VEventAdapter = container.lookupFactory('adapter:vevent');
        var modelsolve = container.lookupFactory('utility:modelsolve');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @adapter pbehavior
         */
        var adapter = VEventAdapter.extend({

            buildURL: function(type, id, record) {
                void(type);
                void(record);

                var result = '/pbehavior';

                if (!isNone(id)) {
                    result += '/' + id;
                }

                return result;
            },

            findCalendarPBehavior: function(type, query){
                var url = "/pbehavior/calendar";
                console.log('display query findAll', query);
                //return this.ajax(url, 'GET', query);

                return new Ember.RSVP.Promise(function(resolve, reject) {
                    var funcres = modelsolve.gen_resolve(resolve);
                    var funcrej = modelsolve.gen_reject(reject);
                    $.get(url, query).then(funcres, funcrej);
                });
            },

            // When New UI, use this.
            /*createRecord: function (store, type, query) {
                void (store);
                var url = '/api/v2/pbehavior'
                return this.ajax(url, 'POST', { data: query });
            },*/

            deleteRecord: function (store, type, id) {
                var url = '/api/v2/pbehavior/' + id

                return this.ajax(url, 'DELETE');
            },

            findQuery: function (store, type, id) {
                var url = '/api/v2/pbehavior_byeid/' + id

                return this.ajax(url, 'GET')
            }

        });

        application.register('adapter:pbehavior', adapter);
    }
});
