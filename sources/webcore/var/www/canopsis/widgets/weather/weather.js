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
	icon_state_source: 'default',

	initComponent: function() {
		this.firstNodeIds = [];
		this.secondNodes = {};
		this.secondNodeIds = [];
		this.external_link_dict = {};

		log.debug('Initialize weather widget', this.logAuthor);
		if (this.exportMode || this.simple_display)
			this.wcontainer_autoScroll = false;

		this.configure();

		//---------------------Process nodes---------------
		for (var i = 0; i < this.nodes.length; i++) {
			var node = this.nodes[i];
			if (node._id)
				this.firstNodeIds.push(node._id);
			else
				this.firstNodeIds.push(node.id);
			if (node.link)
				this.external_link_dict[node.id] = node.link;
		}

		this.callParent(arguments);
	},

	doRefresh: function(from, to) {
		this.getNodes(this.firstNodeIds, this.firstNodesCallback);
	},

	getNodes: function(node_ids,callback) {
		log.debug('+ Get nodes', this.logAuthor);
		Ext.Ajax.request({
			url: this.uri,
			scope: this,
			method: 'GET',
			params: {ids: Ext.encode(node_ids)},
			success: callback,
			failure: function(result, request) {
				log.error('Impossible to get Node', this.logAuthor);
				global.notify.notify(_('Issue'), _("The selected selector can't be found"), 'info');
			}
		});
	},

	firstNodesCallback: function(response) {
		var nodes = Ext.JSON.decode(response.responseText).data;
		this.nodes = nodes;
		if (this.icon_state_source != 'default')
			this.secondNodeCheck();
		else
			this.populate();
	},

	secondNodeCheck: function() {
		//build list of second node ids if not already did
		if (this.secondNodeIds.length == 0) {
			log.debug('Building List of second ids to fetch', this.logAuthor);
			for (var i = 0; i < this.nodes.length; i++) {
				log.debug(' + Check if second node need for: ' + this.nodes[i]._id, this.logAuthor);
				var event_type = this.nodes[i].event_type;

				if (event_type != this.icon_state_source) {
					log.debug('  +  event type different from icon state source', this.logAuthor);
					if (event_type == 'selector') {
						if (this.nodes[i].sla_rk)
							this.secondNodeIds.push(this.nodes[i].sla_rk);
					}else {
						if (this.nodes[i].selector_rk)
							this.secondNodeIds.push(this.nodes[i].selector_rk);
					}
				}
			}
		}

		if (this.secondNodeIds.length > 0) {
			log.debug(' + Fetch secondary nodes', this.logAuthor);
			this.getNodes(this.secondNodeIds, this.secondNodesCallback);
		}else {
			log.debug(' + No need to fetch more nodes, populating', this.logAuthor);
			this.populate();
		}
	},

	secondNodesCallback: function(response) {
		var nodes = Ext.JSON.decode(response.responseText).data;
		for (var i = 0; i < nodes.length; i++)
			this.secondNodes[nodes[i]._id] = nodes[i];
		this.populate();
	},

	populate: function() {
		log.debug('Populate widget with ' + this.nodeId.length + ' elements.', this.logAuthor);
		this.wcontainer.removeAll();

		for (var i = 0; i < this.nodes.length; i++) {
			var node = Ext.clone(this.nodes[i]);

			//-----------------overload values----------------
			if (this.icon_state_source != 'default') {
				log.debug('Attempt to overide values with second node', this.logAuthor);
				if (node.sla_rk)
					var _id = node.sla_rk;
				if (node.selector_rk)
					var _id = node.selector_rk;
				var second_node = this.secondNodes[_id];

				if (second_node) {
					node.state = second_node.state;
					node.last_state_change = second_node.last_state_change;
				}
			}

			//------------------create config----------------
			var config = {
				data: node,
				link: this.external_link_dict[node._id],
				bg_color: (i % 2) ? this.bg_pair_color : this.bg_impair_color
			};
			var weather = Ext.create('widgets.weather.brick', Ext.Object.merge(config, this.base_config));
			this.wcontainer.add(weather);
		}
	},

	configure: function() {
		//-------------------define base config-------------------
		this.base_config = {
				iconSet: this.iconSet,
				state_as_icon_value: this.state_as_icon_value,
				icon_on_left: this.icon_on_left,
				exportMode: this.exportMode,
				display_report_button: this.display_report_button,
				display_derogation_icon: this.display_derogation_icon,
				external_link: this.external_link, //<-- helpdesk, change var name
				linked_view: this.linked_view,
				title_font_size: this.title_font_size,
				simple_display: this.simple_display,
				icon_state_source: this.icon_state_source,
				fullscreenMode: this.fullscreenMode,
				helpdesk: this.helpdesk
			};

		if (this.defaultPadding)
			this.base_config.padding = this.defaultPadding;

		if (this.defaultMargin)
			this.base_config.margin = this.defaultMargin;

		if (this.nodes.length == 1)
			this.base_config.anchor = '100% 100%';

	}

});
