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
    name: 'component-querybuilder',
    initialize: function(container, application) {

        var i18n = container.lookupFactory('utility:i18n');

        var set = Ember.set,
            get = Ember.get,
            __ = Ember.String.loc;

        var build_filter = function(search) {
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
        };


        /**
         * @method extractKeysFromMongoFilter
         * @private
         * @param {object} filter the filter to extract keys from
         * @returns array of keys
         */
        var extractKeysFromMongoFilter = function(filter) {
            var results = [];
            var filterKeys = Ember.keys(filter);
            for (var i = 0; i < filterKeys.length; i++) {
                var currentKey = filterKeys[i];
                if(Ember.isArray(filter[currentKey])) {
                    for (var j = 0; j < filter[currentKey].length; j++) {
                        results = $.unique(results.concat(extractKeysFromMongoFilter(filter[currentKey][j])));
                    }
                } else if(filter[currentKey] === 'object') {
                    results.push(currentKey);
                    results = $.unique(results.concat(extractKeysFromMongoFilter(filter[currentKey])));
                } else {
                    results.push(currentKey);
                }
            }

            results = $(results).not(['$eq', '$regex', '$ne', '$gt', '$gte', '$lt', '$lte', '$in', '$ne', '$nin', '$exists', '$type', '$mod', '$where', '$text']).get();

            return results;
        };

        /**
         * @method cleanFilterEqOperators
         * @private
         * @description cleans every "$eq" operator from a mongo filter
         * @param {object} filter the filter to extract keys from
         * @returns the cleaned filter
         */
        var cleanFilterEqOperators = function(filter) {
			if (filter == null){
				return filter
			}
            var filterKeys = Ember.keys(filter);
            if(filter.$eq !== undefined) {
                return filter.$eq;
            }

            for (var i = 0; i < filterKeys.length; i++) {
                var currentKey = filterKeys[i];

                if(Ember.isArray(filter[currentKey])) {
                    var arrayClause = [];
                    for (var j = 0; j < filter[currentKey].length; j++) {
                        arrayClause.pushObject(cleanFilterEqOperators(filter[currentKey][j]));
                    }
                    filter[currentKey] = arrayClause;
                } else if(typeof filter[currentKey] === 'object') {
                    filter[currentKey] = cleanFilterEqOperators(filter[currentKey]);
                }
            }

            return filter;
        };

        /**
         * @component querybuilder
         * @description A query editor component to handle mongo filters edition, using jquery-querybuilder library.
         * This component also display preview of matching results in a dedicated tab
         */
        var component = Ember.Component.extend({
            /**
             * @property viewTabColumns
             * @type array
             * @description the columns to display in the results preview table
             */
            viewTabColumns: [{
                name:'connector',
                title:'connector'
            }, {
                name:'connector_name',
                title:'connector_name'
            }, {
                name:'component',
                title:'component'
            }, {
                name:'resource',
                title:'resource'
            }],

            /**
             * @property selectionTabColumns
             * @type object
             * @description Handpicked selection table columns
             */
            selectionTabColumns: [{
                name:'connector',
                title:'connector'
            }, {
                name:'connector_name',
                title:'connector_name'
            }, {
                name:'component',
                title:'component'
            }, {
                name:'resource',
                title:'resource'
            },{
                action: 'select',
                actionAll: (get(this, 'multiselect') === true ? 'selectAll' : undefined),
                title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-plus-sign"></span>'),
                style: 'text-align: center;'
            }],

            /**
             * @property helpModal
             * @type object
             * @description Handpicked selection table help modal info
             */
            helpModal: {
                title: __('Syntax'),
                content: ['<ul>',
                    '<li><code>co:regex</code> : ', __('look for a component'), '</li>',
                    '<li><code>re:regex</code> : ', __('look for a resource'), '</li>' ,
                    '<li><code>me:regex</code> : ' , __('look for a metric') , '(<code>me:</code>' , __(' isn\'t needed for this one') , ')</li>' ,
                    '<li>', __('combine all of them to improve your search'),' : <code>co:regex re:regex me:regex</code></li>' ,
                    '<li><code>co:</code>, <code>re:</code>, <code>me:</code> : ', __('look for non-existant field') , '</li>' ,
                    '</ul>'].join('')
            },

            /**
             * @property classNames
             * @type array
             */
            classNames: ['query-builder'],

            /**
             * @property selectionTabSearch
             * @type string
             * @description Handpicked selection table search criterion
             */
            selectionTabSearch: null,

            /**
             * @property filterValue
             * @type string
             * @description contains the computed filter value
             */
            filterValue: '',

            actions: {
                search: function(search) {
                    if(search) {
                        var mfilter = build_filter(search);
                        set(this, 'selectionTabSearch', JSON.stringify(mfilter));
                    }
                    else {
                        set(this, 'selectionTabSearch', null);
                    }
                },

                select: function(selection) {
                    var filterValue = get(this, 'filterValue') || '{}';

                    filterValue = JSON.parse(filterValue);

                    if(!filterValue) {
                        filterValue = {};
                    }

                    if(!filterValue['$or']) {
                        filterValue['$or'] = [];
                    }

                    if(filterValue['$and']) {
                        filterValue['$or'].pushObject({ '$and': filterValue['$and']});
                        filterValue = { '$or': filterValue['$or'] };
                    }

                    var selectedId = get(selection, 'id');
                    if(!filterValue['$or'].findBy('_id', selectedId)) {
                        filterValue['$or'].pushObject({
                            '_id' : get(selection, 'id')
                        });
                    }

                    this.$('.builder').queryBuilder('setRulesFromMongo', filterValue);
                    set(this, 'filterValue', JSON.stringify(filterValue, null, 2));
                }
            },

            /**
             * @method didInsertElement
             * @description contains jquery-querybuilder initialisation and data binding
             */
            didInsertElement: function() {
                var filters = [];
                var schema = window.schemasRegistry.getByName('event').schema;

                for (var i = 0; i < schema.categories.length; i++) {
                    var categoryName = schema.categories[i].title;
                    for (var j = 0; j < schema.categories[i].keys.length; j++) {
                        var key = schema.categories[i].keys[j];
                        var deepkeys = key.split('.'),
                            attribute = schema;

                        $.each(deepkeys, function(idx, key) {
                            if (attribute !== undefined && attribute.properties !== undefined) {
                                attribute = attribute.properties[key];
                            }
                            else {
                                attribute = undefined;
                            }
                        });

                        if (attribute === undefined) {
                            attribute = {'type': 'string'};
                        }

                        var name = key;

                        var filterElementDict = {
                            id: name,
                            label: name,
                            type: attribute.type,
                            optgroup: categoryName
                        };

                        if(attribute.type === 'boolean') {
                            filterElementDict.input = 'radio';
                            filterElementDict.values = {
                                true: 'True',
                                false: 'False'
                            };
                        }

                        if(attribute.type === 'number') {
                            filterElementDict.type = 'integer';
                        }

                        if(attribute.type !== 'object' && attribute.type !== 'array') {
                            filters.pushObject(filterElementDict);
                        }
                    }
                }

                filters.pushObject({
                    id: '_id',
                    label: '_id',
                    type: 'string',
                    optgroup: 'system'
                });

                var existingFiltersKeys = [];
                for (i = 0; i < filters.length; i++) {
                    existingFiltersKeys.push(filters[i].id);
                }

                var filterValue = this.get('filterValue') || undefined;

                if(filterValue && filterValue !== '{}') {
                    filterValue = JSON.parse(filterValue);
                    filterValue = cleanFilterEqOperators(filterValue);
                    var additionnalKeys = extractKeysFromMongoFilter(filterValue);
                    additionnalKeys = $(additionnalKeys).not(existingFiltersKeys).get();

                    for (i = 0; i < additionnalKeys.length; i++) {
						var found = false
						var it = 0
						while (!found && it < filters.length){
							if(filters[it]["id"] == additionnalKeys[i]){
								found = true
							}
							it++
						}
						if (!found){
							filters.pushObject({
								id: additionnalKeys[i],
								label: additionnalKeys[i],
								type: 'string',
								optgroup: 'custom'
							})
						}
                    }
                }

                this.$().find('.builder').queryBuilder({
                    filters: filters,
                    lang_code: i18n.lang,
                    plugins: [
                        // 'sortable', //TODO activate this when events will be triggered at rule drop, probably on version 2.3.1 of querybuilder
                        'key-editable-select'
                    ]
                });

                if(filterValue && filterValue !== '{}') {
                    this.$().find('.builder').queryBuilder('setRulesFromMongo', filterValue);
                }

                var component = this;
                this.$().find('.builder').on('afterUpdateRuleFilter.queryBuilder afterDeleteGroup.queryBuilder afterAddRule.queryBuilder afterDeleteRule.queryBuilder afterUpdateRuleValue.queryBuilder afterUpdateRuleOperator.queryBuilder afterUpdateGroupCondition.queryBuilder', function () {
                    var result = component.$().find('.builder').queryBuilder('getMongo');

                    if (!$.isEmptyObject(result)) {
                        component.set('filterValue', JSON.stringify(result, null, 2));
                    }
                });

                this.$('.btn-reset').on('click', function() {
                    component.$().find('.builder').queryBuilder('reset');
                    component.set('filterValue', '{}');
                });
            }
        });

        application.register('component:component-querybuilder', component);
    }
});
