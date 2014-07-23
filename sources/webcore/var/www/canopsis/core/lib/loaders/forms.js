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

var formsTemplates = [
	'modelform',
	'widgetform'
];

var deps = ['ember'];
var jsDeps = [];
var depsSize = deps.length;


//generate deps
for (var i = 0; i < formsTemplates.length; i++) {
	deps.push('text!app/forms/' + formsTemplates[i] + '/template.html');

	var controllerUrl = 'app/forms/' + formsTemplates[i] + '/controller';
	jsDeps.push(controllerUrl);
}

for (i = 0; i < jsDeps.length; i++) {
	deps.push(jsDeps[i]);
}

console.log({"form dependencies": deps});
define(deps, function(Ember) {
	console.log("load forms", arguments);
	for (var i = 0; i < formsTemplates.length; i++) {
		var templateName = formsTemplates[i];
		Ember.TEMPLATES[templateName] = Ember.Handlebars.compile(arguments[i + depsSize]);
	}
});

