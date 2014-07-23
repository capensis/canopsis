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
	Application.ComponentEventselectorComponent = Ember.Component.extend({

		init: function() {
			this._super();
			this.set("componentDataStore", DS.Store.create({
				container: this.get("container")
			}));
			console.log("Event selector init");

			this.set('selectedEvents', []);

			if (this.get('content') !== undefined) {
				this.initializeEvents();
			}
		},

		initializeEvents: function () {
			rks = this.get('content');

			var that = this;
			var query = this.get("componentDataStore").findQuery(
				'event',
				{
					filter: JSON.stringify({_id: {'$in': rks}}),
					limit: 0,
					noAckSearch: true
				}
			).then(
				function (data) {
					console.log('Fetched initialization data from events', data.content);
					that.set('selectedEvents', data.content);
				}
			);
			void (query);
		},

		findEvents: function() {

			var filter = {};

			var excludeRks = this.getSelectedRks();

			//adding exclusion rks if any loaded
			if (excludeRks.length) {
				filter._id = {'$nin': excludeRks};
			}

			//permissive search throught component and resource
			if (this.get('component')) {
				filter.component = { '$regex' : '.*'+ this.get('component') +'.*', '$options': 'i' };
			}
			if (this.get('resource')) {
				filter.resource = { '$regex' : '.*'+ this.get('resource') +'.*', '$options': 'i' };
			}

			//does user selected selector or topology search
			if (this.get('selectors')) {
				filter.crecord_type = 'selector';
			}

			if (this.get('topologies')) {
				filter.crecord_type = 'topologies';
			}

			if (!filter.resource && !filter.component) {
				this.set('events', []);
				//when user only wants topologies or selectors, query is done anyway with the right crecord type
				if (!this.get('topologies') && !this.get('selectors')) {
					return;
				}
			}

			var query = this.get("componentDataStore").findQuery(
				'event',
				{
					filter: JSON.stringify(filter),
					limit: 10,
					noAckSearch: true
				}
			);

			var that = this;
			query.then(
				function (data) {
					console.log('Fetched data from events', data.content);
					that.set('events', data.content);
			});

			void (query);

		}.observes('component', 'resource'),

		setSelector: function() {
			this.set('topologies', false);
			this.findEvents();
		}.observes('selectors'),

		setTopologies: function() {
			this.set('selectors', false);
			this.findEvents();
		}.observes('topologies'),

		didInsertElement: function() {
		},

		getSelectedRks: function() {
			var selectedEvents = [];
			if (this.get('selectedEvents') !== undefined) {
				for (var i=0; i<this.get('selectedEvents').length; i++) {
					selectedEvents.push(this.get('selectedEvents')[i].id);
				}
			}
			return selectedEvents;
		},

		actions: {

			add: function (event) {
				console.log('Adding event', event);
				this.get('selectedEvents').pushObject(event);
				this.get('events').removeObject(event);
				this.set('content', this.getSelectedRks());
				if (!this.get('events').length) {
					this.findEvents();
				}
			},

			delete: function (event) {
				console.log('Rk to delete', event.id);
				for (var i=0; i<this.get('selectedEvents').length; i++) {
					if (event.id === this.get('selectedEvents')[i].id) {
						console.log('Removing event');
						this.get('selectedEvents').removeAt(i);
						break;
					}
				}
				this.findEvents();
				this.set('content', this.getSelectedRks());
			}
		}
	});

	return Application.ComponentEventselectorComponent;
});