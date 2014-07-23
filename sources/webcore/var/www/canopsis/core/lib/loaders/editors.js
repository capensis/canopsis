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

var editorsTemplates = [
	/* js = 'cvw' : the editor have a Controller, a View and a Webcomponent
	 * js = 'cv' : the editor have a Controller and a View
	 * js = 'w' : the editor have a Webcomponent
	 */
	{ name: 'defaultpropertyeditor', js: 'v' },
	{ name: 'boolean' },
	{ name: 'group', js: 'c' },
	{ name: 'color' },
	{ name: 'array', js: 'w'},
	{ name: 'mail' },
	{ name: 'richtext', js: 'w' },
	{ name: 'timestamp', js: 'v' },
	{ name: 'rights', js: 'v' },
	{ name: 'cmetric', js: 'w' },
	{ name: 'cfilter', js: 'w' },
	{ name: 'cfilterwithproperties'},
	{ name: 'templateSelector' , js: "v" },
	{ name: 'tags' , js: "w" },
	{ name: 'eventselector', js: 'w' },
	{ name: 'state', js: 'w' },
	{ name: 'criticity', js: 'w' },
	{ name: 'actionfilter', js: 'w' },
	{ name: 'simplelist', js: 'v' }
];

var deps = ['ember'];

var depsTemplates = [];

//generate deps
for (var i = 0; i < editorsTemplates.length; i++) {
	var name = editorsTemplates[i].name;
	var files = editorsTemplates[i].js;

	var tmplPos;

	if (files !== undefined) {
		var url;

		if (files.indexOf('c') >= 0) {
			url = 'app/editors/' + name + '/controller';

			deps.push(url);
		}

		if (files.indexOf('v') >= 0) {
			url = 'app/editors/' + name + '/view';

			deps.push(url);
		}

		if (files.indexOf('w') >= 0) {
			url = 'text!app/editors/' + name + '/component.html';

			tmplPos = deps.push(url);
			depsTemplates.push({name: 'components/component-' + name, pos: tmplPos});

			url = 'app/editors/' + name + '/component';
			deps.push(url);
		}
	}

	tmplPos = deps.push('text!app/editors/' + name + '/template.html');
	depsTemplates.push({name: 'editor-' + name, pos: tmplPos});

}

console.log({"editors dependencies": deps});

define(deps, function(Ember) {
	console.log("load editors", arguments);

	for (var i = 0; i < depsTemplates.length; i++) {
		var tmplInfo = depsTemplates[i];

		var template = arguments[tmplInfo.pos - 1];

		console.log("new editor", template, tmplInfo.name);

		Ember.TEMPLATES[tmplInfo.name] = Ember.Handlebars.compile(template);
	}
});

