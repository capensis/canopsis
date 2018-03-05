### requirejs-ember-handlebars
A requirejs plugin for using ember templates

		$ bower install requirejs-ember-handlebars
		$ npm install ember-template-compiler
		
#### Example

	define(
	
		[
			"Ember",
			"App",
			"ehbs!example"
		],
	
		function (Ember, App, template) {
	
			"use strict";
	
			App.ExampleRoute = Ember.Route.extend({});
	
			return App.ExampleRoute;
		}
	);
	
##### What it does

Loads in your "example" template and calls `Ember.TEMPLATES["example"] = Ember.Handlebars.compile(template);` for you automatically.

When you build your project, instead of calling `Ember.Handlebars.compile()` it actually uses [ember-template-compiler](https://github.com/toranb/ember-template-compiler) and precompiles your template for you, then includes it as part of the minified js.


#### Configuration

		require.config({

		    paths: {
		        "jQuery": "//ajax.googleapis.com/ajax/libs/jquery/2.0.3/jquery",
		        "text": "path/to/bower_components/requirejs-text/text",
		        "ehbs" : "path/to/bower_components/requirejs-ember-handlebars/ehbs",
		        "Ember" : "path/to/bower_components/ember/ember",
		        "Handlebars": "path/to/bower_components/handlebars/handlebars"
		    },

		    shim: {
		        "jQuery": {
		            exports: "jQuery"
		        },
		
		        "Handlebars": {
		            exports: "Handlebars"
		        },
		
		        "Ember": {
		            deps: ["jQuery", "Handlebars"],
		            exports: "Ember"
		        }
		    },

		    ehbs : {
		    	extension : "html",					// default : "hbs"
		        templatePath : "app/templates/",	// default : ""
		        ember : "Ember" 					// default : "Ember"
		    }
		});
		
The `paths` and `shim` properties of the config are how I get Ember and AMD to co-exist. The `ehbs` property lists the available config options for this plugin.

##### extension

Specify the file extension you want to use for templates.

##### templatePath

Where your templates live. This is relative to baseUrl.

##### ember

This is where you specify the AMD identifier you are using for Ember. I use "Ember" but you can change this to "Em" or "ember", or whatever else you want to use.
