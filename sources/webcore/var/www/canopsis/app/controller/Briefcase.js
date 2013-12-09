//need:app/lib/controller/cgrid.js,app/lib/view/cmail.js,app/view/Briefcase/Grid.js,app/view/Briefcase/Form.js,app/store/Files.js
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
Ext.define('canopsis.controller.Briefcase', {
	extend: 'canopsis.lib.controller.cgrid',

	requires: [
		'canopsis.lib.view.cmail'
	],

	views: ['Briefcase.Grid', 'Briefcase.Form'],
	stores: ['Files'],
	models: ['File'],

	logAuthor: '[controller][Briefcase]',

	init: function() {
		this.listXtype = 'BriefcaseGrid';
		this.formXtype = 'BriefcaseForm';

		this.modelId = 'File';

		this.callParent(arguments);

	},

	_viewElement: function(view, item) {
		void(view);

		log.debug('Clicked on element, function viewElement', this.logAuthor);
		this.download(item.get('_id'));
	},

	_downloadButton: function() {
		log.debug('clicked deleteButton', this.logAuthor);
		var grid = this.grid;
		var selection = grid.getSelectionModel().getSelection()[0];

		if(selection) {
			this.download(selection.get('_id'));
		}
	},

	download: function(id) {
		url = location.protocol + '//' + location.host + '/files/' + id;
		window.open(url, '_newtab');
	},

	sendByMail: function(record) {
		var config = {
			renderTo: this.grid.id,
			constrain: true,
			attachement: record.get('_id')
		};

		var cmail = Ext.create('canopsis.lib.view.cmail', config);

		cmail.on('finish', function(mail) {
			this._ajaxRequest(mail);
		}, this);

		cmail.show();
	},

	rename: function(item) {
		log.debug(item);
		this._editRecord(this.grid, item);
	},

	_ajaxRequest: function(mail) {
		Ext.Ajax.request({
			type: 'rest',
			url: '/sendreport',
			method: 'POST',
			params: {
				'_id': mail.attachement,
				'recipients': mail.recipients,
				'subject': mail.subject,
				'body': mail.body
				},
			reader: {
				type: 'json',
				root: 'data',
				totalProperty: 'total',
				successProperty: 'success'
			},
			success: function(response) {
				request = Ext.JSON.decode(response.responseText);

				if(request.success) {
					log.debug('Mail have been sent successfuly', this.logAuthor);
					log.info('The server has returned : ' + request.data.output.output);

					if(request.data.output.success === true) {
						global.notify.notify(
							_('Mail sent'),
							_('The mail have been successfuly sent'),
							'success'
						);
					}
				}
				else {
					log.error('Mail have not been sent', this.logAuthor);
					global.notify.notify(_('Failed'), _('Mail not sent, error with celery task'), 'error');
				}
			},
			failure: function() {
				log.debug('Mail request have failed', this.logAuthor);
				global.notify.notify(_('Failed'), _('Mail not sent, webserver error'), 'error');
			}
		});
	}
});
