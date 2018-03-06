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
    name: 'PaginationMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;


        /**
         * @mixin pagination
         *
         * @description
         * Implements pagination in ArrayControllers
         *
         * You should define on the ArrayController:
         *     - the `findOptions` property
         *     - the `findItems()` method
         *
         */
        var mixin = Mixin('pagination', {
            partials: {
                subHeader: ['itemsperpage'],
                actionToolbarButtons: ['pagination'],
                footer: ['pagination-infos'],
                rightFooter: ['pagination']
            },

            init: function() {
                this._super();
            },

            mixinsOptionsReady: function () {
                this._super();
                set(this, 'itemsPerPagePropositionSelected', get(this, 'itemsPerPage'));
            },

            isFirstPage: function () {
                return get(this, 'currentPage') === 1 && get(this, 'paginationFirstItemIndex') === 1;
            }.property('currentPage', 'paginationFirstItemIndex'),

            isLastPage: function () {
                return get(this, 'currentPage') === get(this, 'totalPages');
            }.property('currentPage', 'totalPages'),

            hasOnePage: function () {
                var onepage = get(this, 'totalPages') === 1;
                console.log('Is it a one page ?', onepage);
                return onepage;
            }.property('totalPages'),

            itemsPerPage: function() {
                var itemsPerPage = get(this, 'model.itemsPerPage') || get(this, 'mixinOptions.pagination.defaultItemsPerPage') || 5;

                return itemsPerPage;
            }.property('model.itemsPerPage', 'mixinOptions.pagination.defaultItemsPerPage'),

            paginationMixinContent: function() {
                console.warn('paginationMixinContent should be defined on the concrete class');
            },

            paginationMixinFindOptions: function() {
                console.warn('paginationMixinFindOptions should be defined on the concrete class');
            },

            actions: {
                prevPage: function() {
                    var currentPage = get(this, 'currentPage');

                    if(typeof currentPage === 'string') {
                        currentPage = parseInt(currentPage, 10);
                    }

                    if (currentPage > 1) {
                        set(this, 'currentPage', currentPage - 1);
                    }
                },
                nextPage: function() {
                    var currentPage = get(this, 'currentPage');

                    if(typeof currentPage === 'string') {
                        currentPage = parseInt(currentPage, 10);
                    }

                    if (currentPage < get(this, 'totalPages')) {
                        set(this, 'currentPage', currentPage + 1);
                    }
                },
                firstPage: function() {
                    this.set('currentPage', 1);
                },
                lastPage: function() {
                    if (get(this, 'currentPage') < get(this, 'totalPages')) {
                        set(this, 'currentPage', get(this, 'totalPages'));
                    }
                }
            },

            currentPage: 1,

            itemsDivided: function(){
                return get(this, 'itemsTotal') / get(this, 'itemsPerPage');
            }.property('itemsTotal', 'itemsPerPage'),

            itemsPerPagePropositions : function() {

                var choices = [5, 10, 20, 50];

                var customItemsPerPage = get(this, 'mixinOptions.pagination.customItemsPerPage');

                if (!isNone(customItemsPerPage) && Ember.isArray(customItemsPerPage)) {

                    var length = customItemsPerPage.length;

                    for (var i = 0; i < length; i++) {
                        if(!isNaN(customItemsPerPage[i])) {
                            choices.push(parseInt(customItemsPerPage[i]));
                        }

                    }
                }

                var itemsPerPage = get(this, 'itemsPerPage');
                if(!choices.contains(itemsPerPage)) {
                    choices.pushObject(itemsPerPage);
                }

                choices.sort(function (a, b) {return a - b;});

                return choices;
            }.property('itemsPerPagePropositionSelected'),

            onCurrentPageChanges: function() {
                this.refreshContent();
            }.observes('currentPage'),

            itemsPerPagePropositionSelectedChanged: function() {
                var userSelection = get(this, 'itemsPerPagePropositionSelected');

                if(get(this, 'loaded')) {
                    Ember.setProperties(this, {
                        'model.itemsPerPage': userSelection,
                        'currentPage': 1
                    });

                    this.saveUserConfiguration();
                }

                this.refreshContent();

            }.observes('itemsPerPagePropositionSelected'),

            refreshContent: function() {
                console.group('paginationMixin refreshContent', get(this, 'itemsPerPage'));

                if (get(this, 'paginationMixinFindOptions') === undefined) {
                    set(this, 'paginationMixinFindOptions', {});
                }

                var itemsPerPage = get(this, 'itemsPerPage');
                var start = itemsPerPage * (this.currentPage - 1);

                Ember.setProperties(this, {
                    'paginationMixinFindOptions.start': start,
                    'paginationFirstItemIndex': start + 1,
                    'paginationMixinFindOptions.limit': itemsPerPage
                });

                this._super.apply(this, arguments);

                console.groupEnd();
            },

            paginationLastItemIndex: function () {
                var itemsPerPage = get(this, 'itemsPerPage');

                var start = itemsPerPage * (this.currentPage - 1);

                return Math.min(start + itemsPerPage, get(this, 'itemsTotal'));
            }.property('widgetData', 'itemsTotal'),

            paginationFirstItemIndex: function () {
                var itemsPerPage = get(this, 'itemsPerPage');

                var start = itemsPerPage * (this.currentPage - 1);

                return start + 1;
            }.property('widgetData'),

            itemsTotal: function() {
                return get(this, 'widgetDataMetas').total;
            }.property('widgetDataMetas', 'widgetData'),

            totalPages: function() {
                if (get(this, 'itemsTotal') === 0) {
                    return 0;
                } else {
                    var itemsPerPage = get(this, 'itemsPerPage');
                    return Math.ceil(get(this, 'itemsTotal') / itemsPerPage);
                }
            }.property('itemsTotal')
        });

        application.register('mixin:pagination', mixin);
    }
});
