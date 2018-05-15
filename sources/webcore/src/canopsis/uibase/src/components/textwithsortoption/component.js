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
    name: 'component-textwithsortoption',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component textwithsortoption
         * @description displays a combobox with two choices for letting the user choose a sort direction (Ascending, descending), alongside a text input
         *
         * ![Component preview](../screenshots/component-textwithsortoption.png)
         */
        var component = Ember.Component.extend({

            /**
             * @property directionTypes
             * @type array
             * @description the directions to display in the combobox
             * @default
             */
            directionTypes: ['ASC', 'DESC'],

            /**
             * @property property
             * @description the text displayed in the text input
             * @type string
             */
            property: undefined,

            /**
             * @property direction
             * @description the direction displayed in the combobox
             * @type string
             */
            direction: undefined,

            /**
             * @property content
             * @type object
             */
            content: undefined,

            /**
             * @method init
             */
            init: function () {
                this._super.apply(this, arguments);
                if (!isNone(this.get('content'))) {
                    set(this, 'property', get(this, 'content.property'));
                    set(this, 'direction', get(this, 'content.direction'));
                }
            },

            /**
             * @method onUpdate
             * @description manage changes propagation
             */
            onUpdate: function() {
                this.set('content', {
                    property: get(this, 'property'),
                    direction: get(this, 'direction')
                });
                console.debug('update sortable content', get(this, 'content'));
            }.observes('property', 'direction')
        });

        application.register('component:component-textwithsortoption', component);
    }
});
