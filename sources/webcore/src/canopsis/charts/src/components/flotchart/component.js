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
    name: 'component-flotchart',
    after: 'ValuesUtils',
    initialize: function(container, application) {
        var values = container.lookupFactory('utility:values');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var component = Ember.Component.extend({
            tagName: 'div',
            classNames: 'flotchart',

            init: function() {
                set(this, 'chart', null);

                if (isNone(get(this, 'options'))) {
                    set(this, 'options', null);
                }

                if (isNone(get(this, 'series'))) {
                    set(this, 'series', []);
                }

                if (isNone(get(this, 'human'))) {
                    set(this, 'human', false);
                }

                this._super.apply(this, arguments);

                this.addObserver('series.@each', this.onDataUpdate);
                this.addObserver('options', this.onOptionsUpdate);
            },

            onDataUpdate: function() {
                this.send('renderChart');
            },

            onOptionsUpdate: function() {
                this.send('renderChart');
            },

            didInsertElement: function() {
                this.send('renderChart');
            },

            actions: {
                renderChart: function() {
                    var chart = get(this, 'chart');

                    if(chart !== null) {
                        console.log('Destroy chart');
                        chart.destroy();
                    }

                    this.createChart();
                },

                toggleSerie: function(config) {
                    var label = config[0], serieIndex = config[1];

                    var series = get(this, 'series');
                    var serie = series[serieIndex];

                    console.log('Toggling serie:', label, serieIndex, serie);

                    if (serie.config.lines && serie.lines) {
                        serie.lines.show = !serie.lines.show;
                    }

                    if (serie.config.bars && serie.bars) {
                        serie.bars.show = !serie.bars.show;
                    }

                    if (serie.config.points && serie.points) {
                        serie.points.show = !serie.points.show;
                    }

                    if (serie.config.values && serie.values) {
                        serie.values.show = !serie.values.show;
                    }

                    serie.hidden = !serie.hidden;

                    if (serie.hidden) {
                        serie.color = '#CCCCCC';
                    }
                    else {
                        serie.color = serie.config.color;
                    }

                    this.send('renderChart');

                    // trigger event
                    this.$().trigger('toggleserie', [config]);
                }
            },

            classifiedSeries: function() {
                var series = get(this, 'series'),
                    seriesByAxis = {
                    x: {},
                    y: {}
                };

                for(var i = 0, l = series.length; i < l; i++) {
                    var serie = series[i];

                    serie.index = i;

                    if(seriesByAxis.x[serie.xaxis] === undefined) {
                        seriesByAxis.x[serie.xaxis] = [];
                    }

                    seriesByAxis.x[serie.xaxis].push(serie);

                    if(seriesByAxis.y[serie.yaxis] === undefined) {
                        seriesByAxis.y[serie.yaxis] = [];
                    }

                    seriesByAxis.y[serie.yaxis].push(serie);
                }

                return seriesByAxis;
            }.property('series.@each'),

            createAxes: function() {
                var options = get(this, 'options'),
                    seriesByAxis = get(this, 'classifiedSeries');

                var axis, l,
                    xaxes = Object.keys(seriesByAxis.x),
                    yaxes = Object.keys(seriesByAxis.y);

                for(axis = 0, l = xaxes.length ; axis < l ; axis++) {
                    var key = xaxes[axis];
                    var idx = key - 1;

                    if (options.xaxes[idx] === undefined) {
                        options.xaxes[idx] = {
                            show: true,
                            reserveSpace: true,
                            position: (idx % 2 === 0 ? 'bottom' : 'top'),
                            color: seriesByAxis.x[key][0].color,
                            font: {
                                color: seriesByAxis.x[key][0].color
                            }
                        };
                    }
                }

                var me = this;

                for(axis = 0, l = yaxes.length ; axis < l ; axis++) {
                    var key = yaxes[axis];
                    var idx = key - 1;

                    if (options.yaxes[idx] === undefined) {
                        options.yaxes[idx] = {
                            show: true,
                            reserveSpace: true,
                            position: (idx % 2 == 0 ? 'left' : 'right'),
                            color: seriesByAxis.y[key][0].color,
                            font: {
                                color: seriesByAxis.y[key][0].color
                            }
                        };
                    }

                    options.yaxes[idx].tickFormatter = function(val, axis) {
                        if (get(me, 'human') === true) {
                            var unit = seriesByAxis.y[axis.n][0].unit || '';
                            val = values.humanize(val, unit);
                        }

                        return val;
                    };
                }
            },

            createLegend: function() {
                var options = get(this, 'options'),
                    me = this;

                if(options && options.legend && options.legend.show && options.legend.labelFormatter === undefined) {
                    options.legend.labelFormatter = function(label, serie) {
                        console.log('Format label for serie:', label, serie);

                        var style = 'cursor: pointer; margin-left: 5px;';

                        if(serie.hidden) {
                            style += ' font-style: italic;';
                        }

                        var clickableLabel = $('<span/>', {
                            style: style,
                            onclick: [
                                'Ember.View.views["', me.elementId, '"].send("toggleSerie", [',
                                    '"', label, '",',
                                    serie.index,
                                ']);'
                            ].join(''),
                            text: label
                        });

                        var tmpContainer = $('<div/>');
                        tmpContainer.append(clickableLabel);

                        return tmpContainer.html();
                    };
                }
            },

            createChart: function() {
                console.group('createChart:');

                var plotcontainer = this.$(),
                    series = get(this, 'series'),
                    options = get(this, 'options');

                console.log('options (before):', options);

                if (get(this, 'options') !== null) {
                    this.createAxes();
                    this.createLegend();
                    this.autoscale();

                    /* create plot */
                    console.log('series:', series);
                    console.log('options:', options);

                    set(this, 'chart', $.plot(plotcontainer, series, options));
                }

                console.groupEnd();
            },

            autoscale: function() {
                var seriesByAxis = get(this, 'classifiedSeries'),
                    options = get(this, 'options');

                var yaxes = Object.keys(seriesByAxis.y);

                for(var axis = 0, l = yaxes.length; axis < l; axis++) {
                    var key = yaxes[axis];
                    var idx = key - 1;

                    var n_series = seriesByAxis.y[key].length,
                        min = null, max = null, valuesOnChart = false;

                    for(var serieidx = 0; serieidx < n_series; serieidx++) {
                        var serie = seriesByAxis.y[key][serieidx];

                        if (!serie.hidden) {
                            var boundaries = this.getSerieBoundaries(serie);

                            if (min === null || boundaries[0] < min) {
                                min = boundaries[0];
                            }

                            if (max === null || boundaries[1] > max) {
                                max = boundaries[1];
                            }
                        }

                        if(serie.values.show) {
                            valuesOnChart = true;
                        }
                    }

                    /* calculate new max with margin */
                    var margin = (valuesOnChart ? 115.0 : 105.0);

                    var inc = max + 30;
                    var incperc = (inc * 100.0) / max;
                    var dec;

                    if (incperc > margin) {
                        inc = (max * margin) / 100.0;
                    }

                    /* calculate new min with margin */
                    if(min < 0 || min > 30) {
                        min = -min;
                        dec = min + 30;
                        var decperc = (min * 100.0) / min;

                        if (decperc > margin) {
                            dec = (min * margin) / 100.0;
                        }

                        dec = -dec;
                    }
                    else {
                        dec = 0;
                    }

                    options.yaxes[idx].min = dec;
                    options.yaxes[idx].max = inc;
                }
            },

            refreshChart: function() {
                var chart = get(this, 'chart'),
                    series = get(this, 'series');

                console.log('flotchart refresh series:', series);
                this.autoscale();
                chart.setData(series);
                chart.setupGrid();
                chart.draw();
            },

            getSerieBoundaries: function(serie) {
                var min = null, max = null, options = get(this, 'options');

                for(var i = 0, l = serie.data.length; i < l; i++) {
                    var point = serie.data[i];

                    /* skip points outside timewindow */
                    if (point[0] > options.xaxis.max) {
                        break;
                    }
                    else if (point[0] >= options.xaxis.min) {
                        if (min === null || point[1] < min) {
                            min = point[1];
                        }

                        if (max === null || point[1] > max) {
                            max = point[1];
                        }
                    }
                }

                return [min, max];
            }
        });

        application.register('component:component-flotchart', component);
    }
});
