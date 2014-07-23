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

	Application.ComponentExpandableaddbuttonComponent = Ember.Component.extend({
		tagName: 'span',

		yieldSpanStyle: 'display:none',

		focusIn: function() {
			//check if focused is not true to avoid bubbling to execute in loop the code below
			if (this.get('focused') !== true) {
				this.set('focused', true);

				this.$('.yieldContainer').css({'display':'inline'});
				this.set('focusingChildInput', true);
				this.$(".defaultFocus").focus();
				this.set('focusingChildInput', false);
			}
		},

		focusOut: function() {
			console.log('focusOut', this.get('focusingChildInput'));

			if (! this.get('focusingChildInput')) {
				this.$('.yieldContainer').css({'display':'none'});
				this.set('focused', false);

				var inputValue = this.$('.defaultFocus').val();

				if (!! inputValue) {
					console.log('child input has a value');
					this.sendAction('onAddElement', inputValue);
				}

				this.$('.defaultFocus').val('');
			}
		}
	});

	return Application.ComponentExpandableaddbuttonComponent;
});