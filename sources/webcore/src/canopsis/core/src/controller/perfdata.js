/**
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
    name: 'PerfdataController',
    after: ['ApplicationController', 'PojoAdapter', 'DataUtils'],
    initialize: function(container, application) {
        var dataUtils = container.lookupFactory('utility:data'),
            pojoAdapter = container.lookup('adapter:pojo'),
            applicationController = container.lookup('controller:application');

        var get = Ember.get;

        var controller = Ember.ObjectController.extend({
            fetch: function(metric_id, tstart, tend) {
                applicationController.addConcurrentLoading('perfdata');

                var requestOptions = {
                    'metric_id': metric_id,
                    'timewindow': JSON.stringify({
                        'start': tstart / 1000,
                        'stop': tend / 1000,
                        'timezone': new Date().getTimezoneOffset()
                    }),
                    'timeserie': JSON.stringify({
                        'aggregation': 'NONE'
                    })
                };

                //createRecord is used as it is a POST request
                var promise = pojoAdapter.createRecord('perfdata', undefined, requestOptions);

                promise.then(function() {
                    applicationController.removeConcurrentLoading('perfdata');
                }, function() {
                    applicationController.removeConcurrentLoading('perfdata');
                });

                return promise;
            },

            fetchMany: function(metrics, tstart, tend) {
                return this.fetch(JSON.stringify(metrics), tstart, tend);
            },

            aggregate: function(metric_id, tstart, tend, method, interval) {
                applicationController.addConcurrentLoading('perfdata');

                var requestOptions = {
                    'metric_id': metric_id,
                    'timewindow': JSON.stringify({
                        'start': tstart / 1000,
                        'stop': tend / 1000,
                        'timezone': new Date().getTimezoneOffset()
                    }),
                    'timeserie': JSON.stringify({
                        'aggregation': method,
                        'period': {
                            'second': interval
                        }
                    })
                };

                //createRecord is used as it is a POST request
                var promise = pojoAdapter.createRecord('perfdata', undefined, requestOptions);

                promise.then(function() {
                    applicationController.removeConcurrentLoading('perfdata');
                }, function() {
                    applicationController.removeConcurrentLoading('perfdata');
                });

                return promise;
            },

            aggregateMany: function(metrics, tstart, tend, method, interval) {
                return this.aggregate(JSON.stringify(metrics), tstart, tend, method, interval);
            }
        });

        application.register('controller:perfdata', controller);
    }
});
