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
    name: 'TextWidget',
    after: [
        'WidgetFactory',
        'TimeWindowUtils',
        'EventConsumer',
        'MetricConsumer'
        // 'HumanReadableHelper'
    ],
    initialize: function(container, application) {
        var WidgetFactory = container.lookupFactory('factory:widget'),
            EventConsumer = container.lookupFactory('mixin:eventconsumer'),
            MetricConsumer = container.lookupFactory('mixin:metricconsumer'),
            TimeWindowUtils = container.lookupFactory('utility:timewindow');

        var Handlebars = window.Handlebars,
            __ = Ember.String.loc;

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var widgetOptions = {
            mixins: [EventConsumer, MetricConsumer]
        };

        /**
         * @widget TextWidget
         * @augments Widget
         * @description Displays a text cell, with custom content. The content of the widget can be customized with HTML and Handlebars
         * It is also possible to display information about events and perfdata.
         * # Screenshots
         *
         * ![Simple text](../screenshots/widget-text-simple.png)
         * ![Event custom html](../screenshots/widget-text-customhtml1.png)
         */
        var widget = WidgetFactory('text', {
            /**
             * @method findItems
             */
            findItems: function() {
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

                set(this, 'context', {
                    event: {},
                    serie: {},
                    metric: {}
                });
                set(this, 'context.from', from);
                set(this, 'context.to', to);

                var datafetch = false;

                var query = get(this, 'events');
                if (!isNone(query) && query.length) {
                    this.fetchEvents(query);
                    datafetch = true;
                }

                query = get(this, 'metrics');
                if(get(this, 'appliedDynamicProperties') && get(this, 'appliedDynamicProperties') !== {}) {
                    query = JSON.stringify(query);
                    var context = get(this, 'appliedDynamicProperties');
                    query = Handlebars.compile(query)(context);
                    query = JSON.parse(query);
                }
                if (!isNone(query) && query.length) {
                    this.aggregateMetrics(
                        query,
                        from, to,
                        'last',
                        /* aggregation interval: the whole timewindow for only one point */
                        to - from
                    );
                    datafetch = true;
                }

                query = get(this, 'series');
                if (!isNone(query) && query.length) {
                    this.fetchSeries(query, from, to);
                    datafetch = true;
                }

                /* when no data is requested, just render the template */
                if (datafetch === false) {
                    this.renderTemplate();
                }
            },

            /**
             * @method onEvents
             * @argument events
             * @argument labelsByRk
             */
            onEvents: function(events, labelsByRk) {
                var context = get(this, 'context');

                $.each(events, function(idx, evt) {
                    var rk = get(evt, 'id');
                    var label = labelsByRk[rk];

                    if (!isNone(label)) {
                        set(context, 'event.' + label, evt._data);
                    }
                    else {
                        console.warn('No label found for event, will not be rendered:', rk);
                    }
                });

                set(this, 'context', context);
                this.renderTemplate();
            },

            /**
             * @method onMetrics
             * @argument metrics
             */
            onMetrics: function(metrics) {
                var context = get(this, 'context');

                $.each(metrics, function(idx, metric) {
                    var mid = get(metric, 'meta.data_id').split('/'),
                        points = get(metric, 'points');

                    /* initialize metric value for template context */
                    var npoints = points.length,
                        value = __('No data available');

                    if (npoints) {
                        value = points[npoints - 1][1];
                    }

                    /* compute template context path of metric */
                    var component = undefined,
                        resource = undefined,
                        metricname = undefined;

                    /* "/metric/<connector>/<connector_name>/<component>/[<resource>/]<name>"
                     * once splitted:
                     *  - ""
                     * - "metric"
                     * - "<connector>"
                     * - "<connector_name>"
                     * - "<component>"
                     * - "<resource>" and/or "<name>"
                     */
                    if (mid.length === 6) {
                        component = mid[4];
                        metricname = mid[5];
                    }
                    else {
                        component = mid[4];
                        resource = mid[5];
                        metricname = mid[6];
                    }

                    //fix for metric with . in their names
                    component = component.replace(/\./g,"_");
                    metricname = metricname.replace(/\./g,"_");
                    if(resource !== undefined)
                        resource = resource.replace(/\./g,"_");

                    var varname = 'metric.' + component;

                    if (isNone(get(context, varname))) {
                        set(context, varname, {});
                    }

                    if (!isNone(resource)) {
                        varname += '.' + resource;

                        if (isNone(get(context, varname))) {
                            set(context, varname, {});
                        }
                    }

                    set(context, varname + '.' + metricname, value);
                });

                set(this, 'context', context);
                this.renderTemplate();
            },

            /**
             * @method onSeries
             * @argument series
             */
            onSeries: function(series) {
                var context = get(this, 'context');

                $.each(series, function(idx, serie) {
                    var points = get(serie, 'points'),
                        label = get(serie, 'label');

                    var value = __('No data available');

                    if (points.length) {
                        value = points[points.length - 1][1];
                    }

                    set(context, 'serie.' + label, value);
                });

                set(this, 'context', context);
                this.renderTemplate();
            },

            /**
             * @method makeTemplate
             * @description Make sure template has been compiled.
             */
            makeTemplate: function() {
                var template = undefined;
                try {
                    template = Handlebars.compile(get(this, 'html'));
                }
                catch(err) {
                    template = function() {
                        return '<i>Impossible to render template:</i> ' + err;
                    };
                }
                return template;
            },

            /**
             * @method renderTemplate
             * @description Render compiled template property with context property into the rendered property.
             */
            renderTemplate: function() {
                var template = this.makeTemplate();

                var context = get(this, 'context');

                var appliedDynamicProperties = get(this, 'appliedDynamicProperties') || {};
                $.extend(context, appliedDynamicProperties);


                set(this, 'renderedTemplate', new Ember.Handlebars.SafeString(template(context)));
            }
        }, widgetOptions);

        application.register('widget:text', widget);
    }
});
