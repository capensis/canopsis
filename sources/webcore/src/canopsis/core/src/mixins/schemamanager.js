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
    name:'SchemamanagerMixin',
    after: ['MixinFactory', 'FormsUtils', 'DataUtils', 'SearchmethodsRegistry', 'SchemasLoader'],
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');
        var formsUtils = container.lookupFactory('utility:forms');
        var dataUtils = container.lookupFactory('utility:data');
        var schemasLoader = container.lookupFactory('deprecated:schemasLoader');

        var canopsisConfiguration = window.canopsisConfiguration;

        var searchmethodsRegistry = container.lookupFactory('registry:searchmethods');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;


        /**
         * Mixin allowing to create schemas from the status bar in debug mode
         *
         * @class SchemamanagerMixin
         * @extensionfor ApplicationController
         * @static
         */
        var mixin = Mixin('schemamanager', {
            configuration: canopsisConfiguration,
            available_types : schemasLoader.generatedModels,

            init: function() {
                this.partials.statusbar.pushObject('schemamanagerstatusmenu');
                this._super();
            },

            displayedTypes: function() {
                var searchCriterion = get(this, 'schemaManagerSearchCriterion');
                console.log('recompute displayedTypes');
                if(isNone(searchCriterion)) {
                    return get(this, 'available_types');
                } else {
                    console.log('searchCriterion', searchCriterion);
                    var res = searchmethodsRegistry.getByName('simple').filter(get(this, 'available_types'), {
                        propertyToCheck: 'name',
                        searchCriterion: get(this, 'schemaManagerSearchCriterion')
                    });

                    return res;
                }
            }.property('available_types', 'schemaManagerSearchCriterion'),

            actions: {
                /**
                 * @event addModelInstance
                 * @descriptions Shows a form to a specified model
                 * @param {String} type the model type, lower-cased
                 */
                addModelInstance: function(type) {
                    console.log("add", type);

                    var record = dataUtils.getStore().createRecord(type, {
                        crecord_type: type.underscore()
                    });

                    console.log('temp record', record);

                    var recordWizard = formsUtils.showNew('modelform', record, { title: __("Add ") + type });

                    recordWizard.submit.done(function() {
                        record.save();
                    });
                },
            }
        });

        application.register('mixin:schemamanager', mixin);
    }
});
