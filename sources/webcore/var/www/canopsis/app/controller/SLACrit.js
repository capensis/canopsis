//need:app/lib/controller/cgrid.js,app/view/Rule/Grid.js,app/view/Rule/Form.js,app/store/Rules.js
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
Ext.define('canopsis.controller.SLACrit', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['SLACrit.Grid', 'SLACrit.Form'],
	stores: ['SLACrit'],
	models: ['SLACrit'],

	logAuthor: '[controller][SLACrit]',

	storeDefaultAction: undefined,

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.formXtype = 'SLACritForm';
		this.listXtype = 'SLACritGrid';

		this.modelId = 'SLACrit';

		this.callParent(arguments);
	},

	bindGridEvents: function(grid){
		var slaMacrosCtrl = this.getController('SLAMacros');

		// Bind ldap button
		var btns = Ext.ComponentQuery.query('#' + grid.id + ' button[action=edit-macros-button]');

		for(var i = 0; i < btns.length; i++) {
			btns[i].on('click', slaMacrosCtrl.slaButton, slaMacrosCtrl);
		}
	},

	preSave: function(record) {
		var slatypes = ['warn', 'crit'];
		var slastates = ['ok', 'nok', 'out'];
		var crit = record.get('crit');

		var evt = {
			'connector': 'canopsis',
			'connector_name': 'sla',
			'event_type': 'perf',
			'source_type': 'component',
			'component': '__canopsis__',
			'state': 0,
			'perf_data_array': []
		};

		for(var i = 0; i < slatypes.length; i++) {
			for(var j = 0; j < slastates.length; j++) {
				var metric = 'cps_sla_' + slatypes[i] + '_' + crit + '_' + slastates[j];

				evt.perf_data_array.push({
					'metric': metric,
					'value': 0,
					'type': 'COUNTER'
				});
			}
		}

		log.debug('Send event', this.logAuthor);
		log.debug(evt);

		global.eventsCtrl.sendEvent(evt);

		return record;
	}
});