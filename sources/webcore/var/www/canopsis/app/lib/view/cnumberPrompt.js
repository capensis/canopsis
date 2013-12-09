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

//inspirated by http://stackoverflow.com/questions/13896474/numberfield-prompt-in-extjs

Ext.define('canopsis.lib.view.cnumberPrompt', {
	extend: 'Ext.window.MessageBox',

	MessageBoxMinVal: 0.1,
	defaultValue: 1,

	initComponent: function() {
		this.callParent();

		var index = this.promptContainer.items.indexOf(this.textField);
		this.promptContainer.remove(this.textField);
		this.textField = this._createNumberField();
		this.promptContainer.insert(index, this.textField);
	},

	_createNumberField: function() {
		//copy paste what is being done in the initComonent to create the textfield
		return Ext.widget('numberfield', {
			id: this.id + '-textfield',
			anchor: '100%',
			value:1,
			minValue: this.MessageBoxMinVal,
			enableKeyEvents: true,
			listeners: {
				keydown: this.onPromptKey,
				scope: this
			}
		});
	}
});