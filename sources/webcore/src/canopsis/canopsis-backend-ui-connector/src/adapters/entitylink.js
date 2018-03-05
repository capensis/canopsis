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
    name: 'EntityLinkAdapter',
    after: ['ApplicationAdapter', 'ModelsolveUtils'],
    initialize: function(container, application) {
        var ApplicationAdapter = container.lookupFactory('adapter:application');
        var modelsolve = container.lookupFactory('utility:modelsolve');

        /**
         * @adapter entitylink
         */
        var adapter = ApplicationAdapter.extend({

            init: function () {
                this._super();
            },

            buildURL: function(type, id) {
                void(id);

                return '/entitylink';
            },

            findEventLinks: function(type, query) {
                var url = this.buildURL(type, null);

                console.log('findQuery', query);
                return new Ember.RSVP.Promise(function(resolve, reject) {
                    var funcres = modelsolve.gen_resolve(resolve);
                    var funcrej = modelsolve.gen_reject(reject);
                    $.post(url, query).then(funcres, funcrej);
                });
            }
        });

        application.register('adapter:entitylink', adapter);
    }
});
