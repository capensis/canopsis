//need:app/lib/view/cwidget.js
/*
# Copyright (c) 2014 "Capensis" [http://www.capensis.com]
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

Ext.define('widgets.text.text', {
	extend: 'canopsis.lib.view.cwidget',
	alias: 'widget.text',

	logAuthor: '[widgets][text]',

	initComponent: function() {
		this.callParent(arguments);

		this.wanted_metrics = {};
		this.variables = [];

		this.extractVariables(this.text);
	},

	extractVariables: function(text) {
		/* extract variables from template */
		var vars = text.match(/({.*?})/g);

		if(!vars) {
			return;
		}

		for(var i = 0; i < vars.length; i++) {
			var variable = vars[i].slice(1, -1);
			var parts = variable.split(':');

			/* detect wanted components/resources/metrics from variables */
			if(parts[0] === 'perfdata' || parts[0] === 'perf_data') {
				this.wanted_metrics[variable] = {
					metric: parts[1],
					info: parts[2],
					component: parts[3],
					resource: parts[4]
				};
			}

			if(this.variables.indexOf(variable) === -1) {
				this.variables.push(variable);
			}
		}
	},

	getPerfdataUrl: function(from, to) {
		var url = '/perfstore/values';

		if(from) {
			url = url + '/' + parseInt(from / 1000);

			if(!to) {
				to = Ext.Date.now();
			}

			url = url + '/' + parseInt(to / 1000);
		}

		return url;
	},

	getNodeInfo: function(from, to, advancedFilters) {
		if(this.nodeId && this.nodeId.length > 0) {
			Ext.Ajax.request({
				url: '/rest/events/event',
				method: 'GET',
				params: {
					ids: Ext.JSON.encode(this.nodeId)
				},
				scope: this,

				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);

					this._onRefresh(data.data, from, to, advancedFilters);
				}
			});
		}
		else {
			this._onRefresh([], from, to, advancedFilters);
		}
	},

	onRefresh: function(data, from, to, advancedFilters) {
		var template_data = {};

		for(var i = 0; i < data.length; i++) {
			var evt = Ext.create('canopsis.model.Event', data[i]);
			var co = evt.get('component');
			var re = evt.get('resource');

			var fields = Ext.ModelManager.getModel('canopsis.model.Event').getFields();

			for(var j = 0; j < fields.length; j++) {
				var field = fields[j].name;

				if(field !== 'perf_data' && field !== 'perf_data_array') {
					var key = field + ':' + co;

					if(re) {
						key = key + ':' + re;
					}

					template_data[key] = evt.get(field);
					template_data[field] = evt.get(field);
				}
			}
		}

		var perfRequest = {
			'$or': []
		};

		for(var variable in this.wanted_metrics) {
			var metric = this.wanted_metrics[variable];
			var filter = {
				'me': metric.metric
			};

			if(metric.component) {
				filter.co = metric.component;
			}

			if(metric.resource) {
				filter.re = metric.resource;
			}

			perfRequest['$or'].push(filter);
		}

		if(perfRequest['$or'].length > 0) {
			if(perfRequest['$or'].length === 1) {
				perfRequest = perfRequest['$or'][0];
			}

			Ext.Ajax.request({
				url: '/perfstore/get_all_metrics',
				method: 'GET',
				params: {
					'filter': Ext.JSON.encode(perfRequest),
					'limit': 0,
					'show_internals': true
				},
				scope: this,

				success: function(response) {
					var r_nodes = Ext.JSON.decode(response.responseText);

					var nodes = [];
					var nodesByID = {};

					for(var i = 0; i < r_nodes.data.length; i++) {
						var node = r_nodes.data[i];

						nodes.push({
							id: node._id
						});

						nodesByID[node._id] = node;
					}

					Ext.Ajax.request({
						url: this.getPerfdataUrl(from, to),
						method: 'POST',
						params: {
							'nodes': Ext.JSON.encode(nodes),
							'aggregate_method': 'LAST',
							'aggregate_max_points': 1,
							'timezone': new Date().getTimezoneOffset() * 60
						},
						scope: this,

						success: function(response) {
							var r_vals = Ext.JSON.decode(response.responseText);

							function genKey(prefix, metric, data, component, resource) {
								var key = prefix + ':' + metric + ':' + data;

								if(component) {
									key = key + ':' + component;

									if(resource) {
										key = key + ':' + resource;
									}
								}

								return key;
							}

							for(var i = 0; i < r_vals.data.length; i++) {
								var metric = r_vals.data[i];
								var node = nodesByID[metric.node];

								var datas = {
									'value': metric.values[0][1],
									'unit': metric.bunit,
									'component': node.co,
									'resource': node.re
								};

								for(var k in datas) {
									var key = genKey('perfdata', metric.metric, k);
									template_data[key] = datas[k];

									key = genKey('perfdata', metric.metric, k, node.co, node.re);
									template_data[key] = datas[k];

									key = genKey('perf_data', metric.metric, k);
									template_data[key] = datas[k];

									key = genKey('perf_data', metric.metric, k, node.co, node.re);
									template_data[key] = datas[k];
								}
							}

							this.fillData(template_data, from, to);
							this.computeMathOperations();
						}
					});
				}
			});
		}
		else {
			this.fillData(template_data, from, to);
			this.computeMathOperations();
		}
	},

	fillData: function(data, from, to) {
		console.log(data);
		var text = this.text;

		for(var i = 0; i < this.variables.length; i++) {
			var variable = this.variables[i];

			if(data[variable] !== undefined) {
				text = replaceAll('{' + variable + '}', data[variable], text);
			}
			else {
				text = replaceAll('{' + variable + '}', 'undefined', text);
			}
		}

		try {
			if(from) {
				data.from = rdr_tstodate(parseInt(from / 1000));
			}

			if(to) {
				data.to = rdr_tstodate(parseInt(to / 1000));
			}

			var template = new Ext.XTemplate('<div>' + text + '</div>');

			this.setHtml(template.apply(data));
		}
		catch(err) {
			this.setHtml(err);
		}
	},

	computeMathOperations: function() {
		var math = mathjs();

		$('#' + this.wcontainerId + ' .mathexpression').each(function() {
			$(this).html(math.eval($(this).html()));
		});
	}
});
