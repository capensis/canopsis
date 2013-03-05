//#################
//# Account
//#################

casper.then(function() {
	casper.echo('> Click on Build', 'COMMENT');
	clickMenu("build");
});

casper.then(function() {
	casper.echo('> Click Edit Accounts', 'COMMENT');
	clickMenu("buildAccount");
	waitText("CPS_root");
});

casper.then(function() {
	casper.echo('> Reload Accounts', 'COMMENT');
	click("span.icon-reload");
});

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
	var selector_form = ".x-fieldset";
	fill_field(selector_form, 'user', 'casper');
	fill_field(selector_form, 'passwd', 'casper');
	fill_field(selector_form, 'firstname', 'Casper');
	fill_field(selector_form, 'lastname', 'JS');
	fill_field(selector_form, 'mail', 'capser@js.com');
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