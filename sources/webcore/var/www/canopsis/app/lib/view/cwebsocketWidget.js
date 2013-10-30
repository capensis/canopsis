/*
#--------------------------------
# Copyright (c) 2013 "Capensis" [http://www.capensis.com]
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
Ext.define('canopsis.lib.view.cwebsocketWidget', {
	extend: 'canopsis.lib.view.cwidget',

    /**
     * The amqp queue the widget listens on
     * @protected
     */
	amqp_queue: 'alerts',

	initComponent: function () {
		this.refreshInterval = 5;

		this.callParent(arguments);
	},

	afterContainerRender: function() {
		if (global.websocketCtrl.connected) {
			this.startStream();
		}else {
			this.displayUnavailableMessage()
			global.websocketCtrl.on('transport_up', function() {
				this.wcontainer.removeAll();
				this.startStream();
			}, this, {single: true});
		}
	},

	startStream: function() {
		var me = this;
		this.getHistory(undefined, undefined, function(records) {
			if (records.length > 0)
				me.add_events(records);

			if (! me.reportMode)
				me.subscribe();

			me.ready();
		});
	},

    /**
     * Retreives the history
     * @param timestamp from
     * @param timestamp to
     * @param function onSuccess
     * @protected
     */
	getHistory: function(from, to, onSuccess) {
	},

	/**
     * Add events to the widget
     * @param events the list of events to add
     * @protected
     */
	add_events: function(events) {
	},

    /**
     * Subscribes to an AMQP queue
	 * @uses this.amqp_queue
     * @protected
     */
	subscribe: function() {
		// Subscribe to AMQP channel
		global.websocketCtrl.subscribe('amqp', this.amqp_queue, this.on_event, this);
	},

    /**
     * Unsubscribes to an AMQP queue
	 * @uses this.amqp_queue
     * @protected
     */
	unsubscribe: function() {
		global.websocketCtrl.unsubscribe('amqp', this.amqp_queue, this);
	},

	process_queue: function() {
		// Check burst
		if (! this.in_burst())
			this.purge_queue();
	},

	purge_queue: function() {
		if (this.queue.length) {
			log.debug("Purge event's queue (" + this.queue.length + ')', this.logAuthor);
			// Back to normal, purge queue
			this.add_events(this.queue);
			this.queue = [];
		}
	},

 	beforeDestroy: function() {
		this.unsubscribe();
		this.wcontainer.removeAll(true);

		this.callParent(arguments);
 	},

    /**
     * Publish an event on the amqp queue
	 * @param queue the destination of the amqp message
	 * @param event the content of the message
     * @protected
     */
 	publishEvent: function(queue, event) {
		global.websocketCtrl.publish('amqp', queue, event);
 	},

    /**
     * Displays a message at the center of the widget when the websocket is not availlable
     */
	displayUnavailableMessage: function(){
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
		})
	},

    /**
     * check if publish is possible, otherwise displays a popup to notify the user
     * @returns bool public possiblity
     */
	checkPublishPossible: function(){
		if (! global.websocketCtrl.connected) {
			log.error('Impossible to publish, not connected.', this.logAuthor);
			global.notify.notify(_('Error'), _('Impossible to publish, your are not connected to websocket. Check service or firewall') + ' (port: ' + global.nowjs.port + ')', 'error');
			return false;
		}
		else
			return true;
	},

	/**
     * event triggered when an event comes from the amqp queue
     * @param raw event raw JSON
     * @param rk the routing key of the event
     */

	on_event: function(raw, rk) {

		var id = this.get_event_id(raw);

		var event = Ext.create('widgets.stream.event', {id: id, raw: raw, stream: this});

		if (event.raw.event_type == 'comment') {
			var to_event = this.wcontainer.getComponent(this.id + '.' + event.raw.referer);
			if (to_event) {
				log.debug('Add comment for ' + event.raw.referer, this.logAuthor);
				to_event.comment(event);
			}else {
				log.debug("Impossible to find event '" + event.raw.referer + "' from container, maybe not displayed ?", this.logAuthor);
			}

		}else {
			// Detect Burst or hidden
			if (this.in_burst() || this.isHidden()) {
				this.queue.push(event);

				//Clean queue
				if (this.queue.length > this.max) {
					var event = this.queue.shift();
					event.destroy();
					delete event;
				}
			}else {
				//Display event
				this.process_queue();
				this.add_events([event]);
			}

			this.last_push = new Date().getTime();
		}
	},

	/**
	 * Get the id from a raw event
	 * @param raw the raw event (json)
	 * @returns int the id of the event
	 */
	get_event_id: function(raw) {
		var id = undefined;
		if (raw['_id'])
			id = raw['_id'];

		return id;
	}
});
