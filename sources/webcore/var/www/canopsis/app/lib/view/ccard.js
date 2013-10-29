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
Ext.define('canopsis.lib.view.ccard' , {
	extend: 'Ext.panel.Panel',
	alias: 'widget.ccard',

	layout: {
		type: 'hbox',
		align: 'stretch'
	},

	edit: false,
	activeButton: 0,
	wizardSteps: undefined,

	advanceMode: false,

	contentPanelDefault: {
		xtype: 'form',
		width: '100%',
		border: false,
		padding: '0 5 0 5',
		autoScroll: true,
		layout: {
			type: 'vbox',
			align: 'stretch'
		},
		defaults: {
			autoScroll: true
		}
	},

	buttonDefault: {
		xtype: 'container',
		padding: 0
	},

	buttonTpl: new Ext.Template([
		'<div class="cwizardButton" style="{style}">',
			'{title}',
		'</div>'
	]).compile(),

	items: [
		{
			_name: 'buttonPanel',
			width: 160,
			border: false,
			autoScroll: true,
			html: '<div class="cwizardBorderColor"></div>'
		},{
			flex: 1,
			_name: 'contentPanel',
			border: false,
			layout: 'card'
		}
	],

	initComponent: function() {
		this.callParent(arguments);

		this.panelByindex = {};
		this.buttonByPanelIndex = {};

		//stock ref to panels
		this.buttonPanel = this.down('panel[_name="buttonPanel"]');
		this.contentPanel = this.down('panel[_name="contentPanel"]');

		if(this.wizardSteps !== undefined) {
			this.addNewSteps(Ext.clone(this.wizardSteps));
		}
	},

	afterRender: function() {
		this.callParent(arguments);
		this.addSelectedCss();
		this.switchToSimpleMode();
	},

	//build functions
	addNewSteps: function(configObject) {
		if(!Ext.isArray(configObject)) {
			configObject = [configObject];
		}

		for(var i = 0; i < configObject.length; i++) {
			var button = configObject[i];
			var content = Ext.clone(button);

			if(button.advanced === undefined) {
				button.advanced = false;
			}

			if(!this.advanceMode && button.advanced === true) {
				button.hidden = true;
			}

			button.panelIndex = this.addContent(content);

			/* clean items */
			delete button.items;
			this.addButton(button);
		}
	},

	addButton: function(buttonConfig) {
		buttonConfig = Ext.Object.merge(buttonConfig, this.buttonDefault);

		if(!buttonConfig.buttonIndex) {
			buttonConfig.buttonIndex = this.buttonPanel.items.length;
		}

		buttonConfig.html = this.buttonTpl.apply({title: buttonConfig.title});

		var item = this.buttonPanel.insert(buttonConfig.buttonIndex, buttonConfig);

		this.buttonByPanelIndex[buttonConfig.panelIndex] = item;

		//create click event for item
		if (item.getEl()) {
			item.getEl().on('click', function() {
				this.fireEvent('click', this.buttonIndex, this.panelIndex, this);
			} ,item);
		}
		else {
			item.on('afterrender', function(button) {
				button.getEl().on('click', function() {
					this.fireEvent('click', this.buttonIndex, this.panelIndex, this);
				}, this);
			});
		}

		item.on('click', this.showStep, this);
	},

	addContent: function(extjs_obj_array,data) {
		var _obj = Ext.Object.merge(extjs_obj_array, this.contentPanelDefault);

		//no title header
		_obj.title = undefined;

		//one item
		if(_obj.layout !== 'anchor' && _obj.items && _obj.items.length === 1) {
			_obj.layout = 'fit';
		}

		var panel = this.contentPanel.add(_obj);

		if(data) {
			panel.getForm().setValues(data);
		}

		this.panelByindex[this.contentPanel.items.length - 1] = panel;

		return this.contentPanel.items.length - 1;
	},

	showStep: function(buttonNumber,panelNumber,buttonElement) {
		if(buttonNumber === undefined) {
			return;
		}

		if(buttonElement && buttonElement.hasCls('cwizardDisabledButton')) {
			return;
		}

		if(panelNumber === undefined) {
			panelNumber = this.getButton(buttonNumber).panelIndex;
		}

		this.removeSelectedCss(this.activeButton);
		this.contentPanel.getLayout().setActiveItem(panelNumber);
		this.activeButton = buttonNumber;
		this.addSelectedCss(buttonNumber);

		if(!this.edit) {
			this.fireEvent('buttonChanged');
		}
	},

	cleanPanels: function() {
		var contentToDel = [];
		var buttonToDel = [];

		//stock ref in array, otherwise the array is modified
		//during the removing loop
		var nb_fixed_tab = 2;

		for(var i = nb_fixed_tab; i < this.contentPanel.items.items.length; i++) {
			contentToDel.push(this.contentPanel.items.items[i]);
		}

		for(i = nb_fixed_tab; i < this.buttonPanel.items.items.length; i++) {
			buttonToDel.push(this.buttonPanel.items.items[i]);
		}

		for(i = 0; i < contentToDel.length; i++) {
			this.contentPanel.remove(contentToDel[i], true);
		}

		for(i = 0; i < buttonToDel.length; i++) {
			this.buttonPanel.remove(buttonToDel[i], true);
		}
	},

	nextButton: function(newIndex) {
		if(!newIndex) {
			newIndex = this.activeButton + 1;
		}

		if(newIndex < this.buttonPanel.items.length) {
			if(this._isDisabled(newIndex)) {
				this.nextButton(newIndex + 1);
				return;
			}

			this.showStep(newIndex);
		}
	},

	previousButton: function(newIndex) {
		if(newIndex === undefined) {
			newIndex = this.activeButton - 1;
		}

		if(newIndex >= 0) {
			if(this._isDisabled(newIndex)) {
				this.previousButton(newIndex - 1);
				return;
			}

			this.showStep(newIndex);
		}

		if(!this.edit) {
			this.fireEvent('buttonChanged');
		}
	},

	switchToAdvancedMode: function() {
		this.advanceMode = true;

		for(var i = 1; i < this.buttonPanel.items.items.length; i++) {
			this.getButton(i).show();
		}
	},

	switchToSimpleMode: function() {
		this.advanceMode = false;

		for(var i = 1; i < this.buttonPanel.items.items.length; i++) {
			var button = this.getButton(i);

			if(button.advanced === true) {
				button.hide();
			}
		}

		if(this._isDisabled(this.activeButton)) {
			this.showStep(1);
		}
	},

	//css manipulation
	removeSelectedCss: function(buttonNumber) {
		var item = this.getButton(buttonNumber);

		if(item !== undefined) {
			item.removeCls('cwizardSelectedButton');
		}
	},

	addSelectedCss: function(buttonNumber) {
		var item = this.getButton(buttonNumber);

		if(item !== undefined) {
			this.getButton(buttonNumber).addCls('cwizardSelectedButton');
		}
	},

	_isDisabled: function(buttonNumber) {
		var button = this.getButton(buttonNumber);

		if(button) {
			return button.hidden;
		}
	},

	isLastPanel : function(){
		var button = undefined;
		var i = this.activeButton + 1;

		do {
			button = this.getButton(i);
			i += 1;

		} while (button && (this.advanceMode !== button.advanced));

		if(button) {
			return false;
		}
		else {
			return true;
		}

	},

	//getters

	getButton: function(buttonNumber) {
		if(buttonNumber === undefined) {
			buttonNumber = this.activeButton;
		}

		if(buttonNumber >=  this.buttonPanel.items.items.length) {
			return undefined;
		}

		return this.buttonPanel.items.items[buttonNumber];
	},

	loadData: function() {
	},

	getValue: function() {
		if(this.beforeGetValue) {
			this.beforeGetValue();
		}

		var wizardChilds = this.contentPanel.items.items;
		var _obj = {};

		for(var i = 0; i < wizardChilds.length; i++) {
			_obj = Ext.Object.merge(_obj, wizardChilds[i].getValues(false, false, false, true));
		}

		return _obj;
	},

	setValue: function(data) {
		var wizardChilds = this.contentPanel.items.items;

		for(var i = 0; i < wizardChilds.length; i++) {
			var form = wizardChilds[i].getForm();
			form.setValues(data);
		}
	}
});
