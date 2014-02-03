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

Ext.define('canopsis.view.SLA.BarWidget', {
	extend: 'Ext.container.Container',
	alias: 'widget.slabarwidget',

	logAuthor: '[SLA][Bar]',

	crit: undefined,
	delay: undefined,

	from: function() {
		return undefined;
	},

	to: function() {
		return undefined;
	},

	initComponent: function() {
		this.callParent(arguments);

		this.nodes = {};

		this.refreshNodes();
	},

	refreshNodes: function() {
		Ext.Ajax.request({
			url: '/entities/',
			scope: this,
			method: 'POST',

			jsonData: {
				filter: {
					'$or': [
						{'type': {'$in': ['component', 'resource', 'ack', 'downtime']}},
						{'type': 'metric', 'name': 'cps_state'}
					]
				}
			},

			success: function(response) {
				var data = Ext.JSON.decode(response.responseText);

				var payload = {};
				var key;

				/* parse response from server */
				for(var i = 0; i < data.data.length; i++) {
					var entity = data.data[i];

					if(entity.type === 'component' || entity.type === 'resource') {
						key = entity.type + ':' + entity.name;

						payload[key] = entity;
					}
					else if(entity.type === 'ack' || entity.type === 'downtime') {
						if(entity.resource) {
							key = 'resource:' + entity.resource;
						}
						else {
							key = 'component:' + entity.component;
						}

						if(!payload[key][entity.type]) {
							payload[key][entity.type] = [];
						}

						payload[key][entity.type].push(entity);
					}
					else if(entity.type === 'metric') {
						this.nodes[entity.nodeid] = {
							id: entity.nodeid,
							from: this.from(),
							to: this.to(),
							co: entity.component,
							re: entity.resource,
							me: entity.name
						};
					}
				}

				/* save payload */
				this.payload = [];

				for(key in payload) {
					this.payload.push(payload[key]);
				}

				/* ask nodes */
				var nodes = [];

				for(var nodeid in this.nodes) {
					nodes.push(this.nodes[nodeid]);
				}

				Ext.Ajax.request({
					url: '/perfstore/values',
					scope: this,
					method: 'POST',

					params: {
						'nodes': nodes
					},

					success: function(response) {
						var data = Ext.JSON.decode(response.responseText);

						for(var i = 0; i < data.data.length; i++) {
							var info = data.data[i];
							var node = this.nodes[info.node];
							var entity;

							if(node.re) {
								entity = payload['resource:' + node.re];
							}
							else {
								entity = payload['component:' + node.co];
							}

							entity.states = info.values;

							this.onRefresh(this.payload);
						}
					}
				});
			}
		});
	},

	onRefresh: function(data) {
	},

	getState: function(cps_state) {
		var state = '' + cps_state;

		while(state.length < 3) {
			state = '0' + state;
		}

		return state;
	}
});