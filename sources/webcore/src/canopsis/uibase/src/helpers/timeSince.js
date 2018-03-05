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
    name: 'TimesinceHelper',
    after: 'DatesUtils',
    initialize: function(container, application) {
        void(application);

        var datesUtils = container.lookupFactory('utility:dates');
        var __ = Ember.String.loc;

        Ember.Handlebars.helper('timeSince', function(timestamp , record) {

            if(timestamp || record.timeStampState) {
                timestamp = record.timeStampState || timestamp;


                if (datesUtils.isToday(timestamp)) {
                    //This is today
                    return __('Today');
                } else {
                    var time = datesUtils.durationFromNow(timestamp);
                    return time;
                }
            } else {
                return '';
            }
        });
    }
});
