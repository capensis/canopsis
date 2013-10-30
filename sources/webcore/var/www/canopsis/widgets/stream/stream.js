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

unavailableMessageHTML = Ext.create('Ext.XTemplate',
	'<div>',
    	'<span style="vertical-align: middle;text-align:center;">{text}</span>',
	'</div>',
	{ compiled: true }
);

Ext.define('widgets.stream.stream' , {
	extend: 'canopsis.lib.view.cwebsocketWidget',

	alias: 'widget.stream',

	cls: 'widget-stream',

	max: 10,
	max_comment: 5,

	enable_userinputs: true,
	enable_comments: true,

	tags: '',
	tags_op: true,

	last_push: 0,
	burst_counter: 0,

	burst_interval: 500, //ms
	burst_threshold: 2, //nb events

	wcontainer_autoScroll: true,
	wcontainer_layout: 'anchor',

	showToolbar: true,

	hard_state_only: true,

	compact: true,

	initComponent: function() {
		if (this.fullscreenMode) {
			this.enable_userinputs = false;
			this.enable_comments = false;
		}

		this.queue = [];
		this.nodeId = false;

		if (this.tags != '') {
			this.tags = split_search_box(this.tags);
		}

		if (! this.showToolbar)
			this.enable_userinputs = false;

		if (this.showToolbar && ! this.exportMode) {

			var items = [];

			if (this.enable_userinputs) {
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
									if (e.getKey() == e.ENTER)
										this.publish_event();
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
							if (state) {
								button.setIconCls('icon-control-play');
								this.unsubscribe();
							}else {
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
		logAuthor = '[widget][stream]';
	},

	getHistory: function(from, to, onSuccess) {
		var me = this;
		if (now && global.websocketCtrl.connected) {
			now.stream_getHistory(this.max, this.tags, this.tags_op, from, to, function(records) {
				log.debug('Load ' + records.length + ' events', me.logAuthor);
				if (records.length > 0) {
					for (var i = 0; i < records.length; i++)
							records[i] = Ext.create('widgets.stream.event', {id: me.get_event_id(records[i]), raw: records[i], stream: me});

					if (onSuccess)
						onSuccess(records);
				}
			});
		}else {
				log.error("'now' is undefined, websocket down ?", me.logAuthor);
		}
	},

	publish_event: function() {
		if(this.checkPublishPossible())
		{
			var toolbar = 0;

			this.title ? toolbar = this.getDockedItems()[1] : toolbar = this.getDockedItems()[0];

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

			this.publishEvent('events', event_raw);
		}
	},

	publish_comment: function(event_id, raw, message, orievent) {
		if(this.checkPublishPossible())
		{
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

			this.publishEvent('events', event_raw);
		}
	},

	TabOnShow: function() {
		this.doLayout();
		this.purge_queue();
		this.callParent();
	},

	in_burst: function() {
		if ((this.last_push + this.burst_interval) > new Date().getTime()) {
			if (this.burst_counter < this.burst_threshold) {
				this.burst_counter += 1;
				log.debug('Burst counter: ' + this.burst_counter, this.logAuthor);
				return false;
			}else {
				return true;
			}
		}else {
			this.burst_counter = 0;
			return false;
		}
	},

	add_events: function(events) {
		if (events.length >= this.max)
			this.wcontainer.removeAll(true);

		this.wcontainer.insert(0, events);

		//Remove last components
		while (this.wcontainer.items.length > this.max) {
			var item = this.wcontainer.getComponent(this.wcontainer.items.length - 1);
			this.wcontainer.remove(item.id, true);
		}
	}

});
