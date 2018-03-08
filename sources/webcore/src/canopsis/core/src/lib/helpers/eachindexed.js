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

(function() {

    Ember.Handlebars.registerHelper('eachIndexed', function eachHelper(path, options) {
        var keywordName = 'item',
            fn;
        // Process arguments (either #earchIndexed bar, or #earchIndexed foo in bar)
        if (arguments.length === 4) {
            Ember.assert('If you pass more than one argument to the eachIndexed helper, it must be in the form #eachIndexed foo in bar', arguments[1] === 'in');
            Ember.assert(arguments[0] +' is a reserved word in #eachIndexed', $.inArray(arguments[0], ['index', 'index+1', 'even', 'odd']));
            keywordName = arguments[0];

            options = arguments[3];
            path = arguments[2];
            options.hash.keyword = keywordName;
            if (path === '') { path = 'this'; }
        }

        if (arguments.length === 1) {
            options = path;
            path = 'this';
        }

        // Wrap the callback function in our own that sets the index value
        fn = options.fn;
        function eachFn(){
        var keywords = arguments[1].data.keywords,
            view = arguments[1].data.view,
            index = view.contentIndex,
            list = view._parentView.get('content') || [],
            len = list.length;

            // Set indexes
            keywords.index = index;
            keywords.index_1 = index + 1;
            keywords.first = (index === 0);
            keywords.last = (index + 1 === len);
            keywords.even = (index % 2 === 0);
            keywords.odd = !keywords.even;

            arguments[1].data.keywords = keywords;

            return fn.apply(this, arguments);
        }

        options.fn = eachFn;

        // Render
        options.hash.dataSourceBinding = path;
        if (options.data.insideGroup && !options.hash.groupedRows && !options.hash.itemViewClass) {
            new Ember.Handlebars.GroupedEach(this, path, options).render();
        } else {
            return Ember.Handlebars.helpers.collection.call(this, 'Ember.Handlebars.EachView', options);
        }
    });
})();
