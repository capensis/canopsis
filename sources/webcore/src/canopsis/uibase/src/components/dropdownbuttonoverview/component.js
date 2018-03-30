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
    name: 'component-dropdownbuttonoverview',
    initialize: function(container, application) {

        var get = Ember.get;

        /**
         * @description Component for seeing the chosen color in the dropdownbutton
         * @component Dropdownbuttonoverview
         */
        var component = Ember.Component.extend({
            /**
             * @property classNames {Array}
             * @default
             */
            classNames: ['dropdownbuttonoverview', 'dropdownbuttonoverview-default', 'overview'],

            //update background color of the overview
            attributeBindings: ['style'],

            /**
             * @description return the css style with the right background color
             * @method style
             * @return {String}
             */
            style: function() {
                var code = get(this, 'color');
                return 'background-color:' + code;
            }.property('color')
        });

        application.register('component:component-dropdownbuttonoverview', component);
    }
});
