//need:app/lib/view/cwizard.js
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

Ext.define('canopsis.view.Tabs.ViewWizard' , {
	extend: 'canopsis.lib.view.cwizard',

	title: _('View options wizard'),

	data: undefined,

	logAuthor: '[widget option wizard]',

	initComponent: function() {
		// Build wizard options
		var step1 = {
			title: _('Choose widget'),

			items: [{
				'xtype' : 'colorfield',
				'name' : 'background',
				'fieldLabel': 'Background color',
				'value': 'FFFFFF'
			}]
		};

		this.step_list = [step1];

		this.callParent(arguments);
	},

	cancel_button: function() {
		log.debug('cancel button', this.logAuthor);
		this.fireEvent('cancel', this.widgetId);
		this.close();
	},

	finish_button: function() {
		log.debug('save button', this.logAuthor);
		var variables = this.get_variables();
		this.fireEvent('save', variables);
		this.close();
	}
});
