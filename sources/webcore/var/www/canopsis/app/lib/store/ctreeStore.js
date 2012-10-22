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
Ext.define('canopsis.lib.store.ctreeStore', {
	extend: 'Ext.data.TreeStore',

	autoLocalization: false,

	constructor: function(config) {
        this.callParent(arguments);

        this.proxy.on('exception', this._manage_exception, this);
    },

	//raise an exception if server didn't accept the request
	//and display a popup if the store is modified
	listeners: {
	/*	exception: function(proxy, response, operation){
			Ext.MessageBox.show({
				title: _('REMOTE EXCEPTION'),
				msg: this.storeId + ': ' + _('request failed'),
				icon: Ext.MessageBox.ERROR,
				buttons: Ext.Msg.OK
			});
			log.debug(response);
		},*/
		load: function(store, node, records) {
			if (this.autoLocalization) {
				var i;
				for (i in records) {
					record = records[i];
					record.set('text', _(record.get('text')));
					record.modified = false;
				}
			}
		},
		write: function(store, operation,option) {
			if (operation.action == 'create')
				global.notify.notify(_('Success'), _('Record saved'), 'success');
			if (operation.action == 'destroy')
				global.notify.notify(_('Success'), _('Record deleted'), 'success');
		}
   },

   	_manage_exception: function(store, request, options) {
		if (request.status == 403) {
			global.notify.notify(_('Access denied'), _('You don\'t have the rights to modify this object'), 'error');
			log.error(_('Access denied'));
			this.load();
		}else {
			log.error(_('Error while store synchronisation with server'));
		}
	}
});
