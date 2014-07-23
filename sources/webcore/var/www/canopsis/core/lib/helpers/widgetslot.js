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

define(['ember', 'app/view/widgetslot', 'app/view/widget'], function(Ember, WidgetslotView, WidgetView) {

	/**
	 * Helper to display an editor. Uses the context to get attribute value and options, so take care to where you call this helper.
	 * @param editorType {string} the editor to try to display. If the parameter does not match an existing editor, falls back to the default one
	 *
	 * @author Gwenael Pluchon <info@gwenp.fr>
	 */
	Ember.Handlebars.helper('widgetslot', WidgetslotView);
	Ember.Handlebars.helper('widgethelper', WidgetView);
});


