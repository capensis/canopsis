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

Ember.Handlebars.helper('criticity', function(value, keep_state) {

    var __ = Ember.String.loc;
    //displays keep status information if any onto the state field
    //keep state is generated when a user overrides the criticity of and acknowleged event
    console.log('keep_state', keep_state);

    var display_keep_state = '';

    //default value: unknown
    if (value === undefined) {
        value = 4;
    }

    if (keep_state) {
        var color = '';
        switch(value) {
        case 0: color = 'bg-green'; break;
        case 1: color = 'bg-yellow'; break;
        case 2: color = 'bg-orange'; break;
        case 3: color = 'bg-red'; break;
        }

        display_keep_state = '<span class="badge '+ color +'"><i class="fa fa-male"></i></span>';
    }

    var span;
    switch(value) {
    case 0: span = '<span class="badge bg-green">'+__('Info')+'</span>'; break;
    case 1: span = '<span class="badge bg-yellow">'+__('Minor')+'</span>'; break;
    case 2: span = '<span class="badge bg-orange">'+__('Major')+'</span>'; break;
    case 3: span = '<span class="badge bg-red">'+__('Critical')+'</span>'; break;
    case 4: span = '<span class="badge">'+__('Unknown')+'</span>'; break;
    }
    return span + display_keep_state;

});
