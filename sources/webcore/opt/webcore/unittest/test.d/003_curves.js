//#################
//# curves
//#################

var selector_curve_form = ".x-window-item";

casper.then(function() {
	casper.echo('> Open menu curve', 'COMMENT');
	openMenu('build','buildCurve','cps_pct_by_state_0')
});

casper.then(function() {
	casper.echo('> Reload curves', 'COMMENT');
	click("span.icon-reload");
});

casper.then(function() {
	casper.echo('> Check if curve exist', 'COMMENT');
	casper.waitForText("Casper", function() {
		casper.test.fail("Curve already in store !");
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
	fill_field(selector_curve_form, 'metric', 'Casper');
	fill_field(selector_curve_form, 'label', 'CasperLabel');
	selectComboValue('dashStyle','ShortDashDotDot');
	fill_field(selector_curve_form, 'line_color', 'C79F4B');
	fill_field(selector_curve_form, 'area_color', 'C79F4B');
	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("Casper");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Select created curve', 'COMMENT');
	clickRowLabel('CasperLabel');
});

casper.then(function() {
	casper.echo('> Edit curve', 'COMMENT');
	clickRowLabel('CasperLabel', true);
	wait("span.icon-save");
});

casper.then(function() {
	fill_field(selector_curve_form, 'label', 'CasperLabelModified');
	click("span.icon-save");
	waitWhile("span.icon-save");
	waitText("CasperLabelModified");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Select created curve', 'COMMENT');
	clickRowLabel('CasperLabelModified');
});

casper.then(function() {
	casper.echo('> Remove created Curve', 'COMMENT');
	click("span.icon-delete");
	clickLabel('Yes');
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Reload curve', 'COMMENT');
	click("span.icon-reload");
	casper.waitUntilVisible("div.ui-pnotify-container");
});

casper.then(function() {
	casper.echo('> Check if curve is really deleted', 'COMMENT');

	casper.waitForText("Casper", function() {
		casper.test.fail("curve not deleted");
	}, function(){
		casper.test.pass("Ok");
	}, 500);

});

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});