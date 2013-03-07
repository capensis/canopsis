//#################
//# selectors
//#################

var selector_selector_form = ".x-window-item";

casper.then(function() {
	casper.echo('> Open menu selector', 'COMMENT');
	openMenu('build','buildSelector')
});

casper.then(function() {
	casper.echo('> Reload selectors', 'COMMENT');
	click("span.icon-reload");
});

casper.then(function() {
	casper.echo('> Check if selector exist', 'COMMENT');
	casper.waitForText("Casper", function() {
		casper.test.fail("selector already in store !");
	}, function(){
		casper.test.pass("Ok");
	}, 500);
});

casper.then(function() {
	casper.echo('> Open Add form', 'COMMENT');
	click("span.icon-add");
	wait("span.icon-save");
});

casper.then(function() {
	casper.echo('> Fill fields and Save', 'COMMENT');
	fill_field(selector_selector_form, 'crecord_name', 'Casper');
	fill_field(selector_selector_form, 'display_name', 'CasperDisplayName');

	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Select created selector', 'COMMENT');
	clickRowLabel('Casper');
});

casper.then(function() {
	casper.echo('> Edit selector', 'COMMENT');
	clickRowLabel('Casper', true);
	wait("span.icon-save");
});

casper.then(function() {
	fill_field(selector_selector_form, 'display_name', 'CasperMod');
	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("CasperMod");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Select created selector', 'COMMENT');
	clickRowLabel('Casper');
});

casper.then(function() {
	casper.echo('> Remove created selector', 'COMMENT');
	click("span.icon-delete");
	clickLabel('Yes');
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Reload selector', 'COMMENT');
	click("span.icon-reload");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Check if selector is really deleted', 'COMMENT');

	casper.waitForText("Casper", function() {
		casper.test.fail("selector not deleted");
	}, function(){
		casper.test.pass("Ok");
	}, 500);

});

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});