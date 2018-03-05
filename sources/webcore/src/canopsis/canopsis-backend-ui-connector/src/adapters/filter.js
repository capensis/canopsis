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
    name: 'FilterAdapter',
    after: 'ApplicationAdapter',
    initialize: function(container, application) {
        var ApplicationAdapter = container.lookupFactory('adapter:application');

        var set = Ember.set;

        /**
         * @adapter filter
         */
        var adapter = ApplicationAdapter.extend({
            updateRecord: function(store, type, record) {
                //This value have to be reseted each update for user display purpose
                set(record, 'run_once', false);
                var promise = this._super(store, type, record);
                return promise;
            }
        });

        application.register('adapter:filter', adapter);
    }
});
