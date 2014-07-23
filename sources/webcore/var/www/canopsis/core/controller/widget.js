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
	'ember',
	'app/application',
	'app/lib/utils/userconfiguration',
	'utils'
], function($, Ember, Application, userConfiguration, utils) {

	Application.WidgetController = Ember.ObjectController.extend({

		userParams: {},

		init: function () {

			var widgetController = this;

			console.info("widget init");

			//manage user configuration
			this.set('userConfiguration', userConfiguration.create({widget: this}));

			this.set("container", utils.routes.getCurrentRouteController().container);
			//TODO put spinner in widget when loading

			this.set('isRefreshable', false);
			//TODO load delay from widget configuration
			this.set('refreshDelay', 10000);

			setInterval(function () {
				if (Canopsis.conf.REFRESH_ALL_WIDGETS) {
					var title = widgetController.get('title');
					if (title === undefined) {
						title = '';
					}
					if (!widgetController.get('isRefreshable')) {
						console.log('refreshing widget ' + title);
						widgetController.refreshContent();
					}
				}
			}, widgetController.get('refreshDelay'));


			this.refreshContent();

		},

		onReload: function () {
			console.debug('Reload widget:', this.get('id'));

			if (this.get('widgetData.content') !== undefined) {
				//Allows widget to know how many times they have been repainted
				if (this.get('domReadyCount') === undefined) {
					this.set('domReadyCount', 1);
				} else {
					this.set('domReadyCount', this.get('domReadyCount') + 1);
				}
				this.onDomReady($('#' + this.get('id')));
			}
		},

		onDomReady: function() {
			console.log(this.get('title'), 'widget dom load complete');
			//To override
		},

		stopRefresh: function () {
			this.set('isRefreshable', false);
		},

		startRefresh: function () {
			this.set('isRefreshable', true);
		},

		actions: {
			do: function(action) {
				var params = [];
				for (var i = 1; i < arguments.length; i++) {
					params.push(arguments[i]);
				}

				this.send(action, params);
			},
			creationForm: function(itemType) {
				utils.forms.addRecord(itemType);
			},

			editWidget: function (widget) {
				console.info("edit widget", widget);

				var widgetWizard = utils.forms.show('modelform', widget, { title: "Edit widget" });
				console.log("widgetWizard", widgetWizard);

				var widgetController = this;

				widgetWizard.submit.done(function() {
					console.log('record going to be saved', widget);

					console.log("getCurrentRouteController", utils.routes.getCurrentRouteController());

					//TODO @gwen detect if record is embedded or not HERE
					var userview = utils.routes.getCurrentRouteController().get('content');
					userview.save();
					console.log("triggering refresh");
					widgetController.trigger('refresh');
				});
			},

			removeWidget: function (widget) {
				console.group("remove widget", widget);
				console.log("parent container", this);

				var itemsContent = this.get('content.items.content');

				for (var i = 0, itemsContent_length = itemsContent.length; i < itemsContent_length; i++) {
					console.log(this.get('content.items.content')[i]);
					if (itemsContent[i].get('widget') === widget) {
						itemsContent.removeAt(i);
						console.log("deleteRecord ok");
						break;
					}
				}

				var userview = utils.routes.getCurrentRouteController().get('content');
				userview.save();

				console.groupEnd();
			},

			movedown: function(widgetwrapper) {
				console.group('movedown', widgetwrapper);
				try{
					console.log('context', this);

					var foundElementIndex,
						nextElementIndex;

					for (var i = 0; i < this.get('content.items.content').length; i++) {
						console.log('loop', i, this.get('content.items.content')[i], widgetwrapper);
						console.log(this.get('content.items.content')[i] === widgetwrapper);
						if (foundElementIndex !== undefined && nextElementIndex === undefined) {
							nextElementIndex = i;
							console.log('next element found');
						}

						if (this.get('content.items.content')[i] === widgetwrapper) {
							foundElementIndex = i;
							console.log('searched element found');
						}
					}

					if (foundElementIndex !== undefined && nextElementIndex !== undefined) {
						//swap objects
						var array = Ember.get(this, 'content.items.content');
						console.log('swap objects', array);

						var tempObject = array.objectAt(foundElementIndex);

						array.insertAt(foundElementIndex, array.objectAt(nextElementIndex));
						array.insertAt(nextElementIndex, tempObject);
						array.replace(foundElementIndex + 2, 2);

						console.log('new array', array);
						Ember.set(this, 'content.items.content', array);

						var userview = utils.routes.getCurrentRouteController().get('content');
						userview.save();
					}
				} catch (e) {
					console.error(e.stack, e.message);
				}
				console.groupEnd();
			},

			moveup: function(widgetwrapper) {
				console.group('moveup', widgetwrapper);

				try{
					console.log('context', this);

					var foundElementIndex,
						nextElementIndex;

					for (var i = this.get('content.items.content').length; i >= 0 ; i--) {
						console.log('loop', i, this.get('content.items.content')[i], widgetwrapper);
						console.log(this.get('content.items.content')[i] === widgetwrapper);
						if (foundElementIndex !== undefined && nextElementIndex === undefined) {
							nextElementIndex = i;
							console.log('next element found');
						}

						if (this.get('content.items.content')[i] === widgetwrapper) {
							foundElementIndex = i;
							console.log('searched element found');
						}
					}

					console.log('indexes to swap', foundElementIndex, nextElementIndex);

					if (foundElementIndex !== undefined && nextElementIndex !== undefined) {
						//swap objects
						var array = Ember.get(this, 'content.items.content');
						console.log('swap objects', array);

						var tempObject = array.objectAt(foundElementIndex);

						array.insertAt(foundElementIndex, array.objectAt(nextElementIndex));
						array.insertAt(nextElementIndex, tempObject);
						array.replace(nextElementIndex + 2, 2);

						console.log('new array', array);
						Ember.set(this, 'content.items.content', array);

						var userview = utils.routes.getCurrentRouteController().get('content');
						userview.save();
					}
				} catch (e) {
					console.error(e.stack, e.message);
				}
				console.groupEnd();
			},

		},

		config: Ember.computed.alias("content"),

		itemController: function() {
			return this.get("itemType").capitalize();
		}.property("itemType"),

		refreshContent: function() {
			this._super();

			this.findItems();
		},

		findItems: function() {
			console.warn("findItems not implemented");
		},

		extractItems: function(queryResult) {
			console.log("extractItems", queryResult);

			this._super(queryResult);

			this.set("widgetData", queryResult);
		}


	});

	return Application.WidgetController;
});
