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

Ember.Application.initializer({
    name: 'UtilityClass',
    initialize: function(container, application) {

        var set = Ember.set,
            get = Ember.get;

        /**
         * Class for handling utility objects in canopsis
         * Allow hooks on canopsis utility features
         *
         * @class Utility
         * @memberOf canopsis.frontend.core
         * @extends Ember.Object
         */
        var Utility = Ember.Object.extend({

            init: function() {
                this._super();
                console.log('Registering utility object ' + get(this, 'name'));
            }

        });

        application.register('class:utility', Utility);
    }
});
