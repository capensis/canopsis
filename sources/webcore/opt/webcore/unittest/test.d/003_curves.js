//#################
//# curves
//#################

var selector_curve_form = ".x-window-item";

cgridOpen('build','buildCurve','cps_pct_by_state_0');

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

cgridEditRecord(selector_curve_form,'label', 'CasperLabelModified');

cgridRemoveRecord();

casper.then(function() {
	casper.echo('> Close Tab', 'COMMENT');
	closeTab();
});