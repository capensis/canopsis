/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis. If not, see <http://www.gnu.org/licenses/>.
*/

//TODO implement auto check for mvct file existence and require them automatically

var renderersTemplates = [
    'default',
    'mail',
    'boolean',
    'tags',
    'color',
    'state',
    'criticity',
    'timestamp',
    'percent',
    'ack',
    'crecord-type',
    'rights',
    'actionfilter',
    'cfilter',
    'cfilterwithproperties'
];

var deps = ['ember'];
var depsSize = deps.length;

for (var i = 0; i < renderersTemplates.length; i++) {
    deps.push('text!app/renderers/' + renderersTemplates[i] + '/template.html');
}

define(deps, function(Ember) {
    for (var i = depsSize; i < arguments.length; i++) {
	var templateName = "renderer-" + renderersTemplates[i - depsSize];
	Ember.TEMPLATES[templateName] = Ember.Handlebars.compile(arguments[i]);
    }
});