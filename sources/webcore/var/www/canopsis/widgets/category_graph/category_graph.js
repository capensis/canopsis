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

Ext.define('widgets.timegraph.category_graph', {
	extend: 'canopsis.lib.view.cwidgetGraph',
	alias: 'widget.category_graph',

	logAuthor: '[widgets][category_graph]',

	timeNav: false,
	timeNav_window: global.commonTs.week,

	graph_type: 'pie',
	aggregate_max_points: 1,
	aggregate_method: 'LAST',
	aggregate_round_time: true,

	//Default Options
	max: undefined,

	// Layout options
	other_label: 'Free',

	// pie specific options
	pie_size: 60,
	innerRadius: 0,

	// Bar specific options
	stacked_graph: false,
	vertical_display: false,

	backgroundColor: '#FFFFFF',
	borderColor: '#FFFFFF',
	borderWidth: 0,
	title_fontSize: 15,

	legend_verticalAlign: 'bottom',
	legend_align: 'center',
	legend_layout: 'horizontal',
	legend_backgroundColor: null,
	legend_borderColor: '#909090',
	legend_borderWidth: 1,
	legend_fontSize: 12,
	legend_fontColor: '#3E576F',

	labels: false,
	gradientColor: false,
	pctInLabel: false,

	tooltip: true,

	initComponent: function() {
		this.callParent(arguments);
	},

	setChartOptions: function() {
		this.callParent(arguments);

		$.extend(this.options,
			{
				series: {
				//stack: this.stacked_graph,
					pie: {
						show: (this.graph_type === 'pie'),
						innerRadius: this.innerRadius,
						label: {
							show: this.labels,
						},
						tilt: this.tilt,
					},
					bars: {
						show: (this.graph_type === 'column')
					},
					stack: this.stacked_graph
				},
				grid: {
					hoverable: true,
					clickable: true,
				},
				legend: {
					hideable: true
				},
				xaxis: {
					show: false
				},
				yaxis: {
					min: 0,
					show: (this.graph_type === 'column')
				}
			}
		);

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
					show: false
				},

				yaxis: {
					show: false
				},

				legend: {
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

		return serie;
	},

	addPoint: function(serieId, value) {
		var serie = this.series[serieId];
		serie.data = this.graph_type=='pie'? value[1] : [[0, value[1]]];
		serie.last_timestamp = value[0] * 1000;
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
	}
});