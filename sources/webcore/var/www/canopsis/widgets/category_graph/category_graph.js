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

Ext.define('widgets.category_graph.category_graph', {
	extend: 'canopsis.lib.view.cwidgetGraph',
	alias: 'widget.category_graph',

	logAuthor: '[widgets][category_graph]',

	timeNav: false,
	timeNav_window: global.commonTs.week,

	diagram_type: 'pie', // 'or column'
	aggregate_max_points: 1,
	aggregate_method: 'LAST',
	aggregate_round_time: true,

	//Default Options
	max: 0,

	// Layout options
	other_label: 'Free',

	// pie specific options
	pie_size: 0.8,
	startAngle: 1.5,
	innerRadius: 0,

	// Bar specific options
	stacked_graph: false,
	verticalDisplay: true,

	backgroundColor: '#FFFFFF',
	borderColor: '#FFFFFF',
	borderWidth: 0,
	title_fontSize: 15,

	// legend
	legend: true,
	legend_verticalAlign: 'bottom', // TODO: this property is not managed by flotchart
	legend_align: 'center', // TODO: this property is not managed by flotchart
	legend_layout: 'horizontal', // TODO: this property is not managed by flotchart
	legend_backgroundColor: null, // TODO: this property is not managed by flotchart
	legend_borderColor: '#909090', // TODO: this property is not managed by flotchart
	legend_borderWidth: 1, // TODO: this property is not managed by flotchart
	legend_fontSize: 12, // TODO: this property is not managed by flotchart
	legend_fontColor: '#3E576F', // TODO: this property is not managed by flotchart

	// label
	labels: true,
	nameInLabelFormatter: false,
	pctInLabel: true,
	labels_size: "x-small",

	gradientColor: false, // TODO: this property is not managed by flotchart

	tooltip: true,

	setChartOptions: function() {
		this.callParent(arguments);

		var me = this;

		$.extend(this.options,
			{
				series: {
					lines: {
						show: false,
					},
					pie: {
						show: (this.diagram_type === 'pie'),
						innerRadius: this.innerRadius,
						label: {
							show: this.labels,
							size: this.labels_size
						},
						tilt: this.tilt,
						stroke:{
							color: "fff",
							width: this.stroke_width
						},
						startAngle: this.startAngle,
						radius: this.pie_size / 100
					},
					bars: {
						show: (this.diagram_type === 'column'),
						horizontal: !this.verticalDisplay,
						barWidth: 1,
						zero: true,
						dataLabels: this.labels,
						showNumbers: (this.labels && this.diagram_type === 'column')
					},
					stack: this.stacked_graph
				},
				grid: {
					hoverable: true,
					clickable: true,
				},
				legend: {
					hideable: true,
					show: this.legend,
					labelFormatter: function(label, series) {
						result = me.nameInLabelFormatter ? ("<b>" + label + "</b><br/>") : "";
						result += me.pctInLabel ? (series.data[0] + "%") : yval; // calculate percent
						return result;
					}
				},
				xaxis: {
					show: (this.diagram_type === 'column' && !this.verticalDisplay)
				},
				yaxis: {
					show: (this.diagram_type === 'column' && this.verticalDisplay)
				},
				tooltip: this.tooltip,
				tooltipOpts: {
					content: function(label, xval, yval, flotItem) {
						if(me.humanReadable) {
							for(var serieId in me.series) {
								var serie = me.series[serieId];

								if(serie.label === label) {
									yval = rdr_humanreadable_value(yval, serie.node.bunit);
									break;
								}
							}
						}

						return "<b>" + label + ":</b>" + yval;
					}
				}
			}
		);
	},

	prepareData: function() {
		var label_axis_group = {};
		var label_axis_group_count = 0;

		for(var node_id in this.series) {
			var serie = this.series[node_id];
		}
	},

	addPoint: function(serieId, value, serieIndex) {
		var serie = this.series[serieId];
		var x = serie.label_axis_group;

		log.debug("addPoint");
		log.dump(serie);
		if(serie.node.xpos !== undefined && serie.node.xpos != "")
		{
			log.debug("x pos : " + serie.node.xpos + "for serie :");
			log.dump(serie);

			x = serie.node.xpos;
		}
		else if(this.stacked_graph !== undefined && this.stacked_graph === true) {
			x = 0;
		}
		else {
			x = serieIndex;
		}

		var point = [x, value[1]];

		this.series[serieId].data = [point];
	},

	createChart: function() {
		/* Computes the difference between the series sum and max user input value and then add it to series if max value > series sum */
		//clean series first
		delete this.series['max_diff'];

		if(this.max > 0) {
			var total = 0;
			var serie_count = 0;
			var node_id = undefined;
			var serie = undefined;

			//gets series sum
			for(node_id in this.series) {
				serie_count++;

				if((this.series[node_id].show === undefined || this.series[node_id].show) && this.series[node_id].data.length === 1) {
					total += this.series[node_id].data[0][1];
				}
			}

			if(total < this.max) {
				//data is a list of point with a single computed point (barchart)
				var label = this.other_label;
				var stacked = this.stacked_graph;
				var max = this.max;

				this.series.max_diff = {
					label: label,
					data: [[stacked ? 0: serie_count, max - total]],
					node: {},
				};
			}
		}

		this.callParent(arguments);
	}
});
