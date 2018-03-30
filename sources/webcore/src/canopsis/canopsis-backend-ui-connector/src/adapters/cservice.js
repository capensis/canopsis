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
    name: 'CserviceAdapter',
    after: ['ApplicationAdapter'],
    initialize: function(container, application) {
        var ApplicationAdapter = container.lookupFactory('adapter:application');

        /**
         * @adapter cservice
         */
        var adapter = ApplicationAdapter.extend({
            buildURL: function(type, id) {
                type = 'cservice';

                return this._super(type, id);
            }
        });

        var schemasregistry = window.schemasRegistry;

        for(var sname in schemasregistry.all) {
            var schema = schemasregistry.getByName(sname);

            //TODO: do not use userPreferencesModelName
            var modelname = schema.modelDict.userPreferencesModelName;

            if(modelname.indexOf('crecord.cservice.') === 0) {
                console.log('Add adapter:', sname);

                var initializerName = sname.capitalize() + 'Adapter';
                Ember.Application.initializer({
                    name: initializerName,
                    initialize: function(container, application) {
                        application.register('adapter:' + sname, adapter.extend());
                    }
                });
            }
        }

        application.register('adapter:cservice', adapter);
    }
});
