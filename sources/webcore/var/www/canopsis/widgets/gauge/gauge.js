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
Ext.define('widgets.gauge.gauge' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.gauge',

	logAuthor: '[gauge]',

	lines: 12, // The number of lines to draw
	angle: 4, // The length of each line
	lineWidth: 0.3, // The line thickness
	pointer: {
		length: 0.9, // The radius of the inner circle
		strokeWidth: 0.035, // The rotation offset
		color: '#000000' // Fill color
	},
	colorStart: '#6FADCF',   // Colors
	colorStop: '#8FC0DA',    // just experiment with them
	strokeColor: '#EEEEEE',   // to see which ones work best for you
	generateGradient: true,
	textSize : 40,

	maxValue:100,
	animationSpeed:32,

	afterContainerRender: function() {
		this.callParent(arguments);

		var width = this.wcontainer.getWidth()
		var height = this.getHeight() - (this.textSize + 10)
		var canvasId = this.wcontainerId + '-canvas'
		var textId = this.wcontainerId + '-text'

		var textHTML = '<div id="'+textId+'" style="font-size: '+this.textSize+'px;text-align:center;"></div>'
		var canvasHTML = '<canvas width="'+width+'" height="'+height+'" id="'+canvasId+'"></canvas>'
		
		var target = this.wcontainer.update(textHTML+canvasHTML)

		var opts = {
			lines: this.lines,
			angle: this.angle/100,
			lineWidth: this.lineWidth,
			pointer: this.pointer,
			colorStart: this.colorStart,
			colorStop: this.colorStop,
			strokeColor: this.strokeColor,
			generateGradient: this.generateGradient
		}

		this.gauge = new Gauge(document.getElementById(canvasId));
		this.gauge.setOptions(opts)
		this.gauge.setTextField(document.getElementById(textId));

		this.gauge.maxValue = this.maxValue;
		this.gauge.animationSpeed = this.animationSpeed
		
	},

	onResize: function() {
		log.debug('onRezize', this.logAuthor);
		this.wcontainer.removeAll()
		this.afterContainerRender()
	},

	getNodeInfo: function(from,to) {
		this.processNodes()
		if (this.nodeId) {
			Ext.Ajax.request({
				url: this.makeUrl(from,to),
				scope: this,
				params: this.post_params,
				method: 'POST',
				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);
					if (this.nodeId.length > 1)
						data = data.data;
					else
						data = data.data[0];
					this._onRefresh(data);
				},
				failure: function(result, request) {
					log.error('Impossible to get Node informations, Ajax request failed ... (' + request.url + ')', this.logAuthor);
				}
			});
		}

	},

	makeUrl: function(from, to) {
		var url = '/perfstore/values';
		if (! to)
			url += '/' + from;
		if (from && to) 
			url += '/' + from + '/' + to;
		return url;
	},

	processNodes: function() {
		var post_params = [];
		for (var i in this.nodes) 
			post_params.push({
				id: this.nodes[i].id,
				metrics: this.nodes[i].metrics
			});
		
		this.post_params = {
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : this.aggregate_method,
			'aggregate_interval': this.aggregate_interval,
			'aggregate_max_points': this.aggregate_max_points
			};
	},


	onRefresh: function(data) {
		log.debug('onRefresh', this.logAuthor);
		console.log(data)
		if(data.max)
			this.gauge.maxValue = data.max
		if(data.values)
			this.gauge.set(data.values[data.values.length - 1][1])
	},

});
