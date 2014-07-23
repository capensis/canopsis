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

	Application.ComponentEditorComponent = Ember.Component.extend({
		tagName: 'span',
		init: function() {
			console.log("init editor compo");

			this._super();

			this.set('templateData.keywords.attr', Ember.computed.alias('content'));
		},

		editorType: function() {
			var type = this.get('content.model.type');
			var role = this.get('content.model.options.role');
			console.log('editorType', this.get('content'));
			console.log('editorType', type, role);
			var editorName;

			if (role) {
				editorName = 'editor-' + role;
			} else {
				editorName = 'editor-' + type;
			}

			if (Ember.TEMPLATES[editorName] === undefined) {
				editorName = 'editor-defaultpropertyeditor';
			}

			return editorName;
		}.property('content.type', 'content.role'),

		attr: Ember.computed.alias("content")
	});

	return Application.ComponentEditorComponent;
});