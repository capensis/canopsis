//need:app/lib/view/cwidgetGraph.js
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

Ext.define('widgets.category_graph.category_graph', {
	extend: 'canopsis.lib.view.cwidgetGraph',
	alias: 'widget.category_graph',

	logAuthor: '[widgets][category_graph]',

	timeNav: false,
	timeNav_window: global.commonTs.week,

	diagram_type: 'pie', // 'or column'
	aggregate_max_points: 1,
	aggregate_method: 'LAST',
	aggregate_round_time: true,

	//Default Options
	max: 0,

	// Layout options
	other_label: 'Free',

	// pie specific options
	pie_size: 0.8,
	startAngle: 1.5,
	innerRadius: 0,

	// Bar specific options
	stacked_graph: false,
	verticalDisplay: true,

	backgroundColor: '#FFFFFF',
	borderColor: '#FFFFFF',
	borderWidth: 0,
	title_fontSize: 15,

	// legend
	legend: true,
	legend_verticalAlign: 'bottom', // TODO: this property is not managed by flotchart
	legend_align: 'center', // TODO: this property is not managed by flotchart
	legend_layout: 'horizontal', // TODO: this property is not managed by flotchart
	legend_backgroundColor: null, // TODO: this property is not managed by flotchart
	legend_borderColor: '#909090', // TODO: this property is not managed by flotchart
	legend_borderWidth: 1, // TODO: this property is not managed by flotchart
	legend_fontSize: 12, // TODO: this property is not managed by flotchart
	legend_fontColor: '#3E576F', // TODO: this property is not managed by flotchart

	// label
	labels: true,
	nameInLabelFormatter: false,
	pctInLabel: true,
	labels_size: "x-small",

	gradientColor: false, // TODO: this property is not managed by flotchart

	tooltip: true,

	initComponent: function() {
		this.callParent(arguments);

		// list nodes by categories
		this.categories = [null];
		this.nodesByCategories = {};
		this.nodesNoCategory = [];

		for(var nodeId in this.nodesByID) {
			var node = this.nodesByID[nodeId];

			// classify the node if category is set
			if(node.category) {
				// try to find the category index
				var catIdx = this.categories.indexOf(node.category);

				if(catIdx === -1) {
					catIdx = this.categories.push(node.category) - 1;
				}

				node.categoryIndex = catIdx;

				// create the category if needed
				if(!this.nodesByCategories[node.category]) {
					this.nodesByCategories[node.category] = [];
				}

				this.nodesByCategories[node.category].push(node);
			}
			else {
				node.categoryIndex = 0;
				this.nodesNoCategory.push(node);
			}
		}

		if(this.categories.length > 1) {
			this.stacked_graph = false;
		}
	},

	setChartOptions: function() {
		this.callParent(arguments);

		var me = this;

		$.extend(this.options,
			{
				series: {
					lines: {
						show: false,
					},
					pie: {
						show: (this.diagram_type === 'pie'),
						innerRadius: this.innerRadius,
						label: {
							show: this.labels,
							size: this.labels_size,
							formatter: function(label, slice) {
								var outer = $('<div/>');
								var inner = $('<div/>');

								inner.css({
									'font-size': 'x-small',
									'text-align': 'center',
									'padding': '2px',
									'color': slice.color
								});

								outer.append(inner);

								// generate result
								var result = '';

								if(me.nameInLabelFormatter) {
									result += '<b>' + label + '</b><br/>';
								}

								if(me.pctInLabel) {
									result += slice.percent.toFixed(1) + '%';
								}
								else if(me.humanReadable) {
									result += rdr_humanreadable_value(slice.data[0][1], slice.node.bunit);
								}
								else {
									result += slice.data[0][1];
								}

								// build HTML
								inner.html(result);
								return outer.html();
							}
						},
						tilt: this.tilt,
						stroke:{
							color: "fff",
							width: this.stroke_width
						},
						startAngle: this.startAngle,
						radius: this.pie_size / 100
					},
					bars: {
						show: (this.diagram_type === 'column'),
						horizontal: !this.verticalDisplay,
						barWidth: 1,
						zero: true,
						dataLabels: this.labels,
						showNumbers: (this.labels && this.diagram_type === 'column')
					},
					stack: this.stacked_graph
				},
				grid: {
					hoverable: true,
					clickable: true,
				},
				legend: {
					hideable: true,
					show: this.legend
				},
				xaxis: {
					show: (this.diagram_type === 'column' && !this.verticalDisplay)
				},
				yaxis: {
					show: (this.diagram_type === 'column' && this.verticalDisplay),
					tickFormatter: function(val, axis) {
						if(me.humanReadable) {
							return rdr_humanreadable_value(val, axis.options.unit);
						}
						else {
							return val + ' ' + axis.options.unit;
						}
					}
				},
				tooltip: this.tooltip,
				tooltipOpts: {
					content: function(label, xval, yval, item) {
						var val = item.series.data[item.dataIndex][1];

						if(me.categories.length > 1) {
							for(var i = 0; i < me.categories.length; i++) {
								var category = me.categories[i];
								var serie = me.getSerieForCategory(category);
								var nodes = (category ? me.nodesByCategories[category] : me.nodesNoCategory);

								if(serie.label === label) {
									var output = '<p><b>Category: ' + label + '</b></p><p>';

									// add total
									output += '<b>Total:</b> ';

									if(me.humanReadable) {
										val = rdr_humanreadable_value(val, serie.node.bunit);
									}

									output += val;

									// add each categorized series value
									for(var j = 0; j < nodes.length; j++) {
										var node = nodes[j];
										var serieId = node.id + '.' + node.metrics[0];
										var serie = me.series[serieId];

										output += '<br/><b>' + serie.label + ':</b> ';

										var subval = serie.data[0][1];

										if(me.humanReadable) {
											subval = rdr_humanreadable_value(subval, serie.node.bunit);
										}

										output += subval;
									}

									// finalize output
									output += '</p>';

									return output;
								}
							}
						}
						else {
							for(var serieId in me.series) {
								var serie = me.series[serieId];

								if(serie.label === label) {
									if(me.humanReadable) {
										val = rdr_humanreadable_value(val, serie.node.bunit);
									}

									return '<b>' + label + ':</b> ' + val;
								}
							}
						}
					}
				}
			}
		);
	},

	addPoint: function(serieId, value, serieIndex) {
		var point = [(this.stacked_graph ? 0 : serieIndex), value[1]];

		this.series[serieId].data = [point];
	},

	getSeriesConf: function() {
		// check if categories are set
		// By default, it contains the null object, so empty is 1
		if(this.categories.length === 1) {
			return this.callParent(arguments);
		}

		var series = [];

		for(var i = 0; i < this.categories.length; i++) {
			var category = this.categories[i];
			var serie = this.getSerieForCategory(category);

			series.push(serie);
		}

		return series;
	},

	getSerieForCategory: function(category) {
		// fetch nodes
		var nodes;

		if(category === null) {
			nodes = this.nodesNoCategory;
			category = 'No Category';
		}
		else {
			nodes = this.nodesByCategories[category];
		}

		// generate series
		var cat_serie = {
			node: {
				id: 'category',
				metrics: [category]
			},
			label: category,
			data: [],
			xaxis: 1,
			yaxis: 1
		};

		// consolidates points in the same category
		var val = 0;
		var idx = 0;

		for(var i = 0; i < nodes.length; i++) {
			var node = nodes[i];
			var serieId = node.id + '.' + node.metrics[0];
			var serie = this.series[serieId];

			if(!serie) {
				serie = this.getSerieForNode(node.id);
			}

			if(serie.data.length > 0) {
				val += serie.data[0][1];
			}

			if(node.categoryIndex) {
				idx = node.categoryIndex;
			}

			cat_serie.node.bunit = node.bunit;
		}

		cat_serie.data = [[idx, val]];

		// add new serie to series
		var catSerieId = 'category.' + category;
		this.series[catSerieId] = cat_serie;

		return cat_serie;
	},

	createChart: function() {
		/* Computes the difference between the series sum and max user input value and then add it to series if max value > series sum */
		//clean series first
		delete this.series['internal.max_diff'];

		if(this.max > 0) {
			var total = 0;
			var serie_count = 0;
			var node_id = undefined;
			var serie = undefined;

			//gets series sum
			for(node_id in this.series) {
				serie_count++;

				if((this.series[node_id].show === undefined || this.series[node_id].show) && this.series[node_id].data.length === 1) {
					total += this.series[node_id].data[0][1];
				}
			}

			if(total < this.max) {
				//data is a list of point with a single computed point (barchart)
				var label = this.other_label;
				var stacked = this.stacked_graph;
				var max = this.max;

				this.series['internal.max_diff'] = {
					label: label,
					data: [[stacked ? 0: serie_count, max - total]],
					node: {
						id: 'internal',
						metrics: ['max_diff']
					},
				};
			}
		}

		this.callParent(arguments);
	}
});
