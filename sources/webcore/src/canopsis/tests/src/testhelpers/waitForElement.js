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
/** @method waitForElement
 * @argument {jQuerySelector} element A jquery selector that match an element not present in the DOM
 * @description returns a promise that allows to wait until a specified element is present on the DOM
 * @returns Ember.Test.promise
 * @example waitForElement('.modal-dialog').then(function() {
    // code here will be executed when the element with the class "modal-dialog" is on the DOM
});
 */
Ember.Test.registerAsyncHelper('waitForElement', function(app, element) {
    return Ember.Test.promise(function(resolve) {
        var count = 0;
        Ember.Test.adapter.asyncStart();
        var interval = setInterval(function(){
            if($(element).length>0){
                clearInterval(interval);
                Ember.Test.adapter.asyncEnd();
                Ember.run(null, resolve, true);
            } else {
                count++;
                if(count >= 1000) {
                    clearInterval(interval);
                    Ember.Test.adapter.asyncEnd();
                    console.error('element', element, 'not found');
                }
            }
        }, 100);
    });
});
