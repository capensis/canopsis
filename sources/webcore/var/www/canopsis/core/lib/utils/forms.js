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

define(['app/application'], function(Application) {

	var formUtils = {
		instantiateForm: function(formName, formContext, options) {
			void (formContext);
			console.log('try to instantiate form', formName, options.formParent);
			var classDict = options;

			classDict.target = Canopsis.utils.routes.getCurrentRouteController();
			classDict.container = Application.__container__;

			var formController = Canopsis.forms.all[formName].EmberClass.create(classDict);

			return formController;
		},

		show: function(formName, formContext, options) {
			if (options === undefined) {
				options = {};
			}

			console.log("Form generation", formName);

			var formController = this.instantiateForm(formName, formContext, options);
			console.log("formController", formController);

			Canopsis.utils.routes.getCurrentRouteController().send('showEditFormWithController', formController, formContext, options);
			return formController;
		},

		editRecord: function(record) {
			var widgetWizard = Canopsis.utils.forms.show('modelform', record);
			console.log("widgetWizard", widgetWizard);

			widgetWizard.submit.then(function() {
				console.log('record saved');

				record.save();
				widgetWizard.trigger('hidePopup');
				widgetWizard.destroy();
			});

			return widgetWizard;
		},

		addRecord: function(record_type) {
			Canopsis.utils.routes.getCurrentRouteController().send('show_add_crecord_form', record_type);
		}
	};

	return formUtils;
});