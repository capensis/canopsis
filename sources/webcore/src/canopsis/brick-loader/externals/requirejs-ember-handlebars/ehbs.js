define(

	[
		"module",
		"text"
	],

	function (module, text) {

		var CONFIG,
			parseConfig = function (config) {
				config.ehbs = config.ehbs || {};

				var extension = (config.ehbs.extension !== false) ? "." + (config.ehbs.extension || "hbs") : "",
					templatePath = config.ehbs.templatePath ? config.ehbs.templatePath + (config.ehbs.templatePath[config.ehbs.templatePath.length-1] === "/" ? "" : "/") : "",
					ember = config.ehbs.ember || "Ember";

				return {
					extension : extension,
					templatePath : templatePath,
					ember : ember
				};
			};

		return {

			load: function (name, req, load, config) {

				var ehbsConfig = parseConfig(config);

				if (!config.isBuild) {

					req(["text!" + ehbsConfig.templatePath + name + ehbsConfig.extension], function (val) {

						define(module.id + "!" + name, [ehbsConfig.ember], function (Ember) {
							var t = Ember.TEMPLATES[name] = Ember.Handlebars.compile(val);
							return t;
						});

						req([module.id + "!" + name], function (val) {
							load(val);
						});

					});

				}
				else {
					CONFIG = ehbsConfig;
					load("");
				}
			},

			loadFromFileSystem : function (plugin, name) {

				var fs = nodeRequire('fs'),
					file = require.toUrl(CONFIG.templatePath + name) + CONFIG.extension,
					compiler = nodeRequire(CONFIG.etcPath || 'ember-template-compiler'),
					template = compiler.precompile(fs.readFileSync(file, { encoding: 'utf8' })),
					output = "define('" + plugin + "!" + name  + "', ['" + CONFIG.ember + "'], function (Ember) {\nvar t = Ember.TEMPLATES['" + name + "'] = Ember.Handlebars.template(" + template + ");\nreturn t;\n});\n";

				return output;
			},

			write: function (pluginName, moduleName, write) {
				write(this.loadFromFileSystem(pluginName, moduleName));
			}

		};
	}
);