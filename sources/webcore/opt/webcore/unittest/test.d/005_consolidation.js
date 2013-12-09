//#################
//# consolidations
//#################

var selector_consolidation_form = ".x-window-item";

//open and check
cgridOpen('build','buildConsolidation');

casper.then(function() {
	casper.echo('> Fill fields and Save', 'COMMENT');
	fill_field(selector_consolidation_form, 'crecord_name', 'Casper');
	fill_field(selector_consolidation_form, 'component', 'CasperComponent');
	fill_field(selector_consolidation_form, 'resource', 'CasperResource');

	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

cgridEditRecord(selector_consolidation_form, 'component', 'CasperComponentModified');

cgridRemoveRecord();

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});