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
    name: 'ResponsivelistMixin',
    after: ['MixinFactory', 'FormsUtils', 'HashUtils'],
    initialize: function(container, application) {

        var Mixin = container.lookupFactory('factory:mixin');

        var get = Ember.get,
            set = Ember.set,
            isNone = Ember.isNone;


        function getColumnIndexesPriorities(viewMixin) {
            var stackableColumnsPriority = get(viewMixin, 'controller.mixinOptions.responsivelist.stackableColumnsPriority');
            var controller = get(viewMixin, 'controller');
            //FIXME undefined
            var shownColumns = get(controller, 'shown_columns');

            var columnStackingPriority = Ember.A();

            console.log('stackableColumnsPriority', stackableColumnsPriority, shownColumns);
            if(stackableColumnsPriority) {
                for (var i = 0, l = stackableColumnsPriority.length; i < l; i++) {
                    var currentColumn = shownColumns.findBy('field', stackableColumnsPriority[i]);
                    console.log('currentColumn', currentColumn);
                    if(currentColumn !== undefined) {
                        console.log('currentColumn', currentColumn);
                        var columnIndex = Ember.get(currentColumn, 'index');
                        console.log('columnIndex', columnIndex);
                        columnStackingPriority.pushObject(columnIndex);
                    }
                }
            }

            console.log('stackableColumnsPriority@end', columnStackingPriority);

            return columnStackingPriority;
        }

        function hideColumn(viewMixin, columnToHide) {
            console.log('hideColumn', columnToHide);

            if(!isNone(columnToHide)) {
                this.$('th.' + columnToHide.field).css('display', 'none');
                this.$('td.' + columnToHide.field).css('display', 'none');
                get(viewMixin, 'groupedColumns').pushObject(columnToHide);
                viewMixin.notifyPropertyChange('groupedColumns');
            }
        }


        function checkToToggleStackedDisplay(viewMixin, thresholds, tableContainerWidth, tableWidth) {
            var isTableOverflowing = tableContainerWidth < tableWidth;

            var controller = get(viewMixin, 'controller');
            var shownColumns = get(controller, 'shown_columns');
            var columnStackingPriority = getColumnIndexesPriorities(viewMixin);
            console.log('columnStackingPriority', columnStackingPriority);

            if(isTableOverflowing) {
                console.group('table overflowing, starting shrink loop');
                var i = -1;
                var newTableWidth = tableWidth;

                while(isTableOverflowing) {
                    i++;

                    console.log('while shownColumns[i]', i, shownColumns, shownColumns[i], columnStackingPriority);

                    if(shownColumns[i] !== undefined) {
                        if(columnStackingPriority[i] === undefined) {
                            console.log('no stacking priority, return', i, columnStackingPriority);
                            return;
                        }

                        hideColumn(viewMixin, shownColumns[columnStackingPriority[i]]);

                        newTableWidth -= shownColumns[i].width;
                        isTableOverflowing = tableContainerWidth < newTableWidth;
                    }
                }
                console.groupEnd();
            }
        }


        function getTableThresholds(viewMixin) {
            var controller = get(viewMixin, 'controller');

            var shownColumns = get(controller, 'shown_columns');
            for (var i = 0, l = shownColumns.length; i < l; i++) {
                set(shownColumns[i], 'width', this.$('th.' + shownColumns[i].field).width());
                set(shownColumns[i], 'index', i);
            }
        }

        var viewMixin = Ember.Mixin.create({
            classNames: ['list'],
            groupedColumns: Ember.A(),

            /**
             * Indicates the number of invisible cells to generate in the stackedcolumns view, to prevent some draggableColumnsMixin bugs
             */
            invisibleCellsCount: function() {
                var shownColumns = get(this, 'controller.shown_columns');
                return shownColumns.length - 1;
            }.property('controller.shown_columns'),

            didInsertElement: function() {
                this._super.apply(this, arguments);

                var viewMixin = this;

                get(viewMixin, 'groupedColumns').clear();

                var thresholds = getTableThresholds(viewMixin);

                console.log('resize call');
                var tableContainerWidth = viewMixin.$('.table-responsive table').parent().width();
                var tableWidth = viewMixin.$('.table-responsive table').width();

                checkToToggleStackedDisplay(viewMixin, thresholds, tableContainerWidth, tableWidth);
            }
        });


        /**
         * @mixin responsivelist
         */
        var mixin = Mixin('responsivelist', {
            partials: {
                //TODO check if still used
                alternativeListDisplay: ['groupedrowslistlayout'],
                subRowInformation: ['stackedcolumns']
            },

            stackedColumns: Ember.A(),

            groupedColumnsBindingActivator: 0,

            init:function() {
                console.log('init responsivelist');
                this.addMixinView(viewMixin);

                var mixinsOptions = get(this, 'content.mixins');

                if(mixinsOptions) {
                    var responsivelistOptions = get(this, 'content.mixins').findBy('name', 'responsivelist');
                    this.mixinOptions.responsivelist = responsivelistOptions;
                }

                this._super.apply(this, arguments);
            }
        });

        application.register('mixin:responsivelist', mixin);
    }
});

