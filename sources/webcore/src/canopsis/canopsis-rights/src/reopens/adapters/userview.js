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
    name: 'CanopsisRightsUserviewAdapterReopen',
    after: ['DataUtils', 'UserviewAdapter'],
    initialize: function(container, application) {
        var dataUtils = container.lookupFactory('utility:data');
        var UserviewAdapter = container.lookupFactory('adapter:userview');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @class UserviewAdapter
         * @extends ApplicationAdapter
         * @constructor
         * @description UserviewAdapter reopen
         */
        UserviewAdapter.reopen({
            /**
             * @method updateRecord
             * @param {DS.Store} store
             * @param {DS.Model} type
             * @param {DS.Model} userview
             * @return {Promise} promise
             *
             * Manage right creation and modification on view creation or update.
             * Note that the createRecord method is never used with the userview adapter.
             */
            updateRecord: function(store, type, userview) {
                rightsRegistry = dataUtils.getEmberApplicationSingleton().__container__.lookupFactory('registry:rights');

                formattedViewId = get(userview, 'id').replace('.', '_');

                if(isNone(rightsRegistry.getByName(formattedViewId))) {
                    //The right does not exists, assume that the view is brand new

                    var right = dataUtils.getStore().createRecord('action', {
                          enable: true,
                          crecord_type: 'action',
                          type: 'RW',
                          _id: formattedViewId,
                          id: formattedViewId,
                          crecord_name: formattedViewId,
                          desc: 'Rights on view : ' + get(userview, 'crecord_name')
                    });
                    right.save();

                    rightsRegistry.add(right, get(right, 'crecord_name'));

                    //TODO Add the correct right to the current user, to allow him to display the view
                    var loginController = dataUtils.getLoggedUserController();

                    var rights = get(loginController, 'record.rights');

                    set(rights, formattedViewId, { checksum : 7 });
                    var record = get(loginController, 'record');
                    record.save();
                } else {
                    //TODO the right already exists, it's an update
                    //TODO replace the userview name if it has changed
                }

                return this._super.apply(this, arguments);
            }
        });
    }
});
