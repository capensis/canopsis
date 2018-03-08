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
    name: 'componentRrule',
    initialize: function(container, application) {

        var get = Ember.get,
            RRule = window.RRule;

        /**
         * @component rrule
         * @description A rrule component to display rrules
         */
        var component = Ember.Component.extend({
            /**
             * @property {String} ruleValue Rrule value, property binded in
             * component
             */
            rruleValue: undefined,

            /**
             * @property {String} rruleTooltip Rrule to display in tooltip. Add
             * white space to wordwrap on several line.
             */
            rruleTooltip: function(){
                var value = get(this, 'rruleValue');
                return value.replace(/;/g,'; ');
            }.property('rruleValue'),

            /**
             * @property {String} rruleText Rrule value, property binded in
             * template to render the rrule
             */
            rruleText: function(){
                var value = get(this, 'rruleValue');
                var text = '';
                try{
                    var rruleObject = RRule.fromString(value);
                    text = rruleObject.toText();
                } catch(err) {
                    text = value;
                }
                return text.charAt(0).toUpperCase() + text.slice(1);
            }.property('rruleValue')
        });

        application.register('component:component-rrule', component);

    }
});
