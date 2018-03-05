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

Ember.Handlebars.helper('stateview', function(state) {
    var __ = Ember.String.loc;


    var statelist = {
        0: {color: 'green', text: 'Info'},
        1: {color: 'yellow', text: 'Minor'},
        2: {color: 'orange', text: 'Major'},
        3: {color: 'red', text: 'Critical'}
    };

    var stateSelection = statelist[state];
    var state_template = '<span class="badge bg-%@">%@</span>'.fmt(stateSelection.color, __(stateSelection.text));

    return new Ember.Handlebars.SafeString(state_template);
});
