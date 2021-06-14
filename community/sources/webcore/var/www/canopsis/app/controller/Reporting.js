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
Ext.define('canopsis.controller.Reporting', {
	extend: 'Ext.app.Controller',

	views: [],

	logAuthor: '[controller][Reporting]',

	init: function() {
		log.debug('Initialize ...', this.logAuthor);
		this.callParent(arguments);
	},

	launchReport: function(view_id, from, to, mail, orientation, pagesize, timezone) {
		log.debug('Launch Report on view ' + view_id, this.logAuthor);

		//if no date given
		if(!to) {
			to = parseInt(Ext.Date.now()/1000);
		}
		else {
			to = parseInt(to/1000);
		}

		if(!from) {
			from = -1;
		}
		else {
			from = parseInt(from/1000);
		}

		var url = '/reporting/'+ from + '/' + to + '/' + view_id + '/';

		if(mail !== undefined) {
			url += mail + '/';
			if (timezone !== undefined) {
				url += timezone + '/';
			}
		}

		global.notify.notify(_('Please Wait'), _('Your document is rendering, It will be ready in about 1 minute in the briefcase or via email. Please wait...'));

		Ext.Ajax.request({
			url: url,
			scope: this,
			params: {
				'orientation': orientation,
				'pagesize': pagesize
			},
			success: function(response) {
				var data = Ext.JSON.decode(response.responseText);
				log.dump(data);

				if(data.success === true) {
					var id = data.data[0].id;

					global.notify.notify(
						_('Export ready'),
						_('You can download your document') + ' <a href="' + location.protocol + '//' + location.host + '/files/' + id + '"  target="_blank">' + _('here') + '</a>',
						'success'
					);
				}
				else {
					global.notify.notify('Still working', 'The report generation is currently working, please wait a while...', 'success');
					log.info('The report generation is currently working, please wait a while...', this.logAuthor);
				}
			},
			failure: function() {
				global.notify.notify('Still working', 'The report generation is currently working, please wait a while...', 'success');
				log.info('The report generation is currently working, please wait a while...', this.logAuthor);
			}
		});
	},

	downloadReport: function(id) {
		url = location.protocol + '//' + location.host + '/files/' + id;
		window.open(url, '_newtab');
	},

	openHtmlReport: function(view, from, to) {
		log.debug('Open html report : ' + view, this.logAuthor);

		var url = Ext.String.format(
			'http://{0}{1}?exportMode=true&view_id={2}&from={3}&to={4}',
			window.location.host,
			window.location.pathname,
			view,
			parseInt(from / 1000),
			parseInt(to / 1000)
		);

		log.debug('url is : ' + url);
		window.open(url, '_newtab');
	}
});
