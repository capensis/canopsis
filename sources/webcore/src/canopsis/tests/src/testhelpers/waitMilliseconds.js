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

/**
 * @class TestHelpers
 */
/** @method waitMilliseconds
 * @argument {integer} milliseconds Time to wait, in milliseconds
 * @description returns a promise that allows to wait during a specified time
 * @returns Ember.Test.promise
 * @example waitMilliseconds(500).then(function() {
    // code here will be executed after a 500ms pause
});
 */
Ember.Test.registerAsyncHelper('waitMilliseconds', function(app, milliseconds) {
    return Ember.Test.promise(function(resolve) {
        Ember.Test.adapter.asyncStart();
        var interval = setInterval(function(){
            clearInterval(interval);
            Ember.Test.adapter.asyncEnd();
            Ember.run(null, resolve, true);
        }, milliseconds);
    });
});
