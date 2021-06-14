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
Ext.define('canopsis.controller.Websocket', {
	extend: 'Ext.app.Controller',

	views: [],
	stores: [],

	logAuthor: '[controller][Websocket]',

	autoconnect: true,
	connected: false,

	subscribe_cache: {},
	auto_resubscribe: true,

	jsLoader: undefined,
	jsLoaded: false,

	init: function() {
		global.websocketCtrl = this;
		this.loadLibs();
	},

	loadLibs: function(){
		if(this.jsLoader) {
			this.jsLoader.destroy();
		}

		this.jsLoader = new Ext.Element(document.createElement('script'));

		if(Ext.isIE) {
			Ext.defer(function() {
				if(typeof(now) !== 'undefined') {
					this.jsLoaded = true;
					this.bindNow();

					if(this.autoconnect) {
						this.connect();
					}
				}
			}, 2000, this);
		}
		else {
			this.jsLoader.on('load', function() {
				if(typeof(now) !== 'undefined') {
					this.jsLoaded = true;
					this.bindNow();

					if(this.autoconnect) {
						this.connect();
					}
				}
			}, this, {single: true});
		}

		this.jsLoader.set({
			type: "text/javascript",
			src: global.nowjs.proto + '://' + global.nowjs.hostname + ':' + global.nowjs.port + '/nowjs/now.js'
		});

		document.getElementById('nowjs').appendChild(this.jsLoader.dom);
	},

	bindNow: function() {
		var me = this;

		now.authToken = global.account.authkey;
		now.authId = global.account._id;

		now.ready(function() {
			log.debug(' + Now ready', me.logAuthor);

			now.core.socketio.on('disconnect', function() {
				if(me.connected) {
					me.connected = false;
					me.transport_down();
					me.fireEvent('transport_down', this);
				}
			});

			now.auth(function() {
				log.debug(' + Authed', me.logAuthor);
				if(!me.connected) {
					me.connected = true;
					me.transport_up();
					me.fireEvent('transport_up', this);
				}
			});
		});
	},

	connect: function() {
		log.debug('Connect Websocket ...', this.logAuthor);

		if(this.connected) {
			log.debug(' + Already connected', this.logAuthor);
			return;
		}

		if(!this.jsLoaded) {
			log.error('NowJS Client not loaded. Try to load it.', this.logAuthor);
			this.loadLibs();
			return;
		}
		else {
			log.debug(' + All is Ok, NowJS do reconnect automatically', this.logAuthor);
		}
	},

	transport_down: function() {
		log.info('Transport Down', this.logAuthor);

		if(global.notify) {
			global.notify.notify(_('Info'), _('Disconnected from websocket.'), 'info');
		}
	},

	transport_up: function() {
		log.info('Transport Up', this.logAuthor);

		if(global.notify) {
			global.notify.notify(_('Success'), _('Connected to websocket'), 'success');
		}

		//Re-open channel
		if(this.subscribe_cache && this.auto_resubscribe) {
			for(var i = 0; i < this.subscribe_cache.length; i++) {
				var s = this.subscribe_cache[i];
				delete this.subscribe_cache[i];

				for(var j = 0; j < s.subscribers.length; j++) {
					var t = s.subscribers[j];
					this.subscribe(s.type, s.channel, t.on_message, t.scope);
				}
			}
		}
	},

	subscribe: function(type, channel, on_message, scope) {
		if (this.connected) {
			if (!scope) {
				scope = this;
			}

			log.info(' + Subscribe to ' + type + '.' + channel + ' (' + scope.id + ')', this.logAuthor);

			id = type + '-' + channel;

			// Open one channel by id and distribute messages
			if(!this.subscribe_cache[id]) {
				this.subscribe_cache[id] = {
					type: type,
					channel: channel,
					subscribers: {}
				};

				this.subscribe_cache[id].subscribers[scope.id] = {
					on_message: on_message,
					scope: scope
				};

				var me = this;
				var callback = function(message, rk) {
					Ext.Object.each(me.subscribe_cache[id].subscribers, function(key, subscriber) {
						void(key);

						subscriber.on_message.apply(subscriber.scope, [message, rk]);
					}, me);
				};

				//Register callback
				now[id] = callback;

				//subscribe to group
				now.subscribe(type, channel);
			}
			else {
				this.subscribe_cache[id].subscribers[scope.id] = {
					on_message: on_message,
					scope: scope
				};
			}
		}
	},

	unsubscribe: function(type, channel, scope) {
		if(this.connected) {
			if(!scope) {
				scope = this;
			}

			log.info(' + Unsubscribe to ' + type + '.' + channel + ' (' + scope.id + ')', this.logAuthor);

			id = type + '-' + channel;

			if(this.subscribe_cache[id]) {
				delete this.subscribe_cache[id].subscribers[scope.id];

				if(isEmpty(this.subscribe_cache[id].subscribers)) {
					log.info("  + Delete cache '" + id + "' and unsubscribe from remote queue", this.logAuthor);
					delete this.subscribe_cache[id];

					// Unsubscribe from group
					now.unsubscribe(type, channel);

					//Delete callback
					delete now[id];
				}
			}
			else {
				log.error("  + Invalid queue id '" + id + "'", this.logAuthor);
			}
		}
	},

	publish_event: function() {
	},

	publish: function(type, channel, message) {
		now.publish(type, channel, message);
	},

	on_event: function() {
	},

	on_pv: function() {
	}
});
