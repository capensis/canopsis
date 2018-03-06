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
    name: 'component-selectedmetricheader',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        var component = Ember.Component.extend({
            /**
             * @property rk
             * @description the rk of the metric
             * @type string
             */
            rk: undefined,

            splittedRk: function() {
                return get(this, 'rk').split('/');
            }.property('rk'),

            connectorType: function() {
                return get(this, 'splittedRk')[2];
            }.property('splittedRk'),

            component: function() {
                return get(this, 'splittedRk')[4];
            }.property('splittedRk'),

            resource: function() {
                var splittedRk = get(this, 'splittedRk');
                if(splittedRk.length === 7) {
                    return splittedRk[5];
                } else {
                    return '';
                }
            }.property('splittedRk'),

            metricName: function() {
                var splittedRk = get(this, 'splittedRk');
                return splittedRk[splittedRk.length - 1];
            }.property('splittedRk')
        });

        application.register('component:component-selectedmetricheader', component);
    }
});
