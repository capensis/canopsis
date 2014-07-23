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

//add your custom schemas files here.
var schemaFiles = [
];

var schemasDeps = ['ember-data', 'app/application', 'utils'];

for (var i = 0; i < schemaFiles.length; i++) {
	schemasDeps.push('text!schemas/' + schemaFiles[i] + '.json');
}

define(schemasDeps, function(DS, Application, utils) {

	//TODO nuke this
	if (Application.allModels === undefined) {
		Application.allModels = {};
	}

	/**
	 * Loop over localStorage's schemas
	 */
	function loadSchemasFromLocalStorage() {
		console.group("loadSchemasFromLocalStorage", arguments);

		if (typeof(Storage)==="undefined") {
			console.warn("The browser is not supporting localStorage, don't try to load models from localStorage");
		} else {
			// Retrieve the models from the LS
			var localStorageSchemas = localStorage.getItem("canopsis.schemas");

			if (localStorageSchemas !== undefined) {
				//create Ember models from json
				console.log("found schemas in localStorage");
				console.log(localStorageSchemas);
				localStorageSchemas = JSON.parse(localStorageSchemas);

				for (var key in localStorageSchemas) {
					var schema = localStorageSchemas[key];
					var schemaName = key;

					var schemaInheritance = schemaName.split(".");

					addSchema(schemaInheritance, schemaName, schema);
				}
			}
		}
		console.groupEnd();
	}

	/**
	 * Loop over json schemas stored as files to load them
	 * @param {json document} schemasDepsLength
	 * @param {string} schemaFiles
	 * @param {string} moduleArgs the list of files required by this module
	 */
	function loadSchemasFromJsonFiles(schemasDepsLength, schemaFiles, moduleArgs) {
		console.group("loadSchemasFromJsonFiles", arguments);

		for (var i = schemasDepsLength; i < moduleArgs.length; i++) {
			var schemaIndex =  i - schemasDepsLength;
			console.groupCollapsed("loading schema", schemaFiles[schemaIndex]);

			var schema = JSON.parse(moduleArgs[i]);
			var schemaName = schemaFiles[schemaIndex].capitalize();
			var schemaInheritance = schemaFiles[schemaIndex].split(".");

			addSchema(schemaInheritance, schemaName, schema);

		}

		console.groupEnd();
	}


	/**
	 * Loop over json schemas stored as files to load them
	 * @param {api_result document} contains schemas and meta for schemas
	 */
	function loadSchemasFromApiJson(api_result) {
		console.group("loadSchemasFromApi", arguments);

		//Object that contains both shema names as key and information about inheritance solver
		var schemasDict = {};
		var schemaName;

		for (var i=0; i<api_result.length; i++) {
			var schema = api_result[i].schema;
			schemaName = api_result[i].id.capitalize();
			console.log('Loading schema...', schemaName);

			schemasDict[schemaName] = {schema: schema, solved: false};
		}

		for (schemaName in schemasDict) {
			solveDependancy(schemaName, schemasDict);
		}


		console.groupEnd("loadSchemasFromApi", arguments);
	}


	/**
	*	Allow solving dependancies if any
	*/
	function solveDependancy (currentSchemaName, schemasDict) {

		console.log("solveDependancy", currentSchemaName, schemasDict);
		if(schemasDict[currentSchemaName] === undefined) {
			throw (currentSchemaName + ' not referenced in schemas dict');
		}

		var currentSchema = schemasDict[currentSchemaName].schema;

		if (currentSchema.properties !== undefined) {
			for (var propertyName in currentSchema.properties) {
				var property = currentSchema.properties[propertyName];

				if (property.relationship !== undefined && property.model !== undefined) {
					console.log(currentSchemaName, "has dependancy", property.model, "in its relationships");
					solveDependancy(property.model.capitalize(), schemasDict);
				}
			}
		}

		var schemaInheritanceStringList = currentSchemaName.split(".");

		//if inheritance then try to solve
		if (schemaInheritanceStringList.length > 1) {

			//cut the head
			var pop = schemaInheritanceStringList.pop();
			var parentName = schemaInheritanceStringList.join('.');
			if (schemasDict[parentName] && !schemasDict[parentName].solved) {
				solveDependancy(parentName, schemasDict);
			}
			//if parent exists in schema objects
			if (schemasDict[parentName] && schemasDict[parentName].solved) {
				//Try to solve parent case
				schemaInheritanceStringList.push(pop);
				addSchema(schemaInheritanceStringList, currentSchemaName, currentSchema);
				schemasDict[currentSchemaName].solved = true;
			}
		} else if (schemasDict[currentSchemaName]) {
			addSchema(schemaInheritanceStringList, currentSchemaName, currentSchema);
			schemasDict[currentSchemaName].solved = true;
		}
	}

	/**
	 * Build an EmberJS model from a json schema
	 * @param {json document} schema
	 * @param {string} schemaName
	 * @param {DS.Model} parentModelClass the parent model
	 * @param {string} parentModelClassName name of the parent model
	 */
	function generateEmberModelFromSchema(schema, schemaName, parentModelClass, parentModelClassName) {
		console.group("generate model", schemaName, schema);

		var modelDict = { "categories": schema.categories };
		//TODO check if relationship options are ok
		for (var name in schema.properties) {

			var property = schema.properties[name];
			var propertyType = property.type;
			delete property.type;

			console.group('model\'s attributes and relationships');

			if (property.relationship === undefined) {
				//The property is not a relation
				modelDict[name] = DS.attr(propertyType, property);
			} else if (property.model === undefined) {
				throw "property is relationship but no model defined";
			} else {
				var model = property.model;
				model = model.split('.');
				model = model[model.length - 1];

				if (property.relationship === "belongsTo" && model !== undefined) {
					console.log("creating belongsTo with", property, model);
					modelDict[name] = DS.belongsTo(model, property);
				}
				if (property.relationship === "hasMany" && model !== undefined) {
					console.log("creating hasMany with", property, model);
					modelDict[name] = DS.hasMany(model, property);
				}
			}

			console.groupEnd();
		}

		console.log(schemaName, "inherits from",parentModelClass);

		modelDict = inheritance(modelDict, parentModelClassName, schemaName);

		Application.allModels[schemaName] = modelDict;
		var newModel = parentModelClass.extend(modelDict);

		console.log("newModel store", newModel.store);
		console.groupEnd();

		return newModel;
	}

	/**
	 * Add attribute of a model in another one
	 * @param {} modelDict
	 * @param {string} parentModelClassName name of  the parent model
	 * @param {string}  schemaName of the schema
	 */
	function inheritance(modelDict, parentModelClassName, schemaName) {

		console.group('inherited attributes and relationships');

		for (var keys in Application.allModels[parentModelClassName]) {
			if (Application.allModels[parentModelClassName].hasOwnProperty(keys))
			{
				if (!modelDict.hasOwnProperty(keys))
				{
					var val = Application.allModels[parentModelClassName][keys]._meta.options;

					if (val.relationship === "hasMany" && val.model !== undefined) {
						console.log("Add hasMany relationship : " + keys +' = ', val, Application.allModels[parentModelClassName][keys]._meta.type);
						modelDict[keys] = DS.hasMany(val.model, val);
					} else if (val.relationship === "belongsTo" && val.model !== undefined) {
						console.log("Add belongsTo relationship : " + keys +' = ', val, Application.allModels[parentModelClassName][keys]._meta.type);
						modelDict[keys] = DS.belongsTo(val.model, val);
					} else {
						console.log("Add attribute : " + keys +' = ', val, Application.allModels[parentModelClassName][keys]._meta.type);
						modelDict[keys] = DS.attr(Application.allModels[parentModelClassName][keys]._meta.type, val);
					}

				} else if (modelDict[keys] !== undefined) {

					var oldkeys = Application.allModels[parentModelClassName][keys];
					var newkeys = modelDict[keys];

					if (oldkeys !== undefined) {

						var oldkeysAttribute = oldkeys._meta;
						var newkeysAttribute = newkeys._meta;

						var oldOptions = oldkeysAttribute.options;
						var newOptions = newkeysAttribute.options;

						newOptions = merge(oldOptions, newOptions, schemaName);
						modelDict[keys] = DS.attr(newkeysAttribute.type, newOptions);
					}
				}
			}
		}

		console.groupEnd();

		return modelDict;
 	}

	/**
	*	Processes loaded schemas and adds them to the application scope
	*/
	function addSchema(schemaInheritance, schemaName, schema) {
		console.log("addSchema", schemaName, schemaInheritance, schema);

		var parentModelClass;
		var parentModelClassName = "";

		if (utils.schemaList === undefined) {
			utils.schemaList = {};
		}
		utils.schemaList[schemaName] = schema;

		//retreive the good model class the new model should inherit from
		if (schemaInheritance.length > 1) {
			parentModelClassName = schemaInheritance[schemaInheritance.length - 2].capitalize();
			parentModelClass = Application[parentModelClassName];

			schemaName = schemaInheritance[schemaInheritance.length - 1].capitalize();
		}
		else {
			parentModelClass = DS.Model;
		}

		//does schema already exists ?
		if (Application[schemaName]) {
			console.warn('schema already registered', schemaName);
		}

		//overrided by default following main thread rules (schema load order)
		console.log(
			'Adding schema', {
				schemaName: schemaName,
				parentModelClass: parentModelClass,
				parentModelClassName: parentModelClassName
			}, schema);
		Application[schemaName] = generateEmberModelFromSchema(
			schema,
			schemaName,
			parentModelClass,
			parentModelClassName);


		Application.available_types.push(schemaName);
		console.groupEnd();
	}


	/**
	 * Merge options
	 * @param {string} OldOptions Parent model's options
	 * @param {string} NewOptions Current model's options
	 */
	function merge(oldOptions, newOptions, schemaName) {

		for (var options in oldOptions) {
			//if option isn't in current model's options
			if (oldOptions.hasOwnProperty(options) && newOptions[options] === undefined) {
			newOptions[options] = oldOptions[options];
			console.log ("Added "+ options + " = " + oldOptions[options] + " in "+ schemaName);
			}
		}
		return newOptions ;
	}

	//Module's main thread: init schemas from many sources
	Application.available_types = [];

	$.ajax({
		url: '/rest/schemas',
		data: {limit: 1000},
		success: function(data) {
			if (data.success) {
				console.log('Api schema data',data);
				loadSchemasFromApiJson(data.data);
			} else {
				console.error('Unable to load schemas from API');
			}
		},
		async: false
	});


});