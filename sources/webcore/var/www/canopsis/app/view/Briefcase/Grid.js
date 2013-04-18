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
Ext.define('canopsis.view.Briefcase.Grid' , {
	extend: 'canopsis.lib.view.cgrid',

	alias: 'widget.BriefcaseGrid',

	model: 'File',
	store: 'Files',

	opt_multiSelect: true,

	opt_bar_add: false,
	opt_view_element: true,
	opt_menu_rights: true,
	//opt_bar_download: true,
	//opt_menu_send_mail:true,

	opt_menu_rename: true,

	opt_db_namespace: 'files',

	opt_bar_search: true,
	opt_bar_search_field: ['file_name'],

	opt_bar_customs: [{
		text: 'Add', xtype: 'button', iconCls: 'icon-add', handler: function() {
			var addFileWindow = Ext.create('Ext.window.Window', {
				title: 'Add new file to briefcase',
				height: 110,
		    	width: 300,
		    	layout: 'fit',
		    	items: [
		    		Ext.create('Ext.form.Panel', {
						bodyPadding: '5 5 0',

						defaults: {
							allowBlank: false
						},

				        items: [{
				            xtype: 'filefield',
				            id: 'form-file',
				            emptyText: 'Select a file',
				            fieldLabel: 'File',
				            name: 'file-path',
				            buttonText: 'Browse',
				            width: 275
				        }],

				        buttons: [
				        	{
					            text: 'Upload',
					            handler: function(){
					                var form = this.up('form').getForm();
					                if (form.isValid()) {
					                	global.notify.notify('Uploading your file...');
					                    form.submit({
					                        url: '/files',
					                        success: function(fp, o) {
					                        	console.log(o);
					                        	global.notify.notify(_('Success'), _('File uploaded'), 'success');
					                            var store = Ext.getStore('Files')
					                            store.load();
					                            addFileWindow.close();
					                        },
					                        failure: function(fp, o) {
					                        	var code = o.result.data.code;
					                        	var msg = _('Unknown error');
					                        	if (code === 415) {
					                        		var msg = _('Unsupported Media Type');
					                        	} else if (code === 500) {
					                        		var msg = _('Internal server error');
					                        	} else if (code === 400) {
					                        		var msg = _('Bad request');
					                        	}
					                        	global.notify.notify(_('Failed'), msg, 'error');
					                        	addFileWindow.close();
					                        }
					                    });
					                }
					            }
					        }
					    ]
			    	})
		    	]
			});
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
		var state = state || false;
		var tbar = this.getTbar();
		var items = tbar.items.items;
		
		for (var y=0;y <items.length;y++) {
			var item = items[y];
			if (item.xtype == "button") {
				if (item != button) {
					item.toggle(state);
				}
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
				if (state) {
					button.filter_id = this.store.addFilter(
						{'content_type': 'application/pdf'}
					);
					this.toggleSearchBarButtons(button, false);
				} else {
					if (button.filter_id)
						this.store.deleteFilter(button.filter_id);
				}
				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-mimetype-png',
			pack: 'end',
			tooltip: _('Show png'),
			enableToggle: true,
			scope: this,
			toggleHandler: function(button, state) {
				if (state) {
					button.filter_id = this.store.addFilter(
						{'content_type': 'application/png'}
					);
					this.toggleSearchBarButtons(button, false);
				} else {
					if (button.filter_id)
						this.store.deleteFilter(button.filter_id);
				}
				this.store.load();
			}
		},{
			xtype: 'button',
			iconCls: 'icon-unknown',
			pack: 'end',
			tooltip: _('Show unknown'),
			enableToggle: true,
			scope: this,
			toggleHandler: function(button, state) {
				if (state) {
					button.filter_id = this.store.addFilter(
						{'content_type': null}
					);
					this.toggleSearchBarButtons(button, false);
				} else {
					if (button.filter_id)
						this.store.deleteFilter(button.filter_id);
				}
				this.store.load();
			}
		},'-'],
		this.ctrl = Ext.create('canopsis.lib.controller.cgrid');
		this.callParent(arguments);
	}
});
