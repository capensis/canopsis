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
    name: 'ValidationMixin',
    after: ['MixinFactory', 'FormsUtils'],
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');
        var formUtils = container.lookupFactory('utility:forms');

        var formUtils;

        var get = Ember.get,
            set = Ember.set;

        /**
         * Implements Validation in form
         * You should define on the validationFields
         * @mixin
         */
        var mixin = Mixin('validation', {

            validationFields: function() {
                console.warn("Property \"validationFields\" must be defined on the concrete class.");

                return "<validationFields is null>";
            },

            changeTAB: function( name , active){
                var toFind = "#" + name + "_tab";
                if(active)
                    $(toFind).addClass("active");
                else
                    $(toFind).removeClass("active");


                var id = "#" + name;
                if(active)
                    $(id).addClass("active");
                else
                    $(id).removeClass("active");
            },

            set_tab: function(last_field_error){
                var categories = this.categories,
                    i,
                    current;

                for (i = 0 ; i < categories.length ; i++){
                    current = categories[i];
                    this.changeTAB( current.slug , false );
                    Ember.set(current, "isDefault", false);
                }
        outer:  for (i = 0 ; i < categories.length ; i++){
                    current = categories[i];

                    for (var j = 0 ; j < current.keys.length ; j++){
                        var key = current.keys[j];
                        var field = key.field;

                        if (field === last_field_error ){
                            this.changeTAB( current.slug , true );
                            Ember.set(current, "isDefault", true);

                            break outer;
                        }
                    }
                }
            },

            empty_validationFields: function() {
                set(this, 'validationFields' , Ember.A() );
            },

            validation: function() {
                console.group("Form validation");

                var validationFields = get(this, "validationFields");
                var isValid = true;
                var error_array = [];
                var last_field_error = "";
                var form = this;

                if (validationFields) {
                    for (var z = 0, l = validationFields.length; z < l; z++) {
                        console.log("check if field is valid", get(validationFields[z], 'attr.field'));
                        var current = validationFields[z].validate();

                        if (current.valid !== true) {
                            error_array.push(current);
                            console.log("Attribute not valid", validationFields[z]);
                            last_field_error = validationFields[z].attr.field || validationFields[z].attr.parent.attr.field;
                            isValid =  false ;

                            if (validationFields[z].removedFromDOM){
                                //var form = validationFields[z].get('form');
                                form.validateOnInsert = true;
                                formUtils.showInstance(form);
                                break;
                            }
                        }
                    }
                }
                if( !isValid ){
                    console.log('form not valid');
                    this.set_tab(last_field_error);
                }

                console.groupEnd();
                return isValid;
            }
        });

        application.register('mixin:validation', mixin);
    }
});
