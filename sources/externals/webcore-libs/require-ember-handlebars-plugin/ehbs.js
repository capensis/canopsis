define(["exports"], function(exports) {
  "use strict";

  var codeLoadingHelpers = ["view", "template", "render"];
  var ignore, camelize, underscore, classify, path, ext;

  var paths = {
    views: "views/",
    templates: "templates/",
    controllers: "controllers/",
    helpers: "helpers/"
  };
  var casing = "camel";

  function enforceCase(str) {
    if (casing === "snake" || casing === "underscore") {
      return underscore(str);
    } else if (casing === "class") {
      return classify(str);
    } else {
      return camelize(str);
    }
  }

  function join(uriParts) {
    var parts = [];
    uriParts.forEach(function (part) {
      parts.push.apply(parts, part.split("/"));
    });
    return parts.filter(function(part) {
      return !!part;
    }).join("/");
  }

  function readConfig(config) {
    if (config.ehbs) {
      if (config.ehbs.paths) {
        ["views", "templates", "controllers", "helpers"].forEach(function (type) {
          if (config.ehbs.paths.hasOwnProperty(type)) {
            paths[type] = config.ehbs.paths[type];
          }
        });
      }

      if (config.ehbs.casing) {
        casing = config.ehbs.casing;
      }
    }
  }

  function getNamespaceAndNameFromStatement(statement) {
    if (statement.params[0]) {
      var parts = statement.params[0].string.split(".");
      var namespace;
      var name;

      if (parts.length === 1) {
        namespace = null;
        name = parts[0];
      } else {
        namespace = parts.shift();
        name = parts.join(".");
      }

      return [namespace, name];
    } else {
      return [null, statement.id.string];
    }
  }

  function isIgnoreableNamespace(namespace) {
    return ["Ember", "Em", null].indexOf(namespace) !== -1;
  }

  function shouldIgnore(helper, namespace) {
    if (ignore.indexOf(helper) !== -1) {
      if (codeLoadingHelpers.indexOf(helper) !== -1 && !isIgnoreableNamespace(namespace)) {
        return false;
      } else {
        return true;
      }
    } else if (namespace) {
      return isIgnoreableNamespace(namespace);
    } else {
      return false;
    }
  }

  function getDeps(ast, parentRequire) {
    var deps = ast.statements.reduce(function(deps, statement) {
      var helperName = statement.id && statement.id.string;
      var dep, parts, namespace, arg, nextStatement;

      if (statement.isHelper) {
        parts = getNamespaceAndNameFromStatement(statement);
        namespace = parts[0];
        arg = parts[1];

        if (!shouldIgnore(helperName, namespace)) {
          if (helperName == "view") {
            path = paths.views;
            ext = ".js";
          } else if (helperName == "partial" || helperName == "render") {
            path = "ehbs!";
            ext = "";
          } else if (helperName == "control") {
            path = paths.controllers;
            ext = ".js";
          } else {
            path = paths.helpers;
            arg = statement.id.string;
            ext = ".js";
          }

          arg = enforceCase(arg);

          deps.push(parentRequire.toUrl(path + arg + ext));
        }
      }

      nextStatement = Ember.get(statement, "program.statements");
      if (nextStatement) {
        deps.push.apply(deps, getDeps(statement.program, parentRequire));
      }

      nextStatement = Ember.get(statement, "program.inverse.statements");
      if (nextStatement) {
        deps.push.apply(deps, getDeps(statement.program.inverse, parentRequire));
      }

      return deps;
    }, []);

    return deps.filter(function(dep) {
      return !(parentRequire.defined(dep) || parentRequire.specified(dep));
    });
  }

  var compileOptions = {
    data: true,
    stringParams: true
  };

  exports.load = function(name, parentRequire, onload, config) {
    var parts = name.split(":");
    var path;
    if (parts.length == 2) {
      path = parts[0];
      name = parts[1];
    } else {
      path = name = parts[0];
    }

    // Bail out early during build.
    if (config.isBuild) {
      return onload();
    }

    readConfig(config);
    parentRequire(["text!" + join([paths.templates, path + ".hbs"]), "ember"], function (template, Ember) {
      // Set these from Ember now.
      ignore = Ember.keys(Ember.Handlebars.helpers);
      camelize = Ember.String.camelize;
      underscore = Ember.String.underscore;
      classify = Ember.String.classify;

      var ast = Ember.Handlebars.parse(template);
      var deps = getDeps(ast, parentRequire);

      // This stuff is taken right from Ember.Handlebars.compile()
      var environment = new Ember.Handlebars.Compiler().compile(ast, compileOptions);
      var templateSpec = new Ember.Handlebars.JavaScriptCompiler().compile(environment, compileOptions, undefined, true);
      Ember.TEMPLATES[name] = Ember.Handlebars.template(templateSpec);

      if (deps.length) {
        parentRequire(deps, onload);
      } else {
        onload();
      }
    });
  };

  exports.write = function(pluginName, moduleName, write) {
    write("");
  };
});
