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
Ext.define('canopsis.controller.Keynav', {
	extend: 'Ext.app.Controller',

	logAuthor: '[controller][Keynav]',

	init: function() {
		log.debug('Map keynav', this.logAuthor);

		this.map = Ext.create('Ext.util.KeyMap', Ext.getBody(), [
			{
				key: Ext.EventObject.N,
				ctrl: true,
				scope: this,
				fn: function(key, event) {
					void(key);

					event.stopEvent();
					// New view
					log.debug('ctrl + n', this.logAuthor);
				}
			},{
				key: Ext.EventObject.R,
				ctrl: true,
				scope: this,
				fn: function(key, event) {
					void(key);

					event.stopEvent();
					// Reload view
					log.debug('ctrl + r', this.logAuthor);

					var ctrl = this.getController('Tabs');
					ctrl.reload_active_view();

				}
			},{
				key: Ext.EventObject.E,
				ctrl: true,
				scope: this,
				fn: function(key, event) {
					void(key);

					event.stopEvent();
					// Edit view
					log.debug('ctrl + e', this.logAuthor);

					var ctrl = this.getController('Tabs');
					ctrl.edit_active_view();
				}
			},{
				key: Ext.EventObject.S,
				ctrl: true,
				scope: this,
				fn: function(key, event) {
					void(key);

					event.stopEvent();
					// Save view
					log.debug('ctrl + s', this.logAuthor);

					var ctrl = this.getController('Tabs');
					ctrl.save_active_view();
				}

			},{
				key: Ext.EventObject.X,
				ctrl: true,
				scope: this,
				fn: function(key, event) {
					void(key);

					event.stopEvent();
					// cancel view
					log.debug('ctrl + x', this.logAuthor);

					var ctrl = this.getController('Tabs');
					ctrl.cancel_active_view();
				}
			},{
				key: Ext.EventObject.F,
				ctrl: true,
				scope: this,
				fn: function(key, event) {
					void(key);

					event.stopEvent();
					// Search view
					log.debug('ctrl + f', this.logAuthor);
				}
			}
		]);

		this.callParent(arguments);
	}
});
