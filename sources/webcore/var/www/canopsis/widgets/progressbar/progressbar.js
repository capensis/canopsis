//need:app/lib/view/cperfstoreValueConsumerWidget.js
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

Ext.define('widgets.progressbar.progressbar', {
	extend: 'canopsis.lib.view.cperfstoreValueConsumerWidget',

	alias: 'widget.progressbar',
	logAuthor: '[progressBarWidget]',
	innerText: undefined,
	refresh_number: 0,
	wcontainer_autoScroll: true,
	wcontainer_layout: {type:'anchor'},
	bodyPadding: '5 5 5 5',

	colorBg: '#EEEEEE',
	colorStart: '#1BE01B',
	colorMid: '#FFCD43',
	colorEnd: '#E0251B',

	dispGrad: true,

	gHeight: 20,

	initComponent: function() {
		log.debug("initComponent", this.logAuthor);

		if(Ext.isIE) {
			this.dispGrad = false;
		}

		this.nodesByID = parseNodes(this.nodes);
		log.debug('nodesByID:', this.logAuthor);
		log.dump(this.nodesByID);

		// Color Scaling
		var colors = [this.colorStart, this.colorMid, this.colorEnd];
		this.colorScale = chroma.scale(colors);

		this.progressBars = {};

		this.callParent(arguments);

		log.debug('nodesByID:', this.logAuthor);
		log.dump(this.nodesByID);
	},

	getNodeInfo: function(from, to) {
		this.processNodes();

		this.refreshNodes(from, to);
	},

	setGradient: function(_id, value) {
		log.debug('setGradient: ' + _id + ", value: " + value, this.logAuthor);

		var lowColor = this.colorScale(value).hex();
		var hightColor = chroma.hex(lowColor).brighten(20).hex();

		pbEl = $('#'+_id + " .x-progress-bar");
		pbEl.css('background-image', "none");

		if(this.dispGrad) {
			pbEl.css('background', '-webkit-gradient(linear, left top, left bottom, from(' + hightColor + '), to(' + lowColor + '))');
			pbEl.css('background', '-webkit-linear-gradient(' + hightColor + ', ' + lowColor + ')');
			pbEl.css('background', '-moz-linear-gradient(' + hightColor + ', ' + lowColor + ')');
			pbEl.css('background', '-ms-linear-gradient(' + hightColor + ', ' + lowColor + ')');
			pbEl.css('background', '-o-linear-gradient(' + hightColor + ', ' + lowColor + ')');
			pbEl.css('background', 'linear-gradient(to bottom, '    + hightColor +', ' + lowColor + ')');
		}
		else{
			pbEl.css('background-color', lowColor);
		}
	},

	onRefresh: function(data) {
		var progress_bar = this;
		var pb_update = function(pb, value) {
			progress_bar.setGradient(pb.id, value);
		};

		var pb_afterrender = function(pb) {
			progress_bar.setGradient(pb.id, pb.value);
		};

		for(var i = 0; i < data.length; i++) {
			var item = data[i];
			var _id = item.node;

			// Create it
			log.debug('Item: ' + _id, this.logAuthor);
			log.dump(item);

			var node = this.nodesByID[_id];
			log.debug('Node:', this.logAuthor);
			log.dump(node);

			var label = node.label;

			log.debug(' + Label: '+label, this.logAuthor);

			var value = item.values[0][1];
			var max = item.max;

			//Extra field
			if(node.max) {
				max = node.max;
			}

			log.debug(' + Value: ' + value, this.logAuthor);
			log.debug(' + Max:   ' + max, this.logAuthor);

			var pct = 0;
			if(value && max) {
				pct = (value * 100) / max;
				pct = roundSignifiantDigit(pct, 2);
			}

			log.debug(' + Pct:   ' + pct, this.logAuthor);

			var text = pct + "%";
			var pb = undefined;

			// Check if pb already exist
			if(this.progressBars[_id]) {
				pb = this.progressBars[_id];
				log.debug('Update: ' + _id, this.logAuthor);
				pb.updateText(text);
				pb.updateProgress(pct/100);
				return;
			}

			pb = Ext.create('Ext.ProgressBar', {
				text: text,
				value: pct/100,
				flex:1,
				height: this.gHeight,
				cls:'widgets-progressbar',
				border: 1
			});

			pb.on("update", pb_update);
			pb.on("afterrender", pb_afterrender, this, {single: 1});

			this.wcontainer.add({
				layout: {
					type: 'hbox'
				},
				border: 0,
				margin: 4,
				height: this.gHeight,
				items: [
					{ border: 0, html: String(label), flex:1 },
					pb
				]
			});

			this.progressBars[_id] = pb;
		}
	},

	processNodes: function() {
		var post_params = [];

		Ext.Object.each(this.nodesByID, function(id, node) {
			post_params.push({
				id: id,
				metrics: node.metrics
			});
		}, this);

		this.post_params = {
			'nodes': Ext.JSON.encode(post_params),
			'aggregate_method' : this.aggregate_method,
			'aggregate_interval': this.aggregate_interval,
			'aggregate_max_points': this.aggregate_max_points
		};
	}
});
