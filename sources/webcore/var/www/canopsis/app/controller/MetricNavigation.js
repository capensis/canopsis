//need:app/view/MetricNavigation/MetricNavigation.js,widgets/line_graph/line_graph.js
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
Ext.define('canopsis.controller.MetricNavigation', {
	extend: 'Ext.app.Controller',

	stores: [],
	models: [],

	logAuthor: '[controller][MetricNavigation]',

	views: ['MetricNavigation.MetricNavigation'],

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.control({
			'MetricNavigation': {
				afterrender: this._bindMetricNavigation
			}
		});

		this.callParent(arguments);
	},

	_bindMetricNavigation: function(panel) {
		log.debug('Binding events', this.logAuthor);

		this.panel         = panel;
		this.tabPanel      = panel.tabPanel;
		this.renderPanel   = panel.renderPanel;
		this.metricTab     = panel.metricTab;
		this.renderContent = panel.renderContent;

		// Button bindings
		panel.buttonCancel.on('click', this._buttonCancel, this);
		panel.buttonDisplay.on('click', this._buttonDisplay, this);
		this.tabPanel.on('collapse', this._refreshLayout, this);
		this.tabPanel.on('expand', this._refreshLayout, this);

		// if nodes specified with creation set it automaticaly
		if(panel.nodes.length !== 0) {
			this._addGraph(panel.nodes);
		}
	},

	_buttonCancel: function() {
		log.debug('Click on cancel button', this.logAuthor);
		tab = Ext.getCmp('main-tabs').getActiveTab();
		tab.close();
	},

	_buttonDisplay: function() {
		log.debug('Click on display button', this.logAuthor);
		metrics = this.metricTab.getValue();

		this._addGraph(metrics);
	},

	_addGraph: function(metrics) {
		this.renderContent.removeAll(true);

		timePeriod = this._getTime();

		//add one graph per node
		for(var i = 0; i < metrics.length; i++) {
			var item = this._createGraph([metrics[i]]);

			//set time after first render (avoid useless ajax request)
			item.nodes = [metrics[i]];
			item.processNodes();
			item._doRefresh(timePeriod.from * 1000, timePeriod.to * 1000);
		}

	},

	_createGraph: function() {
		log.debug('Adding graph', this.logAuthor);

		var config = {
			SeriesType: 'line',
			reportMode: true,
			extend: 'Ext.container.Container',
			width: '49%',
			height: 200,
			layout: 'fit'
		};

		var graph = Ext.widget('line_graph', config);
		this.renderContent.add(graph);

		return graph;
	},

	_refreshLayout: function() {
		for(var i = 0; i < this.renderContent.items.length; i++) {
			this.renderContent.items.items[i].onResize();
		}
	},

	_getTime: function() {
		log.debug('Set time period on graphs', this.logAuthor);

		//get time values
		var fromDate = this.panel.fromDate.getValue();
		var toDate   = this.panel.toDate.getValue();
		var fromHour = this.panel.fromHour.getSubmitData().fromHour;
		var toHour   = this.panel.toHour.getSubmitData().toHour;

		//compute from Hour
		arrayFromHour = fromHour.split(':');
		fromHour      = (arrayFromHour[0] * global.commonTs.hours) + (arrayFromHour[1] * 60);

		log.debug('from Hour ts : ' + fromHour);

		//compute to Hour
		arrayToHour = toHour.split(':');
		toHour      = (arrayToHour[0] * global.commonTs.hours) + (arrayToHour[1] * 60);

		log.debug('from Hour ts : ' + toHour);

		fromDate = Ext.Date.format(fromDate, 'U');
		toDate   = Ext.Date.format(toDate, 'U');

		var from = parseInt(fromDate) + parseInt(fromHour);
		var to   = parseInt(toDate) + parseInt(toHour);

		return {
			from: from,
			to: to
		};
	}
});
