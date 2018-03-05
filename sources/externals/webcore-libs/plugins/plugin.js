/*
# Copyright (c) 2015 "Capensis" [http://www.capensis.com]
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
define(function(require) {

    var plugins = {};

    Array.prototype.flatten = function() {
        var toReturn = this;
        if (toReturn.length > 0)
            toReturn =  (this.reduce(function(a, b) { return a.concat(b); }));
        return toReturn;
    };

    Array.prototype.findIndex = function (to_find) {
        for (var i in this)
        if (this[i].name == to_find) {
            return i;
        }
        return -1;
    };

    function formatPath(path) {
        return path.slice(-1) != '/' ? path + '/' : path;
    }

    plugins.Plugins = {

        /**
        * plug, dep_plug -> plugin objects
        * Check if two plugins depends on each other
        */
        checkConflicts: function(plug, dep_plug) {
        if (dep_plug.dependancies)
            for (var i = 0; i < dep_plug.dependancies.length; ++i)
            if (!~plug.dependancies.indexOf(dep_plug.dependancies[i]))
                return 'Conflicting dependancies between "' +
                    dep_plug.name +
                    '" and "' +
                    plug.name + '"';
        },

        /**
        * plugins -> list of plugin objects (see getPugins)
        * Sort the list to load the needed plugins first
        */
        resolveDependancies: function(plugins) {
        for (var i = 0; i < plugins.length; ++i) {
            var c_plugin = plugins[i];

            if (c_plugin.dependancies)
            for (var j = 0; j < c_plugin.dependancies.length; ++j) {
                var elem = c_plugin.dependancies[j];
                var p_pos = plugins.findIndex(c_plugin.name);
                var dep_pos = plugins.findIndex(elem);
                var error = plugins.Plugins.checkConflicts(c_plugin, plugins[dep_pos]);

                if (error) {
                    throw (error);
                }
                if (dep_pos == -1) {
                    throw ("Missing plugin: " + elem);
                }
                else if (dep_pos > p_pos) {
                    plugins.splice(dep_pos+1, 0, c_plugin);
                    plugins.splice(p_pos, 1);
                    i = -1;
                }
            }
        }
        return plugins;
        },

        /**
        * Get list of plugin objects
        */
        getPlugins: function(path, enabledModules) {
        path = formatPath(path);

        return (enabledModules.map(function(elem) {
            var path_ = "text!"+ path + elem +"/files/manifest.json";
            return (JSON.parse(require(path_)));
        }).flatten());
        }
    };

    /**
    * Test deep equality (i.e. every field besides arrays) between two objects
    */
    function areEquals(fst, scnd) {
        for (var key_ in fst) {
            if (!Array.isArray(fst[key_]) && fst[key_] != scnd[key_]) {
                return false;
            }
        }
        return true;
    }

    plugins.Manifest = {

        /**
        * original -> array to filter
        * Remove duplicates and merge sub-level arrays recursively
        */
        format: function(original) {
            var cObj, tObj = {};
            var merged = [];
            var unique = false;

            while (original.length) {
                cObj = original.pop();

                for (var i = 0; i < original.length; ++i) {
                    tObj = original[i];

                    for (var key in cObj) {
                        if (Array.isArray(cObj[key]) && areEquals(cObj, tObj)) {
                            tObj[key] = this.format(cObj[key].concat(tObj[key]));
                            original[i] = tObj;
                        }
                        else if (tObj[key] != cObj[key]) {
                            unique = true; break;
                        }
                        unique = false;
                    }
                    if (!unique) break;
                }
                if (unique) {
                    merged.push(cObj);
                }
                unique = true;
            }
            return merged.reverse();
        },

        /**
         * Gets the plugin route url file
         *
         * @param {dict} plugin the plugin dict
         * @param {string} path the program prefix
         *
         * @returns {string} the url of the file
         */
        getRouteFileUrlForPlugin: function(plugin, path){
            var routeFile = formatPath(path) + plugin.name +"/files/"+ plugin.routes;
            return routeFile;
        },

        /**
        * list -> list of plugin objects
        * Fetch JSON routes files from the manifest of each and every plugin
        * Merge them in a depth-1 array without duplicates
        */
        fetchRoutes: function(pluginsList, path) {
            console.log("Fetch routes for plugins", pluginsList);

            var routesListPerPlugin = [];

            for (var i = 0; i < pluginsList.length; i++) {
                var currentPlugin = pluginsList[i];
                var routeFileRequirement = "text!"+ this.getRouteFileUrlForPlugin(currentPlugin, path);

                var routeJson = JSON.parse(require(routeFileRequirement));
                routesListPerPlugin.push(routeJson);
            }

            var routesList = routesListPerPlugin.flatten();
            routesList = this.format(routesList);

            console.log("Final application routes list", routesList);

            return routesList;
        },

        /**
        * path -> root of the plugin, files_list -> JSON list of files
        * Go through the directories recursively and
        * return an array of the paths to the files to load
        */
        getFiles: function(path, files_list) {
            var files = [];

            for (var i = 0; i < files_list.length; ++i) {
                var file = files_list[i];

                if (file.type && file.type === 'dir' && file.files) {
                    var path_dir = path + file.name + '/';
                    var sub_dir = plugins.Manifest.getFiles(path_dir,  file.files);

                    files = files.concat(sub_dir);
                } else if (file.name) {
                    files.push(path + file.name);
                }
            }

            return files;
        },

        /**
        * list -> list of plugin objects
        * Fetch available files from the manifest of each and every plugin
        */
        fetchFiles: function(list, path_) {
            var files = [];
            path_ = formatPath(path_);

            for (var i = 0; i < list.length; ++i) {
                var plugin = list[i];
                var path = path_ + plugin.name + '/';

                console.log('Will load plugin ' +plugin.files);
                if (plugin.files) {
                    var file = require("text!"+ path +"files/"+ plugin.files);

                    file = JSON.parse(file);
                    file = plugins.Manifest.custom_files_rules(file, plugin.files);
                    console.log(file);

                    files.push(plugins.Manifest.getFiles(path, file));
                }
            }
            console.log(files);

            return files.length ? files.flatten() : [];
        },

        //special rules cases to apply to some particular files on load
        custom_files_rules: function (file, filename) {
            console.log(file);
            if (filename === 'files.json'){
                //Load only core sub file array
                plugins.FILES = file;
                for (var i=0; i<file.length; i++) {
                    if (file[i].name === 'core'){
                        console.log(file[i].name);
                        file = file[i].files;
                    }
                }
            }
            //MAYBE TODO some other custom cases
            return file;
        }

    };

    return plugins;
});
