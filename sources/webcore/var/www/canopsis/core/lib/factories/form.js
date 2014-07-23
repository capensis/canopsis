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
	'app/application',
	'app/controller/form',
	"app/lib/formsmanager",
	'app/view/form'
], function(Application, FormController, FormsManager) {

	/**
	 * Form factory. Creates a controller, stores it in Application
	 * @param formName {string} the name of the new form. lowercase
	 * @param classdict {dict} the controller dict
	 * @param options {dict} options :
	 *			- subclass: to handle form's controller inheritance: default is FormController
	 *			- templateName: to use another template in the editor
	 *
	 * @author Gwenael Pluchon <info@gwenp.fr>
	 */
	function Form(formName, classdict, options) {
		console.log("new form", arguments);

		var extendArguments = [];

		if (options === undefined) {
			options = {};
		}

		if (options.subclass === undefined) {
			options.subclass = FormController;
		}

		if (options.mixins !== undefined) {
			for (var i = 0; i < options.mixins.length; i++) {
				extendArguments.push(options.mixins[i]);
			}
		}

		extendArguments.push(classdict);

		var formControllerName = formName.camelize().capitalize() + "Controller";
		var formViewName = formName.camelize().capitalize() + "View";
		console.log("extendArguments", extendArguments);
		console.log("subclass", options.subclass);


		Application[formViewName] = Application.FormView.extend();
		Application[formControllerName] = options.subclass.extend.apply(options.subclass, extendArguments);

		FormsManager.all[formName] = {
			EmberClass: Application[formControllerName]
		};

		return Application[formControllerName];
	}

	console.log("factory form loaded");

	return Form;
});