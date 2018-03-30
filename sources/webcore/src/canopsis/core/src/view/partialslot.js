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
    name: 'PartialslotView',
    initialize: function(container, application) {
        /**
         * @class PartialslotView
         * @extends Ember.View
         * @constructor
         */
        var view = Ember.View.extend({

            /**
             * @property canopsisConfiguration
             * @type object
             * @description the canopsis frontend configuration object
             */
            canopsisConfiguration: canopsisConfiguration,

            /**
             * @property showPartialslots
             * @type boolean
             * @default Ember.computed.alias('canopsisConfiguration.showPartialslots')
             */
            showPartialslots: Ember.computed.alias('canopsisConfiguration.showPartialslots'),

            /**
             * @property templateName
             * @type string
             */
            templateName: 'partialslot',

            /**
             * @property tagName
             * @type string
             */
            tagName: 'span',

            /**
             * @property slot
             * @type string
             * @required
             */
            slot: Ember.required()
        });

        application.register('view:partialslot', view);
        Ember.Handlebars.helper('partialslot', view);
    }
});
