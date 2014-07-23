define(['consolejs'], function() {

	delete console.init;

	console.debug = console.log;

	return console;
});