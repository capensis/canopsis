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
    name: 'EmbeddedRecordSerializerMixin',
    after: ['MixinFactory', 'DataUtils', 'HashUtils'],
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var schemasregistry = window.schemasRegistry;

        var hashUtils = container.lookupFactory('utility:hash');
        var dataUtils = container.lookupFactory('utility:data');


        var isNone = Ember.isNone,
            get = Ember.get;

        function keyForRelationship(key) {
            key = key.decamelize();
            return key;
        }

        /**
         * Test routine, every widget should implement this, as userview and widgetwrapper
         *
         * var c = getCanopsis();
         * var widgets = c.registries.widgets;
         *
         * for(var i = 0; i<widgets.length; i++) {
         *    var n = widgets[i].name;
         *    console.log(n.capitalize());
         *    var s = c.Application[n.capitalize() + "Serializer"].create();
         *    console.log(c.Application.EmbeddedRecordSerializerMixin.detect(s));
         *    if(!c.Application.EmbeddedRecordSerializerMixin.detect(s)) alert("stop");
         * }
         *
         * var s = c.Application.UserviewSerializer.create();
         * if(!c.Application.EmbeddedRecordSerializerMixin.detect(s)) alert("stop");
         * var s = c.Application.WidgetwrapperSerializer.create();
         * if(!c.Application.EmbeddedRecordSerializerMixin.detect(s)) alert("stop");
         */

        var mixin = Mixin('embeddedRecordSerializer', {
            /**
             * Sideload a JSON object to the payload
             *
             * @method sideloadItem
             * @param {Object} payload JSON object representing the payload
             * @param {DS.Model} type The DS.Model class of the item to be sideloaded
             * @param {Object} item JSON object representing the record to sideload to the payload
             */
            sideloadItem: function(payload, type, item, parentJSON) {
                try {
                    console.log("sideLoad", type, item.xtype);
                    console.log("payload before sideLoad", payload);

                    if(schemasregistry.getByName(item.xtype).EmberModel === undefined) {
                        console.error(payload, 'bad xtype for widget :' + item.xtype);
                        return undefined;
                    }

                    // The key for the sideload array
                    var sideloadKey = item.xtype.pluralize();
                    // the ID property key
                    var primaryKey = Ember.get(this, 'primaryKey');
                    var id = item[primaryKey];

                    // The sideload array for this item
                    var sideloadArr = payload[sideloadKey] || [];

                    console.log("sideloadKey", sideloadKey);

                    if (isNone(id)) {
                        console.log('generateId', item.xtype, item.id);
                        id = hashUtils.generateId('item');
                        item[primaryKey] = id;
                    }

                    if (isNone(item.meta)) {
                        item.meta = {};
                    }

                    console.log('filling item', item, 'meta with parent', parentJSON.xtype);

                    item.meta.embeddedRecord = true;
                    item.meta.parentId = parentJSON.id;

                    if (!isNone(parentJSON.xtype)) {
                        item.meta.parentType = parentJSON.xtype;
                    } else if (!isNone(parentJSON.crecord_type)) {
                        item.meta.parentType = parentJSON.crecord_type;
                    }

                    // Don't add if already side loaded
                    if (!isNone(sideloadArr.findBy('id', id))) {
                        return payload;
                    }

                    console.log("pushing item", sideloadArr, "into", sideloadKey);

                    // Add to sideloaded array
                    sideloadArr.push(item);
                    payload[sideloadKey] = sideloadArr;

                    return payload;
                } catch (e) {
                    console.error(e.message, e.stack);
                }
            },

            /**
             * Extract relationships from the payload and sideload them. This function recursively
             * walks down the JSON tree
             *
             * @method extractRelationships
             * @param {Object} payload JSON object representing the payload
             * @param {Object} recordJSON JSON object representing the current record in the payload to look for relationships
             * @param {Object} primaryType The DS.Model class of the record object
             * @param {Object} parentType
             */
            extractRelationships: function(payload, recordJSON, primaryType, parentType) {
                console.group('extractRelationships', recordJSON, primaryType);

                try {
                    console.log("payload before extractRelationships", payload);

                    console.log('primaryType', primaryType, recordJSON.xtype);
                    if (primaryType === schemasregistry.getByName('widget').EmberModel) {
                        var concreteWidgetType = schemasregistry.getByName(recordJSON.xtype).EmberModel;
                        primaryType = concreteWidgetType;
                    }

                    if (isNone(primaryType.store)) {
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
                                var sideloadItemResult = this.sideloadItem(payload, type, related, recordJSON);

                                if(!isNone(sideloadItemResult)) {
                                    // Replace object with ID
                                    recordJSON[key] = related.id;
                                    recordJSON[key + "Type"] = related.xtype;
                                    // Find relationships in this record
                                    this.extractRelationships(payload, related, type, primaryType);
                                }
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
                                    var sideloadItemResult = this.sideloadItem(payload, type, item, recordJSON);

                                    if(!isNone(sideloadItemResult)) {
                                        // Replace object with ID
                                        related[index] = item.id;

                                        // Find relationships in this record
                                        this.extractRelationships(payload, item, type, primaryType);
                                    }
                                }, this);
                            }

                        }
                    }, this);

                    console.groupEnd();

                    return payload;
                } catch (e) {
                    console.groupEnd();
                    console.error(e.message, e.stack);
                }
            },

            isRecordEmbedded: function(record) {
                var res;
                if (!isNone(record._data.meta) && record._data.meta.embeddedRecord === true) {
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

                if (!isNone(options) && options.lookForDocumentRoot === false) {
                    lookForDocumentRoot = false;
                }

                if (lookForDocumentRoot && this.isRecordEmbedded(record)) {
                    this.getTopmostNotEmbeddedRecordFor(record).save(options);
                    return;
                }

                var res = this._super(record, options);
                res['id'] = record.get('id');

                console.log(' - serialized record', record, res);

                return res;
            },

            serializeBelongsTo: function(record, json, relationship) {
                console.log('serializeBelongsTo', arguments);

                var attr = relationship.key;
                console.log('serializeBelongsTo', relationship, record);

                var key = keyForRelationship(attr);
                console.log("key", key);

                var content = record.get(key);

                console.log('content', content);

                if (!isNone(content)) {
                    var serializedSubDocument = content.serialize({ lookForDocumentRoot : false });
                    console.log("serializedSubDocument", serializedSubDocument);
                    json[key] = serializedSubDocument;
                } else {
                    console.error('Content is empty', content);
                }

            },

            serializeHasMany: function(record, json, relationship) {

                var attr = relationship.key;
                console.log('serializeHasMany', relationship, record, json, attr);

                var key = keyForRelationship(attr);
                console.log("keyForRelationship", key);

                var subDocuments = record.get(key).get('content');
                console.log("subDocuments", subDocuments);

                json[key] = [];

                var sublen = subDocuments.length;

                if (isNone(relationship) || isNone(relationship.type) || isNone(relationship.type.typeKey)) {
                    console.error('Error while retrieving typeKey from model is it is none.');
                }

                for (var i = 0, l = sublen; i < l; i++) {
                    if (!isNone(subDocuments[i])) {
                        var serializedSubDocument = subDocuments[i].serialize({ lookForDocumentRoot : false });
                        serializedSubDocument.xtype = relationship.type.typeKey;
                        json[key].push(serializedSubDocument);
                    } else {
                        console.error('subdocument is none for index', i, subDocuments);
                    }
                }

                console.log("serializedSubDocuments", json[key]);
            },

            /**
             * @method normalize
             *
             * @todo check if useful
             *
             * @param typeClass
             * @param hash
             * @param prop
             */
            normalize: function(typeClass, hash, prop) {
                typeClass.eachRelationship(function(key, relationship) {
                    if(hash[key]) {
                        hash[key].type = hash[key].xtype;
                    }
                });
                return this._super.apply(this, arguments);
            },

            /**
             * @method pushPayload
             *
             * Used in record duplications
             *
             * @param {DS.Store} store
             * @param {Object} payload
             */
            pushPayload: function(store, payload) {
                var payload = dataUtils.cleanJSONIds(payload);

                var payloadKeys = Ember.keys(payload);

                for (var i = 0, l = payloadKeys.length; i < l; i++) {
                    var currentKey = payloadKeys[i];
                    var typeClass = schemasregistry.getByName(currentKey).EmberModel;
                    this.extractRelationships(payload, payload[currentKey], typeClass);
                }

                payload[currentKey].id = hashUtils.generateId(payload[currentKey].xtype || payload[currentKey].crecord_type || 'item');

                return this._super(store, payload);
            }
        });


        application.register('mixin:embedded-record-serializer', mixin);
    }
});
