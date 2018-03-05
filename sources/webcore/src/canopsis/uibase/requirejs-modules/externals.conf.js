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

if(window.bricks.uibase.envMode === "production") {
    require.config({
        paths: {
            'bootstrap': 'canopsis/uibase/lib/externals/bootstrap/dist/js/bootstrap.min',
            'datetimepicker': 'canopsis/uibase/libwrappers/datetimepicker',
            'icheck': 'canopsis/uibase/lib/externals/iCheck/icheck.min',
            'codemirror': 'canopsis/uibase/lib/externals/codemirror/lib/codemirror',
            'summernote': 'canopsis/uibase/lib/externals/summernote/dist/summernote',
            'ember-summernote': 'canopsis/uibase/lib/externals/ember-summernote/lib/component',
            'daterangepicker': 'canopsis/uibase/lib/externals/bootstrap-daterangepicker/daterangepicker',
            'rrule': 'canopsis/uibase/lib/externals/rrule/lib/rrule',
            'nlp': 'canopsis/uibase/lib/externals/rrule/lib/nlp',
            'underscore' : 'canopsis/uibase/libwrappers/underscore',

            'moment': 'canopsis/uibase/lib/externals/moment/min/moment-with-locales.min',
            'jsoneditorlib': 'canopsis/uibase/lib/externals/jsoneditor/jsoneditor.min',
            'ember-jsoneditor-lib': 'canopsis/uibase/lib/externals/ember-jsoneditor/ember-jsoneditor',
            'd3': 'canopsis/uibase/lib/externals/d3/d3.min'
        },
        shim: {
            'rrule': {
                deps: ['jquery', 'underscore']
            },
            'nlp': {
                'deps': ['jquery', 'rrule', 'underscore']
            },
            'bootstrap': {
                deps: ['jquery']
            },
            'datetimepicker': {
                deps: ['jquery', 'moment', 'bootstrap']
            },
            'icheck': {
                deps: ['jquery']
            }
        }
    });

    define([
        'd3',
        'moment',
        'jsoneditorlib',
        'link!canopsis/uibase/lib/externals/jsoneditor/jsoneditor.min.css',
        'link!canopsis/uibase/lib/externals/fontawesome/css/font-awesome.min.css',
        'link!canopsis/uibase/lib/externals/bootstrap-daterangepicker/daterangepicker-bs3.css',
        'canopsis/uibase/lib/externals/stacktable/stacktable',
        'link!canopsis/uibase/lib/externals/stacktable/stacktable.css',
        'daterangepicker',
        'canopsis/uibase/lib/externals/summernote/dist/summernote.min',
        'link!canopsis/uibase/lib/externals/codemirror/lib/codemirror.css',
        'canopsis/uibase/libwrappers/codemirror',
        'ember-summernote',
        'canopsis/uibase/lib/externals/underscore/underscore-min',
        'underscore',
        'nlp',
        'bootstrap',
        'canopsis/uibase/lib/externals/ion.rangeslider/js/ion.rangeSlider.min',
        'link!canopsis/uibase/lib/externals/ion.rangeslider/css/ion.rangeSlider.css',
        'link!canopsis/uibase/lib/externals/ion.rangeslider/css/ion.rangeSlider.skinHTML5.css',
        'canopsis/uibase/lib/externals/ember-datetimepicker/lib/component',
        'canopsis/uibase/lib/externals/ember-icheck/lib/component',
        'canopsis/uibase/lib/externals/ember-tooltip/lib/component',
        'canopsis/uibase/lib/externals/ember-durationcombo/lib/component',
        'link!canopsis/uibase/lib/externals/bootstrap/dist/css/bootstrap.min.css',
        'codemirror',
        'link!canopsis/uibase/lib/externals/codemirror/theme/ambiance.css',
        'link!canopsis/uibase/lib/externals/codemirror/lib/codemirror.css',
        'canopsis/uibase/lib/externals/colpick/js/colpick',
        'link!canopsis/uibase/lib/externals/colpick/css/colpick.css',
        'link!canopsis/uibase/lib/externals/eonasdan-bootstrap-datetimepicker/lib/css/bootstrap-datetimepicker.min.css',
        'canopsis/uibase/lib/externals/eonasdan-bootstrap-datetimepicker/lib/js/bootstrap-datetimepicker.min',
        'canopsis/uibase/lib/externals/iCheck/icheck',
        'link!canopsis/uibase/lib/externals/iCheck/skins/all.css'
    ], function (d3, moment, jsoneditor) {
        window.d3 = d3;
        window.moment = moment;
        window.jsoneditor = { JSONEditor: jsoneditor };

        require(['ember-jsoneditor-lib'], function() {});

        require(['rrule'], function () {});
    });
} else {
        require.config({
        paths: {
            'bootstrap': 'canopsis/uibase/lib/externals/bootstrap/dist/js/bootstrap.min',
            'datetimepicker': 'canopsis/uibase/libwrappers/datetimepicker',
            'icheck': 'canopsis/uibase/lib/externals/iCheck/icheck',
            'codemirror': 'canopsis/uibase/lib/externals/codemirror/lib/codemirror',
            'summernote': 'canopsis/uibase/lib/externals/summernote/dist/summernote',
            'ember-summernote': 'canopsis/uibase/lib/externals/ember-summernote/lib/component',
            'daterangepicker': 'canopsis/uibase/lib/externals/bootstrap-daterangepicker/daterangepicker',
            'rrule': 'canopsis/uibase/lib/externals/rrule/lib/rrule',
            'nlp': 'canopsis/uibase/lib/externals/rrule/lib/nlp',
            'underscore' : 'canopsis/uibase/libwrappers/underscore',

            'moment': 'canopsis/uibase/lib/externals/moment/min/moment-with-locales.min',
            'jsoneditorlib': 'canopsis/uibase/lib/externals/jsoneditor/jsoneditor',
            'ember-jsoneditor-lib': 'canopsis/uibase/lib/externals/ember-jsoneditor/ember-jsoneditor',
            'd3': 'canopsis/uibase/lib/externals/d3/d3'
        },
        shim: {
            'rrule': {
                deps: ['jquery', 'underscore']
            },
            'nlp': {
                'deps': ['jquery', 'rrule', 'underscore']
            },
            'bootstrap': {
                deps: ['jquery']
            },
            'datetimepicker': {
                deps: ['jquery', 'moment', 'bootstrap']
            },
            'icheck': {
                deps: ['jquery']
            }
        }
    });

    define([
        'd3',
        'moment',
        'jsoneditorlib',
        'link!canopsis/uibase/lib/externals/jsoneditor/jsoneditor.css',
        'link!canopsis/uibase/lib/externals/fontawesome/css/font-awesome.min.css',
        'link!canopsis/uibase/lib/externals/bootstrap-daterangepicker/daterangepicker-bs3.css',
        'canopsis/uibase/lib/externals/stacktable/stacktable',
        'link!canopsis/uibase/lib/externals/stacktable/stacktable.css',
        'daterangepicker',
        'canopsis/uibase/lib/externals/summernote/dist/summernote.min',
        'link!canopsis/uibase/lib/externals/codemirror/lib/codemirror.css',
        'canopsis/uibase/libwrappers/codemirror',
        'ember-summernote',
        'canopsis/uibase/lib/externals/underscore/underscore',
        'underscore',
        'nlp',
        'bootstrap',
        'canopsis/uibase/lib/externals/ion.rangeslider/js/ion.rangeSlider.min',
        'link!canopsis/uibase/lib/externals/ion.rangeslider/css/ion.rangeSlider.css',
        'link!canopsis/uibase/lib/externals/ion.rangeslider/css/ion.rangeSlider.skinHTML5.css',
        'canopsis/uibase/lib/externals/ember-datetimepicker/lib/component',
        'canopsis/uibase/lib/externals/ember-icheck/lib/component',
        'canopsis/uibase/lib/externals/ember-tooltip/lib/component',
        'canopsis/uibase/lib/externals/ember-durationcombo/lib/component',
        'link!canopsis/uibase/lib/externals/bootstrap/dist/css/bootstrap.min.css',
        'codemirror',
        'link!canopsis/uibase/lib/externals/codemirror/theme/ambiance.css',
        'link!canopsis/uibase/lib/externals/codemirror/lib/codemirror.css',
        'canopsis/uibase/lib/externals/colpick/js/colpick',
        'link!canopsis/uibase/lib/externals/colpick/css/colpick.css',
        'link!canopsis/uibase/lib/externals/eonasdan-bootstrap-datetimepicker/lib/css/bootstrap-datetimepicker.min.css',
        'canopsis/uibase/lib/externals/eonasdan-bootstrap-datetimepicker/lib/js/bootstrap-datetimepicker.min',
        'canopsis/uibase/lib/externals/iCheck/icheck',
        'link!canopsis/uibase/lib/externals/iCheck/skins/all.css'
    ], function (d3, moment, jsoneditor) {
        window.d3 = d3;
        window.moment = moment;
        window.jsoneditor = { JSONEditor: jsoneditor };

        require(['ember-jsoneditor-lib'], function() {});

        require(['rrule'], function () {});
    });
}
