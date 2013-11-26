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
Ext.define('canopsis.lib.view.cwidget', {
	extend: 'Ext.panel.Panel',

	border: false,
	layout: 'fit',

	nodeId_refresh: true,
	nodeData: {},
	nodes: [],

	data: {},

	defaultHtml: '<center><span class="icon icon-loading" /></center>',

	refreshInterval: 0,

	baseUrl: '/rest/events/event',
	uri: '/rest/events/event',

	wcontainer_layout: 'fit',
	wcontainer_autoScroll: false,

	task: false,

	reportMode: false,
	exportMode: false,
	fullscreenMode: false,

	autoTitle: false,
	humanReadable: true,

	barHeight: 27,

	time_window: global.commonTs.day,

	time_window_offset: 0,

	useLastRefresh: true,

	lastRefresh: undefined,

	active: false,

	initComponent: function() {
		this.callParent(arguments);

		this.active = true;

		if(!this.logAuthor) {
			this.logAuthor = '[widgets][' + this.xtype + ']';
		}

		log.debug('Initialize component:', this.logAuthor);
		log.dump({
			id: this.id,
			type: this.xtype,
			reportMode: this.reportMode,
			exportMode: this.exportMode
		});

		if(this.title === '') {
			this.title = false;
		}

		this.wcontainerId = this.id + '-content';

		this.wcontainer = Ext.create('Ext.container.Container', {
			id: this.wcontainerId,
			border: false,
			layout: this.wcontainer_layout,
			autoScroll: this.wcontainer_autoScroll
		});

		this.items.add(this.wcontainer);

		this.wcontainer.on('afterrender', function() {
			log.debug('SetHeight of wcontainer', this.logAuthor);
			this.wcontainer.setHeight(this.getHeight());

			this.afterContainerRender();
		}, this);

		if(this.reportMode) {
			this.refreshInterval = false;
		}

		//Compatibility
		if(this.nodes && this.nodes.length > 0) {
			log.debug('Nodes:', this.logAuthor);
			log.dump(this.nodes);
			this.nodeId = this.nodes;
		}

		if(this.inventory) {
			this.nodeId = this.inventory;
		}

		if(Ext.isArray(this.nodes)) {
			this.nodesByID = parseNodes(this.nodes);
		}
		else {
			this.nodesByID = expandAttributs(this.nodes);
		}

		//if reporting
		if(!this.exportMode) {
			if(this.refreshInterval) {
				log.debug(' + Refresh Interval: ' + this.refreshInterval, this.logAuthor);

				this.task = {
					run: this._doRefresh,
					interval: this.refreshInterval * 1000,
					scope: this,
					args: [undefined,undefined],
					active: false
				};
			}
		}
	},

	afterContainerRender: function() {
		log.debug(' + Ready', this.logAuthor);
		this.ready();
	},

	getHeight: function() {
		var height = this.callParent();

		var docks = this.getDockedItems();

		if(docks) {
			height -= docks.length * 2;

			for(var i = 0; i < docks.length; i++) {
				if (docks[i].dock === 'top' || docks[i].dock === 'bottom') {
					height -= this.barHeight;
				}
			}
		}

		if(this.border) {
			height -= this.border * 2;
		}

		return height;
	},

	ready: function() {
		if(this.task) {
			this.startTask();
		}
		else {
			if(this.exportMode) {
				this._doRefresh(this.export_from, this.export_to);
			}
			else {
				this._doRefresh(undefined, parseInt(Ext.Date.now()));
			}
		}
	},

	startTask: function() {
		if(!this.reportMode) {
			if(this.task && ! this.task.active) {
				log.debug('Start task, interval:  ' + this.refreshInterval + ' seconds', this.logAuthor);
				Ext.TaskManager.start(this.task);
				this.task.active = true;
			}
			else {
				this._doRefresh(undefined, undefined);
			}
		}
	},

	stopTask: function() {
		if(this.task && this.task.active) {
			log.debug('Stop task', this.logAuthor);
			Ext.TaskManager.stop(this.task);
			this.task.active = false;
		}
	},


	TabOnShow: function() {
		log.debug('Show', this.logAuthor);
		this.active = true;

		if(!this.isDisabled()) {
			this.startTask();
		}
	},

	TabOnHide: function() {
		log.debug('Hide', this.logAuthor);
		this.active = false;

		this.stopTask();
	},

	_doRefresh: function(from, to) {
		var now = parseInt(Ext.Date.now());

		if(!to) {
			to = now - (this.time_window_offset * 1000);
		}

		if(!from && this.useLastRefresh && this.lastRefresh) {
			from = this.lastRefresh;
		}

		if(!from) {
			from = to - (this.time_window * 1000);
		}

		this.doRefresh(from, to);
		this.lastRefresh = to;
	},

	doRefresh: function(from, to) {
		this.getNodeInfo(from, to);
	},

	_onRefresh: function(data, from, to) {
		this.data = data;
		this.onRefresh(data, from, to);
	},

	onRefresh: function() {
		log.debug('onRefresh', this.logAuthor);
	},

	onResize: function() {
		log.debug('onRezize', this.logAuthor);
	},

	getNodeInfo: function(from, to) {
		if(this.nodeId) {

			var nodeInfoParams = this.getNodeInfoParams(from, to);

			Ext.Ajax.request({
				url: this.baseUrl + '/events' + (this.nodeId && this.nodeId.length? ('/' + this.nodeId) : ''),
				scope: this,
				params: nodeInfoParams,
				method: 'GET',
				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);

					if(this.nodeId.length > 1) {
						data = data.data;
					}
					else {
						data = data.data[0];
					}

					this._onRefresh(data, from, to);
				},
				failure: function(result, request) {
					void(result);

					log.error('Impossible to get Node informations, Ajax request failed ... (' + request.url + ')', this.logAuthor);
				}
			});
		}
	},

	getNodeInfoParams: function(from, to) {
		void(from);
		void(to);
		return {};
	},

	setHtml: function(html) {
		log.debug('setHtml in widget', this.logAuthor);

		this.wcontainer.removeAll();

		this.wcontainer.add({
			html: html,
			border: false
		});

		this.wcontainer.doLayout();
	},

	setHtmlTpl: function(tpl, data) {
		log.debug('setHtmlTpl in div ' + this.wcontainerId, this.logAuthor);
		tpl.overwrite(this.wcontainerId, data);
	},

	beforeDestroy: function() {
		log.debug('Destroy ...', this.logAuthor);
		this.stopTask();
		this.callParent(arguments);
 	}
});
