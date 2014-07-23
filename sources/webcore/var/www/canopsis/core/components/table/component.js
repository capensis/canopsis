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
	'app/application',
	'app/mixins/pagination'
], function(Ember, DS, Application, PaginationMixin) {

	Application.ComponentTableComponent = Ember.Component.extend({
		model: undefined,
		modelfilter: undefined,
		data: undefined,

		columns: [],
		items: [],

		onDataChange: function() {
			this.refreshContent();
		}.observes('data.@each'),

		onModelFilterChange: function() {
			this.set('currentPage', 1);
			this.refreshContent();
		}.observes('modelfilter'),

		init: function() {
			this._super(arguments);

			if (this.get('model') !== undefined) {
				this.set('store', DS.Store.create({
					container: this.get('container')
				}));
			}
		},

		didInsertElement: function() {
			this.refreshContent();
		},

		refreshContent: function() {
			this._super(arguments);

			this.findItems();

			console.log(this.get('widgetDataMetas'));
		},

		findItems: function() {
			var me = this;

			var store = this.get('store');

			var query = {
				start: this.get('paginationMixinFindOptions.start'),
				limit: this.get('paginationMixinFindOptions.limit')
			};

			if (this.get('model') !== undefined) {			
				if(this.get('modelfilter') !== null) {
					query.filter = this.get('modelfilter');
				}

				store.findQuery(this.get('model'), query).then(function(result) {
					me.set('widgetDataMetas', result.meta);
					me.set('items', result.get('content'));

					me.extractItems(result);
				});
			}
			else {
				var items = this.get('data').slice(
					query.start,
					query.start + query.limit
				);

				this.set('widgetDataMetas', {total: this.get('data').length});
				this.set('items', items);

				me.extractItems({
					meta: this.get('widgetDataMetas'),
					content: this.get('items')
				});
			}
		},

		actions: {
			do: function(action, item) {
				this.targetObject.send(action, item);
			}
		}
	});

	Application.ComponentTableComponent = Application.ComponentTableComponent.extend(PaginationMixin);

	return Application.ComponentTableComponent;
});