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

// initComponent -> afterContainerRender 	-> setchartTitle -> ready -> doRefresh -> onRefresh -> addDataOnChart
//                                			-> setOptions                             			-> getSerie
//											-> createChart


Ext.define('widgets.categorized_graph.categorized_graph' , {
	extend: 'widgets.pie.pie',

	alias: 'widget.categorized_graph',

	logAuthor: '[categorized_graph]',

	legend: false,
	labels:true,
	graph_type : 'column',

	hide_other_column: true,

	interval: global.commonTs.hours,
	aggregate_method: 'LAST',
	aggregate_interval: 0,
	aggregate_max_points: 1,

	getSerie: function(){
		return  {
					id: 'serie',
					type: this.graph_type,
					data: []
				};
	},

	setAxis: function(data){
		var metrics = []
		for(var i in data)
			if(data[i].metric)
				metrics.push(data[i].metric)

		this.chart.xAxis[0].setCategories(metrics, false)

	},

	setOptions: function() {
		this.callParent(arguments);

		if(this.graph_type == 'column'){
			this.options.legend = {enabled:false}
			this.options.yAxis = [{title: { text: null }}]

			if(this.labels){
				this.options.plotOptions.column.dataLabels = {
					enabled: true,
					formatter: function() {
						if(this.y)
							return '<b>'+this.y +'</b>';
						else
							return ''
					}
				}
			}
		}

	},

	processNodes: function() {
		var post_params = [];
		for (var i in this.nodes) {
			post_params.push({
				id: this.nodes[i].id,
				metrics: this.nodes[i].metrics
			});
		}
		this.post_params = {
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : this.aggregate_method,
			'aggregate_max_points': this.aggregate_max_points
		};

		if(this.aggregate_interval)
			this.post_params['aggregate_interval'] = this.aggregate_interval
	},

	tooltipFunction: function() {
		return this.key+': ' + this.y
	}

});