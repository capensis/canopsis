/*
#--------------------------------
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
# ---------------------------------
*/

// initComponent -> afterContainerRender 	-> setchartTitle -> ready -> doRefresh -> onRefresh -> addDataOnChart
//                                			-> setOptions                             			-> getSerie
//											-> createChart

flag_tootlip_template = Ext.create('Ext.XTemplate',
	'<table>',
		'<tr>',
			'<td style="margin:3px;">',
				'<tpl if="event_type == \'user\'">',
					'<img src="widgets/stream/logo/ui.png" style="width: 32px;"></img>',
				'</tpl>',
				'<tpl if="event_type != \'user\'">',
					'<img src="widgets/stream/logo/{icon}.png" style="width: 32px;"></img>',
				'</tpl>',
			'</td>',
			'<td>',
				'<div style="margin:3px;">',
					'<tpl if="display_name">',
							'<b>{display_name}</b>',
					'</tpl>',
					'<tpl if="display_name == undefined">',
						'<b>{component}</b>',
						'<tpl if="resource">',
							'<b> - {resource}</b>',
						'</tpl>',
					'</tpl>',
					' <span style="color:grey;font-size:10px">{date}</span>',
					'<br/>{text}',
				'</div>',
			'</td>',
		'</tr>',
	'</table>',
	{compiled: true}
);

