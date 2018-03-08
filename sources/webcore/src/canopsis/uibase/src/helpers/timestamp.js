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
    name: 'TimestampHelper',
    after: 'DatesUtils',
    initialize: function(container, application) {
        void(application);

        var datesUtils = container.lookupFactory('utility:dates');

        var get = Ember.get,
            isNone = Ember.isNone;


        Ember.Handlebars.helper('timestamp', function(value, attr, record) {
            if (!isNone(record)) {
                value = get(record, 'timeStampState') || value;
            }

            if(value && !isNone(attr)) {
                format = get(attr, 'options.format');
            }

            var format;
            if (datesUtils.isToday(value)) {
                format = 'timeOnly';
            }

            if (isNone(value)) {
                return "";
            }
            if (value.length === 0) {
                return "";
            }
            if (typeof value === "object") {
                value = value.sort();
                value = value[0];
            }
            var time = datesUtils.timestamp2String(value, format, true);

            return new Ember.Handlebars.SafeString(time);
        });
    }
});
