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

	Application.EntityAdapter = DS.RESTAdapter.extend({
		gen_resolve: function(callback) {
			return function(data) {
				for (var i = 0; i < data.data.length; i++) {
					data.data[i].id = data.data[i]._id;
					delete data.data[i]._id;
				}

				Ember.run(null, callback, data);
			};
		},

		gen_reject: function(callback) {
			return function(xhr) {
				xhr.then = null;
				Ember.run(null, callback, xhr);
			};
		},

		buildURL: function(type, id) {
			return '/entities/' + type + (id ? ('/' + id) : '');
		},

		createRecord: function() {
			Canopsis.utils.notification.error('Impossible to create entity');
		},

		updateRecord: function() {
			Canopsis.utils.notification.error('Impossible to update entity');
		},

		deleteRecord: function() {
			Canopsis.utils.notification.error('Impossible to delete entity');
		},

		find: function(store, model, id) {
			void(store);
			var me = this;

			return new Ember.RSVP.Promise(function(resolve, reject) {
				var url = me.buildURL(model.typeKey, id);
				var funcres = me.gen_resolve(resolve);
				var funcrej = me.gen_reject(reject);

				$.get(url).then(funcres, funcrej);
			});
		},

		findMany: function(store, model, ids) {
			void(store);
			var me = this;

			return new Ember.RSVP.Promise(function(resolve, reject) {
				var funcres = me.gen_resolve(resolve);
				var funcrej = me.gen_reject(reject);

				var mfilter = {'type': model.typeKey};

				if (type === 'downtime') {
					mfilter['downtime_id'] = {'$in': ids};
				}
				else {
					mfilter['name'] = {'$in': ids};
				}

				mfilter = JSON.stringify(mfilter);

				$.post('/entities/', {filter: mfilter}).then(funcres, funcrej);
			});
		},

		findAll: function(store, model, options) {
			void(store);
			var me = this;

			return new Ember.RSVP.Promise(function(resolve, reject) {
				var funcres = me.gen_resolve(resolve);
				var funcrej = me.gen_reject(reject);

				var promise = undefined;

				if (options && options.mfilter) {
					var mfilter = JSON.stringify({
						'$and': [
							{'type': model.typeKey},
							options.mfilter
						]
					});

					promise = $.post('/entities/', {filter: mfilter});
				}
				else {
					var url = me.buildURL(model.typeKey);

					promise = $.get(url);
				}

				promise.then(funcres, funcrej);
			});
		},

		findQuery: function(store, model, query) {
			void(store);
			var me = this;

			return new Ember.RSVP.Promise(function(resolve, reject) {
				var funcres = me.gen_resolve(resolve);
				var funcrej = me.gen_reject(reject);

				if('filter' in query) {
					query.filter.type = model.typeKey;
					query.filter = JSON.stringify(query.filter);
				}
				else {
					query.filter = JSON.stringify({'type': model.typeKey});
				}

				console.log('findQuery: ', query);

				$.post('/entities/', query).then(funcres, funcrej);
			});
		}
	});

	return Application.EntityAdapter;
});
