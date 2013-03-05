//#################
//# Groups
//#################

var selector_account_form = ".x-fieldset";

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

/*
casper.then(function() {
	casper.echo('> Check if account exist', 'COMMENT');
	casper.waitForText("Casper", function() {
		casper.test.fail("Account already in store !");
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
	fill_field(selector_account_form, 'user', 'casper');
	fill_field(selector_account_form, 'passwd', 'casper');
	fill_field(selector_account_form, 'firstname', 'Casper');
	fill_field(selector_account_form, 'lastname', 'JS');
	fill_field(selector_account_form, 'mail', 'capser@js.com');
	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Select created Account', 'COMMENT');
	clickRowLabel('Casper');
});

casper.then(function() {
	casper.echo('> Edit Account', 'COMMENT');
	clickRowLabel('Casper', true);
	wait("span.icon-save");
});

casper.then(function() {
	fill_field(selector_account_form, 'lastname', 'Edited');
	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Edited");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Select created Account', 'COMMENT');
	clickRowLabel('Casper');
});

casper.then(function() {
	casper.echo('> Remove created Account', 'COMMENT');
	click("span.icon-delete");
	clickLabel('Yes');
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Reload Account', 'COMMENT');
	click("span.icon-reload");
});

casper.then(function() {
	casper.echo('> Check if account is realy deleted', 'COMMENT');

	casper.waitForText("Casper", function() {
		casper.test.fail("Account not deleted");
	}, function(){
		casper.test.pass("Ok");
	}, 500);
});

*/

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});

