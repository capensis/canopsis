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
    'app/controller/form',
    'app/mixins/inspectableitem',
    'app/lib/loaders/schema-manager'
], function(Ember, Application, FormController, InspectableitemMixin) {

    Application.CrecordformController = FormController.extend(InspectableitemMixin, {
	modalId: 'crecordform-modal',

	validationFields: Ember.A(),
	ArrayFields: Ember.A(),

	inspectedDataItem: function() {
	    return this.get('editedRecordController');
	}.property('editedRecordController'),

	inspectedItemType: function() {
	    return this.get('crecord_type');
	}.property('crecord_type'),

	actions: {
	    submit: function() {
		console.log("addRecord from CrecordformController, editMode:", this.editMode);

		console.log(this);
		console.log("saveRecord edited crecord controller", this.editedRecordController);



		//will execute callback from options if any given

		var options = this.get('options');
		var override_inverse = {};

		//TODO @eric refactor this in InspectableitemMixin?
		if (options && options.override_labels) {
		    for (key in options.override_labels) {
			override_inverse[options.override_labels[key]] = key;
		    }
		}

		var categories = this.get("categorized_attributes");


		//TODO @momo refactor this in ValidationMixin
		// Validation
		var validationFields = this.get("validationFields");
		if (validationFields)
		{
			for (var z = 0; z < validationFields.length; z++) {
			    console.log("validate on : ", validationFields[z]);
			    // Check if a field's validate function return false
			    if (validationFields[z].validate() !== true) {

				console.log("Can't validate on attr ",validationFields[z]);
				// for now just stop and return (fields error messages have been updated)
				return;
			    }
			}
		}

		//TODO @Momo is this still used? if yes, can't we use an inspectable mixin?
		// Array
		var ArrayFields = this.get("ArrayFields");
		if (ArrayFields) {
			for (var w = 0; w < ArrayFields.length; w++) {
			    console.log("ArrayFields  : ", ArrayFields[w]);
			    ArrayFields[w].onUpdate();
			}
		}

		var newRecord;

		for (var i = 0; i < categories.length; i++) {
		    var	category = categories[i];
		    for (var j = 0; j < category.keys.length; j++) {
			var attr = category.keys[j];
			var field = attr.field;
			//set back overried value to original field
			if (override_inverse[attr.field]) {
			    field = override_inverse[attr.field];
			}
			newRecord[field] = attr.value;

		    }
		}

		if (this.editMode === "add") {
		    //sets data to record from default values described in model
		    //sets extra keys from options parameter if any
		    if (options && options.set) {
			for (key in options.set) {
			    newRecord[key] = options.set[key];
			}
		    }
		    console.log(' -> newRecord :', newRecord);
		    //TODO ugly
		    var mainCrecordController = Canopsis.utils.routes.getCurrentRouteController();
		    mainCrecordController.send('addRecord', this.crecord_type, newRecord, options);
		}
		else if (this.editMode === "edit") {
		    Canopsis.utils.routes.getCurrentRouteController().send("editAndSaveModel", this.editedRecordController, newRecord, options);
		}
		else {
		    console.log("bad record form mode");
		}

		//reset editmode to avoid unpredictable behaviour later
		this.editMode = undefined;
		this.trigger("validate");
		this.submit.resolve();
	    }
	}
    });

    return Application.CrecordformController;
});
