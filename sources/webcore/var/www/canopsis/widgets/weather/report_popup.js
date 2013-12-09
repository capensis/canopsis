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

Ext.define('widgets.weather.report_popup', {
	extend: 'canopsis.lib.view.cpopup',
	alias: 'widget.weather.report_popup',

	_component: undefined,
	referer: undefined,
	width: 300,
	helpdesk: undefined,
	display_name: undefined,

	base_event: {
		'connector_name': 'widget-weather',
		'connector': 'canopsis',
		'event_type': 'log',
		'source_type': 'resource',
		'component': undefined,
		'resource': 'user_problem',
		'referer': undefined,
		'author': undefined,
		'state': 1,
		'display_name': this.display_name,
		'state_type': 1,
		'output': ''
	},


	initComponent: function() {
		this.base_event['author'] = global.account.firstname + ' ' + global.account.lastname;
		this.callParent(arguments);
	},

	_buildForm: function() {
		this._form.add({
			xtype: 'displayfield',
			value: _(this.textAreaLabel)
		});

		this.input_textArea = this._form.add({
			xtype: 'textarea',
			width: '100%'
		});
	},

	buildBar: function() {
		var item_bar = [];

		if(this.helpdesk) {
			item_bar.push(Ext.create('Ext.button.Button', {
				xtype: 'button',
				handler: this.go_to_helpdesk,
				scope: this,
				text: _('Go to Helpdesk'),
				minWidth: 75
			}));
		}

		item_bar.push('->');

		item_bar.push(Ext.create('Ext.button.Button', {
			xtype: 'button',
			handler: this.ok_button_function,
			scope: this,
			text: _('Ok'),
			minWidth: 75
		}));

		item_bar.push(Ext.create('Ext.button.Button', {
			xtype: 'button',
			handler: function() {
				this.close();
			},
			scope: this,
			text: _('Cancel'),
			minWidth: 75
		}));

		var bar = new Ext.toolbar.Toolbar({
			ui: 'footer',
			dock: 'bottom',
			items: item_bar
		});

		return bar;
	},

	go_to_helpdesk: function() {
		if(this.helpdesk.indexOf('http://') === -1 && this.helpdesk.indexOf('https://') === -1 && this.helpdesk.indexOf('mailto:') === -1) {
			this.helpdesk = 'http://' + this.helpdesk;
		}

		log.debug('Go to the helpdesk ' + this.helpdesk, this.logAuthor);
		window.open(this.helpdesk, '_newtab');
	},

	ok_button_function: function() {
		log.debug('Send Event', this.logAuthor);
		var event = Ext.clone(this.base_event);
		event.output = this.input_textArea.getValue();

		if(this._component) {
			event.component = this._component;
		}

		if(this.referer) {
			event.referer = this.referer;
		}

		log.dump(event);

		global.eventsCtrl.sendEvent(event);
		this.close();
	}
});
