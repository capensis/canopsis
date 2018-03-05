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
    name: 'component-tabs',
    initialize: function(container, application) {
        /**
         * @component tabs
         * @description Base component of a tab component composition. To be used jointly with the following components : "tabscontentgroup", "tabsheadergroup", "tabcontent", "tabheader"
         * @example
         *   {{#component-tabs}}
         *       {{#component-tabsheadergroup}}
         *           {{component-tabheader ref="filter" label="Filter" active=true}}
         *           {{component-tabheader ref="output" label="Generated filter"}}
         *           {{component-tabheader ref="result" label="Result"}}
         *       {{/component-tabsheadergroup}}
         *       {{#component-tabscontentgroup}}
         *           {{#component-tabcontent ref="filter" active=true}}
         *               filter
         *           {{/component-tabcontent}}
         *           {{#component-tabcontent ref="output"}}
         *               output
         *           {{/component-tabcontent}}
         *           {{#component-tabcontent ref="result"}}
         *               result
         *           {{/component-tabcontent}}
         *       {{/component-tabscontentgroup}}
         *   {{/component-tabs}}
         */
        var component = Ember.Component.extend({});

        application.register('component:component-tabs', component);
    }
});
