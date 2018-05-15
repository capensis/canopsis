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
    name: 'WidgetslotView',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set;

        /**
         * @class WidgetslotView
         * @extends Ember.View
         * @constructor
         */
        var view = Ember.View.extend({
            /**
             * @method init
             */
            init: function() {
                console.log('widgetslot init', get(this, 'controller.content.widgetslotTemplate'));

                var widgetslotTemplate = get(this, 'controller.content.widgetslotTemplate');

                if(widgetslotTemplate !== undefined && widgetslotTemplate !== null && Ember.TEMPLATES[widgetslotTemplate] !== undefined) {
                    set(this, 'templateName', widgetslotTemplate);
                }
                this._super.apply(this, arguments);
            },

            /**
             * @property templateName
             * @type string
             */
            templateName:'widgetslot-default',

            /**
             * @property classNames
             * @type Array
             */
            classNames: ['widgetslot']
        });

        Ember.Handlebars.helper('widgetslot', view);

        application.register('view:widgetslot', view);
    }
});
