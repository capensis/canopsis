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

Ext.define('widgets.mini_chart.mini_chart', {
	extend: 'canopsis.lib.view.cperfstoreValueConsumerWidget',

	alias: 'widget.mini_chart',
	requires: [ 'canopsis.lib.view.csparkline' ],
	wcontainer_layout: {
		type: 'vbox',
		align: 'center'
	},

	initComponent: function() {
		this.callParent(arguments);

		this.logAuthor = '[widgets][mini_chart]';

		log.debug('initComponent', this.logAuthor);
		log.debug('nodesByID:', this.logAuthor);

		if(Ext.isArray(this.nodes)) {
			this.nodesByID = parseNodes(this.nodes);
		}
		else {
			this.nodesByID = expandAttributs(this.nodes);
		}

		log.dump(this.nodesByID);

		this.series = {};
		this.charts = {};
	},

	doRefresh: function (from, to) {
		this.from = from;
		this.to = to;

		log.debug('Get values from ' + new Date(from) + ' to ' + new Date(to), this.logAuthor);

		this.refreshNodes(from, to);

		this.callParent(arguments);
	},

	buildParams: function(oFrom, oTo) {
		var post_params = [];

		Ext.Object.each(this.nodesByID, function(id, node) {
			var nodeId = id;
			var from = oFrom;
			var to = oTo;

			if (this.aggregate_interval) {
				var aggregate_interval = this.aggregate_interval * 1000;

				if (this.aggregate_interval < global.commonTs['month']) {
					from = Math.floor(from / aggregate_interval) * aggregate_interval;
				}
				else if(this.aggregate_interval >= global.commonTs['month']) {
					from = moment.unix(from / 1000).startOf('month').unix() * 1000;
				}
				else if(this.aggregate_interval >= global.commonTs['year']) {
					from = moment.unix(from / 1000).startOf('year').unix() * 1000;
				}
			}

			post_params.push({
				id: nodeId,
				metrics: node.metrics,
				from: parseInt(from / 1000),
				to: parseInt(to / 1000)
			});
		}, this);

		return {
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : this.aggregate_method,
			'aggregate_interval': this.aggregate_interval,
			'aggregate_max_points': this.aggregate_max_points,
			'consolidation_method': this.consolidation_method
		};
	},

	makeUrl: function(from, to) {
		return '/perfstore/values' + '/' + parseInt(from / 1000) + '/' + parseInt(to / 1000);
	},

	buildOptions : function(info, values, serie_panel, i) {
		var node = info['node'];

		//Find the print label
		var label;

		if(this.nodesByID[node]['label']) {
			label = this.nodesByID[node]['label'];
		}
		else {
			label = info['metric'];
		}

		//Find the unit
		var unit = '';

		if(this.nodesByID[node]['display_pct']) {
			unit = '%';
		}
		else if(this.nodesByID[node]['u']) {
			unit = this.nodesByID[node]['u'];
		}
		else if(info.bunit) {
			unit = info['bunit'];
		}

		//Find Colors for curve
		var colors = global.curvesCtrl.getRenderColors(label, i);
		var curve_color;

		if(this.nodesByID[node]['curve_color']) {
			curve_color = this.nodesByID[node]['curve_color'];
		}
		else {
			curve_color = colors[0];
		}

		var area_color;

		if(this.nodesByID[node]['area_color']) {
			area_color = this.nodesByID[node]['area_color'];
		}
		else if(Ext.isIE) {
			area_color = curve_color;
		}
		else {
			area_color = this.lightenDarkenColor(curve_color, 50);
		}

		var options = {
			width: serie_panel.getWidth(),
			height: serie_panel.getHeight(),
			chartRangeMinX: values[0][0],
			chartRangeMaxX: values[values.length - 1][0],
			lineColor: curve_color,
			fillColor: area_color,
			barColor: area_color,
			tooltipClassname: 'tooltip-minichart',
			metric: label,
			unit: unit,
			chart_type: this.chart_type,
			original_values: values,

			tooltipFormatter: function(sparkline, options, fields) {
				void(sparkline);

				$(document).find('.tooltip-minichart').css('border', '2px solid ' + curve_color);

				if(options.userOptions.chart_type === 'line_graph') {
					return '<b>' + rdr_tstodate(Math.round(fields['x']/1000)) + '</b><br>' + options.userOptions.metric + ': ' + fields['y'] + ' ' + options.userOptions.unit;
				}

				return '<b>' + rdr_tstodate(Math.round(options.userOptions.original_values[fields[0].offset][0] / 1000)) + '</b><br/>' + options.userOptions.metric + ' : ' + fields[0].value + ' ' + options.userOptions.unit;
			}
		};

		return options;
	},

	onResize: function() {
		Ext.Object.each(this.series, function(id, serie) {
			void(id);

			serie.setWidth(this.getWidth());
			serie.setHeight(this.getHeight() / Ext.Object.getSize(this.nodesByID));
		}, this);
	},

	parseValues: function(serie, values) {
		//MAKE A BETTER LOOP, JUST FOR TEST
		for (var i = 0; i < values.length; i++) {
			values[i][0] = values[i][0] * 1000;

			if(this.nodesByID[serie.node]['display_pct'] && serie.max > 0) {
				values[i][1] = getPct(values[i][1], serie.max);
			}

			if(serie.invert) {
				values[i][1] = -values[i][1];
			}

			values[i][0] = values[i][0] / 1000;
		}

		return values;
	},

	onRefresh: function(data) {
		if(Ext.Object.getSize(this.charts) > 0 || data.length > 0 ) {
			for (var i = 0; i < data.length; i++) {
				var info = data[i];

				var node = info['node'] ;
				var values = this.parseValues( info, info['values'] ) ;

				if(Ext.ComponentQuery.query('#'+this.series[node].getId() +" > csparkline").length  === 0 ) {
					//Find the print label
					var label;

					if(this.nodesByID[node]['label']) {
						label = this.nodesByID[node]['label'];
					}
					else {
						label = info['metric'];
					}

					info['metric'] = label;

					//Add a component with the print label
					this.series[node].add( {
						xtype:"panel",
						flex: 2,
						bodyCls: "valigncenter",
						html: "<div><b>" + label + "</b></div>",
						border: false
					});

					//We add the serie panel
					this.series[node].add({
						xtype: "csparkline",
						values: Ext.clone(values),
						node: Ext.clone(this.nodesByID[node]),
						info: Ext.clone(info),
						flex: 4,
						chart_type: this.chart_type,
						border: false
					});

					this.charts[node] = true;

					if(this.nodesByID[node]['printed_value']) {
						this.series[node].add( {
							xtype: "panel",
							flex: 2,
							bodyCls: "valigncenter padding-left",
							border: false,
							html: "<div><b>" + values[values.length - 1][1] + "</b></div>"
						});
					}
				}
				//We display the last value or the evolution
				else {
					//we redraw only the graph and update the last value
					Ext.ComponentQuery.query('#' + this.series[node].getId() + " > csparkline")[0].addValues(values);

					if(Ext.ComponentQuery.query('#' + this.series[node].getId() + " > panel").length === 3) {
						Ext.ComponentQuery.query('#' + this.series[node].getId() + " > panel")[2].update("<div><b>" + values[values.length - 1][1] + "</b></div>");
					}
				}
			}
		}
		else {
			this.getEl().mask(_('No data on interval'));
		}
	},

	afterContainerRender: function() {
		var me = this;

		Ext.Object.each(this.nodesByID, function(id) {
			var pTop = 1;
			var pBottom = 5;
			var gHeight = me.getSize().height - pTop - pBottom - 5;

			if(me.header !== undefined) {
				gHeight = gHeight - me.header.getSize().height;
			}

			var serie = {
				layout: {
					type: "hbox",
					align: "stretch"
				},
				style: {
					paddingLeft: '10px',
					paddingRight: '5px',
					paddingTop: pTop + 'px',
					paddingBottom: pBottom + 'px'
				},
				width: me.getSize().width,
				height: gHeight / Ext.Object.getSize(me.nodesByID),
				border: false
			};

			me.series[id] = me.wcontainer.add(serie) ;
		});

		this.callParent();
	},

	beforeDestroy: function() {
		Ext.Object.each(this.series, function(node_id, serie) {
			void(node_id);

			serie.removeAll();
			serie.destroy();
		});

		this.wcontainer.removeAll();
		this.callParent();
	}
});
