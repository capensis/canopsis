/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 *  (at your option) any later version.
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
    name: 'AlertsAdapter',
    after: ['ApplicationAdapter'],
    initialize: function(container, application) {
        var ApplicationAdapter = container.lookupFactory('adapter:application');

        // Adapter beginning
        var adapter = ApplicationAdapter.extend({

            // Method called in the controller to fetch data
            findQuery: function(store, query, q) {
                var url = '/alerts/get-alarms';
                if (query == 'get-current-alarm') {
                    return this.ajax('/alerts/get-current-alarm', 'GET', {data: q});
                    // return this.ajax('/alerts/get-alarms', 'GET', {data: q});
                } else {
                    if (query["sort_key"] != undefined && query["sort_key"].substring(0, 2) != "v."){
                        query["sort_key"] = "v." + query["sort_key"];
                    }
                    return this.ajax(url, 'GET', {data: query});
                };
            }
        });

        application.register('adapter:alerts', adapter);
    }
});
