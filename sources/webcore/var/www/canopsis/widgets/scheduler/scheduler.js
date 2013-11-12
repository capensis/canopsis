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
	legend: true,

	steps: [],
	checks: [],

	initComponent: function() {
		this.callParent(arguments);

		if(this.inventory) {
			this.treeRk = this.inventory[0];
		}

		log.debug('Tree RK:', this.logAuthor);
		log.dump(this.treeRk);
	},

	afterContainerRender: function() {
		this.callParent(arguments);

		this.setOptions();
		this.createChart();
	},

	setOptions: function() {
		this.options = {
			reportMode: this.reportMode,
			cwdiget: this,

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
				text: this.treeRk
			},
			legend: {
				enabled: this.legend
			},
			yAxis: {
				type: 'datetime',
				tickmarkPlacement: 'on',
				title: 'Duration'
			},
			xAxis: {
				categories: ['Test']
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
				type:'column',
				data: [{
					color: '#FF0000',
					low: 1364374000000,
					high: 1364374300000
				},{
					color: '#00FF00',
					low: 1364375000000,
					high: 1364375300000
				}]
			}],
			plotOptions: {
				series: {
					animation: false,
					shadow: false
				},
				column: {}
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
		this.event = this.getEvent(tree.rk);

		console.log(this.event);
	}
});