/*
# Copyright (c) 2011 "Capensis" [http://www.capensis.com]
#
# This file is part of Canopsis.
#
# Canopsis is free software: you can redistribute it and/or modify
# it under the terms of the GNU Affero General Public License as published by
# the Free Software Foundation, either version 3 of the License, or
# (at your option) any later version.
#
# Canopsis is distributed in the hope that it will be useful,
# but WITHOUT ANY WARRANTY; without even the implied warranty of
# MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
# GNU Affero General Public License for more details.
#
# You should have received a copy of the GNU Affero General Public License
# along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
*/

function normalApp() {
	var app = Ext.application({
		name: 'canopsis',
		appFolder: 'app',

		controllers: [
			'Notify',
			'View',
			'Websocket',
			'Mainbar',
			'Widgets',
			'Group',
			'Account',
			'Ldap',
			'Reporting',
			'ReportingBar',
			'Keynav',
			'Schedule',
			'Briefcase',
			'Curves',
			'MetricNavigation',
			'Events',
			'Selector',
			'Statemap',
			'Derogation',
			'Perfdata',
			'Topology',
			'Consolidation',
			'Rule',
			'Tabs'
		],

		//autoCreateViewport: true,
		launch: function() {
			if(global.minimified) {
				this.createViewport();
			}
			else {
				this.getController('Widgets').on('loaded', this.createViewport, this, {single : true});
			}
		},

		createViewport: function() {
			Ext.create('canopsis.view.Viewport');
			log.debug('Remove mask ...',"[app]");

			if(Ext.get('loading')) {
				Ext.get('loading').remove();
			}

			if(Ext.get('loading-mask')) {
				Ext.get('loading-mask').remove();
			}
		}
	});

	return app;
}

function exportApp() {
	// parse dates
	var to   = new Date().getTime();

	// 1 day by default
	var from = to - (86400 * 1000);

	if(ENV["to"]) {
		to = ENV["to"] * 1000;
	}

	if(ENV["from"] === -1) {
		from = undefined;
	}
	else {
		if(ENV["from"]) {
			from = ENV["from"] * 1000;
		}
	}

	log.debug('Exporting options:');

	if(from) {
		var from_date = new Date(from);
		log.debug(' + From: ' + from_date + '(' + from + ')');
	}

	var to_date = new Date(to);
	log.debug(' + To:   ' + to_date + '(' + to + ')');

	var title = '';

	if(from && from !== to) {
		title =  Ext.String.format(
			'<span>' + _('From') + ' {0} ' + _('To') + ' {1}</span>',
			Ext.Date.format(from_date, 'Y-m-d'),
			Ext.Date.format(to_date, 'Y-m-d')
		);
	}
	else {
		title = Ext.String.format(
			'<span>{0}</span>',
			Ext.Date.format(to_date, 'Y-m-d')
		);
	}

	title = '<div id="interval_header" name="interval_header" style="margin-left:400px;">' + title + '</div>';

	$("body").append(title);
	$("body").append("<div id='container'></div>");

	var app = Ext.application({
		name: 'canopsis',
		appFolder: 'app',
		controllers: [
			'Account',
			'Widgets',
			'Curves',
			'Websocket'
		],
		launch: function() {
			this.getController('Widgets').on('loaded', this.createView);
		},
		createView: function() {
			var content = Ext.create('canopsis.view.Tabs.Content', {
				renderTo: 'container',
				view_id: ENV['view_id'],
				autoshow: true,
				reportMode: true,
				exportMode: true,
				export_from: from,
				export_to: to
			});

			// Hack fix manual height
			content.on('loaded', function() {
				var options = $('#' + content.container.id).jqGridable('get_options');
				log.dump(options);

				var jqg_height = options.rows * options.widget_height;
				var header_height = $('#interval_header').height();

				log.debug('[Reporting] jqg_height: ' + jqg_height);
				log.debug('[Reporting] header_height: ' + header_height);

				var total_height = jqg_height + header_height;

				$('#' + content.container.id).height(total_height);
				$('#' + content.id).height(total_height);
				$('body').height(total_height);

				log.debug('[Reporting] Body height : ' + $('body').height());
			});

			// aware wkhtml that loading is finished
			var task = new Ext.util.DelayedTask(function() {
				window.status = 'ready';
			});

			task.delay(10000);
		}
	});

	return app;
}

