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
Ext.define('canopsis.lib.view.cform', {
	extend: 'Ext.form.Panel',

	alias: 'widget.cform',

	requires: ['Ext.form.field.Text'],

	title: '',
	bodyStyle: 'padding:5px 5px 0',
	border: 0,

	EditMethod: 'tab',

	defaultType: 'textfield',

	logAuthor: '[view][cform]',

	initComponent: function() {
		var tbar = [
			{
				iconCls: 'icon-save',
				text: _('Save'),
				action: 'save'
			},{
				iconCls: 'icon-cancel',
				text: _('Cancel'),
				action: 'cancel'
			}
		];

		var bbar = [
			{
				iconCls: 'icon-cancel',
				text: _('Cancel'),
				action: 'cancel'
			},{
				xtype: 'tbfill'
			},{
				iconCls: 'icon-save',
				text: _('Save'),
				action: 'save',
				pack: 'end'
			}
		];

		if(this.EditMethod === 'tab') {
			this.on('beforeclose', this.beforeclose);
			this.tbar = tbar;
		}
		else {
			this.bbar = bbar;
		}

		this.callParent();
	},

	beforeclose: function() {
		log.debug('Active previous tab', this.logAuthor);
		old_tab = Ext.getCmp('main-tabs').old_tab;

		if(old_tab) {
			Ext.getCmp('main-tabs').setActiveTab(old_tab);
		}
	},

	beforeDestroy: function() {
		log.debug('Destroy items ...', this.logAuthor);
		Ext.form.Panel.superclass.beforeDestroy.call(this);
		log.debug(this.id + ' Destroyed.', this.logAuthor);
	}
});
