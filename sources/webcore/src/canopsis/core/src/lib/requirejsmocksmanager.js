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
    name: 'MocksRegistry',
    after: 'FormsUtils',
    initialize: function(container, application) {

        var formsUtils = container.lookupFactory('utility:forms');

        var set = Ember.set;

        /**
         * Singleton for handling requirejs module mocking
         *
         * @class RequirejsMockManager
         * @memberOf canopsis.frontend.core
         * @static
         */
        var MocksRegistryClass = Ember.Object.extend({
            init: function() {
                this._super();

                this.feedMocks();
            },

            /**
             * Refresh mocks list
             *
             * @method feedMocks
             */
            feedMocks : function() {
                var mocks = localStorage.getItem('files_mocks');
                var mocksArray = Ember.A();

                try {
                    mocks = JSON.parse(mocks);

                    if(!!mocks) {
                        for(var key in mocks) {
                            mocks[key].modulePath = key;
                            mocksArray.push(mocks[key]);
                        }
                    }
                } catch (e) {
                    console.error('Impossible to parse mocks JSON.');
                }

                set(this, 'mocks', mocksArray);
            },

            /**
             * Causes a mock to be edited, by showing an appropriate form, and handling the persistance of the mock
             *
             * @method editMock
             * @param mock {object} the mock to edit
             */
            editMock: function(mock) {
                var mocksManager = this;

                var recordWizard = formsUtils.showNew('modelform', mock, {
                    title: __('Edit mock') + ' ' + mock.modulePath,
                    inspectedItemType: 'requirejsmock'
                });

                recordWizard.submit.then(function(form) {
                    var mocks = localStorage.getItem('files_mocks');
                    try {
                        mocks = JSON.parse(mocks);

                        var mockModulePath = mock.modulePath;
                        delete mock.modulePath;
                        mock.module = mock.module;
                        mocks[mockModulePath] = mock;

                        localStorage.setItem('files_mocks', JSON.stringify(mocks));
                    } catch (e) {
                        console.error('Impossible to parse mocks JSON.');
                    }

                    mocksManager.feedMocks();
                });
            },

            /**
             * Causes a mock to be added, by showing an appropriate form, and handling the persistance of the mock
             *
             * @method addMock
             */
            addMock: function() {
                var mocksManager = this;

                var mock = {
                    modulePath: "",
                    module: "",
                    requirements: []
                };

                var recordWizard = formsUtils.showNew('modelform', mock, {
                    title: __('Add mock'),
                    inspectedItemType: 'requirejsmock'
                });

                recordWizard.submit.then(function(form) {
                    var mocks = localStorage.getItem('files_mocks');
                    try {
                        mocks = JSON.parse(mocks);

                        var mockModulePath = mock.modulePath;
                        delete mock.modulePath;
                        mock.module = mock.module;
                        mocks[mockModulePath] = mock;

                        localStorage.setItem('files_mocks', JSON.stringify(mocks));
                    } catch (e) {
                        console.error('Impossible to parse mocks JSON.');
                    }

                    mocksManager.feedMocks();
                });
            },

            /**
             * Causes a mock to be deleted and the persistance done
             *
             * @method deleteMock
             * @param {String} modulePath The full path of the JS module
             */
            deleteMock: function(modulePath) {
                var mocks = localStorage.getItem('files_mocks');

                try {
                    mocks = JSON.parse(mocks);
                    delete mocks[modulePath];
                    localStorage.setItem('files_mocks', JSON.stringify(mocks));
                } catch (e) {
                    console.error('Impossible to parse mocks JSON.');
                }

                this.feedMocks();
            }
        });

        var registry = MocksRegistryClass.create();
        application.register('registry:mocks', registry);
    }
});
