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
    name:'component-password',
    after: 'HashUtils',
    initialize: function(container, application) {
        var hash = container.lookupFactory('utility:hash');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component password
         * @description A simple password input component, that automatically converts the typed password to a hash, using "sha1" or "md5"
         */
        var component = Ember.Component.extend({
            /**
             * @property method
             * @description the hash method to be used to convert the password
             */
            method: undefined,

            /**
             * @property content
             * @description the hash of the password
             * @type string
             */
            content: undefined,

            /**
             * @method init
             */
            init: function () {
                this._super.apply(this, arguments);

                var allowed_methods = ['sha1', 'md5'];
                var method_name = get(this, 'method');

                if (!isNone(method_name) && allowed_methods.indexOf(method_name) === -1) {
                    console.warning('Invalid method, using sha1:', method_name);
                    set(this, 'method', 'sha1');
                }
            },

            /**
             * @method onUpdate
             * @description observes the password property
             */
            onUpdate: function () {
                var pass = get(this, 'password');
                var method_name = get(this, 'method');

                if (!isNone(method_name)) {
                    var method = get(hash, method_name);

                    pass = method(pass);
                }

                set(this, 'content', pass);

            }.observes('password')

        });

        application.register('component:component-password', component);
    }
});
