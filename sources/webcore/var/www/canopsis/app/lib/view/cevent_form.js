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
Ext.define('canopsis.lib.view.cevent_form' , {
	extend: 'Ext.window.Window',

	alias: 'widget.crights',

	logAuthor: '[cevent_form]',

	layout: 'fit',

	constrain: true,

	resizable: false,

	border: false,

	title: _('Editing rights'),

	initComponent: function() {
		log.debug('Initializing...', this.logAuthor);

		// creating bbar
		this.saveButton = Ext.widget('button', {
			text: _('Save'),
			iconCls: 'icon-save',
			iconAlign: 'right'
		});

		this.cancelButton = Ext.widget('button', {
			text: _('Cancel'),
			iconCls: 'icon-cancel'
		});

		this.bbar = [this.cancelButton, '->', this.saveButton];

		// binding events
		this.saveButton.on('click', function() {
			this._save(this.data);
		}, this);

		this.cancelButton.on('click', function() {
			this.close();
		}, this);

		// create inner form
		var items = [
			{
				xtype: 'textfield',
				fieldLabel: _('Connector'),
				name: 'connector',
				allowBlank: false
			},{
				xtype: 'textfield',
				fieldLabel: _('Connector Name'),
				name: 'connector_name',
				allowBlank: false
			},{
				xtype: 'textfield',
				fieldLabel: _('Event type'),
				name: 'event_type',
				allowBlank: false
			},{
				xtype: 'textfield',
				fieldLabel: _('Source type'),
				name: 'source_type',
				allowBlank: false
			},{
				xtype: 'textfield',
				fieldLabel: _('Component'),
				name: 'component',
				allowBlank: false
			},{
				xtype: 'textfield',
				fieldLabel: _('Resource'),
				name: 'resource',
				allowBlank: false
			},{
				xtype: 'textfield',
				fieldLabel: _('State'),
				name: ('state'),
				allowBlank: false
			},{
				xtype: 'textfield',
				fieldLabel: _('State type'),
				name: 'state_type'
			},{
				xtype: 'textfield',
				fieldLabel: _('Message'),
				name: 'output'
			},{
				xtype: 'textfield',
				fieldLabel: _('Long message'),
				name: 'long_output'
			}
		];

		var config = {
			border: false,
			padding: 2,
			items: items
		};

		this.form = Ext.widget('form', config);

		this.items = [this.form];

		this.callParent(arguments);
	},

	_save: function() {
		if(this.form.form.isValid()) {
			event = this.form.getValues();

			global.eventsCtrl.sendEvent(event);
		}
		else {
			global.notify.notify(_('Invalid form'), _('Please correct the form'), 'error');
		}
	}
});
