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
    name:'RequirejsmocksmanagerMixin',
    after: ['MocksRegistry', 'MixinFactory'],
    initialize: function(container, application) {

        var Mixin = container.lookupFactory('factory:mixin');
        var requirejsmocksmanager = container.lookupFactory('registry:mocks');

        /**
         * Mixin allowing to mock js code from inside the UI, adding a dedicated statusbar button into the app statusbar only in debug mode
         *
         * @class RequirejsmocksmanagerMixin
         * @extensionfor ApplicationController
         * @static
         */
        var mixin = Mixin('Requirejsmocksmanager', {

            requirejsmocksmanager: requirejsmocksmanager,

            init: function() {
                this.partials.statusbar.pushObject('requirejsmockingstatusmenu');
                this._super();
            },

            actions: {
                /**
                 * Causes a mock to be deleted
                 *
                 * @event deleteMock
                 * @param {String} modulePath The full path of the JS module
                 */
                deleteMock: function(modulePath) {
                    requirejsmocksmanager.deleteMock(modulePath);
                },

                /**
                 * Causes a mock to be edited, by showing an appropriate form, and handling the persistance of the mock
                 *
                 * @event editMock
                 * @param {Object} mock The mock to edit
                 */
                editMock: function(mock) {
                    requirejsmocksmanager.editMock(mock);
                },

                /**
                 * Causes a mock to be added, by showing an appropriate form, and handling the persistance of the mock
                 *
                 * @event addMock
                 */
                addMock: function() {
                    requirejsmocksmanager.addMock();
                }
            }
        });

        application.register('mixin:requirejsmocksmanager', mixin);
    }
});
