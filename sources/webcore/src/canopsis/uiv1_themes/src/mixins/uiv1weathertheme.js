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
    name:'Uiv1weatherthemeMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get;

        /**
         * @mixin uiv1weathertheme
         */
        var mixin = Mixin('uiv1weathertheme', {
            partials: {
                weatherTheme: ['uiv1_themes_weather']
            },

            /**
             * @property stateImage
             */
            stateImage: function() {
                return '/static/canopsis/uiv1_themes/images/state_' + get(this, 'worst_state') + '.png';
            }.property('worst_state')
        });

        application.register('mixin:uiv1weathertheme', mixin);
    }
});
