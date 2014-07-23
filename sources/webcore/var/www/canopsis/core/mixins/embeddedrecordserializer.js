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

	function keyForRelationship(key, kind) {
		void (kind);

		key = key.decamelize();
		return key;
	}

	Application.EmbeddedRecordSerializerMixin = Ember.Mixin.create({
		/**
		The current ID index of generated IDs
		@property
		@private
		*/
		_generatedIds: 0,

		/**
		 Sideload a JSON object to the payload

		 @method sideloadItem
		 @param {Object} payload JSON object representing the payload
		 @param {subclass of DS.Model} type The DS.Model class of the item to be sideloaded
		 @paraam {Object} item JSON object representing the record to sideload to the payload
		*/
		sideloadItem: function(payload, type, item, parentJSON) {
			try {
				console.log("sideLoad", type, item.xtype);
				console.log("payload before sideLoad", payload);

				if(item.xtype === undefined) {
					this.addMessage(payload, 'no xtype for widget');
					return payload;
				}
				if(Application[item.xtype.capitalize()] === undefined) {
					this.addMessage(payload, 'bad xtype for widget');
					return payload;
				}

				// The key for the sideload array
				var sideloadKey = item.xtype.pluralize();
				// the ID property key
				var primaryKey = Ember.get(this, 'primaryKey');
				var id = item[primaryKey];

				// The sideload array for this item
				var sideloadArr = payload[sideloadKey] || [];

				console.log("sideloadKey", sideloadKey);

				// Missing an ID, give it one
				if (typeof id === 'undefined') {
					id = 'generated-'+ (++this._generatedIds);
					item[primaryKey] = id;
				}

				if (Ember.isNone(item.meta)) {
					item.meta = {};
				}

				console.log('filling item', item, 'meta with parent', parentJSON.xtype);

				item.meta.embeddedRecord = true;
				item.meta.parentId = parentJSON.id;

				if (!Ember.isNone(parentJSON.xtype)) {
					item.meta.parentType = parentJSON.xtype;
				} else if (!Ember.isNone(parentJSON.crecord_type)) {
					item.meta.parentType = parentJSON.crecord_type;
				}

				// Don't add if already side loaded
				if (!Ember.isNone(sideloadArr.findBy('id', id))) {
					return payload;
				}

				console.log("pushing item", sideloadArr, "into", sideloadKey);

				// Add to sideloaded array
				sideloadArr.push(item);
				payload[sideloadKey] = sideloadArr;

				return payload;
			} catch (e) {
				console.log(e.message, e.stack);
			}
		},

		/**
		 Extract relationships from the payload and sideload them. This function recursively
		 walks down the JSON tree

		 @method sideloadItem
		 @param {Object} payload JSON object representing the payload
		 @param {Object} recordJSON JSON object representing the current record in the payload to look for relationships
		 @param {Object} primaryType The DS.Model class of the record object
		*/
		extractRelationships: function(payload, recordJSON, primaryType, parentType) {
			try {
				console.group('extractRelationships', recordJSON, primaryType);
				console.log("payload before extractRelationships", payload);

				if (primaryType.store === undefined) {
					primaryType.store = parentType.store;
				}

				console.log('primaryType', primaryType, recordJSON.xtype);
				if (primaryType === Application.Widget) {
					var concreteWidgetType = Application[recordJSON.xtype.capitalize()];
					primaryType = concreteWidgetType;
				}

				if (primaryType.store === undefined) {
					primaryType.store = parentType.store;
				}

				primaryType.eachRelationship(function(key, relationship) {
					console.log('eachRelationship', key, recordJSON, recordJSON[key]);
					console.log('relationship', relationship);

					// The record at this relationship
					var related = recordJSON[key],
						// belongsTo or hasMany
						type = relationship.type;

					console.log("related", related);

					if (related) {

						// One-to-one
						if (relationship.kind === 'belongsTo') {
							console.group("belongsTo relationship");
							console.log('related', related);
							console.log('relationship', relationship);
							console.log('recordJSON', recordJSON);
							console.groupEnd();

							// Sideload the object to the payload
							this.sideloadItem(payload, type, related, recordJSON);

							// Replace object with ID
							recordJSON[key] = related.id;
							recordJSON[key + "Type"] = related.xtype;
							// Find relationships in this record
							this.extractRelationships(payload, related, type, primaryType);
						}

						// Many
						else if (relationship.kind === 'hasMany') {
							console.group("hasMany relationship");
							console.log('related', related);
							console.log('relationship', relationship);
							console.log('recordJSON', recordJSON);
							console.groupEnd();

							// Loop through each object
							related.forEach(function(item, index) {
								console.log("sideLoad items in", recordJSON);
								// Sideload the object to the payload
								this.sideloadItem(payload, type, item, recordJSON);

								// Replace object with ID
								related[index] = item.id;

								// Find relationships in this record
								this.extractRelationships(payload, item, type, primaryType);
							}, this);
						}

					}
				}, this);

				console.groupEnd();

				return payload;
			} catch (e) {
				console.log(e.message, e.stack);
			}
		},

		isRecordEmbedded: function(record) {
			var res;
			if (record._data.meta !== undefined && record._data.meta.embeddedRecord === true) {
				res = true;
			} else {
				res = false;
			}
			console.log("isRecordEmbedded", record, res);

			return res;
		},

		getTopmostNotEmbeddedRecordFor: function(record, options) {
			void (options);

			console.log("getTopmostNotEmbeddedRecordFor", record, record.data, record._data);
			var recordCursor = record;

			while(this.isRecordEmbedded(recordCursor)) {
				var parentType = recordCursor._data.meta.parentType;
				var parentId = recordCursor._data.meta.parentId;

				//TODO dynamize
				if (parentType === 'view') {
					parentType = "userview";
				}

				recordCursor = recordCursor.store.getById(parentType, parentId);
			}

			return recordCursor;
		},

		serialize: function(record, options) {
			console.log("serialize", record);
			var lookForDocumentRoot = true;

			if (options !== undefined && options.lookForDocumentRoot === false) {
				lookForDocumentRoot = false;
			}

			if (lookForDocumentRoot && this.isRecordEmbedded(record)) {
				this.getTopmostNotEmbeddedRecordFor(record).save(options);
				return;
			}
			return this._super(record, options);
		},

		serializeBelongsTo: function(record, json, relationship) {
			console.log('serializeBelongsTo', arguments);

			var attr = relationship.key;
			console.log('serializeBelongsTo', relationship, record);

			var key = keyForRelationship(attr, relationship.kind);
			console.log("key", key);

			if (record.get(key) !== undefined && record.get(key) !== null) {
				var serializedSubDocument = record.get(key).serialize({ lookForDocumentRoot : false });
				console.log("serializedSubDocument", serializedSubDocument);
				json[key] = serializedSubDocument;
			}

		},

		serializeHasMany: function(record, json, relationship) {
			console.log('serializeHasMany', arguments);

			var attr = relationship.key;
			console.log('serializeHasMany', relationship, record);

			var key = keyForRelationship(attr, relationship.kind);
			console.log("key", key);

			var subDocuments = record.get(key).get('content');
			console.log("subDocuments", subDocuments);

			json[key] = [];

			for (var i = 0; i < subDocuments.length; i++) {
				if (subDocuments[i] !== undefined && subDocuments[i] !== null) {
					var serializedSubDocument = subDocuments[i].serialize({ lookForDocumentRoot : false });
					serializedSubDocument.xtype = relationship.type.typeKey;
					json[key].push(serializedSubDocument);
				}
			}

			console.log("serializedSubDocuments", json[key]);
		}
	});

	return Application.EmbeddedRecordSerializerMixin;
});
