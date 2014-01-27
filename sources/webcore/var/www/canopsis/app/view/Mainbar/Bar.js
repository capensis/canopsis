//need:app/store/Files.js,app/view/Briefcase/Uploader.js,app/lib/menu/cspinner.js
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
Ext.define('canopsis.view.Mainbar.Bar', {
	extend: 'Ext.toolbar.Toolbar',

	alias: 'widget.Mainbar',

	requires: [
		'canopsis.store.Files',
		'canopsis.view.Briefcase.Uploader',
		'canopsis.lib.menu.cspinner'
	],

	border: false,

	layout: {
		type: 'hbox',
		align: 'stretch'
	},

	baseCls: 'Mainbar',

	initComponent: function() {
		this.localeSelector = Ext.create('Ext.form.field.ComboBox', {
			id: 'localeSelector',
			action: 'localeSelector',
			queryMode: 'local',
			displayField: 'text',
			valueField: 'value',
			fieldLabel: _('Language'),
			value: global.locale,
			store: {
				xtype: 'store',
				fields: ['value', 'text'],
				data: [
					{'value': 'fr', 'text': 'Français'},
					{'value': 'en', 'text': 'English'}
				]
			},
			iconCls: 'no-icon'
		});

		this.clockTypeSelector = Ext.create('Ext.form.field.ComboBox', {
			id: 'clockTypeSelector',
			action: 'clockTypeSelector',
			queryMode: 'local',
			displayField: 'text',
			valueField: 'value',
			fieldLabel: _('Clock type'),
			value: (global.account.clock_type) ? global.account.clock_type : 'auto',
			store: {
				xtype: 'store',
				fields: ['value', 'text'],
				data: [
					{'value': 'auto', 'text': 'Auto'},
					{'value': '24h', 'text': '24h'},
					{'value': '12h', 'text': '12h'}
				]
			},
			iconCls: 'no-icon'
		});

		this.viewSelector = Ext.create('Ext.form.field.ComboBox', {
			id: 'viewSelector',
			action: 'viewSelector',
			store: Ext.getStore('Views'),
			displayField: 'crecord_name',
			valueField: 'id',
			typeAhead: false,
			hideLabel: true,
			minChars: 2,
			queryMode: 'remote',
			emptyText: _('Select a view') + ' ...',
			width: 200
		});

		// Retrieve Files store add apply images filter
		var avatarSelectorStore = Ext.create('canopsis.store.Files',{
			storeId: 'Avatar',
			autoLoad: true
		});

		avatarSelectorStore.addFilter(
			{'content_type': { $in: ['image/png', 'image/jpeg', 'image/gif', 'image/jpg']}}
		);

		var avatarSelectorComboBox = Ext.create('Ext.form.field.ComboBox', {
			id: 'avatarSelector',
			action: 'avatarSelector',
			store: avatarSelectorStore,
			displayField: 'file_name',
			fieldLabel: _('Choose avatar'),
			valueField: '_id',
			value: global.account.avatar_id,
			iconCls: 'no-icon',
			width: 257 - 22 - 2
		});

		// Set if avatar change
		avatarSelectorStore.on('load', function() {
			avatarSelectorComboBox.select(global.account.avatar_id);
		}, {single: true});

		var avatarSelectorAdd = Ext.create('Ext.Button', {
			iconCls: 'icon-add',
			margin: '0 0 0 2',
			listeners: {
				click: function(me) {
					var uploader = Ext.create('canopsis.view.Briefcase.Uploader', {
						callback: function(file_id, filename) {
							global.accountCtrl.setAvatar(file_id, filename);
						}
					});

					me.up('menu').hide();
					uploader.show();
				}
			}
		});

		this.avatarSelector = Ext.create('Ext.container.Container', {
			iconCls: 'no-icon',
			layout: {
				type: 'hbox',
				align: 'right'
			},
			items: [avatarSelectorComboBox, avatarSelectorAdd]
		});

		this.dashboardSelector = Ext.create('Ext.form.field.ComboBox', {
			iconCls: 'icon-mainbar-dashboard',
			id: 'dashboardSelector',
			action: 'dashboardSelector',
			store: Ext.data.StoreManager.lookup('Views'),
			displayField: 'crecord_name',
			valueField: 'id',
			typeAhead: true,
			fieldLabel: _('Home view'),
			minChars: 2,
			queryMode: 'local',
			emptyText: _('Select a view') + ' ...',
			value: global.account['dashboard'],
			width: 200
		});

		// Hide  menu when item are selected
		this.viewSelector.on('select', function() {
			var menu = this.down('menu[name="Run"]');
			menu.hide();
		}, this);

		var menu_build = [];
		var menu_run = [];
		var menu_reporting = [];
		var menu_preferences = [];
		var menu_configuration = [];

		//Root build menu
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_account_admin')) {
			menu_build = menu_build.concat([
				{
					iconCls: 'icon-mainbar-edit-account',
					text: _('Accounts'),
					action: 'editAccount',
					viewId: 'view.account_manager'
				},{
					iconCls: 'icon-mainbar-edit-group',
					text: _('Groups'),
					action: 'editGroup',
					viewId: 'view.group_manager'
				}
			]);
		}

		//Build menu Curves Admin
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_curve_admin')) {
			menu_build = menu_build.concat([
				{
					iconCls: 'icon-mainbar-colors',
					text: _('Curves'),
					action: 'openViewMenu',
					viewId: 'view.curves'
				}
			]);
		}

		//Build menu Curves Admin
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_perfdata_admin')) {
			menu_build = menu_build.concat([
				{
					iconCls: 'icon-mainbar-perfdata',
					text: _('Perfdata'),
					action: 'openViewMenu',
					viewId: 'view.perfdata'
				}
			]);
		}

		//Root selector menu
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_selector_admin')) {
			menu_build = menu_build.concat([
				{
					iconCls: 'icon-mainbar-selector',
					text: _('Selectors'),
					action: 'openViewMenu',
					viewId: 'view.selector_manager'
				}
			]);
		}

		//Root selector menu
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_consolidation_admin')) {
			menu_build = menu_build.concat([
				{
					iconCls: 'icon-mainbar-consolidation',
					text: _('Consolidation'),
					action: 'openViewMenu',
					viewId: 'view.consolidation_manager'
				}
			]);
		}

		//Filter Rules menu
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_rule_admin')) {
			menu_build = menu_build.concat([
				{
					iconCls: 'icon-mainbar-filter',
					text: _('Filter Rules'),
					action: 'openViewMenu',
					viewId: 'view.rules_manager'
				}
			]);
		}

		//Topology menu
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_topology_admin')) {
			menu_build = menu_build.concat([
				{
					iconCls: 'icon-mainbar-topology',
					text: _('Topologies'),
					action: 'openViewMenu',
					viewId: 'view.topology_manager'
				}
			]);
		}

		//Build menu
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_view_admin') || global.accountCtrl.checkGroup('group.CPS_view')) {
			menu_build = menu_build.concat([
				{
					iconCls: 'icon-mainbar-edit-view',
					text: _('Edit active view'),
					action: 'editView'
				},{
					iconCls: 'icon-mainbar-new-view',
					text: _('New view'),
					action: 'newView'
				}
			]);
		}

		//Reporting menu
		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_reporting_admin')) {
			menu_reporting = menu_reporting.concat([
				{
					iconCls: 'icon-mimetype-pdf',
					text: _('Export active view'),
					action: 'exportView'
				}
			]);
		}


		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_schedule_admin')) {
			menu_reporting = menu_reporting.concat([
				{
					iconCls: 'icon-mainbar-add-task',
					text: _('Schedule active view export'),
					action: 'ScheduleExportView'
				},{
					iconCls: 'icon-mainbar-edit-task',
					text: _('Schedules'),
					action: 'editSchedule',
					viewId: 'view.schedule_manager'
				}
			]);
		}

		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_reporting_admin')) {
			menu_reporting = menu_reporting.concat([{
					iconCls: 'icon-mainbar-reporting',
					text: _('Switch to live reporting'),
					action: 'reportingMode'
				}
			]);
		}

		menu_reporting = menu_reporting.concat([{
			iconCls: 'icon-mainbar-eventLog',
			text: _('Event navigation'),
			action: 'eventLog_navigation',
			viewId: 'view.eventLog_navigation'
		}]);

		menu_reporting = menu_reporting.concat([{
			iconCls: 'icon-mainbar-briefcase',
			text: _('Briefcase'),
			action: 'openBriefcase',
			viewId: 'view.briefcase'
		}]);

		//Run menu
		menu_run = menu_run.concat([
			{
				iconCls: 'icon-mainbar-dashboard',
				text: _('Dashboard'),
				action: 'openDashboard'
			},{
				iconCls: 'icon-mainbar-viewdetails',
				text: _('Components'),
				action: 'openViewMenu',
				viewId: 'view.components'
			},{
				iconCls: 'icon-mainbar-viewdetails',
				text: _('Resources'),
				action: 'openViewMenu',
				viewId: 'view.resources'
			}
		]);

		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_view_admin') || global.accountCtrl.checkGroup('group.CPS_view')) {
			menu_run = menu_run.concat(
				[
					{
						iconCls: 'icon-mainbar-run',
						text: _('Views manager'),
						action: 'openViewsManager',
						viewId: 'view.view_manager'
					}
				]
			);
		}

		if(global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_derogation_admin') || global.accountCtrl.checkGroup('group.CPS_derogation_admin')) {
			menu_run = menu_run.concat(
				[
					{
						iconCls: 'icon-crecord_type-derogation',
						text: _('Derogation manager'),
						action: 'openDerogationManager',
						viewId: 'view.derogation_manager'
					}
				]
			);
		}

		if (global.accountCtrl.checkRoot() || global.accountCtrl.checkGroup('group.CPS_statemap_admin')) {
			menu_run = menu_run.concat(
				[
					{
						iconCls: 'icon-crecord_type-statemap',
						text: _('Statemap manager'),
						action: 'openStatemapManager',
						viewId: 'view.statemap_manager',
					}
				]
			);
		}

		menu_run = menu_run.concat([
			'-', this.viewSelector
		]);

		//Configuration menu
		menu_configuration = menu_configuration.concat([
			{
				iconCls: 'icon-clear',
				text: _('Clear tabs cache'),
				action: 'cleartabscache'
			},{
				iconCls: 'icon-access',
				text: _('Authentification key'),
				action: 'authkey'
			}
		]);

		this.changePass = undefined;

		if(!global.account.external) {
			this.changePass = {
				iconCls: 'no-icon',
				text: _('Change your password'),
				onClick: function() {
					var win = Ext.create("canopsis.view.Account.Password");
					win.show();
				}
			};
		}

		//Preferences menu
		menu_preferences = menu_preferences.concat([
			this.localeSelector,
			this.clockTypeSelector,
			'-',
			this.dashboardSelector,
			this.avatarSelector,
			this.changePass
		]);

		//Set Items
		this.items = [
			{
				iconCls: 'icon-mainbar-build',
				text: _('ITIL.Build'),
				menu: {
					items: menu_build
				}
			},{
				iconCls: 'icon-mainbar-run',
				text: _('ITIL.Run'),
				menu: {
					name: 'Run',
					showSeparator: true,
					items: menu_run
				}
			},{
				iconCls: 'icon-mainbar-report',
				text: _('ITIL.Report'),
				menu: {
					name: 'Report',
					showSeparator: true,
					items: menu_reporting
				}
			},'-', {
				xtype: 'tbtext',
				text: 'Canopsis',
				cls: 'cps-title',
				flex: 1
			},{
				xtype: 'container',
				name: 'clock',
				align: 'strech',
				flex: 4
			},'->',{
				xtype: 'container',
				html: "<div class='cps-account' >" + global.account.firstname + ' ' + global.account.lastname + '</div>',
				flex: 2.3
			},
				this.userPreferences
			,'-',{
				icon: '/account/getAvatar',
				iconCls: 'icon-mainbar icon-avatar-bar',
				width: 36,
				menu: {
					items: menu_preferences
				}

			},'-', {
				iconCls: 'icon-mainbar icon-preferences',
				width: 36,
				menu: {
					name: 'Preferences',
					showSeparator: true,
					items: menu_configuration
				}
			},{
				iconCls: 'icon-mainbar icon-about',
				width: 36,
				menu: {
					name: 'About',
					showSeparator: true,
					items: [
						{
							iconCls: 'icon-documentation',
							text: _('Documentation'),
							onClick: function() {
								window.open('https://github.com/capensis/canopsis-doc/wiki', '_blank');
							}
						},{
							iconCls: 'icon-community',
							text: _('Community'),
							onClick: function() {
								window.open('http://www.canopsis.org', '_blank');
							}
						},{
							iconCls: 'icon-forum',
							text: _('Forum'),
							onClick: function() {
								window.open('http://forums.monitoring-fr.org/index.php/board,127.0.html', '_blank');
							}
						},'-', {
							iconCls: 'icon-mainbar-sources',
							text: '<b>Commit</b>: ' + global.commit.substr(0, 10),
							onClick: function() {
								if(global.commit) {
									window.open('https://github.com/capensis/canopsis/commit/' + global.commit, '_blank');
								}
							}
						},{
							iconCls: 'icon-issue',
							text: _('Report a issue'),
							onClick: function() {
								window.open('https://github.com/capensis/canopsis/issues', '_blank');
							}
						},{
							iconCls: 'icon-github',
							text: _('Fork Me') + ' !',
							onClick: function() {
								window.open('https://github.com/capensis/canopsis', '_blank');
							}
						}
					]
				}
			},{
				iconCls: (global.websocketCtrl.connected) ? 'icon-mainbar icon-bullet-green' : 'icon-mainbar icon-bullet-red',
				id: 'Mainbar-menu-Websocket',
				onClick: function() {
					global.websocketCtrl.connect();
				}
			},
			Ext.create('canopsis.lib.menu.cspinner'), {
				iconCls: 'icon-mainbar icon-bootstrap-off',
				action: 'logout',
				tooltip: _('Logout')
			}
		];

		this.callParent(arguments);
	},

	afterRender: function() {
		this.callParent(arguments);

		Ext.create('Ext.tip.ToolTip', {
			target: Ext.getCmp('Mainbar-menu-Websocket').el,
			renderTo: Ext.getBody(),
			_htmlRender: function() {
				return Ext.String.format(
					'{0}: <b>{1}</b><br/><i>{2}</i>',
					_('Websocket is now'),
					global.websocketCtrl.connected ? _('connected') : _('disconnected'),
					global.websocketCtrl.connected ? '' : _('Check if websocket daemon is started or check your firewall port ') + global.nowjs.port
				);
			},
			listeners: {
				beforeshow: function() {
					this.update(this._htmlRender());
				}
			}
		});
	}
});
