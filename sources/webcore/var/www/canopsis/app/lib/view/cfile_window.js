//need:app/lib/view/cpopup.js
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

Ext.define('canopsis.lib.view.cfile_window' , {
	extend: 'canopsis.lib.view.cpopup',
	alias: 'widget.cfile_window',

	logAuthor: '[cfile_window]',

	title: _('Select a file'),

	width: 400,
	contrain: true,

	_name: _('file'),
	_fieldLabel: _('File'),
	_textAreaLabel: _('Or copy/paste'),
	_buttonText: _('Select file'),

	copyPasteZone: false,

	_buildForm: function() {
		this._fileField = this._form.add(Ext.create('Ext.form.field.File', {
			name: this._name,
			fieldLabel: this._fieldLabel,
			labelWidth: 80,
			msgTarget: 'side',
			anchor: '100%',
			buttonText: this._buttonText
		}));

		this._textArea = this._form.add(Ext.widget('textarea', {
			name: this._name,
			fieldLabel: this._textAreaLabel,
			labelWidth: 80,
			height: 200,
			anchor: '100%'
		}));

		return this._form;
	},

	ok_button_function: function() {
		log.debug('clicked on ok button', this.logAuthor);
		var fileList = this._fileField.fileInputEl.dom.files;
		var text = this._textArea.getValue();

		if(fileList.length > 0) {
			this.fireEvent('save', {
				file: true,
				value: fileList
			});
		}
		else {
			if(text !== undefined) {
				try {
					this.fireEvent('save', {
						file: false,
						value: Ext.decode(text)
					});
				}
				catch (err) {
					log.dump(err);
					global.notify.notify(_('Wrong json'), _('This json is malformed'), 'info');
				}
			}
		}
	}
});
