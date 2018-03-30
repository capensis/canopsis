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
if(window.bricks.charts.envMode === "production") {
    require.config({
        paths: {
            'flotchart': 'canopsis/charts/lib/externals/flot/jquery.flot',
            'flotchart-canvas': 'canopsis/charts/lib/externals/flot/jquery.flot.canvas',
            'flotchart-categories': 'canopsis/charts/lib/externals/flot/jquery.flot.categories',
            'flotchart-crosshair': 'canopsis/charts/lib/externals/flot/jquery.flot.crosshair',
            'flotchart-errorbars': 'canopsis/charts/lib/externals/flot/jquery.flot.errorbars',
            'flotchart-fillbetween': 'canopsis/charts/lib/externals/flot/jquery.flot.fillbetween',
            'flotchart-image': 'canopsis/charts/lib/externals/flot/jquery.flot.image',
            'flotchart-navigate': 'canopsis/charts/lib/externals/flot/jquery.flot.navigate',
            'flotchart-pie': 'canopsis/charts/lib/externals/flot/jquery.flot.pie',
            'flotchart-resize': 'canopsis/charts/lib/externals/flot/jquery.flot.resize',
            'flotchart-selection': 'canopsis/charts/lib/externals/flot/jquery.flot.selection',
            'flotchart-stack': 'canopsis/charts/lib/externals/flot/jquery.flot.stack',
            'flotchart-symbol': 'canopsis/charts/lib/externals/flot/jquery.flot.symbol',
            'flotchart-threshold': 'canopsis/charts/lib/externals/flot/jquery.flot.threshold',
            'flotchart-time': 'canopsis/charts/lib/externals/flot/jquery.flot.time',
            'flotchart-valuelabel': 'canopsis/charts/lib/externals/flot-plugins/custom/jquery.flot.valuelabel',
            'flotchart-tooltip': 'canopsis/charts/lib/externals/flot.tooltip/js/jquery.flot.tooltip',
            'flotchart-chartvalues': 'canopsis/charts/lib/externals/flot-plugins/custom/jquery.flot.chartvalues',
            'c3': 'canopsis/charts/lib/externals/c3/c3.min'
        },

        shim: {
            'flotchart': {
                deps: ['jquery'],
            },

            'flotchart-canvas': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-categories': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-crosshair': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-errorbars': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-fillbetween': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-image': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-navigate': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-pie': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-resize': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-selection': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-stack': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-symbol': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-threshold': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-time': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-valuelabel': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-tooltip': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-chartvalues': {
                deps: ['jquery', 'flotchart']
            }
        }
    });

    define([
        'canopsis/charts/src/libwrappers/flotchart',
        'c3',
        'link!canopsis/charts/lib/externals/c3/c3.min.css',
        'flotchart',
        'flotchart-canvas',
        'flotchart-categories',
        'flotchart-crosshair',
        'flotchart-errorbars',
        'flotchart-fillbetween',
        'flotchart-image',
    //    'flotchart-navigate',
        'flotchart-pie',
        'flotchart-resize',
        'flotchart-selection',
        'flotchart-stack',
        'flotchart-symbol',
        'flotchart-threshold',
        'flotchart-time',
        'flotchart-valuelabel',
        'flotchart-tooltip',
        'flotchart-chartvalues'
    ], function (flotchart, c3) {
        window.c3 = c3;
    });
} else {
    require.config({
        paths: {
            'flotchart': 'canopsis/charts/lib/externals/flot/jquery.flot',
            'flotchart-canvas': 'canopsis/charts/lib/externals/flot/jquery.flot.canvas',
            'flotchart-categories': 'canopsis/charts/lib/externals/flot/jquery.flot.categories',
            'flotchart-crosshair': 'canopsis/charts/lib/externals/flot/jquery.flot.crosshair',
            'flotchart-errorbars': 'canopsis/charts/lib/externals/flot/jquery.flot.errorbars',
            'flotchart-fillbetween': 'canopsis/charts/lib/externals/flot/jquery.flot.fillbetween',
            'flotchart-image': 'canopsis/charts/lib/externals/flot/jquery.flot.image',
            'flotchart-navigate': 'canopsis/charts/lib/externals/flot/jquery.flot.navigate',
            'flotchart-pie': 'canopsis/charts/lib/externals/flot/jquery.flot.pie',
            'flotchart-resize': 'canopsis/charts/lib/externals/flot/jquery.flot.resize',
            'flotchart-selection': 'canopsis/charts/lib/externals/flot/jquery.flot.selection',
            'flotchart-stack': 'canopsis/charts/lib/externals/flot/jquery.flot.stack',
            'flotchart-symbol': 'canopsis/charts/lib/externals/flot/jquery.flot.symbol',
            'flotchart-threshold': 'canopsis/charts/lib/externals/flot/jquery.flot.threshold',
            'flotchart-time': 'canopsis/charts/lib/externals/flot/jquery.flot.time',
            'flotchart-valuelabel': 'canopsis/charts/lib/externals/flot-plugins/custom/jquery.flot.valuelabel',
            'flotchart-tooltip': 'canopsis/charts/lib/externals/flot.tooltip/js/jquery.flot.tooltip',
            'flotchart-chartvalues': 'canopsis/charts/lib/externals/flot-plugins/custom/jquery.flot.chartvalues',
            'c3': 'canopsis/charts/lib/externals/c3/c3.min'
        },

        shim: {
            'flotchart': {
                deps: ['jquery'],
            },

            'flotchart-canvas': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-categories': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-crosshair': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-errorbars': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-fillbetween': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-image': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-navigate': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-pie': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-resize': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-selection': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-stack': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-symbol': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-threshold': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-time': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-valuelabel': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-tooltip': {
                deps: ['jquery', 'flotchart']
            },

            'flotchart-chartvalues': {
                deps: ['jquery', 'flotchart']
            }
        }
    });

    define([
        'canopsis/charts/src/libwrappers/flotchart',
        'c3',
        'link!canopsis/charts/lib/externals/c3/c3.css',
        'flotchart',
        'flotchart-canvas',
        'flotchart-categories',
        'flotchart-crosshair',
        'flotchart-errorbars',
        'flotchart-fillbetween',
        'flotchart-image',
    //    'flotchart-navigate',
        'flotchart-pie',
        'flotchart-resize',
        'flotchart-selection',
        'flotchart-stack',
        'flotchart-symbol',
        'flotchart-threshold',
        'flotchart-time',
        'flotchart-valuelabel',
        'flotchart-tooltip',
        'flotchart-chartvalues'
    ], function (flotchart, c3) {
        window.c3 = c3;
    });
}
