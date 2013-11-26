//need:app/lib/view/cform.js
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

Ext.define('canopsis.view.Tabs.View_form' , {
	extend: 'canopsis.lib.view.cform',

	logAuthor: '[view option wizard]',
	EditMethod: 'win',

	items: [{
		'xtype': 'fieldset',
		'title': _('Report options'),
		'defaults': {anchor: '100%'},
		'collapsible': false,
		'items': [{
			'xtype': 'combobox',
			'name': 'orientation',
			'disable': false,
			'fieldLabel': _('Orientation'),
			'queryMode': 'local',
			'displayField': 'text',
			'valueField': 'value',
			'value': 'portrait',
			'store': {
				'xtype': 'store',
				'fields': ['value', 'text'],
				'data' : [
					{'value': 'portrait', 'text': _('Portrait')},
					{'value': 'landscape', 'text': _('Landscape')}
				]
			}
		},{
			'xtype': 'combobox',
			'name': 'pageSize',
			'fieldLabel': _('Page size'),
			'disable': false,
			'queryMode': 'local',
			'displayField': 'text',
			'valueField': 'value',
			'value': 'A4',
			'store': {
				'xtype': 'store',
				'fields': ['value', 'text'],
				'data' : [
					{'value': 'A3', 'text': 'A3'},
					{'value': 'A4', 'text': 'A4'}
				]
			}
		}]
	}],

	initComponent: function() {
		this.callParent(arguments);

		// binding
		var btn = this.down('button[action=cancel]');
		btn.on('click', this.cancel_button, this);
		btn = this.down('button[action=save]');
		btn.on('click', this.save_button, this);
	},

	cancel_button: function() {
		log.debug('cancel button', this.logAuthor);
		this.fireEvent('close');
	},

	save_button: function() {
		log.debug('save button', this.logAuthor);
		var variables = this.getValues();
		this.fireEvent('save', variables);
	}
});
