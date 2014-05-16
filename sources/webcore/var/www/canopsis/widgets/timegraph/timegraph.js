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

	flagRadius: 10,
	flagLineWidth: 1,
	flagLineColor: '#003300',
	flagHeight: 35,
	flagTooltip: undefined,

	initComponent: function() {
		this.callParent(arguments);

		this.logevents = [];
	},

	doRefresh: function(from, to) {
		if(this.aggregate_method) {
			from = to - (this.time_window * 1000);
		}

		this.refreshNodes(from, to);

		if(this.flagFilter) {
			var filter = {'$and': [
				{
					'timestamp': {
						'$gte': parseInt(from / 1000),
						'$lte': parseInt(to / 1000)
					}
				},{
					'state_type': 1
				},
				Ext.JSON.decode(this.flagFilter)
			]};

			Ext.Ajax.request({
				url: '/rest/events_log',
				method: 'GET',
				params: {
					filter: Ext.JSON.encode(filter),
					limit: this.nbMaxEventsDisplayed
				},
				scope: this,

				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);
					data = data.data;

					// add events to the list
					for(var i = 0; i < data.length; i++) {
						this.logevents.push({
							x: data[i].timestamp * 1000,
							y: 0,
							event: data[i]
						});
					}

					// then shift it
					if(this.logevents.length > this.nbMaxEventsDisplayed) {
						this.logevents = this.logevents.slice(this.logevents.length - this.nbMaxEventsDisplayed);
					}

					this.chart.triggerRedrawOverlay();
				}
			});
		}
	},

	createChart: function() {
		this.callParent(arguments);
		var me = this;

		this.plotcontainer.click(function(e) {
			if(me.flagTooltip !== undefined) {
				me.flagTooltip.remove();
				me.flagTooltip = undefined;
			}

			for(var i = 0; i < me.logevents.length; i++) {
				var evt = me.logevents[i];
				var coord = me.chart.p2c(evt);

				var d2 = Math.pow(coord.left - e.offsetX, 2) + Math.pow(coord.top - me.flagHeight - e.offsetY, 2);

				if(d2 <= Math.pow(me.flagRadius, 2)) {
					me.flagTooltip = $('<div/>', {
						id: me.wcontainerId + '-flag-tooltip'
					});

					me.flagTooltip.css({
						position: 'absolute',
						padding: '5px',
						'border-radius': '5px',
						left: coord.left,
						top: coord.top - me.flagHeight - me.flagRadius,
						border: '1px solid black',
						background: '#FFFFFF'
					});

					me.flagTooltip.append('<p><b>' + evt.event.display_name + '</b></p><ul>');

					if(evt.event.source_type === 'component') {
						me.flagTooltip.append('<li><em>Source:</em> ' + evt.event.component + '</li>');
					}
					else {
						me.flagTooltip.append('<li><em>Source:</em> ' + evt.event.component + '/' + evt.event.resource + '</li>');						
					}

					me.flagTooltip.append('<li><em>Message:</em> ' + (evt.event.output ? evt.event.output : '') + '</li></ul>');

					if(evt.event.long_output) {
						me.flagTooltip.append('<p>' + evt.event.long_output + '</p>');
					}

					$('body').append(me.flagTooltip);
					return;
				}
			}
		});

		this.chart.hooks.drawOverlay.push(function(plot, ctx) {
			me.addLogEventsToChart(ctx);
		});
	},

	addLogEventsToChart: function(ctx) {
		var state_colors = [
			'green',  // 0 : OK
			'yellow', // 1 : WARNING
			'red',    // 2 : CRITICAL
			'orange'  // 3 : UNKNOWN
		];

		for(var i = 0; i < this.logevents.length; i++) {
			var evt = this.logevents[i];
			var coord = this.chart.p2c(evt);

			ctx.lineWidth = this.flagLineWidth;
			ctx.strokeStyle = this.flagLineColor;

			ctx.beginPath();
			ctx.moveTo(coord.left, coord.top - this.flagHeight);
			ctx.lineTo(coord.left, coord.top + 7);
			ctx.stroke();

			ctx.beginPath();
			ctx.arc(coord.left, coord.top - this.flagHeight, this.flagRadius, 0, 2 * Math.PI, false);
			ctx.fillStyle = state_colors[evt.event.state];
			ctx.fill();
			ctx.stroke();
		}
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
						mode: 'time',
						timezone: 'browser'
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

		for(var serieId in this.series) {
			var serie = this.series[serieId];
			var idx = serie.yaxis - 1;

			if(!this.options.yaxes[idx]) {
				this.options.yaxes[idx] = {
					position: ((idx % 2) === 0) ? 'left' : 'right'
				};
			}

			this.options.yaxes[idx].font = {
				color: serie.color
			};
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

	prepareData: function(serieId) {
		this.callParent(arguments);

		if(this.aggregate_method) {
			this.series[serieId].data = [];
		}
	},

	addPoint: function(serieId, value) {
		// insert point only if it appends after the last of the serie.
		//gets invert information
		var style = global.curvesCtrl.getRenderInfo(this.series[serieId].node.label)
		var invert = false;
		if (style && style.data.invert) {
			invert = true;
		}

		var points = this.series[serieId].data,
			last_point = points[points.length - 1],
			value_ts = value[0] * 1000;

		if (last_point === undefined || last_point[0] < value_ts) {
			//invert metric depending on it s curve information
			var value = invert ? -value[1] : value[1];
			this.series[serieId].data.push([value_ts, value]);
			this.series[serieId].last_timestamp = (this.aggregate_method ? undefined : value_ts);
		}

	},

	shiftSerie: function(serieId) {
		var now = Ext.Date.now();
		var timestamp;

		if (this.series[serieId].data.length > 0) {

			if(this.timeNav) {
				timestamp = now - this.timeNav_window * 1000;
			}
			else {
				timestamp = now - this.time_window * 1000;
			}

			while(this.series[serieId].data.length > 0 && this.series[serieId].data[0][0] < timestamp) {
				this.series[serieId].data.shift();
			}

		}
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
