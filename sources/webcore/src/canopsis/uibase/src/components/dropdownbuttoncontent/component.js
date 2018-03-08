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
    name: 'component-dropdownbuttoncontent',
    initialize: function(container, application) {

        var get = Ember.get;

        /**
         * @description Component for switching between display and hide of the content
         * @component dropdownbuttoncontent
         */
        var component = Ember.Component.extend({
            /**
             * @property classNames {Array}
             * @default
             */
            classNames: ['dropdownbuttoncontent'],
            /**
             * @property classNameBindings {Array}
             * @default
             */
            classNameBindings: ['dropdownContentMenu'],

            /**
             * @property align {String}
             * @description the alignment of the content container. must be "right" or "left"
             */
            align:'left',

            /**
             * @description Initiate the component
             * @method init
             */
            init: function() {
                var align = get(this, 'align');

                this._super();

                if(align === 'right') {
                    get(this, 'classNames').pushObject('dropdown-menu-right');
                }
            },

            /**
             * @description Method for defining a boolean value on dropdownContentMenu thanks to opened attribute
             * @method dropdownContentMenu
             * @return {boolean}
             */
            dropdownContentMenu: function(){
                return get(this, 'parentView.opened');
            }.property('parentView.opened'),

            attributeBindings: ['aria-labelledby'],
            'aria-labelledby': 'dropdownMenu1'
        });

        application.register('component:component-dropdownbuttoncontent', component);
    }
});
