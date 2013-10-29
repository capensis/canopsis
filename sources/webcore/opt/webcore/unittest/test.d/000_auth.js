//#################
//# Auth
//#################

var login              = "root";
var password           = "root";
var selector_form      = 'div[id="auth_form"]';
var selector_submitBtn = 'button[id="submitbutton-btnEl"]';
var selector_viewport  = 'div.widget-wrapper';
var selector_websocket = '#Mainbar-menu-Websocket-btnIconEl[class~="icon-bullet-green"]';

casper.then(function() {
	casper.echo('> Fill Auth Form', 'COMMENT');
	wait(selector_form, timeout, "Form is loaded");
	fill_field(selector_form, 'login', login);
	fill_field(selector_form, 'password', password);
});

casper.then(function() {
	casper.echo('> Submit Form', 'COMMENT');
	click(selector_submitBtn);
});

casper.then(function() {
	casper.echo('> Wait viewport', 'COMMENT');
	wait(selector_viewport, timeout, "Auth is Ok and Application is loaded");
});

casper.then(function() {
	casper.echo('> Wait Websocket', 'COMMENT');
	wait(selector_websocket, timeout, "Websocket is connected");
});