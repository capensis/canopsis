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
 * 
 * 
 * IMPORTANT
 *  The periodic refresh loads once per widget in the view hierarchy. So if you have 2 widgets, each with
 *  an insance of this mixin, they might run twice on each widget. 
 *  As a general rule, try to have only one periodic refresh instance per view to get the expected refresh rate.
 */


Ember.Application.initializer({
    name: 'PeriodicrefreshMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;


        var viewMixin = Ember.Mixin.create({

            willInsertElement: function() {
                console.log('init periodicrefresh viewMixin');
                this._super();

                //widget refresh management
                var widgetController = get(this, 'controller');

                var previousInterval = get(this, 'mixinOptions.periodicrefresh.refreshInterval');
                if(previousInterval) {
                    clearInterval(previousInterval);
                }

                var interval = get(this, 'widgetRefreshInterval');
                var mixin = this;

                Ember.run(function(){
                    // We get the refreshInteval from the local config.      
                    refreshIntvalConfigValue = widgetController.get('mixinOptions.periodicrefresh.refreshInterval') * 1000;
                    
                    // If no value is fetched, that means that this particular instance should not run
                    // so we return early to avoid setInterval to be called in a 0-delay loop, causing 
                    // a full UI crash.
                    if (isNaN(refreshIntvalConfigValue)) {
                        return;
                    }

                    interval = setInterval(function () {
                        console.log('refreshing widget ' + get(widgetController, 'title'), widgetController.get('mixinOptions.periodicrefresh.refreshInterval'), widgetController);
                        //FIXME periodicrefresh deactivated in testing mode because it throws global failures
                        if(window.environment !== 'test') {
                            widgetController.refreshContent();
                        }
                    },refreshIntvalConfigValue);

                    //keep track of this interval
                    set(mixin, 'widgetRefreshInterval', interval);
                });
            },


            willDestroyElement: function() {
                clearInterval(get(this, 'widgetRefreshInterval'));

                this._super();
            }

        });

        /**
         * @mixin periodicrefresh
         */
        var mixin = Mixin('periodicrefresh', {

            init:function() {
                console.log('init periodicrefresh');
                this.addMixinView(viewMixin);

                var mixinsOptions = get(this, 'content.mixins');

                if(mixinsOptions) {
                    var periodicrefreshOptions = get(this, 'content.mixins').findBy('name', 'periodicrefresh');
                    this.mixinOptions.periodicrefresh = periodicrefreshOptions;
                }

                this._super.apply(this, arguments);
                this.startRefresh();

                //setting default/minimal reload delay for current widget
                if (get(this, 'mixinOptions.periodicrefresh.refreshInterval') < 10 || isNone(get(this, 'mixinOptions.periodicrefresh.refreshInterval'))) {
                    set(this, 'mixinOptions.periodicrefresh.refreshInterval', 10);
                }
            }
        });

        application.register('mixin:periodicrefresh', mixin);
    }
});
