//need:app/store/SLA.js,app/lib/view/cform.js,app/lib/form/field/cdate.js,app/lib/form/field/cduration.js
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

var sla_table_template = Ext.create('Ext.XTemplate',
	'<table class="sla-table">',
		'<thead>',
			'<tr>',
				'<th></th>',
				'<th></th>',
				'<th>SLA OK</th>',
				'<th>SLA NOK</th>',
				'<th>NOT ACK</th>',
				'<th>SLA OK Rate</th>',
			'</tr>',
		'</thead>',
		'<tbody>',
			'<tpl for="crits">',
				'<tr>',
					'<td rowspan="2">{name}</td>',
					'<td>Warning</td>',
					'<td>{nb_warn_sla_ok}</td>',
					'<td>{nb_warn_sla_nok}</td>',
					'<td>{nb_warn_not_ack}</td>',
					'<td>{pct_warn_sla_ok}</td>',
				'</tr>',
				'<tr>',
					'<td>Critical</td>',
					'<td>{nb_crit_sla_ok}</td>',
					'<td>{nb_crit_sla_nok}</td>',
					'<td>{nb_crit_not_ack}</td>',
					'<td>{pct_crit_sla_ok}</td>',
				'</tr>',
			'</tpl>',
		'</tbody>',
	'</table>'
);

