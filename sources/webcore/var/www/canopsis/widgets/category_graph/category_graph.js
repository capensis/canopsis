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

	diagram_type: 'pie', // 'column'
	aggregate_max_points: 1,
	aggregate_method: 'LAST',

	//Default Options
	max: 0,

	// Layout options
	other_label: 'Free',

	// pie specific options
	pie_size: 60,
	startAngle: 1.5,
	radius: 0.8,
	innerRadius: 0,

	// Bar specific options
	stacked_graph: false,
	verticalDisplay: true,

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

	labels: true,
	labels_size: "x-small",
	legend: true,
	gradientColor: false,
	pctInLabel: false,
	tooltip: true,

	setChartOptions: function() {
		this.callParent(arguments);

		$.extend(this.options,
			{
				series: {
				//stack: this.stacked_graph,
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
						startAngle: this.startAngle,
						radius: this.radius
					},
					bars: {
						show: (this.diagram_type === 'column'),
						horizontal: !this.verticalDisplay,
						barWidth: 1,
						zero: true,
						dataLabels: this.labels
					},
					stack: this.stacked_graph
				},
				grid: {
					hoverable: true,
					clickable: true,
				},
				legend: {
					hideable: true,
					show: this.legend
				},
				xaxis: {
					show: (this.diagram_type === 'column' && !this.verticalDisplay)
				},
				yaxis: {
					show: (this.diagram_type === 'column' && this.verticalDisplay)
				},
				tooltip: this.tooltip,
				tooltipOpts : {
					content: function(label, xval, yval, flotItem) {
						void(xval);
						void(flotItem);
                        return "<b>" + label + "<br/></b>" + yval;
                    }
				}
			}
		);
	},

	addPoint: function(serieId, value, serieIndex) {
		var point = [this.diagram_type === 'column' && !this.stacked_graph? serieIndex : 0, value[1]];
		this.series[serieId].data = [point];
	},

	getSeriesConf: function() {
		var series = this.callParent(arguments);

		if (this.max > 0) {
			var total = 0;
			for (var index=0; index < series.length; index++) {
				var serie = series[index];
				if (serie.show && serie.data.length === 1) {
					total += series[index].data[0][1];
				}
			}
			// add other serie if total is less than this.max
			if (this.max > total) {
				// set x=0 only if graph is in stacking mode and is column
				var point = [this.diagram_type === 'column' && !this.stacked_graph? series.length: 0, this.max - total];
				var other_serie = {
					label: this.other_label,
					data: [point]
				};
				series.push(other_serie);
			}
		}

		return series;
	},

});
