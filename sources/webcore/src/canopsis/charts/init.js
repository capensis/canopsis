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

 require.config({
    paths: {
        'components/component-c3categorychart': 'canopsis/charts/src/components/c3categorychart/template',
        'components/component-c3js': 'canopsis/charts/src/components/c3js/template',
        'components/component-flotchart': 'canopsis/charts/src/components/flotchart/template',
        'components/component-metricitemeditor': 'canopsis/charts/src/components/metricitemeditor/template',
        'components/component-metricselector': 'canopsis/charts/src/components/metricselector/template',
        'components/component-metricselector2': 'canopsis/charts/src/components/metricselector2/template',
        'components/component-selectedmetricheader': 'canopsis/charts/src/components/selectedmetricheader/template',
        'components/component-serieitemeditor': 'canopsis/charts/src/components/serieitemeditor/template',
        'editor-metricitem': 'canopsis/charts/src/editors/editor-metricitem',
        'editor-metricselector2': 'canopsis/charts/src/editors/editor-metricselector2',
        'editor-serieitem': 'canopsis/charts/src/editors/editor-serieitem',
        'titlebarbutton-resetzoom': 'canopsis/charts/src/templates/titlebarbutton-resetzoom',
        'categorychart': 'canopsis/charts/src/widgets/categorychart/categorychart',
        'timegraph': 'canopsis/charts/src/widgets/timegraph/timegraph',

    }
});

define([
    'canopsis/charts/src/components/c3categorychart/component',
    'ehbs!components/component-c3categorychart',
    'canopsis/charts/src/components/c3js/component',
    'ehbs!components/component-c3js',
    'canopsis/charts/src/components/flotchart/component',
    'ehbs!components/component-flotchart',
    'canopsis/charts/src/components/metricitemeditor/component',
    'ehbs!components/component-metricitemeditor',
    'canopsis/charts/src/components/metricselector/component',
    'ehbs!components/component-metricselector',
    'canopsis/charts/src/components/metricselector2/component',
    'ehbs!components/component-metricselector2',
    'canopsis/charts/src/components/selectedmetricheader/component',
    'ehbs!components/component-selectedmetricheader',
    'canopsis/charts/src/components/serieitemeditor/component',
    'ehbs!components/component-serieitemeditor',
    'ehbs!editor-metricitem',
    'ehbs!editor-metricselector2',
    'ehbs!editor-serieitem',
    'canopsis/charts/src/libwrappers/flotchart',
    'ehbs!titlebarbutton-resetzoom',
    'ehbs!categorychart',
    'canopsis/charts/src/widgets/categorychart/controller',
    'canopsis/charts/src/widgets/timegraph/controller',
    'ehbs!timegraph',
    'canopsis/charts/requirejs-modules/externals.conf'
], function () {
    
});