Ext.define('canopsis.view.SLA.View', {
	extend: 'Ext.panel.Panel',
	alias: 'widget.SLAView',
	layout: 'fit',

	requires: [
		'canopsis.store.SLA',
		'canopsis.model.SLA',
		'canopsis.lib.view.cform',
		'canopsis.lib.form.field.cdate',
		'canopsis.lib.form.field.cduration'
	],

	logAuthor: '[SLA]',

	initComponent: function() {
		this.callParent(arguments);

		this.macro = null;
		this.period = null;

		/* load configuration */
		this.store_crit = Ext.create('canopsis.store.SLA', {
			filters: [{
				property: 'objclass',
				value: 'crit'
			}]
		});

		this.store_crit.load(this.buildView.bind(this));

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

	buildView: function() {
		this.graphs = Ext.create('Ext.container.Container', {
			layout: 'vbox',
		});

		this.store_crit.each(function(record) {
			var graph = Ext.create('Ext.container.Container', {
				layout: 'hbox',

				items: [{
					xtype: 'label',
					text: record.get('crit') + ' acknowledgement'
				},{
					xtype: 'button',
					text: 'Test'
				},{
					xtype: 'button',
					text: 'Test2'
				}]
			});

			this.graphs.add(graph);
		}, this);

		this.add(this.graphs);
	},

	buildConfigCrit: function() {
		this.critGridEdit = Ext.create('Ext.grid.plugin.RowEditing', {
			autoCancel: false
		});

		this.critGrid = Ext.widget({
			xtype: 'grid',
			title: 'Criticality',
			multiSelect: true,

			tbar: [{
				xtype: 'button',
				text: _('Add'),
				iconCls: 'icon-add',
				action: 'add'
			},{
				xtype: 'button',
				text: _('Reload'),
				iconCls: 'icon-reload',
				action: 'reload'
			},{
				xtype: 'button',
				text: _('Delete'),
				iconCls: 'icon-delete',
				action: 'delete'
			},{
				xtype: 'button',
				text: _('Save'),
				iconCls: 'icon-save',
				action: 'sync'
			}],

			columns: [{
				text: 'Criticality',
				dataIndex: 'crit',
				flex: 1,
				editor: {
					allowBlank: false
				}
			},{
				text: 'Delay',
				dataIndex: 'delay',
				renderer: rdr_time_interval,
				flex: 1,
				editor: {
					xtype: 'cduration',
					allowBlank: false
				}
			}],

			plugins: [this.critGridEdit],
			store: this.store_crit
		});

		/* bind actions */
		var actions = {
			add: this.addCrit,
			reload: this.reloadCrit,
			delete: this.removeCrit,
			sync: this.saveCrit
		};

		for(var name in actions) {
			var btns = Ext.ComponentQuery.query('#' + this.critGrid.id + ' [action=' + name + ']');

			for(var i = 0; i < btns.length; i++) {
				btns[i].on('click', actions[name], this);
			}
		}
	},

	buildConfigMacro: function() {
		this.macroForm = Ext.widget({
			title: 'Macro',
			xtype: 'cform',

			items: [{
				xtype: 'textfield',
				name: 'warn',
				fieldLabel: 'Warning',
				value: this.macro.get('mWarn')
			},{
				xtype: 'textfield',
				name: 'crit',
				fieldLabel: 'Critical',
				value: this.macro.get('mCrit')
			}]
		});

		/* bind form actions */
		var actions = {
			save: this.saveMacro,
			cancel: this.restoreMacro
		};

		for(var name in actions) {
			var btns = Ext.ComponentQuery.query('#' + this.macroForm.id + ' [action=' + name + ']');

			for(var i = 0; i < btns.length; i++) {
				btns[i].on('click', actions[name], this);
			}
		}
	},

	buildConfigPeriod: function() {
		this.periodForm = Ext.widget({
			title: 'Period',
			xtype: 'cform',

			items: [{
				xtype: 'cdate',
				name: 'from',
				label_text: 'From',
				value: this.period.get('from')
			},{
				xtype: 'cdate',
				name: 'to',
				label_text: 'To',
				value: this.period.get('to')
			}]
		});

		/* bind form actions */
		var actions = {
			save: this.savePeriod,
			cancel: this.restorePeriod
		};

		for(var name in actions) {
			var btns = Ext.ComponentQuery.query('#' + this.periodForm.id + ' [action=' + name + ']');

			for(var i = 0; i < btns.length; i++) {
				btns[i].on('click', actions[name], this);
			}
		}
	},

	buildConfigWindow: function() {
		this.buildConfigCrit();
		this.buildConfigMacro();
		this.buildConfigPeriod();

		/* create configuration window */

		this.winconfig = Ext.create('Ext.window.Window', {
			title: 'Configure SLA',
			closeAction: 'hide',

			layout: 'fit',
			width: 350,
			height: 400,

			items: {
				xtype: 'tabpanel',
				items: [
					this.critGrid,
					this.macroForm,
					this.periodForm
				]
			}
		});
	},

	configureSla: function() {
		this.winconfig.show();
	},

	addCrit: function() {
		log.debug('Adding crit', this.logAuthor);

		this.critGridEdit.cancelEdit();

		var crit = Ext.create('canopsis.model.SLA', {
			objclass: 'crit'
		});

		this.store_crit.insert(0, crit);
		this.critGridEdit.startEdit(0, 0);
	},

	reloadCrit: function() {
		log.debug('Reloading crit store', this.logAuthor);
		this.store_crit.load();
	},

	removeCrit: function() {
		log.debug('Removing selected crits', this.logAuthor);

		var sm = this.critGrid.getSelectionModel();

		this.critGridEdit.cancelEdit();

		this.store_crit.remove(sm.getSelection());

		if(this.store_crit.getCount() > 0) {
			sm.select(0);
		}

		this.store_crit.sync();
	},

	saveCrit: function() {
		log.debug('Synchronize crit store', this.logAuthor);
		this.store_crit.sync();
	},

	saveMacro: function() {
		var record = this.macroForm.getValues();

		this.macro = this.store.findRecord('objclass', 'macro');
		this.macro.set('mWarn', record.warn);
		this.macro.set('mCrit', record.crit);

		this.store.sync();
	},

	restoreMacro: function() {
		this.macroForm.getForm().setValues({
			warn: this.macro.get('mWarn'),
			crit: this.macro.get('mCrit')
		});
	},

	savePeriod: function() {
		var record = this.periodForm.getValues();

		this.period = this.store.findRecord('objclass', 'period');
		this.period.set('from', record.from);
		this.period.set('to', record.to);

		this.store.sync();
	},

	restorePeriod: function() {
		this.periodForm.getForm().setValues({
			from: this.period.get('from'),
			to: this.period.get('to')
		});
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