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
        'components/component-querybuilder': 'canopsis/brick-querybuilder/src/components/querybuilder/template',
        'editor-querybuilder': 'canopsis/brick-querybuilder/src/editors/editor-querybuilder',

    }
});

 define([
    'canopsis/brick-querybuilder/src/components/querybuilder/component',
    'ehbs!components/component-querybuilder',
    'ehbs!editor-querybuilder',
    'link!canopsis/brick-querybuilder/src/style.css',
    'canopsis/brick-querybuilder/requirejs-modules/externals.conf'
], function (templates) {
    templates = $(templates).filter('script');
for (var i = 0, l = templates.length; i < l; i++) {
var tpl = $(templates[i]);
Ember.TEMPLATES[tpl.attr('data-template-name')] = Ember.Handlebars.compile(tpl.text());
};
});

