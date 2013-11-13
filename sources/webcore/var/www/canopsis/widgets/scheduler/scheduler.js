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
Ext.define('widgets.scheduler.scheduler', {
	extend: 'canopsis.lib.view.cwidget',
	alias: 'widget.scheduler',

	logAuthor: '[widget][scheduler]',

	treeRk: undefined,

	options: {},
	chart: undefined,

	initComponent: function() {
		this.callParent(arguments);

		if(this.inventory) {
			this.treeRk = this.inventory[0];
		}

		this.rootRk = this.treeRk;

		log.debug('Handling schedule tree: ' + this.treeRk, this.logAuthor);
	},

	afterContainerRender: function() {
		this.callParent(arguments);

		this.setOptions();
		this.createChart();
		this.setChartTitle();
	},

	setOptions: function() {
		var me = this;

		this.options = {
			reportMode: this.reportMode,
			cwidget: function() {
				return me;
			},

			chart: {
				renderTo: this.wcontainerId,
				height: this.getHeight(),
				reflow: false,
				animation: false,
				defaultSeriesType: 'column',
				inverted: true,
				zoomType: 'y'
			},
			title: {
				text: undefined
			},
			legend: {
				enabled: false
			},
			yAxis: {
				type: 'datetime',
				tickmarkPlacement: 'on',
				title: 'Duration'
			},
			xAxis: {
				categories: []
			},
			tooltip: {
				formatter: function() {
					var start  = this.point.low;
					var stop   = this.point.high;
					var length = Math.round(((stop - start) / 1000) / 60);

					var tooltip = '<b>' + this.point.category + '</b><br/>';
					tooltip += 'Start: ' + new Date(start) + '<br/>';
					tooltip += 'Stop: &nbsp;' + new Date(stop) + '<br/>';
					tooltip += 'Duration: ' + length + 'min';

					return tooltip;
				},
				useHTML: true
			},
			series: [{
				name: 'jobs',
				type:'columnrange',
				data: []
			}],
			plotOptions: {
				series: {
					animation: false,
					shadow: false
				},
				columnrange: {
					cursor: 'pointer',
					point: {
						events: {
							click: function() {
								var me = this.series.chart.options.cwidget();

								if(this.node.child_nodes.length > 0) {
									me.treeRk = this.node.rk;
								}
								else {
									me.treeRk = me.rootRk;
								}

								me.setChartTitle();
								me.getNodeInfo(undefined, undefined);
							}
						}
					}
				}
			}
		};
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

	clearGraph: function() {
		log.debug('Cleaning graph', this.logAuthor);

		this.chart.xAxis[0].setCategories([], false);
		this.chart.series[0].setData([], false);
	},

	setChartTitle: function() {
		var me = this;

		var comps = this.treeRk.split('.');
		var title = undefined;

		for(var i = 0; i < comps.length; i++) {
			var current_rk = comps.slice(0, 1 + i).join('.');

			var button = '<span';
			button += ' onclick="Ext.getCmp(\'' + this.id + '\').chartTitleButtonClick(\'' + current_rk + '\');"';
			button += ' onmouseover="this.style.textDecoration=\'underline\';"';
			button += ' onmouseout="this.style.textDecoration=\'none\';"';
			button += '>' + comps[i] + '</span>';

			if(!title) {
				title = button;
			}
			else {
				title += ' / ' + button;
			}
		}

		log.debug('Set title: ' + title, this.logAuthor);
		this.chart.setTitle({text: title, useHTML: true});
	},

	chartTitleButtonClick: function(current_rk) {
		log.debug('Load schedule tree: ' + current_rk, this.logAuthor);

		this.treeRk = current_rk;
		this.setChartTitle();
		this.getNodeInfo(undefined, undefined);
	},

	getUrl: function() {
		return '/rest/events_trees/' + this.treeRk;
	},

	getNodeInfo: function(from, to) {
		if(this.treeRk) {
			Ext.Ajax.request({
				url: this.getUrl(),
				scope: this,
				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);

					if(!data.success) {
						log.error('Impossible to get events tree: ' + this.tree, this.logAuthor);
					}
					else {
						this._onRefresh(data.data, from, to);
					}
				},
				failure: function(result, request) {
					void(result);

					log.error('AJAX request failed (' + request.url + ')', this.logAuthor);
				}
			});
		}
	},

	getEvent: function(rk) {
		var response = Ext.Ajax.request({
			url: this.baseUrl + '/' + rk,
			async: false
		});

		var data = Ext.JSON.decode(response.responseText);

		if(data.success) {
			return data.data[0];
		}

		return null;
	},

	onRefresh: function(tree, from, to) {
		void(from, to);

		this.clearGraph();

		for(var i = 0; i < tree.child_nodes.length; i++) {
			var node = tree.child_nodes[i];
			var evt  = this.getEvent(node.rk);

			log.debug('Adding data for event:', this.logAuthor);
			log.dump(evt);

			var color = global.curvesCtrl.getRenderColors(evt.component, i)[0];

			/* get metrics */
			var startts = undefined;
			var endts   = undefined;

			for(var j = 0; j < evt.perf_data_array.length; j++) {
				var perf = evt.perf_data_array[j];

				if(perf.metric === 'cps_ts_start') {
					startts = perf.value * 1000;
				}
				else if(perf.metric === 'cps_ts_end') {
					endts = perf.value * 1000;
				}
			}

			/* one of the metrics isn't present, so skip this one */
			if(startts === undefined || endts === undefined) {
				continue;
			}

			/* get data */
			var point_name = evt.component.split(this.treeRk + '.')[1];
			var point = {
				node: node,
				name: point_name,
				color: color,
				low: startts,
				high: endts
			};

			log.debug('Add point:', this.logAuthor);
			log.dump(point);

			/* add data on chart */
			this.chart.xAxis[0].categories.push(point_name + ' (' + node.child_nodes.length + ')');

			var serie = this.chart.series[0];
			serie.addPoint(point, false);
		}

		log.debug('Redraw chart.', this.logAuthor);
		this.chart.redraw();
	},

	onResize: function() {
		log.debug('onRezize', this.logAuthor);

		if(this.chart) {
			this.chart.setSize(this.getWidth(), this.getHeight() , false);
		}
	}
});