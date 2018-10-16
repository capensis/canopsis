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
        'doT': 'canopsis/brick-querybuilder/externals/doT/doT.min',
        'jQuery.extendext': 'canopsis/brick-querybuilder/externals/jquery-extendext/jQuery.extendext.min',
        'query-builder': 'canopsis/brick-querybuilder/externals/jQuery-QueryBuilder/dist/js/query-builder',
        'querybuilder-editablekey' : 'canopsis/brick-querybuilder/externals/querybuilder-editablekey/querybuilder-editablekey'
    },
    shim: {
        'querybuilder-editablekey': {
            'deps': ['query-builder', 'jquery']
        }
    }
});

define([
    'query-builder',
    'link!canopsis/brick-querybuilder/externals/jQuery-QueryBuilder/dist/css/query-builder.default.min.css',
    'canopsis/brick-querybuilder/externals/jquery-editable-select/source/jquery.editable-select',
    'link!canopsis/brick-querybuilder/externals/jquery-editable-select/source/jquery.editable-select.min.css',
    'canopsis/brick-querybuilder/externals/querybuilder-editablekey/querybuilder-editablekey'
], function () {
    require(['canopsis/brick-querybuilder/externals/jQuery-QueryBuilder/dist/i18n/query-builder.' + window.i18n.lang]);
});

