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
	'app/routes/authenticated',
	'app/lib/loaders/templates'
], function(Ember, Application, AuthenticatedRoute) {

	Application.ApplicationRoute = AuthenticatedRoute.extend({
		actions: {
			showView: function(id) {
				console.log("ShowView action", arguments);
				this.transitionTo('userview', id);
			},

			showEditFormWithController: function(formController, formContext, options) {
				if (formController.ArrayFields) {
					while(formController.ArrayFields.length > 0) {
						formController.ArrayFields.pop();
					}
				}

				var formName = formController.constructor.toString().slice(1, "Controller".length * -1).toLowerCase();
				console.log("showEditFormWithController", formController, formName, formContext, options);

				var formwrapperController = this.controllerFor('formwrapper');
				Ember.set('Canopsis.formwrapperController', formwrapperController);

				formController.set('formwrapper', formwrapperController);
				formController.set('formContext', formContext);
				formwrapperController.set('form', formController);
				formwrapperController.set('formName', formName);

				formwrapperController.send('show');

				return formController;
			}
		},

		model: function() {
			return {
				title: 'Canopsis'
			};
		},

		renderTemplate: function() {
			this.render();

			//getting the generated controller
			var notificationsController = this.controllerFor('notifications');
			var formwrapperController = this.controllerFor('formwrapper');

			//assigning the model
			notificationsController.set('content', this.store.find("notification"));

			this.render('notifications', {
			    outlet: 'notifications',
			    into: 'application',
			    controller: notificationsController
			});

			this.render('formwrapper', {
			    outlet: 'formwrapper',
			    into: 'application',
			    controller: formwrapperController
			});

		}
	});

	return Application.ApplicationRoute;
});
