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
    name: 'HumanReadableHelper',
    after: ['ValuesUtils'],
    initialize: function(container, application) {
        void(application);
        var ValuesUtils = container.lookupFactory('utility:values');

        var get = Ember.get,
            isNone = Ember.isNone,
            __ = Ember.String.loc,
            Handlebars = window.Handlebars;


        /**
         * @function HumanReadableHelper
         * @param {object} value - Value with unit to humanize
         * @returns {string} Humanized value.
         * Handlebars helper used to humanize variables using ValuesUtility
         */
        var helper = function(value, options) {
            var val = value,
                ret;

            if (isNone(options)) {
                options = value;
                val = get(options, 'hash.value');
            }

            var unit = get(options, 'hash.unit') || '';

            if (isNaN(val)) {
                ret = __('Not a valid number');
            }
            else {
                ret = ValuesUtils.humanize(val, unit);
            }

            return new Ember.Handlebars.SafeString(ret);
        };

        Handlebars.registerHelper('hr', helper);
        Ember.Handlebars.helper('hr', helper);
    }
});
