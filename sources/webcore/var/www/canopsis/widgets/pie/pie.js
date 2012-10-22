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

	options: {},
	chartTitle: null,
	chart: undefined,
	serie: undefined,

	//Default Options
	max: undefined,
	other_label: 'Free',

	backgroundColor: '#FFFFFF',
	borderColor: '#FFFFFF',
	borderWidth: 0,

	exporting_enabled: false,

	title_fontSize: 15,

	pie_size: 60,
	legend_verticalAlign: 'bottom',
	legend_align: 'center',
	legend_layout: 'horizontal',
	legend_backgroundColor: null,
	legend_borderColor: '#909090',
	legend_borderWidth: 1,
	legend_fontSize: 12,
	legend_fontColor: '#3E576F',
	//

	nb_node: 0,


	initComponent: function() {
		this.backgroundColor	= check_color(this.backgroundColor);
		this.borderColor	= check_color(this.borderColor);
		this.legend_fontColor	= check_color(this.legend_fontColor);
		this.legend_borderColor = check_color(this.legend_borderColor);
		this.legend_backgroundColor	= check_color(this.legend_backgroundColor);

		this.nodesByID = {};
		//Store nodes in object
		for (var i in this.nodes) {
			var node = this.nodes[i];

			//hack for retro compatibility
			if (!node.dn)
				node.dn = [node.component, node.resource];

			if (this.nodesByID[node.id]) {
				this.nodesByID[node.id].metrics.push(node.metrics[0]);
			}else {
				this.nodesByID[node.id] = Ext.clone(node);
				this.nb_node += 1;
			}
		}
		log.debug('nodesByID:', this.logAuthor);
		log.dump(this.nodesByID);

		//Set title
		if (this.autoTitle) {
			this.setchartTitle();
			this.title = '';
		}else {
			if (! this.border) {
				this.chartTitle = this.title;
				this.title = '';
			}
		}

		this.callParent(arguments);
	},

	afterContainerRender: function() {
		log.debug('Initialize Pie', this.logAuthor);

		// Clean this.nodes
		if (this.nodes)
			this.processNodes();

		this.setOptions();
		this.createChart();

		this.ready();
	},

	setchartTitle: function() {
		var title = '';
		if (this.nb_node) {
			var component = this.nodes[0].dn[0];
			var source_type = this.nodes[0].source_type;

			if (source_type == 'resource') {
				var resource = this.nodes[0].dn[1];
				title = resource + ' ' + _('pie.on') + ' ' + component;
			}else {
				title = component;
			}
		}
		this.chartTitle = title;
	},

	setOptions: function() {
		this.options = {
			chart: {
				renderTo: this.wcontainerId,
				defaultSeriesType: 'pie',
				height: this.getHeight(),
				reflow: false,
				animation: false,
				borderColor: this.borderColor,
				borderWidth: this.borderWidth,
				backgroundColor: this.backgroundColor
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
					animation: false,
					size: this.pie_size + '%'
				}
			},
			tooltip: {
				formatter: function() {
					return this.point.name + ': ' + Math.round(this.percentage * 1000) / 1000 + ' %';
					}
			},
			title: {
				text: this.chartTitle,
				floating: true,
				style: {
					fontSize: this.title_fontSize
				}
			},
			exporting: {
				enabled: this.exporting_enabled,
				filename: this.chartTitle,
				type: 'image/svg+xml',
				url: '/export_svg',
				buttons: {
					exportButton: {
						enabled: true,
						menuItems: null,
						onclick: function() {
							this.exportChart();
						}
					},
					printButton: {
						enabled: false
					}
				}
			},
			symbols: [],
			credits: {
				enabled: false
			},
			legend: {
				enabled: this.legend,
				verticalAlign: this.legend_verticalAlign,
				align: this.legend_align,
				layout: this.legend_layout,
				backgroundColor: this.legend_backgroundColor,
				borderWidth: this.legend_borderWidth,
				borderColor: this.legend_borderColor,
				itemStyle: {
					fontSize: this.legend_fontSize,
					color: this.legend_fontColor
				}
			},
			series: []
		};

		//specifique options to add
		if (this.exportMode) {
			this.options.plotOptions.pie.enableMouseTracking = false;
			this.options.plotOptions.tooltip = {};
			this.options.plotOptions.pie.shadow = false;
		}
	},

	createChart: function() {
		this.chart = new Highcharts.Chart(this.options);
	},

	processNodes: function() {
		var post_params = [];
		for (var i in this.nodes) {
			post_params.push({
				id: this.nodes[i].id,
				metrics: this.nodes[i].metrics
			});
		}
		this.post_params = {
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : 'LAST',
			'aggregate_max_points': 1
		};
	},

	doRefresh: function(from, to) {
	/*	if (Ext.Date.now() < this.lastRefresh + (this.interval*1000)){
			log.debug('Wait right interval to refresh', this.logAuthor)
			return false
		}*/

		// Get last point only
		if (this.interval) {
			to = Ext.Date.now();
			from = to - this.interval * 1000;
		}else {
			from = to;
		}

		log.debug('Get values from ' + new Date(from) + ' to ' + new Date(to), this.logAuthor);

		if (this.nodes) {
			if (this.nodes.length != 0) {

				var url = '/perfstore/values/' + from + '/' + to;

				Ext.Ajax.request({
					url: url,
					scope: this,
					params: this.post_params,
					method: 'POST',
					success: function(response) {
						var data = Ext.JSON.decode(response.responseText);
						data = data.data;
						this.onRefresh(data);
					},
					failure: function(result, request) {
						log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
					}
				});
			} else {
				log.debug('No nodes specified', this.logAuthor);
			}
		}
	},

	onRefresh: function(data) {
		if (this.chart && data.length != 0) {
			var myEl = this.getEl();
			if (myEl.isMasked && !this.isDisabled())
				myEl.unmask();

			// Remove old series
			this.removeSerie();

			var serie = {
				id: 'pie',
				type: 'pie',
				data: []
			};

			var other_unit = '';

			for (var index in data) {
				info = data[index];

				var node = info['node'];
				var metric = info['metric'];

				var value = undefined;

				//----------------Process value-----------------
			/*	if(info.type == 'COUNTER'){
					if(this.counter_value[metric]){
						value = this.counter_value[metric]
					}else{
						this.counter_value[metric] = 0
						value = 0
					}

						value += info['values'][0][1]

					log.debug('The value for ' + metric + ' is ' + value,
												this.logAuthor)
				}else{
					var info_length = info['values'].length*/
				if (info['values'].length >= 1)
					value = info['values'][0][1];

				//------------------

				var unit = info['bunit'];
				var max = info['max'];

				if (max == null)
					max = this.max;

				if (unit == '%' && ! max)
					max = 100;

				var metric_name = metric;

				var colors = global.curvesCtrl.getRenderColors(metric_name, index);
				var curve = global.curvesCtrl.getRenderInfo(metric_name);

				// Set Label
				var label = undefined;
				if (curve)
					label = curve.get('label');
				if (! label)
					label = metric_name;

				var metric_long_name = '<b>' + label + '</b>';

				if (unit) {
					metric_long_name += ' (' + unit + ')';
					other_unit += ' (' + unit + ')';
				}

				serie.data.push({ id: metric, name: metric_long_name, y: value, color: colors[0] });

			}

			if (data.length == 1) {
				var other_label = '<b>' + this.other_label + '</b>' + other_unit;
				var colors = global.curvesCtrl.getRenderColors(this.other_label, 1);
				serie.data.push({ id: 'pie_other', name: other_label, y: max - value, color: colors[0] });
			}

			if (serie.data) {
				this.serie = serie;
				this.displaySerie();
			}else {
				log.debug('No data to display', this.logAuthor);
			}
		}else {
			this.getEl().mask(_('No data on interval'));
		}

	},

	removeSerie: function() {
		var serie = this.chart.get('pie');
		if (serie)
			serie.destroy();
	},

	displaySerie: function() {
		if (this.serie)
			this.chart.addSeries(Ext.clone(this.serie));
	},

	reloadSerie: function() {
		this.removeSerie();
		this.displaySerie();
	},

	onResize: function() {
		log.debug('onRezize', this.logAuthor);
		if (this.chart) {
			this.chart.setSize(this.getWidth(), this.getHeight() , false);
			this.reloadSerie();
		}
	}

});
