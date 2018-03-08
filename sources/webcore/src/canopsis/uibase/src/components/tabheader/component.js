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
    name: 'component-tabheader',
    initialize: function(container, application) {
        var get = Ember.get;

        /**
         * @component tabheader
         * @description tabs subcomponent. Can be used to display tabs, and handle their content. See the "tabs" component for more information.
         */
        var component = Ember.Component.extend({
            /**
             * @property tagName
             * @type string
             * @default
             */
            tagName: 'li',

            /**
             * @property attributeBindings
             * @type array
             * @default
             */
            attributeBindings: ['data-toggle', 'role', 'data-ref'],

            /**
             * @property data-toggle
             * @type string
             * @default
             */
            'data-toggle': 'tab',

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
             * @property href
             * @type string
             * @description Computed property to generate the anchor link for the tab header label. Dependent on the tabContainer id and the ref property
             */
            href: function() {
                var id = get(this, 'tabContainer.elementId');
                return '#' + id + '-' + get(this, 'ref');
            }.property('tabContainer', 'ref'),

            /**
             * @property tabContainer
             * @type Ember.Component
             * @default Ember.computed.alias('parentView.parentView')
             * @description the root "tabs" component
             */
            tabContainer: Ember.computed.alias('parentView.parentView'),

            /**
             * @property active
             * @description whether the tab header must be active or not
             * @type boolean
             * @default
             */
            active: false,

             /**
             * @method init
             * @description check if tab must be assigned an "active" class
             */
            init: function() {
                if(get(this, 'active') && get(this, 'active') === true) {
                    get(this, 'classNames').pushObject('active');
                } else {
                    get(this, 'classNames').removeObject('active');
                }
                this._super();
            }
        });

        application.register('component:component-tabheader', component);
    }
});
