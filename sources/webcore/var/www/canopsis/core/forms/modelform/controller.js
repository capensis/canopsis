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
	'app/application',
	'app/lib/factories/form',
	'app/mixins/inspectableitem',
	'app/mixins/validation',
	'app/lib/loaders/schema-manager'
], function(Ember, Application, FormFactory, InspectableitemMixin, ValidationMixin) {

	var formOptions = {
		mixins: [
			InspectableitemMixin,
			ValidationMixin
		]
	};

	FormFactory('modelform', {

  		validationFields: Ember.computed(function() {return Ember.A();}),
  		ArrayFields: Ember.A(),

  		onePageDisplay: function () {
  			//TODO search this value into schema
  			return false;
  		}.property(),

		inspectedDataItem: function() {
			return this.get('formContext');
		}.property('formContext'),

		inspectedItemType: function() {
			console.log('recompute inspectedItemType', this.get('formContext'));

			if (this.get('formContext.xtype')) {
				return this.get('formContext.xtype');
			} else {
				return this.get('formContext.crecord_type');
			}

		}.property('formContext'),

		actions: {
			submit: function() {
				if (this.validation !== undefined && !this.validation()) {
					return;
				}

				console.log("submit action");

				var	newRecord = {};
				var override_inverse = {};

				//will execute callback from options if any given
				var options = this.get('options');

				if (options && options.override_labels) {
					for (var key in options.override_labels) {
						override_inverse[options.override_labels[key]] = key;
					}
				}

				var	categories = this.get("categorized_attributes");

				console.log("setting fields");

				for (var i = 0; i < categories.length; i++) {
					var category = categories[i];
					for (var j = 0; j < category.keys.length; j++) {
						var attr = category.keys[j];
						var field = attr.field;
						//set back overried value to original field
						if (override_inverse[attr.field]) {
							field = override_inverse[attr.field];
						}
						newRecord[field] = attr.value;
						this.set('formContext.' + field, attr.value);
					}
				}
				//Update value of array
				var ArrayFields = this.get("ArrayFields");
				if (ArrayFields !== undefined) {
					for (var w = 0; w < this.ArrayFields.length; w++) {
						console.log("ArrayFields  : ", this.ArrayFields[w]);
						this.ArrayFields[w].onUpdate();
					}
				}

				console.log("this is a widget", this.get('formContext'));
				this._super(this.get('formContext'));

			}
		}
	},
	formOptions);

	return Application.ModelformController;
});
