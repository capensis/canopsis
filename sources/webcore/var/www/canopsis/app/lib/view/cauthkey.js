/*
#--------------------------------
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
# ---------------------------------
*/
Ext.define('canopsis.lib.view.cauthkey' , {
	extend: 'Ext.window.Window',

	alias: 'widget.crights',

	title: _('Authentification key'),

	constrain: true,

	account: undefined,

	//layout : 'vbox',

	logAuthor: '[cauthkey]',

	initComponent: function() {
		log.debug('Initializing...', this.logAuthor);

		//set title
		if (this.account)
			this.title += ' : ' + this.account;

		/*
		this.helperTemplate = new Ext.Template(
			'<div name="helper" style="{style}">',
				'<h2>You can auto Logging in canopsis with this link</h2>',
				'<div >',
				'<br />',
				_('Bookmark me ! :'),
				'<a href={link}>',
				'{tinyLink}',
				'</a></div>',
			'</div>',
			{compiled: true}
		);
		*/
		//-----------------------Build inner form----------------
		var config = {
			readOnly: true,
			width: 450
		};
		this.authkey_field = Ext.widget('textfield', config);

		var buttonConfig = {
			tooltip: _('Ask for a new key'),
			iconCls: 'icon-reload',
			width: 26
		};
		this.refreshButton = Ext.widget('button', buttonConfig);

		//-------------------------Build form--------------------

		if (global.accountCtrl.checkRoot() || this.checkDisplayButton() )
			var form_width = config.width + buttonConfig.width;
		else
			var form_width = config.width;

		var formConfig = {
			border: false,
			layout: 'hbox',
			width: form_width,
			margin: 3
			//height : 22,
		};
		this._form = Ext.create('Ext.panel.Panel', formConfig);
		this._form.add([this.authkey_field]);
		if (global.account.user == 'root' || this.checkDisplayButton())
			this._form.add(this.refreshButton);

		//------------------------build link helper--------------
		var configHelper = {
			//height : 200,
			margin: '20 0 0 0'
		};
		this.panelHelper = Ext.create('Ext.panel.Panel', formConfig);

		this.items = Ext.create('Ext.panel.Panel', {
				//items:[this._form,this.panelHelper],
				items: [this._form],
				height: 28,
				border: false
			});

		this.callParent(arguments);

		//------------------------set authkey value-------------
		if (this.account)
			this.getAccountKey();
		else
			this.updateTextBox(global.account.authkey);

		//this.updateHelper()
		//-----------------------binding events-------------------
		this.refreshButton.on('click', this._new_authkey, this);
	},


	_new_authkey: function() {
		log.debug('Asking for a new authentification key', this.logAuthor);

		Ext.MessageBox.confirm(_('Confirm'), _('If you generate a new authentification key, the old one will NOT work anymore. Do want to update the key now ?'),
			function(btn, text) {
				if (btn == 'yes')
					if (this.account)
						global.accountCtrl.new_authkey(this.account, this.updateTextBox, this);
					else
						global.accountCtrl.new_authkey(global.account.user, this.updateTextBox, this);
				else
					log.debug('cancel new key generation', this.logAuthor);
			},this);
	},

	updateTextBox: function(text) {
		if (text != undefined) {
			this.authkey_field.setValue(text);
			//this.updateHelper()
		}else {
			global.notify.notify(_('Error'), _('An error have occured during the updating process'), 'error');
		}
	},
    
    checkDisplayButton: function(){
        if(global.accountCtrl.checkGroup('group.CPS_authkey') && !this.account)
            return true
        else
            return false
    },

	updateHelper: function() {
		var url = location.origin + '/' + global.account.authkey;
		this.panelHelper.update(this.helperTemplate.apply({
				style: 'text-align:center;',
				link: location.origin + '/' + global.account.authkey,
				tinyLink: url.slice(0, 40) + '...'
			}));
	},

	getAccountKey: function() {
		log.debug('Get account Authkey', this.logAuthor);
		global.accountCtrl.get_authkey(this.account, this.updateTextBox, this);
	}

});

