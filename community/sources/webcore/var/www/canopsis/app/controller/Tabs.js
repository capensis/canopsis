//need:app/store/Tabs.js,app/view/Tabs/View.js,app/view/Tabs/Content.js
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
Ext.define('canopsis.controller.Tabs', {
	extend: 'Ext.app.Controller',

	logAuthor: '[controller][tabs]',

	stores: ['Tabs'],
	views: ['Tabs.View', 'Tabs.Content'],

	tabpanel_rendered: false,

	init: function() {
		this.control({
			'tabpanel': {
				tabchange: this.on_tabchange,

				afterrender: function() {
					if(!this.tabpanel_rendered) {
						this.open_dashboard();
						this.open_saved_views();
						this.tabpanel_rendered = true;
					}
				},
				reload_active_view: this.reload_active_view,
				AutoRotateView: this.auto_rotate_view
			}
		});

		this.store = Ext.data.StoreManager.lookup('Tabs');

		if(!Ext.isIE) {
			this.store.proxy.id = this.store.proxy.id + '.' + global.account.user;
			this.store.load();
		}

		global.tabsCtrl = this;

		//taskRotate
		this.taskRotate = {
			run: this.openNextTab,
			interval: 5 * 60 * 1000,
			scope: this
		};

		this.taskRotateIndex = 0;
	},

	auto_rotate_view:function(_switch, delay) {
		if(delay) {
			this.taskRotate.interval = delay * 60 * 1000;
		}

		if(_switch) {
			this.taskRotateIndex = 0;
			Ext.TaskManager.start(this.taskRotate);
		}
		else {
			Ext.TaskManager.stop(this.taskRotate);
		}
	},

	openNextTab: function(){
		var maintabs = Ext.getCmp('main-tabs');
		var maintabLength =  maintabs.items.length;

		if((this.taskRotateIndex) === maintabLength) {
			maintabs.setActiveTab(0);
			this.taskRotateIndex = 0;
		}
		else {
			maintabs.setActiveTab(this.taskRotateIndex);
			this.taskRotateIndex++;
		}
	},

	clearTabsCache: function() {
		this.store.proxy.clear();
	},

	doRedraw: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();
		tab.doRedraw();
	},

	on_tabchange: function(tabPanel, new_tab, old_tab) {
		void(new_tab);

		tabPanel.old_tab = old_tab;
	},

	reload_active_view: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();

		// In edit mode
		if(tab.edit) {
			return;
		}

		log.debug('Reload active view', this.logAuthor);

		if (tab.displayed) {
			tab.removeAll(true);
			tab.displayed = false;
			tab.autoshow = true;
			tab.getView();
		}
		else {
			tab.autoshow = true;
			tab.getView();
		}
	},

	open_dashboard: function() {
		var dashboard_id = this.getController('Account').getConfig('dashboard', 'view._default_.dashboard');
		log.debug('Open dashboard: ' + dashboard_id, this.logAuthor);

		var title = _('Dashboard');

		// Get original Title
		var view_store = Ext.data.StoreManager.lookup('Views');

		if(view_store) {
			if(view_store.loaded) {
				var view = view_store.getById(dashboard_id);

				if(view) {
					title = view.get('crecord_name');
				}

				var tab = this.open_view({
					view_id: dashboard_id,
					title: title,
					closable: false,
					save: false,
					iconCls: 'icon-bullet-green'
				}, 0);

				// Set history
				Ext.getCmp('main-tabs').old_tab = tab;

				return tab;
			}
			else {
				// Load dashboard after view store is loaded
				view_store.on("load", this.open_dashboard, this, {single: true});
				view_store.load();
			}
		}

		return false;
	},

	open_saved_views: function() {
		var views = [];

		this.store.each(function(record) {
			var options = record.get('options');
			views.push(options);
		}, this);

		this.store.proxy.clear();

		log.debug('Load saved tabs:', this.logAuthor);

		for(var i = 0; i < views.length; i++) {
			var options = views[i];
			log.debug(' + ' + options.title + ' (' + options.view_id + ')', this.logAuthor);
			options.autoshow = false;

			var tab = this.open_view(options);

			if(!tab) {
				log.debug('Invalid view options:', this.logAuthor);
				log.dump(options);
			}
		}

		var maintabs = Ext.getCmp('main-tabs');
		maintabs.setActiveTab(0);
	},

	open_view: function(args, index) {
		//default options
		var default_options = {
			title: _('no title'),
			closable: true,
			options: {},
			autoshow: true,
			save: true,
			tab_id: undefined,
			index: undefined,
			iconCls: 'icon-bullet-orange'
		};

		var options = extend(default_options, args);

		log.debug("Open view tab '" + options.view_id + "'", this.logAuthor);

		if(options.view_id) {
			var maintabs = Ext.getCmp('main-tabs');
			var tab_id = undefined;

			if(options.tab_id) {
				tab_id = options.view_id + options.tab_id + '.tab';
			}
			else {
				tab_id = options.view_id + '.tab';
			}

			var tab = Ext.getCmp(tab_id);

			if(tab) {
				log.debug(' + Tab already openned, just show it', this.logAuthor);
				maintabs.setActiveTab(tab);
			}
			else {
				log.debug(' + Create tab ...', this.logAuthor);

				var localstore_record = false;

				if(options.save) {
					// archive tab in store
					log.debug("Add '" + options.title + "' ('" + options.view_id + "') in localstore ...", this.logAuthor);
					localstore_record = this.store.add({options: options});
					this.store.save();
				}

				tab = {
					title: options.title,
					id: tab_id,
					iconCls: [options.iconCls],
					view_id: options.view_id,
					xtype: 'TabsContent',
					closable: options.closable,
					options: options.options,
					autoshow: options.autoshow,
					localstore_record: localstore_record
				};

				log.debug(' + Tab options:', this.logAuthor);
				log.dump(tab);

				if(index !== undefined) {
					tab = maintabs.insert(index, tab);
				}
				else {
					tab = maintabs.add(tab);
				}

				if(options.autoshow){
					maintabs.setActiveTab(tab);
					tab.show();
				}

				// Dashboard If only one tabs
				if(options.autoshow && !tab.displayed && index === 0 && maintabs.items.length === 1) {
					tab.getView();
				}

				return tab;
			}
		}
	},

	save_active_view: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();

		if(tab.edit) {
			var right = this.getController('Account').check_right(tab.view, 'w');

			if(right === true) {
				tab.saveJqGridable();
			}
			else {
				global.notify.notify(_('Access denied'), _('You don\'t have the rights to modify this object'), 'error');
			}
		}
	},

	cancel_active_view: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();

		if(tab.edit) {
			tab.cancel();
		}
	},

	create_new_view: function() {
		this.getController('View').create_new_view();
	},

	edit_active_view: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();

		// Already in edit mode
		if(tab.edit) {
			return;
		}

		if(!tab.view) {
			global.notify.notify(_('Information'), _("View isn't fully loaded."), 'info');
			return;
		}

		if(tab.report_window) {
			global.notify.notify(_('Information'), _('Please close reporting before editing the view'), 'info');
			return;
		}

		var right = this.getController('Account').check_right(tab.view, 'w');
		right = right && !tab.view.internal;

		if(right === true) {
			tab.editMode();
		}
		else {
			global.notify.notify(_('Access denied'), _("You don't have the rights to modify this object"), 'error');
		}
	}
});
