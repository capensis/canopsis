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

	Application.JournalController = Ember.Controller.extend({
		needs: ['login'],

		init: function() {
			this._super();
		},

		actions: {
			publish: function(action, component) {
				var user = this.get('controllers.login').get('username');

				var ev = {
					timestamp: Date.now() / 1000,
					connector: 'canopsis',
					connector_name: 'canopsis-ui',
					event_type: 'uiaction',

					source_type: 'component',
					component: component,

					author: user,
					action: action,

					state: 0,
					state_type: 1,

					output: user + ' did ' + action + ' on ' + component
				};

				//FIXME and add event controller to needs
				// var eventController = this.get('controllers.event');
				// eventController.send('send_event', ev);
			}
		}
	});

	return Application.JournalController;
});