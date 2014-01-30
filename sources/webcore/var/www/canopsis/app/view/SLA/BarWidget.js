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

Ext.define('canopsis.view.SLA.BarWidget', {
	extend: 'Ext.container.Container',
	alias: 'widget.slabarwidget',

	logAuthor: '[SLA][Bar]',

	crit: undefined,
	delay: undefined,

	from: function() {
		return undefined;
	},

	to: function() {
		return undefined;
	},

	initComponent: function() {
		this.callParent(arguments);

		this.nodes = {};

		this.refreshNodes();
	},

	refreshNodes: function() {
		Ext.Ajax.request({
			url: '/perfstore/get_all_metrics',
			method: 'GET',
			scope: this,

			params: {
				show_internals: true,
				limit: 0
			},

			success: function(response) {
				var data = Ext.JSON.decode(response.responseText);

				var nodes = [];

				for(var i = 0; i < data.data.length; i++) {
					var metric = data.data[i];

					if(metric.me === 'cps_state') {
						/* get component associated to the metric */
						var response = null;

						if(!('re' in metric)) {
							response = Ext.Ajax.request({
								url: '/entities/',
								method: 'POST',
								async: false,

								jsonData: {
									filter: {
										'type': 'component',
										'name': metric.co
									}
								}
							});
						}
						else {
							response = Ext.Ajax.request({
								url: '/entities/',
								method: 'POST',
								async: false,
								
								jsonData: {
									filter: {
										'type': 'resource',
										'name': metric.re,
										'component': metric.co
									}
								}
							});
						}

						response = Ext.JSON.decode(response.responseText);

						if(response.data.length >= 1) {
							/* add node only if the criticity match */
							//if(response.data[0].mCrit === this.crit || response.data[0].mWarn === this.crit) {
								this.nodes[metric._id] = response.data[0];

								nodes.push({
									id: metric._id,
									from: this.from(),
									to: this.to()
								});
							//}
						}
					}
				}

				Ext.Ajax.request({
					url: '/perfstore/values',
					scope: this,
					method: 'POST',

					params: {
						'nodes': Ext.JSON.encode(nodes),
						'timezone': new Date().getTimezoneOffset() * 60
					},

					success: function(response) {
						var data = Ext.JSON.decode(response.responseText);

						this.onRefresh(data.data);
					}
				});
			}
		});
	},

	onRefresh: function(data) {
		var sla = {
			warn: {
				ok: 0,
				nok: 0,
				out: 0
			},
			crit: {
				ok: 0,
				nok: 0,
				out: 0
			}
		};

		for(var i = 0; i < data.length; i++) {
			var metric = data[i];
			var node = this.nodes[metric.node];

			var warning_periods = [];
			var critical_periods = [];

			var last_perf = [-1, -1];
			var inperiod = false;

			/* get warning/critical periods */
			for(var j = 0; j < metric.values.length; j++) {
				var perfdata = metric.values[j];
				var state = parseInt(this.getState(perfdata[1])[0]);

				/* a problem occured, save the perfdata */
				if(state !== 0) {
					last_perf = [perfdata[0], state];
					inperiod = true;
				}
				/* went back to normal, calculate period */
				else {
					inperiod = false;

					/* it's a warning period ? */
					if(last_perf[1] === 1 && node.mWarn === this.crit) {
						warning_periods.push({
							from: last_perf[0],
							to: perfdata[0],
							duration: perfdata[0] - last_perf[0]
						});
					}
					/* or a critical period ? */
					else if(last_perf[1] >= 2 && node.mCrit === this.crit) {
						critical_periods.push({
							from: last_perf[0],
							to: perfdata[0],
							duration: perfdata[0] - last_perf[0]
						});
					}

					/* NB: unknown status are critical periods too */
				}
			}

			console.log(warning_periods);
			console.log(critical_periods);

			/* get acknowledgements */
			Ext.Ajax.request({
				url: '/entities/',
				scope: this,
				method: 'POST',

				jsonData: {
					filter: {
						'type': 'ack',
						'component': node.component,
						'resource': node.resource
					}
				},

				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);

					for(var k = 0; k < data.data.length; k++) {
						var ack = data.data[k];

						/* check SLA for warning periods */
						for(var l = 0; l < warning_periods.length; l++) {
							var period = warning_periods[l];

							if(!(period.from <= ack.timestamp && ack.timestamp <= period.to)) {
								sla.warn.out++;
							}
							else if(period.duration >= this.delay) {
								sla.warn.nok++;
							}
							else {
								sla.warn.ok++;
							}
						}

						/* check SLA for critical periods */
						for(var l = 0; l < critical_periods.length; l++) {
							var period = critical_periods[l];

							if(!(period.from <= ack.timestamp && ack.timestamp <= period.to)) {
								sla.crit.out++;
							}
							else if(period.duration >= this.delay) {
								sla.crit.nok++;
							}
							else {
								sla.crit.ok++;
							}
						}
					}

					console.log(sla);
				}
			});
		}
	},

	getState: function(cps_state) {
		var state = '' + cps_state;

		while(state.length < 3) {
			state = '0' + state;
		}

		return state;
	}
});