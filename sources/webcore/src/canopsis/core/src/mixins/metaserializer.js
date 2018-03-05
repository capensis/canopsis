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
    name: 'MetaSerializerMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var isNone = Ember.isNone;

        /**
         * @mixin This mixin should be used with serializers
         * It aims to handle request metadata (total, errors, ...)
         */
        var mixin = Mixin('metaSerializer', {
            extractMeta: function(store, type, payload) {
                console.log("extractMeta", store, type, payload);
                if (payload && payload.meta === undefined) {
                    payload.meta = {};
                }

                if (payload && payload.total !== undefined) {
                    payload.meta.total = payload.total;
                }

                if (payload && payload.messages) {
                    payload.meta.totalmessages = payload.messages;
                }

                if (payload && payload.success !== undefined) {
                    payload.meta.success = payload.success;
                }

                delete payload.total;
                delete payload.messages;
                delete payload.success;

                console.log('normalizePayload', arguments);

                if (isNone(type) || isNone(type.typeKey)) {
                    console.error('Error while retrieving typeKey from model is it is none.');
                }

                var typeKey = type.typeKey,
                    typeKeyPlural = typeKey.pluralize();

                payload[typeKeyPlural] = payload.data;
                delete payload.data;

                console.log('payload', typeKeyPlural, payload);
                // Many items (findMany, findAll)
                if (payload && payload[typeKeyPlural] && typeof payload[typeKeyPlural] !== 'undefined') {
                    payload[typeKeyPlural].forEach(function(item) {
                            this.extractRelationships(payload, item, type);
                    }, this);
                }

                // Single item (find)
                else if (typeof payload[typeKey] !== 'undefined') {
                    this.extractRelationships(payload, payload[typeKey], type);
                }
                return this._super(store, type, payload);
            },

            extractRelationships: function(payload, item, type){
                this._super.apply(this, arguments);
            }

        });

        application.register('mixin:meta-serializer', mixin);
    }
});
