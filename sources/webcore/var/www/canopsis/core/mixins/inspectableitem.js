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
	'jquery',
	'ember',
	'app/application'
], function($, Ember, Application) {

	/**
	 * Provides an "attributes" property, dependant on content, to iterate on model's attributes, with the value and schema's properties
	 *
	 * Warning :the parent controller MUST have attributesKeys property!
	 * @mixin
	 */

	Application.InspectableItemMixin = Ember.Mixin.create({

		/**
			@required
		*/
		inspectedDataItem: function() {
			console.error("This must be defined on the base class. Assuming inspected data is content");

			return "content";
		}.property(),

		/**
			@required
		*/
		inspectedItemType: function() {
			console.error("This must be defined on the base class. Assuming inspected data is content.xtype");

			return "content.xtype";
		}.property(),

		/**
			@required
		*/
		inspectedItemInstance: function() {
			console.error("Not mandatory, but attr.value field will not be set");

			return "content";
		}.property(),

		//getting attributes (keys and values as seen on the form)
		categorized_attributes: function() {
			console.log("recompute categorized_attributes", this.get('inspectedDataItem'));
			if (this.get('inspectedDataItem') !== undefined) {
				console.log("inspectedDataItem attributes", this.get('inspectedDataItem').get('attributes'));

				var me = this;

				if (this.get('inspectedItemType') !== undefined) {
					var itemType;

					if (this.get('inspectedItemType') === "view") {
						itemType = "userview";
					} else {
						itemType = this.get('inspectedItemType');
					}

					console.log("inspected itemType", itemType.capitalize());
					var referenceModel = Application[itemType.capitalize()];

					if (referenceModel === undefined || referenceModel.proto() === undefined) {
						console.error("There does not seems to be a registered schema for", itemType.capitalize());
					}
					if (referenceModel.proto().categories === undefined) {
						console.error("No categories in the schema of", itemType);
					}

					var options = this.get('options');
					var filters = [];

					//Allows showing only some fields in the form.
					if (options && options.filters) {
						filters = options.filters;
					}
					console.log(' + filters ', filters);

					//Enables field label override in form from options.
					var override_labels = {};
					if (options && options.override_labels) {
						override_labels = options.override_labels;
					}

					this.categories = [];

					var modelAttributes = Ember.get(referenceModel, 'attributes');

					for (var i = 0; referenceModel.proto().categories &&
					     i < referenceModel.proto().categories.length; i++) {
						var category = referenceModel.proto().categories[i];
						var createdCategory = [];
						createdCategory.title = category.title;
						createdCategory.keys = [];

						for (var j = 0; j < category.keys.length; j++) {
							var key = category.keys[j];

							if (typeof key === "object") {
								key = key.field;
							}

							if (key !== undefined && modelAttributes.get(key) === undefined) {
								console.error("An attribute that does not exists seems to be referenced in schema categories", key, referenceModel);
							}

							//TODO refactor the 20 lines below in an utility function "getEditorForAttr"
							//find appropriate editor for the model property
							var editorName;
							var attr = modelAttributes.get(key);
							console.log("attr", attr, key);

							//defines an option object explicitely here for next instruction
							if (attr.options === undefined) {
								attr.options = {};
							}

							//hide field if not filter specified or if key match one filter element.
							if (filters.length === 0 || $.inArray(key, filters) !== -1) {
								Ember.set(attr, 'options.hiddenInForm', false);
							} else {
								Ember.set(attr, 'options.hiddenInForm', true);
							}

							if (attr.options !== undefined && attr.options.role !== undefined) {
								editorName = "editor-" + attr.options.role;
							} else {
								editorName = "editor-" + attr.type;
							}

							if (Ember.TEMPLATES[editorName] === undefined) {
								editorName = "editor-defaultpropertyeditor";
							}

							//enable field label override.
							var label = key;
							if (override_labels[key]) {
								label = override_labels[key];
							}

							createdCategory.keys[j] = {
								field: label,
								model: modelAttributes.get(key),
								editor: editorName
							};

							if (me.get('inspectedDataItem') !== undefined) {
								createdCategory.keys[j].value = me.get('inspectedDataItem').get(key);
							} else {
								createdCategory.keys[j].value = undefined;
							}

							console.log("category key ", category.keys[j].value);
						}

						this.categories.push(createdCategory);
					}

					console.log("categories", this.categories);
					return this.categories;
				}
				else {
					return undefined;
				}
			}
		}.property("inspectedDataItem", "inspectedItemType")
	});

	return Application.InspectableItemMixin;
});
