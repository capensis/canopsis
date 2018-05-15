/*
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

define([
    'app/lib/factories/widget',
    'app/lib/promisesmanager'
], function(WidgetFactory, promisesManager) {
    var get = Ember.get,
        set = Ember.set;

    var widget = WidgetFactory('requestsmanager', {

        originalHandlePromise: promisesManager.handlePromise,

        init: function() {
            console.error('handlePromise addObserver', promisesManager);

            var requestsmanager = this;
            promisesManager.handlePromise = function(promise){
                requestsmanager.originalHandlePromise.apply(promisesManager, arguments);
                requestsmanager.onPromise(promise);
            };
        },

        onPromise: function(promise) {
            console.error('onPromise', arguments);
            var newRequest = Ember.Object.create({
                promise: promise,
                state: '0'
            });

            promise.then(function(answer) {
                set(newRequest, 'state', 1);
                set(newRequest, 'message', answer);
            }, function(reason) {
                set(newRequest, 'state', 2);
                set(newRequest, 'message', reason);
            });

            get(this, 'widgetData').pushObject(newRequest);
        },

        widgetData: Ember.A()
    });

    return widget;
});
