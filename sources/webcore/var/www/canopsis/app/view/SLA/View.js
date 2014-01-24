//need:app/store/SLA.js;app/lib/view/cform.js;app/lib/form/field/cdate.js
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
Ext.define('canopsis.view.SLA.View', {
	extend: 'Ext.panel.Panel',
	alias: 'widget.SLAView',

	requires: [
		'canopsis.store.SLA',
		'canopsis.model.SLA',
		'canopsis.lib.view.cform',
		'canopsis.lib.form.field.cdate'
	],

	logAuthor: '[SLA]',

	initComponent: function() {
		this.callParent(arguments);

		this.macro = null;
		this.period = null;

		/* load configuration */
		this.store = Ext.create('canopsis.store.SLA', {});
		this.store.load(function() {
			/* get macro configuration */
			this.macro = this.store.findRecord('objclass', 'macro');

			if(this.macro === null) {
				this.macro = Ext.create('canopsis.model.SLA', {
					objclass: 'macro'
				});

				this.store.add(this.macro);
			}

			/* get period configuration */
			this.period = this.store.findRecord('objclass', 'period');

			if(this.period === null) {
				var to = new Date() / 1000;
				var from = new Date(to * 1000 - global.commonTs.week * 1000) / 1000;

				this.period = Ext.create('canopsis.model.SLA', {
					objclass: 'period',
					from: from,
					to: to
				});

				this.store.add(this.period);
			}

			log.debug('Configuration loaded:', this.logAuthor);
			log.dump(this.macro);
			log.dump(this.period);

			/* synchronize the store */
			this.store.sync();

			/* build UI */
			this.buildConfigWindow();
		}.bind(this));

		this.buildToolbar();
	},

	buildToolbar: function() {
		this.slaSelector = Ext.widget({
			xtype: 'combobox',
			displayField: 'title',
			valueField: 'type',
			queryMode: 'local',
			value: 'periods',
			store: {
				xtype: 'store',
				fields: ['title', 'type'],
				data: [{
					title: 'Periods',
					type: 'periods',
					handler: this.slaPerPeriods.bind(this)
				},{
					title: 'Hostgroups',
					type: 'hostgroups',
					handler: this.slaPerHostgroups.bind(this)
				}]
			},
			listeners: {
				select: this.selectSlaType.bind(this)
			}
		});

		this.addDocked({
			xtype: 'toolbar',
			dock: 'top',
			items: [{
				xtype: 'label',
				text: 'SLA per:'
			},
			this.slaSelector,
			{
				xtype: 'tbfill'
			},{
				xtype: 'button',
				text: 'Configure',
				iconCls: 'icon-mainbar icon-preferences',
				handler: this.configureSla.bind(this)
			}]
		});
	},

	buildConfigWindow: function() {
		this.critGrid = Ext.widget({
			xtype: 'grid',
			border: false,
			title: 'Criticality',

			tbar: [{
				xtype: 'button',
				text: 'Add',
				iconCls: 'icon-add',
				handler: this.addCriticality.bind(this)
			},{
				xtype: 'button',
				text: 'Delete',
				iconCls: 'icon-delete',
				handler: this.removeCriticality.bind(this)
			}],

			columns: [{
				text: 'Criticality',
				dataIndex: 'crit',
				flex: 1
			},{
				text: 'Delay',
				dataIndex: 'delay',
				renderer: rdr_time_interval,
				flex: 1
			}],

			store: Ext.create('canopsis.store.SLA', {
				filters: [{
					property: 'objclass',
					value: 'crit'
				}]
			})
		});

		this.winconfig = Ext.create('Ext.window.Window', {
			title: 'Configure SLA',
			closeAction: 'hide',

			layout: 'fit',
			width: 350,
			height: 400,

			items: {
				xtype: 'tabpanel',
				items: [this.critGrid, {
					title: 'Macro',
					xtype: 'cform',

					items: [{
						xtype: 'textfield',
						fieldLabel: 'Warning',
						value: this.macro.get('mWarn')
					},{
						xtype: 'textfield',
						fieldLabel: 'Critical',
						value: this.macro.get('mCrit')
					}]
				},{
					title: 'Period',
					xtype: 'cform',

					items: [{
						xtype: 'cdate',
						label_text: 'From',
						value: this.period.get('from')
					},{
						xtype: 'cdate',
						label_text: 'To',
						value: this.period.get('to')
					}]
				}]
			}
		});
	},

	configureSla: function() {
		this.winconfig.show();
	},

	addCriticality: function() {
		console.log('Add');
	},

	removeCriticality: function()  {
		console.log('Remove');
	},

	selectSlaType: function(combo, records) {
		void(combo);

		var selected = records[0]; /* no multi-select so records.length always equal to 1 */

		/* call the handler associated to the value */
		selected.raw.handler.call(this);
	},

	slaPerPeriods: function() {
		console.log('Periods');
	},

	slaPerHostgroups: function() {
		console.log('Hostgroups');
	},

	beforeDestroy: function() {
		this.winconfig.destroy();
	}
});