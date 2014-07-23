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
	'app/mixins/inspectableitem',
	'app/lib/loaders/schema-manager'
], function(Ember, Application, InspectableItem) {

	Application.CrecordController = Ember.ObjectController.extend(InspectableItem, {
		actions: {
			showEditForm: function() {
				var crecord_type = this.get("model.constructor.typeKey");
				console.log("Form generation for", crecord_type);

				var crecordformController = Application.CrecordformController.create();
				crecordformController.set("crecord_type", crecord_type);
				crecordformController.set("editMode", "edit");
				crecordformController.set("editedRecordController", this);

				//Delete old validationFields (should be done on close)
				if (crecordformController.validationFields) {
					while(crecordformController.validationFields.length > 0) {
						crecordformController.validationFields.pop();
					}
				}

				//Delete old ArrayFields (should be done on close)
				if (crecordformController.ArrayFields) {
					while(crecordformController.ArrayFields.length > 0) {
						crecordformController.ArrayFields.pop();
					}
				}

				this.send('showEditFormWithController', crecordformController);
			},

			editRecord: function(record_raw, callback) {
				console.log("editRecord", record_raw);
				this.get("model").setProperties(record_raw);
				var promise = this.get("model").save();
				if (callback !== undefined) {
					promise.then(callback);
				}
			}
		},

		remove: function() {
			this.get("model").deleteRecord();
			this.get("model").save();
		},

		//This is where to get data from the crecord. It should not be changed, and is for internal use only
		dataAccessKey: "content._data",

		inspectedDataItem: function() {
			return this.get('widgetData');
		}.property('widgetData')
	});

	return Application.CrecordController;
});