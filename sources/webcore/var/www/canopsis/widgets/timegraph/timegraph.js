//need:app/lib/view/cwidgetGraph.js
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

		this.options.xaxes.push({
			position: 'bottom',
			mode: 'time',
			timeformat: '%d %b - %H:%M:%S'
		});

		if( !this.displayVerticalLines)
		{
			this.options.xaxis.tickLength = 0;
		}

		if( !this.displayHorizontalLines)
		{
			this.options.yaxis.tickLength = 0;
		}

		this.options.series = {
			shadowSize :0,
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
		};

		if(this.tooltip) {
			this.options.tooltip = this.tooltip;

			this.options.tooltipOpts = {};
			this.options.tooltipOpts.content = "<b>%x<br/>%s :</b> %y";
		}

		this.options.legend = {
			hideable: true
		};

		if(this.timeNav) {
			/* copy options, but override some */
			this.options_overview = {
				crosshair: {
					mode: 'x'
				},

				selection: {
					mode: 'x'
				},

				grid: {
					borderWidth: {
						top: 0,
						bottom: 0,
						right: 0,
						left: 0
					},
					hoverable: true,
					clickable: true
				},

				xaxis: {
					min: now - this.timeNav_window * 1000,
					max: now,
					show: false,
				},
				yaxis: {
					show: false,
				},

				legend: {
					hideable: true,
					show: false
				}
			};
		}
	},

	createChart: function() {
		var me = this;

		/* initialize time navigation parameters if needed */
		if(this.timeNav) {
			var overview_h = this.getHeight() / 5;

			// NB: this.plotcontainer doesn't exist yet.
			var plotcontainer = $('#' + this.wcontainerId);
			plotcontainer.nextAll().remove();

			this.plotoverview = $('<div/>');
			this.plotoverview.width(this.wcontainer.getWidth());
			this.plotoverview.height(overview_h);

			this.plotoverview.attr('class', plotcontainer.attr('class'));
			plotcontainer.height(this.getHeight() - overview_h);

			plotcontainer.after(this.plotoverview);
		}

		/* create chart with modified plotcontainer */
		this.callParent(arguments);

		/* create the overview chart */
		if(this.timeNav) {
			this.chart_overview = $.plot(this.plotoverview, this.getSeriesConf(), this.options_overview);

			/* connect the two charts */
			this.plotcontainer.bind('plotselected', function(event, ranges) {
				void(event);

				console.log("Selected Range: " + ranges.xaxis.from + ' -> ' + ranges.xaxis.to);

				me.chart.getOptions().xaxes[0].min = ranges.xaxis.from;
				me.chart.getOptions().xaxes[0].max = ranges.xaxis.to;
				me.chart.clearSelection(true);

				me.chart.setupGrid();
				me.chart.draw();
				me.chart.autoScale();

				me.chart_overview.setSelection(ranges, true);

			});

			this.plotoverview.bind('plotselected', function(event, ranges) {
				void(event);

				me.chart.setSelection(ranges);
			});
		}
	},

	renderChart: function() {
		/* update container size before rendering */
		if(this.timeNav) {
			var overview_h = this.getHeight() / 5;

			this.plotcontainer.height(this.getHeight() - overview_h);
		}

		this.callParent(arguments);

		/* now render overview chart */
		if(this.timeNav) {
			this.chart_overview.setData(Ext.clone(this.getSeriesConf()));
			this.chart_overview.setupGrid();
			this.chart_overview.draw();
		}
	},

	destroyChart: function() {
		this.callParent(arguments);

		if(this.timeNav) {
			this.chart_overview.destroy();
		}
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
			fill: (curve_type === 'area')
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

	shiftSerie: function(serieId) {
		var now = Ext.Date.now();
		var timestamp = now - this.timeNav_window * 1000;

		while(this.series[serieId].data[0][0] < timestamp) {
			this.series[serieId].data.shift();
		}
	},

	doRefresh: function(from, to) {
		if(this.timeNav) {
			var now = Ext.Date.now();

			to = now;
			from = now - this.timeNav_window * 1000;
		}

		this.refreshNodes(from, to);
	},

	dblclick: function() {
		log.debug('Zoom Out', this.logAuthor);

		this.chart.getOptions().xaxes[0].min = this.chart.getOptions().xaxis.min;
		this.chart.getOptions().xaxes[0].max = this.chart.getOptions().xaxis.max;

		this.chart.setupGrid();
		this.chart.draw();
		this.chart.autoScale();
	}
});