//#################
//# Groups
//#################

var selector_group_form = ".x-window-item";

casper.then(function() {
	casper.echo('> Click on Build', 'COMMENT');
	clickMenu("build");
});

casper.then(function() {
	casper.echo('> Click Edit Groups', 'COMMENT');
	clickMenu("buildGroup");
	waitText("CPS_root");
});

casper.then(function() {
	casper.echo('> Reload Groups', 'COMMENT');
	click("span.icon-reload");
});

casper.then(function() {
	casper.echo('> Check if group exist', 'COMMENT');
	casper.waitForText("Casper", function() {
		casper.test.fail("Group already in store !");
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
	fill_field(selector_group_form, 'crecord_name', 'Casper');
	fill_field(selector_group_form, 'description', 'CasperDescription');
	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Select created group', 'COMMENT');
	clickLabel('CasperDescription')
	wait(".x-editor");
});

casper.then(function() {
	casper.echo('> Edit group', 'COMMENT');
	fillGridEditableField('.x-editor input', 'Modified');
	click("span.icon-reload");
	waitText("CasperDescriptionModified")
});


casper.then(function() {
	casper.echo('> Select created Group', 'COMMENT');
	clickRowLabel('Casper');
});

casper.then(function() {
	casper.echo('> Remove created Group', 'COMMENT');
	click("span.icon-delete");
});

casper.then(function() {
	clickLabel('Yes');
	waitWhile("div.ui-pnotify-container");
})

casper.then(function() {
	casper.echo('> Reload Group', 'COMMENT');
	click("span.icon-reload");
	waitWhile("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Check if group is really deleted', 'COMMENT');
	casper.waitForText("Casper", function() {
		casper.test.fail("Group not deleted");
	}, function(){
		casper.test.pass("Ok");
	}, 500);
});

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});

