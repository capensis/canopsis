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

var widgetsTemplates = [
	{ name:'list', url:'app/widgets/list', hasJSPart: true },
	{ name:'canvas', url:'app/widgets/canvas', hasJSPart: true },
	{ name:'hbox', url:'app/widgets/hbox', hasJSPart: true },
	{ name:'lighthbox', url:'app/widgets/lighthbox', hasJSPart: true },
	{ name:'vbox', url:'app/widgets/vbox', hasJSPart: true },
	{ name:'tabmanager', url:'app/widgets/tabmanager', hasJSPart: true }
];

var deps = ['ember'];
var jsDeps = [];
var depsSize = deps.length;


//generate deps
for (var i = 0; i < widgetsTemplates.length; i++) {
	deps.push('text!' + widgetsTemplates[i].url + '/template.html');

	if (widgetsTemplates[i].hasJSPart === true) {
		var viewUrl = widgetsTemplates[i].url + '/controller';
		console.log("adding view", viewUrl);

		jsDeps.push(viewUrl);
	}
}

for (i = 0; i < jsDeps.length; i++) {
	deps.push(jsDeps[i]);
}

console.log({"widget dependencies": deps});
define(deps, function(Ember) {
	console.log({"load widgets": arguments});
	for (var i = 0; i < widgetsTemplates.length; i++) {
		var templateName = widgetsTemplates[i].name;
		Ember.TEMPLATES[templateName] = Ember.Handlebars.compile(arguments[i + depsSize]);
	}
});

