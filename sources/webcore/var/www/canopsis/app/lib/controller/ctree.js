//need:app/lib/view/crights.js
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
Ext.define('canopsis.lib.controller.ctree', {
	extend: 'Ext.app.Controller',

	requires: [
		'canopsis.lib.view.crights'
	],

	logAuthor: '[controller][ctree]',

	init: function() {
		log.debug('[controller][ctree] - ' + this.id + ' Initialize ...');

		var control = {};
		control[this.listXtype] = {
		                afterrender: this._bindTreeEvents
		};
		this.control(control);

		this.callParent(arguments);

	},

	_bindTreeEvents: function(tree) {
		var btns = undefined;
		var i = undefined;

		var id = tree.id;
		this.tree = tree;

		log.debug('[controller][ctree] - Bind events "' + id + '" ...');

		// Bind Context Menu
		if(tree.contextMenu) {
			tree.on('itemcontextmenu', this._showMenu, this);

			//Duplicate button
			btns = Ext.ComponentQuery.query('#' + tree.contextMenu.id + ' [action=duplicate]');

			for(i = 0; i < btns.length; i++) {
				btns[i].on('click', this._duplicateButton, this);
			}

			//DeleteButton
			btns = Ext.ComponentQuery.query('#' + tree.contextMenu.id + ' [action=delete]');

			for(i = 0; i < btns.length; i++) {
				btns[i].on('click', this._deleteButton, this);
			}

			//edit rights
			btns = Ext.ComponentQuery.query('#' + tree.contextMenu.id + ' [action=rights]');

			for(i = 0; i < btns.length; i++) {
				btns[i].on('click', this._editRights, this);
			}

			//rename button
			btns = Ext.ComponentQuery.query('#' + tree.contextMenu.id + ' [action=rename]');

			for(i = 0; i < btns.length; i++) {
				btns[i].on('click', this._renameButton, this);
			}

			//rename button
			btns = Ext.ComponentQuery.query('#' + tree.contextMenu.id + ' [action=export]');

			for(i = 0; i < btns.length; i++) {
				btns[i].on('click', this._exportButton, this);
			}
		}

		// toolbar buttons
		//delete
		btns = Ext.ComponentQuery.query('#' + id + ' button[action=delete]');

		for(i = 0; i < btns.length; i++) {
			btns[i].on('click', this._deleteButton, this);
		}

		// Reload button
		btns = Ext.ComponentQuery.query('#' + id + ' button[action=reload]');

		for(i = 0; i < btns.length; i++) {
			btns[i].on('click', this._reloadButton, this);
		}

		// Add buttons
		btns = Ext.ComponentQuery.query('#' + id + ' button[action=add_leaf]');

		for(i = 0; i < btns.length; i++) {
			btns[i].on('click', this._addLeafButton, this);
		}

		btns = Ext.ComponentQuery.query('#' + id + ' button[action=add_directory]');

		for(i = 0; i < btns.length; i++) {
			btns[i].on('click', this._addDirectoryButton, this);
		}

		//duplicate
		btns = Ext.ComponentQuery.query('#' + id + ' button[action=duplicate]');

		for(i = 0; i < btns.length; i++) {
			btns[i].on('click', this._duplicateButton, this);
		}

		// general binding
		tree.on('selectionchange', this._selectionchange, this);
		tree.on('itemdblclick', this._itemDoubleClick, this);

		// before drop fonction listening
		tree.getView().on('beforedrop', function(n,d,o) {
			var stop_event = this._check_right_on_drop(n, d, o);

			if(stop_event === true) {
				return false;
			}
		}, this);


		// keep memory of last expanded node, expand it again after load
		tree.store.on('expand', function(node) {
			this.currentNode = node.getPath();
		}, this);

		tree.on('load', function() {
			if(this.currentNode) {
				this.tree.expandPath(this.currentNode);
			}
		}, this);

		if(this.bindTreeEvent) {
			this.bindTreeEvent();
		}

	},

	_check_right_on_drop: function(node, data, overModel, dropPosition, dropFunction, opts) {
		var stop_event;

		if(this.check_right_on_drop) {
			stop_event = this.check_right_on_drop(node, data, overModel, dropPosition, dropFunction, opts);
		}

		if(stop_event) {
			return true;
		}
	},

	_selectionchange: function(view, records) {
		var tree = this.tree;

		//Enable delete Button
		var btns = Ext.ComponentQuery.query('#' + tree.id + ' button[action=delete]');

		for(var i = 0; i < btns.length; i++) {
			btns[i].setDisabled(records.length === 0);
		}

		if(this.selectionchange) {
			this.selectionchange(view, records);
		}
	},

	_addDirectoryButton: function() {
		log.debug('add directory', this.logAuthor);

		if(this.addDirectoryButton) {
			this.addDirectoryButton();
		}
	},

	_addLeafButton: function() {
		log.debug('add leaf', this.logAuthor);

		if(this.addLeafButton) {
			this.addLeafButton();
		}
	},

	_duplicateButton: function() {
		log.debug('duplicate', this.logAuthor);

		if(this.duplicateButton) {
			this.duplicateButton();
		}
	},

	_deleteButton: function() {
		log.debug('[controller][ctree] - clicked deleteButton', this.logAuthor);
		var tree = this.tree;

		var verification = true;

		var selection = tree.getSelectionModel().getSelection();

		for(var i = 0; i < selection.length; i++) {
			if(this.getController('Account').check_record_right(selection[i], 'w')) {
				if(selection[i].childNodes.length > 0) {
					global.notify.notify(_('Directory not empty'), _('The directory must be empty if you want to remove it'), 'error');
					verification = false;
				}
				else if(this.checkOpen && this.checkOpen(selection[0].data.id)) {
					verification = false;
				}
			}
		}

		if(verification === true) {
			Ext.MessageBox.confirm(_('Confirm'), _('Are you sure you want to delete') + ' ' + selection.length + ' ' + _('items') + ' ?',
				function(btn) {
					if(btn === 'yes') {
						for(i = 0; i < selection.length; i++) {
							selection[i].remove();
						}
					}

					tree.store.sync();
				}
			);
		}

		if(this.deleteButton) {
			this.deleteButton(button, grid, selection);
		}
	},

	_reloadButton: function() {
		log.debug('[controller][ctree] - Reload store "' + this.tree.store.storeId + '" of ' + this.tree.id);

		if(!this.tree.store.isLoading()) {
			this.tree.store.load();
		}
	},

	_renameButton: function() {
		var tree = this.tree;
		var selection = tree.getSelectionModel().getSelection()[0];

		if(this.checkOpen && this.checkOpen(selection.data.id)) {
			return;
		}

		if(this.getController('Account').check_record_right(selection, 'w')) {
			Ext.Msg.prompt(_('View name'), _('Please enter view name:'), function(btn, new_name) {
				if (btn === 'ok') {
					selection.set('crecord_name', new_name);
					this.tree.store.sync();
				}
			}, this, undefined, selection.get('crecord_name'));
		}
		else {
			global.notify.notify(_('Access denied'), _('You don\'t have the rights to modify this object'), 'error');
		}
	},

	_editRights: function() {
		log.debug('Edit rights', this.logAuthor);
		var tree = this.tree;
		var selection = tree.getSelectionModel().getSelection()[0];

		//create form
		if(this.getController('Account').check_record_right(selection, 'w')) {
			var config = {
				data: selection,
				renderTo: tree.id,
				constrain: true
			};

			var crights = Ext.create('canopsis.lib.view.crights', config);

			//listen to save event to refresh store
			crights.on('save', function() {
				tree.store.load();
			}, this);

			crights.show();
		}
		else {
			global.notify.notify(_('Access denied'), _('You don\'t have the rights to modify this object'), 'error');
		}
	},

	_showMenu: function(view, rec, node, index, e) {
		void(node, index);

		var selection = this.tree.getSelectionModel().getSelection();

		if(selection.length < 2) {
			view.select(rec);
		}

		this.tree.contextMenu.showAt(e.getXY());
		return false;
	},

	_itemDoubleClick: function(record) {
		if(this.itemDoubleClick) {
			this.itemDoubleClick(record);
		}
	},

	_exportButton: function() {
		log.debug('Export function', this.logAuthor);
		var tree = this.tree;
		var selection = tree.getSelectionModel().getSelection()[0];

		if(this.exportButton) {
			this.exportButton(selection);
		}
	}
});
