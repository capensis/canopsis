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
	useLastRefresh: false,
	aggregate_method: 'LAST',

	initComponent: function() {
		this.callParent(arguments);

		this.wanted_metrics = {};
		this.templateVariables = [];

		this.extractTemplateVariables(this.text);
	},

	extractTemplateVariables: function(text) {
		/* extract variables from template */
		var extractedTemplateVars = text.match(/({.*?})/g);

		if(!extractedTemplateVars) {
			return;
		}

		for(var i = 0; i < extractedTemplateVars.length; i++) {
			var templateVariable = extractedTemplateVars[i].slice(1, -1);
			var parts = templateVariable.split(':');

			/* detect wanted components/resources/metrics from variables */
			if(parts[0] === 'perfdata' || parts[0] === 'perf_data') {
				this.wanted_metrics[templateVariable] = {
					metric: parts[1],
					info: parts[2],
					component: parts[3],
					resource: parts[4]
				};
			}

			if(this.templateVariables.indexOf(templateVariable) === -1) {
				this.templateVariables.push(templateVariable);
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

	onRefresh: function(data, from, to) {
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

		for(var templateVariable in this.wanted_metrics) {
			var metric = this.wanted_metrics[templateVariable];
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

					var requestParams = {
							'nodes': Ext.JSON.encode(nodes),
							'aggregate_method': this.aggregate_method,
							'aggregate_max_points': 1,
							'timezone': new Date().getTimezoneOffset() * 60,
						};

					if (this.subset_selection) {

						log.debug('Adding live reporting advanced filter to post param');
						requestParams['subset_selection'] = Ext.JSON.encode(this.subset_selection);
						//remove subset selection to avoid further side effects
						this.subset_selection = undefined;
					}

					Ext.Ajax.request({
						url: this.getPerfdataUrl(from, to),
						method: 'POST',
						params: requestParams,
						scope: this,

						success: function(response) {
							log.debug("+++ request success");
							log.dump(response);
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

								var value = metric.values[0][1];

								if(this.humanReadable) {
									value = rdr_humanreadable_value(value, metric.bunit);
								}

								var datas = {
									'value': value,
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

		this.add_csv_download_button();
	},

	fillData: function(data, from, to) {

		var text = this.text;

		for(var i = 0; i < this.templateVariables.length; i++) {
			var templateVariable = this.templateVariables[i];

			if(data[templateVariable] !== undefined) {
				text = replaceAll('{' + templateVariable + '}', data[templateVariable], text);
			}
			else {
				text = replaceAll('{' + templateVariable + '}', 'undefined', text);
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
			this.tmpldata = data;

			this.setHtml(template.apply(data));
		}
		catch(err) {
			this.setHtml(err);
		}
	},

	computeMathOperations: function() {
		var math = mathjs();
		var that = this;
		var expressionCounter = 1;

		$('#' + this.wcontainerId + ' .mathexpression').each(function() {
			try {
				var expression = $(this).html();

				expression = expression.replace(/\ /g,'');
				expression = expression.replace(/undefined/g,'0');
				expression = math.eval(expression);
				if(typeof expression === 'object' || isNaN(expression)){
					expression = 0;
				}
				if (that.round_math_expression !== undefined) {
					expression = parseFloat(expression).toFixed(that.round_math_expression);
				}

				$(this).html(expression);
			} catch(err) {
				log.warning('unable to compute math expression' + $(this).html());
				console.error(err.stack);
			}
		});
	},

	add_csv_download_button: function() {
		log.debug("@@@ add_csv_download_button");
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

		$('#' + this.wcontainerId).append(this.get_csv_data);

		var that = this;

		$('#' + this.wcontainerId).mouseenter(function() {
			if(that.get_csv_data !== undefined) {
				$(that.get_csv_data).show();
			}
		});
		$('#' + this.wcontainerId).mouseleave(function() {
			if(that.get_csv_data !== undefined) {
				$(that.get_csv_data).hide();
			}
		});

		$(this.get_csv_data).hide();
		var that = this;

		this.get_csv_data.click(function (){
			if($('#' + that.wcontainerId + " table.exportable").size() === 0) {
				global.notify.notify(_('Issue'), _("No table with class \"exportable\" found in the text cell."), 'info');
			}
			else
			{
				var csv = $('#' + that.wcontainerId + " table.exportable").table2CSV({delivery:'value'});
				window.location.href = 'data:text/csv;charset=UTF-8,' + encodeURIComponent(csv);
			}
		});
	}

});
