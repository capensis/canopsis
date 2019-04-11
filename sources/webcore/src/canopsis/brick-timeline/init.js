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
        'components/component-timeline': 'canopsis/brick-timeline/src/components/timeline/template',

    }
});

define([
    'canopsis/brick-timeline/src/adapters/timelineadapter',
    'canopsis/brick-timeline/src/components/timeline/component',
    'ehbs!components/component-timeline',
    'link!canopsis/brick-timeline/src/style.css',
    'canopsis/brick-timeline/requirejs-modules/i18n'
], function () {
    
});
