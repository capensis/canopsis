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

	Application.WidgetView = Ember.View.extend({
		templateName:'widget',
		classNames: ['widget'],

		/**
		 * Used to visually display error messages to the user (in the widget template)
		 */
		errorMessages : Ember.A(),

		init: function() {
			this.set('target', this.get('controller'));

			this._super();
			if (!! this.widget) {
				this.setupController(this.widget);
			} else {
				console.error("no correct widget found for view", this);
				this.errorMessages.pushObject('No correct widget found');
			}
		},

		setupController: function(widget) {
			console.group('set controller for widget', widget);
			this.set("controller", this.instantiateCorrectController(widget));
			this.set("templateName", widget.get("xtype"));
			this.registerHooks();
			console.groupEnd();
		},

		//Controller -> View Hooks
		registerHooks: function() {
			console.log("registerHooks", this.get("controller"), this.get("controller").on);
			this.get("controller").on('refresh', this, this.rerender);
		},

		unregisterHooks: function() {
		},

		rerender: function() {
			console.info('refreshing widget');
			this._super.apply(this, arguments);
		},

		instantiateCorrectController: function(widget) {
			//for a widget that have xtype=widget, controllerName=WidgetController
			var controllerName = widget.get("xtype").capitalize() + "Controller";
			var widgetController;
			console.log("controllerName", controllerName, Application[controllerName], this.get('target'));

			if (Application[controllerName] !== undefined) {
				//var mixinClass = Application.SearchableMixin
				 widgetController =  Application[controllerName].createWithMixins(Ember.Evented, {
					content: widget,
					target: this.get('target')
				});
			} else {
				 widgetController =  Application.WidgetController.createWithMixins(Ember.Evented, {
					content: widget,
					target: this.get('target')
				});
			}
			var mixinsName = widget["_data"].mixins;

			if (  mixinsName  ){
				for (var i = 0 ; i < mixinsName.length ; i++ ){
					var currentName =  mixinsName[i];
					var currentMixin = Application.SearchableMixin.all[currentName];
					if ( currentMixin ){
						currentMixin.apply(widgetController);
					}
				}
			}
			return widgetController;
		},

		didInsertElement : function() {
			console.log("inserted widget, view:", this);

			this.registerHooks();
			var result = this._super.apply(this, arguments);
			this.get('controller').onReload();
			//TODO put this somewhere on list widget
			// this.$('input').iCheck({checkboxClass: 'icheckbox_minimal-grey', radioClass: 'iradio_minimal-grey'});

			return result;
		},

		willClearRender: function() {
			this.unregisterHooks();
			return this._super.apply(this, arguments);
		}

	});

	return Application.WidgetView;
});