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
                    interval = setInterval(function () {
                        console.log('refreshing widget ' + get(widgetController, 'title'), widgetController.get('mixinOptions.periodicrefresh.refreshInterval'), widgetController);
                        //FIXME periodicrefresh deactivated in testing mode because it throws global failures
                        if(window.environment !== 'test') {
                            widgetController.refreshContent();
                        }
                    }, widgetController.get('mixinOptions.periodicrefresh.refreshInterval') * 1000);

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
