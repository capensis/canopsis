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
	 * Helper to display a renderer. Uses the context to get attribute value and options, so take care to where you call this helper.
	 * @param rendererType {string} the renderer to try to display. If the parameter does not match an existing renderer, falls back to the default one
	 *
	 * @author Gwenael Pluchon <info@gwenp.fr>
	 */
	Ember.Handlebars.registerHelper('renderer', function(options) {
		console.log("renderer helper", arguments);
		console.log(options);

		var attr = options.data.keywords.attr;
		var value;

		if (!Ember.isNone(attr.value)) {
			value = attr.value;
		} else if (!Ember.isNone(options.data.keywords) && !Ember.isNone(options.data.keywords.crecord)) {
			value = options.data.keywords.crecord.get(attr.field);
		} else {
			value = '';
		}

		var rendererName;

		console.log(attr);
		console.log(value);

		var rendererContext = {
			attr : attr,
			value : value
		};

		options.contexts.push(rendererContext);

		if (!Ember.isNone(attr.options) && !Ember.isNone(attr.options.role)) {
			rendererName = attr.options.role;
		} else if (!Ember.isNone(attr.model) && !Ember.isNone(attr.model.options) && !Ember.isNone(attr.model.options.role)) {
			rendererName = attr.model.options.role;
		} else {
			rendererName = attr.type;
		}

		rendererName = "renderer-" + rendererName;

		//trying to find if the required renderer is an helper or a template
		var foundRenderer = "renderer-default";

		//if a template matches the renderer, select it, else keep the standard one
		if (Ember.TEMPLATES[rendererName] !== undefined) {
			foundRenderer = rendererName;
		}

		if(foundRenderer === 'renderer-default') {
			if (Ember.isNone(value)) {
				value = '';
			} else {
				value = Ember.Handlebars.Utils.escapeExpression(value);
			}

			return new Ember.Handlebars.SafeString(value);
		} else {
			console.log("call renderer partial", foundRenderer);
			return Ember.Handlebars.helpers.partial.apply(this, [foundRenderer, options]);
		}

	});

});
