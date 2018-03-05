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
    name: 'UserviewAdapter',
    after: 'ApplicationAdapter',
    initialize: function(container, application) {
        var ApplicationAdapter = container.lookupFactory('adapter:application');

        /**
         * @adapter userview
         */
        var adapter = ApplicationAdapter.extend({
            buildURL: function(type, id, arg3, method) {
                if(method === 'DELETE') {
                    return '/userview/' + id;
                }

                return this._super('view', id);
            },

            /**
             * @method deleteRecord
             */
            deleteRecord: function(store, type, record) {
                var id = Ember.get(record, 'id');

                return this.ajax(this.buildURL('userview', id, undefined, 'DELETE'), 'DELETE');
            }
        });

        application.register('adapter:userview', adapter);
    }
});
