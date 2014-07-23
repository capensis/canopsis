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

var utilsFiles = [
	'forms',
	'routes',
	'data',
	'notification',
	'i18n',
	'hash',
	'dates'
];

var deps = [];

for (var i = 0; i < utilsFiles.length; i++) {
	deps.push('app/lib/utils/' + utilsFiles[i]);
}

define(deps, function() {
	var utils = {};
	console.log("Begin load utils", arguments);
	for (var i = 0; i < arguments.length; i++) {
		var utilName = utilsFiles[i];
		console.log("load util", utilName);
		utils[utilName] = arguments[i];
	}
	window.ctools = utils;
	return utils;
});
