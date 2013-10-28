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
Ext.define('canopsis.lib.store.ctreeStore', {
	extend: 'Ext.data.TreeStore',

	autoLocalization: false,

	constructor: function(config) {
		this.callParent([config]);
		this.proxy.on('exception', this._manage_exception, this);
	},

	listeners: {
		move: function() {
			this.sync();
		},
		write: function(store, operation) {
			void(store);

			this.displaySuccess(operation);

			if(operation.success && this.afterCorrectWrite) {
				this.afterCorrectWrite();
			}
		}
	},

	displaySuccess: function(operation) {
		if (operation.success) {
			if(operation.action === 'create') {
				global.notify.notify(_('Success'), _('Record saved'), 'success');
			}
			else if (operation.action === 'destroy') {
				global.notify.notify(_('Success'), _('Record deleted'), 'success');
			}
			else if (operation.action === 'update') {
				global.notify.notify(_('Success'), _('Record updated'), 'success');
			}
		}
	},

	_manage_exception: function(store, request) {
		void(store);

		if(request.status === 403) {
			global.notify.notify(_('Access denied'), _('You don\'t have the rights to modify this object'), 'error');
			log.error(_('Access denied'));
			this.load();
		}
		else {
			log.error(_('Error while store synchronisation with server'));
		}
	}
});
