/**
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
    name: 'MixinFactory',
    after: 'MixinsRegistry',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var mixinsRegistry = container.lookupFactory('registry:mixins');


        /**
         * Mixin factory. Creates a controller, stores it in Application
         * @param name {string} the name of the new mixin. lowercase
         * @param classdict {dict} the controller dict
         *
         * @author Gwenael Pluchon <info@gwenp.fr>
         */
        function Mixin(name, mixindict) {
            console.tags.add('factory');
            console.group('mixin factory call', arguments);

            var mixinName = name.camelize().capitalize() + 'Mixin';

            var mixin = Ember.Mixin.create(mixindict);

            var registryEntry = Ember.Object.create({
                name: name,
                EmberClass: mixin
            });

            mixinsRegistry.add(registryEntry, name);

            console.groupEnd();
            console.tags.remove('factory');

            return mixin;
        }

        application.register('factory:mixin', Mixin);
    }
});
