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
    name: 'component-stateeditor',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component stateeditor
         * @description Displays buttons to change an event state. There are 4 states (info, minor, major, critical). The button corresponding to the event's current state is not displayed.
         *
         * ![Component preview](../screenshots/component-changestate.png)
         *
         * @example {{component-stateeditor content=attr.value title=attr.field showAll=attr.model.options.showAll}}
         */
        var component = Ember.Component.extend({
            /**
             * @property previousContent
             * @type integer
             * @description A backup of the initial event' state.
             */
            previousContent: undefined,

            /**
             * @property content
             * @type integer
             * @description the event' state.
             */
            content: undefined,
            /**
             * @method init
             */
            init: function() {
                this._super();
                set(this, 'previousContent', get(this, 'content'));
                if(isNone(get(this, 'hidePrevious'))) {
                    //arbitrary default value currently used for change criticity action.
                    set(this, 'hidePrevious', true);
                }
            },

            /**
             * @property isInfo
             * @type boolean
             * @description Computed property dependent on "content". Returns true if the event' state is info.
             */
            isInfo:function () {
                return get(this, 'content') === 0;
            }.property('content'),

            /**
             * @property isMinor
             * @type boolean
             * @description Computed property dependent on "content". Returns true if the event' state is minor.
             */
            isMinor:function () {
                return get(this, 'content') === 1;
            }.property('content'),

            /**
             * @property isMajor
             * @type boolean
             * @description Computed property dependent on "content". Returns true if the event' state is major.
             */
            isMajor:function () {
                return get(this, 'content') === 2;
            }.property('content'),

            /**
             * @property isCritical
             * @type boolean
             * @description Computed property dependent on "content". Returns true if the event' state is critical.
             */
            isCritical:function () {
                return get(this, 'content') === 3;
            }.property('content'),

            /**
             * @method previousIs
             * @param {integer} state the state to check
             * @returns boolean
             * @description Returns true if the event' state is the state specified in the method parameter.
             */
            previousIs: function (state) {
                if (get(this, 'showAll')) {
                    return false;
                }
                return get(this, 'hidePrevious') && get(this, 'previousContent') === state;
            },

            /**
             * @property previousIsInfo
             * @type boolean
             * @description Computed property dependent on "previousContent". is "True" if the event's previous state is info.
             */
            previousIsInfo:function () {
                return this.previousIs(0);
            }.property('previousContent'),

            /**
             * @property previousIsMinor
             * @type boolean
             * @description Computed property dependent on "previousContent". is "True" if the event's previous state is minor.
             */
            previousIsMinor:function () {
                return this.previousIs(1);
            }.property('previousContent'),

            /**
             * @property previousIsMajor
             * @type boolean
             * @description Computed property dependent on "previousContent". is "True" if the event's previous state is major.
             */
            previousIsMajor:function () {
                return this.previousIs(2);
            }.property('previousContent'),

            /**
             * @property previousIsCritical
             * @type boolean
             * @description Computed property dependent on "previousContent". is "True" if the event's previous state is critical.
             */
            previousIsCritical:function () {
                return this.previousIs(3);
            }.property('previousContent'),


            actions: {
                /**
                 * @method actions_setState
                 * @param {integer} state
                 * @description changes the state of the event.
                 */
                setState:function (state) {
                    set(this, 'content', parseInt(state));
                }
            }
        });

        application.register('component:component-stateeditor', component);
    }
});
