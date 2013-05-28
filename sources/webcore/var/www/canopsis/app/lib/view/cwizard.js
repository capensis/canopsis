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
Ext.define('canopsis.lib.view.cwizard' , {
	extend: 'canopsis.lib.view.ccard',

	alias: 'widget.cwizard',

	requires: [
				'canopsis.lib.form.field.cinventory',
				'canopsis.lib.form.field.cmetric',
				'canopsis.lib.form.field.ccustom',
				'canopsis.lib.form.field.cfilter',
				'canopsis.lib.form.field.ctag',
				'canopsis.lib.form.field.cfieldset',
				'canopsis.lib.form.field.cdate',
				'canopsis.lib.form.field.cduration',
				'canopsis.lib.form.field.cduration',
				'canopsis.lib.form.field.ccolorfield'
			],

	wizardSteps: {
				title: _('Choose widget'),
				//description : _('choose the type of widget you want, its title, and refresh interval'),
				items: [{
							xtype:'cfieldset',
							title:_('General options'),
							items:[{
								xtype: 'combo',
								store: 'Widgets',
								queryMode: 'local',
								forceSelection: true,
								fieldLabel: _('Type'),
								name: 'xtype',
								editable: false,
								displayField: 'name',
								valueField: 'xtype',
								//autoLoad: true,
								//value: 'empty',
								allowBlank: false
							},{
								xtype: 'displayfield',
								name: 'description',
								isFormField: false,
								fieldLabel: _('Description')
							},{
								xtype: 'checkbox',
								fieldLabel: _('Show border'),
								checked: false,
								name: 'border',
								uncheckedValue: false
							},{
								xtype: 'checkbox',
								fieldLabel: _('Auto title') + ' ' + _('if available'),
								checked: true,
								inputValue: true,
								uncheckedValue: false,
								name: 'autoTitle'
							},{
								xtype: 'textfield',
								fieldLabel: _('Title') + ' (' + _('optional') + ')',
								name: 'title'
							},{
								xtype: 'combobox',
								name: 'refreshInterval',
								fieldLabel: _('Refresh interval'),
								queryMode: 'local',
								editable: false,
								displayField: 'text',
								valueField: 'value',
								value: 300,
								store: {
									xtype: 'store',
									fields: ['value', 'text'],
									data: [
										{value: 0,	text: 'None'},
										{value: 60,	text: '1 minutes'},
										{value: 300,	text: '5 minutes'},
										{value: 600,	text: '10 minutes'},
										{value: 900,	text: '15 minutes'},
										{value: 1800,	text: '30 minutes'},
										{value: 3600,	text: '1 hour'}
									]
								}
							}]
						}]
			},
	
	afterRender: function(){
		this.callParent(arguments);

		this.childStores = {}

		if(this.data){
			this.edit = true
			this.bbarNextButton.hide()
			this.bbarNextButton.hide()
			this.bbarPreviousButton.hide()
            this.bbarFinishButton.show()
        }

		var combo = this.down('combobox[name=xtype]')
		if(combo.rendered)
			this.comboAfterRender()
		else
			combo.on('afterRender',this.comboAfterRender,this)
	},

	comboAfterRender: function(){
		var combo = this.down('combobox[name=xtype]')
		combo.on('select', this.addOptionPanel, this);
		if(this.data){
			combo.setValue(this.data.xtype)
			var record = combo.store.findRecord('xtype',this.data.xtype,undefined,false,false,true)
			this.addNewSteps(Ext.clone(record.raw.options))
			this.setValue(this.data)
			combo.setDisabled(true)
		}
	},

	cleanPanels: function(){
		var contentToDel = [], buttonToDel = []
		
		//stock ref in array, otherwise the array is modified
		//during the removing loop
		for(var i = 1; i < this.contentPanel.items.items.length; i++)
			contentToDel.push(this.contentPanel.items.items[i])
		for(var i = 1; i < this.buttonPanel.items.items.length; i++)
			buttonToDel.push(this.buttonPanel.items.items[i])

		for(var i = 0; i < contentToDel.length; i++)
			this.contentPanel.remove(contentToDel[i],true)
		for(var i = 0; i < buttonToDel.length; i++)
			this.buttonPanel.remove(buttonToDel[i],true)
	},

	addOptionPanel: function(combo,records,opts){
		this.cleanPanels()
		this.addNewSteps(Ext.clone(records[0].raw.options))
	},

	beforeGetValue: function(){
		this.down('combobox[name=xtype]').setDisabled(false)
	}

});
