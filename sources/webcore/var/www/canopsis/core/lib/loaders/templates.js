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

var templates = [
    { name: 'application' },
    { name: 'item' },
    { name: 'notifications' },
    { name: 'containerwidget' },
    { name: 'containervbox' },
    { name: 'containerhbox' },
    { name: 'container2' },
    { name: 'crecordform' },
    { name: 'formwrapper' },
    { name: 'menu' },
    { name: 'widgettitlebar' },
    { name: 'userview' },
    { name: 'widget' },
    { name: 'listline' },
    { name: 'widgetslot-default' },
    { name: 'widgetslot-grey' },
    { name: 'actionbutton-ack', classes: ["action"], icon : "ok" , label : "Ack"},
    { name: 'actionbutton-cancel', classes: ["action"], icon : "ban-circle" , label : "Cancel"},
    { name: 'actionbutton-remove', classes: ["action"], icon : "trash", label : "Remove"},
    { name: 'actionbutton-show', classes: ["action"], icon : "eye-open" ,label : "Show" },
    { name: 'actionbutton-info', classes: ["action"], icon : "info-sign" ,label : "Info"},
    { name: 'actionbutton-create', classes: ["action"], icon : "plus-sign" , label : "Create" },
    { name: 'actionbutton-removeselection', classes: ["action", "toolbar"], icon : "trash"  , label : "Remove-selection" },
    { name: 'actionbutton-history', classes: ["action"],icon : "time"  , label : "History" },
    { name: 'formbutton-submit', classes: ["formbutton"]   },
    { name: 'titlebarbutton-moveup', classes: ["formbutton"]  },
    { name: 'titlebarbutton-movedown', classes: ["formbutton"]  },
    { name: 'titlebarbutton-moveleft', classes: ["formbutton"]  },
    { name: 'titlebarbutton-moveright', classes: ["formbutton"] },
    { name: 'titlebarbutton-minimize', classes: ["formbutton"]   },
    { name: 'formbutton-cancel', classes: ["formbutton"]  }
];

var deps = ['ember'];
var depsSize = deps.length;

for (var i = 0; i < templates.length; i++) {
	deps.push('text!app/templates/' + templates[i].name + '.html');
}

define(deps, function(Ember) {
	var templatesLoaded = {};
	templatesLoaded.all = [];
	templatesLoaded.byClass = {};

	for (var i = depsSize; i < arguments.length; i++) {
		var currentTemplate = templates[i - depsSize];
		Ember.TEMPLATES[currentTemplate.name] = Ember.Handlebars.compile(arguments[i]);

		if (currentTemplate.classes !== undefined) {
			for (var j = 0; j < currentTemplate.classes.length; j++) {
				var currentClass = currentTemplate.classes[j];

				if (templatesLoaded.byClass[currentClass] === undefined) {
					templatesLoaded.byClass[currentClass] = [];
				}

				templatesLoaded.byClass[currentClass].push(currentTemplate);
			}
		}

	    templatesLoaded.all.push(currentTemplate);
	}

	return templatesLoaded;
});
