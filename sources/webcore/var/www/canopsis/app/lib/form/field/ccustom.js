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
	extend: 'canopsis.lib.view.ccard',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.ccustom',

	border: false,

	sharedStore : undefined,

	name : 'ccustom',

	customForm: undefined,

    contentPanelDefault: {
                    xtype:'form',
                    width: '100%',
                    border:false,
                    //padding: '0 5 0 5',
                    margin: '30px',
                    autoScroll: true,
                    layout: {
                                type: 'vbox',
                                align: 'stretch'
                            },
                },

    buttonTpl: new Ext.Template([
                '<div class="ccustomButton">',
                    '{title}',
                '</div>'
            ]).compile(),

	initComponent: function() {
        this.callParent(arguments);

        this.panelIdByNode = {}

        this.customForm.push({
								"xtype":"hiddenfield",
								"name":"_id"
							},{
								"xtype":"hiddenfield",
								"name":"titleInWizard"
							})

    },

	afterRender: function() {
		this.callParent(arguments);

		//vars
		this.sourceStore = this.findParentByType('cwizard').childStores[this.sharedStore]

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
		var panelConfig = {items:Ext.clone(this.customForm)}

		data.titleInWizard = this.buildTitle(data)

		var panelNumber = this.addContent(panelConfig,data)
		if(data._id)
			this.panelIdByNode[data._id] = panelNumber
		var button = {title:data.titleInWizard,panelIndex: panelNumber}

		this.addButton(button)
	},

	removePanels: function(records){
		//if someone find something to do this in extjs, build this again
		if(!Ext.isArray(records))
			records = [records]
		for(var i = 0; i < records.length; i++){
			var nodeId = records[i].data.id
			if(this.panelIdByNode[nodeId] != undefined){
				var panelNumber = this.panelIdByNode[nodeId]
				var panel = this.panelByindex[panelNumber]
				var button = this.buttonByPanelIndex[panelNumber]
				console.log('panel and button')
				console.log(this.panelByindex)
				if(panel && button){
					panel.setDisabled(true)
					panel.hide()
					button.setDisabled(true)
					button.hide()
					delete this.panelByindex[panelNumber]
					delete this.panelIdByNode[nodeId]
					delete this.buttonByPanelIndex[panelNumber]
				}
			}
		}
			

	},

	setValue: function(data){
		console.log(data)
		Ext.Object.each(data, function(key, value, myself) {
			value._id = key
			this.addPanel(this.buildTitle(value),value)
		},this)
	},

	getValue: function(){
		var wizardChilds = this.contentPanel.items.items
        var output = {}
        for(var i = 0; i < wizardChilds.length; i++){
        	if(!wizardChilds.disabled){
	        	var values = wizardChilds[i].getValues(false, false, false, true)
	        	if(values._id){
		        	output[values._id] = values
		        	output[values._id]['titleInWizard'] = this.buildTitle(values)
		        }
	        }
        }
		//prevent sub item form to be submit individualy
		this.disable()

		return output
	},

	buildTitle: function(data){
		if(data.titleInWizard)
			return data.titleInWizard

		if(!data.co || !data.me)
			return data._id

		var title = data.co
		if(data.re)
			title += ' ' + data.re
		title += ' ' + data.me

		return title
	}

});