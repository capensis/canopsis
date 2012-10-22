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
Ext.define('canopsis.controller.Selector', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Selector.Grid', 'Selector.Form'],
	stores: ['Selectors'],
	models: ['Selector'],

	logAuthor: '[controller][selector]',

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.formXtype = 'SelectorForm';
		this.listXtype = 'SelectorGrid';

		this.modelId = 'Selector';

		this.callParent(arguments);

		//needed for weather widget
		global.selectorCtrl = this;
	},

	beforeload_EditForm: function(form, item) {
		var name = Ext.ComponentQuery.query('#' + form.id + ' textfield[name=crecord_name]')[0];
		if (name)
			name.setReadOnly(true);
	},

	preSave: function(record, data, form) {
		var _id = record.get('_id');
		record.set('id', _id);
		record.set('loaded', false);

		if (record.get('dosla'))
			record.set('sla_timewindow', record.get('sla_timewindow_value') * record.get('sla_timewindow_unit'));
		else
			record.set('sla_timewindow', undefined);

		record.set('state', undefined);
		record.set('sla_state', undefined);
		record.set('sla_timewindow_perfdata', undefined);

		return record;
	},

	ajaxValidation: function(record,edit) {
		if (edit) {
			this._save(record, true);
			return;
		}

		this.backup_record = record;
		filter = {
					filter: Ext.encode({'crecord_name': record.get('crecord_name')}),
					limit: 1
					};

		Ext.Ajax.request({
			method: 'GET',
			scope: this,
			params: filter,
			url: '/rest/object/selector',
			success: function(response, opts) {
				var data = Ext.decode(response.responseText).data;

				if (data.length == 0) {
					this._save(this.backup_record, false);
					this.backup_record = undefined;
					this.backup_edit = undefined;
				}else {
					global.notify.notify(_('Bad name'), _('This selector name already exist'), 'warning');
				}
			}
		});
	},

	change_selector_output: function(_id,type,message) {
		log.debug('Change selector/sla output', this.logAuthor);
		log.debug('_id: ' + _id, this.logAuthor);
		log.debug('message: ' + message, this.logAuthor);

		var data = {
			loaded: false
		};

		if (type == 'selector')
			data.output_tpl = message;
		else
			data.sla_output_tpl = message;

		updateRecord('object', 'selector', 'canopsis.model.selector', _id, data,
			function() {
				global.notify.notify(_('Message updated'), 'The message will be display in few minutes', 'success');
			});
	}

});
