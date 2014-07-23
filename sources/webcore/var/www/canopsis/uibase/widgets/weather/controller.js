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
	'jquery',
	'app/lib/factories/widget'
], function($, WidgetFactory) {

	var widget = WidgetFactory('weather', {
		init: function() {
			this._super();
			this.set('worst_state', 0);
			this.set('sub_weather', []);
			this.fetchStates();
			console.log('Setting up weather widget : ' + this.get('config.title'));
		},

		//generate and refresh the title
		title: function () {
			return this.get('config.title');
		}.property('config.title'),

		icon: function () {
			return this.class_icon(this.get('worst_state'));
		}.property('worst_state'),

		//generate weather class depending on status
		class_icon: function (status) {
			return 'ion-ios7-' + [
				'sunny',
				'cloudy',
				'thunderstorm',
				'thunderstorm'][status] + '-outline';
		},

		//generate and refresh background property for widget weather display
		background: function () {
			return this.class_background(this.get('worst_state'));
		}.property('worst_state'),

		//generate weather class depending on status
		class_background: function (status) {
			return [
				'bg-green',
				'bg-orange',
				'bg-red',
				'bg-red'][status];
		},

		fetchStates: function () {
			var that = this;
			var rks = that.get('config.event_selection');

			if (!rks || !rks.length) {
				console.warn('Widget weather ' + this.get('title') + ' No rk found, the widget may not be configured properly');
				return;
			}

			var params = {
				limit: 0,
				ids: JSON.stringify(rks)
			};

			$.ajax({
				url: '/rest/events',
				data: params,
				success: function(data) {
					if (data.success) {
						that.computeWeather(data.data);
					} else {
						console.error('Unable to load event information for weather widget from API');
					}
					that.trigger('refresh');
					console.log(' + Weather content', that.get('config.event_selection'));
				}
			});
		},

		computeWeather: function (data) {
			console.log(' + computing weathers');
			var worst_state = 0;
			var sub_weathers = [];
			for (var i=0; i<data.length; i++) {

				//computing worst state for general weather display
				if (data[i].state > worst_state) {
					worst_state = data[i].state;
				}

				//compute sub item title depending on if resource exists
				var resource = '';
				if (data[i].resource) {
					resource = ' ' + data[i].resource;
				}
				//building the data structure for sub parts of the weather
				sub_weathers.push({
					title: data[i].component + resource,
					custom_class: this.class_background(data[i].state)
				});

			}
			console.log('weather content', {sub_weathers: sub_weathers, worst_state: worst_state});
			this.set('sub_weather', sub_weathers);
			this.set('worst_state', worst_state);
			this.trigger('refresh');
		},

		refreshContent: function () {
			this.fetchStates();
			this._super();
		}
	});

	return widget;
});