//#################
//# selectors
//#################

var selector_selector_form = ".x-window-item";

//open and check
cgridOpen('build','buildSelector');

casper.then(function() {
	casper.echo('> Fill fields and Save', 'COMMENT');
	fill_field(selector_selector_form, 'crecord_name', 'Casper');
	fill_field(selector_selector_form, 'display_name', 'DisplayName');

	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

cgridEditRecord(selector_selector_form, 'display_name', 'CasperModified');

cgridRemoveRecord();

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});