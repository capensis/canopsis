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
    name: 'TrapAdapter',
    after: 'ApplicationAdapter',
    initialize: function(container, application) {
        var ApplicationAdapter = container.lookupFactory('adapter:application');

        /**
         * @adapter trap
         */
        var adapter = ApplicationAdapter.extend({

            buildURL: function(type, id) {
                void(id);
                return '/trap';
            },

            findQuery: function(store, type, query) {

                var url = '/trap';

                if (query.skip !== undefined){
                    query.start = query.skip;
                    delete query.skip;
                }

                return this.ajax(url, 'GET', { data: query });
            }
        });

        application.register('adapter:trap', adapter);
    }
});
