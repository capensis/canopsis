define([
	'ember'
], function(Ember) {

	/**
	 * TODO : Add validators to crecord.attributes
	 * Scan attr's options in order to retrieve all needed validators (Ember.validators)
	 * @return validators : array of validators
	 */
	function GetValidators(attr) {

		console.log("attr = ", attr);
		var options = attr.model.options;
		var validators = [];

		if (options !== undefined) {

			for (var key in options) {
				if ( options.hasOwnProperty( key ) ) {
					var validatorName = (key === "role")? options[key] : key;
					var validator = Ember.validators[validatorName];

					if (validator !== undefined) {

						console.log("pushed : ", validatorName);
						validators.push(Ember.validators[validatorName]);
					}
				}
			}
		}

		return validators;
	}
	/**
	 * Create struct (not really needed)
	 */
	function makeStruct(attributes) {

		var names = attributes.split(' ');
		var count = names.length;
		function constructor() {

			for (var i = 0; i < count; i++) {
				this[names[i]] = arguments[i];
			}
		}
		return constructor;
	}

	/**
	* Check attr's value with all needed validators
	* @return valideStruct : struct containing result of validation(boolean) and message(string) .
	*/
	function Validator(attr) {

		var errorMessage = "";
		var valideStruct = makeStruct("valid error");
		var toReturn = new valideStruct(true, errorMessage);

		var validators = GetValidators(attr);

		for (var i = 0; i < validators.length; i++) {

			validator = validators[i];
			toReturn = validator(attr, toReturn);

			if (toReturn.valid === false) {
				return toReturn;
			}
		}
		return toReturn;
	};

	return Validator;
});
