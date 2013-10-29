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

unavailableMessageHTML = Ext.create('Ext.XTemplate',
	'<div>',
	'<span style="vertical-align: middle;text-align:center;">{text}</span>',
	'</div>',
	{ compiled: true }
);

Ext.define('widgets.stream.stream', {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.stream',
	logAuthor: '[widget][stream]',

	cls: 'widget-stream',

	max: 10,
	max_comment: 5,

	enable_userinputs: true,
	enable_comments: true,

	tags: '',
	tags_op: true,

	last_push: 0,
	burst_counter: 0,

	//ms
	//nb events
	burst_interval: 500,
	burst_threshold: 2,

	wcontainer_autoScroll: true,
	wcontainer_layout: 'anchor',

	showToolbar: true,

	amqp_queue: 'alerts',
	hard_state_only: true,

	compact: true,

	initComponent: function() {
		if(this.fullscreenMode) {
			this.enable_userinputs = false;
			this.enable_comments = false;
		}

		this.queue = [];
		this.nodeId = false;
		this.refreshInterval = 5;

		if(this.tags !== '') {
			this.tags = split_search_box(this.tags);
		}

		if(!this.showToolbar) {
			this.enable_userinputs = false;
		}

		if(this.showToolbar && !this.exportMode) {
			var items = [];

			if(this.enable_userinputs) {
				items = items.concat([
					{
						xtype: 'tbtext',
						text: "<img src='widgets/stream/logo/ui.png' height='19' width='19'></img>"
					},{
						xtype: 'combobox',
						id: this.id + '-state',
						queryMode: 'local',
						displayField: 'text',
						valueField: 'value',
						width: 70,
						value: 0,
						store: {
							xtype: 'store',
							fields: ['value', 'text'],
							data: [
								{value: 0, text: 'Ok'},
								{value: 1, text: 'Warning'},
								{value: 2, text: 'Critical'}
							]
						}
					},{
						xtype: 'textfield',
						emptyText: _('Leave a') + ' ' + _('event') + ' ?',
						id: this.id + '-message',
						width: 300,
						listeners: {
							specialkey: {
								fn: function(field, e) {
									void(field);

									if(e.getKey() === e.ENTER) {
										this.publish_event();
									}
								},
								scope: this
							}
						}
					}
				]);
			}

			items = items.concat([
				'->', {
						iconCls: 'icon-control-repeat',
						tooltip: _('Clear tray'),
						scope: this,
						handler: function() {
							this.wcontainer.removeAll(true);
						}
					},{
						iconCls: 'icon-control-pause',
						tooltip: _('Pause stream'),
						scope: this,
						enableToggle: true,
						toggleHandler: function(button, state) {
							if(state) {
								button.setIconCls('icon-control-play');
								this.unsubscribe();
							}
							else {
								button.setIconCls('icon-control-pause');
								this.subscribe();
							}
						}
					}
			]);

			this.tbar = Ext.create('Ext.toolbar.Toolbar', {
				items: items
			});
		}

		this.callParent(arguments);
	},

	afterContainerRender: function() {
		if(global.websocketCtrl.connected) {
			this.startStream();
		}
		else {
			this.displayUnavailableMessage();

			global.websocketCtrl.on('transport_up', function() {
				this.wcontainer.removeAll();
				this.startStream();
			}, this, {single: true});
		}
	},

	displayUnavailableMessage: function() {
		this.wcontainer.add({
			xtype:'panel',
			anchor:'100% 100%',
			border: 0,
			layout: {
				align: 'middle',
				pack: 'center',
				type: 'hbox'
			},
			items:[{
				xtype:'panel',
				unstyled: true,
				html:unavailableMessageHTML.apply({text:_('Websocket Unavailable')})
			}]
		});
	},

	startStream: function() {
		var me = this;

		this.getHistory(undefined, undefined, function(records) {
			if(records.length > 0) {
				me.add_events(records);
			}

			if(!me.reportMode) {
				me.subscribe();
			}

			me.ready();
		});
	},

	getHistory: function(from, to, onSuccess) {
		var me = this;

		if(now && global.websocketCtrl.connected) {
			now.stream_getHistory(this.max, this.tags, this.tags_op, from, to, function(records) {
				log.debug('Load ' + records.length + ' events', me.logAuthor);

				if(records.length > 0) {
					for(var i = 0; i < records.length; i++) {
						records[i] = Ext.create('widgets.stream.event', {
							id: me.get_event_id(records[i]),
							raw: records[i],
							stream: me
						});
					}

					if(onSuccess) {
						onSuccess(records);
					}
				}
			});
		}
		else {
			log.error("'now' is undefined, websocket down ?", me.logAuthor);
		}
	},

	subscribe: function() {
		// Subscribe to AMQP channel
		global.websocketCtrl.subscribe('amqp', this.amqp_queue, this.on_event, this);
	},

	unsubscribe: function() {
		// Unsubscribe
		global.websocketCtrl.unsubscribe('amqp', this.amqp_queue, this);
	},

	publish_event: function() {
		if(!global.websocketCtrl.connected) {
			log.error('Impossible to publish, not connected.', this.logAuthor);
			global.notify.notify(_('Error'), _('Impossible to publish, your are not connected to websocket. Check service or firewall') + ' (port: ' + global.nowjs.port + ')', 'error');
			return;
		}

		var toolbar = 0;

		if(this.title) {
			toolbar = this.getDockedItems()[1];
		}
		else {
			toolbar = this.getDockedItems()[0];
		}

		var message = toolbar.getComponent(this.id + '-message').getValue();
		toolbar.getComponent(this.id + '-message').reset();

		var state = toolbar.getComponent(this.id + '-state').getValue();

		var event_raw = {
			'connector_name': 'widget-stream',
			'source_type': 'component',
			'event_type': 'user',
			'component': global.account.id,
			'output': message,
			'display_name': global.account.firstname + ' ' + global.account.lastname,
			'author': global.account.firstname + ' ' + global.account.lastname,
			'state': state,
			'state_type': 1,
			'tags': this.tags
		};

		global.websocketCtrl.publish('amqp', 'events', event_raw);
	},

	publish_comment: function(event_id, raw, message) {
		if(!global.websocketCtrl.connected) {
			log.error('Impossible to publish, not connected.', this.logAuthor);
			global.notify.notify(_('Error'), _('Impossible to publish, your are not connected to websocket. Check service or firewall') + ' (port: ' + global.nowjs.port + ')', 'error');
			return;
		}

		log.debug(event_id + ' -> ' + message, this.logAuthor);

		var event_raw = {
			'connector_name': 'widget-stream',
			'source_type': raw.source_type,
			'event_type': 'comment',
			'component': raw.component,
			'resource': raw.resource,
			'output': message,
			'referer': event_id,
			'author': global.account.firstname + ' ' + global.account.lastname,
			'state': 0,
			'state_type': 1,
			'tags': raw.tags
		};

		global.websocketCtrl.publish('amqp', 'events', event_raw);
	},

	doRefresh: function(from, to) {
		if(this.reportMode) {
			this.unsubscribe();
			this.purge_queue();
			this.wcontainer.removeAll(true);

			var me = this;
			this.getHistory(parseInt(from/1000), parseInt(to/1000), function(records) {
				if(records.length > 0) {
					me.add_events(records);
				}
			});
		}
		else {
			this.process_queue();

			//refresh time
			for(var i = 0; i < this.wcontainer.items.length; i++) {
				var event = this.wcontainer.getComponent(i);

				if(event) {
					event.update_time();
				}
			}
		}
	},

	TabOnShow: function() {
		this.doLayout();
		this.purge_queue();
		this.callParent();
	},

	process_queue: function() {
		// Check burst
		if(!this.in_burst()) {
			this.purge_queue();
		}
	},

	purge_queue: function() {
		if(this.queue.length) {
			log.debug("Purge event's queue (" + this.queue.length + ')', this.logAuthor);
			// Back to normal, purge queue
			this.add_events(this.queue);
			this.queue = [];
		}
	},

	in_burst: function() {
		if((this.last_push + this.burst_interval) > new Date().getTime()) {
			if (this.burst_counter < this.burst_threshold) {
				this.burst_counter += 1;
				log.debug('Burst counter: ' + this.burst_counter, this.logAuthor);
				return false;
			}
			else {
				return true;
			}
		}
		else {
			this.burst_counter = 0;
			return false;
		}
	},

	get_event_id: function(raw) {
		var id = undefined;

		if(raw['_id']) {
			id = raw['_id'];
		}

		return id;
	},

	on_event: function(raw) {
		//Only hard state
		if(raw.state_type === 0 && this.hard_state_only) {
			return;
		}

		var i;

		// Check tags
		if(this.tags && raw.tags) {
			if(this.tags_op) {
				// AND
				for(i = 0; i < this.tags.length; i++) {
					if(!Ext.Array.contains(raw.tags, this.tags[i])) {
						return;
					}
				}
			}
			else {
				// OR
				var show = false;

				for (i = 0; i < this.tags.length; i++) {
					if(Ext.Array.contains(raw.tags, this.tags[i])) {
						show = true;
					}
				}

				if(!show) {
					return;
				}
			}
		}

		var id = this.get_event_id(raw);

		var event = Ext.create('widgets.stream.event', {
			id: id,
			raw: raw,
			stream: this
		});

		if(event.raw.event_type === 'comment') {
			var to_event = this.wcontainer.getComponent(this.id + '.' + event.raw.referer);

			if(to_event) {
				log.debug('Add comment for ' + event.raw.referer, this.logAuthor);
				to_event.comment(event);
			}
			else {
				log.debug("Impossible to find event '" + event.raw.referer + "' from container, maybe not displayed ?", this.logAuthor);
			}
		}
		else {
			// Detect Burst or hidden
			if(this.in_burst() || this.isHidden()) {
				this.queue.push(event);

				//Clean queue
				if(this.queue.length > this.max) {
					event = this.queue.shift();
					event.destroy();
					delete event;
				}
			}
			else {
				//Display event
				this.process_queue();
				this.add_events([event]);
			}

			this.last_push = new Date().getTime();
		}
	},

	add_events: function(events) {
		if(events.length >= this.max) {
			this.wcontainer.removeAll(true);
		}

		this.wcontainer.insert(0, events);

		//Remove last components
		while(this.wcontainer.items.length > this.max) {
			var item = this.wcontainer.getComponent(this.wcontainer.items.length - 1);
			this.wcontainer.remove(item.id, true);
		}
	},

 	beforeDestroy: function() {
		this.unsubscribe();
		this.wcontainer.removeAll(true);

		this.callParent(arguments);
 	}
});
