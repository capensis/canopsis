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

var helpers = [
    'color',
    'colorview',
    'date-fromnow',
    'duration',
    'editor',
    'enableview',
    'format-date',
    'group',
    'glyphicon',
    'log',
    'percent',
    'renderwidget',
    'stateview',
    'criticityview',
    'timestamp',
    'timeSince',
    'menu',
    'rights',
    'widgetslot',
    // 'renderer', //Deprecated
    'ifcond',
    'set',
    'getfield',
    'unset',
    'i18n',
    'tooltip',
    'json2html',
    'ack'
];

var deps = ['app/application', 'app/lib/wrappers/swag'];

for (var i = 0; i < helpers.length; i++) {
    deps.push("app/lib/helpers/" + helpers[i]);
}

define(deps, function() {

    return helpers;
});