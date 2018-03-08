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
    name: 'Compare',
    initialize: function(container, application) {
        void(application);

        var Handlebars = window.Handlebars;

        var compare = function(lvalue, rvalue, options) {
            if (arguments.length < 3)
                throw new Error('Handlebars Helper \'compare\' needs 2 parameters');

            var operator = options.hash.operator || '==';

            var operators = {
                '==':       function(l,r) { return l == r; },
                '===':  function(l,r) { return l === r; },
                '!=':       function(l,r) { return l != r; },
                'lt':        function(l,r) { return l < r; },
                'gt':        function(l,r) { return l > r; },
                'lte':       function(l,r) { return l <= r; },
                'gte':       function(l,r) { return l >= r; },
                'typeof':   function(l,r) { return typeof l == r; }
            };

            if (!operators[operator])
                throw new Error('Handlebars Helper "compare" doesn\'t know the operator '+ operator);

            var result = operators[operator](lvalue,rvalue);

            if( result ) {
                return options.fn(this);
            } else {
                return options.inverse(this);
            }
        };

        Handlebars.registerHelper('compare', compare);
        Ember.Handlebars.helper('compare', compare);
    }
});



