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

Ember.Application.initializer({
    name:"component-c3categorychart",
    after: ['HashUtils', 'ValuesUtils'],
    initialize: function(container, application) {

    var HashUtils = container.lookupFactory('utility:hash');
    var ValuesUtils = container.lookupFactory('utility:values');


    var get = Ember.get,
        set = Ember.set,
        isNone = Ember.isNone,
        __ = Ember.String.loc;

    /**
     * @description Component for instantiate C3 charts
     * @component c3categorychart
     * @example
     * {{#if parentController.options.allow_user_display}}
     *   {{#component-dropdownbutton}}
     *       {{#component-dropdownbuttonheader}}
     *           {{#component-dropdownbuttontitle}}
     *               {{#component-tooltip content="Display"}}
     *                   <i class="fa fa-eye"></i>
     *               {{/component-tooltip}}
     *           {{/component-dropdownbuttontitle}}
     *       {{/component-dropdownbuttonheader}}
     *       {{#component-dropdownbuttoncontent}}
     *           <ul class="list-group">
     *               <li class="list-group-item pointer" {{action "transform" "gauge"}}>
     *                   <i class="fa fa-tachometer"></i> {{tr "gauge"}}
     *               </li>
     *               <li class="list-group-item pointer" {{action "transform" "pie"}}>
     *                   <i class="fa fa-pie-chart"></i> {{tr "pie"}}
     *               </li>
     *               <li class="list-group-item pointer" {{action "transform" "donut"}}>
     *                   <i class="fa fa-dot-circle-o"></i> {{tr "donut"}}
     *               </li>
     *               <li class="list-group-item pointer" {{action "transform" "progressbar"}}>
     *                   <i class="fa fa-battery-three-quarters"></i> {{tr "progressbar"}}
     *               </li>
     *               <li class="list-group-item pointer" {{action "transform" "bar"}}>
     *                   <i class="fa fa-bar-chart"></i> {{tr "bar"}}
     *               </li>
     *           </ul>
     *       {{/component-dropdownbuttoncontent}}
     *   {{/component-dropdownbutton}}
     *   {{/if}}
     *
     *   <div {{bind-attr id="uuid"}}></div>
     *
     *   {{#unless chart}}
     *       <h2 class="text-center">No data</h2>
     *   {{/unless}}
     */
    var component = Ember.Component.extend({
        /**
         * @description instantiate component
         * @method init
         */
        init: function() {
            this._super();
            Ember.setProperties(this, {
                'uuid': HashUtils.generateId('categoryChart'),
                'parentController.chartComponent': this,
            });
        },

        /**
         * @description Destroy each event handled before in the component
         * @method willDestroyElement
         */
        willDestroyElement: function() {
            var chart = get(this, 'chart');
            if (!isNone(chart)) {
                chart.destroy();
            }
        },

        /**
         * @description Define correctly the maxValue
         * @method maxValue
         * @returns {Integer} maxValue
         */
        maxValue: function () {

            /**
            User max value may be defined
            if it is defined, it must be greater than series values sum
            max value is used only if configuration allowed it
            **/

            var maxValue = get(this, 'parentController.options.max_value'),
                useMaxValue = get(this, 'parentController.options.use_max_value'),
                seriesSum = get(this, 'seriesSum');

            if (useMaxValue && !isNone(maxValue) && maxValue > seriesSum) {
                return maxValue;
            } else {
                return seriesSum;
            }

        }.property('parentController.options.max_value'),

        /**
         * @description Compute all series values sum
         * @method seriesSum
         * @returns {Integer} sum
         */
        seriesSum: function () {

            var sum = 0,
                series = get(this, 'seriesWithComputedNames');
            for (var i=0; i<series.length; i++) {
                sum += series[i][1];
            }
            return sum;

        }.property('series'),

        /**
         * @description Get the list of each distinct serie name.
         * @method seriesNames
         * @returns {Array} seriesNames
         */
        seriesNames: function () {

            var seriesNames = [];
            var series = get(this, 'c3series');
            var length = series.length;
            for (var i=0; i<length; i++) {
                seriesNames.push(series[i][0]);
            }
            return seriesNames;

        }.property('series'),


        /**
         * @description Generate a new array for series with metric names computed from user template.
         * @method seriesWithComputedNames
         * @returns {Array} namedSeries
         */
        seriesWithComputedNames: function () {

            var context = ['type', 'connector','connector_name', 'component','resource', 'metric'],
                seriesWithMeta = $.extend(true, [], get(this, 'series')),
                namedSeries = [],
                i,
                j;

            console.log('seriesWithMeta', seriesWithMeta);

            var length = seriesWithMeta.length;

            for (i=0; i<length; i++) {
                var serie = seriesWithMeta[i].serie,
                    id = seriesWithMeta[i].id;

                console.log('meta serie', id, serie);

                var seriesInfo = id.split('/'),
                    templateContext = {};
                var lengthSeriesInfo = seriesInfo.length;

                //Build template context
                for (j=0; j<lengthSeriesInfo; j++) {
                    //+1 is for preceding /
                    templateContext[context[j]] = seriesInfo[j + 1];
                }
                console.log('Template context', templateContext, 'for metric', id);

                var template = get(this, 'parentController.options.metric_template');

                try {
                    serie[0] = Handlebars.compile(template)(templateContext);
                } catch (err) {
                    console.log('could not proceed template feed', err);
                }


                namedSeries.push(serie);

            }

            return namedSeries;

        }.property(),

        /**
         * @description Generate chart required values to be displayed
         * @method c3series
         * @returns {Array} series
         */
        c3series: function () {

            console.log('chart series is now', get(this, 'series'));

            var restValue = get(this, 'maxValue') - get(this, 'seriesSum'),
                series = get(this, 'seriesWithComputedNames');
                //base series data deep copied

            //Compute difference between max value and series values sum
            if (restValue > 0) {
                var leftValueLabel = get(this, 'parentController.options.text_left_space');
                series.push([leftValueLabel, restValue]);
            }

            //Sort series for clean display
            series.sort(function(a, b) {
                return b[1] - a[1];
            });

            return series;

        }.property('series'),

        /**
         * @description Compute the color dict for nice chart display from options
         * @method colors
         * @returns {Object} colors
         */
        colors: function () {

            var seriesNames = get(this, 'seriesNames');

            var colors = {
                leftValueLabel: '#EEEEEE',
            };

            for (var i=0; i<seriesNames.length; i++) {
                colors[seriesNames[i]] = ['#FF0000', '#F97600'][i];
            }
            return colors;

        }.property('seriesNames'),

        /**
         * @description Tells the component is inserted into the dom
         * @method didInsertElement
         */
        didInsertElement: function () {

            set(this, 'domready', true);
        },

        /**
         * @decription Uses series and chart options to insert a C3js chart element in the dom
         * @method generateChart
         */
        generateChart: function () {

            if(isNone(get(this, 'domready'))) {
                console.log('Dom is not ready for category chart, cannot draw');
                return;
            }

            if(isNone(get(this, 'parentController.options'))) {
                console.log('Chart options are not ready cannot draw');
                return;
            }


            var domElement = '#' + get(this, 'uuid'),
                seriesSum = get(this, 'seriesSum'),
                seriesNames = get(this, 'seriesNames'),
                c3series = get(this, 'c3series'),
                colors = get(this, 'colors'),
                maxValue = get(this, 'maxValue'),
                leftValueLabel = get(this, 'parentController.options.text_left_space'),
                chartType = get(this, 'parentController.options.display'),
                showLegend = get(this, 'parentController.options.show_legend'),
                tooltip = get(this, 'parentController.options.show_tooltip'),
                showLabels = get(this, 'parentController.options.show_labels'),
                stacked = get(this, 'parentController.options.stacked'),
                humanReadable = get(this, 'parentController.options.human_readable'),
                rotated = false,
                showAxes = true,
                isBarChart = true;

            seriesSum = humanReadable ? ValuesUtils.humanize(seriesSum, ''): seriesSum.toFixed(2);

            var label = {
                show : showLabels,
                format: function(value, ratio){
                    return  showLabels ? seriesSum: '';
                }
            };

            var gauge = {
                    label:label
                },
                pie = {},
                donut = {
                    title: (showLabels && seriesSum) ? seriesSum: ''
                };

            console.log('seriesNames', seriesNames);
            console.log('c3series', c3series);

            if (chartType == 'progressbar') {
                //cheating alias becrause progress bar does not exists in c3 js yet
                isBarChart = false;
                chartType = 'bar';
                rotated = true;
                showAxes = false;
                showLabels = false;
            }

            //max value may be equal to 0 when series did not fetch points.
            if (maxValue > 0) {
                //define the max value of the chart and wether or not a delta serie is created
                gauge.max = maxValue;
                if (maxValue > seriesSum && $.inArray(leftValueLabel, seriesNames) === -1) {
                    seriesNames.push(leftValueLabel);
                }
            }

            var options = {
                bindto: domElement,
                groups: seriesNames,
                tooltip: {show: tooltip},
                legend: {show: showLegend},
                data: {
                    columns: c3series,
                    type: chartType,
                    groups: [seriesNames],
                    labels: {
                        format: function (v, id, i, j) {
                            v = humanReadable ? ValuesUtils.humanize(v, '') : parseFloat(v).toFixed(2);
                            return showLabels ? id + ' : ' + v : '';
                        }
                    },
                    empty: {
                        label: {
                            text: __('No Data')
                        }
                    }
                },
                //color: colors,
                gauge: gauge,
                donut: donut,
                pie: pie,
                axis: { //for bar mode
                  rotated:rotated,
                  x: {
                    show: showAxes
                  },
                  y: {
                    tick: {
                        format: function (v) {
                            return humanReadable ? ValuesUtils.humanize(v, '') : parseFloat(v).toFixed(2);
                        }
                    },
                    show: showAxes
                  }
                },
            };

            if (chartType === 'bar' && !stacked) {
                //no more stacked view bor bar charts
                delete options.groups;
                delete options.data.groups;
            }

            var chart = c3.generate(options);

            set(this, 'chart', chart);

        },

        /**
         * @description Update the chart display with new values. Insert a new chart if it does not exists yet.
         * @method update
         */
        update: function () {

            var chart = get(this, 'chart');

            if (isNone(chart)) {
                this.generateChart();
            } else {
                console.log('refreshing c3 chart with series', get(this, 'c3series'));

                var previousSeriesNames = get(this, 'seriesNames');

                if (previousSeriesNames.length) {
                    chart.unload({
                        ids: previousSeriesNames
                    });
                }

                chart.load({
                    columns: get(this, 'c3series'),
                    groups: [get(this, 'seriesNames')]
                });
            }


        }.observes('series', 'ready', 'parentController.options'),

        actions: {
            /**
             * @description Remove the actual chart to replace it by a new one
             * @method actions_transform
             */
            transform: function (type) {
                get(this, 'chart').destroy();
                set(this, 'parentController.options.display', type);
                this.generateChart();
            }

        }

    });

    application.register('component:component-c3categorychart', component);

    }
});
