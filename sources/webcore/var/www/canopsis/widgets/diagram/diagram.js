//need:app/lib/view/cperfstoreValueConsumerWidget.js
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
Ext.define('widgets.diagram.diagram', {
	extend: 'canopsis.lib.view.cperfstoreValueConsumerWidget',
	alias: 'widget.diagram',

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

	labels: false,
	gradientColor: false,

	interval: global.commonTs.hours,
	aggregate_method: 'LAST',
	aggregate_interval: 0,
	aggregate_max_points: 1,
	aggregate_round_time: true,

	nb_node: 0,
	hide_other_column: false,

	diagram_type: 'pie',

	nameInLabelFormatter: false,
	pctInLabel: true,

	useLastRefresh: false,

	labelFormatter: function() {
		if(this.y === 0) {
			return;
		}

		var me = this.series.chart.options.cwidget();

		var prefix = "";

		var formatter = function(value, unit) {
			if(me.humanReadable) {
				return rdr_humanreadable_value(value, unit);
			}
			else {
				if(unit) {
					return value + ' ' + unit;
				}
				else {
					return value;
				}
			}
		};

		if(me.nameInLabelFormatter) {
			if(this.x) {
				prefix = '<b>' + this.x + ':</b> ';
			}
			else {
				prefix = '<b>' + this.point.metric + ':</b> ';
			}
		}

		if(me.pctInLabel && this.percentage !== undefined) {
			return prefix + formatter(this.percentage, "%");
		}
		else {
			return prefix + formatter(this.y, this.point.bunit);
		}
	},

	initComponent: function() {
		this.callParent(arguments);

		this.logAuthor = '[widgets][diagram]';

		this.backgroundColor        = check_color(this.backgroundColor);
		this.borderColor            = check_color(this.borderColor);
		this.legend_fontColor       = check_color(this.legend_fontColor);
		this.legend_borderColor     = check_color(this.legend_borderColor);
		this.legend_backgroundColor = check_color(this.legend_backgroundColor);

		//retrocompatibility
		if(Ext.isArray(this.nodes)) {
			this.nodesByID = parseNodes(this.nodes);
		}
		else {
			this.nodesByID = expandAttributs(this.nodes);
		}

		this.nb_node = Ext.Object.getSize(this.nodesByID);

		this.series_array = [];
		this.series_list = undefined;
		this.categories = [];
		this.nodesByMetricAndCategory = {};

		log.debug('nodesByID:', this.logAuthor);
		log.dump(this.nodesByID);

		Ext.Object.each(this.nodesByID, function(id, node) {
			// initialize categories
			if(node.category) {
				if(Ext.Array.indexOf(this.categories, node.category) === -1) {
					this.categories.push(node.category);
				}

				this.nodesByMetricAndCategory[id] = {
					category: node.category,
					metric: node.label
				};
			}
		}, this);

		// at least one metric is categorized
		if(this.categories.length > 0) {
			// put all other metrics in default category
			Ext.Object.each(this.nodesByID, function(id, node) {
				if(!node.category) {
					node.category = 'unaffected';

					if(Ext.Array.indexOf(this.categories, node.category) === -1) {
						this.categories.push(node.category);
					}

					this.nodesByMetricAndCategory[id] = {
						category: node.category,
						metric: node.label
					};
				}
			}, this);
		}

		//Set title
		if(this.autoTitle) {
			this.setchartTitle();
			this.title = '';
		}
		else if(!this.border) {
			this.chartTitle = this.title;
			this.title = '';
		}
	},

	afterContainerRender: function() {
		log.debug('Initialize Pie', this.logAuthor);

		this.setOptions();
		this.createChart();

		this.ready();
	},

	setchartTitle: function() {
		var title = '';

		if(this.nb_node === 1) {
			var firstKey = Ext.Object.getKeys(this.nodesByID)[0];
			var firstNode = this.nodesByID[firstKey];
			var component = firstNode.component;
			var source_type = firstNode.source_type;

			if(source_type === 'resource') {
				var resource = firstNode.resource;
				title = resource + ' ' + _('on') + ' ' + component;
			}
			else {
				title = component;
			}
		}

		this.chartTitle = title;
	},

	setOptions: function() {
		var me = this;

		this.options = {
			cwidget: function() {
				return me;
			},

			chart: {
				renderTo: this.wcontainerId,
				defaultSeriesType: 'pie',
				height: this.getHeight(),
				reflow: false,
				animation: false,
				borderColor: this.borderColor,
				borderWidth: this.borderWidth,
				backgroundColor: this.backgroundColor,
				inverted: (this.diagram_type === 'column') ? this.verticalDisplay : false
			},
			colors: [],
			plotOptions: {
				pie: {
					allowPointSelect: true,
					cursor: 'pointer',
					dataLabels: {
						enabled: this.labels,
						color: '#000000',
						connectorColor: '#000000',
						formatter: this.labelFormatter
					},
					showInLegend: true,
					animation: false,
					size: this.pie_size + '%'
				},
				column: {
					animation: false,
					dataLabels: {
						enabled: this.labels,
						color: '#000000',
						connectorColor: '#000000',
						formatter: this.labelFormatter
					}
				}
			},
			tooltip: {
				formatter: this.tooltip_formatter
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
				enabled: (this.diagram_type === 'column') ? false : this.legend,
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
			series: [],
			xAxis: {
				title: { text: null },
				labels: {enabled: (this.nameInLabelFormatter) ? false : true}
			},
			yAxis: {
				title: { text: null },
				labels: {formatter: this.y_formatter}
			}
		};

		//specifique options to add
		if(this.exportMode) {
			this.options.plotOptions.pie.enableMouseTracking = false;
			this.options.plotOptions.tooltip = {};
			this.options.plotOptions.pie.shadow = false;
		}

		if(this.diagram_type === 'column' && this.categories.length > 0) {
			this.options.xAxis.categories     = this.categories;
			this.options.xAxis.labels.enabled = true;
			this.options.legend.enabled       = this.legend;
		}
	},

	createChart: function() {
		this.chart = new Highcharts.Chart(this.options);
		Highcharts.setOptions({
			lang: {
				months: [_('January'), _('February'), _('March'), _('April'), _('May'), _('June'), _('July'), _('August'), _('September'), _('October'), _('November'), _('December')],
				weekdays: [_('Sunday'), _('Monday'), _('Tuesday'), _('Wednesday'), _('Thursday'), _('Friday'), _('Saturday')],
				shortMonths: [_('Jan'), _('Feb'), _('Mar'), _('Apr'), _('May'), _('Jun'), _('Jul'), _('Aug'), _('Sept'), _('Oct'), _('Nov'), _('Dec')]
			}
		});
	},

	processPostParams: function(post_params) {
		post_params['aggregate_max_points'] = 1;
	},

	doRefresh: function(from, to) {
		log.debug('Get values from ' + new Date(from) + ' to ' + new Date(to), this.logAuthor);

		this.refreshNodes(from, to);
	},

	onRefresh: function(data) {
		var myEl = this.getEl();
		var i = undefined;

		if(this.chart && data.length > 0) {
			if(myEl && myEl.isMasked && !this.isDisabled()) {
				myEl.unmask();
			}

			// Remove old series
			this.removeSerie();

			var other_unit  = '';
			var series_list = undefined;
			var serie_conf  = undefined;
			var serie       = undefined;
			var info        = undefined;
			var node        = undefined;

			if(this.diagram_type === 'column' && this.categories.length > 0) {
				var j = 0;

				series_list = {};

				for(i = 0; i < data.length; i++) {
					info = data[i];
					node = this.nodesByID[info.node];

					var metric   = this.nodesByMetricAndCategory[info.node].metric;
					var category = this.nodesByMetricAndCategory[info.node].category;

					if(!series_list[metric]) {
						series_list[metric] = this.getSerie(data, metric);
						j++;
					}

					if(!series_list[metric].data) {
						series_list[metric].data = [];
					}

					if(node.label) {
						info.metric = node.label;
					}

					serie_conf = this.getSerieConf(info, node, j);
					serie_conf.category = category;

					if(series_list[metric].data.length > 0) {
						serie_conf.color = undefined;
					}

					var idcat = Ext.Array.indexOf(this.categories, category);
					series_list[metric].data[idcat] = serie_conf;

					// Make sure there is no undefined field
					for(var k = 0; k < series_list[metric].data.length; k++) {
						if(series_list[metric].data[k] === undefined) {
							// When there is no data, Highcharts expect null object
							series_list[metric].data[k] = null;
						}
					}

					if(!series_list[metric].name) {
						series_list[metric].name = serie_conf.name;
					}

					if(!series_list[metric].color) {
						series_list[metric].color = serie_conf.color;
					}
				}
			}
			else {
				serie = this.getSerie(data);

				for(i = 0; i < data.length; i++) {
					info = data[i];
					node = this.nodesByID[info.node];

					if(node.label) {
						data[i].metric = node.label;
					}

					serie_conf = this.getSerieConf(info, node, i);
					serie.data.push(serie_conf);
				}
			}

			if(this.setAxis && this.diagram_type === 'column' && this.series_list === undefined) {
				if(series_list !== undefined && Ext.Object.getSize(series_list) > 0) {
					var first_serie = series_list[Object.keys(series_list)[0]];

					this.setAxis(first_serie.data);
				}
				else {
					this.setAxis(serie.data);
				}
			}

			if(data.length === 1 && !this.hide_other_column && this.diagram_type === 'pie' && serie_conf._max) {
				var other_label = '<b>' + this.other_label + '</b>' + other_unit;
				var colors = global.curvesCtrl.getRenderColors(this.other_label, 1);

				var color = (this.gradientColor ? this.getGradientColor(colors[0]) : colors[0]);

				serie.data.push({
					id: 'pie_other',
					name: other_label,
					metric: this.other_label,
					y: serie_conf._max - value,
					color: color
				});
			}

			if((serie && serie.data) || Ext.Object.getSize(series_list) > 0) {
				if(series_list !== undefined && Ext.Object.getSize(series_list) > 0) {
					this.series_list = series_list;
				}
				else {
					this.serie = serie;
				}

				this.displaySerie();
			}
			else {
				log.debug('No data to display', this.logAuthor);
			}
		}
		else {
			myEl.mask(_('No data on interval'));
		}

	},

	removeSerie: function() {
		var serie = undefined;

		if(this.series_array.length > 0) {
			for(var i = 0; i < this.series_array.length; i++) {
				serie = this.chart.get(this.series_array[i]);

				if(serie) {
					serie.destroy();
				}
			}
		}
		else {
			serie = this.chart.get('serie');

			if(serie) {
				serie.destroy();
			}
		}
	},

	displaySerie: function() {
		if(this.series_list !== undefined && Ext.Object.getSize(this.series_list) > 0) {
			Ext.Object.each(this.series_list, function(id, serie) {
				void(id);

				this.chart.addSeries(Ext.clone(serie));
			}, this);
		}
		else if(this.serie) {
			this.chart.addSeries(Ext.clone(this.serie));
		}
	},

	reloadSerie: function() {
		this.removeSerie();
		this.displaySerie();
	},

	getSerie: function(data, metric) {
		var bunit = undefined;

		if(data.length > 0) {
			for(var i = 0; i < data.length; i++) {
				if(data[i].bunit) {
					bunit = data[i].bunit;
				}
			}
		}

		if(metric === undefined) {
			return {
				id: 'serie',
				type: this.diagram_type,
				shadow: false,
				data: [],
				bunit: bunit
			};
		}
		else {
			this.series_array.push('serie_' + metric);

			return {
				id: 'serie_' + metric,
				type: this.diagram_type,
				shadow: false,
				data: [],
				bunit: bunit
			};
		}
	},

	getSerieConf: function(info, node, i) {
		var metric = undefined;
		var value  = undefined;
		var unit   = undefined;
		var max    = undefined;

		if(info !== undefined) {
			metric = info.metric;
			unit   = info.bunit;
			max    = info.max;

			if(info.values !== undefined && info.values.length >= 1) {
				value = info.values[0][1];
			}
			else {
				value = 0;
			}
		}

		if(!max) {
			max = this.max;
		}

		if(unit === '%' && !max) {
			max = 100;
		}

		var colors = global.curvesCtrl.getRenderColors(metric, i);
		var curve  = global.curvesCtrl.getRenderInfo(metric);

		// Set label
		var label = undefined;

		if(curve) {
			label = curve.get('label');
		}

		if(!label) {
			label = metric;
		}

		metric = label;

		var metric_long_name = '<b>' + label + '</b>';

		if(unit) {
			metric_long_name += ' (' + unit + ')';
		}

		var color = (node !== undefined && node.curve_color ? node.curve_color : colors[0]);
		color = (this.gradientColor ? this.getGradientColor(color) : color);

		return {
			id: metric,
			name: metric_long_name,
			metric: metric,
			y: value,
			color: color,
			bunit: unit,
			_max: max
		};
	},

	getGradientColor: function(color) {
		return {
			radialGradient: {
				cx: 0.5,
				cy: 0.3,
				r: 0.7
			},
			stops: [
				[0, color],
				[1, Highcharts.Color(color).brighten(-0.3).get('rgb')]
			]
		};
	},

	onResize: function() {
		log.debug('onRezize', this.logAuthor);

		if(this.chart) {
			this.chart.setSize(this.getWidth(), this.getHeight() , false);
			this.reloadSerie();
		}
	},

	tooltip_formatter: function() {
		var me = this.series.chart.options.cwidget();

		var formatter = function(options, value) {
			if(options.invert) {
				value = -value;
			}

			if(me.humanReadable) {
				value = rdr_humanreadable_value(value, options.bunit);
			}
			else if (options.bunit) {
				value = value + ' ' + options.bunit;
			}

			return '<b>' + options.metric + '</b>: ' + value;
		};

		var s = '';

		if(this['points']) {
			// Shared
			$.each(this.points, function(i, point) {
				void(i);

				s += formatter(point.options, point.y);
			});
		}
		else {
			s += formatter(this.point.options, this.y);
		}

		return s;
	},

	setAxis: function(data) {
		var metrics = [];

		for(var i = 0; i < data.length; i++) {
			if(!data[i]) {
				metrics.push('');
			}
			else if(data[i].metric) {
				metrics.push(data[i].metric);
			}
		}

		if(this.categories.length === 0) {
			this.chart.xAxis[0].setCategories(metrics, false);
		}
	},

	y_formatter: function() {
		var me = this.chart.options.cwidget();

		if(this.chart.series.length) {
			var bunit = this.chart.series[0].options.bunit;

			if(me.humanReadable) {
				return rdr_humanreadable_value(this.value, bunit);
			}
			else if(bunit) {
				return this.value + ' ' + bunit;
			}
		}

		return this.value;
	}
});
