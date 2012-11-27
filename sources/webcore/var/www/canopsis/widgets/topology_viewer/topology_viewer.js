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
Ext.define('widgets.topology_viewer.topology_viewer' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.topology_viewer',

	//geometric options

	FirstChildAngle: Math.PI / 2,
	LastChildAngle: 3 * Math.PI / 2,
	radiusOffset: 0.4,

	usable_angle: Math.PI / 2,

	//sigma options
	defaultLabelColor: '#ccc',
	sigmaFont: 'Arial',
	edgeColor: 'source',
	defaultEdgeType: 'line',

	minNodeSize: 6,
	maxNodeSize: 8,

    minEdgeSize: 2,
    maxEdgeSize: 4,

	labelThreshold: 7,
	labelBackground: true,

	firstPointColor: '#ff0000',
	secondPointColor: '#00ff00',

	baseUrl: '/rest/events',
	background_color: undefined,

	state_color: ['green', 'orange', 'red', 'grey'],

	lastUpdate: undefined,
	canvasContext: undefined,
	canvas: undefined,

	//-----------------widget functions----------------

	initComponent: function() {
		if (this.background_color)
			this.bodyStyle = {'background-color': this.background_color};



		this.callParent(arguments);
	},

	doRefresh: function() {
		if (this.sigmaContainer)
			this.sigmaContainer.empty();

		this.initSigma();
	},

	//-------------------Sigma related functions------------------

	initSigma: function() {
		log.debug('Init Sigma.js');
		var sigma_root = this.wcontainer.getEl().id;
		this.sigmaContainer = sigma.init(document.getElementById(sigma_root));
		this.sigmaContainer.drawingProperties({
		  defaultLabelColor: this.defaultLabelColor,
		  labelThreshold: this.labelThreshold,
		  font: this.sigmaFont,
		  edgeColor: this.edgeColor,
		  defaultEdgeType: this.defaultEdgeType,
		  labelHoverShadow: this.labelBackground
		}).graphProperties({
		  minNodeSize: this.minNodeSize,
			maxNodeSize: this.maxNodeSize,
			minEdgeSize: this.minEdgeSize,
			maxEdgeSize: this.maxEdgeSize
		}).mouseProperties({
			maxRatio: 4
		});

		if (this.nodes) {
			var id = this.nodes[0];
			Ext.Ajax.request({
				url: this.baseUrl + '/' + id,
				scope: this,
				method: 'GET',
				params: {ids: Ext.encode(this.nodeId)},
				success: function(response) {
					var nodes = Ext.JSON.decode(response.responseText).data;
					var nestedTree = nodes[0]['nestedTree'];
					this.lastUpdate = nodes[0]['crecord_creation_time'];
					this.drawRecursiveTree(nestedTree);
					this.sigmaDraw();
				}
			});
		}
	},

	computeAnglesPosition: function(number_of_point,usable_angle,start_angle) {
		//center to the middle of the angle
		start_angle = start_angle - (usable_angle / 2);

		//tweak, if 2 point, use first and last point angle
		if (number_of_point == 2 && usable_angle != 2 * Math.PI)
			return [start_angle, start_angle + usable_angle];

		//otherwise compute offset step for each point
		var offset = usable_angle / number_of_point;
		var tab = [];
		for (var i = 0; i < number_of_point; i++)
			tab.push(offset * i + start_angle);

		return tab;
	},

	drawRecursiveTree: function(tree,depth,angle,referent_coord) {
		//console.log('Recursive tree '+tree.name+' ' + depth +' '+ angle)
		if (!depth) depth = 0;
		if (!angle) angle = Math.PI * 2;
		if (!referent_coord) referent_coord = {x: 0, y: 0};

		var depth_coef = (depth / 100);

		var radius = this.radiusOffset - depth_coef;

		var coord = {
			x: this.getX(referent_coord.x, radius, angle),
			y: this.getY(referent_coord.y, radius, angle)
		};

		var node_params = {
			label: tree.name,
			x: coord.x,
			y: coord.y,
			shape: 'square'
		};


		if (depth == 0) {
			node_params.shape = 'square',
			node_params.size = this.maxNodeSize * 1.5;
		}

		if (tree.state != undefined)
			node_params.color = this.state_color[tree.state];

		this.addNode(tree._id, node_params);

		if (depth == 0)
			var usuable_angle = 2 * Math.PI;
		else
			var usuable_angle = this.usable_angle - (depth_coef * 50);

		if (tree.childs && tree.childs.length > 0) {
			angle = this.computeAnglesPosition(tree.childs.length, usuable_angle, angle);
			for (var i in tree.childs) {
				var child = tree.childs[i];
				var _id = this.drawRecursiveTree(child, depth + 1, angle[i], coord);

				var params = {
				};

				if (child.state != undefined)
					params.color = this.state_color[child.state];

				this.linkNode(_id + '-' + tree._id, _id, tree._id, params);
			}
		}

		return tree._id;
	},

	addNode: function(name,config) {
		this.sigmaContainer.addNode(name, config);
	},

	linkNode: function(link_name,first_node,second_node,params) {
		this.sigmaContainer.addEdge(link_name, first_node, second_node, params);
	},

	displayLastUpdate: function() {
		if (this.canvasContext && this.lastUpdate) {
			//console.log(this.canvas.width)
			//console.log(this.canvas.height)
			this.canvasContext.fillText(rdr_elapsed_time(this.lastUpdate), 10, 10);
		}
	},

	sigmaDraw: function() {
		this.sigmaContainer.draw();
		this.canvas = document.getElementById(this.sigmaContainer._core.domRoot.lastChild.id);
		this.canvasContext = this.canvas.getContext('2d');
		this.canvasContext.font = 'bold 12px sans-serif';
		this.displayLastUpdate();

		this.sigmaContainer._core.mousecaptor.bind(
							'stopinterpolate',
							this.displayLastUpdate.bind(this));
	},

	//----------------trigo functions-----------------
	degTorad: function(val) {
		return val * (Math.PI / 180);
	},

	radTodeg: function(val) {
		return val * (180 / Math.PI);
	},

	getXY: function(x,y,radius,angle) {
		var x = this.getX(x, radius, angle);
		var y = this.getY(y, radius, angle);
		return {x: x, y: y};
	},

	getX: function(x,radius,angle) {
		return x + radius * Math.cos(angle);
	},

	getY: function(y,radius,angle) {
		return y + radius * Math.sin(angle);
	}

});
