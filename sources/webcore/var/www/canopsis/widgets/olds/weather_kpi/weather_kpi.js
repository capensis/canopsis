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
Ext.define('widgets.weather_kpi.weather_kpi' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.weather_kpi',

	logAuthor: '[weather_kpi]',

	iconset: '1',

	initComponent: function() {
		log.debug('Init Weather kpi ' + this.id, this.logAuthor);
		log.debug(' + NodeId: ' + this.nodeId, this.logAuthor);

		this.callParent(arguments);


	},

	onRefresh: function(data) {

		if (this.iconset < 10) {
			var icon = '0' + this.iconset;
		}else {
			var icon = this.iconset;
		}

		if (data.state == 2) {
			this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_0-10'/></center>");
		}else if (data.state == 1) {
			this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_50-60'/></center>");
		} else {
			this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_90-100'/></center>");
		}
	}


});
