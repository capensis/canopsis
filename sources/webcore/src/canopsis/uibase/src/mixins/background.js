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
    name:'BackgroundMixin',
    after: ['MixinFactory', 'FormsUtils', 'HashUtils'],
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get;


        var viewMixin = Ember.Mixin.create({
            didInsertElement: function() {
                var imageUrl = get(this, 'controller.mixinOptions.background.imageUrl'),
                    position = get(this, 'controller.mixinOptions.background.position'),
                    repeat = get(this, 'controller.mixinOptions.background.repeat'),
                    size = get(this, 'controller.mixinOptions.background.size'),
                    backgroundcolor = get(this, 'controller.mixinOptions.background.backgroundcolor'),
                    element = this.$();

                if(imageUrl) {
                    element.css('background-image', 'url(' + imageUrl + ')');
                }
                if(position) {
                    element.css('background-position', position);
                }
                if(size) {
                    element.css('background-size', size);
                }
                if(repeat) {
                    element.css('background-repeat', repeat);
                }
                if(backgroundcolor) {
                    element.css('background-color', backgroundcolor);
                }

                this._super();
            }
        });

        /**
         * @mixin background
         * @description Allows to change the background color or image of widgets.
         * ![Preview](../screenshots/mixin-background.png)
         */
        var mixin = Mixin('background', {
            /**
             * @method init
             */
            init: function() {
                this._super();
                this.addMixinView(viewMixin);
            }
        });

        application.register('mixin:background', mixin);
    }
});
