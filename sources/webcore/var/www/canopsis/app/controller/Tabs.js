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
				//add: this.on_add,
				//remove: this.on_remove,
				afterrender: function() {
					if (! this.tabpanel_rendered) {
						this.open_dashboard();
						this.open_saved_view();
						this.tabpanel_rendered = true;
					}
				}
			}
		});

		this.store = Ext.data.StoreManager.lookup('Tabs');

		if (!Ext.isIE) {
			this.store.proxy.id = this.store.proxy.id + '.' + global.account.user;
			this.store.load();
		}

		global.tabsCtrl = this;
	},

	clearTabsCache: function() {
		this.store.proxy.clear();
	},

	doRedraw: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();
		tab.doRedraw();
	},

  	on_tabchange: function(tabPanel, new_tab, old_tab, object) {
		//log.debug('Tabchange', this.logAuthor);
		tabPanel.old_tab = old_tab;
	},

	reload_active_view: function() {
		log.debug('Reload active view', this.logAuthor);
		var tab = Ext.getCmp('main-tabs').getActiveTab();
		tab.removeAll(true);
		tab.displayed = false;
		tab.autoshow = true;
		tab.getView();
	},

	open_dashboard: function() {
		var dashboard_id = this.getController('Account').getConfig('dashboard', 'view._default_.dashboard');
		log.debug('Open dashboard: ' + dashboard_id, this.logAuthor);
		return this.open_view({ view_id: dashboard_id, title: _('Dashboard'), closable: false, save: false }, 0);
	},

	open_saved_view: function() {
		var views = [];

		this.store.each(function(record) {
			var options = record.get('options');
			views.push(options);
		}, this);

		this.store.proxy.clear();

		log.debug('Load saved tabs:', this.logAuthor);
		for (var i=0; i < views.length; i++) {
			var options = views[i];
			log.debug(' + ' + options.title + '(' + options.view_id + ')', this.logAuthor);
			options.autoshow = false;

			var tab = this.open_view(options);
			if (! tab) {
				log.debug('Invalid view options:', this.logAuthor);
				//log.dump(options)
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

		if (options.view_id) {
			var maintabs = Ext.getCmp('main-tabs');

			if (options.tab_id) {
				var tab_id = options.view_id + options.tab_id + '.tab';
			} else {
				var tab_id = options.view_id + '.tab';
			}

			var tab = Ext.getCmp(tab_id);

			if (tab) {
				log.debug(' - Tab allerady open, just show it', this.logAuthor);
				maintabs.setActiveTab(tab);
			}else {
				log.debug(' - Create tab ...', this.logAuthor);
				log.debug('    - Get view config (' + options.view_id + ') ...', this.logAuthor);

				var localstore_record = false;
				if (options.save) {
					// archive tab in store
					log.debug("Add '" + options.title + "' ('" + options.view_id + "') in localstore ...", this.logAuthor);
					localstore_record = this.store.add({options: options});
					this.store.save();
				}

				var tab = {
					title: _(options.title),
					id: tab_id,
					iconCls: [options.iconCls],
					view_id: options.view_id,
					xtype: 'TabsContent',
					closable: options.closable,
					options: options.options,
					autoshow: options.autoshow,
					localstore_record: localstore_record
				};

				if (index != undefined) {
					tab = maintabs.insert(index, tab);
				}else {
					tab = maintabs.add(tab);
				}

				if (options.autoshow) {
					tab.show();
				}

				return tab;
			}
		}
	},

	save_active_view: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();
		if (tab.edit) {
			var right = this.getController('Account').check_right(tab.view, 'w');
			if (right == true) {
				tab.saveJqGridable();
			}else {
				global.notify.notify(_('Access denied'), _('You don\'t have the rights to modify this object'), 'error');
			}
		}
	},

	create_new_view: function() {
		this.getController('View').create_new_view();
	},

	edit_active_view: function() {
		var tab = Ext.getCmp('main-tabs').getActiveTab();
		if (! tab.edit) {
			if (!tab.report_window) {
				var right = this.getController('Account').check_right(tab.view, 'w');

				right = right && ! tab.view.internal;

				if (right == true) {
					tab.editMode();
				}else {
					global.notify.notify(_('Access denied'), _('You don\'t have the rights to modify this object'), 'error');
				}
			}else {
				global.notify.notify(_('Information'), _('Please close reporting before editing the view'), 'info');
			}
		}
	}

  	/*on_add: function(component, index, object){
		log.debug('Added', this.logAuthor);
	},

	on_remove: function(component, object){
		log.debug('Removed', this.logAuthor);
	}*/

});
