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
    name:'GridlayoutMixin',
    after: 'MixinFactory',
    initialize: function(container, application) {
        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;

        var viewMixin = Ember.Mixin.create({
            didInsertElement: function() {
                //iteration over all content widget to set them the appropriage css class
                var wrappers = get(this, 'controller.items.content');
                //if view update, push it to db
                var haveToSaveView = false;

                var containerGridLayout,
                    forcelegacy = false;

                var containerMixins = get(this, 'controller.mixins');
                if (!isNone(containerMixins)) {
                    containerGridLayout = containerMixins.findBy('name', 'gridlayout');
                    if (!isNone(containerGridLayout)) {
                        forcelegacy = get(containerGridLayout, 'forcelegacy');
                    }
                }

                var gridLayoutMixin = get(this, 'controller.mixins').findBy('name', 'gridlayout');
                var columnXS = '12';

                if(!isNone(get(gridLayoutMixin, 'columnXS'))) {
                    columnXS = get(gridLayoutMixin, 'columnXS');
                }

                var columnMD = '6';
                if(!isNone(get(gridLayoutMixin, 'columnMD'))) {
                    columnMD = get(gridLayoutMixin, 'columnMD');
                }

                var columnLG = '3';
                if(!isNone(get(gridLayoutMixin, 'columnLG'))) {
                    columnLG = get(gridLayoutMixin, 'columnLG');
                }


                var offset = get(gridLayoutMixin, 'offset') || '0';

                var classValue = [
                    'col-md-',
                    columnMD,
                    ' col-xs-',
                    columnXS,
                    ' col-lg-',
                    columnLG,
                    ' col-md-offset-',
                    offset
                ].join('');

                if(get(gridLayoutMixin, 'hiddenLG')) {
                    classValue += ' hidden-lg';
                }
                if(get(gridLayoutMixin, 'hiddenMD')) {
                    classValue += ' hidden-md';
                }
                if(get(gridLayoutMixin, 'hiddenXS')) {
                    classValue += ' hidden-xs hidden-sm';
                }

                set(this, 'controller.defaultItemCssClass', classValue);

                if(wrappers) {
                    for (var i = wrappers.length - 1; i >= 0; i--) {
                        //Dynamic mixin values setting
                        var currentWrapperMixins = get(wrappers[i], 'mixins');
                        if (isNone(currentWrapperMixins)) {
                            set(wrappers[i], 'mixins', []);
                            currentWrapperMixins = get(wrappers[i], 'mixins');
                        }

                        if (isNone(currentWrapperMixins.findBy('name', 'gridlayout'))) {
                            //Legacy management, defining sub wrapper container mixin data.
                            if (isNone(containerGridLayout)) {
                                containerGridLayout = {'name': 'gridlayout'};
                            }
                            //Push fresh information
                            currentWrapperMixins.pushObject(containerGridLayout);
                            haveToSaveView = true;
                        }

                        //Manage legacy data from mixin to content wrappers
                        if (forcelegacy) {
                            //Get it now because it may have just been added before
                            var existingGridLayoutMixin = currentWrapperMixins.findBy('name', 'gridlayout');
                            //Clean previous information
                            if (!isNone(existingGridLayoutMixin)) {
                                currentWrapperMixins.removeObject(existingGridLayoutMixin);
                            }
                            //avoid side effects
                            delete containerGridLayout.forcelegacy;

                            //Add legacy information
                            currentWrapperMixins.pushObject(containerGridLayout);
                            haveToSaveView = true;
                        }

                        //Computes class value
                        classValue = get(this, 'controller').getSection(currentWrapperMixins);

                        Ember.setProperties(wrappers[i], {
                            'classValue': classValue
                        });
                    }
                }

                if (haveToSaveView) {
                    get(this, 'controller.viewController.content').save();
                }

                this._super();
            }
        });


        /**
         * Provides a responsive grid layout to a container widget.
         *
         * The grid is managed by bootstrap CSS classes.
         *
         * @mixin Gridlayout
         */
        var mixin = Mixin('gridlayout', {
            partials: {
                layout: ['gridlayout']
            },

            init: function() {
                //Attach view to the mixin
                this._super();
                this.addMixinView(viewMixin);
            },

            isGridLayout: function () {
                //Tells the controller of this mixin that it is a grid layout
                return true;
            }.property(),

            /**
             *   Builds css classes for the widget wrapper that allow responsive parametrized diplay
             *  depending on legacy/overriden values.
             **/
            getSection: function (currentWrapperMixins) {
                var gridLayoutMixin = currentWrapperMixins.findBy('name', 'gridlayout');
                var columnXS = '12';
                if(!isNone(gridLayoutMixin.columnXS)) {
                    columnXS = gridLayoutMixin.columnXS;
                }
                var columnMD = '6';
                if(!isNone(gridLayoutMixin.columnMD)) {
                    columnMD = gridLayoutMixin.columnMD;
                }
                var columnLG = '3';
                if(!isNone(gridLayoutMixin.columnLG)) {
                    columnLG = gridLayoutMixin.columnLG;
                }

                var offset = gridLayoutMixin.offset || '0';

                var classValue = [
                    'col-md-',
                    columnMD,
                    ' col-xs-',
                    columnXS,
                    ' col-lg-',
                    columnLG,
                    ' col-md-offset-',
                    offset
                ].join('');

                if(gridLayoutMixin.hiddenLG) {
                    classValue += ' hidden-lg';
                }
                if(gridLayoutMixin.hiddenMD) {
                    classValue += ' hidden-md';
                }
                if(gridLayoutMixin.hiddenXS) {
                    classValue += ' hidden-xs hidden-sm';
                }

                set(this, 'defaultItemCssClass', classValue);

                return classValue;
            }
        });

        application.register('mixin:gridlayout', mixin);
    }
});
