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

/*

+ initComponent
|--- setChartTitle
+ afterContainerRender
|--- setOptions
|--- createChart
\--+ ready
   \--+ doRefresh
      \--+ refreshNodes
         \--+ onRefresh
            \--+ addDataOnChart
               \--- getSerie

*/

Ext.define('widgets.bar_graph.bar_graph', {
	extend: 'widgets.line_graph.line_graph',

	alias: 'widget.bar_graph',

	layout: 'fit',

	first: false,

	last_from: false,
	pushPoints: false,

	logAuthor: '[bar_graph]',

	options: {},
	chart: false,

	params: {},

	metrics: [],

	chartTitle: null,

	//Default Options
	time_window: global.commonTs.day,
	zoom: true,
	legend: true,
	tooltip: true,
	tooltip_crosshairs: true,
	tooltip_shared: false,

	backgroundColor: '#FFFFFF',
	borderColor: '#FFFFFF',
	borderWidth: 0,

	exporting_enabled: false,

	showWarnCritLine: true,

	marker_symbol: null,
	marker_radius: 2,

	title_fontSize: 15,

	chart_type: 'column',

	legend_verticalAlign: 'bottom',
	legend_align: 'center',
	legend_layout: 'horizontal',
	legend_backgroundColor: null,
	legend_borderColor: '#909090',
	legend_borderWidth: 1,
	legend_fontSize: 12,
	legend_fontColor: '#3E576F',
	// 10 minutes
	maxZoom: 60 * 10,

	interval: global.commonTs.hours,
	aggregate_method: 'MAX',
	aggregate_interval: 0,
	aggregate_max_points: 500,

	SeriesType: 'area',
	SeriePercent: false,
	lineWidth: 1,

	//trends
	data_trends: [],
	trend_lines: false,
	trend_lines_type: 'ShortDot',

	nb_node: 0,
	same_node: true,

	//column chart specific
	columnDatalabels: false,
	verticalDisplay: false,


	setOptions: function() {
		this.callParent(arguments);

		if(this.columnDatalabels) {
			this.options.plotOptions.column.dataLabels = {
				enabled: true,
				formatter: function() {
					if(this.y) {
						if(this.humanReadable) {
							return rdr_humanreadable_value(this.y, this.point.bunit);
						}
						else {
							if(this.point.bunit) {
								return this.y + ' ' + this.point.bunit;
							}
							else {
								return this.y;
							}
						}
					}

					return '';
				}
			};
		}

		if(this.verticalDisplay) {
			this.options.chart.inverted = true;
		}
	}
});
