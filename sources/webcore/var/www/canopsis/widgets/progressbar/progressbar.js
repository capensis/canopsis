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

Ext.define('widgets.progressbar.progressbar' , 
{
	extend: 'canopsis.lib.view.cwidget',
	alias: 'widget.progressbar',
	logAuthor: '[progressBarWidget]',
	innerText: undefined,
	refresh_number: 0,
	wcontainer_autoScroll: true,
	wcontainer_layout: {type:'anchor'},
	bodyPadding: '5 5 5 5',

	colorBg: '#EEEEEE',
	colorStart: '#1BE01B',
	colorMid: '#FFCD43',
	colorEnd: '#E0251B',

	dispGrad: true,
	
	//boldText: true,
	//fontSize: 100,
	
	gHeight: 20,

	initComponent: function()
	{
		log.debug("initComponent", this.logAuthor)

		this.nodesByID = parseNodes(this.nodes);
		log.debug('nodesByID:', this.logAuthor);
		log.dump(this.nodesByID);

		// Color Scaling
		var colors = [this.colorStart, this.colorMid, this.colorEnd];
		this.colorScale = chroma.scale(colors);

		this.progressBars = {};

		this.callParent(arguments);
	},

	getNodeInfo: function(from,to) 
	{
		this.processNodes();
		if (this.nodesByID) 
		{
			Ext.Ajax.request(
			{
				url: '/perfstore/values' 
					+ '/' + parseInt(to/1000)+ '/' + parseInt(to/1000),
				scope: this,
				params: this.post_params,
				method: 'POST',
				success: function(response) 
				{
					var data = Ext.JSON.decode(response.responseText);
					data = data.data;
					this.displayBars(data);
				}, 
				failure: function(result, request) 
				{
					log.error
					(
						'get Node info, Ajax req failed ... (' 
						+ request.url + ')', this.logAuthor
					);
				}
			});
		}
	},

	setGradient: function(_id, value){
		log.debug('setGradient: '+_id+", value: "+value, this.logAuthor);

		var lowColor = this.colorScale(value).hex()
		var hightColor = chroma.hex(lowColor).brighten(20).hex()

		pbEl = $('#'+_id + " .x-progress-bar");
		pbEl.css('background-image', "none");

		if (this.dispGrad){
			pbEl.css('background', '-webkit-gradient(linear, left top, left bottom, from('	+ hightColor +'), to('+lowColor+'))');
			pbEl.css('background', '-webkit-linear-gradient(' + hightColor +', '+lowColor+')');
			pbEl.css('background', '-moz-linear-gradient('	+ hightColor +', '+lowColor+')');
			pbEl.css('background', '-ms-linear-gradient('	+ hightColor +', '+lowColor+')');
			pbEl.css('background', '-o-linear-gradient('	+ hightColor +', '+lowColor+')');
			pbEl.css('background', 'linear-gradient(to bottom, '	+ hightColor +', '+lowColor+')');
		}else{
			pbEl.css('background-color', lowColor);
		}

		//pbEl.height(this.gHeight);

	},

	displayBars: function(data)
	{
		for (var i = 0; i < data.length; i++){

			var item = data[i];
			var _id = item.node;


			// Create it
			log.debug('Item: ' + _id, this.logAuthor);
			log.dump(item);

			var node = this.nodesByID[_id];
			log.debug('Node:', this.logAuthor);
			log.dump(node);
			
			var label = node.label;

			log.debug(' + Label: '+label, this.logAuthor);

			var value = item.values[0][1];
			var max = item.max;

			//Extra field
			if (node.max)
				max = node.max

			log.debug(' + Value: '+value, this.logAuthor);
			log.debug(' + Max:   '+max, this.logAuthor);

			var pct = 0;
			if (value && max){
				pct = (value * 100)/max;
				pct = roundSignifiantDigit(pct, 2);
			}

			log.debug(' + Pct:   '+pct, this.logAuthor);

			var text = pct + "%";

			// Check if pb already exist
			if (this.progressBars[_id]){
				var pb = this.progressBars[_id];
				log.debug('Update: ' + _id, this.logAuthor)
				pb.updateText(text);
				pb.updateProgress(pct/100);
				return
			}

			var pb = Ext.create('Ext.ProgressBar', {
				text: text,
				value: pct/100,
				flex:1,
				height: this.gHeight,
				cls:'widgets-progressbar',
				border: 1,
			});

			pb.on("update", function(pb, value){ this.setGradient(pb.id, value); }, this);
			pb.on("afterrender", function(pb){ this.setGradient(pb.id, pb.value); }, this, {single: 1});

			this.wcontainer.add({
				layout: {
					type: 'hbox'
				},
				border: 0,
				margin: 4,
				height: this.gHeight,
				items: [
					{ border: 0, html: String(label), flex:1 },
					pb
				]
			});

			this.progressBars[_id] = pb;
		}
	},

	processNodes: function() 
	{
		var post_params = [];
		Ext.Object.each(this.nodesByID, function(id, node, obj) {
			post_params.push({
				id: id,
				metrics: node.metrics
			});
		},this)

		this.post_params = 
		{
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : this.aggregate_method,
			'aggregate_interval': this.aggregate_interval,
			'aggregate_max_points': this.aggregate_max_points
		};
	}
});
