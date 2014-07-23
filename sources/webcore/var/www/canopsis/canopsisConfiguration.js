define([], function () {

	/*
	* Here is the canopsis UI main configuration file.
	* It is possible to add properies and values that are reachable
	* from the whole application throught the namespace Canopsis.conf.PROPERTY
	*/
	var canopsisConfiguration = {
		DEBUG: true,
		VERBOSE: 1,
		DISPLAY_SCHAMA_MANAGER: true,
		REFRESH_ALL_WIDGETS: false,
		TRANSLATE: false
	};

	return canopsisConfiguration;
});