//need:app/lib/view/cwidget.js
/*
# Copyright (c) 2013 "Capensis" [http://www.capensis.com]
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
Ext.define('canopsis.lib.view.cperfstoreValueConsumerWidget', {
	extend: 'canopsis.lib.view.cwidget',

	getUrl: function(from, to) {
		var url = '/perfstore/values' + (from !== undefined ? ('/' + parseInt(from / 1000) + '/' + parseInt(to / 1000)) : '');

		return url;
	},

	refreshNodes: function(from, to, advancedFilters) {
		console.log("@@@@@ cpsvc refreshNodes");
		console.log(advancedFilters);
		if(this.nodesByID && Ext.Object.getSize(this.nodesByID) != 0) {
			var url = this.getUrl(from, to);

			Ext.Ajax.request({
				url: url,
				scope: this,
				params: this.getParams(from, to),
				method: 'POST',

				success: function(response) {
					var data = Ext.JSON.decode(response.responseText);
					data = data.data;

					if(data.length > 0) {
						if(this.nodesByID[data[0]['node']]['order'] !== undefined) {
							var that = this;

							data.sort(function(a,b) {
								return that.nodesByID[a['node']]['order']-that.nodesByID[b['node']]['order'];
							});

							console.log("ONREFRESH     advancedFilters");
							console.log(data);
							//TODO bad check (exclusions should actually not be there)
							if(advancedFilters !== undefined)
							{
								for (var i = 0; i < data.length; i++) {
									var node = data[i];

									console.log(" ++++ node");
									console.log(node);

									for (var j = node.values.length - 1; j >= 0; j--) {
										var value = node.values[j];
										for (var k = advancedFilters.exclusions.length - 1; k >= 0; k--) {
											var exclusion = advancedFilters.exclusions[k];
											if(value[0] > exclusion.from && value[0] < exclusion.to)
											{
												//remove the point from array
												// data[i].values.splice(j, 1);
											}
										}
									}

									for (var i = advancedFilters.keyFilter.length - 1; i >= 0; i--) {
										var filter = advancedFilters.keyFilter[i];

										console.log("comparing keys");
										console.log(filter);
										console.log(node[filter.key]);
										console.log(filter.value);
										if(node[filter.key] != filter.value) {
											console.log("delete key ");
											console.log(data[i]);
											data.splice(i, 1);
										}
									}
								}
							}
						}
					}

					this.onRefresh(data, from, to);
				},

				failure: function(result, request) {
					void(result);

					log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
				}
			});
		}
		else {
			log.debug('No nodes specified', this.logAuthor);

			this.getChart().showLoading(_('Please choose a valid metric in wizard'));
		}
	},

	getChart: function() {
		if(this.chart === undefined) {
			throw new Exception("chart field is not defined in " + this);
		}

		return this.chart;
	},

	getParams: function(from, to) {
		var post_params = [];

		Ext.Object.each(this.nodesByID, function(id, node) {
			var nodeId = id;
			var serieId = nodeId + '.' + node.metrics[0];
			var serie = this.series !== undefined ? this.series[serieId] : undefined;

			if(from) {
				if(!this.reportMode && !this.timeNav) {
					if(serie && serie['last_timestamp']) {
						from = serie['last_timestamp'];
					}

					if(from < (to - (this.time_window * 1000))) {
						from = to - (this.time_window * 1000);
					}
				}

				if(this.aggregate_interval) {
					var aggregate_interval = this.aggregate_interval * 1000;

					if(this.aggregate_interval < global.commonTs['month']) {
						from = Math.floor(from / aggregate_interval) * aggregate_interval;
					}
					else if(this.aggregate_interval >= global.commonTs['year']) {
						from = moment.unix(from / 1000).startOf('year').unix() * 1000;
					}
					else if(this.aggregate_interval >= global.commonTs['month']) {
						from = moment.unix(from / 1000).startOf('month').unix() * 1000;
					}

					var tzOffset = new Date().getTimezoneOffset();
					log.debug('TZ Offset: ' + tzOffset, this.logAuthor);
					from += tzOffset * 60 * 1000;
				}

				log.debug('Serie ' + nodeId + ' ' + node.metrics + ':', this.logAuthor);
				log.debug(' + From: ' + new Date(from) + ' (' + from + ')', this.logAuthor);
				log.debug(' + To:   ' + new Date(to) + ' (' + to + ')', this.logAuthor);

			}

			post_param = {
				id: nodeId,
				metrics: node.metrics
			}

			if(from) {
				post_param['from'] = parseInt(from / 1000);
			}
			if(to) {
				post_param['to'] = parseInt(to / 1000);
			}

			this.processPostParam(post_param);

			post_params.push(post_param);
		}, this);

		post_params = {
			'nodes': Ext.JSON.encode(post_params),
			'timezone': new Date().getTimezoneOffset() * 60
		};

		if(this.aggregate_method) {
			post_params['aggregate_method'] = this.aggregate_method;
		}

		if(this.aggregate_interval) {
			post_params['aggregate_interval'] = this.aggregate_interval;
		}

		if(this.aggregate_max_points) {
			post_params['aggregate_max_points'] = this.aggregate_max_points;
		}

		if(this.aggregate_round_time) {
			post_params['aggregate_round_time'] = this.aggregate_round_time;
		}

		if(this.consolidation_method) {
			post_params['consolidation_method'] = this.consolidation_method;
		}

		this.processPostParams(post_params);

		return post_params;
	},

	processPostParam: function(post_param) {
		void(post_param);

		return;
	},

	processPostParams: function(post_params) {
		void(post_params);

		return;
	}
});
