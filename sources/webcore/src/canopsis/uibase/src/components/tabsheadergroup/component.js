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
    name: 'component-tabsheadergroup',
    initialize: function(container, application) {

        /**
         * @component tabsheadergroup
         * @description tabs subcomponent. Can be used to display tabs, and handle their content. See the "tabs" component for more information.
         */
        var component = Ember.Component.extend({
            /**
             * @property tagName
             * @type string
             * @default
             */
            tagName: 'ul',

            /**
             * @property classNames
             * @type array
             * @default
             */
            classNames: ['nav', 'nav-tabs']
        });

        application.register('component:component-tabsheadergroup', component);
    }
});
