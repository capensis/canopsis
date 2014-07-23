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

var routes;

require([
	'plugins',
	'text!canopsis/enabled.json',
	'text!canopsis/core/files/manifest.json',
	'text!canopsis/core/files/routes.json',
	'text!canopsis/core/files/files.json',
	'text!canopsis/uibase/files/manifest.json',
	'text!canopsis/uibase/files/routes.json',
	'text!canopsis/uibase/files/files.json'
], function(plugins_tool) {
	var plugins = [];
	var path = "canopsis/";
	var files;

	try {
		plugins = plugins_tool.Plugins.getPlugins(path);
		plugins = plugins_tool.Plugins.resolveDependancies(plugins);
	} catch (e) {
		console.log("PluginError: " + e);
	}

	routes = plugins_tool.Manifest.fetchRoutes(plugins, path);
	files = plugins_tool.Manifest.fetchFiles(plugins, path);
	files = files.map(function(e) {
		return e.replace("canopsis/core/", "app/");
	});
	require(files);
});

