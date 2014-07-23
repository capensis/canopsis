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
	'app/application'
], function(Ember, DS, Application) {

	Application.ApplicationAdapter = DS.RESTAdapter.extend({
		buildURL: function(type, id) {
			return ("/rest/object/" + type + (!!id ? "/" + id : ""));
		},

		createRecord: function(store, type, record) {
			var data = {};
			var serializer = store.serializerFor(type.typeKey);

			data = serializer.serializeIntoHash(data, type, record, "POST", { includeId: true });

			return this.ajax(this.buildURL(type.typeKey, record.id), "POST", { data: data });
		},

		updateRecord: function(store, type, record) {
			var data = {};
			var serializer = store.serializerFor(type.typeKey);

			data = serializer.serializeIntoHash(data, type, record, "PUT");

			var id = Ember.get(record, 'id');

			return this.ajax(this.buildURL(type.typeKey, id), "PUT", { data: data });
		}
	});

	return Application.ApplicationAdapter;
});
