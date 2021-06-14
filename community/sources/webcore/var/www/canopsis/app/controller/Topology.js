//need:app/lib/controller/cgrid.js,app/view/Topology/Grid.js,app/view/Topology/Form.js,app/store/Topologies.js
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
Ext.define('canopsis.controller.Topology', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Topology.Grid', 'Topology.Form'],
	stores: ['Topologies'],

	models: ['Topology'],

	logAuthor: '[controller][topology]',

	EditMethod: 'tab',

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.formXtype = 'TopologyForm';
		this.listXtype = 'TopologyGrid';

		this.modelId = 'Topology';

		this.callParent(arguments);
	},

	preSave: function(record) {
		record.set('loaded', false);
		return record;
	},

	ajaxValidation: function(record, edit) {
		if(edit) {
			this._save(record, true);
			return;
		}

		isRecordExist('object', 'topology', 'crecord_name', record, function(ctrl, record, exist) {
			if(exist) {
				ctrl._save(record, false);
			}
			else {
				global.notify.notify(_('Bad name'), _('This topology name already exist'), 'warning');
			}
		}, this);
	},

	afterload_EditForm: function(form, item_copy) {
		form.rightPanel.setValue({
			'nodes': item_copy.get('nodes'),
			'conns': item_copy.get('conns'),
			'root': item_copy.get('root')
		});
	}
});
