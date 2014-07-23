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
	'app/view/editor'
], function(Ember, Application, EditorView) {

	/**
	 * Editors factory. Creates a view, store it in Application and an appropriate Handlebars helper
	 * @param editorName {string} the name of the new editor. lowercase
	 * @param classdict {dict} the view dict
	 * @param options {dict} options :
	 *			- subclass: to handle editor's view inheritance: default is EditorView
	 *			- templateName: to use another template in the editor
	 *
	 * @author Gwenael Pluchon <info@gwenp.fr>
	 */
	function Editor(editorName, classdict, options) {
		if (options === undefined) {
			options = {};
		}

		if (options.subclass === undefined) {
			options.subclass = EditorView;
		}

		if (classdict.templateName === undefined) {
			//FIXME this is not working
			//TODO write a dict in Application, with all templates stored, and an api to search for {editor, renderers, regular} templates
			classdict.templateName = "editor-" + editorName;
		}

		var editorViewName = editorName.camelize().capitalize();

		Application[editorViewName] = options.subclass.extend(classdict);
		Ember.Handlebars.helper('editor-' + editorName, Application[editorViewName]);

		return Application[editorViewName];
	}

	console.log("factory editor loaded");

	return Editor;
});