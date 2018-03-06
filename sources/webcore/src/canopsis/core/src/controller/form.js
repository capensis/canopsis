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
    name: 'FormController',
    after: ['FormsUtils', 'DebugUtils'],
    initialize: function(container, application) {
        var formUtils = container.lookupFactory('utility:forms');
        var debugUtils = container.lookupFactory('utility:debug');

        var get = Ember.get,
            set = Ember.set;

        //TODO refactor this
        var eventedController = Ember.Controller.extend(Ember.Evented, {

            mergedProperties: ['partials'],

            _partials: {},

            refreshPartialsList: function() {
                console.log('refreshPartialsList', get(this, 'partials'));

                var partials = get(this, 'partials'),
                    mixins = get(this, 'content.mixins');

                set(this, '_partials', partials);

                if(Ember.isArray(mixins)) {
                    for (var i = 0, l = mixins.length; i < l; i++) {
                        partials = this.mergeMixinPartials(mixins[i], partials);
                    }
                }

                console.log('set partials for ', this, ' --> ', partials);
                set(this, '_partials', partials);
            },

            mergeMixinPartials: function(Mixin, partials) {
                var me = this;

                console.log("mergeMixinPartials mixin:", Mixin);
                if(mixinsRegistry.getByName(Mixin.decamelize())) {
                    var partialsToAdd = mixinsRegistry.getByName(Mixin.decamelize()).EmberClass.mixins[0].properties.partials;

                    for (var k in partialsToAdd) {
                        if (partialsToAdd.hasOwnProperty(k)) {
                            var partialsArray = partialsToAdd[k];

                            var partialKey = '_partials.' + k;
                            set(this, partialKey, union_arrays(get(this, partialKey), partialsArray));
                        }
                    }
                    return partials;
                }
            }
        });
        /**
         * @class FormController
         * @constructor
         * @description
         * Default is to display all fields of a given model if they are referenced into category list (in model)
         * options: is an object that can hold a set dictionnary of values to override
         *   - filters: is a list of keys to filter the fields that can be displayed
         *   - override_labels is an object that helps translate fields to display in form
         *   - callback, witch is called once form sent
         *   - plain ajax contains information that will be used insted of ember data mechanism
         */
        var controller = eventedController.extend({
            needs: ['application'],

            /**
             * @method init
             */
            init: function() {
                var formParent = get(this, 'formParent');
                set(this, 'previousForm', formParent);

                this._super.apply(this, arguments);
            },


            /**
             * @property confirmation
             * @type boolean
             */
            confirmation: false,

            /**
             * @property submit
             * @type $.Deferred
             * @static
             * @description
             *
             * Deferred to help manage form callbacks. You can implement :
             *  - done
             *  - always
             *  - fail
             *
             * with the form :
             * myForm.submit.done(function(){ [code here] })
             *
             * Caution: FormController#submit is NOT FormController#_actions#submit
             */
            submit: $.Deferred(),

            actions: {
                /**
                 * @event previousForm
                 * @description rollback to the previous form (if applicable)
                 */
                previousForm: function() {
                    var previousForm = get(this, 'previousForm');

                    console.log('previousForm', previousForm, this);
                    formUtils.showInstance(previousForm);
                },


                /**
                 * @event show
                 * @description
                 *
                 * Action triggered when the form is shown.
                 * By default it is used to reinitialize the "submit" deferred
                 */
                show: function() {
                    //reset submit defered
                    this.submit = $.Deferred();
                },

                /**
                 * @event submit
                 * @description
                 *
                 * Action triggered when the form is submit.
                 * If there is a parent form, its submit action is also called
                 */
                submit: function() {
                    console.log("onsubmit", this.formParent);

                    if (this.formParent !== undefined) {
                        this.formParent.send('submit', arguments);
                    }
                    else {
                        console.log("resolve modelform submit");
                        if ( this.confirmation ){
                            var record = this.formContext;
                            formUtils.showNew('confirmform', record , { title : " confirmation "  , newRecord : arguments[0]});
                        } else {
                            this.submit.resolve(this, arguments);
                            get(this, 'formwrapper').trigger("hide");
                        }
                    }
                },


                /**
                 * @event abort
                 * @description
                 *
                 * Action triggered when the form is aborted.
                 * If there is a parent form, its abort action is also called
                 */
                abort: function() {
                    if(this.formParent !== undefined) {
                        this.formParent.send('abort', arguments);
                    } else {
                        console.log('rejecting submit promise');
                        this.submit.reject();
                    }
                },

                /**
                 * @event inspectForm
                 * @description
                 *
                 * Inspect form in console and put it in the global $E variable
                 */
                inspectForm: function() {
                    console.group('inspectForm');
                    console.log('form:', this);

                    debugUtils.inspectObject(get(this, 'categorized_attributes'));

                    console.log('categorized_attributes available in $E');

                    console.groupEnd();
                }
            },

            partials: {
                buttons: ["formbutton-cancel"],
                debugButtons: ['formbutton-inspectform']
            },

            /**
             * @property title
             * @description
             *
             * The title of the form, usually displayed in the formwrapper header
             */
            title: function() {
                console.warn("Property \"title\" must be defined on the concrete class.");

                return "<Untitled form>";
            }.property()
        });

        application.register('controller:form', controller);
    }
});
