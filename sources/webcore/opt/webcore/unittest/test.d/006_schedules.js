//#################
//# schedules
//#################

var selector_consolidation_form = ".x-window-item";

//open and check
cgridOpen('report','buildSchedule')

casper.then(function() {
	casper.echo('> Fill fields and Save', 'COMMENT');
	fill_field(selector_consolidation_form, 'crecord_name', 'Casper');
	selectComboValue('exporting_viewName','Dashboard')
	fill_field(selector_consolidation_form, 'crontab_hours','00:00')

	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

cgridEditRecord(selector_consolidation_form, 'crecord_name', 'CasperMod')

cgridRemoveRecord('CasperMod')

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});