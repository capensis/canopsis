//need:app/lib/form/cfield.js
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

Ext.define('canopsis.lib.form.field.cwlist' , {
	extend: 'Ext.panel.Panel',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.cwlist',

	border: false,

	nodes: undefined,

	cls: "cwlist",

	items: {
		xtype: "dataview",
		store: 'Widgets',
		tpl: [
			'<tpl for=".">',
				'<div class="thumb-wrap" id="{xtype}">',
					'<div class="thumb"><img height=60px width=60px src="{thumb}" title="{name}"></div>',
					'<span class="x-editable">{name}</span>',
				'</div>',
			'</tpl>',
			'<div class="x-clear"></div>'
		],
		multiSelect: false,
		trackOver: true,
		overItemCls: 'x-item-over',
		itemSelector: 'div.thumb-wrap',

		listeners: {
			selectionchange: function(dv, nodes){
				void(dv);

				var field = this.up('panel');
				field.nodes = nodes;
				field.fireEvent("select", field, nodes);
			}
		}
	},

	getValue: function() {
		if(this.nodes) {
			return this.nodes[0].raw.xtype;
		}
		else {
			return undefined;
		}
	},

	setValue: function(xtype) {
		this.nodes = [Ext.getStore("Widgets").findRecord('xtype', xtype, undefined, false, false, true)];
	}
});
