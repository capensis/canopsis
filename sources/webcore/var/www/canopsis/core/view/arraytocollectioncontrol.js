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
	'app/mixins/arraymixin',
	'app/view/crecords'
], function(Ember, Application) {


    Application.ArrayToCollectionControlView = Ember.CollectionView.extend(Application.ArrayMixin,{

	itemViewClass: Ember.View.extend({
	    tagName: '',
	    template: Ember.Handlebars.compile(" <button  data-hint={{ unbound template.label}} {{action 'modify' template target='view.parentView' }} {{bind-attr class='template.CSSclass'}}  > {{glyphicon template.icon}}   </button> ")
	}),

	/*
	 *  modify template's CSSClass and value (called when button is pressed).
	 */
	actions: {
	    modify: function(template) {
			var value = this.get("value");
			var isPresent = this.checkIfAContainB(value,template);
			console.log("isPresent = ",isPresent, " value = ",value," and template =", template);

			if (!isPresent) {
			    value.pushObject(template.name);
			} else {
			    value.removeObject(template.name);
			}
			this.changeCssClass(template,value);
	    }
	}
    });
    return Application.ArrayToControlView;

});