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
// TODO: just make a function from this
	Application.mixinArrayMixin = Ember.Mixin.create({

		onInit : function ( contentREF , _self ){

			function getAndPushMixinNames(classToGet , contentREF){
				var currentClass = SearchableMixin.byClass[classToGet];
				for ( var i = 0 ; i < currentClass.length ; i++ ) {
					var nameMixin = { name : currentClass[i] };
					contentREF.push(nameMixin);
				}
			}

			var formController  =  Canopsis.formwrapperController.form;
            if ( formController ){
				var classToGet = _self.templateData.keywords.controller.content.model.options.mixinClass;
				var SearchableMixin = Canopsis.Application.SearchableMixin;

				if (classToGet !== undefined) {
					getAndPushMixinNames( classToGet , contentREF );
				}
				else {
					for ( var attribut in SearchableMixin.byClass ) {
						if ( SearchableMixin.byClass.hasOwnProperty( attribut ) ) {
							getAndPushMixinNames( attribut , contentREF );
						}
					}
				}
			}
			_self.set("select", 1 );
		}
	});

	return Application.mixinArrayMixin;
});