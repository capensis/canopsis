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

define([], function() {
    var promisesmanager;

    var get = Ember.get,
        isNone = Ember.isNone;

    var entities = ['nagios', 'shinken'];

    /**
     * @adapter application
     */
    var adapter = DS.RESTAdapter.extend({

        /**
         * @method ajax
         * Override allowing to use the promisemanager
         */
        ajax: function(url, type, hash) {
            var promise = this._super(url, type, hash);

            promise.url = url;
            promise.type = type;
            promisesmanager.handlePromise(promise);

            promise.then(function() {
                promisesmanager.promiseSuccess(promise);
            });

            promise.catch(function() {
                promisesmanager.promiseFail(promise);
            });

            promise.finally(function() {
                promisesmanager.promiseFinally(promise);
            });

            return promise;
        },

        /**
         * @method buildURL
         * @argument type
         * @argument id
         * @return {string} the url of the query
         */
        buildURL: function(type, id) {
            var namespace = get(this, 'namespace');

            if(isNone(namespace)) {
                namespace = ( entities.contains(type) ) ? 'entities' :'object' ;
            }

            return ('/rest/'+ namespace + '/' + type + (id ? '/' + id : ''));
        },

        /**
         * @method createRecord
         */
        createRecord: function(store, type, record) {
            var data = {};
            if (isNone(type) || isNone(type.typeKey)) {
                console.error('Error while retrieving typeKey from type is it is none.');
            }
            var serializer = store.serializerFor(type.typeKey);

            data = serializer.serializeIntoHash(data, type, record, 'POST', { includeId: true });

            return this.ajax(this.buildURL(type.typeKey, record.id, undefined, 'POST'), 'POST', { data: data });
        },

        /**
         * @method updateRecord
         */
        updateRecord: function(store, type, record) {
            var data = {};
            if (isNone(type) || isNone(type.typeKey)) {
                console.error('Error while retrieving typeKey from type is it is none.');
            }
            var serializer = store.serializerFor(type.typeKey);

            data = serializer.serializeIntoHash(data, type, record, 'PUT');

            var id = Ember.get(record, 'id');

            return this.ajax(this.buildURL(type.typeKey, id, undefined, 'PUT'), 'PUT', { data: data });
        }
    });

    Ember.Application.initializer({
        name: 'ApplicationAdapter',
        after: 'PromisesRegistry',
        initialize: function(container, application) {
            promisesmanager = container.lookupFactory('registry:promises');

            application.register('adapter:application', adapter);
        }
    });

    return adapter;
});
