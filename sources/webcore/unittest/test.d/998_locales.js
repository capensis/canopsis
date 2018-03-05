//#################
//# Locales
//#################

casper.then(function() {
	casper.echo('> Check untranslated strings', 'COMMENT');

	var dump = casper.evaluate(function() {
		return global.dump_untranslated();
	});

	if(dump) {
		console.log(dump);
		casper.test.fail("Some strings are not translated");
	}
	else {
		casper.test.pass("All strings are translated");
	}
});