Ext.define('widgets.mini_chart.mini_chart' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.mini_chart',

	wcontainer_layout: {
        type: 'vbox',
        align: 'center'
    },


	logAuthor: '[mini_chart]',
	initComponent: function() {
		log.debug('initComponent', this.logAuthor) ;	
		log.debug('nodesByID:', this.logAuthor);
		log.dump(this.nodesByID);
		if(Ext.isArray(this.nodes))
			this.nodesByID = parseNodes(this.nodes);
		else
			this.nodesByID = expandAttributs(this.nodes)

		this.series = {} ;
		this.charts = { } ;
		this.callParent(arguments);
 	},
	doRefresh: function (from, to ) {
		this.from = from;	
		this.to = to;
		log.debug('Get values from ' + new Date(from) + ' to ' + new Date(to), this.logAuthor);
		if (this.nodesByID) {
			if (Ext.Object.getSize(this.nodesByID) != 0) {
				url = this.makeUrl(from, to);
				Ext.Ajax.request({
					url: url,
					scope: this,
					params: this.buildParams(from, to),
					method: 'POST',
					success: function(response) {
						var data = Ext.JSON.decode(response.responseText);
						data = data.data;
						this.onRefresh(data);
					},
					failure: function(result, request) {
						log.error('Ajax request failed ... (' + request.url + ')', this.logAuthor);
					}
				});
			} else {
				log.debug('No nodes specified', this.logAuthor);
				this.chart.showLoading(_('Please choose a valid metric in wizard'));
			}
		}
		this.callParent(arguments);
	},
	buildParams: function(oFrom, oTo) {
		//TODO: Rebuild this with new format !

		var now = Ext.Date.now();
		var post_params = [];

		Ext.Object.each(this.nodesByID, function(id, node, obj) {
			var nodeId = id;
			var from = oFrom;
			var to = oTo;
			if (this.aggregate_interval) {
				var aggregate_interval = this.aggregate_interval * 1000;

				if (this.aggregate_interval < global.commonTs['month']) {
					from = Math.floor(from / aggregate_interval) * aggregate_interval;
				}else {
					if (this.aggregate_interval >= global.commonTs['month'])
						from = moment.unix(from / 1000).startOf('month').unix() * 1000;
					if (this.aggregate_interval >= global.commonTs['year'])
						from = moment.unix(from / 1000).startOf('year').unix() * 1000;
				}
			}
			post_params.push({
				id: nodeId,
				metrics: node.metrics,
				from: parseInt(from / 1000),
				to: parseInt(to / 1000)
			});

		},this)

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
	buildOptions : function(info, values, serie_panel, i ) {
		
		var node = info['node'] ;
		
		//Find the print label
		var label ;
		if ( this.nodesByID[node]['label'] ) 
			label = this.nodesByID[node]['label'] ;
		else 
			label = info['metric'] ;

		//Find the unit
		var unit = '';
		if ( this.SeriePercent && info['max'] )
			unit = '%';
		else if (this.nodesByID[node]['u'] ) 
			unit = this.nodesByID[node]['u'] ;
		else if ( info.bunit )
			unit = info['bunit']; 

		//Find Colors for curve
		var colors = global.curvesCtrl.getRenderColors(label, i);
		var curve_color;
		if ( this.nodesByID[node]['curve_color'] ) 
			curve_color = this.nodesByID[node]['curve_color'] ;
		else
			curve_color = colors[0] ;

		var area_color ;
		if ( this.nodesByID[node]['area_color'] ) 
			area_color = this.nodesByID[node]['area_color'] ;
		else if ( Ext.isIE ) 
			area_color = curve_color;
		else
			area_color = this.lightenDarkenColor( curve_color, 50 );
	
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
						$(document ).find('.tooltip-minichart').css('border', '2px solid '+curve_color);
						if ( options.userOptions.chart_type == 'line_graph' ) 
							return '<b>' + rdr_tstodate(Math.round(fields['x']/1000)) + '</b><br>' + options.userOptions.metric + ': ' + fields['y'] + ' ' + options.userOptions.unit;
						return '<b>'+rdr_tstodate(Math.round( options.userOptions.original_values[fields[0].offset][0] / 1000 ) )+'</b><br />'+options.userOptions.metric+' : '+fields[0].value+' '+options.userOptions.unit ; 
					} 
		} ;
		return options;

	},
	lightenDarkenColor: function (col,amt,usePound,num,f,h,r,b,g) {
		if (col[0]=="#") {
			col = col.slice(1);
			usePound = (usePound==undefined?true:usePound)
		}
 
		num = parseInt(col,16);
		f=function(n) { return n>255?255:(n<0?0:n) }
		h=function(n) { return n.length<2?"0"+n:n }
   
		r = h(f((num >> 16) + amt).toString(16));
		b = h(f(((num >> 8) & 0x00FF) + amt).toString(16));
		g = h(f((num & 0x0000FF) + amt).toString(16));
 
		return (usePound?"#":"") + r + b + g;
	},
	onResize: function() {
		Ext.Object.each ( this.series, function(id, serie, obj) {
			serie.setWidth( this.getWidth() );
			serie.setHeight( this.getHeight() / Ext.Object.getSize(this.nodesByID) );
			var last;
			if ( serie.items.items.length == 2 ) 
			{
				serie.items.items[1].destroy();
				var serie_panel = serie.add ( {
					xtype: "panel",
					flex: 5,
					border: false
				} ) ;
				this.charts[id]['options']['width'] = serie_panel.getWidth() ;
				this.charts[id]['options']['height'] = serie_panel.getHeight() ;
				
				$('#'+serie_panel.getId() ).sparkline(this.charts[id]['values'], this.charts[id]['options'] );
			} else if ( serie.items.items.length == 3 ) {
				var last_cpn = serie.items.items[2].cloneConfig() ;
				serie.items.items[2].destroy();
				serie.items.items[1].destroy();
				var serie_panel = serie.add ( {
					xtype: "panel",
					flex: 5,
					border: false
				} ) ;
				this.charts[id]['options']['width'] = serie_panel.getWidth() ;
				this.charts[id]['options']['height'] = serie_panel.getHeight() ;
				$('#'+serie_panel.getId() ).sparkline(this.charts[id]['values'], this.charts[id]['options'] );
				serie.add(last_cpn);
			}
			//serie.items.items[1].el.id +" canvas").height(this.getHeight)  ;
		}, this ) ;	
	},
	parseValues: function(serie, values) {
	
		//MAKE A BETTER LOOP, JUST FOR TEST
		for (var i = 0; i < values.length; i++) {
			values[i][0] = values[i][0] * 1000;

			if (this.SeriePercent && serie.max > 0) { 
				values[i][1] = getPct(values[i][1], serie.max);
			}

			if (serie.invert)
				values[i][1] = - values[i][1];
		}
		return values;
	},
	onRefresh: function(data) {
		if ( Ext.Object.getSize(this.charts) > 0 || data.length > 0 ) {
				for ( var i=0; i < data.length; i++ ) {
						var info = data[i];

						var node = info['node'] ;
						var values = this.parseValues( info, info['values'] ) ;

						//clean the series
						if ( this.series[node] ) 
								this.series[node].removeAll();

						//Find the print label
						var label ;
						if ( this.nodesByID[node]['label'] ) 
								label = this.nodesByID[node]['label'] ;
						else 
								label = info['metric'] ;

						//Add a component with the print label
						this.series[node].add( {
							xtype:"panel",
							flex: 2,	
							bodyCls: "valigncenter",
							html: "<div><b>"+label+"</b></div>",
							border: false 
						});
			
						//We add the serie panel
						var serie_panel = this.series[node].add ( {
							xtype: "panel",
							flex: 4,
							border: false,
							style: {
								width: '95%',
								marginLeft: 'auto',
								marginRight:'auto'
							}
						} ) ;
			
						//we get options and build the serie panel
						var options = this.buildOptions(info, values, serie_panel, i) ;
						if ( this.chart_type == 'column' ) {
								var new_values = new Array();
								for ( var i=0; i < values.length; i++ ){
										new_values[i] = values[i][1]  ;
									}
								options.type = 'bar' ;
								this.charts[node] = { 'values': new_values, 'options': options } ;
								$('#'+serie_panel.getId() ).sparkline(new_values, options );
						} else {
								this.charts[node] = { 'values': values, 'options': options } ;
								$('#'+serie_panel.getId() ).sparkline(values, options );
						}
			
						if (this.nodesByID[node]['printed_value'] == 'last_value' ) {
								this.series[node].add( {
									xtype: "panel",
									flex: 2,
									bodyCls: "valigncenter padding-left",
									border: false,
									html : "<div><b>"+values[ values.length - 1][1]+" "+options['unit']+"</b></div>"
								}) ;
						} else if ( this.nodesByID[node]['printed_value'] == 'trend' ) {
								var nByID = { } ;
								nByID[node] = this.nodesByID[node] ;
								this.series[node].add( {
									xtype: "trends",
									flex: 3,
									bodyCls: "valigncenter-child padding-left",
									bodyStyle: {
										padding:"0px",
										margin:"0px",
										width: "100%",
									},
									height: this.series[node].getHeight(),
									border: false,
									nodes: nByID,
									nodesByID: nByID,
									margin:0,
									aggregate_method : this.aggregate_method,
									aggregate_interval: this.aggregate_interval,
									aggregate_max_points: this.aggregate_max_points,
									display_pct: this.nodesByID[node]['display_pct'],
									colorLow: this.nodesByID[node]['colorLow'],
									colorMid: this.nodesByID[node]['colorMid'],
									colorHight: this.nodesByID[node]['colorHight'],	
									from: this.from, 
									to: this.to,
									caller: "mini_chart"
								}) ; 
						}
						//We display the last value or the evolution
				}
		} else {
			this.getEl().mask(_('No data on interval'));
		}
		//this.callParent(arguments);
	},
	afterContainerRender: function() {
		var me = this ;
		Ext.Object.each(this.nodesByID, function(id, node, obj) {
			var serie  = {
				layout: {
					type: "hbox",
					align: "stretch"
				},
				style: {
					paddingLeft: '10px',
					paddingRight: '0px',
					paddingTop: '1px',
					paddingBottom: '5px'
				},
				width: me.getSize().width,
				height: me.getSize().height / Ext.Object.getSize(me.nodesByID),
				border: false
			};
			me.series[id] = me.wcontainer.add( serie ) ;
		} ) ;
		
		
		this.callParent();
	},
	beforeDestroy : function() {	
		Ext.Object.each( this.series, function( node_id, serie, obj) {
			serie.removeAll();
			serie.destroy();
		});
		this.wcontainer.removeAll();
		this.callParent();
	}

});
