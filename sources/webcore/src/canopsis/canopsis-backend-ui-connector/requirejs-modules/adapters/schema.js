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

define(['canopsis/canopsis-backend-ui-connector/requirejs-modules/adapters/application'], function (ApplicationAdapter) {
    var shemasLimit = 200;

    /**
     * @adapter schema
     */
    var adapter = ApplicationAdapter.extend({
        findAll: function (callback) {
            return this.findAllSyncronous(callback);
        },

        //TODO make this asyncronous
        findAllSyncronous: function (successCallback) {
            return $.ajax({
                url: '/rest/schemas',
                data: {limit: shemasLimit},
                success: successCallback,
                async: false
            });

        }
    });

    Ember.Application.initializer({
        name: 'SchemaAdapter',
        initialize: function(container, application) {
            application.register('adapter:schema', adapter);
        }
    });

    return adapter;
});
