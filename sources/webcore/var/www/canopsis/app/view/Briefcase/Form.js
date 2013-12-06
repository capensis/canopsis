//need:app/lib/view/cform.js
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
Ext.define('canopsis.view.Briefcase.Form', {
	extend: 'canopsis.lib.view.cform',

	alias: 'widget.BriefcaseForm',

	initComponent: function() {
		var _id = Ext.widget('textfield', {
			fieldLabel: _('File name'),
			name: '_id',
			hidden: true
		});

		var filename = Ext.widget('textfield', {
			fieldLabel: _('File name'),
			name: 'file_name'
		});

		this.items = [_id, filename];

		this.callParent();
	}
});