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
    name: 'CategorychartWidget',
    after: ['WidgetFactory', 'MetricConsumer', 'TimeWindowUtils'],
    initialize: function(container, application) {
        var WidgetFactory = container.lookupFactory('factory:widget'),
            MetricConsumer = container.lookupFactory('mixin:metricconsumer'),
            TimeWindowUtils = container.lookupFactory('utility:timewindow');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var widgetOptions = {
            mixins: [MetricConsumer]
        };

        /**
         * @widget Categorychart
         * @augments Widget
         */
        var Widget = WidgetFactory('categorychart', {
            init: function() {
                this._super.apply(this, arguments);
            },

            findItems: function() {
                var props = [
                    'display',
                    'allow_user_display',
                    'use_max_value',
                    'max_value',
                    'show_legend',
                    'show_tooltip',
                    'show_labels',
                    'metric_template',
                    'stacked',
                    'text_left_space',
                    'human_readable',
                ];

                var options = {},
                    me = this;

                props.forEach(function(prop) {
                    set(options, prop, get(me, prop));
                });

                set(this, 'options', options);
                set(this, 'chartSeries', {});

                var tw = TimeWindowUtils.getFromTo(
                    get(this, 'time_window'),
                    get(this, 'time_window_offset')
                );
                var from = tw[0],
                    to = tw[1];

                /* live reporting support */
                var liveFrom = get(this, 'from'),
                    liveTo = get(this, 'to');

                if (!isNone(liveFrom)) {
                    from = liveFrom;
                }

                if (!isNone(liveTo)) {
                    to = liveTo;
                }

                var query = get(this, 'metrics');
                if (!isNone(query) && query.length) {
                    this.aggregateMetrics(
                        query,
                        from, to,
                        'last',
                        /* aggregation interval: the whole timewindow for only one point */
                        to - from
                    );
                }

                query = get(this, 'series');
                if (!isNone(query) && query.length) {
                    this.fetchSeries(query, from, to);
                }
            },

            /**
             * @method updateChart
             * @memberof CategoryChartWidget
             * Update inner chart component series.
             */
            updateChart: function() {
                var chartSeries = [];

                $.each(get(this, 'chartSeries'), function(key, serie) {
                    chartSeries.push(serie);
                });

                set(this, 'chartComponent.series', chartSeries);
            },

            onMetrics: function(metrics) {
                var chartSeries = get(this, 'chartSeries');

                $.each(metrics, function(idx, metric) {
                    var mid = get(metric, 'meta.data_id'),
                        points = get(metric, 'points');

                    /* initialize metric value */
                    var npoints = points.length,
                        value = 0;

                    if (npoints) {
                        value = points[npoints - 1][1];
                    }

                    /* compute metric name */
                    var component = undefined,
                        resource = undefined,
                        metricname = undefined,
                        midsplit = mid.split('/');

                    /* "/metric/<connector>/<connector_name>/<component>/[<resource>/]<name>"
                     * once splitted:
                     *  - ""
                     * - "metric"
                     * - "<connector>"
                     * - "<connector_name>"
                     * - "<component>"
                     * - "<resource>" and/or "<name>"
                     */
                    if (midsplit.length === 6) {
                        component = midsplit[4];
                        metricname = midsplit[5];
                    }
                    else {
                        component = midsplit[4];
                        resource = midsplit[5];
                        metricname = midsplit[6];
                    }

                    var label = component;

                    if (!isNone(resource)) {
                        label += '.' + resource;
                    }

                    label += '.' + metricname;

                    set(chartSeries, label.replace(/\./g, '_'), {
                        id: mid,
                        serie: [label, value]
                    });
                });

                set(this, 'chartSeries', chartSeries);
                this.updateChart();
            },

            onSeries: function (series) {
                var chartSeries = get(this, 'chartSeries');

                $.each(series, function(idx, serie) {
                    var points = get(serie, 'points'),
                        label = get(serie, 'label'),
                        value = 0;

                    if (points.length) {
                        value = points[points.length - 1][1];
                    }

                    set(chartSeries, label, {
                        id: label,
                        serie: [label, value]
                    });
                });

                set(this, 'chartSeries', chartSeries);
                this.updateChart();
            }
        }, widgetOptions);

        application.register('widget:categorychart', Widget);
    }
});
