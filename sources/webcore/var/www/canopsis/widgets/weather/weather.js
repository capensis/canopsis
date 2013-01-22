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

Ext.define('widgets.weather.weather' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.weather',
	logAuthor: '[widget][weather]',
	border: false,

	cls: 'widget-weather',

	wcontainer_autoScroll: true,
	wcontainer_layout: 'anchor',

	selector_record: undefined,
	sla_id: undefined,

	//brick options
	iconSet: '01',
	icon_on_left: false,
	defaultHeight: undefined,
	defaultPadding: undefined,
	defaultMargin: undefined,
	state_as_icon_value: false,
	selector_state_as_icon_value: false,
	bg_impair_color: undefined,
	bg_pair_color: '#FFFFFF',

	base_config: undefined,

	simple_display: false,
	title_font_size: 14,
	display_report_button: false,
	display_derogation_icon: true,
	external_link: undefined,
	linked_view: undefined,
	helpdesk: undefined,


	initComponent: function() {
		if (this.exportMode || this.simple_display)
			this.wcontainer_autoScroll = false;
		this.callParent(arguments);
		log.debug('Initialize weather widget', this.logAuthor);

		//---------------------Process nodes---------------
		var ids = [];
		this.metric_options = {};
		if (this.nodeId && this.nodeId.length != 0) {
			if (typeof(this.nodeId[0]) != 'string') {
				for (var i = 0; i < this.nodeId.length; i++) {
					var node = this.nodeId[i];
					ids.push(node.id);
					if (node.link)
						this.metric_options[node.id] = node.link;
				}
				this.nodeId = ids;
			}
		}
	},

	afterContainerRender: function() {
		this.configure();
		this.callParent(arguments);
	},

	doRefresh: function(from, to) {
		log.debug('Do refresh', this.logAuthor);
		if (this.nodeId) {
			if (!this.reportMode || this.exportMode)
				this.getNodes();
			else
				this.getPastNodes(from, to);
		}
	},

	getNodes: function() {
		log.debug('+ Get nodes', this.logAuthor);
		Ext.Ajax.request({
			url: this.uri,
			scope: this,
			method: 'GET',
			params: {ids: Ext.encode(this.nodeId)},
			success: function(response) {
				var nodes = Ext.JSON.decode(response.responseText).data;
				var nodes_obj = {};

				for (var i = 0; i < nodes.length; i++) {
					nodes[i].nodeId = getMetaId(nodes[i].component, nodes[i].resource, nodes[i].metric);
					nodes_obj[nodes[i]._id] = nodes[i];
				}

				this.nodes = nodes_obj;

				if (this.selector_state_as_icon_value)
					this.getSelectorNodes(nodes_obj);
				else
					this.populate(nodes_obj);
			},
			failure: function(result, request) {
				log.error('Impossible to get Node', this.logAuthor);
				global.notify.notify(_('Issue'), _("The selected selector can't be found"), 'info');
			}
		});
	},

	getSelectorNodes: function(nodes) {
		var selector_list = [];

		for (var i = 0; i < nodes.length; i++)
			if (nodes[i].selector_rk)
				selector_list.push(nodes[i].selector_rk);


		Ext.Ajax.request({
			url: this.uri,
			scope: this,
			method: 'GET',
			params: {ids: Ext.encode(selector_list)},
			success: function(response) {
				var nodes = Ext.JSON.decode(response.responseText).data;
				node_dict = {};

				for (var i = 0; i < nodes.length; i++)
					node_dict[nodes[i]._id] = nodes[i];

				this.selector_nodes = node_dict;
				this.populate(this.nodes);
			},
			failure: function(result, request) {
				log.error('Impossible to get Node', this.logAuthor);
				global.notify.notify(_('Issue'), _("The selected selector can't be found"), 'info');
			}
		});

	},

	getPastNodes: function(from, to) {
		log.debug(' + Request data from: ' + from + ' to: ' + to, this.logAuthor);
		//log.dump(this.nodes)
		//--------------------Prepare post params-----------------

		var post_params = [];
		var me = this;
		Ext.Object.each(this.nodes, function(key, value, myself) {
			post_params.push({id: me.nodes[key].node_meta_id})
		});

		//-------------------------send request--------------------
		Ext.Ajax.request({
			url: '/perfstore/values/' + parseInt(from/1000) + '/' + parseInt(to/1000),
			params: {'nodes': Ext.JSON.encode(post_params)},
			scope: this,
			success: function(response) {
				var data = Ext.JSON.decode(response.responseText);
				data = data.data;
				this.report(data);
			},
			failure: function(result, request) {
				log.error('Impossible to get sla informations on the given time period', this.logAuthor);
			}
		});
	},

	generate_node_meta_id: function() {
		var me = this;
		Ext.Object.each(this.nodes, function(node_id, node, myself) {
			//build selector get id or node id
			if (me.selector_state_as_icon_value && me.selector_nodes[node.selector_rk]) {

				var selector 	= me.selector_nodes[node.selector_rk];
				var component 	= selector.component;
				var resource 	= selector.resource;
				var metric 		= 'cps_state';

				if (resource)
					var selector_id = getMetaId(component, resource, metric);
				else
					var selector_id = getMetaId(component, undefined, metric);

				node.selector_meta_id = selector_id;

			}else {
				var component 	= node.component;
				var resource 	= node.resource;

				if (node.event_type == 'selector')
					var metric = 'cps_state';
				else
					var metric = 'cps_pct_by_state_0';

				if (resource)
					var node_meta_id = getMetaId(component, resource, metric);
				else
					var node_meta_id = getMetaId(component, undefined, metric);

				node.node_meta_id = node_meta_id;
			}
		});
	},

	configure: function() {
		//-------------------define base config-------------------
		this.base_config = {
				iconSet: this.iconSet,
				state_as_icon_value: this.state_as_icon_value,
				icon_on_left: this.icon_on_left,
				exportMode: this.exportMode
			};
		
		if (this.defaultPadding)
			this.base_config.padding = this.defaultPadding;

		if (this.defaultMargin)
			this.base_config.margin = this.defaultMargin;

		if (this.nodes.length == 1) 
			this.base_config.anchor = '100% 100%';
		/*else
			if (this.defaultHeight)
				this.base_config.height = parseInt(this.defaultHeight, 10);
		*/
	},

	populate: function(data) {
		log.debug('Populate widget with ' + this.nodeId.length + ' elements.', this.logAuthor);

		this.generate_node_meta_id();

		this.wcontainer.removeAll();
		var debug_loop_count = 0;

		for (var i = 0; i < this.nodeId.length; i++) {
			var node_id = this.nodeId[i];

			if (data[node_id]) {
				var node = data[node_id];

				log.debug('Build brick for node ' + node._id, this.logAuthor);

				var config = {
					nodeId: node.nodeId,
					data: node,
					display_report_button: this.display_report_button,
					display_derogation_icon: this.display_derogation_icon,
					brick_number: i,
					external_link: this.external_link,
					linked_view: this.linked_view,
					title_font_size: this.title_font_size,
					simple_display: this.simple_display,
					selector_state_as_icon_value: this.selector_state_as_icon_value,
					link: this.metric_options[node_id],
					fullscreenMode: this.fullscreenMode,
					helpdesk: this.helpdesk
				};

				if (node.node_meta_id)
					config.node_meta_id = node.node_meta_id;

				if (node.selector_meta_id)
					config.selector_meta_id = node.selector_meta_id;

				//if selector_rk, put previously get selector in new object
				if (node.selector_rk)
					if (this.selector_nodes)
						if (this.selector_nodes[node.selector_rk])
							config.selector = this.selector_nodes[node.selector_rk];


				if ((i % 2) == 0)
					config.bg_color = this.bg_pair_color;
				else
					config.bg_color = this.bg_impair_color;

				var meteo = Ext.create('widgets.weather.brick', Ext.Object.merge(config, this.base_config));
				this.wcontainer.add(meteo);
				debug_loop_count += 1;
			}
		}

		log.debug('Finished to populate weather widget with ' + debug_loop_count + ' elements', this.logAuthor);

		if (this.exportMode) {
			log.debug('Exporting mode enable, fetch data', this.logAuthor);
			this.getPastNodes(export_from, export_to);
		}
	},

	report: function(data) {
		log.debug(' + Enter report function', this.logAuthor);
		var bricks = this.wcontainer.items.items;
		var dataById = {}

		for (var i = 0; i < data.length; i++)
			dataById[data[i].node] = data[i]

		Ext.Array.each(bricks, function(brick) {
			var new_values = dataById[brick.selector_meta_id];

			if (! new_values)
				new_values = dataById[brick.node_meta_id];

			if (new_values && new_values.values.length > 0) {
				log.debug(' + New values for ' + brick.event_type + ' ' + brick.component, this.logAuthor);
				brick.buildReport(new_values);
			}else {
				log.debug(' + No data recieved for ' + brick.event_type + ' ' + brick.component, this.logAuthor);
				brick.buildEmpty();
			}
		},this);
	}

});
