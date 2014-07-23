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
	'app/routes/application',
	'utils'
	], function(Ember, Application, utils) {

	Application.ApplicationController = Ember.ObjectController.extend({
		needs: ['login'],

		username: function () {
			return this.get('controllers.login').get('username');
		}.property('controllers.login.username'),

		actions: {
			toggleEditMode: function() {
				if (Canopsis.editMode === true) {
					console.info('Entering edit mode');
					Ember.set('Canopsis.editMode', false);
				} else {
					console.info('Leaving edit mode');
					Ember.set('Canopsis.editMode', true);
				}
			},

			addNewView: function () {
				var type = "userview";
				var applicationController = this;
				console.log("add", type);

				var containerwidgetId = utils.hash.generateId('container');

				var containerwidget = Canopsis.utils.data.getStore().createRecord('vbox', {
					xtype: 'vbox',
					id: containerwidgetId
				});

				var userview = Canopsis.utils.data.getStore().push(type, {
					id: utils.hash.generateId('userview'),
					crecord_type: 'view',
					containerwidget: containerwidgetId,
					containerwidgetType: 'vbox'
				});

				console.log('temp record', userview);

				var recordWizard = Canopsis.utils.forms.show('modelform', userview, { title: "Add " + type });

				function transitionToView(userview) {
					console.log('userview saved, switch to the newly created view');
					applicationController.send('showView', userview.get('id'));
				}

				recordWizard.submit.done(function() {
					console.log("userview.save()");
					console.log(userview.save());
					userview.save().then(transitionToView);
				});
			},

			openTab: function(url) {
				this.transitionToRoute(url);
			},

			addModelInstance: function(type) {
				console.log("add", type);

				var record = Canopsis.utils.data.getStore().createRecord(type, {
					crecord_type: type
				});
				console.log('temp record', record);

				var recordWizard = Canopsis.utils.forms.show('modelform', record, { title: "Add " + type });

				recordWizard.submit.done(function() {
					record.save();
				});
			},

			logout: function() {
				this.get('controllers.login').setProperties({
					'authkey': null,
					'errors': []
				});

				window.location = '/logout';
			}
		}

	});

	void (utils);
	return Application.ApplicationController;
});