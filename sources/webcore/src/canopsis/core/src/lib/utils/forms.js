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
    name: 'FormsUtils',
    after: ['UtilityClass', 'FormsRegistry', 'RoutesUtils', 'DataUtils'],
    initialize: function(container, application) {
        var Utility = container.lookupFactory('class:utility');

        var routesUtils = container.lookupFactory('utility:routes');
        var dataUtils = container.lookupFactory('utility:data');

        var formsregistry = container.lookupFactory('registry:forms');

        var get = Ember.get,
            set = Ember.set,
            __ = Ember.String.loc;

        var formUtils = Utility.create({
            name: 'forms',

            instantiateForm: function(formName, formContext, options) {
                void (formContext);

                if(this.formController) {
                    this.formController.destroy();
                }

                console.log('try to instantiate form', formName, options.formParent);
                var classDict = options;

                options.formName = formName;
                classDict.target = routesUtils.getCurrentRouteController();
                classDict.container = dataUtils.getEmberApplicationSingleton().__container__;

                if(container.lookupFactory('form:' + formName) === undefined) {
                    console.error('the form', formName, 'was not found');
                }

                var formController = container.lookupFactory('form:' + formName).create(classDict);
                formController.refreshPartialsList();

                this.formController = formController;

                return formController;
            },

            showInstance: function(formInstance) {
                formInstance.empty_validationFields();

                set(formsregistry.formwrapper, 'form.validateOnInsert', false);
                set(formsregistry.formwrapper, 'form', formInstance);
                set(formsregistry.formwrapper, 'formName', formInstance.formName);
            },

            showNew: function(formName, formContext, options) {
                if (options === undefined) {
                    options = {};
                }

                if (Ember.isNone(formContext) || Ember.isNone(get(formContext, 'crecord_type'))) {
                    console.warn('There is no crecord_type in the given record. Form may not display properly.');
                }

                console.log("Form generation", formName);

                var formController = this.instantiateForm(formName, formContext, options);
                console.log("formController", formController);

                routesUtils.getCurrentRouteController().send('showEditFormWithController', formController, formContext, options);

                return formController;
            },

            editRecord: function(record, callback) {
                var widgetWizard = this.showNew('modelform', record);
                console.log('widgetWizard', widgetWizard);

                widgetWizard.submit.then(function() {
                    console.log('record saved');
                    var save = record.save();
                    if (!isNone(callback)) {
                        save.then(callback);
                    }

                    widgetWizard.trigger('hidePopup');
                    widgetWizard.destroy();
                });

                return widgetWizard;
            },

            editSchemaRecord: function (schemaName, container, callback) {

                var forms = this;

                var dataStore = DS.Store.create({
                    container: container//get(this, "container")
                });

                dataStore.findQuery(schemaName, {}).then(function(queryResults) {

                    console.log('queryResults', queryResults);

                    var record = get(queryResults, 'content')[0];

                    //it is always translated this way
                    var errorMessage = [
                        schemaName,
                        ' ',
                        __('information not registered in database.'),
                        ' ',
                        __('Please contact your administrator.'),
                    ].join('');

                    if (record) {
                        forms.editRecord(record, callback);
                     } else {
                        forms.error(errorMessage);
                    }


                });
            },

            addRecord: function(record_type) {
                routesUtils.getCurrentRouteController().send('show_add_crecord_form', record_type);
            }
        });

        application.register('utility:forms', formUtils);
    }
});
