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

Ext.define('canopsis.lib.view.cpopup' , {
	extend: 'Ext.window.Window',
	alias: 'widget.cpopup',

	logAuthor: '[cpopup]',


	constrain: true,
	resizable: false,

	textAreaLabel: 'Type your message here:',

	initComponent: function() {
		this.logAuthor = '[' + this.id + ']';

		log.debug('Initialize cpopup', this.logAuthor);

		var form = this.buildForm();
		var bar = this.buildBar();

		this.dockedItems = [bar];
		this.items = [form];

		this.callParent(arguments);
	},

	show: function() {
		this.callParent(arguments);

		var elem = this.down('.textfield');

		if(!elem) {
			elem = this.down('.textarea');
		}
		else if (elem.focus) {
			elem.focus(false, 200);
		}
	},

	buildForm: function() {
		this._form = Ext.create('Ext.form.Panel', {
			layout: 'anchor',
			bodyStyle: {
				background: '#ededed'
			},
			bodyPadding: 10,
			border: false
		});

		if(this._buildForm) {
			this._buildForm();
		}

		return this._form;
	},

	buildBar: function() {
		var button_ok = Ext.create('Ext.button.Button', {
			xtype: 'button',
			handler: this.ok_button_function,
			scope: this,
			text: _('Ok'),
			minWidth: 75
		});

		var button_cancel = Ext.create('Ext.button.Button', {
			xtype: 'button',
			handler: this.cancel_button_function,
			scope: this,
			text: _('Cancel'),
			minWidth: 75
		});

		var bar = new Ext.toolbar.Toolbar({
			ui: 'footer',
			dock: 'bottom',
			layout: {
				pack: 'end'
			},
			items: [button_cancel, button_ok]
		});

		if(this._buildBar) {
			this.buildBar(bar);
		}

		return bar;
	},

	ok_button_function: function() {
		log.debug('clicked on ok button', this.logAuthor);

		if(this._ok_button_function) {
			this._ok_button_function();
		}
	},

	cancel_button_function: function(){
		this.close();
	}
});
