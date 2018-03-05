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
    name: 'LoginController',
    after: 'DataUtils',
    initialize: function(container, application) {

        var DataUtils = container.lookupFactory('utility:data');

        var set = Ember.set,
            get = Ember.get,
            isNone = Ember.isNone;

        /**
         * @class LoginController
         * @extends Ember.ObjectController
         * @constructor
         */

        var controller = Ember.ObjectController.extend({

            /**
             * @property content
             */
            content: {},

            needs: ['application'],

            /**
             * @method init
             */
            init: function() {
                this._super.apply(this, arguments);

                //TODO delete store in this#destroy
                var store = DS.Store.create({
                    container: get(this, "container")
                });

                set(this, 'store', store);
            },

            /**
             * @property authkey
             */
            authkey: function () {
                var authkey = localStorage.cps_authkey;
                if (authkey === 'undefined') {
                    authkey = undefined;
                }
                return authkey;
            }.property('authkey'),

            /**
             * @method authkeyChanged
             * @observes authkey
             */
            authkeyChanged: function() {
                localStorage.cps_authkey = get(this, 'authkey');
            }.observes('authkey'),

            /**
             * @method sessionStart
             * Tells the backend the user is logging in
             **/

            sessionStart: function () {
                loggedaccountAdapter = DataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:loggedaccount');
                window.moment.locale(get(this, 'record.ui_language'));
                username = get(this, 'record._id');
                loggedaccountAdapter.sessionStart(username);

                //Keep alive every 4"30'
                setInterval(function () {
                    loggedaccountAdapter.keepAlive(username);
                },1000 * 60 * 4.5);
            }

        });

        application.register('controller:login', controller);
    }
});