function fullscreenApp() {
	$("body").append("<div id='container'></div>");

	var app = Ext.application({
		name: 'canopsis',
		appFolder: 'app',
		controllers: [
			'Account',
			'Widgets',
			'Curves',
			'Websocket',
			'Notify'
		],
		launch: function() {
			this.getController('Widgets').on('loaded', this.createView);
		},
		createView: function() {
			var content = Ext.create('canopsis.view.Tabs.Content', {
				renderTo: 'container',
				width: ($('#container').width()- 1),
				view_id : ENV['view_id'],
				autoshow : true,
				fullscreenMode: true
			});
		}
	});

	return app;
}

function checkLocale() {
	log.debug("Loading locale:", "[app]");

	log.debug(" + ENV:     " + ENV['locale'], "[app]");
	global.locale = ENV['locale'];

	var cookie_locale = Ext.util.Cookies.get('locale');
	log.debug(" + Cookie:  "+ cookie_locale, "[app]");

	if(global.account && global.account['locale']) {
		var account_locale = global.account['locale'];
		log.debug(" + Account: "+ account_locale, "[app]");
	}

	// For export mode
	if(ENV["exportMode"] && global.locale !== account_locale) {
		Ext.fly('locale').set({src:'/' + account_locale + '/static/canopsis/locales.js'});
	}

	log.debug("Locale:     " + global.locale, "[app]");
}

function createApplication(account) {
	void(account);

	log.debug("Remove auth form ...", "[app]");

	if(Ext.get('auth')) {
		Ext.get('auth').remove();
	}

	if(global.account.id === undefined) {
		global.account.id = global.account._id;
	}

	// URL Decode
	var url_options = Ext.Object.fromQueryString(window.location.search);
	log.debug(" + Url options: ", "[app]");
	log.dump(url_options);

	// Set env
	ENV["fullscreenMode"] = url_options['fullscreenMode'] === 'true' ? true : false;
	ENV["exportMode"] = url_options['exportMode'] === 'true' ? true : false;
	ENV["reportMode"] = url_options['reportMode'] === 'true' ? true : false;

	checkLocale();

	if(url_options['from']) {
		try {
			ENV["from"] = parseInt(url_options['from']);
		}
		catch(err) {
			log.error("Impossible to parse: " + url_options['from'], "[app]");
		}
	}

	if(url_options['to']) {
		try {
			ENV["to"] = parseInt(url_options['to']);
		}
		catch(err) {
			log.error("Impossible to parse: " + url_options['to'], "[app]");
		}
	}

	if(url_options['view_id']) {
		ENV["view_id"] = url_options['view_id'];
	}

	log.debug(" + ENV: ", "[app]");
	log.dump(ENV);

	// Answer to every error
	Ext.Ajax.on('requestexception', function(conn, response, options) {
		void(conn, options);

		log.error('requestexception: ' + response.status + ': ' + response.statusText, "[app]");

		if(response.status === 403) {
			global.notify.notify(_('Server notification'), _('You have no sufficient rights'), 'info');
		}

		if(response.status === 500) {
			global.notify.notify(_('Server error'),_('Unexpected server error'),'error');
		}

		if(response.status === 400) {
			global.notify.notify(_('Server error'), _(response.statusText), 'error');
		}
	});

	log.debug("Start ExtJS application ...", "[app]");

	var app = undefined;

	// Route
	if(ENV["exportMode"] || ENV["reportMode"]) {
		app = exportApp();
	}
	else if(ENV["fullscreenMode"]) {
		app = fullscreenApp();
	}
	else {
		app = normalApp();
	}

	log.debug("Application started", "[app]");

	return app;
}

