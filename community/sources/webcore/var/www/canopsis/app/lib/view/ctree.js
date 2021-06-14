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
Ext.define('canopsis.lib.view.ctree' , {
	extend: 'Ext.tree.Panel',

	animCollapse: false,

	useArrows: true,
	rootVisible: false,
	multiSelect: true,

	leafName: 'leaf',

	//options
	opt_bar: true,
	opt_bar_add: true,
	opt_bar_add_directory: true,
	opt_bar_delete: true,
	opt_bar_duplicate: true,
	opt_bar_reload: true,

	opt_bar_export: false,

	opt_menu_rights: true,
	opt_menu_rename: true,

	initComponent: function() {
		var bar_child = [];

		if(this.opt_bar_add) {
			bar_child.push({
				xtype: 'button',
				iconCls: 'icon-leaf-add',
				text: _('Add') + ' ' + this.leafName,
				action: 'add_leaf'
			});
		}

		if(this.opt_bar_add_directory) {
			bar_child.push({
				xtype: 'button',
				iconCls: 'icon-folder-add',
				text: _('Add directory'),
				action: 'add_directory'
			});
		}

		if(this.opt_bar_reload) {
			bar_child.push({
				xtype: 'button',
				iconCls: 'icon-reload',
				text: _('Reload'),
				action: 'reload'
			});
		}

		if(this.opt_bar_delete) {
			bar_child.push({
				xtype: 'button',
				iconCls: 'icon-delete',
				text: _('Delete'),
				disabled: true,
				action: 'delete'
			});
		}

		// creating toolbar
		if(this.opt_bar_bottom) {
			this.bbar = Ext.create('Ext.toolbar.Toolbar', {
				items: bar_child
			});

			this.dockedToolbar = this.bbar;
		}
		else {
			this.tbar = Ext.create('Ext.toolbar.Toolbar', {
				items: bar_child
			});

			this.dockedToolbar = this.tbar;
		}

		// Context menu
		if(this.opt_bar) {
			var myArray = [];

			if(this.opt_bar_delete) {
				myArray.push(
					Ext.create('Ext.Action', {
						iconCls: 'icon-delete',
						text: _('Delete'),
						action: 'delete'
					})
				);
			}

			if(this.opt_bar_duplicate) {
				myArray.push(
					Ext.create('Ext.Action', {
						iconCls: 'icon-copy',
						text: _('Duplicate'),
						action: 'duplicate'
					})
				);
			}

			if(this.opt_menu_rights === true) {
				myArray.push(
					Ext.create('Ext.Action', {
						iconCls: 'icon-access',
						text: _('Rights'),
						action: 'rights'
					})
				);
			}

			if(this.opt_menu_rename === true) {
				myArray.push(
					Ext.create('Ext.Action', {
						iconCls: 'icon-edit',
						text: _('Rename'),
						action: 'rename'
					})
				);
			}

			if(this.opt_bar_export === true) {
				myArray.push(
					Ext.create('Ext.Action', {
						iconCls: 'icon-mainbar-add-task',
						text: _('Schedule export'),
						action: 'export'
					})
				);
			}

			if(this.add_to_context_menu) {
				this.add_to_context_menu(myArray);
			}

			if(myArray.length !== 0) {
				this.contextMenu = Ext.create('Ext.menu.Menu', {
					items: myArray
				});
			}
		}

		this.callParent(arguments);

		//fix the non destroy plugin with manualy instantiate it
		var innerView = this.getView();

		this.ddplugin = Ext.PluginManager.create({
			ptype: 'treeviewdragdrop'
		});

		this.ddplugin.init(innerView);

		innerView.animate = false;
	},

	destroy: function() {
		//manualy destroy the plugin
		this.ddplugin.destroy();
		Ext.tree.Panel.superclass.destroy.call(this);
	}
});
