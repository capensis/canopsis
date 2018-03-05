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
    name: 'ContextAdapters',
    after: 'ApplicationAdapter',
    initialize: function(container, application) {
        var ApplicationAdapter = container.lookupFactory('adapter:application');

        var isNone = Ember.isNone,
            get = Ember.get;

        /**
         * @adapter context
         */
        var adapter = ApplicationAdapter.extend({

            buildURL: function(type, id) {
                var url = '/context';

                if (!isNone(type)) {
                    url += '/' + type;
                }

                if (!isNone(id)) {
                    url += '/' + id;
                }

                return url;
            },

            createRecord: function(store, type, record) {
                return this.updateRecord(store, type, record);
            },

            updateRecord: function(store, model, record) {
                var id = record.get('_data._id');

                var serializer = store.serializerFor(model.typeKey),
                    url = '/context',
                    data = {};

                data = serializer.serializeIntoHash(
                        data, model, record, 'PUT'
                );

                // rewrite id because it comes with 'genertedId' prefix
                data._id = id
                data.id = id
                // set values that are not defind in the schima
                data.infos = record._data.infos
                data.links = record._data.links
                data.depends = record._data.depends
                data.impact = record._data.impact
                
                

                var query = {
                    _type: 'context',
                    entity: data
                };

                return this.ajax(url, 'PUT', {data: query});
            },

            deleteRecord: function(store, type, record) {
                var url = this.buildURL();

                var id = get(record, 'id');
                var query = {data: {ids: id}};

                return this.ajax(url, 'DELETE', query);
            },

            findQuery: function(store, model, query) {
                void(store);
                var url = this.buildURL();

                if(typeof (query.filter) !== 'string') {
                    query._filter = query.filter;
                }
                else {
                    query._filter = query.filter;
                }

                if(isNone(query.limit)) {
                    query.limit = 5;
                }

                delete query.filter;

                console.log('findQuery: ', query);

                return this.ajax(url, 'POST', {data: query});
            }
        });

        application.register('adapter:context', adapter);
        application.register('adapter:entity', adapter);

        application.register('adapter:ctxconnector', adapter.extend({
            buildURL: function(type, id) {
                return '/context/connector' + (id ? ('/' + id) : '');
            }
        }));


        application.register('adapter:ctxconnectorname', adapter.extend({
            buildURL: function(type, id) {
                return '/context/connector_name' + (id ? ('/' + id) : '');
            }
        }));

        application.register('adapter:ctxcomponent', adapter.extend({
            buildURL: function(type, id) {
                return '/context/component' + (id ? ('/' + id) : '');
            }
        }));

        application.register('adapter:ctxresource', adapter.extend({
            buildURL: function(type, id) {
                return '/context/resource' + (id ? ('/' + id) : '');
            }
        }));

        application.register('adapter:ctxmetric', adapter.extend({
            buildURL: function(type, id) {
                return '/api/context/metric' + (id ? ('/' + id) : '');
            }
        }));

        application.register('adapter:ctxhostgroup', adapter.extend({
            buildURL: function(type, id) {
                return '/context/hostgroup' + (id ? ('/' + id) : '');
            }
        }));

        application.register('adapter:ctxservicegroup', adapter.extend({
            buildURL: function(type, id) {
                return '/context/servicegroup' + (id ? ('/' + id) : '');
            }
        }));

        application.register('adapter:ctxtopology', adapter.extend({
            buildURL: function(type, id) {
                return '/context/topo' + (id ? ('/' + id) : '');
            }
        }));

        application.register('adapter:ctxselector', adapter.extend({
            buildURL: function(type, id) {
                return '/context/selector' + (id ? ('/' + id) : '');
            }
        }));
    }
});
