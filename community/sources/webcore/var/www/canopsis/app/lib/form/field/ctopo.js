//need:app/lib/form/cfield.js
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

Ext.define('canopsis.lib.form.field.ctopo' , {
	extend: 'Ext.container.Container',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.ctopology',
	name: 'ctopo',
	autoScroll: true,
	style: {
		borderColor: '#000000',
		borderStyle: 'solid',
		borderWidth: '1px'
	},

	padding: 5,

	logAuthor: '[lib][ctopo]',
	width: '100%',
	height: '100%',

	node_default: {
		nodeBorderWidth: 0,
		nodeBorderColor: '#333',
		nodeBackgroundColor: '#eee',
		nodeColor: '#333',
		cls: '',
		html: ''
	},

	layout: 'absolute',

	cls: 'ctopo-dropzone',
	tpl: [
		'<div id="{id}" class="{cls} {nodeType}" style="border:{nodeBorderWidth}px solid {nodeBorderColor}; background-color:{nodeBackgroundColor}; color:{nodeColor}; padding:1em; border-radius:1em; box-shadow:0.2em 0.2em 0.1em 0.1em #999;">',
			'<center>',
				'<tpl if="event_type == \'operator\'">',
					'<tpl if="connector == \'topology\'">',
						'<img class="ctopo-point" style="width: 32px; height: 32px;" src="/static/canopsis/themes/canopsis/resources/images/topology.png" />',
					'</tpl>',
					'<tpl if="connector != \'topology\'">',
						'<img class="ctopo-point" style="width: 32px; height: 32px;" src="/static/canopsis/themes/canopsis/resources/images/Tango-Blue-Materia/24x24/categories/applications-accessories.png" />',
					'</tpl>',
				'</tpl>',
				'<tpl if="event_type != \'operator\'">',
					'<img class="ctopo-point" style="width: 32px; height: 32px;" src="widgets/stream/logo/{connector}.png" />',
				'</tpl>',
			'</center>',
			'<center>',
				'<span>{label}</span>',
			'</center>',
		'</div>'
	].join('\n'),

	windowOperatorOption: 0,

	getSubmitData: function() {
		return this.getValue();
	},

	getValue: function() {
		var nodes = {};
		var conns = [];

		log.debug('Dump ' + this.items.length + ' nodes:', this.logAuthor);
		if(this.rootNode) {
			for(var i = 0; i < this.items.length; i++) {
				var node = this.getComponent(i);
				node.data['x'] = node.getPosition()[0];
				node.data['y'] = node.getPosition()[1];

				log.debug(' + ' + node._id, this.logAuthor);

				for(var j = 0; j < node.conns.length; j++) {
					var conn = node.conns[j];

					if(conn) {
						var target = Ext.getCmp(conn.targetId);
						var source = Ext.getCmp(conn.sourceId);

						if(target && source) {
							log.debug('   -> ' + source.id, this.logAuthor);
							conns.push([source.id, target.id]);
						}
					}
				}

				//spring cleaning
				if(node.data.conns) {
					delete node.data.conns;
				}

				if(node.data.data) {
					delete node.data.data;
				}

				nodes[node.id] = node.data;
			}


			var r = {
				nodes: nodes,
				conns: conns,
				root: this.rootNode
			};

			log.dump(r);
			return r;
		}
		else {
			log.error('error : no root node', this.logAuthor);
		}
	},

	setValue: function(data) {
		log.debug('Set values', this.logAuthor);

		var nodes = data.nodes;
		var corresp = {};
		var me = this;

		this.removeAll();

		log.debug(' + Create nodes', this.logAuthor);

		Ext.Object.each(nodes, function(key, value) {
			var el = me.createNode(value);
			me.add(el);
			corresp[key] = el.id;
		});

		log.debug('   + Done', this.logAuthor);

		this.rootNode = corresp[data.root];
		log.debug(' + Root Node: ' + this.rootNode, this.logAuthor);

		log.debug(' + Connect nodes', this.logAuthor);

		for(var i = 0; i < data.conns.length; i++) {
			var source = data.conns[i][0];
			var target = data.conns[i][1];
			log.debug('     + ' + source + ' -> ' + target, this.logAuthor);

			this.jsPlumbInstance.connect({
				source: corresp[source],
				target: corresp[target]
			});
		}

		log.debug('   + Done', this.logAuthor);

		this.jsPlumbInstance.repaintEverything();
	},

	initComponent: function() {
		if(!this.topoName) {
			this.topoName = 'topology';
		}

		this.tpl = new Ext.XTemplate(this.tpl, {
			compiled: true,
			operatorNotRootNode: function(event_type, connector) {
				return (event_type === 'operator' && !connector);
			}
		});

		this.callParent(arguments);
	},
	verifyNodeConn: function(connexion) {
		var node_source = Ext.getCmp(connexion.connection.sourceId);

		if(!node_source) {
			return false;
		}

		var connections_out = this.jsPlumbInstance.getConnections({
			source: node_source.id
		});

		var node_target = Ext.getCmp(connexion.targetId);

		if(!node_target) {
			return false;
		}

		var connections_in = this.jsPlumbInstance.getConnections({ source: node_target.id});

		if(connexion.sourceId === connexion.targetId) {
			log.warning('Impossible to self connect', this.logAuthor);
			return false;
		}

		if(node_source.data.nodeMaxOutConnexion !== undefined && connections_out.length >= node_source.data.nodeMaxOutConnexion) {
			log.warning('No OUT slot available', this.logAuthor);
			return false;
		}

		if(node_target.data.nodeMaxInConnexion !== undefined && connections_in.length >= node_target.data.nodeMaxInConnexion) {
			log.warning('No IN slot available', this.logAuthor);
			return false;
		}

		var existingConnexion = this.jsPlumbInstance.getConnections({
			source: node_source.id,
			target: node_target.id
		});

		if(existingConnexion.length > 0) {
			log.debug('Already connected', this.logAuthor);
			return false;
		}

		log.debug('Connect two nodes', this.logAuthor);
		return true;
	},

	setRootNodeName: function(name) {
		var root = Ext.getCmp(this.rootNode);
		var data = {};
		var data = Ext.Object.merge(data, this.node_default);

		data.connector = 'topology';
		data.event_type = 'operator';
		data.label = name;
		data['html'] = this.tpl.apply(data);

		root.data.label = name;
		root.update(data.html);

		this.topologyName = name;
	},

	setRootNodeDescription: function(description) {
		var root = Ext.getCmp(this.rootNode);
		var data = {};
		var data = Ext.Object.merge(data, this.node_default);

		data.connector = 'topology';
		data.event_type = 'operator';
		data.label = this.topologyName + '<br />' + description;
		data['html'] = this.tpl.apply(data);

		root.data.label = data.label;
		root.update(data.html);
		this.topologyDescription = description;
	},

	initJsPlumb: function()  {
		log.debug('initialization of jsPlumb Library', this.logAuthor);
		var me = this;

		this.jsPlumbInstance = jsPlumb.getInstance({
			Container: this.id,
			Endpoint: ['Dot', {radius: 3}],
			HoverPaintStyle: {
				strokeStyle: '#42a62c',
				lineWidth: 2
			},
			PaintStyle: {
				strokeStyle: '#aaa',
				lineWidth: 2
			},
			ConnectionsDetachable: false,
			ConnectionOverlays: [
				['Arrow', {
					location: 1,
					id: 'arrow',
					length: 14,
					foldback: 0.8
				}],
				['Label', { id: 'label' }]
			]
		});

		//to prevent endpoint on right click
		this.jsPlumbInstance.bind('contextmenu', function(a) {
			event.preventDefault();
			var cmp = Ext.getCmp($(a.getElement()).parent().parent().parent().attr('id'));
			me.buildFormOperatorOption(cmp);
			me.jsPlumbInstance.deleteEndpoint(a);
		});

		this.jsPlumbInstance.bind('jsPlumbConnection', function(conn) {
			var target = Ext.getCmp(conn.targetId);
			var source = Ext.getCmp(conn.sourceId);
			target.conns.push(conn);
		});

		this.jsPlumbInstance.bind('beforeDrop', function(conn) {
			return (me.verifyNodeConn(conn));
		});

		this.jsPlumbInstance.bind('click', function(c) {
			log.debug('Remove connection', me.logAuthor);
			me.removeConnection(c);
		});
	},

	removeConnection: function(c ) {
		var target = Ext.getCmp(c.targetId);
		var source = Ext.getCmp(c.sourceId);

		for(var i = 0; i < target.conns.length; i++) {
			if(target.conns[i] && target.conns[i].connection && target.conns[i].connection.id === c.id) {
				delete (target.conns[i]);
			}
		}

		// Remove endpoints
		for(i = 0; i < c.endpoints.length; i++) {
			this.jsPlumbInstance.deleteEndpoint(c.endpoints[i]);
		}

		// Remove connection
		this.jsPlumbInstance.detach(c);
	},

	buildFormOperatorOption: function(nodeEl) {
		var me = this;
		var extForm = nodeEl.extForm;

		if(nodeEl.form) {
			if(extForm && ! extForm.isDestroyed) {
				log.debug('Show form', me.logAuthor);
				extForm.show();
			}
			else {
				log.debug('Create form', me.logAuthor);

				// Translate form
				var form = nodeEl.form;

				if(!form.translated){
					for(var i = 0; i < form.items.length; i++) {
						var item = form.items[i];

						if(item.fieldLabel) {
							item.fieldLabel = _(item.fieldLabel);
						}

						if(item.xtype === 'combobox') {
							for(var j = 0; j < item.store.data.length; j++) {
								if(item.store.data[j].text) {
									item.store.data[j].text = _(item.store.data[j].text);
								}
							}
						}

						form.items[i] = item;
					}

					nodeEl.form.translated = true;
				}

				form.bodyStyle = 'padding: 5px;';
				var form = Ext.create('Ext.form.Panel', form);

				// Load form
				if(nodeEl.options) {
					form.getForm().setValues(nodeEl.options);
				}

				var bbar = [
					{
						xtype: 'button',
						text: _('Cancel'),
						handler: function() {
							log.debug(' + Cancel form', this.logAuthor);
							nodeEl.extForm.close();
						},
						iconCls: 'icon-cancel'
					},'->', {
						xtype: 'button',
						text: _('Save'),
						handler: function() {
							log.debug(' + Save form', this.logAuthor);

							var form = nodeEl.extForm.getComponent(0).getForm();
							var values = form.getValues();

							nodeEl.data.options = values;
							nodeEl.options = values;
							nodeEl.extForm.close();
						},
						iconCls: 'icon-save',
						iconAlign: 'right'
					}
				];

				delete nodeEl.extForm;

				nodeEl.extForm = Ext.create('Ext.window.Window', {
					title: _('Change operator option'),
					items: form,
					bbar: bbar,
					resizable: false,
					closeAction: 'destroy',
					constrain: true,
					renderTo: Ext.getCmp('main-tabs').getActiveTab().id
				}).show().center();
			}
		}
	},

	createNode: function(data) {
		var me = this;
		//we do a deep copy of data, in order to prevent to have a node which is a tree
		var orig_data = {};

		$.extend(true, orig_data, data);
		var data = Ext.Object.merge(data, this.node_default);

		data['data'] = orig_data;
		data['conns'] = [];

		if(data['display_name']) {
			data['label'] = data['display_name'];
		}
		else if(data['source_type'] === 'resource') {
			data['label'] = data['component'] + '<br>' + data['resource'];
		}

		else if(data['source_type'] === 'component') {
			data['label'] = data['component'];
		}

		if(data['event_type'] === 'operator' && data['connector'] === undefined) {
			data['connector'] = undefined;
		}

		data['html'] = this.tpl.apply(data);

		var nodeId = data['_id']
		delete data['id'];

		var nodeEl = Ext.create('Ext.Component', data);
		nodeEl.on('afterRender', function() {
			me.jsPlumbInstance.draggable(nodeEl.id, {
				containment: 'parent'
			});

			var point = $('#' + nodeEl.id + ' .ctopo-point');
			var node = $('#' + nodeEl.id);

			me.jsPlumbInstance.makeSource(point, {
				parent: node,
				anchor: 'Continuous',
				connector: ['StateMachine', { curviness: 20 }]
			});

			me.jsPlumbInstance.makeTarget(node, {
				dropOptions: { hoverClass: 'dragHover' },
				anchor: 'Continuous'
			});

			// Destroy node
			nodeEl.getEl().on('dblclick', function() {
				log.debug('Destroy node ' + nodeEl._id, me.logAuthor);

				if(me.rootNode !== nodeEl.id && !nodeEl.notRemovable) {
					var connections_out = me.jsPlumbInstance.getConnections({
						source: nodeEl.id
					});

					var connections_in = me.jsPlumbInstance.getConnections({
						target: nodeEl.id
					});

					for(var i = 0; i < connections_out.length; i++) {
						me.removeConnection(connections_out[i]);
					}

					for(i = 0; i < connections_in.length; i++) {
						me.removeConnection(connections_in[i]);
					}

					nodeEl.destroy();
				}
				else {
					log.warning('impossible to delete rootNode', this.logAuthor);
				}
			});

			if(data['event_type'] === 'operator') {
				nodeEl.getEl().on('contextmenu', function(e) {
					e.preventDefault();
					me.buildFormOperatorOption(nodeEl);
				});
			}
		});

		return nodeEl;
	},

	afterRender: function() {
		var me = this;

		// JSPlumb
		this.initJsPlumb();

		// Drag and Drop
		log.debug('Init DropZone', this.logAuthor);

		this.dropZone = new Ext.dd.DropZone(this.getEl(), {
			ddGroup: 'search_grid_DNDGroup',

			getTargetFromEvent: function(e) {
				return e.getTarget('.ctopo-dropzone');
			},

			onNodeOver: function() {
				return Ext.dd.DropZone.prototype.dropAllowed;
			},

			onNodeDrop: function(target, dd, e, data) {
				void(target, dd);

				log.debug('Item was dropped', me.logAuthor);

				var form = data.records[0].raw.form;
				var options = data.records[0].raw.options;
				var data = data.records[0].data;

				if(data['event_type'] === 'check') {
					data['nodeMaxInConnexion'] = 0;
				}

				data['x'] = e.getX() - me.getEl().getX() - 20;
				data['y'] = e.getY() - me.getEl().getY() - 20;
				data['form'] = form;
				data['options'] = options;

				var node = me.createNode(data);
				// Add div to container
				me.add(node);
			}
		});

		//Create Root Node
		var root = this.createNode({
			label: 'Topology',
			_id: 'worst_state',
			event_type: 'operator',
			source_type: 'operator',
			nodeMaxOutConnexion: 0,
			connector: 'topology'
		});

		this.add(root);
		this.rootNode = root.id;
	},

	beforeDestroy: function() {
		log.debug('Clean jsPlumb', this.logAuthor);
		this.jsPlumbInstance.deleteEveryEndpoint();
		delete this.jsPlumbInstance;

		this.callParent(arguments);
 	}
});
