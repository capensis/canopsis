//#################
//# Groups
//#################

var selector_group_form = ".x-window-item";

cgridOpen('build','buildGroup','CPS_root');

casper.then(function() {
	casper.echo('> Fill fields and Save', 'COMMENT');
	fill_field(selector_group_form, 'crecord_name', 'Casper');
	fill_field(selector_group_form, 'description', 'CasperDescription');
	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Select created group', 'COMMENT');
	clickLabel('CasperDescription');
	wait(".x-editor");
});

casper.then(function() {
	casper.echo('> Edit group', 'COMMENT');
	fillGridEditableField('.x-editor input', 'Modified');
	click("span.icon-reload");
	waitText("CasperDescriptionModified");
});

cgridRemoveRecord();

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});

