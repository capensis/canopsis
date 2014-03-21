//need:app/lib/view/cperfstoreValueConsumerWidget.js
/*
 * Copyright (c) 2013 "Capensis" [http://www.capensis.com]
 *
 * This file is part of Canopsis.
 *
 * Canopsis is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Canopsis is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Canopsis.  If not, see <http://www.gnu.org/licenses/>.
 */

Ext.define('canopsis.lib.view.cwidgetGraph', {
	extend: 'canopsis.lib.view.cperfstoreValueConsumerWidget',

	logAuthor: '[widgets][graph]',
	layout: 'vbox',

	time_window: global.commonTs.day,
	time_window_offset: 0,

	lineWidth: 1,

	initComponent: function() {
		this.callParent(arguments);

		this.clearHooks();
		this.initHooks();

		this.setChartOptions();
		this.series = {};

		this.chart = undefined;

		this.on('boxready', this.createChart, this);

		this.on('resize', function() {
			if(this.chart !== undefined) {
				this.renderChart();
			}
		}, this);
	},

	initHooks: function() {
		var me = this;

		me.baseMethods = {};

		function initHook(name) {
			me.baseMethods[name] = me[name];

			me[name] = function() {
				var currentHook = {
					params: arguments,
					result: undefined,
					cancel: false
				};

				me.runHook(name, true, [currentHook]);

				if(currentHook.cancel) {
					return currentHook.result;
				}

				currentHook.result = me.baseMethods[name].apply(me, currentHook.params);

				me.runHook(name, false, [currentHook]);

				return currentHook.result;
			}
		}

		for(var hook in this.hooks) {
			initHook(hook);
		}
	},

	addHook: function(name, callback, prehook, id) {
		if(name in this.hooks) {
			log.debug('Add ' + (!!prehook ? 'pre' : 'post') + '-hook: ' + name, this.logAuthor);

			this.hooks[name][id] = {
				prehook: !!prehook,
				func: callback
			};
		}
	},

	clearHooks: function() {
		log.debug('Clear hooks', this.logAuthor);

		this.hooks = {
			insertGraphExtraComponents: {},
			setChartOptions: {},
			createChart: {},
			renderChart: {},
			destroyChart: {},
			getSerieForNode: {},
			updateSeriesConfig: {},
			doRefresh: {},
			onRefresh: {},
			dblclick: {}
		};
	},

	runHook: function(name, prehook, args) {
		if(name in this.hooks) {
			var selectPrehook = !!prehook;

			log.debug('Run ' + (!!prehook ? 'pre' : 'post') + '-hooks: ' + name, this.logAuthor);

			for(var hookId in this.hooks[name]) {
				var hook = this.hooks[name][hookId];

				if(hook.prehook === selectPrehook) {
					hook.func.apply(this, args);
				}
			}
		}
	},

	afterContainerRender: function() {
		this.callParent(arguments);

		if(this.chart !== undefined) {
			this.renderChart();
		}
	},

	setChartOptions: function() {
		var now = Ext.Date.now();

		if(this.options === undefined) {
			this.options = {};
		}

		this.options.cwidgetId = this.id;
		this.options.cwidget = function() {
			return this;
		}.bind(this);
	},

	insertGraphExtraComponents: function() {
	},

	createChart: function() {
		this.plotcontainer = $('#' + this.wcontainerId);

		this.insertGraphExtraComponents();

		/* create the main chart */
		this.chart = $.plot(this.plotcontainer, this.getSeriesConf(), this.options);

		/* initialize plugins */
		this.chart.initializeCurveStyleManager(this);

		if(this.tooltip) {
			this.chart.initializeTooltip(this);
		}

		this.chart.initializeHiddenGraphs(this);
		this.chart.initializeThresholds(this);
		this.chart.initializeDowntimes(this);
		this.chart.initializeGraphStyleManager(this);
		this.chart.initializeLegendManager(this);
		this.chart.initializeAutoScale(this);
		this.chart.initializeTimeNavigation(this);
		this.chart.initializeBarLabels(this);
	},

	renderChart: function() {
		this.chart.recomputePositions(this);

		// this.chart.setData(this.getSeriesConf());
		this.chart.setupGrid();
		this.chart.draw();
		this.add_csv_download_button();
	},

	add_csv_download_button: function() {
		//@see jqgridable for mouseover
		if(this.get_csv_data !== undefined) {
			this.get_csv_data.remove();
			this.get_csv_data = undefined;
		}

		this.get_csv_data = $('<div/>', {
			class: 'chart_button',
			text: 'download as csv'
		});

		this.get_csv_data.css({
			display: 'inline-block',
			position: 'absolute',
			top: 5,
			left: 30,
		});

		this.plotcontainer.parent().append(this.get_csv_data);

		var that = this;

		$(this.plotcontainer.parent()).mouseenter(function() {
			if(that.get_csv_data !== undefined) {
				$(that.get_csv_data).show();
			}
		});
		$(this.plotcontainer.parent()).mouseleave(function() {
			if(that.get_csv_data !== undefined) {
				$(that.get_csv_data).hide();
			}
		});

		$(this.get_csv_data).hide();
		var that = this,
			csv_content = '"component";"resource";"metric";"type"<br>';

		this.get_csv_data.click(function (){
			var serie,
				line_timestamps,
				line_values,
				line_start;

			for (var serieId in that.series) {
				serie = that.series[serieId];
				var node = serie.node,
					position,
					values = ['values'],
					timestamps = ['timestamps'],
					point,
					head;

				for (var position in serie.data) {
					point = serie.data[position];
					timestamps.push(point[0]);
					values.push(point[1])
				}

				head = [node.component, node.resource, node.metric];
				line_start = '"' + head.join('";"') + '";"';
				line_values 	= line_start + values.join('";"') + '"<br>';
				line_timestamps = line_start + timestamps.join('";"') + '"<br>';

				csv_content += line_values + line_timestamps;

			}
			postDataToURL('/echo', [
				{'name': 'filename',	'value': head.join('-') + '.csv'},
				{'name': 'content',		'value': csv_content},
				{'name': 'header',		'value': 'text/csv'}
			]);
		});


	},

	destroyChart: function() {
		this.chart.destroy();
	},

	getSeriesConf: function() {

		var series = [];

		Ext.Object.each(this.series, function(serieId, serie) {
			series.push(serie);
		});

		if(this.groupby_metric) {
			series = series.sort(function (a, b) {
				return a.label.localeCompare( b.label );
			});
		}

		return series;
	},

	getSerieForNode: function(nodeid) {
		var node = this.nodesByID[nodeid];

		var serie = {
			node: node,
			label: node.label,
			data: [],
			last_timestamp: -1,
			xaxis: 1
		};

		return serie;
	},

	addPoint: function(serieId, value, serieIndex) {
		void(serieIndex);
		this.series[serieId].data.push(Ext.clone(value));
	},

	shiftSerie: function(serieId) {
		void(serieId);
		return;
	},

	doRefresh: function(from, to) {
		this.refreshNodes(from, to);
	},

	onRefresh: function(data, from, to) {
		this.callParent(arguments);

		log.debug('Received data:');
		log.dump(data);

		if(data.length > 0) {
			for(var i = 0; i < data.length; i++) {
				var info = data[i];
				var node = this.nodesByID[info.node];
				var serieId = info.node + '.' + node.metrics[0];

				/* create the serie if it doesn't exist */
				if(!(serieId in this.series) || this.series[serieId] === undefined) {
					this.series[serieId] = this.getSerieForNode(info.node);

					//Show component/resource in labels ?
					var serie = this.series[serieId];
					var types = ['resource', 'component'];
					for (type in types) {
						if (serie.node[types[type]] && this[types[type] + '_in_label']) {
							serie.label = serie.node[types[type]] + ' ' + serie.label;
						}
					}
				}

				if(this.reportMode || this.exportMode) {
					this.series[serieId].data = [];
				}

				this.prepareData(serieId);

				/* add data to the serie */
				for(var j = 0; j < info.values.length; j++) {
					var value = info.values[j];

					this.addPoint(serieId, value, i);
				}

				/* shifting serie */
				this.shiftSerie(serieId);
			}
		}

		this.destroyChart();

		this.setChartOptions();

		this.updateAxis(from, to);
		this.updateSeriesConfig();

		this.createChart();
		this.renderChart();
	},

	prepareData: function(serieId) {},
	updateSeriesConfig: function() {},
	dblclick: function() {},

	updateAxis: function(from, to) {
		void(from);
		void(to);
	}

});
