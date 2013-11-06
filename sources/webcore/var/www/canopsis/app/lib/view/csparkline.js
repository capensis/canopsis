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
Ext.define('canopsis.lib.view.csparkline' , {
	extend: 'Ext.panel.Panel',
	layout: "fit",
	alias: 'widget.csparkline',

	listeners: {
		afterlayout: function() {
			this.buildSparkline();
		},
		resize: function() {
			this.buildSparkline();
		}
	},

	initComponent: function() {
		this.callParent(arguments);
	},

	buildOptions: function() {
		//Find the print label
		var label;

		if(this.node['label']) {
			label = this.node['label'];
		}
		else {
			label = this.info['metric'];
		}

		//Find the unit
		var unit = '';

		if(this.node['display_pct']) {
			unit = '%';
		}
		else if(this.node['u']) {
			unit = this.node['u'];
		}
		else if(this.info.bunit) {
			unit = this.info['bunit'];
		}

		//Find Colors for curve
		var colors = global.curvesCtrl.getRenderColors(label, 0);
		var curve_color;

		if(this.node['curve_color']) {
			curve_color = this.node['curve_color'];
		}
		else {
			curve_color = colors[0];
		}

		var area_color;

		if(this.node['area_color']) {
			area_color = this.node['area_color'] ;
		}
		else if(Ext.isIE) {
			area_color = curve_color;
		}
		else {
			area_color = chroma.hex(curve_color).brighten(20).hex();
		}

		if(!this.chart_type) {
			this.chart_type = 'line_graph';
		}

		var options = {
			width: this.getWidth(),
			height: this.getHeight(),
			chartRangeMinX: this.values[0][0],
			chartRangeMaxX: this.values[this.values.length - 1][0],
			lineColor: curve_color,
			fillColor: area_color,
			barColor: area_color,
			tooltipClassname: 'tooltip-sparkline',
			metric: label,
			unit: unit,
			chart_type: this.chart_type,
			original_values: Ext.clone(this.values),
			tooltipFormatter: this.tooltipFormatter
		};

		this.options = options;
	},

	addValues: function(values) {
		this.values = this.values.slice(values.length);

		for(var i = 0; i < values.length; i++) {
			this.values.push(values[i]);
		}

		this.buildSparkline();
	},

	tooltipFormatter: function(sparkline, options, fields) {
		void(sparkline);

		$('.tooltip-sparkline').css('border-color', options.userOptions.lineColor);

		var html;

		if(options.userOptions.chart_type === 'line_graph') {
			html = '<b>' + rdr_tstodate(Math.round(fields['x'])) + '</b><br>' + options.userOptions.metric + ': ' + fields['y'] + ' ' + options.userOptions.unit;
		}
		else {
			html = '<b>' + rdr_tstodate(Math.round(options.userOptions.original_values[fields[0].offset][0] / 1000)) + '</b><br />' + options.userOptions.metric + ' : ' + fields[0].value + ' ' + options.userOptions.unit;
		}

		return html;
	},

	buildSparkline: function() {
		this.buildOptions();

		if(this.chart_type === 'column') {
			var new_values = [];

			for(var i = 0; i < this.values.length; i++) {
				new_values[i] = this.values[i][1];
			}

			this.options.type = 'bar' ;
			this.charts = {
				'values': new_values,
				'options': this.options
			};

			$('#'+this.getId()).sparkline(new_values, this.options);
		}
		else {
			this.charts = {
				'values': this.values,
				'options': this.options
			};

			$('#'+this.getId()).sparkline(Ext.clone(this.values), this.options);
		}
	}
});
