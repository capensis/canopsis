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
Ext.define('canopsis.lib.view.cmail' , {
	extend: 'Ext.window.Window',

	alias: 'widget.crights',

	title: 'Email edition',

	bodyHtml: false,

	constrain: true,

	attachement: false,

	logAuthor: '[cmail]',

	initComponent: function() {
		log.debug('Initializing...', this.logAuthor);

		// building bbar
		this.bbar = Ext.create('Ext.toolbar.Toolbar');

		this.cancelButton = this.bbar.add({
			xtype: 'button',
			text: _('Cancel'),
			action: 'cancel',
			iconCls: 'icon-cancel'
		});

		this.bbar.add('->');

		this.finishButton = this.bbar.add({
			xtype: 'button',
			text: _('Finish'),
			action: 'finish',
			iconCls: 'icon-save',
			iconAlign: 'right'
		});

		// building reciepients options

		this.recipientsOptions = Ext.widget('container', {
			xtype: 'fieldset',
			layout: 'hbox',
			frame: true,
			border: 0,
			padding: 0,
			width: 500,
			collapsible: false
		});

		this.to = Ext.widget('textfield', {
			fieldLabel: _('To'),
			width: 275,
			name: 'recipients'
		});

		this.comboUser = Ext.widget('combo', {
			margin: '0 0 0 2',
			queryMode: 'local',
			displayField: 'user',
			valueField: 'mail',
			width: 120,
			store: 'Accounts'
		});

		this.addUserButton = Ext.widget('button', {
			xtype: 'button',
			margin: '0 0 0 2',
			text: _('Add')
		});

		this.recipientsOptions.add([this.to, this.comboUser, this.addUserButton]);

		// mail information
		this.subject = Ext.widget('textfield', {
			fieldLabel: _('subject'),
			width: 400,
			name: 'subject'
		});


		if(this.bodyHtml === true) {
			this.mailbody = Ext.widget('htmleditor', {
				fieldLabel: _('body'),
				name: 'body'
			});
		}
		else {
			this.mailbody = Ext.widget('textareafield', {
				fieldLabel: _('body'),
				width: 400,
				height: 200,
				name: 'body'
			});
		}

		// Building window
		this._form = Ext.create('Ext.form.Panel', {border: false});
		this._form.add([this.recipientsOptions, this.subject, this.mailbody]);
		this.items = this._form;

		this.callParent(arguments);

		//  Binding events
		this.addUserButton.on('click', this._addUser, this);

		this.cancelButton.on('click', function() {
			this.close();
		}, this);

		this.finishButton.on('click', function() {
			var values = this._form.getValues();

			if(this.attachement) {
				values.attachement = this.attachement;
			}

			this.fireEvent('finish', values);
			this.close();
		}, this);
	},

	_addUser: function() {
		log.debug('clicked on adduser', this.logAuthor);
		var recipientsValue = this.to.getValue();

		if(recipientsValue === '') {
			this.to.setValue(this.comboUser.getValue());
		}
		else {
			this.to.setValue(recipientsValue + ',' + this.comboUser.getValue());
		}
	}
});
