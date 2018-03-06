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
    name: 'CrudContextAdapter',
    after: 'ApplicationAdapter',
    initialize: function(container, application) {
        var ApplicationAdapter = container.lookupFactory('adapter:application');

        /**
         * @adapter alarm
         */
        var adapter = ApplicationAdapter.extend({

            generateUUID: function() {
                var d = new Date().getTime();
                var uuid = 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, function(c) {
                    var r = (d + Math.random()*16)%16 | 0;
                    d = Math.floor(d/16);
                    return (c=='x' ? r : (r&0x3|0x8)).toString(16);
                });
                return uuid;
            },

            findQuery: function(store, model, query) {

                return $.ajax({
                    method: "GET",
                    url: "/api/v2/context",
                    dataType: "json"
                });
            },

            createRecord: function (store, model, record) {
                var data = {};
                var serializer = store.serializerFor(model.typeKey);
                var url = '/context'

                data = serializer.serializeIntoHash(
                    data, model, record, 'POST', { includeId: true }
                );

                var query = {
                    entity: data[0]
                };

                return new Ember.RSVP.Promise(function(resolve, reject) {
                    $.ajax({
                        url: url,
                        type: 'PUT',
                        data: JSON.stringify(query)
                        
                    }).then(resolve, reject);
                });
            },

            deleteRecord: function(store, model, record) {
                var id = record.get('_data._id');

                return this.ajax('/context', 'DELETE', {data: {ids: id}});
            },

            updateRecord: function(store, model, record) {
                // if (isNone(type) || isNone(type.typeKey)) {
                //     throw new Error('Cannot determine type of record: type or type.typeKey is None');
                // }

                var id = record.get('_data._id');

                

                var serializer = store.serializerFor(model.typeKey),
                    url = '/context',
                    data = {};

                data = serializer.serializeIntoHash(
                        data, model, record, 'PUT'
                );

                data._id = id
                data.depends = record._data.depends
                data.impact = record._data.impact
                data.type = record._data.type
                data.infos = record._data.infos
                data.links = record._data.links
                data.measurements = record._data.measurements

                var query = {
                    _type: 'crudcontext',
                    entity: data
                };

                return this.ajax(url, 'PUT', {data: query});
            }
        });

        application.register('adapter:crudcontext', adapter);
    }
});