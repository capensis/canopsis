//need:app/lib/controller/cgrid.js,app/view/Consolidation/Grid.js,app/view/Consolidation/Form.js,app/store/Consolidations.js
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
Ext.define('canopsis.controller.Consolidation', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Consolidation.Grid', 'Consolidation.Form'],
	stores: ['Consolidations'],
	models: ['Consolidation'],

	logAuthor: '[controller][Consolidation]',

	init: function() {
		log.debug('Initialize ...', this.logAuthor);
		this.formXtype = 'ConsolidationForm';
		this.listXtype = 'ConsolidationGrid';

		this.modelId = 'Consolidation';
		this.callParent(arguments);
	},

	_saveForm: function(form) {
		if(form.record !== undefined) {
			form.record.loaded = false;
			form.record.nb_items = undefined;
			form.record.output_engine = undefined;
		}

		this.callParent(arguments);
	},

	afterload_EditForm: function(form, item_copy) {
		// checkboxgroup don't tick boxes, this code do.
		var operators = item_copy.get('consolidation_method');

		if(!Ext.isArray(operators)) {
			operators = [operators];
		}

		for(var i = 0; i < operators.length; i++) {
			form.down('checkbox[inputValue=' + operators[i] + ']').setValue(true);
		}
	},

	afterload_DuplicateForm: function(form, item_copy) {
		this.afterload_EditForm(form, item_copy);
	}

});
