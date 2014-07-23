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
	'app/application',
	'app/mixins/arraymixin',
], function($, Ember, Application) {

	Application.ComponentArrayComponent = Ember.Component.extend(Application.ArrayMixin,{
		valueRefPath:"content.value",
		valuePath:"value",

		init: function() {
			this._super.apply(this, arguments);
			var value = this.getValue();
			var values = this.get('value');

			var me = this;

			Ember.run(function() {
				me.set('arrayAttributes', Ember.A());

				if(values !== undefined) {
					for (var i = 0; i < values.length; i++) {
						me.get('arrayAttributes').pushObject(me.generateVirtualAttribute(i));
					};
				}
			});

		},

		didInsertElement: function() {
			var sortableElements = this.$(".sortable");
			var componentArrayComponent = this;

			if(sortableElements.length >= 0) {
				this.$(".sortable").sortable({
					items: '> tr',
					handle: ".handle",
					axis: "y",
					// forceHelperSize: true,
					// forcePlaceholderSize: true,
					helper: function(e, tr) {
						void (e);
						var $originals = tr.children();
						var $helper = tr.clone();
						$helper.children().each(function(index) {
							// Set helper cell sizes to match the original sizes
							$(this).width($originals.eq(index).width());
						});
						return $helper;
					},
					start: function(event, ui) {
						// creates a temporary attribute on the element with the old index
						$(ui.item.context).attr('data-previndex', $(ui.item.context).index('tr'));
					},
					update: function(event, ui) {
						console.log('update', ui.item);
						console.log('update', $(ui.item.context));
						var newIndex = $(ui.item.context).index('tr');
						var oldIndex = parseInt($(ui.item.context).attr('data-previndex'), 10);

						ui.item.remove();
						componentArrayComponent.moveItem(oldIndex, newIndex);
					}
				});
			}
		},

		itemEditorType: function(){
			var type = this.get('content.model.options.items.type');
			var role = this.get('content.model.options.items.role');
			console.log('editorType', this.get('content'), type, role);
			var editorName;

			if (role) {
				editorName = 'editor-' + role;
			} else {
				editorName = 'editor-' + type;
			}

			if (Ember.TEMPLATES[editorName] === undefined) {
				editorName = 'editor-defaultpropertyeditor';
			}

			return editorName;
		}.property('content.model.options.items.type', 'content.model.options.items.role'),

		generateVirtualAttribute: function(itemIndex) {
			var values = this.get('value');
			var content = this.get('content.model.options.items');
			var componentArrayComponent = this;


			console.log('generateVirtualAttribute', itemIndex, values[itemIndex]);

			var currentGeneratedAttr = Ember.Object.create({
				parent: componentArrayComponent,
				index: itemIndex,
				value : values[itemIndex]
			});

			//apply options to virtual attribute
			Ember.set(currentGeneratedAttr, 'model', Ember.Object.create());
			Ember.set(currentGeneratedAttr, 'model.options', Ember.Object.create());

			for (var key in content) {
				if (key !== 'value') {
					Ember.set(currentGeneratedAttr, 'model.options.' + key, content[key]);
				}
			}
			console.log("generateVirtualAttribute virtual attribute", currentGeneratedAttr);

			Ember.addObserver(currentGeneratedAttr, 'value', function(attr) {
				console.log('value changed', attr.value, attr.index);
				componentArrayComponent.get('value').removeAt(attr.index);
				componentArrayComponent.get('value').insertAt(attr.index, attr.value);
				console.log('content changed', componentArrayComponent.get('value'));
			});

			console.log("generateVirtualAttribute@3");

			return currentGeneratedAttr;
		},

		contentChanged: function() {
			console.log('recompute contentAttributeArray');
		},

		actions: {
			addItem: function() {
				console.log('addItem', this.get('value'));

				if(this.get('content.model.options.items.type') === 'object') {
					this.get('value').pushObject(Ember.Object.create(this.get('content.model.options.items.objectDict')));
				} else {
					this.get('value').pushObject(undefined);
				}
				var newIndex = this.get('value').length -1;
				this.get('arrayAttributes').pushObject(this.generateVirtualAttribute(newIndex));
			},
			removeItem: function(item) {
				console.log('removeItem', this.get('value'));

				var arrayAttributes = this.get('arrayAttributes');
				this.get('value').removeAt(item.index);
				arrayAttributes.removeAt(item.index);
				for (var i = item.index; i < arrayAttributes.length; i++) {
					arrayAttributes.objectAt(i).set('index', arrayAttributes.objectAt(i).get('index') - 1);
				}
			}
		},
		moveItem: function(oldIndex, newIndex) {
			console.log('moveItem', arguments);
			// this.get('arrayAttributes').moveElement(oldIndex, newIndex);

			//update virtual attributes
			var array = this.get('arrayAttributes');
			array.arrayContentWillChange(oldIndex, 1, 0);
			var oldIndex_value = array[oldIndex];
			array.splice(oldIndex, 1);
			array.splice(newIndex, 0, oldIndex_value);

			for (var i = 0; i < array.length; i++) {
				array[i].index = i;
			};

			array.arrayContentDidChange(newIndex, 0, 1);

			//update component value
			var array = this.get('value');
			array.arrayContentWillChange(oldIndex, 1, 0);
			var oldIndex_value = array[oldIndex];
			array.splice(oldIndex, 1);
			array.splice(newIndex, 0, oldIndex_value);

			array.arrayContentDidChange(newIndex, 0, 1);

		}
	});

	return Application.ComponentArrayComponent;
});