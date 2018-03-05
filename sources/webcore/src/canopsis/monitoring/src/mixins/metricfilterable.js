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
    name: 'MetricFilterable',
    after: ['MixinFactory'],
    initialize: function(container, application) {
        var MixinFactory = container.lookupFactory('factory:mixin');

        var get = Ember.get;

        /**
         * @mixin MetricFilterable
         * @augments Mixin
         * Provide metric filter support to widget.
         */
        var mixin = MixinFactory('metricfilterable', {
            /**
             * @method getMetricFilter
             * @memberof MetricFilterableMixin
             * Using ``metric_filter`` field, query the context to find metric IDs.
             */
            getMetricFilter: function() {
                var metric_regex = get(this, 'metric_filter');
                var regex_parts = metric_regex.split(' ');

                var regex = {
                    component: [],
                    resource: [],
                    name: []
                };

                $.each(regex_parts, function(idx, part) {
                    var prefix = part.slice(0, 3),
                        spec = {'$regex': part.slice(3)};

                    if (prefix === 'co:') {
                        regex.component.push(spec);
                    }
                    else if (prefix === 're:') {
                        regex.resource.push(spec);
                    }
                    else if (prefix === 'me:') {
                        regex.name.push(spec);
                    }
                    else {
                        spec = {'$regex': part};

                        regex.component.push(spec);
                        regex.resource.push(spec);
                        regex.name.push(spec);
                    }
                });

                var mfilter = {'$and': []};

                $.each(regex, function(key, items) {
                    if (items.length > 0) {
                        var local_mfilter = {'$or': []};

                        $.each(items, function(idx, item) {
                            var spec = {};
                            spec[key] = item;

                            local_mfilter['$or'].push(spec);
                        });

                        if (local_mfilter['$or'].length === 1) {
                            local_mfilter = local_mfilter['$or'][0];
                        }

                        mfilter['$and'].push(local_mfilter);
                    }
                });

                if (mfilter['$and'].length === 0) {
                    mfilter = {};
                }
                else if (mfilter['$and'].length === 1) {
                    mfilter = mfilter['$and'][0];
                }

                return mfilter;
            }
        });

        application.register('mixin:metricfilterable', mixin);
    }
});
