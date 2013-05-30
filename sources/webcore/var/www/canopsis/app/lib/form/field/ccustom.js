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

Ext.define('canopsis.lib.form.field.ccustom' , {
	extend: 'Ext.panel.Panel',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.ccustom',

	border: false,

	sharedStore : undefined,

	layout:'card',

	tbar: [
			{xtype:'button',iconCls : 'icon-previous'},
			'->',
			{xtype:'button',iconCls : 'icon-next'},
		],

	customForm: undefined,

	initComponent: function() {
        this.callParent(arguments);
       	this.matchingDict = {}
    },

	afterRender: function() {
		this.callParent(arguments);

		//vars
		this.sourceStore = this.findParentByType('cwizard').childStores[this.sharedStore]


		//bindings
		this.down('button[iconCls=icon-previous]').on('click',function(){
			var panel = this.getLayout().getPrev()
			if(panel)
				this.getLayout().prev()
		},this)
		this.down('button[iconCls=icon-next]').on('click',function(){
			var panel = this.getLayout().getNext()
			if(panel)
				this.getLayout().next()
		},this)

		this.sourceStore.on('add',function(store,records){
			this.addPanels(records)
		},this)
		this.sourceStore.on('remove',function(store,records){
			this.removePanels(records)
		},this)

	},

	addPanels: function(records){
		if(!Ext.isArray(records))
			records = [records]
		for(var i = 0; i < records.length; i++)
			this.addPanel(records[i].data.id, records[i].data)
	},

	addPanel: function(nodeId, data){
		var elem = this.add({
			xtype:'form',
			nodeId: nodeId,
			margin: '10 5 5 5',
			items: Ext.clone(this.customForm),
			layout: {
				type: 'vbox',
				align : 'left',
				//pack  : 'start',
			},
			border: false,
		/*	fieldDefaults: {
				//anchor:'50%'
			},*/
			dataToLoad: Ext.clone(data),
			listeners:{
				afterrender: function(){
								this.getForm().setValues(this.dataToLoad)
							}
			}
		})

		this.matchingDict[nodeId] = elem
	},

	removePanels: function(records){
		if(!Ext.isArray(records))
			records = [records]
		for(var i = 0; i < records.length; i++){
			var nodeId = records[i].data.id
			if(this.matchingDict[nodeId])
				this.remove(this.matchingDict[nodeId])
		}
	},

	setValue: function(data){
		Ext.Object.each(data, function(key, value, myself) {
			this.addPanel(key,value)
		},this)
	},

	getValue: function(){
		var output = {}
		Ext.Object.each(this.matchingDict, function(key, value, myself) {
			output[key] = value.getForm().getValues()
		},this)

		//prevent sub item form to be submit individualy
		this.disable()

		return output
	}

});