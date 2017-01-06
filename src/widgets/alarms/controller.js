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

    				    set(this, 'store', DS.Store.extend({
                    container: get(this, 'container')
                }));

                this.fetchAlarms();
            },

            /**
             * Set the reload to true in order to redraw events
             * extend the native refreshContent method
             * @method refreshContent
             */
            refreshContent: function () {},

            /**
             * Get events statistics from backend
             * @method fetchEvents
             */
            fetchAlarms: function() {
              var controller = this;
              var query = {};

              var adapter = dataUtils.getEmberApplicationSingleton().__container__.lookup('adapter:alerts');
              adapter.findQuery('alerts', query).then(function (result) {
                  // onfullfillment
                  var alerts = get(result, 'data');
                  console.error('alerts', alerts);
              }, function (reason) {                                                                                                                     console.error('ERROR in the adapter: ', reason);
                  // onrejection
              });
            }

        }, listOptions);

        application.register('widget:alarms', widget);
    }
});
