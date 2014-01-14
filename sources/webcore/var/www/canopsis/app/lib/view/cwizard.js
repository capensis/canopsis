//need:app/lib/form/field/cinventory.js,app/lib/form/field/ccomponentlist.js,app/lib/form/field/cmetric.js,app/lib/form/field/ccustom.js,app/lib/form/field/cfilter.js,app/lib/form/field/cwlist.js,app/lib/form/field/ctag.js,app/lib/form/field/cfieldset.js,app/lib/form/field/cdate.js,app/lib/form/field/cduration.js,app/lib/form/field/cduration.js,app/lib/form/field/ccolorfield.js,app/lib/view/ccard.js
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
Ext.define('canopsis.lib.view.cwizard', {
	extend: 'Ext.window.Window',
	alias: 'widget.cwizard',

	layout: 'fit',

	title: 'Wizard',
	closable: true,
	closeAction: 'destroy',
	constrain: true,
	width: 800,
	height: 500,

	cwlist: undefined,

	bbar: [
		{iconCls: 'icon-previous', text:'previous', _name: 'bbarPrevious'},
		{iconCls: 'icon-advanced', text: _('Advanced mode'), _name: 'bbarAdvance', enableToggle: true }, '->',
		{iconCls: 'icon-next', iconAlign: 'right', text: 'Next', _name: 'bbarNext'},
		{iconCls: 'icon-save', iconAlign: 'right',  text: 'Finish', _name: 'bbarFinish', hidden: true, disabled: true}
	],

	requires: [
		'canopsis.lib.form.field.cinventory',
		'canopsis.lib.form.field.ccomponentlist',
		'canopsis.lib.form.field.cmetric',
		'canopsis.lib.form.field.ccustom',
		'canopsis.lib.form.field.cfilter',
		'canopsis.lib.form.field.cwlist',
		'canopsis.lib.form.field.ctag',
		'canopsis.lib.form.field.cfieldset',
		'canopsis.lib.form.field.cdate',
		'canopsis.lib.form.field.cduration',
		'canopsis.lib.form.field.ccolorfield',
		'canopsis.lib.view.ccard'
	],

	items: [{
		xtype: 'ccard',
		wizardSteps: [{
			title: _('Choose widget'),
			items: [{
				xtype: 'cwlist',
				name: 'xtype'
			}]
		},{
			title: _('General'),
			items: [{
				xtype:'cfieldset',
				title:_('General options'),
				items: [{
					xtype: 'textfield',
					fieldLabel: _('Title') + ' (' + _('optional') + ')',
					name: 'title'
				},{
					xtype: 'checkbox',
					fieldLabel: _('Auto title') + ' ' + _('if available'),
					checked: true,
					inputValue: true,
					uncheckedValue: false,
					name: 'autoTitle'
				},{
					xtype: 'checkbox',
					fieldLabel: _('Show border'),
					checked: false,
					name: 'border',
					uncheckedValue: false
				},{
					xtype: 'checkbox',
					fieldLabel: _('Human readable values'),
					checked: true,
					inputValue: true,
					uncheckedValue: false,
					name: 'humanReadable'
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
							{value: 0,      text: _('None')},
							{value: 1,      text: '1 '  + _('second')},
							{value: 10,     text: '10'  + _('seconds')},
							{value: 30,     text: '30'  + _('seconds')},
							{value: 60,     text: '1 '  + _('minute')},
							{value: 300,    text: '5 '  + _('minutes')},
							{value: 600,    text: '10 ' + _('minutes')},
							{value: 900,    text: '15 ' + _('minutes')},
							{value: 1800,   text: '30 ' + _('minutes')},
							{value: 3600,   text: '1 '  + _('hour')}
						]
					}
				},{
					xtype: "cduration",
					name: "time_window",
					value: 86400,
					fieldLabel: _('Time window')
				}]
			}]
		}]
	}],

	initComponent: function() {
		this.callParent(arguments);

		this.ccard = this.down('ccard');
		this.bbarNextButton = this.down('button[_name=bbarNext]');
		this.bbarFinishButton = this.down('button[_name=bbarFinish]');
		this.bbarPreviousButton = this.down('button[_name="bbarPrevious"]');
		this.bbarAdvanceButton = this.down('button[_name="bbarAdvance"]');

		//binding events
		if(this.bbarFinishButton) {
			this.bbarFinishButton.on('click', function() {
				this.fireEvent('save',this.widgetId,this.getValue());
				this.destroy();
			}, this);
		}

		if(this.bbarPreviousButton) {
			this.bbarPreviousButton.on('click', function() {
				this.previousButton();
			}, this.ccard);
		}

		if(this.bbarNextButton) {
			this.bbarNextButton.on('click', function() {
				this.nextButton();
			}, this.ccard);
		}

		if(this.bbarAdvanceButton) {
			this.bbarAdvanceButton.on('toggle', function(button,toggle){
				void(button);

				if(toggle) {
					this.switchToAdvancedMode();
				}
				else {
					this.switchToSimpleMode();
				}
			}, this.ccard);
		}

		this.ccard.on('buttonChanged', this.checkDisplayFinishButton, this);
	},

	afterRender: function() {
		this.callParent(arguments);

		this.childStores = {};

		if(this.data){
			this.edit = true;
			this.bbarNextButton.hide();
			this.bbarNextButton.hide();
			this.bbarPreviousButton.hide();
			this.bbarFinishButton.show();
		}

		this.cwlist = this.down('cwlist');

		this.cwlist.on('select', this.addOptionPanel, this);

		if(this.data) {
			this.cwlist.setValue(this.data.xtype);
			this.cwlist.setDisabled(true);

			var options = Ext.clone(this.cwlist.nodes[0].raw.options);
			this.ccard.addNewSteps(options);
			this.ccard.setValue(this.data);

			// Hide "choose widgets" in edit mode
			this.ccard.getButton(0).hide();
			this.ccard.showStep(1);
			this.bbarFinishButton.enable();
		}
	},

	addOptionPanel: function(cwlist, records) {
		void(cwlist);

		this.ccard.cleanPanels();

		this.bbarFinishButton.enable();

		this.ccard.addNewSteps(Ext.clone(records[0].raw.options));

		//additionnal options
		var options = {
			"refreshInterval": records[0].raw.refreshInterval,
			"border" : records[0].raw.border
		};

		this.ccard.setValue(options);
	},

	beforeGetValue: function() {
		this.cwlist.setDisabled(false);
	},

	checkDisplayFinishButton: function() {
		if(this.edit) {
			return;
		}

		if(this.ccard.isLastPanel()) {
			this.bbarNextButton.hide();
			this.bbarFinishButton.show();
		}
		else {
			this.bbarNextButton.show();
			this.bbarFinishButton.hide();
		}
	},

	getValue: function() {
		if(this.beforeGetValue) {
			this.beforeGetValue();
		}

		var wizardChilds = this.ccard.contentPanel.items.items;

		var mergedNode = {};
		var cmetricName = undefined;
		var outputObj = {};

		//for each panel
		for(var i = 0; i < wizardChilds.length; i++) {
			var form = wizardChilds[i];
			var values = form.getValues(false, false, false, true);

			var formType = form.items.items[0].xtype;

			if(formType === 'cmetric'){
				//the key is the cmetric given name (in widget.json)
				cmetricName =  Ext.Object.getKeys(values)[0];

				for(var j = 0; j < values[cmetricName].length; j++) {
					var node = values[cmetricName][j];
					mergedNode[node.id] = node;
				}
			}
			else {
				outputObj = Ext.Object.merge(outputObj, values);
			}
		}

		if(mergedNode.length !== 0) {
			outputObj[cmetricName] = mergedNode;
		}

		return outputObj;
	}
});
