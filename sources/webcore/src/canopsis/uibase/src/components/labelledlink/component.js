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
    name: 'component-labelledlink',
    initialize: function(container, application) {
        var __ = Ember.String.loc;

        /**
         * @component labelledlink
         * @description displays two inputs, to edit an url, with a link associated to it
         */
        var component = Ember.Component.extend({
            /**
             * @property label_placeholder
             * @description the placeholder for the label input
             * @type string
             * @default __('label')
             */
            label_placeholder: __('label'),

            /**
             * @property url_placeholder
             * @description the placeholder for the url input
             * @type string
             * @default __('url')
             */
            url_placeholder: __('url'),

            /**
             * @property content
             * @description an object that must have an "url" property, and a "label" property
             * @type object
             */
            content: undefined
        });

        application.register('component:component-labelledlink', component);
    }
});
