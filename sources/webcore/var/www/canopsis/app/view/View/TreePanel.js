//need:app/lib/view/ctree.js
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
Ext.define('canopsis.view.View.TreePanel' , {
	extend: 'canopsis.lib.view.ctree',

	alias: 'widget.ViewTreePanel',

	store: 'TreeStoreViews',
	model: 'View',

	leafName: 'view',

	reporting: true,
	opt_bar_export: true,

	listeners: {
		selectionchange: function(selectionModel, selected) {
			var store = selectionModel.getStore();
			var all_selected = selected.length == store.count();

			var cb = Ext.getCmp(this.cb_select_all_id);

			if(cb.getValue() != all_selected) {
				cb.onSelectionChange = true;
				cb.setValue(all_selected);
			}
		}
	},

	initComponent: function() {
		this.cb_select_all_id = Ext.id();

		this.columns = [{
			xtype: 'treecolumn',
			text: _('Name'),
			flex: 5,
			dataIndex: 'crecord_name',
			renderer: function(value, metaData, record) {
				void(metaData);

				return "<span name='view." + record.get('crecord_name') + "'></span>" + value;
			}
		},{
			text: _('Export Options'),
			flex: 1,
			menuDisabled: true,
			dataIndex: 'view_options',
			renderer: function(val, metaData, record) {
				void(metaData);

				if(val && record.raw && record.raw.crecord_type !== 'view_directory') {
					return val.pageSize + ' - ' + _(val.orientation);
				}
			}
		},{
			flex: 1,
			dataIndex: 'aaa_owner',
			renderer: rdr_clean_id,
			text: _('Owner')
		},{
			flex: 1,
			dataIndex: 'aaa_group',
			renderer: rdr_clean_id,
			text: _('Group')
		},{
			width: 80,
			align: 'center',
			text: _('Owner'),
			dataIndex: 'aaa_access_owner',
			renderer: rdr_access
		},{
			width: 60,
			align: 'center',
			text: _('Group'),
			dataIndex: 'aaa_access_group',
			renderer: rdr_access
		},{
			width: 60,
			align: 'center',
			text: _('Others'),
			dataIndex: 'aaa_access_other',
			renderer: rdr_access
		},{
			width: 130,
			align: 'center',
			text: _('Last modified'),
			dataIndex: 'crecord_write_time',
			renderer: rdr_tstodate
		},{
			width: 130,
			align: 'center',
			text: _('Creation date'),
			dataIndex: 'crecord_creation_time',
			hidden: true,
			renderer: rdr_tstodate
		},{
			xtype: 'actioncolumn',
			width: 25,
			menuDisabled: true,
			align: 'center',
			tooltip: _('Dump'),
			icon: './themes/canopsis/resources/images/Tango-Blue-Materia/16x16/actions/gtk-indent.png',
			handler: function(tree, rowIndex) {
				var rec = tree.getStore().getAt(rowIndex).raw;

				if(rec.crecord_type === 'view') {
					tree.fireEvent('getViewFile', rec._id);
				}
			}
		},{
			xtype: 'actioncolumn',
			width: 25,
			align: 'center',
			menuDisabled: true,
			tooltip: _('URL'),
			icon: './themes/canopsis/resources/images/icons/page_white_code.png',
			handler: function(tree, rowIndex) {
				var rec = tree.getStore().getAt(rowIndex).raw;

				if(rec.crecord_type === 'view') {
					var view = rec._id;
					var authkey = global.account.authkey;
					var url = Ext.String.format(
						'http://{0}{1}?fullscreenMode=true&view_id={2}&authkey={3}',
						window.location.host,
						window.location.pathname,
						view,
						authkey
					);

					var _window = Ext.widget('window', {
						resizable: false,
						constrain: true,
						constrainTo: Ext.getCmp('main-tabs').getActiveTab().id,
						title: _('URL to') + ' ' + rec.crecord_name,
						items: [{
							xtype: 'form',
							border: false,
							layout: {
								type: 'hbox',
								align: 'stretch'
							},
							items: [{
								margin: 3,
								xtype: 'textfield',
								border: false,
								readOnly: true,
								width: 600,
								value: url
							},{
								xtype: 'button',
								tooltip: _('Go to the page'),
								iconCls: 'icon-page-go',
								height: '100%',
								margin: 3,
								handler: function() {
									window.open(url, '_blank');
								}
							}]
						}]
					});

					_window.show();
				}
			}
		},{
			xtype: 'actioncolumn',
			width: 25,
			menuDisabled: true,
			align: 'center',
			tooltip: _('Options'),
			icon: './themes/canopsis/resources/images/icons/cog.png',
			handler: function(tree, rowIndex) {
				var rec = tree.getStore().getAt(rowIndex).raw;

				if(rec.crecord_type === 'view') {
					tree.fireEvent('OpenViewOption', rec);
				}

			}
		}];

		if(global.reporting === true) {
			this.columns.push({
				width: 20,
				renderer: rdr_export_button
			});
		}

		this.columns.push({
			width: 16
		});

		this.callParent(arguments);

		var me = this;
		var buttons = [
			{
				xtype: 'button',
				iconCls: 'icon-import',
				text: _('Import view'),
				disabled: false,
				action: 'import'
			},{
				xtype: 'button',
				iconCls: 'icon-export',
				text: _('Export view as JSON'),
				disabled: false,
				action: 'exportjson'
			},{
				xtype: 'checkboxfield',
				id: this.cb_select_all_id,
				boxLabel: _('Select all'),
				disabled: false,
				onSelectionChange: false,
				handler: function(checkbox, checked) {
					if(!checkbox.onSelectionChange) {
						if(checked) {
							me.getSelectionModel().selectAll();
						}
						else {
							me.getSelectionModel().deselectAll();
						}
					}

					checkbox.onSelectionChange = false;
				}
			}
		];

		for(var i = 0; i < buttons.length; i++) {
			this.dockedToolbar.add(buttons[i]);
		}
	},

	export_pdf: function(view) {
		this.fireEvent('exportPdf', view);
	},

	add_to_context_menu: function(item_array) {
		item_array.push(
			Ext.create('Ext.Action', {
				iconCls: 'icon-preferences',
				text: _('Options'),
				action: 'OpenViewOption'
			})
		);

		return item_array;
	}
});
