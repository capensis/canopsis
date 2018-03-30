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
    name: 'component-searchbar',
    initialize: function(container, application) {
        var get = Ember.get;


        /**
         * @component searchbar
         * @description Search bar component
         *
         * Includes 3 tabs :
         *  - All : allow selection of every field
         *  - Indexed : only for indexed fields
         *  - filter : cfilter embedding
         *
         * This component is a WIP, it only supports basic search at the moment
         */
        var component = Ember.Component.extend({
            /**
             * @property showSearchOptions
             * @type boolean
             */
            showSearchOptions: false,

            /**
             * @property tagName
             * @type string
             */
            tagName: 'span',

            actions: {
                /**
                 * @method actions_searchInputAction
                 * @description Action to handle search input changes
                 */
                searchInputAction: function() {
                    var searchPhrase = get(this, 'value');
                    console.log('searchItems', this, this.controller, searchPhrase);

                    get(this, 'controller').target.set('searchCriterion', searchPhrase);
                },

                /**
                 * @method actions_clearSearch
                 * @description Action to handle search form clear.
                 */
                clearSearch: function () {
                    console.log('clear search field');
                    //clear text field
                    this.set('value', '');
                    //set search field
                    get(this, 'controller').target.set('searchFieldValue', '');
                    //trigger search
                    this.send('searchInputAction', '');

                }

            },

            /**
             * @property tabAllId
             * @type string
             * @description Computed property. generated id for the DOM element. It allows to handle children elements by their DOM ids without having id collision while several identical components are present on the DOM.
             */
            tabAllId: function() {
                console.log('tabAllId');

                return get(this, 'elementId') + 'TabAll';
            }.property('elementId'),

            /**
             * @property tabIndexedId
             * @type string
             * @description Computed property. generated id for the DOM element. It allows to handle children elements by their DOM ids without having id collision while several identical components are present on the DOM.
             */
            tabIndexedId: function() {
                return get(this, 'elementId') + 'TabIndexed';
            }.property('elementId'),

            /**
             * @property tabFilterId
             * @type string
             * @description Computed property. generated id for the DOM element. It allows to handle children elements by their DOM ids without having id collision while several identical components are present on the DOM.
             */
            tabFilterId: function() {
                return get(this, 'elementId') + 'TabFilter';
            }.property('elementId'),

            /**
             * @property tabAllHref
             * @type string
             * @description Computed property. generated href for the "all" link.
             */
            tabAllHref: function() {
                return '#' + get(this, 'elementId') + 'TabAll';
            }.property('elementId'),

            /**
             * @property tabIndexedHref
             * @type string
             * @description Computed property. generated href for the "indexed" link.
             */
            tabIndexedHref: function() {
                return '#' + get(this, 'elementId') + 'TabIndexed';
            }.property('elementId'),

            /**
             * @property tabFilterHref
             * @type string
             * @description Computed property. generated href for the "filter" link.
             */
            tabFilterHref: function() {
                return '#' + get(this, 'elementId') + 'TabFilter';
            }.property('elementId')
        });

        application.register('component:component-searchbar', component);
    }
});
