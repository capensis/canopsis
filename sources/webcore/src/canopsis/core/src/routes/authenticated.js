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
    name: 'AuthenticatedRoute',
    after: ['DataUtils', 'ActionsUtils'],
    initialize: function(container, application) {
        var dataUtils = container.lookupFactory('utility:data');
        var actionsUtils = container.lookupFactory('utility:actions');

        var get = Ember.get,
            set = Ember.set;

        /**
         * @class ApplicationRoute
         * @extends Ember.Route
         * @constructor
         */
        var route = Ember.Route.extend({
            /**
             * @method setupController
             * @param controller
             */
            setupController: function(controller) {
                this.controllerFor('application').onIndexRoute = true;
                actionsUtils.setDefaultTarget(controller);

                return this._super.apply(this, arguments);
            },

            actions: {
                editAndSaveModel: function(model, record_raw, callback) {
                    console.log("editAndSaveModel", record_raw);

                    model.setProperties(record_raw);

                    var promise = model.save();

                    if (callback !== undefined) {
                        promise.then(callback);
                    }
                },

                addRecord: function(crecord_type, raw_record, options) {
                    raw_record[crecord_type] = crecord_type;

                    if (options !== undefined && options.customFormAction !== undefined) {
                        options.customFormAction(crecord_type, raw_record);
                    } else {
                        var record = this.store.createRecord(crecord_type, raw_record);
                        var promise = record.save();
                        if (options !== undefined && options.callback !== undefined) {
                            promise.then(options.callback);
                        }
                    }
                    //send the new item via the API
                    //optional callback may ba called from crecord form
                },
                error: function(reason, transition) {
                    if (reason.status === 403 || reason.status === 401) {
                        this.loginRequired(transition);
                    }
                }
            }
        });
        application.register('route:authenticated', route);
    }
});
