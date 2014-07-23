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

	Application.WidgetslotView = Ember.View.extend({
		init: function() {
			console.log('widgetslot init', this.get('controller.content.widgetslotTemplate'));

			var widgetslotTemplate = this.get('controller.content.widgetslotTemplate');

			if(widgetslotTemplate !== undefined && widgetslotTemplate !== null && Ember.TEMPLATES[widgetslotTemplate] !== undefined) {
				this.set('templateName', widgetslotTemplate);
			}
			this._super.apply(this, arguments);
		},

		templateName:'widgetslot-default',
		classNames: ['widgetslot'],

		actions: {
			minimize: function() {
				console.log('minimize action', arguments);
				if(this.get('minimized') === true) {
					this.set('minimized', false);
				} else {
					this.set('minimized', true);
				}
			}
		}
	});

	return Application.WidgetslotView;
});