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
    name: 'component-table',
    after: 'PaginationMixin',
    initialize: function(container, application) {
        var PaginationMixin = container.lookupFactory('mixin:pagination');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component table
         */
        var component = Ember.Component.extend(PaginationMixin, {
            model: undefined,
            modelfilter: undefined,
            data: undefined,

            columns: [],
            items: [],

            mixinOptions: {
                pagination: {
                    itemsPerPage: 5
                }
            },

            onDataChange: function() {
                this.refreshContent();
            }.observes('data.@each'),

            onModelFilterChange: function() {
                set(this, 'currentPage', 1);
                this.refreshContent();
            }.observes('modelfilter'),

            init: function() {
                this._super(arguments);

                if (!isNone(get(this, 'model'))) {
                    set(this, 'store', DS.Store.create({
                        container: get(this, 'container')
                    }));
                }

                set(this, 'widgetDataMetas', {total: 0});

                if (isNone(get(this, 'items'))) {
                    set(this, 'items', []);
                }
            },

            didInsertElement: function() {
                this.refreshContent();
            },

            refreshContent: function() {
                this._super(arguments);

                this.findItems();

                console.log(get(this, 'widgetDataMetas'));
            },

            findItems: function() {
                //TODO: clean this try/catch
                try {
                    var me = this;

                    var store = get(this, 'store'),
                        model = get(this, 'model'),
                        modelfilter = get(this, 'modelfilter');

                    var query = {
                        limit: get(this, 'paginationMixinFindOptions.limit')
                    };

                    var queryStartOffsetKeyword = get(this, 'queryStartOffsetKeyword') || 'skip';
                    query[queryStartOffsetKeyword] = get(this, 'paginationMixinFindOptions.start');

                    if (model !== undefined) {
                        if(modelfilter !== null) {
                            query.filter = modelfilter;
                        }

                        store.findQuery(model, query).then(function(result) {
                            console.log('Received data for table:', result);

                            set(me, 'widgetDataMetas', get(result, 'meta'));
                            set(me, 'items', get(result, 'content'));

                        });
                    }
                    else {
                        var items = get(this, 'data').slice(
                            query.skip,
                            query.skip + query.limit
                        );

                        set(this, 'widgetDataMetas', {
                            total: get(this, 'data.length')
                        });

                        set(this, 'items', items);

                    }
                } catch(err) {
                    console.warn('extractItems not updated:', err);
                }
            },

            actions: {
                do: function(action, item) {
                    this.targetObject.send(action, item);
                }
            }
        });

        application.register('component:component-table', component);
    }
});

