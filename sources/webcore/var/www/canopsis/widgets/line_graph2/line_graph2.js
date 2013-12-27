//need:app/lib/view/cwidgetGraph.js
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
Ext.define('widgets.line_graph2.line_graph2', {
	extend: 'canopsis.lib.view.cwidgetGraph',
	alias: 'widget.line_graph2',

	logAuthor: '[widgets][graph][line]',

	initComponent: function() {
		this.callParent(arguments);
	},

	setChartOptions: function() {
		this.callParent(arguments);

		this.options.series = {
			lines: {
				show: true,
				fill: true
			},
			points: {
				show: false
			}
		};
	},

	getSerieForNode: function(nodeid) {
		var node = this.nodesByID[nodeid];

		var serie = {
			label: node.label,
		}
	}
});