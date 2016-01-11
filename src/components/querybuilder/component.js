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
            viewTabColumns: [{ name:'connector', title:'connector' }, { name:'connector_name', title:'connector_name' }, { name:'component', title:'component' }, { name:'resource', title:'resource' }],

            /**
             * @property filterValue
             * @type string
             * @description contains the computed filter value
             */
            filterValue: '',
            /**
             * @method didInsertElement
             * @description contains jquery-querybuilder initialisation and data binding
             */
            didInsertElement: function() {
                var filters = [];
                Ember.get(container.lookupFactory('model:event'), 'attributes').forEach(function(attribute, name) {
                    console.error(attribute);
                    var filterElementDict = {
                        id: name,
                        label: name,
                        type: attribute.type
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
                });
                this.$().find('.builder').queryBuilder({
                    filters: filters,
                    lang_code: i18n.lang,
                    plugins: [
                        'sortable',
                        'key-editable-select'
                    ]
                });

                if(this.get('filterValue')) {
                    this.$().find('.builder').queryBuilder('setRulesFromMongo', JSON.parse(this.get('filterValue')));
                }

                var component = this;
                this.$().find('.builder').on('afterUpdateRuleFilter.queryBuilder afterDeleteGroup.queryBuilder afterAddRule.queryBuilder afterDeleteRule.queryBuilder afterUpdateRuleValue.queryBuilder afterUpdateRuleOperator.queryBuilder afterUpdateGroupCondition.queryBuilder', function () {
                    var result = component.$().find('.builder').queryBuilder('getMongo');

                    if (!$.isEmptyObject(result)) {
                        component.set('filterValue', JSON.stringify(result, null, 2));
                    }
                });
            }
        });

        application.register('component:component-querybuilder', component);
    }
});
