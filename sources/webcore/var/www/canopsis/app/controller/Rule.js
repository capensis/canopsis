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
Ext.define('canopsis.controller.Rule', {
	extend: 'canopsis.lib.controller.cgrid',

	views: ['Rule.Grid', 'Rule.Form'],
	stores: ['Rules', 'DefaultRules'],
	models: ['Rule', 'DefaultRule'],

	logAuthor: '[controller][rules]',

	storeDefaultAction: undefined,

	init: function() {
		log.debug('Initialize ...', this.logAuthor);

		this.formXtype = 'RuleForm';
		this.listXtype = 'RuleGrid';

		this.modelId = 'Rule';

		this.control( {
			'combobox[alias=widget.defaultAction]': {
				select: this.defaultRuleChanged,
				afterrender: this.updateDefaultActionCombo
			}
		});

		this.callParent(arguments);

		this.storeDefaultAction = this.getStore('DefaultRules');

		this.storeDefaultAction.load();
	},

	defaultRuleChanged: function(combo, newVal) {
		void(combo);

		var record = this.storeDefaultAction.getAt(0);

		if(record === undefined) {
			this.storeDefaultAction.add({action: 'pass'});
		}
		else {
			this.storeDefaultAction.getAt(0).set("action", newVal[0].data.value);
		}
	},

	updateDefaultActionCombo: function(combo) {
		var comboValue;

		try {
			comboValue = this.storeDefaultAction.first().get("action");
		}
		catch(err) {
			comboValue = "pass";
		}

		if(comboValue === "pass" || comboValue === undefined) {
			combo.select(combo.getStore().getAt(0));
		}
		else {
			combo.select(combo.getStore().getAt(1));
		}
	}
});