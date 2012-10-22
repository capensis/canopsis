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
Ext.define('widgets.text.text' , {
	extend: 'canopsis.lib.view.cwidget',
	alias: 'widget.text',
	initComponent: function() {
		//Initialisation of ext JS template
		this.myTemplate = new Ext.XTemplate('<div>' + this.text + '</div>');
		//Compilation of template ( to accelerate the render )
		this.myTemplate.compile();
		this.HTML = ''; // contains the html
		this.callParent(arguments); // Initialization globale of the template
	},
	onRefresh: function(data) {
		if (data)
		{
			if (this.nodes.length > 1)
			{
				var htmlArray = new Array();
				for (i in data)
				{
					var obj = data[i];
					obj.timestamp = rdr_tstodate(obj.timestamp);
					htmlArray.push(this.myTemplate.apply(obj));
				}
				this.HTML = htmlArray.join('');
			}
			else
			{
				if (this.nodes.length > 1)
					for (i in data)
						data[i].timestamp = rdr_tstodate(data[i].timestamp);
				//If data exist we apply the template on the node
				data.timestamp = rdr_tstodate(data.timestamp);
				this.HTML = this.myTemplate.apply(data);

			}
		}
		else
		{
			//otherwise we put the text contained in the field
			this.HTML = this.text;
		}
		this.setHtml(this.HTML);
	},
	getNodeInfo: function() {
		//we override the function : if there is'nt any nodeId specified we call the onRefresh function
		if (! this.nodeId)
		{
			this.onRefresh(false);
		}
		//we call the parent which is applied when there is a nodeId specified.
		this.callParent(arguments);
	}
});
