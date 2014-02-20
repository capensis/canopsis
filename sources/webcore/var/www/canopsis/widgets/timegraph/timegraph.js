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
	time_window_offset: 0,

	interval: global.commonTs.hours,
	aggregate_method: 'MEAN',
	aggregate_interval: 0,
	aggregate_max_points: 500,
	aggregate_round_time: true,
	consolidation_method: "",

	legend: true,
	legend_fontSize: 12,
	legend_fontColor: "3E576F",
	legend_borderWidth: 1,
	legend_backgroundColor: "FFFFFF",
	legend_borderColor: "909090",
	legend_verticalAlign: "bottom",
	legend_align: "center",
	legend_layout: "horizontal",

	tooltip: true,
	tooltip_crosshairs: true, // TODO: to manage

	SeriesType: 'area',
	lineWidth: 1,
	marker_symbol: null,
	marker_radius: 0,
	stacked_graph: false,

	tooltip_shared: false,
	zoom: true,
	backgroundColor: "FFFFFF",
	borderColor: "FFFFFF",
	borderWidth: 0,

	displayVerticalLines: false,
	displayHorizontalLines: true,

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
					min: now - (this.time_window_offset + this.time_window) * 1000,
					max: now - this.time_window_offset * 1000
				},

				yaxis: {
					tickFormatter: function(val, axis) {
						if(this.humanReadable) {
							return rdr_humanreadable_value(val, axis.options.unit);
						}
						else {
							return val + ' ' + axis.options.unit;
						}
					}.bind(this)
				},

				xaxes: [
					{
						position: 'bottom',
						mode: 'time'
					}
				],

				yaxes: [],

				legend: {
					hideable: true,
					legend: this.legend,
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
				},

				tooltip: this.tooltip
			}
		);

		if(!this.displayVerticalLines) {
			this.options.xaxis.tickLength = 0;
		}

		if(!this.displayHorizontalLines) {
			this.options.yaxis.tickLength = 0;
		}
	},

	insertGraphExtraComponents: function() {
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
				symbol: node.point_shape ? node.point_shape : undefined,
			}
		};

		serie.bars = {
			show: (curve_type === 'bars')
		};

		return serie;
	},

	getSeriesConf: function() {
		var series = this.callParent(arguments);

		function getY(reg, x) {
			return x * reg[0] + reg[1];
		}

		function getLinearRegressionPoint(reg, point) {
			return [point[0], getY(reg, point[0])];
		}

		for(var series_index = (series.length - 1); series_index >= 0; series_index--) {
			var serie = series[series_index];

			if(serie.node.trend_curve && serie.data.length > 1) {
				var x = [], y = [];

				for(var serie_index = 0; serie_index < serie.data.length; serie_index++) {
					var data = serie.data[serie_index];
					x.push(data[0]);
					y.push(data[1]);
				}

				var ret = linearRegression(x, y);
				var trend_serie = this.getSerieForNode(serie.node.id);

				var data = [
					getLinearRegressionPoint(ret, serie.data[0]),
					getLinearRegressionPoint(ret, serie.data[serie.data.length - 1])
				];

				trend_serie.label += '_trend';
				trend_serie.data = data;

				series.splice(series_index + 1, 0, trend_serie);
			}
		}

		return series;
	},

	addPoint: function(serieId, value) {
		this.series[serieId].data.push([value[0] * 1000, value[1]]);
		this.series[serieId].last_timestamp = value[0] * 1000;
	},

	shiftSerie: function(serieId) {
		var now = Ext.Date.now();
		var timestamp = now - this.timeNav_window * 1000;

		while(this.series[serieId].data[0][0] < timestamp) {
			this.series[serieId].data.shift();
		}
	},

	dblclick: function() {
	},

	updateAxis: function(from, to) {
		if(this.reportMode || this.exportMode) {
			this.options.xaxis.min = from;
			this.options.xaxis.max = to;
		}
		else {
			this.options.xaxis.min = to - (this.time_window + this.time_window_offset) * 1000;
			this.options.xaxis.max = to - this.time_window_offset * 1000;
		}
	}
});
