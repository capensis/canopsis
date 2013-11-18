//need:app/lib/view/cgrid.js,app/view/Briefcase/Uploader.js,app/lib/controller/cgrid.js
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
Ext.define('canopsis.view.Briefcase.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.BriefcaseGrid',

	requires: [
		'canopsis.view.Briefcase.Uploader',
		'canopsis.lib.controller.cgrid'
	],

	model: 'File',
	store: 'Files',

	opt_multiSelect: true,

	opt_bar_add: false,
	opt_view_element: true,
	opt_menu_rights: true,

	opt_menu_rename: true,
	opt_menu_set_avatar: true,

	opt_db_namespace: 'files',

	opt_bar_search: true,
	opt_bar_search_field: ['file_name'],

	opt_bar_customs: [{
		text: 'Add',
		xtype: 'button',
		iconCls: 'icon-add',
		handler: function() {
			var addFileWindow = Ext.create('canopsis.view.Briefcase.Uploader');
			addFileWindow.show();
		}
	}],

	columns: [{
			header: '',
			width: 25,
			sortable: false,
			renderer: rdr_file_type,
			dataIndex: 'content_type'
		},{
			header: _('Creation date'),
			flex: 1,
			sortable: true,
			dataIndex: 'crecord_creation_time',
			renderer: rdr_tstodate
		},{
			header: _('Name'),
			flex: 4,
			sortable: true,
			dataIndex: 'file_name'
		},{
			width: 120,
			dataIndex: 'aaa_owner',
			text: _('Owner'),
			renderer: rdr_clean_id
		},{
			width: 120,
			dataIndex: 'aaa_group',
			text: _('Group'),
			renderer: rdr_clean_id
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
		}
	],

	toggleSearchBarButtons: function(button, state) {
		if(!state) {
			state = false;
		}

		var tbar = this.getTbar();
		var items = tbar.items.items;

		for(var y = 0; y < items.length; y++) {
			var item = items[y];

			if(item.xtype === "button" && item !== button) {
				item.toggle(state);
			}
		}
	},

	initComponent: function() {
		this.bar_search = [{
			xtype: 'button',
			iconCls: 'icon-mimetype-pdf',
			pack: 'end',
			tooltip: _('Show pdf'),
			enableToggle: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(state) {
					button.filter_id = this.store.addFilter(
						{'content_type': 'application/pdf'}
					);

					this.toggleSearchBarButtons(button, false);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-mimetype-png',
			pack: 'end',
			tooltip: _('Show images'),
			enableToggle: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(state) {
					button.filter_id = this.store.addFilter(
						{'content_type': { $in: ['image/png', 'image/jpeg', 'image/gif', 'image/jpg']}}
					);

					this.toggleSearchBarButtons(button, false);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-mimetype-video',
			pack: 'end',
			tooltip: _('Show videos'),
			enableToggle: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(state) {
					button.filter_id = this.store.addFilter(
						{'content_type': 'video/ogg'}
					);

					this.toggleSearchBarButtons(button, false);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-mimetype-unknown',
			pack: 'end',
			tooltip: _('Show unknown'),
			enableToggle: true,
			scope: this,
			toggleHandler: function(button, state) {
				if(state) {
					button.filter_id = this.store.addFilter(
						{'content_type': null}
					);

					this.toggleSearchBarButtons(button, false);
				}
				else if(button.filter_id) {
					this.store.deleteFilter(button.filter_id);
				}

				this.store.load();
			}
		},'-'],

		this.ctrl = Ext.create('canopsis.lib.controller.cgrid');
		this.callParent(arguments);
	}
});
