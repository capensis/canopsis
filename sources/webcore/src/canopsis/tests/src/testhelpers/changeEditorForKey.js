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
/**
 * @method changeEditorForKey
 * @description call this method to use a specified editor for a given key name
 * @argument App not used
 * @argument key the key of the form element
 * @argument editorName the name of the editor to use
 * @return Ember.RSVP.Promise
 */
Ember.Test.registerAsyncHelper('changeEditorForKey', function(app, key, editorName) {
    Ember.Test.adapter.asyncStart();
    return Ember.Test.promise(function(resolve) {
        click('.left-side .frontend-config');
        waitForElement('input[name=title]').then(function(){
            click('.modal-body a[href=#editors]');
            fillIn('.modal-body .tab-pane.active input:first', key);
            fillIn('.modal-body .tab-pane.active input:last', editorName);
            click('.modal-body .btn-success');
            click('.modal-footer .btn-submit');

            waitForElementRemoval('.modal-backdrop').then(function() {
                Ember.Test.adapter.asyncEnd();
                Ember.run(null, resolve, true);
            });
        });
    });
});
