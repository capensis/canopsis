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
    name: 'TaskForm',
    after: ['FormFactory', 'ModelForm', 'FormsUtils'],
    initialize: function(container, application) {
        var FormFactory = container.lookupFactory('factory:form');
        var ModelFormController = container.lookupFactory('form:modelform');
        var formsUtils = container.lookupFactory('utility:forms');

        var get = Ember.get,
            set = Ember.set;

        var formOptions = {
            subclass: ModelFormController
        };

        var form = FormFactory('taskform', {
            title: 'Configure Job settings',
            scheduled: true,

            parentContext: function() {
                return get(this, 'formParent.formContext');
            }.property('formParent'),

            init: function() {
                this._super();

                set(this, 'store', DS.Store.create({
                    container: get(this, "container")
                }));

                if(get(this, 'scheduled') === true) {
                    set(this, 'partials.buttons', ["formbutton-next", "formbutton-cancel"]);
                } else {
                    set(this, 'partials.buttons', ["formbutton-cancel", "formbutton-submit"]);
                }

                var wizard = formsUtils.showNew('scheduleform', get(this, 'parentContext'), {
                    formParent: this,
                    title: 'Configure Schedule'
                });

                set(this, 'nextForm', wizard);
                this.refreshPartialsList();
            },

            actions: {
                next: function() {
                    console.group('configureTask');

                    console.log('parent:', get(this, 'parentContext'));
                    console.log('ctx:', get(this, 'formContext'));
                    console.log('form:', get(this, 'nextForm'));

                    formsUtils.showInstance(get(this, 'nextForm'));

                    console.groupEnd();
                },

                submit: function() {
                    console.group('submitTask');

                    var parentForm = get(this, 'formParent'),
                        ctx = get(this, 'formContext');
                    var job = get(parentForm, 'formContext');

                    set(job, 'params', ctx);

                    console.groupEnd();

                    if(get(this, 'scheduled') === true) {
                        this._super(arguments);
                    }
                    else {
                        this._super([job]);
                    }
                }
            },

            partials: {
                buttons: ["formbutton-cancel", "formbutton-next"]
            }
        }, formOptions);

        application.register('form:taskform', form);

        Ember.TEMPLATES['taskform'] = Ember.TEMPLATES['modelform'];
    }
});
