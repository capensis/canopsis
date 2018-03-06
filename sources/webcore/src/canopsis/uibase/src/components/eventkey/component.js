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
    name: 'component-eventkey',
    after: 'EventUtils',
    initialize: function(container, application) {
        var eventUtils = container.lookupFactory('utility:event');

        var get = Ember.get,
            set = Ember.set,
            __ = Ember.String.loc;


        /**
         * @component eventkey
         */
        var component = Ember.Component.extend({

            init:function () {
                this._super();
                set(this,'selectedMode',__('List'));

                set(this, 'modes',[
                    __('List'),
                    __('Custom')
                ]);

                var eventFields = eventUtils.getFields();
                console.log('eventFields', eventFields);
                var selectableProperties = [];
                var length = eventFields.length;
                for (var i=0; i<length; i++) {
                    selectableProperties.push({field: eventFields[i]});
                }

                set(this, 'selectableProperties', selectableProperties);
            },

            useTextField: function () {
                return get(this, 'selectedMode') === __('Custom');
            }.property('selectedMode'),

            useList: function () {
                return get(this, 'selectedMode') === __('List');
            }.property('selectedMode'),

            testActive: function () {
                var properties = get(this, 'selectableProperties');
                var len = properties.length;
                for (var i=0; i<len; i++) {
                    var isActive = properties[i].field === get(this, 'content');
                    set(properties[i], 'isActive', isActive);
                    console.log('compare', properties[i].field ,'and',  get(this, 'content'));
                }
            }.observes('content'),

            actions : {
                setProperty: function (field) {
                    set(this, 'content', field);
                    console.log('content is now', get(this, 'content'));
                }
            }

        });

        application.register('component:component-eventkey', component);
    }
});
