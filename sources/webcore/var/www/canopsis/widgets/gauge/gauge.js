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

	pointerLength: 0.9,
	pointerWidth: 0.035,
	pointerColor: '#000000',

	colorStart: '#6FADCF',   // Colors
	colorStop: '#8FC0DA',    // just experiment with them
	strokeColor: '#EEEEEE',   // to see which ones work best for you
	generateGradient: true,
	
	textSize : 40,
	labelSize: 25,
	maxValue:100,
	animationSpeed:30,

	afterContainerRender: function() {
		this.callParent(arguments);

		var width = this.wcontainer.getWidth()
		var height = this.getHeight() - (this.textSize + this.labelSize + 10)

		var canvasId = this.wcontainerId + '-canvas'
		var textId = this.wcontainerId + '-text'
		var labelId = this.wcontainerId + '-label'

		log.debug('canvasId: ' + canvasId, this.logAuthor);
		log.debug('textId: ' + textId, this.logAuthor);
		log.debug('labelId: ' + labelId, this.logAuthor);

		var textHTML = '<div id="'+textId+'" style="font-size: '+this.textSize+'px;text-align:center;"></div>'
		var canvasHTML = '<canvas width="'+width+'" height="'+height+'" id="'+canvasId+'"></canvas>'
		var labelHTML = '<div id="'+labelId+'" style="font-size: '+this.labelSize+'px;text-align:center;color:#3E576F"></div>'

		if(this.title)
			var target = this.wcontainer.update(canvasHTML+textHTML)
		else
			var target = this.wcontainer.update(labelHTML+canvasHTML+textHTML)

		var opts = {
			lines: this.lines,
			angle: this.angle/100,
			lineWidth: this.lineWidth,
			pointer: {
				length: this.pointerLength, // The radius of the inner circle
				strokeWidth: this.pointerWidth, // The rotation offset
				color: this.pointerColor // Fill color
			},
			colorStart: this.colorStart,
			colorStop: this.colorStop,
			strokeColor: this.strokeColor,
			generateGradient: this.generateGradient
		}

		this.gauge = new Gauge(document.getElementById(canvasId));
		this.gauge.setOptions(opts)
		this.gauge.setTextField(document.getElementById(textId));

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
				url: '/perfstore/values' + '/' + to + '/' + to,
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

		var label = Ext.get(this.wcontainerId + '-label')

		var fields = undefined
		if(this.nodes[0].extra_field)
			fields = this.nodes[0].extra_field

		//update metric name
		if(label){
			if(fields && fields.label)
				label.update(fields.label)
			else
				label.update(data.metric)
		}
			
		//update metric value
		if(fields && fields.ma){
			this.gauge.maxValue = fields.ma
		}else if(data.max){
			this.gauge.maxValue = data.max
		}else{
			this.gauge.maxValue = this.maxValue;
		}

		try{
			if(data.values)
				this.gauge.set(data.values[data.values.length - 1][1])
		}catch(err){
			log.error('Error while set value:' + err, this.logAuthor)
		}
	},

});
