//need:app/lib/controller/cgrid.js,app/view/Derogation/Form.js,app/view/Derogation/Grid.js,app/store/Derogations.js
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
Ext.define('canopsis.controller.Derogation', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Derogation.Form', 'Derogation.Grid'],

	model: ['Derogation'],
	stores: ['Derogations'],

	logAuthor: '[controller][Derogation]',

	init: function() {
		log.debug('[' + this.id + '] - Initialize ...');

		this.formXtype = 'DerogationForm';
		this.listXtype = 'DerogationGrid';

		this.modelId = 'Derogation';

		this.callParent(arguments);

		global.derogationCtrl = this;
	},

	_preSave: function(record, data, form) {
		log.debug('Process record', this.logAuthor);

		var stopTs = undefined;

		if(data.forTs) {
			stopTs = data.startTs + data.forTs;
			record.set('forTs', data.forTs);
		}
		else {
			stopTs = data.stopTs;
		}

		record.set('stopTs', stopTs);
		record.set('startTs', data.startTs);

		record.set('time_conditions', [{type: 'time_interval', startTs: data.startTs, stopTs: stopTs}]);
		record.set('description', data.description);

		if(form.editing) {
			if(Ext.isArray(form.record.ids)) {
				record.set('ids', form.record.ids);
			}
			else {
				record.set('ids', [form.record.ids]);
			}
		}
		else {
			record.set('ids', form.ids);
		}

		if(data._id) {
			record.set('_id', data._id);
		}
		else {
			record.set('_id', global.gen_id());
		}

		if(data.tags && Ext.isString(data.tags)) {
			record.set('tags', [data.tags]);
		}

		if(form.selector_name) {
			record.set('selector_name', form.selector_name);
		}

		// Actions
		var actions = [];

		if(!Ext.isArray(data.actions)) {
			actions.push(data.actions);
		}
		else {
			actions = actions.concat(data.actions);
		}

		if(data.statemap) {
			actions.push(data.statemap);
		}

		record.set('actions', actions);

		// conditions
		if(data.evfilter) {
			record.set('conditions', data.evfilter);
		}

		return record;
	},

	getEditTitle: function() {
		var name = _(this.modelId);
		return name;
	},

	afterload_EditForm: function(form, item) {
		var data = item.data;

		if(!data.forTs) {
			form.periodTypeCombo.setValue('to');
		}

		var conditions = data.time_conditions[0];

		if(conditions.stopTs) {
			form.stopDate.setValue(conditions.stopTs);
		}

		if(conditions.startTs) {
			form.startDate.setValue(conditions.startTs);
		}

		if(data.actions) {
			for(var i = 0; i < data.actions.length; i++) {
				var action = data.actions[i];
				log.dump(action);

				if(action.type === 'override') {
					var field = action.field;
					var value = action.value;
					form.addNewField(field, value);
				}
				else if(action.type === 'requalificate') {
					form.setRequalification(action.statemap);
				}
			}
		}

		if(data.conditions) {
			form.eventFilter.setValue(data.conditions);
		}
	}
});
