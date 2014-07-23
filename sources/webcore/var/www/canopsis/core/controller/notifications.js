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
	'app/application',
	'utils',
	'app/lib/loaders/schema-manager'
], function(Ember, Application, utils) {

	Application.NotificationsRoute = Ember.Route.extend({
		model: function() {
			return this.store.find('notification');
		}
	});

	//Model initialization
	Application.NotificationAdapter = DS.FixtureAdapter.extend();
	Application.Notification.FIXTURES = [];

	Application.NotificationsController = Ember.ArrayController.extend({

		init: function () {
			this._super();
			utils.notification.controller = this;
		},

		createNotification: function (level, message) {

			var falevel = level;

			if (message === undefined || level === undefined) {
				message = 'missing information for notification';
				falevel = 'warning';
				level = 'warning';
			}
			if (level === 'error') {
				falevel = 'warning';
				level = 'danger';
			}
			var notification = this.store.createRecord('notification',{
				level: level,
				message: message,
				timestamp: new Date().getTime(),
				falevel: 'fa-' + falevel
			});

			notification.save();
		},

		count: function () {
			return this.get('content.length');
		}.property('content.length')
	});


	return Application.NotificationsController;
});