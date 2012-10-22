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
Ext.define('widgets.host_header.host_header' , {
	extend: 'canopsis.lib.view.cwidget',

	alias: 'widget.host_header',

	tpl: '<h1>{connector_name} > {source_type} > {component}</h1>',

	htmlTpl: false,

	initComponent: function() {
		this.refreshInterval = 0;
		this.htmlTpl = new Ext.Template(this.tpl, {compiled: true});
		this.callParent(arguments);
	},

	onRefresh: function(data) {
		this.data = data;
		this.setHtmlTpl(this.htmlTpl, data);
	}
});
