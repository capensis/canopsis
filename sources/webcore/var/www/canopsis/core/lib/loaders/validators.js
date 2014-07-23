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

var validatorsArray = [
	'mail',
	'rights',
	'required',
	'validate'
];

var deps = ['ember'];

for (var i = 0; i < validatorsArray.length; i++) {
	var validatorUrl = 'app/validators/' + validatorsArray[i] + '/validator';
	deps.push(validatorUrl);
}

define(deps, function(Ember) {
	var validators = {};
	console.log("Begin load validators", arguments);
	for (var i = 1; i < arguments.length; i++) {
		var validatorName = validatorsArray[i-1];
		console.log("load validator", validatorName);
		validators[validatorName] = arguments[i];
	}
	Ember.validators = validators;
	return validators;

});
