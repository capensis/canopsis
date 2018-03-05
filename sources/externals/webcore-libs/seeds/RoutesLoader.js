define([
	'module',
	'jquery',
	'ember',
	'ember-data'
], function(module, $, Ember, DS) {

	availableFiles = [];

	var app;

	var isFileAvailable = function(type, name) {
		type += "s";
		if(availableFiles[type] !== undefined) {
			return availableFiles[type].contains(name);
		}
		else {
			return false;
		}
	};

	//Use requireJs to load files (such as models, views, controllers, templates) if they exists
	var getRequirementsOfRoute = function(name) {
		var requirements = [];

		if(isFileAvailable("model", name + ".js")) {
			requirements.push("app/model/" + name);
		}
		if(isFileAvailable("view", name + ".js")) {
			requirements.push("app/view/" + name);
		}
		if(isFileAvailable("controller", name + ".js")) {
			requirements.push("app/controller/" + name);
		}
		if(isFileAvailable("template", name + ".html")) {
			requirements.push('text!app/templates/' + name + '.html');
		}
		return requirements;
	};

	var getRequirementsList = function(currentRequirements, routes) {
		for (var i = 0; i <= routes.length - 1; i++) {
			var currentRoute = routes[i];

			var newRequirements = getRequirementsOfRoute(currentRoute.name);
			currentRequirements = currentRequirements.concat(newRequirements);

		}
		return currentRequirements;
	};

	var initializeRoutesRecursive = function(currentScope, routes, parents) {
		var childRoutes = function(){
			if(currentRoute.children !== undefined) {
				newParents.push({
					type: currentRoute.type,
					name: currentRoute.name,
					icon: currentRoute.icon,
					appears_on: currentRoute.appears_on
				});
				initializeRoutesRecursive(this,currentRoute.children, newParents);
			}
		};

		if(parents === undefined)
			parents = [];

		var newParents = [];
		var i;
		for (i = 0; i < parents.length; i++) {
			newParents.push(parents[i]);
		}

		for (i = 0; i <= routes.length - 1; i++) {
			var currentRoute = routes[i];
			// TODO ? automatic route object creation addRoute(currentScope, currentRoute.name, [], currentRoute);

			var options = {};
			if(currentRoute.path !== undefined){
				options.path = currentRoute.path;
			}
			console.log(currentRoute);
			if(currentRoute.type === "route") {
				currentScope.route(currentRoute.name, options, childRoutes);
			} else if (currentRoute.type === "resource") {
				currentScope.resource(currentRoute.name, options, childRoutes);
			}
		}
	};

	var requirements = [];

	var RoutesLoader = Ember.Object.extend({
		initializeFiles : function(routes, callback) {
			requirements = getRequirementsList([], routes);

			return require(requirements, function() {

				console.log("initialize templates");

				for (var i = 0; i < requirements.length; i++) {
					console.log("requirement: ", requirements[i]);
					if(requirements[i].slice(0,5) === "text!"){
						var name = requirements[i].split("templates/")[1].slice(0,-5);
						Ember.TEMPLATES[name] = Ember.Handlebars.compile(arguments[i]);
					}
				}

				callback();
			});
		},

		initializeRoutes : function(application, routes) {
			app = application;

			application.Router.map(function() {
				initializeRoutesRecursive(this, routes);
			});
		}
	});

	return RoutesLoader.create();
});
