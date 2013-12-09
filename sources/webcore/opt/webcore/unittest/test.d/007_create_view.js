//#################
//# Create View
//#################

casper.then(function() {
	casper.echo('> Open menu and click "new view"', 'COMMENT');
	openMenu('build','buildNewView');
});

casper.then(function() {
	casper.waitUntilVisible('.x-message-box');
});

casper.then(function() {
	click(".x-message-box input");
	casper.sendKeys('.x-message-box input', 'CasperView');
	clickLabel('OK');
});

casper.then(function() {
	casper.waitForSelector('.widget-container-hgrid', function() {
		click('span.jq-tb-save');
		casper.waitUntilVisible("div.ui-pnotify-container");
	});
});

casper.then(function() {
	closeTab();
	openMenu('run','runViewManager','CasperView');
});

casper.then(function() {
	wait('span[name="view.CasperView"]',function() {
		casper.mouseEvent('mouseover', 'span[name="view.CasperView"]');
	});
});

casper.then(function() {
	casper.mouseEvent('mousedown', 'span[name="view.CasperView"]');
});

casper.then(function() {
	click("span.icon-delete");
	clickLabel('Yes');
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.waitForText("CasperView", function() {
		casper.test.fail("record not deleted");
	}, function() {
		casper.test.pass("Ok");
	}, 2000);
});