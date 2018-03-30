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
    name: 'component-slider',
    initialize: function(container, application) {

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        /**
         * @component slider
         * @description Displays an editable slider based on the ionRangeSlider library that allows to edit a numeric value
         *
         * ![Slider preview](../screenshots/component-slider.png)
         * @example {{component-slider content=view.charge options=view.optionsCharge}}
         */
        var component = Ember.Component.extend({
            /**
             * @property content
             * @type integer
             * @description the numeric value to edit
             */
            content: undefined,

            /**
             * @property options
             * @type object
             * @description an option dictionnary for the slider.
             * Options are : min (integer), max (integer), step (integer)
             */
            options: undefined,

            /**
             * @method didInsertElement
             * @description enable the ionRangeSlider
             */
            didInsertElement: function () {
                var sliderComponent = this;

                var options = get(this, 'options');
                var min = get(options, 'min') || 0;
                var max = get(options, 'max') || 100;
                var step = get(options, 'step') || 1;

                var value = parseInt(get(sliderComponent, 'content'));
                if (isNone(value) || isNaN(value)) {
                    value = get(options, 'default') || min;
                }

                console.log('slider options', {
                    min: min,
                    max: max,
                    step: step,
                    value: value
                });

                var slider = sliderComponent.$('#range_slider');

                slider.ionRangeSlider({
                    min: min,
                    max: max,
                    from: value,
                    type: 'single',
                    step: step,
                    prefix: '',
                    onChange: function (data) {
                        set(sliderComponent, 'content', get(data, 'from'));
                    }
                });

                //hack as library does not manage properly the from parameter in this version.
                var mockFrom = function () {
                    var irsLine = sliderComponent.$('.irs-line');
                    if (irsLine !== undefined) {
                        if (irsLine.is(':visible')) {
                            var width = irsLine.width();
                            var proportion = width / max * value;
                            //nice display ajustement....
                            var maxwidth = width - 20;
                            if (proportion > maxwidth) {
                                proportion = maxwidth;
                            }
                            //Manually placing slider and tooltip proportionnaly to the width of the slider.
                            sliderComponent.$('.irs-single').css('left', proportion);
                            sliderComponent.$('.irs-slider').css('left', proportion);
                        } else {
                            setTimeout(mockFrom, 500);
                        }
                    }
                };

                //as from feature doesn t work on this slider...
                mockFrom();
            }
        });
        application.register('component:component-slider', component);
    }
});
