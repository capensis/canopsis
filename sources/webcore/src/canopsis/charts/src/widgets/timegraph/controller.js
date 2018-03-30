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
    name: 'TimegraphWidget',
    after: ['WidgetFactory', 'ValuesUtils', 'DatesUtils'],
    initialize: function(container, application) {
        var WidgetFactory = container.lookupFactory('factory:widget');
        var values = container.lookupFactory('utility:values');
        var dates = container.lookupFactory('utility:dates');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var widgetOptions = {};

        var FlotChartViewMixin = Ember.Mixin.create({
            didInsertElement: function() {
                var ctrl = get(this, 'controller');

                console.group('timegraph init');

                var updateGrid = function(evt, ranges) {
                    // little hack so chartOptions will always notify its observers
                    var chartOptions = get(ctrl, 'chartOptions');

                    var opts = {};
                    $.extend(opts, chartOptions);

                    if(ranges.xaxis) {
                        set(opts, 'xaxis.min', ranges.xaxis.from);
                        set(opts, 'xaxis.max', ranges.xaxis.to);
                    }

                    set(ctrl, 'chartOptions', opts);

                    set(ctrl, 'zooming', true);
                };

                // fill chart options
                this.setDefaultChartOptions();

                var graphcontainer = this.$('.flotchart-plot-container .flotchart');
                graphcontainer.bind('plotselected', updateGrid);

                if (get(ctrl, 'config.timenav')) {
                    this.setDefaultTimenavOptions();

                    var timecontainer = this.$('.flotchart-preview-container .flotchart');
                    timecontainer.bind('plotselected', updateGrid);

                    var component = Ember.View.views[timecontainer.attr('id')];

                    graphcontainer.bind('toggleserie', function(evt, config) {
                        component.send('renderChart');
                    });
                }
                ctrl.on('onRenderer', this, this.onRenderer);
                console.groupEnd();

                this._super.apply(this, arguments);
            },

            willDestroyElement: function() {
                var ctrl = get(this, 'controller');
                var graphcontainer = this.$('.flotchart-plot-container .flotchart');
                graphcontainer.unbind('plotselected');

                if (get(this, 'controller.config.timenav')) {
                    var timecontainer = this.$('.flotchart-preview-container .flotchart');
                    timecontainer.unbind('plotselected');
                    graphcontainer.unbind('toggleserie');
                }
                ctrl.off('onRenderer', this, this.onRenderer);
            },

            setDefaultChartOptions: function() {
                //get the timestamp, and not the date object
                var now = +new Date();

                var ctrl = get(this, 'controller');
                var config = get(ctrl, 'config'),
                    template = get(ctrl,'model.metric_template');

                var chartOptions = {};
                $.extend(chartOptions, get(ctrl, 'chartOptions'));

                if(!get(ctrl, 'zooming')) {
                    $.extend(chartOptions, {
                        xaxis: {
                            min: now - get(ctrl, 'time_window_offset') - get(ctrl, 'time_window'),
                            max: now - get(ctrl, 'time_window_offset')
                        }
                    });
                }

                $.extend(chartOptions, {
                    zoom: {
                        interactive: false
                    },

                    selection: {
                        mode: 'x'
                    },

                    crosshair: {
                        mode: 'x'
                    },

                    grid: {
                        hoverable: true,
                        clickable: true,
                        borderWidth: 2
                    },

                    yaxis: {},

                    xaxes: [{
                        show: true,
                        reserveSpace: true,
                        position: 'bottom',
                        mode: 'time',
                        timezone: 'browser'
                    }],

                    yaxes: [{
                        show: true,
                        reserveSpace: true
                    }],

                    legend: {
                        show: get(config, 'legend'),
                        container: this.$('.flotchart-legend-container')
                    },
                    tooltip: get(config, 'tooltip'),
                    tooltipOpts: {
                        id: this.$().closest('.ember-view').attr('id') + '-tooltip',
                        content: function(label, xval, yval, item) {
                            var date = dates.timestamp2String(xval / 1000, 'f', true);
                            var unit = item.series.unit || '';
                            var html = ['<p>', date, '</p><p><b>', label, ' :</b> '].join('');

                            if (get(ctrl, 'human_readable') === true) {
                                html += values.humanize(yval, unit || '');
                            }
                            else {
                                html += yval + ' ' + unit;
                            }

                            return html;
                        },
                        shifts: {
                            x: -60,
                            y: 25
                        }
                    }
                });

                console.log('Configure chart:', chartOptions);
                set(ctrl, 'chartOptions', chartOptions);
            },

            setDefaultTimenavOptions: function() {
                //get the timestamp, and not the date object
                var now = +new Date();
                var ctrl = get(this, 'controller');
                var config = get(ctrl, 'config');

                var chartOptions = {};
                $.extend(chartOptions, get(ctrl, 'timenavOptions'));
                $.extend(chartOptions, {
                    zoom: {
                        interactive: false
                    },

                    selection: {
                        mode: 'x'
                    },

                    crosshair: {
                        mode: 'x'
                    },

                    grid: {
                        hoverable: true,
                        clickable: true,
                        borderWidth: 2
                    },

                    xaxis: {
                        reserveSpace: true,
                        min: now - get(ctrl, 'time_window_offset') - get(ctrl, 'timenav_window'),
                        max: now - get(ctrl, 'time_window_offset')
                    },

                    yaxis: {},

                    xaxes: [{
                        show: true,
                        reserveSpace: true,
                        position: 'bottom',
                        mode: 'time',
                        timezone: 'browser'
                    }],

                    yaxes: [{
                        show: true,
                        reserveSpace: true
                    }],

                    legend: {
                        show: false
                    },
                    tooltip: false
                });

                console.log('Configure time navigation:', chartOptions);
                set(ctrl, 'timenavOptions', chartOptions);
            },

            resetZoom: function() {
                var ctrl = get(this, 'controller');

                set(ctrl, 'zooming', false);

                this.setDefaultChartOptions();
                this.setDefaultTimenavOptions();

            },

            onRenderer: function() {
                this.rerender();
            },

            actions: {
                stepBack: function() {
                    var ctrl = get(this, 'controller');
                    var step = get(ctrl, 'timestep');

                    var opts = {};
                    $.extend(opts, get(ctrl, 'chartOptions'));

                    opts.xaxis.min -= step;
                    opts.xaxis.max -= step;

                    set(ctrl, 'chartOptions', opts);
                },

                stepForward: function() {
                    var ctrl = get(this, 'controller');
                    var step = get(ctrl, 'timestep');

                    var opts = {};
                    $.extend(opts, get(ctrl, 'chartOptions'));

                    opts.xaxis.min += step;
                    opts.xaxis.max += step;

                    set(ctrl, 'chartOptions', opts);
                }
            },

            //Controller -> View Hooks
            registerHooks: function() {
                console.log("registerHooks", get(this, "controller"), get(this, "controller").on);
                get(this, "controller").on('resetZoom', this, this.resetZoom);
                return this._super();
            },

            unregisterHooks: function() {
                get(this, "controller").off('resetZoom', this, this.resetZoom);
                return this._super();
            }
        });

        var widget = WidgetFactory('timegraph', {
            needs: ['serie', 'perfdata'],

            viewMixins: [
                FlotChartViewMixin
            ],

            partials: {
                titlebarbuttons: [
                    'titlebarbutton-resetzoom'
                ],
                widgetActionButtons: [
                    'timegraphbutton-resetzoom'
                ]
            },

            actions: {
                resetZoom: function () {
                    this.trigger('resetZoom');
                }
            },

            init: function() {
                set(this, 'zooming', false);
                set(this, 'chartOptions', undefined);
                set(this, 'timenavOptions', undefined);
                set(this, 'chartSeries', Ember.Object.create({}));
                set(this, 'dataSeries', Ember.A());

                this._super.apply(this, arguments);
            },

            human_readable: function() {
                return get(this, 'config.human_readable');
            }.property('config.human_readable'),

            time_window: function() {
                return get(this, 'config.time_window') * 1000;
            }.property('config.time_window'),

            time_window_offset: function() {
                return get(this, 'config.time_window_offset') * 1000;
            }.property('config.time_window_offset'),

            timenav_window: function() {
                if(get(this, 'config.timenav')) {
                    return get(this, 'config.timenav_window') * 1000;
                }
                else {
                    return get(this, 'time_window');
                }
            }.property('config.timenav_window'),

            timestep: function() {
                return get(this, 'config.timestep') * 1000;
            }.property('config.timestep'),

            findItems: function() {
                var me = this;

                var replace = true;
                var now = new Date().getTime();
                var to = now - get(this, 'time_window_offset');
                var from = to - get(this, 'timenav_window') - get(this, 'time_window_offset');

                console.log('refresh:', from, to, replace);

                this.updateAxisLimits(from, to);

                console.group('Load stylized series:');
                this.fetchStylizedSeries(from, to, replace);
                console.groupEnd();

                console.group('Load stylized metrics:');
                this.fetchStylizedMetrics(from, to, replace);
                console.groupEnd();

                me.trigger('onRenderer');
            },

            updateAxisLimits: function(from, to) {
                /* update axis limits */
                if(!get(this, 'zooming')) {
                    var opts = {};
                    $.extend(opts, get(this, 'chartOptions'));
                    $.extend(opts, {
                        xaxis: {
                            min: to - get(this, 'time_window') - get(this, 'time_window_offset'),
                            max: to
                        }
                    });

                    set(this, 'chartOptions', opts);

                    if(get(this, 'timenav')) {
                        opts = {};
                        $.extend(opts, get(this, 'timenavOptions'));
                        $.extend(opts, {
                            xaxis: {
                                min: to - get(this, 'timenav_window') - get(this, 'time_window_offset'),
                                max: to
                            }
                        });

                        set(this, 'timenavOptions', opts);
                    }
                }
            },

            fetchStylizedSeries: function(from, to, replace) {
                var store = get(this, 'widgetDataStore');

                /* fetch stylized series */
                var stylizedseries = get(this, 'config.series');
                var series = {};
                var curveIds = [];

                for(var i = 0, l = stylizedseries.length; i < l; i++) {
                    var serieId = get(stylizedseries[i], 'serie');

                    series[serieId] = {
                        style: stylizedseries[i],
                        serie: undefined,
                        curve: undefined
                    };

                    curveIds.push(get(stylizedseries[i], 'curve'));
                }

                var serieIds = JSON.stringify(Object.keys(series));
                curveIds = JSON.stringify(curveIds);

                console.log('series:', serieIds);
                console.log('curves:', curveIds);

                /* load series configuration */
                var me = this;

                Ember.RSVP.all([
                    store.findQuery('serie', {ids: serieIds}),
                    store.findQuery('curve', {ids: curveIds})
                ]).then(function(pargs) {
                    me.genChartConfig(pargs, series, from, to, replace);
                });
            },

            fetchStylizedMetrics: function(from, to, replace) {

                var store = get(this, 'widgetDataStore'),
                    stylizedmetrics = get(this, 'config.metrics'),
                    series = [],
                    seriesById = {},
                    curveIds = [],
                    me = this;

                var contextMetric = ['type', 'connector','connector_name', 'component','resource', 'metric'],
                    namedMetrics = [],
                    j,
                    tmpl;

                /**
                 * Compile template for metric labels with the context given above (contextMetric)
                 * @method tmpl
                 * @param {object} metric Related metric to template with
                 * @return {array} label
                 */
                tmpl = function(metric) {
                    console.log('metric tmpl', metric);
                    //debugger;
                    var serie = metric.id;
                    var seriesInfo = serie.split('/'),
                        templateContext = {};
                    var lengthSeriesInfo = seriesInfo.length;

                    //Build template context
                    for (j=0; j<lengthSeriesInfo; j++) {
                        //+1 is for preceding /
                        templateContext[contextMetric[j]] = seriesInfo[j + 1];
                    }

                    var template = get(me, 'config.metric_template');

                    if (template === "" || isNone(template)) {
                        template = "{{component}} - {{resource}} - {{metric}}";
                    }

                    try {
                        serie = Handlebars.compile(template)(templateContext);
                    } catch (err) {
                        console.log('could not proceed template feed', err);
                    }

                    var label = namedMetrics.pushObject(serie);
                    return label;
                };

                var fetchDone = function(curveIds) {
                    curveIds = JSON.stringify(curveIds);

                    console.log('series:', seriesById);
                    console.log('curves:', curveIds);

                    store.findQuery('curve', {ids: curveIds}).then(function(curveResult) {
                        var virtualResult = {
                            meta: {
                                total: series.length
                            },
                            content: series,
                        };

                        me.genChartConfig([virtualResult, curveResult], seriesById, from, to, replace);
                    });
                };

                if(stylizedmetrics.length > 0) {
                    var metricIds = [];

                    for(var i = 0, l = stylizedmetrics.length ; i < l ; i++) {
                        var metricId = get(stylizedmetrics[i], 'metric');

                        metricIds.push(metricId);
                    }

                    store.findQuery('ctxmetric', {
                        filter: JSON.stringify({
                            _id: {
                                '$in': metricIds
                            }
                        })
                    }).then(function(result) {
                        var metricsById = {},
                            i, l;

                        for(i = 0, l = get(result, 'content.length') ; i < l; i++) {
                            var info = get(result, 'content')[i];

                            metricsById[get(info, 'id')] = info;
                        }

                        for(i = 0, l = stylizedmetrics.length; i < l; i++) {
                            var metricId = get(stylizedmetrics[i], 'metric');
                            var metricInfo = metricsById[metricId];
                            console.log('mectric info', metricInfo);
                            //debugger;
                            var serieconf = Ember.Object.create({
                                id: metricId,
                                virtual: true,
                                crecord_name: tmpl(metricInfo),
                                metrics: [metricId],
                                aggregate_method: 'none',
                                unit: get(stylizedmetrics[i], 'unit')
                            });

                            seriesById[metricId] = {
                                style: stylizedmetrics[i],
                                serie: serieconf,
                                curve: undefined
                            };

                            series.push(serieconf);
                            curveIds.push(get(stylizedmetrics[i], 'curve'));
                        }

                        fetchDone(curveIds);
                    });
                }
                else {
                    fetchDone(curveIds);
                }
            },

            genChartConfig: function(pargs, series, from, to, replace) {
                console.group('Generate Chart series');

                var serieResult = pargs[0]; // arguments of first promise
                var curveResult = pargs[1]; // arguments of second promise

                var i, l;

                console.group('Fetch curves:');
                var curvesById = {};

                for(i = 0, l = curveResult.meta.total; i < l; i++) {
                    var curve = curveResult.content[i];
                    curvesById[curve.id] = curve;
                }

                console.log(curvesById);
                console.groupEnd();

                console.group('Fetch series:');
                for(i = 0, l = serieResult.meta.total; i < l; i++) {
                    var serieconf = serieResult.content[i];
                    var serieId = serieconf.id;

                    console.log(serieconf, serieId);

                    if(series[serieId] !== undefined) {
                        var config = series[serieId];
                        var curveId = get(config, 'style.curve');
                        var curveconf = curvesById[curveId];

                        if(curveconf !== undefined) {
                            set(config, 'curve', curveconf);
                        }

                        console.log(curveId, curveconf);

                        set(config, 'serie', serieconf);
                        this.genChartSerie(config, from, to, replace);
                    }
                }

                console.log('stylizedseries:', series);
                console.groupEnd();

                console.groupEnd();
            },

            genChartSerie: function(config, from, to, replace) {
                console.group('Generating Chart serie:', config);

                var me = this;

                var fillcolor = null;
                var fillopacity = false;

                if (get(config, 'curve.areas')) {
                    var spec = $.color.parse(get(config, 'style.color'));

                    fillopacity = get(config, 'curve.area_opacity');
                    spec.a = fillopacity;

                    fillcolor = spec.toString();
                }

                var chartSerie = {
                    label: get(config, 'serie.crecord_name'),
                    color: get(config, 'style.color'),
                    lines: {
                        show: get(config, 'curve.lines') || get(config, 'curve.areas'),
                        used: get(config, 'curve.lines'),
                        lineWidth: get(config, 'curve.line_width'),
                        fill: fillopacity,
                        fillColor: fillcolor
                    },
                    bars: {
                        show: get(config, 'curve.bars'),
                        used: get(config, 'curve.bars'),
                        lineWidth: get(config, 'curve.line_width'),
                        barWidth: get(config, 'curve.bar_width') * 1000,
                        fill: fillopacity,
                        fillColor: fillcolor
                    },
                    points: {
                        show: get(config, 'curve.points'),
                        used: get(config, 'curve.points'),
                        symbol: get(config, 'curve.point_shape')
                    },
                    values: {
                        show: get(config, 'curve.valueLabels'),
                        labelFormatter: function(series, text) {
                            var n = parseFloat(text);
                            return n.toFixed(2);
                        }
                    },
                    xaxis: parseInt(get(config, 'style.xaxis')) || 1,
                    yaxis: parseInt(get(config, 'style.yaxis')) || 1,
                    clickable: true,
                    hoverable: true,

                    hidden: false,
                    config: {
                        color: get(config, 'style.color'),
                        lines: get(config, 'curve.lines') || get(config, 'curve.areas'),
                        bars: get(config, 'curve.bars'),
                        points: get(config, 'curve.points'),
                        values: get(config, 'curve.valueLabels')
                    },
                    unit: get(config, 'serie.unit')
                };

                var oldSerieId, oldSerie, ctrl, request;

                if(get(config, 'serie.virtual') === true) {
                    ctrl = get(this, 'controllers.perfdata');
                    oldSerieId = get(config, 'style.metric');
                    oldSerie = get(this, 'chartSeries')[oldSerieId];
                    request = get(config, 'serie.id');
                }
                else {
                    ctrl = get(this, 'controllers.serie');
                    oldSerieId = get(config, 'style.serie');
                    oldSerie = get(this, 'chartSeries.' + oldSerieId);
                    request = get(config, 'serie');
                }

                if(oldSerie !== undefined && !replace) {
                    chartSerie.data = oldSerie.data;
                }
                else {
                    chartSerie.data = [];
                }

                chartSerie.serie = oldSerieId;

                console.log('chartserie:', chartSerie);
                console.log('Fetch perfdata and compute serie');

                var aggregation = (
                    (get(config, 'serie.aggregate_method') !== 'none')
                    ||
                    (get(config, 'serie.metrics.length') > 1)
                );

                var min, max;

                if(get(this, 'timenav')) {
                    min = get(this, 'timenavOptions.xaxis.min');
                    max = get(this, 'timenavOptions.xaxis.max');
                }
                else {
                    min = get(this, 'chartOptions.xaxis.min');
                    max = get(this, 'chartOptions.xaxis.max');
                }

                if(aggregation) {
                    from = min;
                    to = max;
                    chartSerie.data = [];
                }

                console.log('call controller fetch', request);

                ctrl.fetch(request, from, to).then(function(data) {
                    console.log('fetch:', data);

                    if(get(config, 'serie.virtual') === true) {
                        var chartSeries = get(me, 'chartSeries');
                        var points = data.data[0].points;

                        for(var i = 0, l = points.length; i < l; i++) {
                            points[i][0] = points[i][0] * 1000;
                        }

                        chartSerie.data = chartSerie.data.concat(points);
                        chartSeries[get(config, 'style.metric')] = chartSerie;

                        set(me, 'chartSeries', chartSeries);
                    }
                    else {
                        chartSerie.data = chartSerie.data.concat(data);

                        set(me, 'chartSeries.' + get(config, 'style.serie'), chartSerie);
                    }

                    me.recomputeDataSeries();
                });

                console.groupEnd();
            },

            recomputeDataSeries: function() {
                var chartSeries = get(this, 'chartSeries');
                var oldSeries = get(this, 'dataSeries');

                var i, l, serieId;

                for(i = 0, l = oldSeries.length; i < l; i++) {
                    var oldSerie = oldSeries[i];
                    serieId = get(oldSerie, 'serie');
                    var chartSerie = get(chartSeries, serieId);

                    if (chartSerie !== undefined) {
                        var hidden = get(oldSerie, 'hidden');
                        set(chartSerie, 'hidden', hidden);

                        if (hidden) {
                            set(chartSerie, 'color', '#CCCCCC');
                            set(chartSerie, 'lines.show', false);
                            set(chartSerie, 'bars.show', false);
                            set(chartSerie, 'points.show', false);
                            set(chartSerie, 'values.show', false);
                        }
                    }
                }

                delete oldSeries;

                var series = Ember.A();
                var serieIds = Object.keys(chartSeries);

                for(i = 0, l = serieIds.length; i < l; i++) {
                    serieId = serieIds[i];

                    series.pushObject(chartSeries[serieId]);
                }

                console.log('dataSeries:', series);

                set(this, 'dataSeries', series);
            }
        }, widgetOptions);

        application.register('widget:timegraph', widget);
    }
});
