Plugin system as bricks
=======================

How to add and configure plugins
--------------------------------

### enabled.json

_Must be located at the root of the directory `plugins/`_

Simple list of enabled plugins
No specific order needed

### manifest.json

_Must be located in the directory `files/` at the root of every plugin directory_
_All `.json` files must be in the same directory_

* name: Name of the plugin ot be usen
* version: Version of the plugin loaded
* routes: Routes to be initalized (see [EmberJS](http://emberjs.com/))
* files: Files to be required (see [RequireJS](http://requirejs.org/))
* dependancies: Plugins needed (they will be loaded first)

The order of the routes and files is kept.

### routes.json

* type: type of the node
* name: name of the node
* description: description of the node
* appears_on: deprecated
* children: children nodes

If a plugin defines a same route as an other,
only the one from the plugin to be loaded first is kept.

_For more informations see [EmberJS doc.](http://emberjs.com/guides/routing/)_

### files.json

* type: `dir` or undefined
* name: name of the file alone
* files: files in the directory if the type is `dir`


How to load enabled plugins
---------------------------

### Dependancies

```javascript
require([
         'plugins',
         'text!plugin_x/files/manifest.json',
         'text!plugin_x/files/files.json',
         'text!plugin_x/files/routes.json'],
         function() { } );
        
```

### Getting the plugins and resolved dependancies

```javascript
var plugins = [];

try {
    plugins = Plugins.getPlugins("./plugins/");
    plugins = Plugins.resolveDependancies(plugins);
} catch (e) {
    console.log("PluginError: " + e);
}
```

### Getting the routes

```javascript
var routes = [];

routes = Manifest.fetchRoutes(plugins);
```

### Getting the files

```javascript
var files = []

files = Manifest.fetchFiles(plugins);
```
