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
    name: 'component-tabcontent',
    initialize: function(container, application) {
        var get = Ember.get,
            set = Ember.set;

        /**
         * @component tabscontent
         * @description tabs subcomponent. Can be used to display tabs, and handle their content. See the "tabs" component for more information.
         */
        var component = Ember.Component.extend({
            /**
             * @property tabContainer
             * @type Ember.Component
             * @default Ember.computed.alias('parentView.parentView')
             * @description the root "tabs" component
             */
            tabContainer: Ember.computed.alias('parentView.parentView'),

            /**
             * @property classNames
             * @type array
             * @default
             */
            classNames: ['tab-pane'],

            /**
             * @property attributeBindings
             * @type array
             * @default
             */
            attributeBindings: ['role', 'data-ref'],

            /**
             * @property role
             * @type string
             * @default
             */
            'role': 'tab',

            /**
             * @property data-ref
             * @type string
             * @default Ember.computed.alias('ref')
             */
            'data-ref': Ember.computed.alias('ref'),

            /**
             * @property active
             * @type boolean
             * @description whether the tab is active by default or not
             * @default
             */
            active: false,

            /**
             * @method init
             * @description assign elementId accordingly to parent tab container elementId, and check if the content must be displayed
             */
            init: function() {
                set(this, 'elementId', get(this, 'tabContainer.elementId') + '-' + get(this, 'ref'));
                if(get(this, 'active') && get(this, 'active') === true) {
                    get(this, 'classNames').pushObject('active');
                } else {
                    get(this, 'classNames').removeObject('active');
                }
                this._super();
            }
        });

        application.register('component:component-tabcontent', component);
    }
});
