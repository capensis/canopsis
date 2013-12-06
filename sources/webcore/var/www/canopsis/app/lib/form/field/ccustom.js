//need:app/lib/view/ccard.js,app/lib/form/cfield.js
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

Ext.define('canopsis.lib.form.field.ccustom', {
	extend: 'canopsis.lib.view.ccard',
	mixins: ['canopsis.lib.form.cfield'],

	alias: 'widget.ccustom',

	border: false,

	sharedStore: undefined,

	name: 'ccustom',

	contentPanelDefault: {
		xtype: 'form',
		width: '100%',
		border: false,
		margin: '30px',
		autoScroll: true,
		layout: {
			type: 'vbox',
			align: 'stretch'
		},
		defaults: {
			listeners: {
				change: function() {
					this.up('form').fireEvent('ccustomChange');
				}
			}
		}
	},

	buttonTpl: new Ext.Template([
		'<div class="ccustomButton">',
			'{title}',
		'</div>'
	]).compile(),

	initComponent: function() {
		this.callParent(arguments);
		this.panelIdByNode = {};
	},

	afterRender: function() {
		this.callParent(arguments);

		this.sourceStore = this.findParentByType('cwizard').childStores[this.sharedStore];

		this.sourceStore.on('add', function(store, records) {
			void(store);

			this.addPanels(records);
		}, this);

		this.sourceStore.on('remove', function(store, records) {
			void(store);

			this.removePanels(records);
		}, this);

		this.sourceStore.each(function(record) {
			this.addPanels(record);
		}, this);
	},

	addPanels: function(records) {
		if(!Ext.isArray(records)) {
			records = [records];
		}

		for(var i = 0; i < records.length; i++) {
			var nodeId = records[i].data.id;

			if(!this.panelIdByNode[nodeId]) {
				this.addPanel(nodeId, records[i]);
			}
		}
	},

	addPanel: function(nodeId, record) {
		void(nodeId);

		var panelConfig = {
			items: Ext.clone(this.customForm)
		};

		var title = this.buildTitle(record.data);
		var panelNumber = this.addContent(panelConfig, record);

		this.panelIdByNode[record.get('id')] = panelNumber;

		var button = {
			title: title,
			panelIndex: panelNumber
		};

		this.addButton(button);
	},

	addContent: function(extjs_obj_array, record) {
		var _obj = Ext.Object.merge(extjs_obj_array, this.contentPanelDefault);

		//add title in front of panel
		var panelTitle = {
			xtype: 'panel',
			border: false,
			padding: '0 0 20 0',
			html: '<center>' + this.buildTitle(record.data) + '</center>'
		};

		_obj.items = Ext.Array.merge(panelTitle,_obj.items);

		var panel = this.contentPanel.add(_obj);

		if(record) {
			panel.getForm().setValues(record.data);
			panel.sourceRecord = record;
		}

		//save record when change made on form
		panel.on('ccustomChange', function() {
			var values = this.getForm().getValues(false, false, false, true);

			if(this.sourceRecord) {
				this.sourceRecord.beginEdit();
				this.sourceRecord.set(values);
				this.sourceRecord.endEdit(true);
			}
		}, panel);

		this.panelByindex[this.contentPanel.items.length - 1] = panel;
		return this.contentPanel.items.length - 1;
	},

	removePanels: function(records) {
		//if someone find something to do this in extjs, build this again
		if(!Ext.isArray(records)) {
			records = [records];
		}

		for(var i = 0; i < records.length; i++) {
			var nodeId = records[i].data.id;

			if(this.panelIdByNode[nodeId] !== undefined) {
				var panelNumber = this.panelIdByNode[nodeId];
				var panel = this.panelByindex[panelNumber];
				var button = this.buttonByPanelIndex[panelNumber];

				if(panel && button) {
					panel.setDisabled(true);
					panel.hide();

					button.setDisabled(true);
					button.hide();

					delete this.panelByindex[panelNumber];
					delete this.panelIdByNode[nodeId];
					delete this.buttonByPanelIndex[panelNumber];
				}
			}
		}
	},

	setValue: function() {
		return;
	},

	getValue: function() {
		return;
	},

	buildTitle: function(data) {
		if(data.titleInWizard) {
			return data.titleInWizard;
		}

		if(!data.co || !data.me) {
			return data._id;
		}

		var title = data.co;

		if(data.re) {
			title += ' ' + data.re;
		}

		title += ' ' + data.me;

		return title;
	}
});
