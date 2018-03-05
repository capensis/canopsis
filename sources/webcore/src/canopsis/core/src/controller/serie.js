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
    name: 'SerieController',
    after: 'HashUtils',
    initialize: function(container, application) {
        var hash = container.lookupFactory('utility:hash');
        var math = window.math;

        var get = Ember.get,
            set = Ember.set;


        var controller = Ember.ObjectController.extend({
            needs: ['perfdata'],

            perfdata: Ember.computed.alias('controllers.perfdata'),

            fetch: function(serie, from, to) {
                if(get(serie, 'metrics.length') > 1 && get(serie, 'aggregate_method') === 'none') {
                    console.group('More than one metric in serie, performing an aggregation');

                    console.log('serie:', serie);
                    console.log('aggregation: average - 60s');

                    set(serie, 'aggregate_method', 'average');
                    set(serie, 'aggregate_interval', 60);

                    console.groupEnd();
                }

                if(get(serie, 'aggregate_method') === 'none') {
                    var promise = get(this, 'perfdata').fetchMany(
                        get(serie, 'metrics'),
                        from, to
                    );
                }
                else {
                    var promise = get(this, 'perfdata').aggregateMany(
                        get(serie, 'metrics'),
                        from, to,
                        get(serie, 'aggregate_method'),
                        get(serie, 'aggregate_interval')
                    );
                }

                var formula = get(serie, 'formula');

                return promise.then(function(result) {
                    console.group('Computing serie:', formula);

                    var nmetric = result.total;
                    var metrics = result.data;

                    // build points dictionnary
                    var points = {};
                    var length = false;
                    var i;

                    for(i = 0; i < nmetric; i++) {
                        var metric = metrics[i];

                        var id = metric.meta.data_id;
                        var mid = 'metric_' + hash.md5(id);
                        var mname = '/' + id.split('/').slice(4).join('/');

                        // replace metric name in formula by the unique id
                        formula = formula.replaceAll(mname, mid);

                        console.log('metric:', mname, mid);

                        points[mid] = metric.points;

                        /* make sure we treat the same amount of points by selecting
                         * the metric with less points.
                         */
                        if(!length || metric.points.length < length) {
                            length = metric.points.length;
                        }
                    }

                    console.log('formula:', formula);
                    console.log('points:', points);

                    var mids = Object.keys(points);
                    var finalSerie = [];

                    // now loop over all points to calculate the final serie
                    for(i = 0; i < length; i++) {
                        var data = {};
                        var ts = 0;
                        var j, l;

                        for(j = 0, l = mids.length; j < l; j++) {
                            var mid = mids[j];

                            // get point value at timestamp "i" for metric "mid"
                            data[mid] = points[mid][i][1];

                            // set timestamp
                            ts = points[mid][i][0];
                        }

                        // import data in math context
                        math.import(data);
                        var pointval = math.eval(formula);

                        // remove data from math context
                        for(j = 0, l = mids.length; j < l; j++) {
                            var mid = mids[j];
                            delete math[mid];
                        }

                        // push computed point in serie
                        finalSerie.push([ts * 1000, pointval]);
                    }

                    console.log('finalserie:', finalSerie);
                    console.groupEnd();

                    return finalSerie;
                });
            }
        });

        application.register('controller:serie', controller);
    }
});
