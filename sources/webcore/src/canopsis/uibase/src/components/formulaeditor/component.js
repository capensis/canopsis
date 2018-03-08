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
    name: 'component-formulaeditor',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        /**
         * @component formulaeditor
         */
        var component = Ember.Component.extend({
            form: Ember.computed.alias('parentView.parentView.controller'),

            metrics_cols: [
                {name: 'component', title: __('Component')},
                {name: 'resource', title: __('Resource')},
                {name: 'id', title: __('Metric')},
                {
                    action: 'insert',
                    actionAll: undefined,
                    title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-plus-sign"></span>'),
                    style: 'text-align: center;'
                }
            ],

            init: function() {
                set(this, 'selectedMetrics', Ember.A());

                this._super.apply(this, arguments);

                var categories = get(this, 'form.categories'),
                    found = false;

                /* look for the field 'metrics' */
                for(var i = 0, l = categories.length; i < l; i++) {
                    var category = categories[i];

                    for(var j = 0, nkeys = category.keys.length; j < nkeys; j++) {
                        var attr = category.keys[j];

                        if (attr.field === 'metrics') {
                            /* notify when user modify this value in another editor */
                            Ember.addObserver(attr, 'value', this, this.metricDidChange);
                            found = true;
                            break;
                        }
                    }

                    if (found) {
                        break;
                    }
                }
            },

            metricDidChange: function(attr, field) {
                var metrics = get(attr, field);
                var metas = [];

                /* convert metric ids to metric metadata */
                for(var i = 0, l = metrics.length; i < l; i++) {
                    /* /metric/<connector>/<connector name>/<component>/<metric>
                     * or
                     * /metric/<connector>/<connector name>/<component>/<resource>/<metric>
                     */

                    var metric = metrics[i].split('/');
                    var nmeta = metric.length;

                    if (nmeta === 7) {
                        metric = {
                            component: metric[nmeta - 3],
                            resource: metric[nmeta - 2],
                            id: metric[nmeta - 1]
                        };
                    }
                    else {
                        metric = {
                            component: metric[nmeta - 2],
                            resource: null,
                            id: metric[nmeta - 1]
                        };
                    }

                    metas.push(metric);
                }

                set(this, 'selectedMetrics', metas);
            },

            actions: {
                insert: function(metric) {
                    var add_to_formula = ' ' + [
                        '',
                        get(metric, 'component'),
                        get(metric, 'resource') || '',
                        get(metric, 'id')
                    ].join('/');

                    var formula = get(this, 'content.value');

                    if (isNone(formula)) {
                        set(this, 'content.value', add_to_formula);
                    }
                    else {
                        set(this, 'content.value', formula + add_to_formula);
                    }
                }
            }
        });

        application.register('component:component-formulaeditor', component);
    }
});
