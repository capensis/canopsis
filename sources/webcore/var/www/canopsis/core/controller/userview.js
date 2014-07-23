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
	'app/controller/crecord',
	'app/routes/userview',
	'app/view/userview',
	'app/serializers/userview'
], function(Ember, Application, CrecordController) {

	Application.UserviewController = CrecordController.extend(Ember.Evented, {
		actions: {
			insertWidget: function(containerController) {
				console.log("insertWidget", containerController);
				var widgetChooserForm = Canopsis.utils.forms.show('widgetform', this);

				var userviewController = this;

				widgetChooserForm.submit.done(function(form, newWidgetWrappers) {
					void (form);
					var newWidgetWrapper = newWidgetWrappers[0];

					console.log('onsubmit, adding widgetwrapper to containerwidget', newWidgetWrapper, containerController);
					console.log('containerwidget items', containerController.get('content.items.content'));
					//FIXME wrapper does not seems to have a widget
					containerController.get('content.items.content').pushObject(newWidgetWrapper);

					console.log("saving view");
					userviewController.get('content').save();
				});

				widgetChooserForm.submit.fail(function() {
					console.info('Widget insertion canceled');
				});
			}
		}
	});

	return Application.UserviewController;
});