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
	'app/application',
	'app/adapters/application',
	'app/serializers/event',
], function(Application, ApplicationAdapter) {
	var adapter = ApplicationAdapter.extend({

		buildURL: function(type, id) {
			return "/event";
		},

		findQuery: function(store, type, query) {
			if (query && query.useLogCollection !== undefined) {
				this.set('useLogCollection', query.useLogCollection);
				delete query.useLogCollection;
			}
			var noAckSearch = false;
			if (query && query.noAckSearch) {
				noAckSearch = true;
				delete query.noAckSearch;
			}
			var log = this.get('useLogCollection') ? '_log' : '';
			console.log('LOG', log);
			var url = "/rest/events" + log;

			return this.ajax(url, 'GET', { data: query });
		}
	});

	Application.EventAdapter = adapter;

	return adapter;
});
