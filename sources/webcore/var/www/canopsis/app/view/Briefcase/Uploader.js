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
Ext.define('canopsis.view.Briefcase.Uploader', {
	extend: 'Ext.window.Window',

	alias: 'widget.Uploader',

	title: 'Add new file to briefcase',
	height: 110,
	width: 300,
	layout: 'fit',

	logAuthor: '[view][Uploader]',

	callback: undefined,

	items: [{
		xtype: 'form',
		bodyPadding: '5 5 0',

		defaults: {
			allowBlank: false
		},

		items: [{
			xtype: 'filefield',
			id: 'form-file',
			emptyText: 'Select a file',
			fieldLabel: 'File',
			name: 'file-path',
			buttonText: 'Browse',
			width: 275
		}],

		buttons: [
			{
				text: 'Upload',
				handler: function() {
					var win = this.up('window');
					var form = this.up('form').getForm();

					if(form.isValid()) {
						form.submit({
							url: '/files',
							success: function(form, action) {
								void(form);

								var file_id = action.result.data.file_id;
								var filename = action.result.data.filename;

								log.debug("File '" + filename + "' uploded with id: '" + file_id + "'", win.logAuthor);

								global.notify.notify(_('Success'), _('File uploaded with sucess'), 'success');

								if(win.callback) {
									win.callback(file_id, filename);
								}

								// Reload stores
								Ext.getStore('Files').load();
								Ext.getStore('Avatar').load();

								// Close windows
								win.close();

							},
							failure: function() {
								log.error("Failed to upload file", win.logAuthor);
								global.notify.notify(_('Failed'), _('Failed to upload file'), 'error');
								win.close();
							}
						});
					}
				}
			}
		]
	}]
});