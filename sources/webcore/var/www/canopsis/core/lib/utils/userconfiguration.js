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
	'utils'
], function(Ember, Application, utils) {

	Application.userConfiguration = Ember.ObjectController.extend({

		needs: ['login'],
		content: {},
		init: function () {
			this.loadUserConfiguration();
		},

		saveUserConfiguration: function (preferences_level) {

			var preferences = this.get('widget.userParams');
			console.debug('Ready to save user configuration', preferences);

			var preference_id = this.get('preference_id');
			if (preference_id === undefined) {
				preference_id = utils.hash.generate_GUID();
			}

			var userConfiguration = {
				preferences_level: preferences_level,
				widget_preferences: preferences,
				crecord_name: utils.session.username,
				widget_id: this.get('widget.id'),
				id: preference_id,
				crecord_type: 'userpreferences'
			};
			//TODO @eric use an adapter
			$.ajax({
				url: '/rest/userpreferences/userpreferences',
				type: 'POST',
				data: JSON.stringify(userConfiguration),
				success: function(data) {
					void (data);
					console.log('User configuration save statement for widget complete');
				}
			});

		},

		loadUserConfiguration: function() {
			var userConfiguration = this;

			console.debug('loading configuration');
			//TODO @eric use an adapter
			$.ajax({
				url: '/rest/userpreferences/userpreferences',
				async: false,
				data: {
					limit: 1,
					filter: JSON.stringify({
						crecord_name: utils.session.username,
						widget_id: this.get('widget.id')
					})
				},
				success: function(data) {
					if (data.success && data.data.length && data.data[0].widget_preferences !== undefined) {
						console.log('User configuration for widget complete', data);
						var preferences = data.data[0].widget_preferences;
						userConfiguration.set('preference_id', data.data[0]._id);
						userConfiguration.set('widget.userParams', preferences);
						for (var key in preferences) {
							console.debug('User preferences: will set key', key, 'in widget', userConfiguration.get('widget.title'));
							userConfiguration.set('widget.' + key, preferences[key]);
						}
					} else {
						console.debug('No user preference exists for widget' + userConfiguration.get('widget.title'));
					}
				}
			}).fail(
				function (error) {
					void (error);
					console.log('No user s preference found for this widget');
				}
			);
		}

	});


	return Application.userConfiguration;
});