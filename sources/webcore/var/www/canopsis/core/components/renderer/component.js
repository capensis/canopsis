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

	Application.ComponentRendererComponent = Ember.Component.extend({
		tagName: 'span',

		value: function() {
			return this.get('record.' + this.get('attr.field'));
		}.property('attr.field', 'record'),

		rendererType: function() {
			var type = this.get('attr.type');
			var role = this.get('attr.options.role');
			var rendererName;
			if (role) {
				rendererName = 'renderer-' + role;
			} else {
				rendererName = 'renderer-' + type;
			}

			if (Ember.TEMPLATES[rendererName] === undefined) {
				rendererName = undefined;
			}

			return rendererName;
		}.property('attr.type', 'attr.role')

	});

	return Application.ComponentEditorComponent;
});