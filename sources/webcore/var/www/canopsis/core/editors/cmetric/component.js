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
	'utils'
], function(Ember, DS, Application, utils) {

	Application.ComponentCmetricComponent = Ember.Component.extend({
		selectedMetrics: undefined,
		metricSearch: null,

		valueChanged: function() {
			var metrics = this.get('selectedMetrics');
			this.set('content', metrics);

		}.observes('selectedMetrics.@each'),

		helpModal: {
			title: 'Syntax',
			content: "<ul>"
			    + "<li><code>co:regex</code> : look for a component</li>"
			    + "<li><code>re:regex</code> : look for a resource</li>"
			    + "<li><code>me:regex</code> : look for a metric (<code>me:</code> isn't needed for this one)</li>"
			    + "<li>combine all of them to improve your search : <code>co:regex re:regex me:regex</code></li>"
			    + "<li><code>co:</code>, <code>re:</code>, <code>me:</code> : look for non-existant field</li>"
			    + "</ul>",

			id: utils.hash.generateId('cmetric-help-modal'),
			label: utils.hash.generateId('cmetric-help-modal-label')
		},

		select_cols: function() {
			return [
				{name: 'component', title: 'Component'},
				{name: 'resource', title: 'Resource'},
				{name: 'name', title: 'Metric'},
				{
					action: 'select',
					actionAll: 'selectAll',
					title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-plus-sign"></span>'),
					style: 'text-align: center;'
				}
			];
		}.property(),

		unselect_cols: function() {
			return [
				{name: 'component', title: 'Component'},
				{name: 'resource', title: 'Resource'},
				{name: 'name', title: 'Metric'},
				{
					action: 'unselect',
					actionAll: 'unselectAll',
					title: new Ember.Handlebars.SafeString('<span class="glyphicon glyphicon-trash"></span>'),
					style: 'text-align: center;'
				}
			];
		}.property(),

		init: function() {
			this._super(arguments);

			this.set('store', DS.Store.create({
				container: this.get('container')
			}));

			var content = this.get('content') || [];
			this.set('selectedMetrics', content);
		},

		build_filter: function(search) {
			var conditions = search.split(' ');
			var i;

			var patterns = {
				component: [],
				resource: [],
				name: []
			};

			for(i = 0; i < conditions.length; i++) {
				var condition = conditions[i];

				if(condition !== '') {
					var regex = condition.slice(3) || null;

					if(condition.indexOf('co:') === 0) {
						patterns.component.push(regex);
					}
					else if(condition.indexOf('re:') === 0) {
						patterns.resource.push(regex);
					}
					else if(condition.indexOf('me:') === 0) {
						patterns.name.push(regex);
					}
					else {
						patterns.name.push(condition);
					}
				}
			}

			var mfilter = {'$and': []};
			var filters = {
				component: {'$or': []},
				resource: {'$or': []},
				name: {'$or': []}
			};

			for(var key in filters) {
				for(i = 0; i < patterns[key].length; i++) {
					var filter = {};
					var value = patterns[key][i];

					if(value !== null) {
						filter[key] = {'$regex': value};
					}
					else {
						filter[key] = null;
					}

					filters[key]['$or'].push(filter)
				}

				var len = filters[key]['$or'].length;

				if(len === 1) {
					filters[key] = filters[key]['$or'][0];
				}

				if(len > 0) {
					mfilter['$and'].push(filters[key]);
				}
			}

			if(mfilter['$and'].length === 1) {
				mfilter = mfilter['$and'][0];
			}

			return mfilter;
		},

		actions: {
			select: function(metric) {
				var selected = this.get('selectedMetrics');

				if (selected.indexOf(metric) < 0) {
					console.log('Select metric:', metric);
					selected.pushObject(metric);
				}

				this.set('selectedMetrics', selected);
			},

			unselect: function(metric) {
				var selected = this.get('selectedMetrics');

				var idx = selected.indexOf(metric);

				if (idx >= 0) {
					console.log('Unselect metric:', metric);
					selected.removeAt(idx);
				}

				this.set('selectedMetrics', selected);
			},

			selectAll: function() {
				var store = this.get('store');
				var metrics = store.findAll('metric');

				this.set('selectedMetrics', metrics);
			},

			unselectAll: function() {
				this.set('selectedMetrics', []);
			},

			search: function(search) {
				if(search) {
					var mfilter = this.build_filter(search);
					this.set('metricSearch', mfilter);
				}
				else {
					this.set('metricSearch', null);
				}
			}
		},

		didInsertElement: function() {
			$('#' + this.get('helpModal.id')).popover({trigger: 'hover'});
		}
	});

	return Application.ComponentCmetricComponent;
});
