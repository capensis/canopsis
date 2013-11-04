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

flag_tootlip_template = Ext.create('Ext.XTemplate',
	'<table>',
		'<tr>',
			'<td style="margin:3px;">',
				'<tpl if="event_type == \'user\'">',
					'<img src="widgets/stream/logo/ui.png" style="width: 32px;"></img>',
				'</tpl>',
				'<tpl if="event_type != \'user\'">',
					'<img src="widgets/stream/logo/{icon}.png" style="width: 32px;"></img>',
				'</tpl>',
			'</td>',
			'<td>',
				'<div style="margin:3px;">',
					'<tpl if="display_name">',
							'<b>{display_name}</b>',
					'</tpl>',
					'<tpl if="display_name == undefined">',
						'<b>{component}</b>',
						'<tpl if="resource">',
							'<b> - {resource}</b>',
						'</tpl>',
					'</tpl>',
					' <span style="color:grey;font-size:10px">{date}</span>',
					'<br/>{text}',
				'</div>',
			'</td>',
		'</tr>',
	'</table>',
	{compiled: true}
);

Ext.define('widgets.line_graph.line_graph', {
	extend: 'canopsis.lib.view.cperfstoreValueConsumerWidget',

	alias: 'widget.line_graph',

	layout: 'fit',

	first: false,
	pushPoints: false,

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

	chart_type: 'line',

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
	aggregate_round_time: true,

	SeriesType: 'area',

	SeriePercent: false,
	lineWidth: 1,

	//trends
	trend_lines: false,
	trend_lines_type: 'ShortDot',

	last_values: [],

	nb_node: 0,
	same_node: true,
	displayLastValue: false,

	consolidation_method: undefined,

	timeNav: false,
	timeNav_window: 604800,
	onDoRefresh: false,

	nbMavEventsDisplayed: 100,

	autoShift: false,
	lastShift: undefined,

	initComponent: function() {
		this.callParent(arguments);

		this.logAuthor = '[widgets][line_graph]';

		this.backgroundColor        = check_color(this.backgroundColor);
		this.borderColor            = check_color(this.borderColor);
		this.legend_fontColor       = check_color(this.legend_fontColor);
		this.legend_borderColor     = check_color(this.legend_borderColor);
		this.legend_backgroundColor = check_color(this.legend_backgroundColor);

		this.lastYaxis = 0;
		this.bunitYaxis = {};

		this.OverlayLegend = [];

		log.debug('nodes:', this.logAuthor);
		log.dump(this.nodes);

		//retro compatibility
		if(Ext.isArray(this.nodes)) {
			this.nodesByID = parseNodes(this.nodes);
		}
		else {
			this.nodesByID = expandAttributs(this.nodes);
		}

		this.nb_node = Ext.Object.getSize(this.nodesByID);

		// Check if same node
		if(this.nb_node === 1) {
			this.same_node = true;
		}
		else {
			var flag = undefined;

			Ext.Object.each(this.nodesByID, function(key, value) {
				void(key);

				var node = value['resource'] + value['component'];

				if(flag === undefined) {
					flag = node;
				}
				else if(flag !== node) {
					this.same_node = false;
				}
			}, this);
		}

		log.debug('nodesByID:', this.logAuthor);
		log.dump(this.nodesByID);
		log.debug('same_node: ' + this.same_node, this.logAuthor);

		if(this.timeNav && this.exportMode) {
			this.timeNav = false;
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
		log.debug('Initialize ' + this.chart_type + '_graph', this.logAuthor);
		log.debug(' + Time window: ' + this.time_window, this.logAuthor);

		this.series = {};

		this.setOptions();
		this.createChart();

		this.ready();
	},

	setchartTitle: function() {
		var title = '';

		if(this.nb_node && this.same_node) {
			var firstKey = Ext.Object.getKeys(this.nodesByID)[0];
			var firstNode = this.nodesByID[firstKey];
			var component = firstNode.component;
			var resource = undefined;

			try {
				resource = firstNode.resource;
			}
			catch(err) {
				;
			}

			if(resource) {
				title = resource + ' ' + _('on') + ' ' + component;
			}
			else {
				title = component;
			}
		}

		this.chartTitle = title;
	},

	setOptions: function() {
		this.options = {
			reportMode: this.reportMode,

			cwidget: this,

			chart: {
				renderTo: this.wcontainerId,
				defaultSeriesType: this.SeriesType,
				height: this.getHeight(),
				reflow: false,
				animation: false,
				borderColor: this.borderColor,
				borderWidth: this.borderWidth,
				backgroundColor: this.backgroundColor
			},
			exporting: {
				enabled: (this.exportMode || this.reportMode) ? false : this.exporting_enabled,
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
			navigator: {
				enabled: this.timeNav,
				baseSeries: 0,
				height: 20
			},
			rangeSelector: {
				enabled: this.timeNav,
				buttons: [{
					type: 'hour',
					count: 1,
					text: _('1h')
				}, {
					type: 'day',
					count: 1,
					text: _('1d')
				}, {
					type: 'week',
					count: 1,
					text: _('1w')
				}],
				inputEnabled: false,
				selected: 1
			},
			scrollbar: {
				enabled: this.timeNav
			},
			colors: [],
			title: {
				text: this.chartTitle,
				floating: true,
				style: {
					fontSize: this.title_fontSize
				}
			},
			tooltip: {
				shared: this.tooltip_shared,
				crosshairs: this.tooltip_crosshairs,
				enabled: this.tooltip,
				formatter: this.tooltip_formatter,
				useHTML: true
			},
			xAxis: {
				id: 'timestamp',
				type: 'datetime',
				tickmarkPlacement: 'on',
				events: {
					afterSetExtremes: this.afterSetExtremes
				}
			},
			loading: {
				labelStyle: {
					top: '2em'
				}
			},
			yAxis: [
				{
					id: 'state',
					title: { text: null },
					labels: { enabled: false },
					max: 100
				},{
					title: { text: null },
					labels: {
						formatter: this.y_formatter
					}
				},{
					title: { text: null },
					labels: {
						formatter: this.y_formatter
					},
					opposite: true
				},{
					title: { text: null },
					labels: {
						formatter: this.y_formatter
					}
				},{
					title: { text: null },
					labels: {
						formatter: this.y_formatter
					},
					opposite: true
				}
			],
			plotOptions: {
				series: {
					animation: false,
					shadow: false
				},
				column: {}
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

		//graph type (for column)
		if(this.chart_type === 'column') {
			this.options.chart.type = this.chart_type;
		}

		// Check marker
		var marker_enable = false;
		if(this.marker_symbol) {
			marker_enable = true;
		}
		else {
			this.marker_symbol = null;
			this.marker_radius = 0;
		}

		// Configure line type
		this.options.plotOptions[this.SeriesType] = {
			lineWidth: this.lineWidth,
			shadow: false,
			cursor: 'pointer',
			turboThreshold: 10,
			marker: {
				enabled: marker_enable,
				symbol: this.marker_symbol,
				radius: this.marker_radius
			}
		};

		//specifique options to add
		if(this.exportMode) {
			this.options.plotOptions.series['enableMouseTracking'] = false;
		}
		else if(this.zoom) {
			this.options.chart.zoomType = 'x';
		}

		//Time Navigation
		if(this.timeNav) {
			var now = Ext.Date.now();
			var data = [[now - (this.timeNav_window * 1000), null], [now, null]];

			this.options.series.push({
				id: 'timeNav',
				name: 'timeNav',
				data: data,
				showInLegend: false
			});

			// Disable legend, see: https://github.com/highslide-software/highcharts.com/issues/567
			this.options.legend.enabled = false;

			this.options.xAxis['min'] = now - (this.time_window * 1000);
			this.options.xAxis['max'] = now;
		}

		// Update axis color with curve color
		for(var id in this.nodesByID) {
			var node = this.nodesByID[id];

			if(!node.yAxis) {
				continue;
			}

			var axis_color = node.curve_color || node.area_color || undefined;

			Ext.merge(this.options.yAxis[node.yAxis], {
				labels: {
					style: {
						color: axis_color
					}
				},
				title: {
					text: (!this.options.legend.enabled ? node.label : null),
					style: {
						color: axis_color
					}
				}
			});
		}
	},

	y_formatter: function() {
		var me = this.chart.options.cwidget;

		if(this.axis.series.length) {
			var bunit = this.axis.series[0].options.bunit;

			if(me.humanReadable) {
				return rdr_humanreadable_value(this.value, bunit);
			}
			else {
				if(bunit) {
					return this.value + ' ' + bunit;
				}
				else {
					return this.value;
				}
			}
		}

		return rdr_yaxis(this.value);
	},

	tooltip_formatter: function() {
		if(this.point && this.point._flag === true) {
			return flag_tootlip_template.applyTemplate(this.point);
		}
		else {
			var me;

			if(this['points']) {
				me = this.points[0].series.chart.options.cwidget;
			}
			else {
				me = this.series.chart.options.cwidget;
			}

			var formatter = function(options, value) {
				if(options.invert) {
					value = - value;
				}

				if(me.humanReadable) {
					value = rdr_humanreadable_value(value, options.bunit);
				}
				else if(options.bunit) {
					value = value + ' ' + options.bunit;
				}

				return '<b>' + options.metric + ':</b> ' + value;
			};

			var s = '<b>' + rdr_tstodate(this.x / 1000) + '</b>';

			if(this['points']) {
				// Shared
				$.each(this.points, function(i, point) {
					void(i);

					s += '<br/>' + formatter(point.series.options, point.y);
				});
			}
			else {
				s += '<br/>' + formatter(this.series.options, this.y);

				if(this.series.eta){
					var eta = this.series.eta;
					var dEta = parseInt(Ext.Date.now()) + eta;
					s += '<br/><b>ETA:</b> ' + rdr_tstodate(dEta/1000) + ' (' + rdr_duration(eta/1000, 2) + ')';
				}
			}

			return s;
		}
	},

	createChart: function() {
		this.chart = new Highcharts.Chart(this.options);
		Highcharts.setOptions({
			global: {
				useUTC: false
			},
			lang: {
				months: [_('January'), _('February'), _('March'), _('April'), _('May'), _('June'), _('July'), _('August'), _('September'), _('October'), _('November'), _('December')],
				weekdays: [_('Sunday'), _('Monday'), _('Tuesday'), _('Wednesday'), _('Thursday'), _('Friday'), _('Saturday')],
				shortMonths: [_('Jan'), _('Feb'), _('Mar'), _('Apr'), _('May'), _('Jun'), _('Jul'), _('Aug'), _('Sept'), _('Oct'), _('Nov'), _('Dec')]
			}
		});
	},

	// CORE

	doRefresh: function(from, to) {
		var now = Ext.Date.now();

		if(this.chart) {

			if(this.timeNav) {
				var time_limit = now - (this.timeNav_window * 1000);

				if(to > now) {
					to = now;
				}

				if(from < time_limit) {
					from = time_limit;
				}

				if(to <= time_limit) {
					this.chart.showLoading(_('Time is out of range') + '...');
					return;
				}

				var time_window = to - from;

				this.onDoRefresh = true;

				var serie = this.chart.get('timeNav');
				var e = serie.xAxis.getExtremes();
				time_window = e.max - e.min;

				if(this.reportMode) {
					this.stopTask();
					serie.xAxis.setExtremes(from, to, false);
				}
				else {
					serie.xAxis.setExtremes(now - time_window, now, false);
				}
			}

			this.refreshNodes(from, to);

			if(this.flagFilter) {
				var filter = [{
					'timestamp': { '$gte': parseInt(from / 1000), '$lte': parseInt(to / 1000) }
				},{
					'state_type': 1
				}];

				var decodedFilter = Ext.decode(this.flagFilter);
				var decodedFilterKeys = Ext.Object.getKeys(decodedFilter);

				if(!decodedFilterKeys[0] || decodedFilterKeys[0] === 'null') {
					return;
				}

				filter.push(decodedFilter);

				Ext.Ajax.request({
					url: '/rest/events_log',
					scope: this,
					params: {
						filter: Ext.encode({'$and': filter}),
						limit: this.nbMavEventsDisplayed
					},
					method: 'Get',
					success: function(response) {
						var data = Ext.JSON.decode(response.responseText).data;
						this.addFlagSerie(data);
					},
					failure: function(result, request) {
						void(result);

						log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
					}
				});
			}
		}
	},

	onRefresh: function(data) {
		if(this.chart) {
			log.debug('On refresh', this.logAuthor);

			this.clearGraph();

			var toggle_max_percent = false;

			if(data.length > 0) {
				this.last_values = [];

				//percent check
				if(this.SeriePercent) {
					if(data.max) {
						toggle_max_percent = true;
					}
					else {
						toggle_max_percent = false;
					}
				}

				for(var i = 0; i < data.length; i++) {
					this.addDataOnChart(data[i]);

					var node_id = data[i].node;
					var node = this.nodesByID[node_id];

					// TODO: Fix with new format
					if(!node.max && data[i].max) {
						node.max = data[i].max;
					}

					// Exclude state lines and timeNav
					if(!this.timeNav
						&& data[i]['metric'] !== 'cps_state'
						&& data[i]['metric'] !== 'cps_state_ok'
						&& data[i]['metric'] !== 'cps_state_warn'
						&& data[i]['metric'] !== 'cps_state_crit'
						&& node.trend_curve) {

						//add/refresh trend lines
						this.addTrendLines(data[i]);
					}

					if(data[i]['values']) {
						var last_value = data[i]['values'][data[i]['values'].length - 1][1];

						this.last_values.push([
							node.label,
							last_value,
							node.bunit
						]);
					}
				}

				if(this.displayLastValue && this.last_values) {
					this.drawLastValue(this.last_values);
				}

				//Disable no data message
				if(this.chartMessage) {
					this.chartMessage.destroy();
					this.chartMessage = undefined;
				}

				//set max
				if(this.SeriePercent && toggle_max_percent) {
					this.chart.yAxis[1].setExtremes(0, 100, false);
				}

				this.chart.hideLoading();

				if(this.autoShift) {
					this.shift();
				}

				this.chart.redraw();

			}
			else {
				log.debug(' + No data', this.logAuthor);

				// if report, cleaning the chart
				if(this.reportMode === true) {
					this.clearGraph();
				}

				if(this.reportMode === true || this.series.length === 0) {
					this.chart.showLoading(_('Unfortunately, there is no data for this period'));
				}

				this.chart.redraw();
			}
		}
	},

	clearGraph: function() {
		for(var i = 0; i < this.chart.series.length; i++) {
			var name = this.chart.series[i].name;

			if(name !== 'timeNav' && name !== 'Navigator') {
				log.debug('Cleaning serie: ' + name, this.logAuthor);
				this.chart.series[i].setData([], false);
			}
		}
	},

	shift: function(tolerance) {
		if(tolerance === undefined) {
			if(this.aggregate_interval) {
				tolerance = this.aggregate_interval * 1000;
			}
			else {
				tolerance = 0;
			}
		}

		if(this.options && this.options.cwidget) {
			me = this.options.cwidget;
		}
		else {
			me = this;
		}

		var now = Ext.Date.now();

		if(!me.lastShift) {
			me.lastShift = now;
		}

		if(me.chart.series.length > 0 && now > (me.lastShift + 50000)) {
			log.debug('Check shifting (' + me.chart.series.length + ' series):', me.logAuthor);

			var timestamp = now - (me.time_window * 1000) - (me.time_window_offset * 1000) - tolerance;

			for(var i = 0; i < me.chart.series.length; i++) {
				var serie = me.chart.series[i];

				log.debug(' + ' + serie.name + ', ' + serie.data.length + ' point(s)', me.logAuthor);

				// Don't shift timeNav or short serie
				if(serie.name === 'timeNav' || serie.name === 'Navigator'
				   || (serie.data.length <= 2 && serie.name !== 'Flags')) {
					continue;
				}

				var fpoint = serie.data[0];
				var removed = 0;

				while(serie.data.length && fpoint.x < timestamp) {
					fpoint.remove(false, false);
					fpoint = serie.data[0];
					removed += 1;
				}

				log.debug('   - ' + removed + ' point(s) removed', me.logAuthor);
			}

			me.lastShift = now;
		}
	},

	onResize: function() {
		log.debug('onRezize', this.logAuthor);

		if(this.chart) {
			this.chart.setSize(this.getWidth(), this.getHeight() , false);

			if (this.displayLastValue && this.last_values) {
				this.drawLastValue(this.last_values);
			}
		}
	},

	dblclick: function() {
		if(this.chart && !this.isDisabled()) {
			this.chart.zoomOut();
		}
	},

	getSerie: function(node_id, metric_name, bunit, min, max, yAxis) {
		var serie_id = node_id + '.' + metric_name;

		var serie = this.chart.get(serie_id);

		if(serie) {
			return serie;
		}

		var node = this.nodesByID[node_id];

		log.debug('  + Create Serie:', this.logAuthor);

		if(bunit === null) {
			bunit = undefined;
		}

		if(this.SeriePercent && max > 0) {
			bunit = '%';
		}

		if(node.bunit) {
			bunit = node.bunit;
		}

		var _yAxis = node.yAxis;

		if(Ext.isNumber(yAxis)) {
			_yAxis = yAxis;
		}

		if(_yAxis === undefined || _yAxis === "default" || _yAxis === "") {
			_yAxis = this.getYaxis(bunit);
		}

		var serie_index = this.chart.series.length;

		log.debug('    + serie id: ' + serie_id, this.logAuthor);
		log.debug('    + serie index: ' + serie_index, this.logAuthor);
		log.debug('    + metric_name: ' + metric_name, this.logAuthor);
		log.debug('    + bunit: ' + bunit, this.logAuthor);
		log.debug('    + yAxis: ' + _yAxis, this.logAuthor);

		var metric_long_name = '';

		if(!this.same_node && !this.consolidation_method && node && (!node.label)) {
			metric_long_name = node.component;

			if(node.source_type === 'resource') {
				metric_long_name += ' - ' + node.resource;
			}

			metric_long_name = '(' + metric_long_name + ') ';
		}

		var colors = global.curvesCtrl.getRenderColors(metric_name, serie_index);
		var curve = global.curvesCtrl.getRenderInfo(metric_name);

		// Set Label
		var label = undefined;

		if(node.label) {
			label = node.label;
		}

		if(curve) {
			label = curve.get('label');
		}

		if(!label) {
			label = metric_name;
		}

		log.debug('    + label: ' + label, this.logAuthor);

		metric_long_name += '<b>' + label + '</b>';

		if(bunit) {
			metric_long_name += ' (' + bunit + ')';
		}

		log.debug('    + legend: ' + metric_long_name, this.logAuthor);
		log.debug('    + color: ' + colors[0], this.logAuthor);

		var _color = colors[0];

		if(node.curve_color) {
			_color = check_color(node.curve_color);
		}

		serie = {
			id: serie_id,
			name: metric_long_name,
			metric: label,
			data: [],
			color: _color,
			min: min,
			max: max,
			yAxis: _yAxis,
			bunit: bunit,
			last_timestamp: undefined
		};

		if(curve) {
			serie['dashStyle'] = curve.get('dashStyle');
			serie['invert'] = curve.get('invert');
		}

		//gets text instead value
		var type = undefined;

		if(node.curve_type !== undefined) {
			type = node.curve_type.toLowerCase();
		}

		if(type && type !== 'default') {
			serie.type = type;
		}

		if(type === 'area' || ( type === 'default' && this.SeriesType === 'area')) {
			if(node.area_color) {
				serie['fillColor'] = check_color(node.area_color);
			}
			else if(curve) {
				serie['fillColor'] = colors[1];
				serie['fillOpacity'] = colors[2] / 100;
				serie['zIndex'] = curve.get('zIndex');
			}
		}

		this.series[serie_id] = serie;

		this.chart.addSeries(Ext.clone(serie), false, false);
		var hcserie = this.chart.get(serie_id);

		return hcserie;
	},

	parseValues: function(serie, values, type) {
		//MAKE A BETTER LOOP, JUST FOR TEST
		for(var i = 0; i < values.length; i++) {
			values[i][0] = values[i][0] * 1000;

			if(this.SeriePercent && serie.options.max > 0) {
				values[i][1] = getPct(values[i][1], serie.options.max);
			}

			if(serie.options.invert) {
				values[i][1] = -values[i][1];
			}
		}

		//type specifique parsing
		if(type === 'COUNTER' && !this.aggregate_interval && !this.reportMode && serie.data.length) {
			var last_point = serie.data[serie.data.length - 1];
			var last_value = last_point.y;
			var new_values = [];

			for(i = 0; i < values.length; i++) {
				if(values[i][1] !== 0) {
					new_values.push([values[i][0], last_value + values[i][1]]);
				}
			}

			values = new_values;
		}

		return values;
	},

	addPlotlines: function(metric_name, value, color) {
		var curve = global.curvesCtrl.getRenderInfo(metric_name);
		var label = undefined;
		var zindex = 10;
		var width = 2;
		var dashStyle = 'Solid';

		if(curve) {
			label = _(curve.get('label'));
			color = global.curvesCtrl.getRenderColors(metric_name, 1)[0];
			zindex = curve.get('zIndex');
			dashStyle = curve.get('dashStyle');
		}

		if(!label) {
			label = metric_name;
		}

		this.chart.yAxis[1].addPlotLine({
			value: value,
			width: width,
			zIndex: zindex,
			color: color,
			dashStyle: dashStyle,
			label: {
				text: label
			}
		});
	},

	getYaxis: function(bunit) {
		var yaxis;

		if(bunit === undefined) {
			bunit = 'default';
		}

		yaxis = this.bunitYaxis[bunit];

		if(yaxis === undefined) {
			yaxis = this.lastYaxis + 1;
			this.bunitYaxis[bunit] = yaxis;
			this.lastYaxis = yaxis;
		}

		return yaxis;
	},

	addDataOnChart: function(data) {
		var metric_name = data['metric'];
		var values = data['values'];
		var bunit = data['bunit'];
		var node_id = data['node'];
		var min = data['min'];
		var max = data['max'];
		var type = data['type'];

		var serie = undefined;
		var value = undefined;

		if(metric_name === 'cps_state_ok' || metric_name === 'cps_state_warn' || metric_name === 'cps_state_crit') {
			serie = this.getSerie(node_id, metric_name, undefined, undefined, undefined, 0);
		}

		if(metric_name === 'cps_state') {
			var states = [0, 1, 2, 3];
			var states_data = [[], [], [], []];

			for(var i = 0; i < data['values'].length; i++) {
				state = parseInt(data['values'][i][1] / 100);

				for(var j = 0; j < states.length; j++) {
					value = 0;

					if(state === j) {
						value = 100;
					}

					states_data[j].push([data['values'][i][0], value]);
				}
			}

			for(i = 0; i < states.length; i++) {
				data['metric'] = 'cps_state_' + i;
				data['values'] = states_data[i];
				data['bunit'] = '%';
				this.addDataOnChart(data);
			}

			return true;

		}
		else {
			serie = this.getSerie(node_id, metric_name, bunit, min, max, undefined);
		}

		if(!serie) {
			log.error('Impossible to get serie, node: ' + node_id + ' metric: ' + metric_name, this.logAuthor);
			return false;
		}

		if(!serie.options) {
			log.error("Impossible to read serie's option", this.logAuthor);
			log.dump(serie);
			return false;
		}

		//Add war/crit line if on first serie
		if(this.chart.series.length === 1 && this.showWarnCritLine) {
			if(data['thld_warn']) {
				value = data['thld_warn'];

				if(this.SeriePercent && serie.options.max > 0) {
					value = getPct(value, serie.options.max);
				}

				this.addPlotlines('pl_warning', value, 'orange');
			}

			if(data['thld_crit']) {
				value = data['thld_crit'];

				if(this.SeriePercent && serie.options.max > 0) {
					value = getPct(value, serie.options.max);
				}

				this.addPlotlines('pl_critical', value, 'red');
			}

			this.showWarnCritLine = false;
		}

		var serie_id = serie.options.id;

		values = this.parseValues(serie, values, type);

		log.debug('  + Add data for ' + node_id + ', metric "' + metric_name + '" ...', this.logAuthor);

		if(values.length <= 0) {
			log.debug('   + No data', this.logAuthor);

			if(this.reportMode) {
				if(serie.visible) {
					serie.setData([], false);
					serie.hide();
				}

				return true;
			}
			else {
				return false;
			}
		}

		if(this.reportMode && !serie.visible) {
			serie.show();
		}

		if(values.length) {
			var last_timestamp = values[values.length - 1][0];

			// Timestamp of new and old point are equal, remove last point for update
			if(this.aggregate_interval && last_timestamp === this.series[serie_id]['last_timestamp'] && serie.data.length) {
				var point = serie.data[serie.data.length - 1];
				log.debug('   + Remove last point', this.logAuthor);
				point.remove(false, false);
			}

			this.series[serie_id]['last_timestamp'] = last_timestamp;
		}

		if(!this.series[serie_id].pushPoints || this.reportMode) {
			this.series[serie_id].pushPoints = true;

			log.debug('   + Set data', this.logAuthor);
			serie.setData(values, false);

		}
		else {
			log.debug('   + Push data', this.logAuthor);

			for(var idx = 0; idx < values.length; idx++) {
				value = values[idx];
				//addPoint (Object options, [Boolean redraw], [Boolean shift], [Mixed animation]) :
				serie.addPoint(value, false, false, false);
			}
		}

		return true;
	},

	getEta: function(y, a, b) {
		var now = parseInt(Ext.Date.now());
		return parseInt((y - b) / a) - now;
	},

	addTrendLines: function(data) {
		log.debug(' + Trend line', this.logAuthor);

		var serie_id = data.node + '.' + data.metric;
		var referent_serie = this.chart.get(serie_id);
		var trend_id = data.node + '.' + data.metric + '-TREND';

		var node = this.nodesByID[data.node];

		//get the trend line
		var trend_line = this.chart.get(trend_id);

		//update/create the trend line
		if(trend_line) {
			log.debug('  +  Trend line found : ' + trend_id, this.logAuthor);

			var line = [];

			for(var i = 0; i < referent_serie.data.length; i++) {
				var point = referent_serie.data[i];
				line.push([point.x, point.y]);
			}

			var reg = fitData(line);

			var y = undefined;

			if(reg.slope > 0 && node.max !== undefined) {
				y = node.max;
			}

			if(reg.slope < 0 && node.min !== undefined) {
				y = node.min;
			}

			if(y !== undefined) {
				trend_line.eta = this.getEta(y, reg.slope, reg.intercept);
			}

			line = reg.data;
			trend_line.setData(line, true);
		}
		else {
			log.debug('  +  Trend line not found : ' + trend_id, this.logAuthor);
			log.debug('  +  Create it', this.logAuthor);

			//name
			var trend_name = referent_serie.name + '-TREND';
			var curve = global.curvesCtrl.getRenderInfo(trend_name);
			var color = undefined;

			if(curve) {
				label = curve.get('label');
				color = referent_serie.color;
			}
			else {
				//check if referent curve have its own curve
				curve = global.curvesCtrl.getRenderInfo(data.metric);

				if(curve) {
					label = curve.get('label') + '-TREND';
				}
				else {
					label = trend_name;
				}
			}

			if(!color) {
				color = referent_serie.options.color;
			}

			var trend_dashStyle = 'ShortDot';

			if(this.trend_lines_type) {
				trend_dashStyle = this.trend_lines_type;
			}
			else if(curve) {
				trend_dashStyle = curve.get('dashStyle');
			}

			//serie
			var serie = {
				id: trend_id,
				type: 'line',
				name: label,
				metric: referent_serie.options.metric + ' (TREND)',
				bunit: referent_serie.options.bunit,
				yAxis: referent_serie.options.yAxis,
				data: [],
				marker: {enabled: false},
				dashStyle: trend_dashStyle
			};

			if(color) {
				serie['color'] = color;
			}

			//push the trendline in hichart, load trend data
			this.chart.addSeries(Ext.clone(serie), false, false);
			var hcserie = this.chart.get(trend_id);

			if(data.values.length > 2) {
				reg = fitData(data.values);
				y = undefined;

				if(reg.slope > 0 && node.max !== undefined) {
					y = node.max;
				}

				if(reg.slope < 0 && node.min !== undefined) {
					y = node.min;
				}

				if(y !== undefined) {
					hcserie.eta = this.getEta(y, reg.slope, reg.intercept);
				}

				line = reg.data;

				//trunc value
				line = this.truncValueArray(line);

				log.debug('  +  set data', this.logAuthor);
				hcserie.setData(line, false);
			}
			else {
				log.debug('  +  not enough data to draw trend line');
			}
		}
	},

	drawLastValue: function(values) {
		var html = '<span style="color:{0};font-size: 1.2em;">{1}</span>';

		var list_string = [];

		var bigLenght = undefined;

		// Push values
		for(var i = 0; i < values.length; i++) {
			var strvalue = values[i][2] ? values[i][2] : '';
			var str = values[i][0] + ': ' + values[i][1] + strvalue;

			// Search most longer string
			if(bigLenght === undefined || bigLenght < str.length) {
				bigLenght = str.length;
			}

			list_string.push(
				Ext.String.format(
					html,
					'dark grey',
					str
				)
			);
		}

		// Clean
		for(i = 0; i < this.OverlayLegend.length; i++) {
			this.OverlayLegend[i].destroy();
		}

		this.OverlayLegend = [];

		// Display text on chart
		var charH = 20;
		var charW = 5;

		var h = this.getHeight();
		var w = this.getWidth();

		var marginTop = h / 4;
		var marginRight = 20;

		var x = w - (bigLenght * charW) - marginRight;
		var y = marginTop;

		for(i = 0; i < list_string.length; i++) {
			var string = list_string[i];
			var chartText = this.chart.renderer.text(
				string,
				x - string.length,
				y + (i * charH)
			);

			this.OverlayLegend.push(chartText);
			chartText.add();
		}
	},

	truncValueArray: function(value_array) {
		for(var i = 0; i < value_array.length; i++) {
			value_array[i][1] = Math.floor(value_array[i][1] * 1000) / 1000;
		}

		return value_array;
	},

	addFlagSerie: function(data) {
		var serie = this.chart.get('x_flags');
		var sData = [];

		if(data) {
			for(var i = 0; i < data.length; i++) {
				var state_color = this.getStateColor(data[i].state);
				sData.push({
					x: data[i].timestamp * 1000,
					text: data[i].output,
					_flag: true,
					component: data[i].component,
					resource: data[i].resource,
					display_name: data[i].display_name,
					event_type: data[i].event_type,
					date: rdr_tstodate(data[i].timestamp),
					icon: data[i].connector.toLowerCase(),
					fillColor: state_color,
					style: {color: state_color},
					states: {
						hover: {
							fillColor: state_color
						}
					},
					title: 'A'
				});
			}

			if(serie) {
				if(this.reportMode) {
					serie.setData(sData);
				}
				else {
					for(i = 0; i < sData.length; i++) {
						serie.addPoint(sData[i], true, false);
					}
				}
			}
			else {
				serie = {
					id: 'x_flags',
					name: 'Flags',
					type: 'flags',
					data: sData,
					shape: 'circlepin',
					width: 17,
					color: 'black',
					zIndex: 2,
					showInLegend: false
				};

				this.series[serie.id] = serie;
				this.chart.addSeries(serie, true, false);
			}
		}
	},

	getStateColor: function(state) {
		if(state === 0) {
			return global.state_colors.ok;
		}
		else if(state === 1) {
			return global.state_colors.warning;
		}
		else if(state === 2) {
			return global.state_colors.critical;
		}
		else {
			return global.state_colors.unknown;
		}
	},

	afterSetExtremes: function(e) {
		var me = this.chart.options.cwidget;

		if(me.onDoRefresh) {
			me.onDoRefresh = false;
			return;
		}

		if(me.timeNav) {
			var from = Math.round(e.min, 0);
			var to = Math.round(e.max, 0);
			log.debug('Highcharts: afterSetExtremes: ' + from + ' -> ' + to, me.logAuthor);

			if(!isNaN(from) && !isNaN(to)) {
				me.chart.showLoading(_('Loading data from server') + '...');
				me.reportMode = true;
				me.doRefresh(from, to);
			}
		}
	},

 	beforeDestroy: function() {
		this.callParent(arguments);

 		if(this.chart) {
			this.chart.destroy();
			log.debug(' + Chart Destroyed', this.logAuthor);
		}
 	},

 	processPostParam: function(post_param) { // patch in waiting that shift method is reused
 		if (post_param['from'] && post_param['to']) {
 			if(this.timeNav) {
				var time_limit = (post_param['to'] - this.timeNav_window);
				post_param['from'] = (post_param['to'] - this.timeNav_window);

				if(post_param['from'] < time_limit) {
					post_param['from'] = time_limit;
				}
			}
			else {
 				post_param['from'] = (post_param['to'] - this.time_window);
 			}
 		}
 	}
});
