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
    name:'UserprofilestatusmenuMixin',
    after: ['MixinFactory', 'FormsUtils', 'DataUtils', 'NotificationUtils'],
    initialize: function(container, application) {

        var Mixin = container.lookupFactory('factory:mixin');

        var formsUtils = container.lookupFactory('utility:forms');
        var dataUtils = container.lookupFactory('utility:data');
        var notificationUtils = container.lookupFactory('utility:notification');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

       /**
         * Mixin allowing to manage the current user profile, adding a button into the app status bar
         *
         * @class UserprofilestatusmenuMixin
         * @extensionfor ApplicationController
         * @static
         */
        var mixin = Mixin('Userprofilestatusmenu', {

            init: function() {
                if (!isNone(this.partials.statusbar)) {
                    this.partials.statusbar.pushObject('userstatusmenu');
                }
                this._super();
            },

            actions: {
                /**
                 * @event showUserProfile
                 * @descriptions Shows a Modelform with the user profile
                 */
                showUserProfile: function () {
                    var applicationController = this;

                    var ouser = get(this, 'controllers.login.record');
                    var recordWizard = formsUtils.showNew('modelform', ouser, {
                        title: get(ouser, '_id') + ' ' + __('profile'),
                        filterFieldByKey: {
                            firstname: {readOnly: true},
                            lastname: {readOnly: true},
                            authkey: {readOnly: true},
                            ui_language: true,
                            mail: true,
                            userpreferences: true,
                        }
                    });

                    recordWizard.submit.then(function(form) {
                        //ouilang = get(ouser, 'ui_language');  Interrestingly, that doesn't work
                        ouilang = ouser['_data']['ui_language'];

                        console.group('submit form:', form);

                        //generated data by user form fill
                        var editedUser = form.get('formContext');
                        console.log('save record:', editedUser);

                        set(editedUser, 'crecord_type', 'user');
                        editedUser.id= editedUser._id;

                        user = editedUser.serialize();
                        delete user['rights'];

                        var editedUserRecord = dataUtils.getStore().createRecord('user', user);
                        editedUserRecord.save();

                        notificationUtils.info(__('profile') + ' ' +__('updated'));

                        uilang = get(editedUser, 'ui_language');

                        console.log('Lang:', uilang, ouilang);
                        if (uilang !== ouilang) {
                            console.log('Language changed, will prompt for application reload');
                            applicationController.send(
                                'promptReloadApplication',
                                'Interface language has changed, reload canopsis to apply changes ?'
                            );
                        }

                        console.groupEnd();
                    });
                },
            }
        });

        application.register('mixin:userprofilestatusmenu', mixin);
    }
});
