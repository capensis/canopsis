//#################
//# schedules
//#################

var selector_consolidation_form = ".x-window-item";

//open and check
cgridOpen('report','buildSchedule');

casper.then(function() {
	casper.echo('> Fill fields and Save', 'COMMENT');
	fill_field(selector_consolidation_form, 'crecord_name', 'Casper');
	fill_field(selector_consolidation_form, 'crontab_hours','00:00');
});

casper.then(function() {
	selectComboValue('exporting_viewName','Dashboard');
});

casper.then(function() {
	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

cgridEditRecord(selector_consolidation_form, 'crecord_name', 'CasperModified');

cgridRemoveRecord('CasperModified');

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});