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

Ext.define('canopsis.lib.view.ccard' , {
    extend: 'Ext.window.Window',
    alias: 'widget.ccard',

    title: 'Wizard',
    closable: true,
    closeAction: 'destroy',
    width: 800,
    height: 500,
    layout: {
        type:'hbox',
        align:'stretch'
    },

    constrain: true,
    edit: false,
    activeButton: 0,
    wizardSteps : undefined,

    contentPanelDefault: {
                        xtype:'form',
                        width: '100%',
                        border:false,
                        padding: '0 5 0 5',
                        autoScroll: true
                    },

    buttonDefault: {
                    xtype:'container',
                    padding: 0,
                    hidden: false
                },

    buttonTpl: new Ext.Template([
                    '<div class="cwizardButton" style="{style}">',
                        '{title}',
                    '</div>'
                ]).compile(),

    items:[
            {
                _name: 'buttonPanel',
                width: 160,
                border:false,
                html: '<div class="cwizardBorderColor"></div>',
            },{
                flex: 1,
                _name: 'contentPanel',
                border:false,
                layout:'card',
            }
    ],

    bbar: [
            {text:'advance mode',_name:'bbarAdvance',enableToggle:true,hidden:true},
            {iconCls : 'icon-previous', text:'previous',_name:'bbarPrevious'},'->',
            {iconAlign: 'right',iconCls : 'icon-next', text: 'Next',_name:'bbarNext'},
            {iconAlign: 'right',iconCls : 'icon-save', text: 'Finish',_name:'bbarFinish', hidden:true}
        ],

    initComponent: function() {
        this.callParent(arguments);

        //stock ref to panels
        this.buttonPanel = this.down('panel[_name="buttonPanel"]')
        this.contentPanel = this.down('panel[_name="contentPanel"]')

        if(this.wizardSteps != undefined)
            this.addNewSteps(Ext.clone(this.wizardSteps))

        this.bbarNextButton = this.down('button[_name=bbarNext]')
        this.bbarFinishButton = this.down('button[_name=bbarFinish]')
        this.bbarPreviousButton = this.down('button[_name="bbarPrevious"]')

        //binding events
        this.bbarFinishButton.on('click',function(){
            this.fireEvent('save',this.widgetId,this.getValue())
            this.destroy()
        },this)
        this.bbarPreviousButton.on('click',function(){this.previousButton()},this)
        this.bbarNextButton.on('click',function(){this.nextButton()},this)
        this.down('button[_name="bbarAdvance"]').on('toggle',
            function(button,toggle){
                if(toggle)
                    this.switchToAdvancedMode()
                else
                    this.switchToSimpleMode()
            },this)
    },

    afterRender : function(){
        this.callParent(arguments);
        this.addSelectedCss()
        this.switchToSimpleMode()
    },

    //build functions
    addNewSteps: function(configObject){
        if(!Ext.isArray(configObject))
            configObject = [configObject]

        for(var i =0; i < configObject.length; i++){
            var button = configObject[i]

            button.panelIndex = this.addContent(Ext.clone(button))

            /*clean items*/
            delete button.items
            this.addButton(button)
        }
    },

    addButton: function(buttonConfig){
        buttonConfig = Ext.Object.merge(buttonConfig,this.buttonDefault)

        if(!buttonConfig.buttonIndex)
            buttonConfig.buttonIndex = this.buttonPanel.items.length

        buttonConfig.html = this.buttonTpl.apply({title:buttonConfig.title})

        var item = this.buttonPanel.insert(buttonConfig.buttonIndex,buttonConfig)

        //create click event for item
        if(item.getEl())
            item.getEl().on('click',
                        function(){this.fireEvent('click',this.buttonIndex,this.panelIndex,this)}
                    ,item)
        else
            item.on('afterrender',function(button){
                    button.getEl().on('click',
                        function(){this.fireEvent('click',this.buttonIndex,this.panelIndex,this)}
                        ,this)
                    })

        item.on('click',this.showStep,this)
    },

    addContent: function(extjs_obj_array){
        var _obj = Ext.Object.merge(extjs_obj_array,this.contentPanelDefault)

        //no title header
        _obj.title = undefined

        //one item
        if(_obj.layout != 'anchor' && _obj.items && _obj.items.length == 1)
            _obj.layout = 'fit'

        this.contentPanel.add(_obj)
        return this.contentPanel.items.length - 1
    },

    showStep: function(buttonNumber,panelNumber,buttonElement){
        if(buttonNumber == undefined )
            return;

        if(buttonElement && buttonElement.hasCls('cwizardDisabledButton'))
            return;

        if(panelNumber == undefined)
            panelNumber = this.getButton(buttonNumber).panelIndex

        this.removeSelectedCss(this.activeButton)
        this.contentPanel.getLayout().setActiveItem(panelNumber)
        this.activeButton = buttonNumber
        this.addSelectedCss(buttonNumber)

        if(!this.edit)
            this.checkDisplayFinishButton()
    },

    checkDisplayFinishButton: function(){
        if(this.activeButton == this.buttonPanel.items.length - 1){
            this.bbarNextButton.hide()
            this.bbarFinishButton.show()
        }else{
            this.bbarNextButton.show()
            this.bbarFinishButton.hide()
        }
    },

    nextButton: function(newIndex){
        if(!newIndex)
            newIndex = this.activeButton +1

        if(newIndex < this.buttonPanel.items.length){
            if(this._isDisabled(newIndex)){
                this.nextButton(newIndex + 1)
                return
            }
            this.showStep(newIndex);
        }
    },

    previousButton: function(newIndex){
        if(newIndex == undefined)
            newIndex = this.activeButton - 1

        if(newIndex >= 0){
            if(this._isDisabled(newIndex)){
                this.previousButton(newIndex - 1)
                return
            }
            this.showStep(newIndex);
        }

        if(!this.edit)
            this.checkDisplayFinishButton()
    },

    switchToAdvancedMode: function(){
        for(var i = 0; i < this.buttonPanel.items.items.length;i++)
            if(this.buttonPanel.items.items[i].advanced == true)
                this.buttonPanel.items.items[i].getEl().removeCls('cwizardDisabledButton')
    },

    switchToSimpleMode: function(){
        for(var i = 0; i < this.buttonPanel.items.items.length;i++)
            if(this.buttonPanel.items.items[i].advanced == true)
                this.buttonPanel.items.items[i].getEl().addCls('cwizardDisabledButton')

        if(this._isDisabled(this.activeButton))
            this.showStep(0)
    },

    //css manipulation
    removeSelectedCss: function(buttonNumber){
        var item = this.getButton(buttonNumber)
        if(item != undefined)
            item.removeCls('cwizardSelectedButton')
    },

    addSelectedCss: function(buttonNumber){
        var item = this.getButton(buttonNumber)
        if(item != undefined)
            this.getButton(buttonNumber).addCls('cwizardSelectedButton')
    },

    _isDisabled: function(buttonNumber){
        var button = this.getButton(buttonNumber)
        if(button)
            return button.getEl().hasCls('cwizardDisabledButton')
    },

    //getters

    getButton: function(buttonNumber){
        if(buttonNumber == undefined)
            buttonNumber = this.activeButton
        return this.buttonPanel.items.items[buttonNumber]
    },

    loadData: function(){

    },

    getValue: function(){
        if(this.beforeGetValue)
            this.beforeGetValue()
        
        var wizardChilds = this.contentPanel.items.items
        var _obj = {}
        for(var i = 0; i < wizardChilds.length; i++)
            _obj = Ext.Object.merge(_obj,wizardChilds[i].getValues(false, false, false, true))
        return _obj
    },

    setValue: function(data){
        var wizardChilds = this.contentPanel.items.items
        for (var i = 0; i < wizardChilds.length; i++) {
            var form = wizardChilds[i].getForm();
            form.setValues(data);
        }
    }

})