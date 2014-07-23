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

var factories = [
	{ name:'editor', url: 'app/lib/factories/editor' },
	{ name:'widget', url: 'app/lib/factories/widget' },
	{ name:'form', url: 'app/lib/factories/form' }
];

var factoriesDeps = ['app/application'];
var factoriesDepsSize = factoriesDeps.length;

for (var i = 0; i < factories.length; i++) {
	factoriesDeps.push(factories[i].url);
}

define(factoriesDeps, function(Application) {
	console.log(arguments);
	Application.factories = {};

	console.log("loading factories", factories, "into", Application.factories);

	for (var i = 0; i < factories.length; i++) {
		Application.factories[factories[i].name.capitalize()] = arguments[i + factoriesDepsSize];
	}

	return Application.factories;
});