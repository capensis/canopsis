/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License,  or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not,  see <http://www.gnu.org/licenses/>.
 */

(function() {

    var isNone = Ember.isNone,
        __ = Ember.String.loc,
        types = {
            check : {color: 'green', icon: 'certificate'},
            consolidation : {color: 'green', icon: 'cogs'},
            log : {color: 'black', icon: 'clock-o'},
            perf : {color: 'grey', icon: 'pie-chart'},
            selector : {color: 'green', icon: 'eye'},
            sla : {color: 'green', icon: 'warning'},
            topology : {color: 'green', icon: 'code-fork'},
            eue : {color: 'green', icon: 'cog'},
            calendar : {color: 'blue', icon: 'calendar'},
            comment : {color: 'blue', icon: 'comment'},
            trap : {color: 'blue', icon: 'wrench'},
            user : {color: 'blue', icon: 'male'},
            ack : {color: 'purple', icon: 'check'},
            downtime : {color: 'yellow', icon: ''},
            cancel : {color: 'yellow', icon: 'trash'},
            uncancel : {color: 'yellow', icon: 'reply'},
            ackremove : {color: 'purple', icon: 'close'},
            assocticket : {color: 'yellow', icon: 'thumb-tack'}
        };

    Ember.Handlebars.helper('eventtype', function(eventType) {
        if (isNone(types[eventType])) {
            return eventType;
        }

        var selection = types[eventType];
        var html = '<span class="badge bg-%@"><i class="fa fa-%@"></i> %@</span>'.fmt(
            selection.color,
            selection.icon,
            __(eventType)
        );

        return new Ember.Handlebars.SafeString(html);
    });

})();
