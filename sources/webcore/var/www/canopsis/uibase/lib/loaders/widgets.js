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

var uibaseWidgetsTemplates = [
	{ name:'weather', url:'canopsis/uibase/widgets/weather', hasJSPart: true },
	{ name:'text', url:'canopsis/uibase/widgets/text', hasJSPart: true }
];

var deps = ['ember'];
var jsDeps = [];
var depsSize = deps.length;

//generate deps
for (var i = 0; i < uibaseWidgetsTemplates.length; i++) {
	deps.push('text!' + uibaseWidgetsTemplates[i].url + '/template.html');

	if (uibaseWidgetsTemplates[i].hasJSPart === true) {
		var viewUrl = uibaseWidgetsTemplates[i].url + '/controller';
		console.log("adding view", viewUrl);

		jsDeps.push(viewUrl);
	}
}

for (i = 0; i < jsDeps.length; i++) {
	deps.push(jsDeps[i]);
}

console.log({"uibase widget dependencies": deps});
define(deps, function(Ember) {
	console.log("load widgets from uibase", arguments);
	for (var i = 0; i < uibaseWidgetsTemplates.length; i++) {
		var templateName = uibaseWidgetsTemplates[i].name;
		Ember.TEMPLATES[templateName] = Ember.Handlebars.compile(arguments[i + depsSize]);
	}
});

