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

// initComponent -> doRefresh -> get_config -> createHighchartConfig -> doRefresh -> addDataOnChart

Ext.define('widgets.stock_graph.stock_graph' , {
	extend: 'widgets.line_graph.line_graph',

	alias: 'widget.stock_graph',

	logAuthor: '[stock_graph]',

	time_window: global.commonTs.year,

	setOptions: function() {
		this.callParent();
		this.options.tooltip.formatter = undefined;

		//this.options.xAxis.maxZoom = 60 * 5 * 1000 //5 minutes

		this.options.rangeSelector = {
			selected: 0,
			inputEnabled: false
		};
	},

	createChart: function() {
		this.chart = new Highcharts.StockChart(this.options);
	}
});
