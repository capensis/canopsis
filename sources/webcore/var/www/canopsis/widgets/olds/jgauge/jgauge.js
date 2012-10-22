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
Ext.define('widgets.jgauge.jgauge' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.jgauge',

	logAuthor: '[jgauge]',

	initComponent: function() {
		log.debug('Init jgauge kpi ' + this.id, this.logAuthor);
		log.debug(' + NodeId: ' + this.nodeId, this.logAuthor);
		this.callParent(arguments);
	},

	onRefresh: function(data) {
		if (!this.jgauge) {
			this.setGauge();
		}

		var health = (this.getHealth(data));

		if (health) {
			this.jgauge.setValue(health);
		}else {
			this.setHtml('<center><div>Impossible to display gauge because</br>input data are invalid (check console)</div></center>');
		}
	},

	setGauge: function() {
		this.setHtml("<div id='jgauge-" + this.id + "' class='jgauge'></div>");
		this.jgauge = new jGauge();
		this.jgauge.id = 'jgauge-' + this.id;

		this.jgauge.ticks.count = 5;
		this.jgauge.ticks.start = 0;
		this.jgauge.ticks.end = 100;

		this.jgauge.label.suffix = '%';

		var orig_start = -200;
		var orig_end = 20;

		var total = 220;

		this.jgauge.range.radius = 50;
		this.jgauge.range.thickness = 10;
		//this.jgauge.range.start = (( 90 * total) / 100) + orig_start ;
		//this.jgauge.range.end = (( 100 * total) / 100) + orig_start;
		this.jgauge.range.color = 'rgba(0, 255, 0, 0.5)';
		this.jgauge.init();
	}
});
