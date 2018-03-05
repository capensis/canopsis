define(["lib/seeds/utils/console"], function(newConsole) {

	proxies = [newConsole];

	var initialisation_done = 0;
	if(initialisation_done === 0) {
		initialisation_done ++;
		console.log("initialisation_done " + initialisation_done);
		var proxiedDefine = define; // Preserving original function

		console.log("define overwrite");
		define = function() {
			for (var i = 0; i < proxies.length; i++) {
				var proxy = proxies[i];

				if(proxy.beforeDefine !== undefined)
					arguments = proxy.beforeDefine.apply(this, arguments);
			}

			if(arguments.length === 2)
			{
				var proxiedCallback = arguments[1];

				arguments[1] = function() {


					for (var i = 0; i < proxies.length; i++) {
						var proxy = proxies[i];

						var args = arguments;

						if(proxy.beforeCallback !== undefined)
							args = proxy.beforeCallback.apply(this, args);
						console.old.log("args");
						console.old.log(args);
					}

					return proxiedCallback.apply(this, arguments);
				};
			}
			return proxiedDefine.apply(this, arguments);
		}
	}
});