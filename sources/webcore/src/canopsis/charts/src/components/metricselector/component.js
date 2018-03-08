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
    name: 'component-metricselector',
    after: 'HashUtils',
    initialize: function(container, application) {
        var hash = container.lookupFactory('utility:hash');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone,
            __ = Ember.String.loc;

        var component = Ember.Component.extend({
            valueChanged: function() {
                var metrics = get(this, 'selectedMetrics') || [];
                var ids = [];

                for(var i = 0, l = metrics.length; i < l; i++) {
                    var m_id = get(metrics[i], 'id');
                    ids.push(m_id);
                }

                if (get(this, 'multiselect') === true) {
                    set(this, 'content', ids);
                }
                else {
                    set(this, 'content', ids[0]);
                }
            }.observes('selectedMetrics.@each'),

            helpModal: {
                title: __('Syntax'),
                content: ['<ul>',
                    '<li><code>co:regex</code> : ', __('look for a component'), '</li>',
                    '<li><code>re:regex</code> : ', __('look for a resource'), '</li>' ,
                    '<li><code>me:regex</code> : ' , __('look for a metric') , '(<code>me:</code>' , __(' isn\'t needed for this one') , ')</li>' ,
                    '<li>', __('combine all of them to improve your search'),' : <code>co:regex re:regex me:regex</code></li>' ,
                    '<li><code>co:</code>, <code>re:</code>, <code>me:</code> : ', __('look for non-existant field') , '</li>' ,
                    '</ul>'].join(''),

                id: hash.generateId('cmetric-help-modal'),
                label: hash.generateId('cmetric-help-modal-label')
            },

            select_cols: function() {
                return [
                    {name: 'component', title: __('Component')},
                    {name: 'resource', title: __('Resource')},
                    {name: 'name', title: __('Metric')},
                    {
                        action: 'select',
                        actionAll: (get(this, 'multiselect') === true ? 'selectAll' : undefined),
                        title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-plus-sign"></span>'),
                        style: 'text-align: center;'
                    }
                ];
            }.property(),

            unselect_cols: function() {
                return [
                    {name: 'component', title: __('Component')},
                    {name: 'resource', title: __('Resource')},
                    {name: 'name', title: __('Metric')},
                    {
                        action: 'unselect',
                        actionAll: 'unselectAll',
                        title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-trash"></span>'),
                        style: 'text-align: center;'
                    }
                ];
            }.property(),

            init: function() {
                if (isNone(get(this, 'selectedMetrics'))) {
                    set(this, 'selectedMetrics', []);
                }

                if (isNone(get(this, 'metricSearch'))) {
                    set(this, 'metricSearch', null);
                }

                if (isNone(get(this, 'multiselect'))) {
                    set(this, 'multiselect', true);
                }

                set(this, 'metrics', []);

                this._super.apply(this, arguments);

                var store = DS.Store.create({
                    container: this.get('container')
                });

                set(this, 'componentDataStore', store);

                var query = {filter: {_id: undefined}}, me = this;

                if (!isNone(get(this,'content'))) {
                    if (get(this, 'multiselect') === true) {
                        var content = get(this, 'content') || [];
                        query.filter._id = {'$in': content};
                    }
                    else {
                        query.filter._id = get(this, 'content');
                    }

                    store.findQuery('ctxmetric', query).then(function(result) {
                        var metrics = get(me, 'selectedMetrics') || [],
                            content = get(result, 'content'),
                            l = get(result, 'meta.total');

                        console.log('Received data:', l, content, metrics);

                        for(var i = 0; i < l; i++) {
                            metrics.pushObject(content[i]);
                        }

                        set(me, 'selectedMetrics', metrics);
                    });
                }
            },

            build_filter: function(search) {
                var conditions = search.split(' ');
                var i;

                var patterns = {
                    component: [],
                    resource: [],
                    name: []
                };

                for(i = 0; i < conditions.length; i++) {
                    var condition = conditions[i];

                    if(condition !== '') {
                        var regex = condition.slice(3) || null;

                        if(condition.indexOf('co:') === 0) {
                            patterns.component.push(regex);
                        }
                        else if(condition.indexOf('re:') === 0) {
                            patterns.resource.push(regex);
                        }
                        else if(condition.indexOf('me:') === 0) {
                            patterns.name.push(regex);
                        }
                        else {
                            patterns.name.push(condition);
                        }
                    }
                }

                var mfilter = {'$and': []};
                var filters = {
                    component: {'$or': []},
                    resource: {'$or': []},
                    name: {'$or': []}
                };

                for(var key in filters) {
                    for(i = 0; i < patterns[key].length; i++) {
                        var filter = {};
                        var value = patterns[key][i];

                        if(value !== null) {
                            filter[key] = {'$regex': value};
                        }
                        else {
                            filter[key] = null;
                        }

                        filters[key].$or.push(filter);
                    }

                    var len = filters[key].$or.length;

                    if(len === 1) {
                        filters[key] = filters[key].$or[0];
                    }

                    if(len > 0) {
                        mfilter.$and.push(filters[key]);
                    }
                }

                if(mfilter.$and.length === 1) {
                    mfilter = mfilter.$and[0];
                }

                return mfilter;
            },

            actions: {
                select: function(metric) {
                    var selected = get(this, 'selectedMetrics');

                    if (selected.indexOf(metric) < 0) {
                        console.log('Select metric:', metric);

                        if (get(this, 'multiselect') === true) {
                            selected.pushObject(metric);
                        }
                        else {
                            selected = [metric];
                        }
                    }

                    set(this, 'selectedMetrics', selected);
                },

                unselect: function(metric) {
                    var selected = get(this, 'selectedMetrics');

                    var idx = selected.indexOf(metric);

                    if (idx >= 0) {
                        console.log('Unselect metric:', metric);
                        selected.removeAt(idx);
                    }

                    set(this, 'selectedMetrics', selected);
                },

                selectAll: function() {
                    if (get(this, 'multiselect') === true) {
                        var metrics = get(this, 'metrics');

                        set(this, 'selectedMetrics', metrics);
                    }
                },

                unselectAll: function() {
                    set(this, 'selectedMetrics', []);
                },

                search: function(search) {
                    if(search) {
                        var mfilter = this.build_filter(search);
                        set(this, 'metricSearch', mfilter);
                    }
                    else {
                        set(this, 'metricSearch', null);
                    }
                }
            }
        });

        application.register('component:component-metricselector', component);
    }
});
