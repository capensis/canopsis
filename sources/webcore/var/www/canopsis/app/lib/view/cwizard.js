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
	extend: 'Ext.window.Window',

	alias: 'widget.ViewBuilderWizard',

	requires: [
				'canopsis.lib.form.field.cinventory',
				'canopsis.lib.form.field.cmetric',
				'canopsis.lib.form.field.cfilter',
				'canopsis.lib.form.field.ctag',
				'canopsis.lib.form.field.cfieldset',
				'canopsis.lib.form.field.cdate',
				'canopsis.lib.form.field.cduration',
				'canopsis.lib.form.field.cduration',
				'canopsis.lib.form.field.ccolorfield'
			],

	title: _('Wizard'),
	closable: true,
	closeAction: 'destroy',
	width: 800,
	//minWidth: 350,
	height: 500,
	layout: 'fit',
	bodyStyle: 'padding: 5px;',
	labelWidth: 200,

	constrain: true,
	constrainHeader: true,

	edit: false,

	step_list: [{
			title: _("i'm empty !"),
			html: _('you must give an object to fill me')
		}],

	initComponent: function() {
		this.logAuthor = '[Wizard ' + this.id + ']';
		log.debug('Create Wizard "' + this.title + '"' , this.logAuthor);


		//-----------------buttons--------------------------

		this.bbar = Ext.create('Ext.toolbar.Toolbar');

		this.cancelButton = this.bbar.add({xtype: 'button', text: _('Cancel'), action: 'cancel', iconCls: 'icon-cancel'});
		this.bbar.add('->');
		this.previousButton = this.bbar.add({xtype: 'button', text: _('Previous'), action: 'previous', disabled: true, iconCls: 'icon-previous'});
		this.nextButton = this.bbar.add({xtype: 'button', text: _('Next'), action: 'next', disabled: true, iconCls: 'icon-next', iconAlign: 'right'});

		this.finishButton = this.bbar.add({xtype: 'button', text: _('Finish'), action: 'finish', iconCls: 'icon-save', iconAlign: 'right'});

		//--
		this.tabPanel = Ext.create('Ext.tab.Panel', {
			layout: 'fit',
			//xtype: 'tabpanel',
			plain: true
			//deferredRender: false,
		});


		if (this.step_list) {
			//var tmp = this.build_step_list(this.step_list)
			log.debug('Wizard steps fully generated', this.logAuthor);
			for (var i = 0; i < this.step_list.length; i++)
				this.add_new_step(this.step_list[i]);
		}

		this.items = [this.tabPanel];

		this.callParent(arguments);

		this.previousButton.setDisabled(true);

	},

	afterRender: function() {
		//needed
		this.callParent(arguments);
		this.tabPanel.setActiveTab(0);
		this.bind_buttons();

		//bind combobox
		if (this.data) {
			this.loadData();
		} else {
			var combo = Ext.ComponentQuery.query('#' + this.id + ' [name=xtype]');
			if (combo.length != 0)
				combo[0].on('select', this.add_option_panel, this);
		}
	},


	bind_buttons: function() {
		log.debug('binding buttons', this.logAuthor);
		//---------------------previous button--------------------
		var btns = Ext.ComponentQuery.query('#' + this.id + ' [action=previous]');
		for (var i = 0; i < btns.length; i++)
			btns[i].on('click', this.previous_button, this);

		//---------------------next button--------------------
		var btns = Ext.ComponentQuery.query('#' + this.id + ' [action=next]');
		for (var i = 0; i < btns.length; i++)
			btns[i].on('click', this.next_button, this);

		//---------------------cancel button--------------------
		var btns = Ext.ComponentQuery.query('#' + this.id + ' [action=cancel]');
		for (var i = 0; i < btns.length; i++)
			btns[i].on('click', this.cancel_button, this);

		//---------------------finish button-------------------
		var btns = Ext.ComponentQuery.query('#' + this.id + ' [action=finish]');
		for (var i = 0; i < btns.length; i++)
			btns[i].on('click', this.finish_button, this);

		this.tabPanel.on('tabchange', this.update_button, this);
	},

	add_new_step: function(step) {
		step.autoScroll = true;
		step.xtype = 'form';

		step.defaults = { labelWidth: this.labelWidth };

		if (step.items.length == 1)
			if (! step.layout && !(step.items[0].xtype == 'fieldset' || step.items[0].xtype == 'cfieldset')) {
				step.layout = 'fit';
				step.defaults.autoScroll = true;
			}

		step.padding = 5;

		return this.tabPanel.add(step);
	},

	loadData: function() {
		this.edit = true;
		if (this.data.xtype) {
			var combo = Ext.ComponentQuery.query('#' + this.id + ' [name=xtype]');

			combo[0].setValue(this.data.xtype);
			var list_tab = this.add_option_panel();
			combo[0].setDisabled(true);
		}

		//hack for delete description
		if (this.data.description)
			delete this.data.description;

		var child_items = this.tabPanel.items.items;
		for (var i = 0; i < child_items.length; i++) {
			var form = child_items[i].getForm();
			form.setValues(this.data);
		}

	},

	reset_steps: function() {
		var tab_childs = this.tabPanel.items.items;
		var tab_length = tab_childs.length;

		//log.debug('child panel : ' + tab_length)
		//log.debug('step list length :' + this.step_list.length)

		var tab_to_remove = [];

		for (var i = this.step_list.length; i < tab_length; i++)
			tab_to_remove.push(tab_childs[i]);

		for (var i = 0; i < tab_to_remove.length; i++)
			this.tabPanel.remove(tab_to_remove[i]);

	},

	isValid: function() {
		var child_items = this.tabPanel.items.items;
		for (var i = 0; i < child_items.length; i++) {
			var form = child_items[i].getForm();
			if (form.hasInvalidField())
				return false;
		}
		return true;
	},

	get_variables: function() {
		var output = {};

		var child_items = this.tabPanel.items.items;
		for (var i = 0; i < child_items.length; i++) {
			var form = child_items[i].getForm();
			var values = form.getValues(false, false, false, true);

			Ext.Object.each(values, function(key, value, myself) {
				output[key] = value;
			});
		}
		return output;
	},

	add_option_panel: function() {
		output = [];
		this.reset_steps();
		var combo = Ext.ComponentQuery.query('#' + this.id + ' [name=xtype]');
		var description_field = Ext.ComponentQuery.query('#' + this.id + ' [name=description]')[0];

		if (combo[0].isValid()) {
			var store = combo[0].getStore();
			var record = store.findRecord('xtype', combo[0].getValue());

			this.set_default_values(record);

			var options = record.get('options');

			if (options) {
				for (var i = 0; i < options.length; i++) {
					for (var j = 0; j < options[i].items.length; j++)
						if (options[i].items[j].xtype == 'fieldset' || options[i].items[j].xtype == 'cfieldset')
							options[i].items[j].defaults = { labelWidth: this.labelWidth };

					output.push(this.add_new_step(options[i]));
				}
			}

			var description = undefined;
			if (global.locale != 'en')
				var description = record.get('description' + '-' + global.locale);
			if (!description)
				var description = record.get('description');

			if (description_field && description)
				description_field.setValue(description);

			this.update_button();
		}
		return output;
	},

	set_default_values: function(record) {
		log.debug('set defaults value to wizard', this.logAuthor);
		var elements = Ext.ComponentQuery.query('#' + this.id + ' [name]');
		for (var i = 0; i < elements.length; i++) {
			var element = elements[i];
			//don't reset xtype
			if (element.name != 'xtype') {
				var value = record.get(element.name);
				if (value != undefined) {
					//log.debug('setting value :' + value)
					//log.debug('to : ' + element.name)
					element.setValue(value);
				}
			}
		}
	},
	//----------------------button action functions-----------------------
	previous_button: function() {
		log.debug('previous button', this.logAuthor);
		panel = this.tabPanel;
		active_tab = this.tabPanel.getActiveTab();
		panel.setActiveTab(panel.items.indexOf(active_tab) - 1);
		this.update_button();
	},

	next_button: function() {
		log.debug('next button', this.logAuthor);

		var panel = this.tabPanel;
		var active_tab = panel.getActiveTab();
		var index = panel.items.indexOf(active_tab);

		panel.setActiveTab(index + 1);
		this.update_button();
	},

	update_button: function() {
		var activeTabIndex = this.tabPanel.items.findIndex('id', this.tabPanel.getActiveTab().id);
		var tabCount = this.tabPanel.items.length;

		if (activeTabIndex == 0) {
			this.previousButton.setDisabled(true);
		} else {
			this.previousButton.setDisabled(false);
		}

		if (activeTabIndex == (tabCount - 1)) {
			this.nextButton.setDisabled(true);
			if (!this.edit) {
				this.finishButton.setDisabled(false);
			}
		} else {
			this.nextButton.setDisabled(false);
			if (!this.edit) {
				this.finishButton.setDisabled(true);
			}
		}
	},

	cancel_button: function() {
		log.debug('cancel button', this.logAuthor);
		this.close();
	},

	finish_button: function() {
		log.debug('save button', this.logAuthor);
		if (this.isValid()) {
			var variables = this.get_variables();
			log.debug('Saved values are:', this.logAuthor);
			log.debug(variables, this.logAuthor);
			this.fireEvent('save', variables);
			this.close();
		}else {
			global.notify.notify('Form error', 'There is incorrect field in form', 'info');
		}
	}

});
