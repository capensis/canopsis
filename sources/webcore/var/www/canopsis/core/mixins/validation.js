/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

define([
	'ember',
	'app/application'
], function(Ember, Application) {

	/**
	 * Implements Validation in form
	 * You should define on the validationFields
	 * @mixin
	 */
	Application.ValidationMixin = Ember.Mixin.create({

		validationFields: function() {
			console.warn("Property \"validationFields\" must be defined on the concrete class.");

			return "<validationFields is null>";
		},

		validation: function() {
			console.log("Enter validation MIXIN");
			var validationFields = this.get("validationFields");
			if (validationFields) {
				for (var z = 0; z < validationFields.length; z++) {
					console.log("validate on : ", validationFields[z]);
					// Check if a field's validate function return false
					if (validationFields[z].validate() !== true) {
						console.log("Can't validate on attr ",validationFields[z]);
						// for now just stop and return (fields error messages have been updated)
						return false ;
					}
				}
			}

			return true;
		}
	});
	return Application.ValidationMixin;
});