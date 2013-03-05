//#################
//# Configs
//#################
var url = 					'http://127.0.0.1:8082/en/static/canopsis/index.debug.html';
///var url = 					'http://demo-devel.canopsis.org/';
var timeout =				5000;

var casper_verbose =		false;
var casper_logLevel =		'debug';
var viewportSize = 			{width: 1366, height: 768};
var capture_interval = 		1000;

//#################
//# Casper
//#################
var fs = require('fs');

var casper = require('casper').create({
    verbose: casper_verbose,
    logLevel: casper_logLevel,
    viewportSize: viewportSize,
    onStepComplete: function(){
    	capture();
    }
});

casper.on('remote.message', function(msg){
	if (casper_verbose)
		console.log(msg);
});


//#################
//# Utils
//#################

var step = 0;
var capture = function(){
	casper.capture('captures/step-'+step+'.png');
	step += 1;
}

var capturer = setInterval(function(){
	capture();
}, capture_interval);

var fill_field = function(selector, name, value){
	var options = {};
	options[name] = value;

	casper.test.assertExists('input[name="'+name+'"]', 'Check if field "'+name+'" exist');
	casper.fill(selector, options, false);
	casper.test.assertField(name, value);
}

var click = function(selector){
	casper.waitForSelector(selector, function() {
		
		casper.test.assertExists(selector, 'Check if '+selector+' exist');

		//casper.mouseEvent('mousemove', selector);
		//casper.mouseEvent('mouseover', selector);
		//casper.mouseEvent('mousedown', selector);

		casper.thenClick(selector, function(){
			casper.test.pass(" + Clicked");
		});

	},function() {
		casper.test.fail("'"+selector+"' not found.");

	}, timeout);
}

var clickRowLabel = function(label){
	var selector = {
		type: 'xpath',
		path: '//*[text()="'+label+'"]'
	}

	casper.waitForText(label, function() {
		casper.mouseEvent('mouseover', selector);
		casper.mouseEvent('mousedown', selector);
		wait("tr.x-grid-row-selected");
	}, timeout)
}

var clickLabel = function(label){
	casper.waitForText(label, function() {
		casper.clickLabel(label);
	}, timeout)
}

var clickMenu = function(name){
	var menu = {
		build:  "span.icon-mainbar-build",
		run:    "span.icon-mainbar-run", 
		report: "span.icon-mainbar-report",
		buildAccount: "img.icon-mainbar-edit-account",
	}
	click(menu[name]);
}

var wait = function(selector, timeout, str_onSuccess, str_onFailed){
	if (str_onSuccess == undefined)  str_onSuccess	= selector + " found";
	if (str_onFailed == undefined)	str_onFailed	= "Impossible to find "+selector;

	casper.waitForSelector(selector, function() {
		casper.test.pass(str_onSuccess);
	}, function() {
		casper.test.fail(str_onFailed);
	}, timeout);
}

var waitWhile = function(selector, timeout, str_onSuccess, str_onFailed){
	if (str_onSuccess == undefined)  str_onSuccess	= selector + " is not found";
	if (str_onFailed == undefined)	str_onFailed	= selector + " found";

	casper.waitWhileSelector(selector, function() {
		casper.test.pass(str_onSuccess);
	}, function() {
		casper.test.fail(str_onFailed);
	}, timeout);
}

var waitText = function(text){
	casper.waitForText(text, function() {
		casper.test.pass("'"+text+"' found.");
	}, function(){
		casper.test.fail("'"+text+"' not found.");
	}, timeout)
}

//#################
//# Load page
//#################
casper.echo('> Load page "'+url+'"', 'COMMENT');
casper.start(url, function() {
	casper.test.assertTitle('Canopsis', 'Check page title');

	casper.evaluate(function() {
		$('body').append('<div id="div-click" style="z-index: 99999999999; position:absolute; overflow:hidden; border-radius:10px; width:10px; height:10px; background-color: #FF0000;"></div>');
		$('#div-click').hide();

		$('body').click(function(event) {
			if (global.divClickTimeout)
				clearTimeout(global.divClickTimeout);

			$('#div-click').css("top", event.pageY);
			$('#div-click').css("left", event.pageX);
			$('#div-click').show();

			global.divClickTimeout = setTimeout(function(){
				$('#div-click').hide();
			}, 300);
		});
	});
});

//#################
//# Exit casper at end
//#################
casper.run(function() {
	clearInterval(capturer);
	capture();

	casper.echo('\n########### END ###########', 'COMMENT');
	casper.echo('> Quit Casper', 'COMMENT');
	//casper.test.done()
	casper.exit();
});


//#################
//# All tests
//#################

