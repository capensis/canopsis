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
    name: 'MailValidator',
    initialize: function(container, application) {

        function mailValidator(attr, valideStruct) {
            var regex = /^(([^<>()[\]\\.,;:\s@\"]+(\.[^<>()[\]\\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/;
            if (Ember.isBlank(attr.value)) {
                if(attr.model.options.required) {
                    valideStruct.valid = false;
                    valideStruct.error = 'Mail is required';
                }
                else {
                    valideStruct.valid = true;
                }
            }
            else if (regex.test(attr.value)) {
                valideStruct.valid = true ;
            } else {
                valideStruct.valid = false ;
                valideStruct.error = "Mail's format should be: X@Y.Z";
            }

            return valideStruct;
        }

        application.register('validator:mail', mailValidator);
    }
});
