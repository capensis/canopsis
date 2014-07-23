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

define(['ember'], function(Ember) {

	/**
	 * Helper to display an editor. Uses the context to get attribute value and options, so take care to where you call this helper.
	 * @param editorType {string} the editor to try to display. If the parameter does not match an existing editor, falls back to the default one
	 *
	 * @author Gwenael Pluchon <info@gwenp.fr>
	 */
	Ember.Handlebars.registerHelper('editorhelper', function(editorType, options) {
		void (editorType);
		console.log("editor helper", arguments, options.data.keywords.attr.editor);

		var editor = options.data.keywords.attr.editor;

		//trying to find if the required editor is an helper or a template
		if (Ember.Handlebars.helpers[editor] !== undefined) {

			//rendering editor by calling the helper
			console.log("call editor helper for type", editor);
			return Ember.Handlebars.helpers[editor].apply(this, [options]);

		} else {
			var foundEditor = "editor-defaultpropertyeditor";

			//if a template matches the editor, select it, else keep the standard one
			if (Ember.TEMPLATES[editor] !== undefined) {
				foundEditor = editor;
			}

			console.log("call editor partial for type", foundEditor);
			return Ember.Handlebars.helpers.partial.apply(this, [foundEditor, options]);
		}

		console.error("editor helper did not find a way to display correctly an editor", arguments);
	});

});
