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
], function(Ember, Application) {
	Application.ComponentStateComponent = Ember.Component.extend({


		init: function() {
			this._super();
		},

		isOk:function () {
			return this.get('content') === 0;
		}.property('content'),

		isWarning:function () {
			return this.get('content') === 1;
		}.property('content'),

		isError:function () {
			return this.get('content') === 2;
		}.property('content'),

		isUnknown:function () {
			return this.get('content') === 3;
		}.property('content'),


		actions: {
			setState:function (state) {
				this.set('content', state);
			}
		},

	});

	return Application.ComponentStateComponent;
});