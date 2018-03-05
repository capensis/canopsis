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
 *
 * @module canopsis-frontend-core
 */

Ember.Application.initializer({
    name: 'ModelForm',
    after: ['FormFactory', 'InspectableItemMixin', 'ValidationMixin', 'SlugUtils'],
    initialize: function(container, application) {
        var FormFactory = container.lookupFactory('factory:form');
        var schemasRegistry = window.schemasRegistry;
        var InspectableitemMixin = container.lookupFactory('mixin:inspectable-item');
        var ValidationMixin = container.lookupFactory('mixin:validation');
        var slugUtils = container.lookupFactory('utility:slug');


        var set = Ember.set,
            get = Ember.get,
            isNone = Ember.isNone;

        var formOptions = {
            mixins: [
                InspectableitemMixin,
                ValidationMixin
            ]
        };

        /**
         * @class Modelform
         * @description Generic form which dynamically generates its content by reading a model's schema
         */
        var form = FormFactory('modelform', {
            needs: ['application'],

            partials: {
                debugButtons: ['formbutton-inspectform']
            },

            validationFields: Ember.computed(function() {return Ember.A();}),
            ArrayFields: Ember.A(),

            init: function() {
                this._super();
                this.set('partials.buttons', ['formbutton-submit']);
            },

            filterUserPreferenceCategory: function (category, keyFilters) {
                var keys = get(category, 'keys');
                set(category, 'keys', []);

                for (var i = 0, l = keys.length; i < l; i++) {
                    console.log('key', keys[i]);

                    if (this.get('userPreferencesOnly')) {
                        //isUserPreference is set to true in the key schema field.
                        if (keys[i].model && keys[i].model.options && keys[i].model.options.isUserPreference) {
                            get(category, 'keys').push(keys[i]);
                        }
                    } else {
                        //Filter from form parameter
                        if (keyFilters[keys[i].field]) {
                            console.log('magic keys', keys[i]);
                            if (keyFilters[keys[i].field].readOnly) {
                                keys[i].model.options.readOnly = true;
                            }
                            get(category, 'keys').push(keys[i]);
                        }
                    }
                }
                return category;
            },


            /**
             * @property categories
             * @type {Array}
             */
            categories: function(){
                var res = get(this, 'categorized_attributes');
                var category_selection = [];
                if(res instanceof Array) {
                    for(var i = 0, l = res.length; i < l; i++) {
                        var category = res[i];

                        category.slug = slugUtils.slug(category.title);
                        console.log('current category', category);

                        if (get(this, 'filterFieldByKey') || get(this, 'userPreferencesOnly')) {
                            //filter on user preferences fields only
                            //if (category)
                            this.filterUserPreferenceCategory(category, get(this, 'filterFieldByKey'));
                            if (category.keys.length) {
                                category_selection.push(res[i]);
                            }

                            console.log('category');
                            console.log(category);
                        } else {
                            //select all
                            category_selection.push(res[i]);
                        }
                    }
                    if (category_selection.length) {
                        set(category_selection[0], 'isDefault', true);
                    }
                    return category_selection;
                }
                else {
                    return [];
                }
            }.property('categorized_attributes'),

            onePageDisplay: function () {
                //TODO search this value into schema
                return false;
            }.property(),

            inspectedDataItem: function() {
                return get(this, 'formContext');
            }.property('formContext'),

            /**
             * @property inspectedItemType
             * @type {string} lowercased model name
             * @description
             *
             * Used to dynamically create form editors and assign values to the edited model.
             * To force editing as a specific model type, override this property.
             */
            inspectedItemType: function() {
                console.log('recompute inspectedItemType', get(this, 'formContext'));

                if (get(this, 'formContext.xtype')) {
                    return get(this, 'formContext.xtype');
                } else {
                    if(get(this, 'formContext.crecord_type') === "user") {
                        return "account";
                    }
                    return get(this, 'formContext.crecord_type') || get(this, 'formContext.connector_type')  ;
                }

            }.property('formContext'),

            actions: {
                submit: function() {
                    if (this.validation !== undefined && !this.validation()) {
                        return;
                    }
                    console.log('submit action', arguments);

                    var override_inverse = {};

                    if(this.isOnCreate && this.modelname){

                        var stringtype = this.modelname.charAt(0).toUpperCase() + this.modelname.slice(1);

                        //TODO use the real schema, not the dict used to create it
                        //retreive the corresponding schema dict
                        var model = schemasRegistry.getByName(stringtype);

                        if(model) {
                            for(var fieldName in model){
                                if(model.hasOwnProperty(fieldName)) {
                                    var field = model[fieldName];
                                    if(field && field._meta &&  field._meta.options){
                                        var metaoptions = field._meta.options;
                                        if( 'setOnCreate' in metaoptions){
                                            var value = metaoptions.setOnCreate;
                                            set(this, 'formContext.' + fieldName, value);
                                        }
                                    }
                                }
                            }
                        }
                    }
                    //will execute callback from options if any given
                    var options = get(this, 'options');

                    if(options && options.override_labels) {
                        for(var key in options.override_labels) {
                            if(options.override_labels.hasOwnProperty(key)) {
                                override_inverse[options.override_labels[key]] = key;
                            }
                        }
                    }

                    var categories = get(this, 'categorized_attributes');

                    console.log('setting fields');
                    for (var i = 0, li = categories.length; i < li; i++) {
                        var category = categories[i];
                        for (var j = 0, lj = category.keys.length; j < lj; j++) {
                            var attr = category.keys[j];
                            var categoryKeyField = attr.field;
                            //set back overried value to original field
                            if (override_inverse[attr.field]) {
                                categoryKeyField = override_inverse[attr.field];
                            }

                            if(attr.field === 'mixins') {
                                var tempValue = [];
                                if(!isNone(attr.value)) {
                                    for (var k = 0; k < attr.value.length; k++) {

                                        var mixinKeys = Ember.keys(attr.value[k]);
                                        var newMixinDict = {}

                                        for (var l = 0; l < mixinKeys.length; l++) {
                                            newMixinDict[mixinKeys[l]] = attr.value[k][mixinKeys[l]];
                                        }
                                        window.$M = newMixinDict;
                                        tempValue.push(newMixinDict);
                                    }
                                }

                                Ember.set(attr, 'value', tempValue);
                            }

                            set(this, 'formContext.' + categoryKeyField, attr.value);
                        }
                    }

                    console.log('this is a widget', get(this, 'formContext'));

                    var args = [get(this, 'formContext')];
                    args.addObjects(arguments);

                    this._super.apply(this, args);
                }
            }
        },
        formOptions);

        application.register('form:modelform', form);
    }
});
