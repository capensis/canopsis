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
    name: 'SortablearrayMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @mixin sortablearray
         *
         * @description
         * Implements sorting in arraycontrollers
         *
         * You should define on the ArrayController:
         *     - the `findOptions` property
         *     - the `refreshContent()` method
         *
         */
        var mixin = Mixin('sortablearray', {


            init: function () {
                this._super();

                set(this, 'sort_direction', 'ASC');

                var sorted = get(this, 'model.user_sortedAttribute');
                if(sorted) {
                    set(this, 'sortedAttribute', sorted);
                    var direction = get(this, 'sortedAttribute.headerClassName') === 'sorting_asc' ? 'ASC': 'DESC';
                    set(this, 'sort_direction', direction);
                }
            },

            actions: {
                sort: function(attribute) {
                    var old_direction = get(this, 'sort_direction');
                    var direction = old_direction === 'ASC' ? 'DESC': 'ASC';
                    set(this, 'sort_direction', direction);

                    if (get(this, 'model.user_sortedAttribute') !== undefined) {
                        set(this, 'model.user_sortedAttribute.headerClassName', 'sorting');
                    }

                    console.log('attribute', attribute);
                    set(attribute, 'headerClassName', 'sorting_' + direction.toLowerCase());

                    set(this, 'sortedAttribute', attribute);
                    set(this, 'model.user_sortedAttribute', attribute);

                    console.log('sortBy', arguments);
                    if (this.findOptions === undefined) {
                        this.findOptions = {};
                    }

                    this.findOptions.sort = JSON.stringify([{
                        property: attribute.field,
                        direction: direction
                    }]);

                    // set(this, 'model.user_sortedAttribute', attribute.field);
                    this.saveUserConfiguration();

                    this.refreshContent();
                }
            },

            computeShownColumns: function(columns) {
                var shown_columns = this._super(columns);

                var sortedAttribute = get(this, 'sortedAttribute'),
                    columnSort = get(this, 'default_column_sort');

                for(var column = 0, l = shown_columns.length; column < l; column++) {
                    //Manage sort icon from default sort
                    if (!isNone(columnSort) &&
                        columnSort.property === shown_columns[column].field &&
                        !isNone(columnSort.direction) &&
                        (isNone(sortedAttribute) || sortedAttribute === {})) {

                        var headerClass = columnSort.direction === 'ASC' ? 'sorting_asc' : 'sorting_desc';
                        shown_columns[column].headerClassName = headerClass;
                    }
                }

                return shown_columns;
            },

            attributesKeys: function() {
                console.log('attributesKeys from sortableArray');
                var keys = this._super.apply(this);
                var sortedAttribute = get(this, 'sortedAttribute');

                console.log('sortedAttribute', sortedAttribute);
                if(sortedAttribute !== undefined)
                {
                    for (var i = 0, li = keys.length; i < li; i++) {
                        var currentKey = keys[i];
                        var sortedAttributeField = get(sortedAttribute, 'field');
                        var sortedAttributeHeaderClassName = get(sortedAttribute, 'headerClassName');

                        if(get(currentKey, 'field') === sortedAttributeField) {
                            set(currentKey, 'headerClassName', sortedAttributeHeaderClassName);
                        } else {
                            set(currentKey, 'headerClassName', 'sorting');
                        }
                    }
                } else {
                    if(!keys) {
                        return [];
                    } else {
                        for (var j = 0, lj = keys.length; j < lj; j++) {
                            set(keys[j], 'headerClassName', 'sorting');
                        }
                    }
                }
                return keys;
            }.property('inspectedProperty', 'inspectedDataArray')
        });

        application.register('mixin:sortablearray', mixin);
    }
});
