//#################
//# Logout
//#################

var selector_logoutBtn =	'button[data-qtip="Logout"]';

casper.then(function() {
	casper.echo('> Logout', 'COMMENT');
	click(selector_logoutBtn);
	wait(selector_form, timeout * 2, "Logout Ok");
});