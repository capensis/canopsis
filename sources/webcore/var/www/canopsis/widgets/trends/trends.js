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
Ext.define('widgets.trends.trends' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.trends',

	logAuthor: '[trends]',

	wcontainer_layout: 'anchor',

	interval: global.commonTs.hours,
	aggregate_method: undefined,
	aggregate_interval: 0,
	aggregate_max_points: 1,

	item_height: 30,

	wcontainer_autoScroll: true,

	colorLow: "#1BE01B",
	colorMid: "#E0E0E0",
	colorHight: "#E0251B",
	display_pct: true,

	initComponent: function() {
		log.debug('nodesByID:', this.logAuthor);
		log.dump(this.nodesByID);

		// Color Scaling
		var colors = [this.colorLow, this.colorLow, this.colorMid, this.colorHight, this.colorHight];
		this.colorScale = chroma.scale(colors);

		this.callParent(arguments);
	},

	getPostParams: function() {
		var post_params = [];
		Ext.Object.each(this.nodesByID, function(id, node, obj) {
			post_params.push({
				id: id,
				metrics: node.metrics
			});
		},this)

		var post_params = {
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : this.aggregate_method,
			'aggregate_interval': this.aggregate_interval,
			'aggregate_max_points': this.aggregate_max_points,
			//'consolidation_method': this.consolidation_method
		};

		//if (this.aggregate_interval)
		//	this.post_params['aggregate_interval'] = this.aggregate_interval;

		return post_params;
	},

	doRefresh: function(from, to) {
		if ( this.from )  
			from = this.from;

		if ( this.to ) 
			to  = this.to;
		log.debug('Get values from ' + new Date(from) + ' to ' + new Date(to), this.logAuthor);

		if (this.nodesByID && Ext.Object.getSize(this.nodesByID) != 0 ) {
			var url = '/perfstore/values/' + parseInt(from / 1000) + '/' + parseInt(to / 1000);

			Ext.Ajax.request({
				url: url,
				scope: this,
				params: this.getPostParams(),
				method: 'POST',
				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);
					data = data.data;
					this.onRefresh(data);
				},
				failure: function(result, request) {
					log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
				}
			});
		}
	},

	onRefresh: function(data) {
		if (! data.length)
			return;

		this.wcontainer.removeAll();

		for (var i=0; i < data.length; i++){
			var _id = data[i].node;
			var values = data[i].values;
			var bunit = data[i].bunit;
			var node = this.nodesByID[_id];

			var max = data[i].max;
			if (node.max)
				max = node.max;

			log.debug("Node: " + _id, this.logAuthor);
			log.debug(" + Max: " + max, this.logAuthor);

			var x = [];
			var y = [];

			for (var j = 0; j < values.length; j++) {
			  if (values[j] != null && values[j][0] != null && values[j][1] != null) {
			    x.push(values[j][0]);
			    y.push(values[j][1]);
			  }
			}

			var ret = linearRegression(x, y);

			var delta = undefined;
			var delta_pct = undefined;
			var hdelta = 'NaN';
			
			if (values.length >= 2){
				var v1 = values[0][1]
				var t1 = values[0][0]
				var v2 = values[values.length-1][1]
				var t2 = values[values.length-1][0]

				v1 = ret[0]*t1 + ret[1]
				v2 = ret[0]*t2 + ret[1]

				//console.log("v1", v1, "v2", v2)

				delta = roundSignifiantDigit(v2 - v1, 2);
				hdelta = rdr_humanreadable_value(delta, bunit)

				if (delta > 0)
					hdelta = "+" + hdelta;

				log.debug(" + Delta: " + delta, this.logAuthor);

				if (max)
					if (delta > 0)
						delta_pct = Math.round((delta * 100) / max);
					else
						delta_pct = -1 * Math.round((-delta * 100) / max);

				log.debug(" + Delta Pct: " + delta_pct, this.logAuthor);
			}


			//delta_pct = 105

			var fill = this.colorScale(0.5).hex();
			var degrees = 0;

			if (delta_pct != undefined){
				degrees = Math.round((-delta_pct * 90) / 100);

				if (degrees > 90)
					degrees = 90;

				if (degrees < -90)
					degrees = -90;

				fill = this.colorScale(0.5+(delta_pct*0.5)/100).hex();
			}

			log.debug(" + Degrees: " + degrees, this.logAuthor);
			log.debug(" + Fill color: " + fill, this.logAuthor);

			var row =  Ext.create('Ext.draw.Component', {
				viewBox: false,
				autoSize: true,
				//height: 40,
				items: [{
					type: "path",
					path: "M 100,50 L 40,0 L 40,30 L 0,30 L 0,70 L 40,70 L 40,100",
					//stroke: "none",
					fill: fill,
					rotate: {
						degrees: degrees
					},
					scale: {
						x: 0.3,
						y: 0.3
					}
				}]
			});

			var text = hdelta;
			
			// Display as pct
			if (this.display_pct && delta_pct != undefined){
				if (delta_pct > 0)
					delta_pct = "+" + delta_pct;

				text = delta_pct + '%'
			}

			log.debug(" + Text: " + text, this.logAuthor);
			if ( this.caller != undefined && this.caller == 'mini_chart') {
				this.wcontainer.add({
					layout: {
						type: 'hbox',
						//align: 'stretch'
					},
					border: 0,
					margin: 1,
					anchor: "100% 100%",
					items: [
						{ border: 0, height: this.item_height, html: String(text),  bodyStyle: { "line-height": this.item_height+ "px", "text-align": "right", "padding-right": "3px"}},
						row
					]
				});

			} else {
				this.wcontainer.add({
					layout: {
						type: 'hbox',
						//align: 'stretch'
					},
					border: 0,
					margin: 1,
					items: [
						{ border: 0, height: this.item_height, html: node.label, flex:1, bodyStyle: { "line-height": this.item_height + "px" } },
						{ border: 0, height: this.item_height, html: String(text), bodyStyle: { "line-height": this.item_height+ "px", "text-align": "right", "padding-right": "3px"}},
						row
					]
				});
			}
		}
	},

});
