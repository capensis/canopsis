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
	'ember-data',
	'app/lib/factories/widget',
	'app/mixins/pagination',
	'app/mixins/inspectablearray',
	'app/mixins/arraysearch',
	'app/mixins/sortablearray',
	'app/mixins/history',
	'app/mixins/sendevent',
	'utils',
	'app/lib/loaders/schema-manager',
	'app/adapters/event',
	'app/adapters/userview',
	'canopsis/core/lib/wrappers/ember-cloaking',
	'app/view/listline',
], function(Ember, DS, WidgetFactory, PaginationMixin, InspectableArrayMixin,
		ArraySearchMixin, SortableArrayMixin, HistoryMixin, SendEventMixin, utils) {

	var listOptions = {
		mixins: [
			InspectableArrayMixin,
			ArraySearchMixin,
			SortableArrayMixin,
			PaginationMixin,
			HistoryMixin,
			SendEventMixin
	]};

	var widget = WidgetFactory('list',
		{
			needs: ['login'],

			init: function() {
				this._super();
			},

			/**
			* @return : a list of record from list which checked value matches isSelected paramterter
			* @param : isSelected determines wether or not the record has to be selected
			**/
			getRecordCheckedTo: function (isSelected) {

				//TODO not ready yet, juste return all for now
				return this.get('widgetData.content');
				/*
				var selectedRecords = [];

				console.log('will iterate over', this.get('widgetData.content'), 'witch lenght is ', this.get('widgetData.content').length, 'isSelected value is', isSelected);

				for (var i=0; i<this.get('widgetData.content').length; i++) {
					console.log('crecord tested is', this.get('widgetData.content')[i], 'Selected value is', this.get('widgetData.content')[i].get('isSelected'));
					if (this.get('widgetData.content')[i].isSelected === isSelected) {
						selectedRecords.push(this.get('widgetData.content')[i]);
					}
				}
				return selectedRecords;
				*/
			},



			itemType: function() {
				var listed_crecord_type = this.get("listed_crecord_type");
				console.info('listed_crecord_type', listed_crecord_type);
				if(listed_crecord_type !== undefined && listed_crecord_type !== null ) {
					return this.get("listed_crecord_type");
				} else {
					return 'event';
				}
			}.property("listed_crecord_type"),

			widgetData: [],

			findOptions: {},
			loaded: false,

			itemsPerPagePropositions : [5, 10, 20, 50],
			userDefinedItemsPerPageChanged : function() {
				this.set('itemsPerPage', this.get('userDefinedItemsPerPage'));
				this.refreshContent();
			}.observes('userDefinedItemsPerPage'),

			//Mixin aliases
			//history
			historyMixinFindOptions: Ember.computed.alias("findOptions.useLogCollection"),
			//inspectedDataItemMixin
			inspectedDataArray: Ember.computed.alias("widgetData"),
			itemsPerPage: Ember.computed.alias("content.itemsPerPage"),
			//pagination
			paginationMixinFindOptions: Ember.computed.alias("findOptions"),

			//TODO put this in widget conf
			searchableAttributes: ["firstname", "lastname"],

			onReload: function () {
				this._super();
			},

			onDomReady: function (element) {
				void(element);
			},

			actions: {
				setFilter: function (filter) {
					this.findOptions.filter = filter;

					if (this.currentPage !== undefined) {
						this.set("currentPage", 1);
					}

					this.refreshContent();
				},

				moveColumn: function (attr, direction) {
					console.log('moving', attr, direction);
					var columns = this.get('shown_columns');
					var col;
					for (var i=0; i<columns.length; i++) {
						if (columns[i].field === attr.field) {
							console.log(attr.field +  ' found at position ' + i);
							col = i;
							break;
						}
					}
					if (col !== undefined) {
						if( !(col === 0 && direction === 'left') && !(col === columns.length && direction === 'right')) {
							var permutation;
							if (direction === 'left') {
								permutation = columns[col - 1];
								columns[col - 1] = columns[col];
								columns[col] = permutation;
							} else {
								permutation = columns[col + 1];
								columns[col + 1] = columns[col];
								columns[col] = permutation;
							}
							console.debug('permuting column to ' + direction);
							this.set('userParams.user_show_columns', columns);
							this.get('userConfiguration').saveUserConfiguration();
							this.trigger('refresh');
						} else {
							console.debug('impossible action for colums switch');
						}
					}

				},

				switchColumnDisplay: function (attr) {
					console.log('column switch display', attr);
					console.log('attribute keys', this.get('attributesKeys'));
					Ember.set(attr, 'options.show', !Ember.get(attr, 'options.show'));
					var columns = this.get('shown_columns');
					console.log('Columns on reload', columns);
					this.set('userParams.user_show_columns', columns);
					this.get('userConfiguration').saveUserConfiguration();
				},

				show: function(id) {
					console.log("Show action", arguments);
					utils.routes.getCurrentRouteController().send('showView', id);
				},

				add: function (recordType) {
					console.log("add", recordType);

					var record = this.get("widgetDataStore").createRecord(recordType, {
						crecord_type: recordType
					});
					console.log('temp record', record);

					var recordWizard = utils.forms.show('modelform', record, { title: "Add " + recordType });

					var listController = this;

					recordWizard.submit.then(function(form) {
						console.log('record going to be saved', record, form);

						record = form.get('formContext');

						record.save();

						listController.trigger('refresh');

						listController.startRefresh();
					});
				},

				edit: function (record) {
					console.log("edit", record);

					var listController = this;
					var recordWizard = utils.forms.show('modelform', record, { title: "Edit " + record.get('crecord_type') });

					recordWizard.submit.then(function(form) {
						console.log('record going to be saved', record, form);

						record = form.get('formContext');

						record.save();
						listController.trigger('refresh');
					});
				},

				removeRecord: function(record) {
					console.info('removing record', record);
					record.deleteRecord();
					record.save();
				},

				removeSelection: function() {
					var selected = this.get("widgetData").filterBy('isSelected', true);
					console.log("remove action", selected);
					for (var i = 0; i < selected.length; i++) {
						var currentSelectedRecord = selected[i];
						this.send("removeRecord", currentSelectedRecord);
					}
				}
			},

			findItems: function() {
				var me = this;

				if (this.get("widgetDataStore") === undefined) {
					this.set("widgetDataStore", DS.Store.create({
						container: this.get("container")
					}));
				}

				var itemType = this.get("itemType");

				console.log("findItems", itemType);

				if (itemType === undefined || itemType === null) {
					console.error ("itemType is undefined for", this);
					return;
				}

				console.tags.add('data');
				console.log("find items of type", itemType, "with options", this.get('findOptions'));
				console.tags.remove('data');

				this.get("widgetDataStore").findQuery(itemType, this.findOptions).then(function(queryResults) {
					console.tags.add('data');
					console.log("got results in widgetDataStore", itemType, "with options", me.get('findOptions'));
					console.tags.remove('data');

					//retreive the metas of the records
					me.set("widgetDataMetas", me.get("widgetDataStore").metadataFor(me.get("listed_crecord_type")));
					me.extractItems.apply(me, [queryResults]);
					me.set('loaded', true);

					me.trigger('refresh');
				});
			},

			attributesKeysDict: function() {
				var res = {};
				var attributesKeys = this.get('attributesKeys');
				var sortedAttribute = this.get('sortedAttribute');

				for (var i = 0; i < attributesKeys.length; i++) {
					if (sortedAttribute !== undefined && sortedAttribute.field === attributesKeys[i].field) {
						res[attributesKeys[i].field] = sortedAttribute;
					} else {
						res[attributesKeys[i].field] = attributesKeys[i];
					}
				}

				return res;
			}.property('attributesKeys'),

			shown_columns: function() {
				console.log("compute shown_columns", this.get('sorted_columns'), this.get('attributesKeys'), this.get('sortedAttribute'));
				if (this.get('user_show_columns') !== undefined) {
					console.log('user columns selected', this.get('user_show_columns'));
					return this.get('user_show_columns');
				}

				var shown_columns = [];
				if (this.get('sorted_columns') !== undefined && this.get('sorted_columns').length > 0) {

					var attributesKeysDict = this.get('attributesKeysDict');

					var sorted_columns = this.get('sorted_columns');

					for (var i = 0; i < sorted_columns.length; i++) {
						if (attributesKeysDict[sorted_columns[i]] !== undefined) {
							attributesKeysDict[sorted_columns[i]].options.show = true;
							shown_columns.push(attributesKeysDict[sorted_columns[i]]);
						}
					}
				} else {
					console.log('no shown columns set, displaying everything');
					shown_columns = this.get('attributesKeys');
				}
				var selected_columns = [];
				for(var column=0; column < shown_columns.length; column++) {

					shown_columns[column].options.show = true;
					if ($.inArray(shown_columns[column].field, this.get('hidden_columns')) === -1) {
						selected_columns.push(shown_columns[column]);
					}
				}
				return selected_columns;

			}.property('attributesKeysDict', 'attributesKeys', 'sorted_columns')
	}, listOptions);

	return widget;
});