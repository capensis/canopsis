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
    name: 'RequiredValidator',
    initialize: function(container, application) {

        function requiredValidator(attr, valideStruct) {
            console.log("requiredValidator :attr = ",attr) ;
            if (attr.model.options.required !== undefined && attr.model.options.required === true  && (attr.value === undefined || attr.value === ""  || attr.value === null)) {
                valideStruct.valid = false;
                valideStruct.error = " Field can't be empty";
            } else {
                valideStruct.valid = true;
                valideStruct.error = "";
            }

            return valideStruct;
        }

        application.register('validator:required', requiredValidator);
    }
});
