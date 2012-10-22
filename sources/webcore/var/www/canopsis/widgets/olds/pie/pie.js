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
Ext.define('widgets.pie.pie' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.pie',

	logAuthor: '[pie]',

	initComponent: function() {
		this.callParent(arguments);

	},

	onRefresh: function(data) {
		this.uri += '/' + this.nodeId;
		if (this.chart) {
			if (data.perf_data_array) {
				this.parseData(data.perf_data_array);
			}
		}else {
			this.createHighchart(data);
		}
	},

	createHighchart: function(data) {
		this.setOptions();

		var title = '';
		if (data.resource) {
			title = data.resource + ' on ';
		}
		if (data.component) {
			title += data.component;
		}
		this.options.title.text = title;

		log.debug(" + set title: '" + title + "'", this.logAuthor);

		if (data.perf_data_array) {
			this.parseData(data.perf_data_array);
		}

		this.chart = new Highcharts.Chart(this.options);
		//this.doRefresh();
	},

	parseData: function(perf_data) {
		//log.dump(perf_data)

		if (this.chart) {
			var serie = this.chart.get('pie');
			serie.remove();
		}

		this.options.series = [];

		var serie = {
			id: 'pie',
			type: 'pie',
			data: []
		};

		if (this.metric) {
			log.debug(" + Use one metric: '" + this.metric + "'", this.logAuthor);
			metric = perf_data[this.metric];

			var metric_max = metric.max;
			if (this.metric_max) {
				log.debug(' + Set max to: ' + this.metric_max, this.logAuthor);
				metric_max = this.metric_max;
			}

			serie.data.push({ name: 'Free', y: metric_max - metric.value, color: global.default_colors[0] });
			serie.data.push({ name: metric.metric, y: metric.value, color: global.default_colors[1]});
		}else {
			log.debug(' + Use Multiple metrics', this.logAuthor);
			var index;
			var total = 0;
			for (index in perf_data) {
				metric = perf_data[index];
				total += metric.value;
			}
			if (total == 0) { total = 1 }
			log.debug('   + Total: ' + total, this.logAuthor);

			var i = 0;
			for (index in perf_data) {
				log.debug("   + Push metric: '" + index + "'", this.logAuthor);
				var metric = perf_data[index];
				var color = global.default_colors[i];
				serie.data.push({ name: metric.metric, y: Math.round(metric.value / total), color: color });
				i += 1;
			}
		}

		this.options.series.push(serie);
		if (this.chart) {
			this.chart.addSeries(serie);
		}
	},

	setOptions: function() {
		this.options = {
			chart: {
				renderTo: this.divId,
				defaultSeriesType: 'pie',
				height: this.divHeight,
				animation: false,
				borderColor: '#FFFFFF'
			},
			exporting: {
				enabled: false
			},
			colors: [],
			plotOptions: {
				pie: {
					allowPointSelect: true,
					cursor: 'pointer',
					dataLabels: {
						enabled: false
					},
					showInLegend: true,
					animation: false
				}
			},
			tooltip: {
				formatter: function() {
					return '<b>' + this.point.name + '</b>: ' + Math.round(this.percentage) + ' %';
					}
			},
			title: {
				text: '',
				floating: true
			},
			symbols: [],
			credits: {
				enabled: false
			},
			series: []
		};

		//specifique options to add
		if (this.exportMode) {
			this.options.plotOptions.pie.enableMouseTracking = false;
			this.options.plotOptions.tooltip = {};
			this.options.plotOptions.pie.shadow = false;
		}

	}



});
