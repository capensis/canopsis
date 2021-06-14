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
Ext.define('canopsis.lib.menu.cclear' , {
	extend: 'Ext.menu.Menu',

	grid: undefined,

	initComponent: function() {
		if (this.grid) {
			this.grid_store = this.grid.getStore();

			this.clearAllButton = Ext.create('Ext.Action', {
				iconCls: 'icon-delete',
				text: _('Clear all'),
				scope: this,
				handler: this.action_clearall
			});

			this.deleteButton = Ext.create('Ext.Action', {
				iconCls: 'icon-delete',
				text: _('Delete selected'),
				scope: this,
				handler: this.action_delete
			});

			this.items = [this.deleteButton, this.clearAllButton];

			// Bind Context Menu
			this.grid.on('itemcontextmenu', this.action_open, this);
		}

		this.callParent(arguments);
	},

	action_open: function(view, rec, node, index, e) {
		void(node, index);

		e.preventDefault();

		//don't auto select if multi selecting
		var selection = this.grid.getSelectionModel().getSelection();

		if(selection.length < 2) {
			view.select(rec);
		}

		this.showAt(e.getXY());
		return false;
	},

	action_delete: function() {
		var selection = this.grid.getSelectionModel().getSelection();

		for(var i = 0; i < selection.length; i++) {
			this.grid_store.remove(selection[i]);
		}
	},

	action_clearall: function() {
		this.grid_store.removeAll();
	}
});
