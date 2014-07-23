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
	Application.ComponentActionfilterComponent = Ember.Component.extend({


		init: function() {
			this._super();
			//default value on load
			this.set('selectedAction', 'pass');
			console.log(' ! --- > content', this.get('content'));
			if (this.get('content') === undefined) {
				this.set('content', []);
			}
		},

		selectedAction: 'pwoute',
		availableactions: ['pass','drop','override','remove'],

		isOverride: function () {
			console.log('isOverride', this.get('selectedAction'), this.get('selectedAction') === 'override');
			return this.get('selectedAction') === 'override';
		}.property('selectedAction'),

		isRoute: function () {
			//not used yet
			return false;
			//console.log('isRoute', this.get('selectedAction'), this.get('selectedAction') === 'route');
			//return this.get('selectedAction') === 'route';
		}.property('selectedAction'),

		isRemove: function () {
			console.log('isRemove', this.get('selectedAction'), this.get('selectedAction') === 'remove');
			return this.get('selectedAction') === 'remove';
		}.property('selectedAction'),


		actions : {
			addAction: function () {
				var action = {
					type: this.get('selectedAction')
				};

				if (this.get('selectedAction') === 'override') {
					action.field = this.get('field');
					action.value = this.get('value');
				}

				if (this.get('selectedAction') === 'remove') {
					action.key = this.get('key');
				}

				console.log('Adding action', action);
				this.get('content').pushObject(action);
			},
			deleteAction: function (action) {
				console.log('Removing action', action);
				this.get('content').removeObject(action);

			}
		}

	});

	return Application.ComponentActionfilterComponent;
});