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
    name: 'AclAdapters',
    after: 'ApplicationAdapter',
    initialize: function(container, application) {

        var ApplicationAdapter = container.lookupFactory('adapter:application');

        var get = Ember.get,
            isNone = Ember.isNone;

        /**
         * @adapter role
         * @adapter group
         * @adapter account
         * @adapter user
         * @adapter right
         * @adapter profile
         */
        var adapter = ApplicationAdapter.extend({

            /**
             * @method buildURL
             * @argument type
             * @argument id
             * @argument record_or_records
             * @argument method
             * @return {string} the url of the query
             */
            buildURL: function(type, id, record_or_records, method) {
                console.log('buildURL', arguments);

                if(type === 'account') {
                    type = 'user';
                }

                if(method === 'GET') {
                    return ('/rest/default_rights/' + type + (id ? '/' + id : ''));
                } else if(type === 'action' && (method === 'POST' || method === 'PUT')) {
                    return ('/rest/default_rights/' + type + (id ? '/' + id : ''));
                } else if(method === 'DELETE') {
                    return ('/account/delete/' + type + (id ? '/' + id : ''));
                } else {
                    return ('/account/' + type + (id ? '/' + id : ''));
                }
            },

            /**
             * @method find
             */
            find: function(store, type, id, record) {
                if (isNone(type) || isNone(type.typeKey)) {
                    console.error('Error while retrieving typeKey from type is it is none.');
                }
                return this.ajax(this.buildURL(type.typeKey, id, record, 'GET'), 'GET');
            },

            /**
             * @method findMany
             */
            findMany: function(store, type, ids, records) {
                if (isNone(type) || isNone(type.typeKey)) {
                    console.error('Error while retrieving typeKey from type is it is none.');
                }
                return this.ajax(this.buildURL(type.typeKey, ids, records, 'GET'), 'GET', { data: { ids: ids } });
            },

            /**
             * @method findQuery
             */
            findQuery: function(store, type, query) {
                if (isNone(type) || isNone(type.typeKey)) {
                    console.error('Error while retrieving typeKey from type is it is none.');
                }
                return this.ajax(this.buildURL(type.typeKey, undefined, undefined, 'GET'), 'GET', { data: query });
            },

            /**
             * @method createRecord
             */
            createRecord: function(store, type, record) {
                var me = this;
                if (isNone(type) || isNone(type.typeKey)) {
                    console.error('Error while retrieving typeKey from type is it is none.');
                }

                return new Ember.RSVP.Promise(function(resolve, reject) {
                    var url = me.buildURL(type.typeKey, undefined, record, 'POST');
                    var hash = me.serialize(record, {includeId: true});

                    var data = {};
                    data[type.typeKey] = JSON.stringify(hash);

                    $.ajax({
                        url: url,
                        type: 'POST',
                        data: data
                    }).then(resolve, reject);
                });
            },

            /**
             * @method updateRecord
             */
            updateRecord: function(store, type, record) {
                var me = this;
                if (isNone(type) || isNone(type.typeKey)) {
                    console.error('Error while retrieving typeKey from type is it is none.');
                }

                return new Ember.RSVP.Promise(function(resolve, reject) {
                    var url = me.buildURL(type.typeKey, undefined, record, 'POST');
                    var hash = me.serialize(record, {includeId: true});

                    var data = {};
                    data[type.typeKey] = JSON.stringify(hash);

                    $.ajax({
                        url: url,
                        type: 'POST',
                        data: data
                    }).then(resolve, reject);
                });
            },

            /**
             * @method deleteRecord
             */
            deleteRecord: function(store, type, record) {
                var id = get(record, 'id');
                if (isNone(type) || isNone(type.typeKey)) {
                    console.error('Error while retrieving typeKey from type is it is none.');
                }
                return this.ajax(this.buildURL(type.typeKey, id, record, 'DELETE'), 'DELETE');
            }
        });

        application.register('adapter:role', adapter);
        application.register('adapter:group', adapter);
        application.register('adapter:account', adapter);
        application.register('adapter:user', adapter);
        application.register('adapter:right', adapter);
        application.register('adapter:profile', adapter);
    }
});
