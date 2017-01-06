/*
 * Copyright (c) 2015 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
 */

Ember.Application.initializer({
    name: 'AlarmsWidget',
    after: ['TimeWindowUtils', 'DataUtils', 'WidgetFactory', 'AlarmsView', 'UserconfigurationMixin', 'SchemasLoader'],
    initialize: function(container, application) {
		var timeWindowUtils = container.lookupFactory('utility:timewindow'),
        	dataUtils = container.lookupFactory('utility:data'),
			WidgetFactory = container.lookupFactory('factory:widget'),
        	viewMixin = container.lookupFactory('view:alarms'),
			UserConfigurationMixin = container.lookupFactory('mixin:userconfiguration');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        // load the viewMixin
        var listOptions = {
            mixins: [
                UserConfigurationMixin
            ]
        };

        /**
         * This widget allows to display events, with or without reccurence rules, on a calendar.
         *
         * @memberOf canopsis.frontend.brick-calendar
         * @mixes UserConfigurationMixin
         * @class WidgetStatsTable
         * @widget calendar
         */
        var widget = WidgetFactory('alarms',{

			viewMixins: [
                 viewMixin
            ],

            /**
             * Create the widget and set widget params into Ember vars
             * @method init
             */
            init: function() {
                this._super.apply(this, arguments);
				//this.setParams();

				set(this, 'store', DS.Store.extend({
                    container: get(this, 'container')
                }));
            },

			/**
			 * @method updateInterval
			 */
			/*updateInterval: function(interval) {
                var from = get(interval, 'timestamp.$gte'),
                    to = get(interval, 'timestamp.$lte');

                if(!isNone(from)) {
                    set(this, 'from', from * 1000);
                }
                else {
                    set(this, 'from', undefined);
                }

                if(!isNone(to)) {
                    set(this, 'to', to * 1000);
                }
                else {
                    set(this, 'to', undefined);
                }

                this.setParams();
            },*/

            /**
             * Set the reload to true in order to redraw events
             * extend the native refreshContent method
             * @method refreshContent
             */
            refreshContent: function () {
				// Not implemented because backend too long, feature not useful for this widget
            },

            /**
             * Set widget params to create an Ember object to manipulate
             * @method setParams
             */
			setParams: function() {
            	var controller = this;

				var tw = timeWindowUtils.getFromTo(
                    get(controller, 'time_window'),
                    get(controller, 'time_window_offset')
                );

				var from = tw[0],
                    to = tw[1];

                /* live reporting support */
                var liveFrom = get(controller, 'from'),
                    liveTo = get(controller, 'to');

				if (!isNone(liveFrom)) {
                    from = liveFrom;
                }

                if (!isNone(liveTo)) {
                    to = liveTo;
                }

				var widgetParams = {
					timewindow: {
						from: parseInt(from / 1000),
                    	to: parseInt(to / 1000)
					},
					title: get(controller, 'model.title'),
					domain: get(controller, 'model.domain'),
					perimeter: get(controller, 'model.perimeter'),
					user: get(controller, 'model.user'),
					show_events: get(controller, 'model.show_events'),
					show_users: get(controller, 'model.show_users')
				};

				set(controller, 'widgetParams', widgetParams);
				if (get(controller, 'model.show_events')) {
					controller.fetchEvents();
				}
				if (get(controller, 'model.show_users')) {
					controller.fetchUser();
				}
            },

            /**
             * Get events statistics from backend
             * @method fetchEvents
             */
            fetchEvents: function() {
            	var controller = this;

				var tags= {
					domain: get(controller, 'widgetParams.domain'),
					perimeter: get(controller, 'widgetParams.perimeter')
				};

				var tag = JSON.stringify(tags);
				var query = {
					tstart: get(controller, 'widgetParams.timewindow.from'),
					tstop: get(controller, 'widgetParams.timewindow.to'),
					tags: tag
				};

				var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:statstableevents');
                adapter.findQuery('statstableevents', query).then(function (result) {
            		// onfullfillment
					var stats = get(result, 'data');
					set(controller, 'statsevent', stats[0]);
					controller.trigger('disableLoading');
                }, function (reason) {                                                                                                                     console.error('ERROR in the adapter: ', reason);
            	});
			},

            /**
             * Get user statistics from backend
             * @method fetchUser
             */
            fetchUser: function() {
            	var controller = this;

				var tags= {
                    domain: get(controller, 'widgetParams.domain'),
                    perimeter: get(controller, 'widgetParams.perimeter')
                };

				var tag = JSON.stringify(tags);
                var query = {
                    tstart: get(controller, 'widgetParams.timewindow.from'),
                    tstop: get(controller, 'widgetParams.timewindow.to'),
                    users: get(controller, 'widgetParams.user'),
					tags: tag
                };

                var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:statstableuser');
            	adapter.findQuery('statstableuser', query).then(function (result) {
                    // onfullfillment
					var stats = get(result, 'data');
					set(controller, 'statsuser', stats);
					controller.trigger('disableLoading');
                }, function (reason) {
                    // onrejection
                    console.error('ERROR in the adapter: ', reason);
                });
            }

        }, listOptions);

        application.register('widget:alarms', widget);
    }
});
