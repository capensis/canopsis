//need:app/lib/view/cwizard.js,app/view/Tabs/View_form.js,app/view/ReportingBar/ReportingBar.js
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
Ext.define('canopsis.view.Tabs.Content' , {
	extend: 'Ext.jq.Gridable',

	alias: 'widget.TabsContent',

	requires: [
		'canopsis.lib.view.cwizard',
		'canopsis.view.Tabs.View_form',
		'canopsis.view.ReportingBar.ReportingBar'
	],

	logAuthor: '[view][tabs][content]',

	autoScroll: true,

	autoshow: true,
	displayed: false,

	items: [],

	border: false,

	//Ext.jq.Gridable
	spotlight: true,
	contextMenu: true,

	debug: false,

	autoScale: true,
	autoDraw: false,
	wizard: 'canopsis.lib.view.cwizard',
	view_option_form: 'canopsis.view.Tabs.View_form',

	// Export an report
	reportMode: false,
	exportMode: false,
	export_from: undefined,
	export_to: undefined,
	fullscreenMode: false,

	record: undefined,

	view_option_window: undefined,

	dump: undefined,
	view: undefined,
	view_options: undefined,

	onClose: false,

	spinner: undefined,
	spinner_delay: 300,
	spinner_options: {
		color: '#808080',
		lines: 11,
		width: 18,
		length: 30,
		shadow: false,
		radius: 40,
		top: 'auto',
		left: 'auto'
	},

	//Locales
	locales: {
		save: _('Save'),
		column: _('Add column'),
		row: _('Add row'),
		editMode: _('Edit mode'),
		viewMode: _('View mode'),
		removeAll: _('Remove all'),
		del: _('Delete'),
		cancel: _('Cancel'),
		edit: _('Edit'),
		duplicate: _('Duplicate'),
		configure: _('Configure'),
		changePageMode: _('Toggle to real page size'),
		options: _('Edit options')
	},

	pageWidth: {
		'portrait': {'A4': 940, 'A3': 1090 },
		'landscape': {'A4': 1090, 'A3': 1585}
	},

	//Logging
	log: function(message) {
		log.debug(message, this.logAuthor);
	},

	//Init
	initComponent: function() {
		this.callParent(arguments);

		// Display timer if view is too long to load
		this.on('afterlayout', function() {
			Ext.defer(function() {
				if(!this.displayed) {
					log.debug('Display view spinner', this.logAuthor);
					tEl = document.getElementById(this.JqgContainerId);
					this.spinner = new Spinner(this.spinner_options).spin(tEl);
				}
			}, this.spinner_delay, this);
		}, this, { single: true });

		log.debug("Display view '" + this.view_id + "' ...", this.logAuthor);

		this.options = {
			reportMode: this.reportMode,
			exportMode: this.exportMode,
			export_from: this.export_from,
			export_to: this.export_to,
			fullscreenMode: this.fullscreenMode
		};

		this.on('save', this.saveView, this);

		this.on('beforeclose', this.beforeclose);

		//binding event to save resources
		this.on('show', this._onShow, this);
		this.on('hide', this._onHide, this);
		this.on('resizeWidget', this.onResizeWidget, this);

		//Apply view options when loaded

		if(this.fullscreenMode || this.exportMode) {
			this.getView();
		}
	},

	getView: function() {
		Ext.Ajax.request({
			url: '/rest/object/view/' + this.view_id,
			scope: this,
			success: function(response) {
				var data = Ext.JSON.decode(response.responseText);
				var view = data.data[0];
				var dump = view.items;

				if(dump === undefined || dump === '') {
					dump = [];
				}

				this.view = view;
				this.dump = dump;

				var view_options = {
					orientation: 'portrait',
					pageSize: 'A4'
				};

				if(view.view_options) {
					view_options = Ext.Object.merge(view_options, view.view_options);
				}

				this.view_options = view_options;

				var width = this.pageWidth[view_options.orientation][view_options.pageSize];

				this.pageModeSize = width;

				log.debug('getView:', this.logAuthor);
				log.debug(' + view_options:', this.logAuthor);
				log.dump(view_options);
				log.debug(' + width: ' + width, this.logAuthor);
				log.debug(' + nb widget: ' + dump.length, this.logAuthor);
				log.debug(' + dump:', this.logAuthor);
				log.dump(dump);

				// Set width for exporting in PDF
				if(this.exportMode) {
					log.debug('Orientation: ' + view_options.orientation, this.logAuthor);
					log.debug('pageSize: ' + view_options.pageSize, this.logAuthor);
					log.debug('width: ' + width, this.logAuthor);

					this.setWidth(width);

				}

				if(this.spinner) {
					log.debug('Remove view spinner', this.logAuthor);
					this.spinner.stop();
					delete this.spinner;
				}

				this.setContent(dump);

			},
			failure: function(result, request) {
				void(result);

				log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
				log.error('Close tab, maybe not exist ...', this.logAuthor);
				global.notify.notify(_('Warning'), _('Impossible to get view options.'), 'warning');
				this.destroy();
			}
		});
	},

	applyViewOptions: function(options) {
		log.debug('Apply view options', this.logAuthor);

		if(!options) {
			options = this.view_options;
		}

		if(options && options['background']) {
			this.body.setStyle('background', '#' + options['background']);
		}
	},

	onResizeWidget: function(cmp) {
		cmp.onResize();
	},

	setContent: function(dump) {
		if(!dump) {
			dump = this.dump;
		}

		if(dump && !this.displayed) {
			log.debug('setContent of ' + this.view_id, this.logAuthor);
			this.load(dump);
		}

		this.displayed = true;
	},

	saveView: function(dump) {
		//ajax request with dump sending
		log.debug('Saving view ' + this.view_id, this.logAuthor);

		if(dump === undefined) {
			dump = this.dumpJqGridable();
		}

		var view_options = undefined;

		if(this.getViewOptions) {
			view_options = this.getViewOptions();
		}

		if(!view_options && this.view_options) {
			view_options = this.view_options;
		}

		if(!view_options) {
			view_options = {
				orientation: 'portrait',
				pageSize: 'A4'
			};
		}

		// Update view
		if(this.view_id) {
			var data = {
				'crecord_name': this.view.crecord_name,
				'items': dump,
				'view_options': view_options
			};

			updateRecord('object', 'view', 'canopsis.model.View', this.view_id, data);
		}

		this.dump = dump;

		this.startAllTasks();
	},

	//Binding
	_onShow: function() {
		if(this.onClose) {
			return;
		}

		log.debug('Show tab ' + this.id, this.logAuthor);

		if(!this.displayed) {
			this.getView();
		}
		else {
			var cmps = this.getCmps();

			for(var i = 0; i < cmps.length; i++) {
				if(cmps[i].TabOnShow) {
					cmps[i].TabOnShow();
				}
			}
		}
	},

	_onHide: function() {
		log.debug('Hide tab ' + this.id, this.logAuthor);

		var cmps = this.getCmps();

		for(var i = 0; i < cmps.length; i++) {
			if(cmps[i].TabOnHide) {
				cmps[i].TabOnHide();
			}
		}
	},

	editMode: function() {
		if(!this.edit) {
			this.stopAllTasks();
			this.callParent(arguments);
		}
	},

	cancel: function() {
		this.callParent(arguments);
	},

	startAllTasks: function() {
		log.debug('Start all tasks', this.logAuthor);

		var cmps = this.getCmps();

		for(var i = 0; i < cmps.length; i++) {
			if(cmps[i].startTask) {
				cmps[i].startTask();
			}
		}
	},

	stopAllTasks: function() {
		log.debug('Stop all tasks', this.logAuthor);

		var cmps = this.getCmps();

		for(var i = 0; i < cmps.length; i++) {
			if(cmps[i].stopTask) {
				cmps[i].stopTask();
			}
		}
	},

	//Reporting
	addReportingBar: function() {
		if(!this.report_window) {
			var config = {
				width: 620,
				border: false,
				title: _('Live reporting toolbar'),
				constrain: true,
				renderTo: this.id,
				resizable: false,
				closable: false
			};

			this.reportingBar = Ext.widget('ReportingBar', {
				reloadAfterAction: true
			});

			this.report_window = Ext.widget('window', config);
			this.report_window.addDocked(this.reportingBar);
			this.report_window.show();

			//switch widget to reporting mode
			var cmps = this.getCmps();

			for(var i = 0; i < cmps.length; i++) {
				if(cmps[i].reportMode === false) {
					cmps[i].reportMode = true;
				}
			}

			this.stopAllTasks();
		}
		else {
			log.debug('Reporting bar already opened', this.logAuthor);
		}
	},

	openOptions: function() {
		this.callParent(arguments);

		var options = this.getViewOptions();

		if(!options) {
			options = this.view_options;
		}

		if(this.view_option_win && options) {
			var form = this.view_option_win.down('form');

			if(form) {
				form.getForm().setValues(options);
			}
		}
	},

	setViewOptions: function(data) {
		this.callParent(arguments);

		if(data.orientation && data.pageSize) {
			this.pageModeSize = this.pageWidth[data.orientation][data.pageSize];
		}
	},

	setReportDate: function(from, to) {
		log.debug('Send report data for widgets', this.logAuthor);

		var cmps = this.getCmps();

		for(var i = 0; i < cmps.length; i++) {
			cmps[i]._doRefresh(from, to);
		}
	},

	//misc
	beforeclose: function() {
		log.debug('Active previous tab', this.logAuthor);
		this.onClose = true;

		old_tab = Ext.getCmp('main-tabs').old_tab;

		if(old_tab) {
			Ext.getCmp('main-tabs').setActiveTab(old_tab);
		}


		if(this.localstore_record) {
			//remove from store
			log.debug('Remove this tab from localstore ...', this.logAuthor);

			var store = Ext.data.StoreManager.lookup('Tabs');
			store.remove(this.localstore_record);
			store.save();
		}
	},

 	beforeDestroy: function() {
		log.debug('Destroy items ...', this.logAuthor);
		canopsis.view.Tabs.Content.superclass.beforeDestroy.call(this);
 		log.debug(this.id + ' Destroyed.', this.logAuthor);
 	}
});
