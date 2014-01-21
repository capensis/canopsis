//need:app/lib/view/cwidgetGraph.js
/*
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
*/

Ext.define('widgets.timegraph.timegraph', {
	extend: 'canopsis.lib.view.cwidgetGraph',
	alias: 'widget.timegraph',

	logAuthor: '[widgets][timegraph]',

	timeNav: false,
	timeNav_window: global.commonTs.week,

	initComponent: function() {
		this.callParent(arguments);
	},

	setChartOptions: function() {
		this.callParent(arguments);

		var now = Ext.Date.now();

		$.extend(this.options,
			{
				zoom: {
					interactive: false
				},

				selection: {
					mode: 'x'
				},

				crosshair: {
					mode: 'x'
				},

				grid: {
					hoverable: true,
					clickable: true
				},

				xaxis: {
					min: now - this.time_window * 1000,
					max: now
				},

				xaxes: [
					{
						position: 'bottom',
						mode: 'time'
					}
				],

				yaxes: [],

				legend: {
					hideable: true
				},

				series: {
					shadowSize: 0,
					stack: this.stacked_graph,
					lines: {
						show: (this.SeriesType === 'area' || this.SeriesType === 'line'),
						fill: (this.SeriesType === 'area'),
						lineWidth: this.lineWidth
					},
					points: {
						show: false
					},
					bars: {
						show: (this.SeriesType === 'bars')
					}
				}
			}
		);

		if( !this.displayVerticalLines)
		{
			this.options.xaxis.tickLength = 0;
		}

		if( !this.displayHorizontalLines)
		{
			if(this.options.yaxis === undefined)
				this.options.yaxis = {};

			this.options.yaxis.tickLength = 0;
		}
	},

	insertGraphExtraComponents: function(){
		this.callParent(arguments);
	},

	createChart: function() {
		var me = this;

		// NB: this.plotcontainer doesn't exist yet.
		if(!!this.plotcontainer)
			this.plotcontainer.nextAll().remove();

		/* create chart with modified plotcontainer */
		this.callParent(arguments);
	},

	renderChart: function() {
		this.callParent(arguments);

		this.chart.setupGrid();
		this.chart.draw();
	},

	destroyChart: function() {
		this.callParent(arguments);
	},

	getSerieForNode: function(nodeid) {
		var serie = this.callParent(arguments);
		var node = serie.node;

		var curve_type = this.SeriesType;

		if(node.curve_type && node.curve_type !== 'default') {
			curve_type = node.curve_type;
		}

		serie.lines = {
			show: (curve_type === 'area' || curve_type === 'line'),
			fill: (curve_type === 'area'),
			points: {
				symbol: node.point_shape? node.point_shape : undefined,
			}
		};

		serie.bars = {
			show: (curve_type === 'bars')
		};

		return serie;
	},

	addPoint: function(serieId, value) {
		this.series[serieId].data.push([value[0] * 1000, value[1]]);
		this.series[serieId].last_timestamp = value[0] * 1000;
	},

	// shiftSerie: function(serieId) {
	// 	var now = Ext.Date.now();
	// 	var timestamp = now - this.timeNav_window * 1000;

	// 	while(this.series[serieId].data[0][0] < timestamp) {
	// 		this.series[serieId].data.shift();
	// 	}
	// },

	dblclick: function() {
	},

	updateAxis: function(from, to) {
		if(this.reportMode || this.exportMode) {
			this.updateAxis();
			this.options.xaxis.min = from;
			this.options.xaxis.max = to;
		}
		else {
			this.options.xaxis.min = to - (this.time_window + this.time_window_offset) * 1000;
			this.options.xaxis.max = to - this.time_window_offset * 1000;
		}
	}

});
