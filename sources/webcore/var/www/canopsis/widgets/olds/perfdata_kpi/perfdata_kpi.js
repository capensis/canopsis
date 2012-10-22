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
Ext.define('widgets.perfdata_kpi.perfdata_kpi' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.perfdata_kpi',

	logAuthor: '[perfdata_kpi]',

	iconset: '1',

	initComponent: function() {
		log.debug('Init Weather kpi ' + this.id, this.logAuthor);
		log.debug(' + NodeId: ' + this.nodeId, this.logAuthor);

		this.callParent(arguments);


	},

	onRefresh: function(data) {

		//formating iconset name
		if (this.iconset < 10) {
			var icon = '0' + this.iconset;
		}else {
			var icon = this.iconset;
		}

		var health = this.getHealth(data);
		//little tweak, because 0 = undefined ...
		if (health == 0) {
			health = 1;
		}

		if (health) {
			//round the result
			var roundHealth = 100 - (Math.round(health / 10) * 10);
			log.debug(roundHealth);
			switch (roundHealth) {
				case 0:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_0-10'/></center>");
					break;
				case 10:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_10-20'/></center>");
					break;
				case 20:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_20-30'/></center>");
					break;
				case 30:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_30-40'/></center>");
					break;
				case 40:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_40-50'/></center>");
					break;
				case 50:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_50-60'/></center>");
					break;
				case 60:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_60-70'/></center>");
					break;
				case 70:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_70-80'/></center>");
					break;
				case 80:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_80-90'/></center>");
					break;
				case 90:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_90-100'/></center>");
					break;
				case 100:
					this.setHtml("<center><span class='kpi kpi_iconSet" + icon + "_90-100'/></center>");
					break;
			}
		} else {
			this.setHtml('<center><div>Metric invalid or data missing, you can set this in the view editor</br>check the console for more details</div></center>');
		}
	}


});
