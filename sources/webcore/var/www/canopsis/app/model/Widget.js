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
Ext.define('canopsis.model.Widget', {
	extend: 'Ext.data.Model',
	fields: [
		{name: 'name'},
		{name: 'description', defaultValue: undefined},
		{name: 'description-fr', defaultValue: undefined},
		{name: 'version', defaultValue: undefined},
		{name: 'author', defaultValue: undefined},
		{name: 'website', defaultValue: undefined},
		{name: 'options', defaultValue: undefined},
		{name: 'xtype'},
		{name: 'thumb'},
		{name: 'colspan', defaultValue: 1},
		{name: 'rowspan', defaultValue: 1},
		{name: 'refreshInterval', defaultValue: undefined},
		{name: 'nodeId', defaultValue: undefined},
		{name: 'title', defaultValue: undefined},
		{name: 'border', defaultValue: false},
		{name: 'rowHeight', defaultValue: undefined},
		{name: 'formWidth', defaultValue: 350},
		{name: 'locales', defaultValue: undefined},
		{name: 'disabled', defaultValue: false},
		{name: 'thirdparty', defaultValue: false}
	]
});
