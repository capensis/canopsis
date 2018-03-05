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

    var get = Ember.get,
        isNone = Ember.isNone;

    /**
     * Helper to log something into the console.
     *
     * @author Gwenael Pluchon <info@gwenp.fr>
     */
    Ember.Handlebars.registerBoundHelper('clog', function(message, options) {
        console.tags.add("helper");

        var typeKey = get(this, 'model.type.typeKey');

        if (!isNone(typeKey)) {
            console.tags.add(typeKey);
        }

        var level = get(options, 'hash.level');

        if (isNone(level)) {
            level = 'log'
        }

        console[level](message);

        console.tags.remove("helper");

        if (!isNone(typeKey)) {
            console.tags.remove(typeKey);
        }


        return '';
    });

})();
