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
    name: 'component-rights-action',
    after: 'RightsRegistry',
    initialize: function(container, application) {
        var rightsRegistry = container.lookupFactory('registry:rights');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        /**
         * @component right-action
         * @description Display a right properly, to embed into a right list
         */
        var component = Ember.Component.extend({
            /**
             * @property value
             * @description the right value
             * @type string
             */
            value: undefined,

            /**
             * @property description
             * @description Computed property, dependant on "value". the right description
             * @type string
             */
            description: function() {
                var value = get(this, 'value');
                //FIXME don't use _data, it might lead to unpredictable behaviours!
                var action = rightsRegistry.getByName(value);
                if(action && action._data) {
                    return action._data.desc;
                }
            }.property('value')
        });

        application.register('component:component-rights-action', component);
    }
});
